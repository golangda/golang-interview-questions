# BigCommerce Go Interview Recovery – Consolidated Interview Question Bank

This document contains a categorized list of interview questions derived directly from the feedback received in the BigCommerce interview. Use this as your Go-focused knowledge base for interviews and system design preparation.

---

## 🟨 Go Fundamentals

1. What is the difference between value and pointer receivers in Go? When would you use each?
2. How do slices, maps, and channels behave when passed between functions?
3. What is the zero value of common Go types like string, int, bool, and interfaces?
4. How do you manage shared state between goroutines safely?
5. What are common pitfalls with channel-based concurrency in Go?
6. How do you store and manage a dynamic list of channels or goroutines?
7. When using interfaces in Go, how do you ensure correctness with type assertions and type switches?
8. How do you structure a Go application with proper modular boundaries?
9. What is the idiomatic way to organize large Go codebases (beyond just controllers/services)?

---

## 🟨 Async / Event-Driven Systems

10. How would you correlate events across multiple services in a pub-sub system?
11. How would you implement end-to-end tracing in an async architecture?
12. How would you test a system where components communicate asynchronously?
13. What mechanisms would you use to prevent or detect message loss in Kafka/RMQ systems?
14. How do you gracefully shut down a long-running goroutine listening to events?
15. What retry mechanisms and patterns are suitable for event consumers?

---

## 🟨 Kafka vs RabbitMQ

16. Compare Kafka and RabbitMQ in terms of durability, scalability, latency, and message semantics.
17. Which would you choose for building an event-sourcing system and why?
18. How would you implement "at least once" delivery semantics in Kafka?
19. What are the pros/cons of using RabbitMQ’s routing keys vs Kafka’s topics/partitions?
20. How do acknowledgments and offset commits differ in Kafka and RabbitMQ?

---

## 🟨 Testing Strategy

21. What is the testing pyramid? How do you apply it in Go-based microservices?
22. What kinds of tests would you prioritize if you had limited time before a release?
23. How would you test a function that produces Kafka events and another that consumes them?
24. How do you mock dependencies (DB, Kafka, external APIs) in Go tests?
25. What tools or techniques do you use for writing integration vs unit vs E2E tests?
26. How do you ensure flaky tests in async systems are handled properly?
27. What’s your approach to testing time-sensitive logic (e.g., retries, timeouts)?

---

## 🟨 Monitoring & Observability

28. How would you monitor an async microservice that processes Kafka events?
29. What metrics would you track to identify performance issues in a Go service?
30. How would you set up alerting for failures in an event processing pipeline?
31. What’s the difference between logs, metrics, and traces — and when to use which?
32. How would you define SLIs/SLOs for a Go-based async service?
33. What are some examples of RED and USE metrics in a microservice setup?
34. How do you handle alert fatigue and false positives in production alerting?

---

## 🟨 Architecture & Design Thinking

35. What is Domain-Driven Design (DDD) and how would you apply it in Go?
36. How does Clean Architecture apply to Go microservices?
37. How would you design a scalable, observable, and fault-tolerant order return system?
38. How would you isolate domain logic from infrastructure in your Go codebase?
39. How would you refactor a legacy Go monolith into smaller services or modules?
40. What design principles would you apply when building APIs for public consumption?
---
