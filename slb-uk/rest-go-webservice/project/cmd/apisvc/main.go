// main.go
// @title Message Service API
// @version 1.0
// @description This is the API server for the distributed message service.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1

//go:generate swag init --parseDependency --parseInternal --dir . --output docs
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type messageBody struct {
	Message string `json:"message"`
}

type acceptedResp struct {
	TraceID string `json:"trace_id"`
	Status  string `json:"status"`
}

type Ack struct {
	TraceID string                 `json:"trace_id"`
	Status  string                 `json:"status"`
	Event   string                 `json:"event"`
	Payload map[string]any         `json:"payload,omitempty"`
	Error   *struct{ Code, Detail string } `json:"error,omitempty"`
}

var (
	resultCache   = make(map[string]Ack)
	resultExpires = make(map[string]time.Time)
	cacheMu       sync.RWMutex
	cacheTTL      = 2 * time.Minute
)

func putAck(a Ack) {
	cacheMu.Lock()
	resultCache[a.TraceID] = a
	resultExpires[a.TraceID] = time.Now().Add(cacheTTL)
	cacheMu.Unlock()
}

func getAck(id string) (Ack, bool) {
	cacheMu.RLock()
	a, ok := resultCache[id]
	exp := resultExpires[id]
	cacheMu.RUnlock()
	if !ok || time.Now().After(exp) {
		return Ack{}, false
	}
	return a, true
}

func sweeper() {
	for range time.Tick(30 * time.Second) {
		cacheMu.Lock()
		for k, t := range resultExpires {
			if time.Now().After(t) {
				delete(resultExpires, k)
				delete(resultCache, k)
			}
		}
		cacheMu.Unlock()
	}
}

// @Summary Create a new message
// @Description Receives a message payload and publishes to Kafka
// @Tags messages
// @Accept json
// @Produce json
// @Param message body messageBody true "Message payload"
// @Success 200 {object} acceptedResp
// @Failure 400 {string} string "invalid body"
// @Router /messages [post]
func createMessageHandler(producer sarama.SyncProducer, cmdTopic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var b messageBody
		if json.NewDecoder(r.Body).Decode(&b) != nil || strings.TrimSpace(b.Message) == "" {
			http.Error(w, "invalid body", 400)
			return
		}
		enqueueCommand(w, producer, cmdTopic, "Create", map[string]any{"message": b.Message})
	}
}

// @Summary Get a message by ID
// @Tags messages
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} Ack
// @Router /messages/{id} [get]
// @Summary Update a message
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param message body messageBody true "Updated message"
// @Success 200 {object} Ack
// @Router /messages/{id} [put]
// @Summary Delete a message
// @Tags messages
// @Param id path string true "Message ID"
// @Success 204
// @Router /messages/{id} [delete]
func messageByIDHandler(producer sarama.SyncProducer, cmdTopic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/v1/messages/")
		switch r.Method {
		case http.MethodGet:
			enqueueCommand(w, producer, cmdTopic, "Read", map[string]any{"id": idStr})
		case http.MethodPut:
			var b messageBody
			if json.NewDecoder(r.Body).Decode(&b) != nil || strings.TrimSpace(b.Message) == "" {
				http.Error(w, "invalid body", 400)
				return
			}
			enqueueCommand(w, producer, cmdTopic, "Update", map[string]any{"id": idStr, "message": b.Message})
		case http.MethodDelete:
			enqueueCommand(w, producer, cmdTopic, "Delete", map[string]any{"id": idStr})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// @Summary Get operation status
// @Tags operations
// @Produce json
// @Param trace_id path string true "Trace ID"
// @Success 200 {object} Ack
// @Success 204 {string} string "No Content"
// @Router /operations/{trace_id} [get]
func operationResultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID := strings.TrimPrefix(r.URL.Path, "/v1/operations/")
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		for {
			if a, ok := getAck(traceID); ok {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(a)
				return
			}
			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusNoContent)
				return
			case <-time.After(200 * time.Millisecond):
			}
		}
	}
}

func enqueueCommand(w http.ResponseWriter, p sarama.SyncProducer, topic, cmd string, payload map[string]any) {
	traceID := uuid.NewString()
	idemp := uuid.NewString()
	m := map[string]any{
		"trace_id": traceID,
		"command":  cmd,
		"resource": "Message",
		"payload":  payload,
	}
	b, _ := json.Marshal(m)

	headers := []sarama.RecordHeader{
		{Key: []byte("trace_id"), Value: []byte(traceID)},
		{Key: []byte("command"), Value: []byte(cmd)},
	}

	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.ByteEncoder(idemp),
		Value:   sarama.ByteEncoder(b),
		Headers: headers,
	}

	if _, _, err := p.SendMessage(msg); err != nil {
		http.Error(w, "enqueue failed", 503)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(acceptedResp{TraceID: traceID, Status: "PENDING"})
}

func startAckConsumer(brokers []string, topic string) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Version = sarama.V2_6_0_0

	group, err := sarama.NewConsumerGroup(brokers, "api-acks", cfg)
	if err != nil {
		log.Fatal(err)
	}

	handler := &ackHandler{}

	go func() {
		for {
			if err := group.Consume(context.Background(), []string{topic}, handler); err != nil {
				log.Println("ack consume error:", err)
				time.Sleep(time.Second)
			}
		}
	}()
}

type ackHandler struct{}

func (ackHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (ackHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (ackHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var a Ack
		if err := json.Unmarshal(msg.Value, &a); err == nil && a.TraceID != "" {
			putAck(a)
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func main() {
	brokers := strings.Split(getenv("KAFKA_BROKERS", "kafka:9092"), ",")
	cmdTopic := getenv("KAFKA_TOPIC_COMMANDS", "messages.commands")
	acksTopic := getenv("KAFKA_TOPIC_ACKS", "messages.acks")
	addr := getenv("API_HTTP_ADDR", ":8080")

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Producer.Return.Successes = true
	cfg.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	go startAckConsumer(brokers, acksTopic)
	go sweeper()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/messages", createMessageHandler(producer, cmdTopic))
	mux.HandleFunc("/v1/messages/", messageByIDHandler(producer, cmdTopic))
	mux.HandleFunc("/v1/operations/", operationResultHandler())

	log.Println("API listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
