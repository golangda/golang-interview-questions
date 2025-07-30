# 🧠 Modern Kafka Architecture (KRaft Mode)

Apache Kafka is a high-throughput, distributed event streaming platform. Starting with Kafka 2.8+, a new **KRaft mode** was introduced to eliminate the dependency on ZooKeeper and simplify the deployment architecture.

---

## 📊 Kafka Modern Architecture (KRaft Mode)

### ✅ Key Components

| Component                 | Description                                                                        |
| ------------------------- | ---------------------------------------------------------------------------------- |
| **Producer**              | Sends records (messages) to Kafka topics.                                          |
| **Topic**                 | A logical channel to which producers write and consumers read.                     |
| **Partition**             | Subdivision of a topic for parallelism and scalability.                            |
| **Broker**                | A Kafka server that stores data and serves clients.                                |
| **Kafka Controller**      | Elected broker that manages metadata and cluster state (in KRaft).                 |
| **Consumer**              | Subscribes to one or more topics and processes records.                            |
| **Consumer Group**        | Multiple consumers working together for scalability.                               |
| **KRaft Metadata Quorum** | Replaces ZooKeeper, a quorum of brokers that manage metadata using Raft consensus. |
| **Log**                   | Each partition is an append-only log file on disk.                                 |
| **Offset**                | Position of a record within a partition.                                           |

---

## 🖼️ Kafka KRaft Mode Architecture Diagram

```
                +-----------------------------+
                |         Producers           |
                +-----------------------------+
                           |
                           v
                 +---------------------+
                 |     Kafka Topics    |
                 |   (with Partitions) |
                 +---------------------+
                   /       |       \        
                  /        |        \      
                 v         v         v
         +---------+ +-----------+ +-----------+
         | Broker1 | |  Broker2  | |  Broker3  |   <-- Kafka Cluster
         | (Leader)| | (Follower)| | (Leader) |   <-- Store partitions
         +---------+ +-----------+ +-----------+
              |           |             |
              |           |             |
         +-----------------------------------+
         |  KRaft Metadata Quorum (3 Brokers)|
         |   (Raft-based consensus)         |
         +-----------------------------------+
                        |
                        v
            +-----------------------------+
            |      Kafka Controller       |
            | (part of one broker node)   |
            +-----------------------------+
                        |
                        v
             +------------------------+
             |     Consumers          |
             | (Consumer Groups)      |
             +------------------------+
```

---

## 🧩 Kafka KRaft vs. ZooKeeper Mode

| Feature                  | Kafka with ZooKeeper | Kafka KRaft Mode          |
| ------------------------ | -------------------- | ------------------------- |
| Metadata Manager         | ZooKeeper            | Kafka Broker (Controller) |
| Consensus Algorithm      | ZAB (ZooKeeper)      | Raft                      |
| Setup Complexity         | Requires ZooKeeper   | Simpler (No ZK needed)    |
| Production-ready (2024+) | Deprecated           | Recommended               |

---

## 🔄 Data Flow Summary (in KRaft Mode)

1. **Producer** sends a message to a **topic**.
2. The message is written to a specific **partition** on a **broker**.
3. Each partition has one **leader broker** and zero or more **followers**.
4. Data is stored on disk in append-only **log files**.
5. **Consumers** read data from partitions, tracking offsets.
6. The **Kafka Controller**, elected from brokers, manages metadata and partition leadership.
7. **KRaft quorum** maintains metadata consistency using Raft consensus protocol.

---

## 🔐 Modern Features & Enhancements

* **Exactly-once semantics (EOS)** via idempotent producers and transactional APIs.
* **Schema Registry** integration for Avro/JSON/Protobuf validation.
* **Kafka Streams** for in-stream processing.
* **Kafka Connect** for integration with external systems (e.g., DBs, S3).
* **Tiered Storage (in preview)** – store old data on cheaper storage (e.g., cloud/S3).

---

## 📘 Useful References

* [Kafka Documentation](https://kafka.apache.org/documentation/)
* [KRaft Mode Intro](https://kafka.apache.org/documentation/#kraft)
* [Kafka Architecture Explained (Confluent)](https://www.confluent.io/blog/kafka-fastest-messaging-system/)

---

## 🧪 Try This Locally

You can simulate a local KRaft-based Kafka cluster using:

```bash
docker-compose -f https://github.com/confluentinc/cp-all-in-one/blob/7.5.0-post/cp-all-in-one-kraft/docker-compose.yml up
```
