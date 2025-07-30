# ✅ 25 Modern Programming and Design Principles Interview Questions (with Detailed Answers)

## 1. What is SOLID? Explain each principle.

SOLID is an acronym for 5 object-oriented design principles:

* **S**: Single Responsibility Principle (SRP) – A class should have one, and only one, reason to change.
* **O**: Open/Closed Principle – Software entities should be open for extension but closed for modification.
* **L**: Liskov Substitution Principle – Subtypes must be substitutable for their base types without altering correctness.
* **I**: Interface Segregation Principle – Clients should not be forced to depend on interfaces they do not use.
* **D**: Dependency Inversion Principle – Depend on abstractions, not concretions.

## 2. What is the DRY principle?

**DRY** stands for **Don’t Repeat Yourself**. Every piece of knowledge must have a single, unambiguous, authoritative representation within a system. Avoid duplicating logic to enhance maintainability.

## 3. What is YAGNI and when should it be applied?

**YAGNI** = "You Aren’t Gonna Need It." Don’t implement a feature until it is actually required. It avoids overengineering.

## 4. What is the KISS principle in software design?

**Keep It Simple, Stupid** – Emphasizes simplicity. Avoid unnecessary complexity; the simplest solution that works is often the best.

## 5. What is the Law of Demeter (LoD)?

Also known as the **Principle of Least Knowledge**, it recommends a method should only interact with:

* Itself
* Objects it owns
* Parameters
  Avoid deep chains like `obj.getX().getY().doSomething()`.

## 6. What is cohesion and why is it important?

Cohesion refers to how closely related the responsibilities of a module/class are. **High cohesion** results in better readability, maintainability, and testability.

## 7. What is coupling?

Coupling measures how dependent modules are on each other. **Low (loose) coupling** is desired to allow components to evolve independently.

## 8. What is separation of concerns?

It advocates dividing a system into distinct sections, where each addresses a specific concern (e.g., UI, business logic, DB). Enhances modularity and clarity.

## 9. How do you ensure your code is testable?

* Use dependency injection
* Break down large functions
* Avoid global state
* Prefer interfaces
* Isolate side effects

## 10. What are design patterns? Name a few used in Go.

Design patterns are standard solutions to common design problems. Common Go patterns:

* Factory
* Strategy (via interfaces)
* Singleton (via sync.Once)
* Decorator (via embedding)
* Adapter
* Observer

## 11. What is immutability? How do you implement it in Go?

Immutability = state doesn't change after creation. In Go:

* Keep struct fields private
* Provide read-only accessors
* Avoid mutation of slices/maps unless safe

## 12. How would you handle dependency injection in Go?

* Use constructor injection
* Define interfaces for dependencies
* Avoid global singletons
* Can use tools like `google/wire`

## 13. What is the difference between composition and inheritance?

* **Composition** (Go's style): Reuse behavior by combining types
* **Inheritance**: Behavior is derived from base classes (not in Go)
  Go promotes **composition over inheritance**.

## 14. How do you apply the Open-Closed Principle in Go?

Design code to use interfaces and inject behavior. New features can be added via new implementations without modifying core logic.

## 15. What is defensive programming?

Defensive programming anticipates and handles unexpected input or usage:

* Input validation
* Error logging
* Panic recovery

## 16. What is code smell? Give examples.

Code smells are signs of poor design:

* God object (large class)
* Long methods
* Deep nesting
* Repeated logic

## 17. What is the benefit of small functions?

* Easier to test
* Reusable
* Improve readability
* Focused logic (SRP)

## 18. Why is interface segregation important in Go?

Smaller interfaces are more flexible:

```go
type Reader interface {
   Read([]byte) (int, error)
}
```

Clients only depend on what they use.

## 19. What is the benefit of using context in Go?

`context.Context` enables:

* Timeouts
* Cancellations
* Propagation across API boundaries
  Used for graceful shutdown and request-scoped values.

## 20. What are the benefits of clean code practices?

* Readability
* Maintainability
* Onboarding new devs
* Fewer bugs
* Easier testing

## 21. What is idempotency in the context of API design?

Calling an idempotent API multiple times results in the same outcome. Helps avoid duplication (e.g., multiple refund triggers).

## 22. What are some API design principles you follow?

* RESTful URL structure
* Use verbs for actions (POST, GET, DELETE)
* Support versioning
* Return appropriate HTTP codes
* Document with OpenAPI/Swagger

## 23. What is the principle of least privilege?

Grant only necessary permissions. For example, limit DB user permissions to prevent accidental DELETEs or schema changes.

## 24. What is optimistic vs pessimistic locking?

* **Optimistic**: No locking. Conflict checked via version/timestamp.
* **Pessimistic**: Locks the record to prevent conflict. Safer but less scalable.

## 25. How do you handle backwards compatibility?

* Never break existing public APIs
* Use API versioning
* Use feature flags for gradual rollout
* Maintain old DB columns alongside new ones

---