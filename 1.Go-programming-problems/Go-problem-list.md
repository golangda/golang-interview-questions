# 🧠 50 Apache Kafka Interview Questions – With Detailed Answers

This document contains 50 carefully selected Kafka interview questions with answers, grouped across foundational, architectural, programming, performance, and real-world use case topics.

---

## ✅ Section 1: Kafka Basics

### 1. ❓ What is Apache Kafka?
**Answer**: Kafka is a distributed, partitioned, replicated commit log service. It is used for building real-time streaming data pipelines and applications. Kafka was originally developed at LinkedIn and is now part of the Apache Software Foundation.

### 2. ❓ What are the main components of Kafka?
**Answer**:
- **Producer**: Sends data to Kafka topics.
- **Consumer**: Reads data from Kafka topics.
- **Broker**: Kafka server that stores data and serves clients.
- **Topic**: A category/feed name to which records are sent.
- **Partition**: Topic split for scalability and parallelism.
- **ZooKeeper** (optional in older versions): Used for managing cluster state.

### 3. ❓ What is a Kafka Topic?
**Answer**: A topic is a logical channel to which data records are published. Consumers subscribe to one or more topics to receive data.

### 4. ❓ What is a Kafka Partition?
**Answer**: Each Kafka topic is split into partitions. Partitions allow Kafka to scale horizontally and achieve parallelism.

### 5. ❓ What is an offset in Kafka?
**Answer**: An offset is a unique identifier (long integer) for each message within a partition. Consumers use offsets to track read progress.

### 6. ❓ How does Kafka ensure message durability?
**Answer**: Kafka persists messages to disk immediately upon arrival and uses replication across brokers to ensure durability.

### 7. ❓ What is the difference between Kafka and traditional message brokers (e.g., RabbitMQ)?
**Answer**:
- Kafka is log-based and supports high throughput.
- Kafka retains messages for a configured duration regardless of consumption.
- Kafka is distributed and horizontally scalable by design.

### 8. ❓ What is the retention period in Kafka?
**Answer**: It's the configured duration for which Kafka retains messages (e.g., 7 days). Controlled via `log.retention.hours` or `log.retention.bytes`.

### 9. ❓ What is the difference between Consumer Group and Consumer?
**Answer**:
- **Consumer**: Reads data from topic partitions.
- **Consumer Group**: A set of consumers coordinating to read partitions of a topic without overlap.

### 10. ❓ Can a partition be read by multiple consumers in the same group?
**Answer**: No. Only one consumer per consumer group reads from a given partition. This avoids duplicate reads.

---

## ✅ Section 2: Producers & Consumers

### 11. ❓ What happens when a Kafka producer sends messages to a topic that doesn’t exist?
**Answer**: If auto topic creation is enabled on the broker, Kafka will create the topic with default settings. Otherwise, the producer receives an error.

### 12. ❓ How does Kafka achieve high throughput on the producer side?
**Answer**:
- Batching messages
- Compression (e.g., Snappy, GZIP)
- Asynchronous sending
- Parallelism with multiple partitions

### 13. ❓ What is `acks=all` in Kafka producer configuration?
**Answer**: It means the leader waits for acknowledgment from all in-sync replicas before considering the write successful. Ensures highest durability.

### 14. ❓ What is the role of a Kafka Serializer?
**Answer**: Serializers convert objects into byte arrays before sending data to Kafka. E.g., `StringSerializer`, `JsonSerializer`.

### 15. ❓ How do Kafka consumers keep track of what messages they’ve read?
**Answer**: Kafka consumers track offsets. These can be committed automatically or manually to Kafka’s internal `__consumer_offsets` topic.

### 16. ❓ What is the difference between earliest and latest offset reset?
**Answer**:
- `earliest`: Start reading from the beginning of the partition.
- `latest`: Start from the latest record (new data only).

### 17. ❓ What happens if a Kafka consumer fails after consuming but before committing the offset?
**Answer**: The message will be reprocessed when the consumer restarts, as the offset wasn’t committed.

### 18. ❓ How does Kafka handle backpressure in consumers?
**Answer**: Kafka itself doesn’t enforce backpressure. It’s up to the consumer app to control flow using batching, pause/resume, and proper offset commits.

### 19. ❓ What is Kafka’s exactly-once semantics (EOS)?
**Answer**: Kafka supports EOS using:
- Idempotent producers
- Transactions (atomic writes to multiple partitions + offsets)
- Enabled with `enable.idempotence=true` and proper transaction setup

### 20. ❓ What is Kafka Connect?
**Answer**: A framework for ingesting data from sources (DBs, APIs) into Kafka and exporting from Kafka to sinks (DBs, files, systems). Uses prebuilt connectors.

---

## ✅ Section 3: Kafka Internals

### 21. ❓ What is a Kafka Broker?
**Answer**: A broker is a Kafka server that stores data and serves producer/consumer requests. A cluster usually has multiple brokers.

### 22. ❓ How does Kafka replicate data?
**Answer**: Each partition has replicas across brokers. One is elected as leader; others are followers. Replication ensures durability and high availability.

### 23. ❓ What is ISR (In-Sync Replica)?
**Answer**: ISR is the set of replicas that are fully caught up with the leader. Kafka writes are only acknowledged if replicated to all ISRs.

