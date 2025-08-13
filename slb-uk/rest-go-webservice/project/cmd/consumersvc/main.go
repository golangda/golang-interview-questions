package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	_ "github.com/go-sql-driver/mysql"
)

type Command struct {
	TraceID  string                 `json:"trace_id"`
	Command  string                 `json:"command"`
	Resource string                 `json:"resource"`
	Payload  map[string]any         `json:"payload"`
}

type Ack struct {
	TraceID string                 `json:"trace_id"`
	Status  string                 `json:"status"`
	Event   string                 `json:"event"`
	Payload map[string]any         `json:"payload,omitempty"`
	Error   *struct{ Code, Detail string } `json:"error,omitempty"`
}

func main() {
	brokers := []string{getenv("KAFKA_BROKERS", "kafka:9092")}
	cmdTopic := getenv("KAFKA_TOPIC_COMMANDS", "messages.commands")
	acksTopic := getenv("KAFKA_TOPIC_ACKS", "messages.acks")
	dsn := getenv("MYSQL_DSN", "root:root@tcp(mysql:3306)/app?parseTime=true")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("db ping:", err)
	}

	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Version = sarama.V2_6_0_0
	cfg.Consumer.Return.Errors = true
	cfg.Producer.Return.Successes = true
	cfg.Producer.Idempotent = true

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "message-worker", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer consumerGroup.Close()

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	handler := &consumerHandler{db: db, producer: producer, ackTopic: acksTopic}

	log.Println("consumer running…")
	for {
		if err := consumerGroup.Consume(nil, []string{cmdTopic}, handler); err != nil {
			log.Println("consume error:", err)
			time.Sleep(time.Second)
		}
	}
}

type consumerHandler struct {
	db       *sql.DB
	producer sarama.SyncProducer
	ackTopic string
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var cmd Command
		if err := json.Unmarshal(msg.Value, &cmd); err != nil {
			log.Println("bad command:", err)
			continue
		}

		status := "SUCCESS"
		event := ""
		payload := map[string]any{}
		var e *struct{ Code, Detail string }

		err := withTx(h.db, func(tx *sql.Tx) error {
			key := string(msg.Key)
			if key == "" {
				key = cmd.TraceID
			}
			processed, err := checkIdempotent(tx, key)
			if err != nil {
				return err
			}
			if processed {
				return nil
			}

			switch cmd.Command {
			case "Create":
				m, _ := cmd.Payload["message"].(string)
				res, err := tx.Exec("INSERT INTO messages(message) VALUES(?)", m)
				if err != nil {
					status = "FAILURE"
					e = &struct{ Code, Detail string }{"DB_ERROR", err.Error()}
					logSaga(tx, cmd.TraceID, "CreateMessage", "FAILURE", "DB_ERROR", err.Error())
					return nil
				}
				id, _ := res.LastInsertId()
				payload["id"] = id
				payload["message"] = m
				event = "MessageCreated"
				logSaga(tx, cmd.TraceID, "CreateMessage", "SUCCESS", "", "")
			case "Read":
				idStr, _ := cmd.Payload["id"].(string)
				id, _ := strconv.ParseInt(idStr, 10, 64)
				row := tx.QueryRow("SELECT id, message FROM messages WHERE id=?", id)
				var mid int64
				var m string
				if err := row.Scan(&mid, &m); err != nil {
					status = "FAILURE"
					e = &struct{ Code, Detail string }{"NOT_FOUND", fmt.Sprintf("id=%d", id)}
					logSaga(tx, cmd.TraceID, "ReadMessage", "FAILURE", "NOT_FOUND", e.Detail)
					return nil
				}
				payload["id"] = mid
				payload["message"] = m
				event = "MessageRead"
				logSaga(tx, cmd.TraceID, "ReadMessage", "SUCCESS", "", "")
			case "Update":
				idStr, _ := cmd.Payload["id"].(string)
				id, _ := strconv.ParseInt(idStr, 10, 64)
				m, _ := cmd.Payload["message"].(string)
				res, err := tx.Exec("UPDATE messages SET message=? WHERE id=?", m, id)
				if err != nil {
					status = "FAILURE"
					e = &struct{ Code, Detail string }{"DB_ERROR", err.Error()}
					logSaga(tx, cmd.TraceID, "UpdateMessage", "FAILURE", "DB_ERROR", err.Error())
					return nil
				}
				affected, _ := res.RowsAffected()
				if affected == 0 {
					status = "FAILURE"
					e = &struct{ Code, Detail string }{"NOT_FOUND", fmt.Sprintf("id=%d", id)}
					logSaga(tx, cmd.TraceID, "UpdateMessage", "FAILURE", "NOT_FOUND", e.Detail)
					return nil
				}
				payload["id"] = id
				payload["message"] = m
				event = "MessageUpdated"
				logSaga(tx, cmd.TraceID, "UpdateMessage", "SUCCESS", "", "")
			case "Delete":
				idStr, _ := cmd.Payload["id"].(string)
				id, _ := strconv.ParseInt(idStr, 10, 64)
				res, err := tx.Exec("DELETE FROM messages WHERE id=?", id)
				if err != nil {
					status = "FAILURE"
					e = &struct{ Code, Detail string }{"DB_ERROR", err.Error()}
					logSaga(tx, cmd.TraceID, "DeleteMessage", "FAILURE", "DB_ERROR", err.Error())
					return nil
				}
				affected, _ := res.RowsAffected()
				if affected == 0 {
					status = "FAILURE"
					e = &struct{ Code, Detail string }{"NOT_FOUND", fmt.Sprintf("id=%d", id)}
					logSaga(tx, cmd.TraceID, "DeleteMessage", "FAILURE", "NOT_FOUND", e.Detail)
					return nil
				}
				payload["id"] = id
				event = "MessageDeleted"
				logSaga(tx, cmd.TraceID, "DeleteMessage", "SUCCESS", "", "")
			default:
				status = "FAILURE"
				e = &struct{ Code, Detail string }{"UNSUPPORTED", "unknown command"}
			}

			return markIdempotent(tx, key, cmd.TraceID, status)
		})

		if err != nil {
			log.Println("tx error:", err)
			status = "FAILURE"
			event = "Error"
			e = &struct{ Code, Detail string }{"INTERNAL", err.Error()}
		}

		ack := Ack{TraceID: cmd.TraceID, Status: status, Event: event, Payload: payload, Error: e}
		b, _ := json.Marshal(ack)
		ackMsg := &sarama.ProducerMessage{
		    Topic: h.ackTopic,
			Key:   sarama.ByteEncoder(msg.Key), // still using the consumer msg's key
			Value: sarama.ByteEncoder(b),
		}

		if _, _, err := h.producer.SendMessage(ackMsg); err != nil {
			log.Println("ack produce:", err)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func withTx(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func checkIdempotent(tx *sql.Tx, key string) (bool, error) {
	row := tx.QueryRow("SELECT 1 FROM idempotency_keys WHERE idempotency_key=?", key)
	var one int
	if err := row.Scan(&one); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func markIdempotent(tx *sql.Tx, key, traceID, status string) error {
	_, err := tx.Exec("INSERT IGNORE INTO idempotency_keys(idempotency_key, last_status, trace_id) VALUES(?,?,?)", key, status, traceID)
	return err
}

func logSaga(tx *sql.Tx, traceID, step, status, code, detail string) {
	_, _ = tx.Exec("INSERT INTO saga_log(trace_id, step, status, error_code, error_detail) VALUES(?,?,?,?,?)", traceID, step, status, code, detail)
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
