
# Kafka Interview Questions with Hints and Links

## 🧩 Section 1: Kafka Basics (Q1–Q8)
1. **What is Kafka and how is it different from a traditional message queue?**  
   🔹 *Hint*: Think distributed, log-based, pull model.  
   🔗 https://kafka.apache.org/documentation/

2. **Define the terms: topic, partition, offset, and consumer group.**  
   🔹 *Hint*: Partitions enable scalability; consumer groups enable parallelism.

3. **What guarantees does Kafka provide regarding message delivery?**  
   🔹 *Hint*: At most once, at least once, exactly once (EOS).  
   🔗 https://kafka.apache.org/documentation/#semantics

4. **What are Kafka brokers, producers, and consumers?**  
   🔹 *Hint*: Brokers store; producers write; consumers read.

5. **What is the role of ZooKeeper in Kafka (pre-KRaft)?**  
   🔹 *Hint*: Metadata management, controller election.  
   🔗 https://www.confluent.io/blog/kafka-kraft/

6. **How does Kafka achieve high throughput and low latency?**  
   🔹 *Hint*: Batching, zero-copy, compression.

7. **Explain retention vs compaction in Kafka.**  
   🔹 *Hint*: Retention is time/size-based; compaction keeps latest by key.  
   🔗 https://kafka.apache.org/documentation/#compaction

8. **What is the role of an ISR (In-Sync Replica)?**  
   🔹 *Hint*: Replicas that are up-to-date and eligible for leader election.

## 🧱 Section 2: Kafka Architecture (Q9–Q18)
9. **Describe Kafka’s partition leader election process.**  
   🔹 *Hint*: Controller broker decides.

10. **What happens when a Kafka broker goes down?**  
    🔹 *Hint*: Leader re-election for partitions hosted on it.

11. **How does Kafka ensure message order?**  
    🔹 *Hint*: Order is guaranteed within a partition.

12. **What’s the difference between Kafka replication and acknowledgment?**  
    🔹 *Hint*: Replication is for durability; acks control delivery guarantee.

13. **How does Kafka handle backpressure or slow consumers?**  
    🔹 *Hint*: Offsets and configurable retention.

14. **Explain the concept of log segments.**  
    🔹 *Hint*: Segments = immutable files written sequentially.

15. **What is the role of Kafka Controller?**  
    🔹 *Hint*: Manages partition leaders, detects broker failures.

16. **Can Kafka partitions be rebalanced? When and how?**  
    🔹 *Hint*: Yes—when brokers are added/removed.  
    🔗 https://docs.confluent.io/platform/current/kafka/rebalance.html

17. **What are rack awareness and broker affinity in Kafka?**  
    🔹 *Hint*: Useful for replica placement in multi-AZ/DC.

18. **What happens when a producer sends data to a non-leader partition replica?**  
    🔹 *Hint*: Fails unless sent to leader.

## 🧪 Section 3: Kafka Producer (Q19–Q25)
19. **Explain how batching works in Kafka producer.**  
    🔹 *Hint*: `linger.ms`, `batch.size`

20. **What are acks=0, 1, all in producer config?**  
    🔹 *Hint*: Controls durability/latency tradeoff.  
    🔗 https://kafka.apache.org/documentation/#producerconfigs_acks

21. **What is idempotence in Kafka producer? How do you enable it?**  
    🔹 *Hint*: Prevents duplicate writes.  
    🔗 https://www.confluent.io/blog/enabling-exactly-once-kafka-streams/

22. **How do retries and timeouts work in Kafka producer?**  
    🔹 *Hint*: `retries`, `request.timeout.ms`, `delivery.timeout.ms`

23. **How does Kafka ensure partition assignment for messages?**  
    🔹 *Hint*: Key-based hashing or custom Partitioner.

24. **What is the maximum size of a Kafka message? How do you handle large ones?**  
    🔹 *Hint*: `max.message.bytes`, chunking.

25. **How does compression work in Kafka? Which codecs are supported?**  
    🔹 *Hint*: gzip, snappy, lz4, zstd.

## 🧲 Section 4: Kafka Consumer (Q26–Q35)
26. **How does a Kafka consumer keep track of offsets?**  
    🔹 *Hint*: `enable.auto.commit` and consumer group coordination.

27. **What is rebalancing in Kafka consumer groups?**  
    🔹 *Hint*: Occurs when group membership changes.

28. **What are static membership and cooperative rebalancing?**  
    🔹 *Hint*: Reduce rebalancing downtime.  
    🔗 https://www.confluent.io/blog/cooperative-rebalancing-in-kafka/

29. **Explain partition assignment strategies.**  
    🔹 *Hint*: Range, RoundRobin, Sticky.

30. **How can you manually commit offsets in a consumer?**  
    🔹 *Hint*: `commitSync()`, `commitAsync()` in client.

31. **What’s the impact of consuming from a compacted topic?**  
    🔹 *Hint*: Only latest records by key.

32. **What happens if a consumer crashes before committing offsets?**  
    🔹 *Hint*: It reprocesses messages.

33. **Can two consumers in the same group read from the same partition?**  
    🔹 *Hint*: No—only one per partition per group.

34. **What is auto.offset.reset?**  
    🔹 *Hint*: earliest, latest, none.

35. **What are dead-letter topics and how are they used with Kafka?**  
    🔹 *Hint*: Failed messages redirected for later processing.

## 🧵 Section 5: Kafka CLI, Admin & APIs (Q36–Q42)
36. **Command to list topics from CLI?**  
    🔹 *Hint*: `kafka-topics.sh --list`

37. **Command to describe a topic’s partition details?**  
    🔹 *Hint*: `--describe` flag

38. **How do you create a Kafka topic with specific partitions and replication factor?**  
    🔹 *Hint*: CLI or Admin API

39. **How to delete a topic from CLI?**  
    🔹 *Hint*: `kafka-topics.sh --delete`

40. **What tools do you use to monitor Kafka health?**  
    🔹 *Hint*: JMX, Prometheus, Confluent Control Center

41. **How to check consumer group lags?**  
    🔹 *Hint*: `kafka-consumer-groups.sh --describe`

42. **What are some useful Kafka configs you’ve tuned in production?**  
    🔹 *Hint*: `linger.ms`, `num.partitions`, `replication.factor`, etc.

## 🔐 Section 6: Kafka Security & Transactions (Q43–Q46)
43. **What are the authentication mechanisms supported by Kafka?**  
    🔹 *Hint*: SASL, SSL, Kerberos

44. **How does Kafka authorize access to topics?**  
    🔹 *Hint*: ACLs via ZooKeeper or KRaft

45. **What are transactional messages in Kafka?**  
    🔹 *Hint*: Atomic writes across topics/partitions.

46. **How do you configure a producer for transactions?**  
    🔹 *Hint*: `transactional.id`, `enable.idempotence = true`

## 🚨 Section 7: Kafka in Production & Troubleshooting (Q47–Q50)
47. **What causes high consumer lag and how do you fix it?**  
    🔹 *Hint*: Slow consumers, network, partition skew.

48. **How do you prevent data loss in Kafka?**  
    🔹 *Hint*: `acks=all`, replication, idempotent producer.

49. **Kafka producer is timing out – what are possible causes?**  
    🔹 *Hint*: Broker down, large batch, unacknowledged writes.

50. **How do you do capacity planning for Kafka?**  
    🔹 *Hint*: Partition count, disk I/O, throughput, retention config.