### 24. ❓ What is the role of ZooKeeper in Kafka (pre-2.8)?
**Answer**: ZooKeeper manages broker metadata, leader election, and cluster configuration. From Kafka 2.8+, ZooKeeper is optional due to KRaft mode.

### 25. ❓ What is the role of the controller in Kafka?
**Answer**: The controller is a broker elected to manage partition leaders, replication status, and administrative events.

### 26. ❓ What are log segments?
**Answer**: Kafka partitions are stored as log segments on disk. Each segment is a file with a range of offsets.

### 27. ❓ What is log compaction?
**Answer**: A mechanism to retain only the latest value for a given key. Useful for changelog-type topics.

### 28. ❓ What is the `__consumer_offsets` topic?
**Answer**: An internal Kafka topic that stores committed offsets for consumer groups.

### 29. ❓ What is idempotence in Kafka producers?
**Answer**: Guarantees that retries do not result in duplicate messages. Enabled by `enable.idempotence=true`.

### 30. ❓ How does Kafka handle leader election?
**Answer**: ZooKeeper or the KRaft controller coordinates leader election for partitions. The controller picks one ISR as the new leader when the current one fails.

---

## ✅ Section 4: Performance & Tuning

### 31. ❓ How to tune Kafka producer for low latency?
**Answer**:
- `linger.ms=0`
- `batch.size` small
- `acks=1`
- Compression for network efficiency

### 32. ❓ How to tune Kafka for throughput?
**Answer**:
- Increase partition count
- Use asynchronous and batched producers
- Use Snappy compression
- Tune socket buffers and `fetch.min.bytes`

### 33. ❓ How many partitions should a topic have?
**Answer**: Depends on desired parallelism and throughput. Generally 2–4× the number of consumers or CPU cores.

### 34. ❓ What happens when a Kafka topic has too many partitions?
**Answer**:
- Increases overhead on brokers and controllers
- Can impact GC and replication latency
- Degrades performance if not managed

### 35. ❓ How does Kafka handle disk I/O efficiently?
**Answer**:
- Sequential writes to disk
- OS page cache utilization
- Batching and compression reduce IOPS

### 36. ❓ What is Kafka’s typical latency?
**Answer**: Low-latency pipelines can achieve 2–10ms end-to-end. Actual latency depends on broker load, batching, network, and consumer processing.

### 37. ❓ Can Kafka guarantee message order?
**Answer**: Yes, but only within a single partition. Messages across partitions may arrive out of order.

### 38. ❓ Can Kafka lose messages?
**Answer**: Not if configured properly (replication, acks=all, durability settings). Misconfigurations or hardware failures without replication can cause loss.

### 39. ❓ What is a dead letter queue (DLQ) in Kafka?
**Answer**: A topic where failed or malformed messages are redirected for later inspection or processing. Helps prevent consumer crashes.

### 40. ❓ What are Kafka Streams and how do they differ from Kafka Consumer API?
**Answer**:
- Kafka Streams is a Java library for processing Kafka data in real-time.
- Offers DSL, windowing, joins, aggregations, and fault tolerance.
- Kafka Consumer API is lower-level and requires manual state tracking.

---

## ✅ Section 5: Real-World Scenarios & Admin

### 41. ❓ How do you monitor Kafka health?
**Answer**:
- Monitor broker metrics (CPU, disk, network)
- Check topic lag using consumer group metrics
- Monitor controller status, under-replicated partitions

### 42. ❓ What tools are used for Kafka monitoring?
**Answer**:
- Prometheus + Grafana
- Confluent Control Center
- LinkedIn Burrow
- Kafka Manager

### 43. ❓ What is Kafka KRaft mode?
**Answer**: KRaft (Kafka Raft Metadata mode) replaces ZooKeeper with an internal quorum-based controller. Available from Kafka 2.8+.

### 44. ❓ How do you reassign partitions in Kafka?
**Answer**:
- Use `kafka-reassign-partitions.sh`
- Provide a reassignment JSON file
- Validate and execute using CLI

### 45. ❓ How do you delete a Kafka topic?
**Answer**:
```bash
kafka-topics.sh --bootstrap-server <host> --delete --topic <topic-name>
```
Make sure `delete.topic.enable=true` is set in broker config.

### 46. ❓ What happens if a broker fails?
**Answer**: Kafka reassigns leadership to another in-sync replica. Producers/consumers continue through other brokers.

### 47. ❓ Can Kafka be used as a database?
**Answer**: Not recommended for transactional workloads. Kafka is append-only and best for logs/event streams. Use with databases or state stores.

### 48. ❓ How do you secure a Kafka cluster?
**Answer**:
- Enable TLS for encryption
- SASL (Kerberos, SCRAM) for authentication
- ACLs for authorization

### 49. ❓ How do you ensure schema consistency in Kafka?
**Answer**: Use a Schema Registry (e.g., Confluent Schema Registry) with Avro or Protobuf. It validates compatibility during schema evolution.

### 50. ❓ What are common Kafka use cases?
**Answer**:
- Log aggregation
- Event sourcing
- Metrics collection
- Real-time analytics
- Microservice communication
- ETL and CDC pipelines

---

For hands-on practice:
- [Kafka Playground](https://www.katacoda.com/courses/kafka)
- [Confluent Cloud](https://www.confluent.io)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
