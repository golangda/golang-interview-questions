# Top 10 Design Patterns in Go

This document covers the **Top 10 Design Patterns in Go** with definitions, sample code, and explanations. Useful for system design and Go interviews such as BigCommerce.

---

## 1. Singleton Pattern

**Definition**: Ensures a class has only one instance and provides a global point of access to it.

```go
package singleton

import (
	"sync"
)

type Config struct {
	AppName string
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{AppName: "BigCommerce Returns"}
	})
	return instance
}
```

**Explanation**: `sync.Once` ensures thread-safe lazy initialization.

---

## 2. Factory Pattern

**Definition**: Provides an interface for creating objects in a superclass, but allows subclasses to alter the type of objects that will be created.

```go
package factory

type ReturnType interface {
	Process() string
}

type Refund struct{}
func (r Refund) Process() string { return "Refund Processed" }

type Exchange struct{}
func (e Exchange) Process() string { return "Exchange Processed" }

func GetReturnType(t string) ReturnType {
	if t == "refund" {
		return Refund{}
	} else if t == "exchange" {
		return Exchange{}
	}
	return nil
}
```

**Explanation**: Abstract object creation with decoupled logic.

---

## 3. Builder Pattern

**Definition**: Builds complex objects step by step using a builder struct.

```go
package builder

type ReturnRequest struct {
	OrderID string
	Reason  string
	Items   []string
}

type ReturnBuilder struct {
	req ReturnRequest
}

func NewBuilder() *ReturnBuilder {
	return &ReturnBuilder{req: ReturnRequest{}}
}

func (b *ReturnBuilder) SetOrderID(id string) *ReturnBuilder {
	b.req.OrderID = id
	return b
}

func (b *ReturnBuilder) AddItem(item string) *ReturnBuilder {
	b.req.Items = append(b.req.Items, item)
	return b
}

func (b *ReturnBuilder) Build() ReturnRequest {
	return b.req
}
```

**Explanation**: Helps in constructing readable and maintainable objects.

---

## 4. Strategy Pattern

**Definition**: Defines a family of algorithms and makes them interchangeable.

```go
package strategy

type RefundStrategy interface {
	Calculate(amount float64) float64
}

type FullRefund struct{}
func (f FullRefund) Calculate(amount float64) float64 { return amount }

type PartialRefund struct{}
func (p PartialRefund) Calculate(amount float64) float64 { return amount * 0.5 }

func ProcessRefund(s RefundStrategy, amt float64) float64 {
	return s.Calculate(amt)
}
```

**Explanation**: Enables runtime strategy switching.

---

## 5. Adapter Pattern

**Definition**: Converts the interface of a class into another interface clients expect.

```go
package adapter

type LegacyAPI struct{}
func (l LegacyAPI) OldReturnMethod() string { return "Legacy Return Processed" }

type NewAPI interface {
	NewReturnMethod() string
}

type Adapter struct {
	legacy LegacyAPI
}

func (a Adapter) NewReturnMethod() string {
	return a.legacy.OldReturnMethod()
}
```

**Explanation**: Useful when integrating with legacy systems.

---

## 6. Observer Pattern

**Definition**: A one-to-many dependency between objects so when one object changes state, all its dependents are notified.

```go
package observer

type Subscriber interface {
	Notify(string)
}

type Publisher struct {
	subscribers []Subscriber
}

func (p *Publisher) Register(s Subscriber) {
	p.subscribers = append(p.subscribers, s)
}

func (p *Publisher) Broadcast(msg string) {
	for _, s := range p.subscribers {
		s.Notify(msg)
	}
}
```

**Explanation**: Perfect for event-driven systems.

---

## 7. Decorator Pattern

**Definition**: Adds behavior to objects dynamically without changing the original structure.

```go
package decorator

import "fmt"

type Processor interface {
	Process() string
}

type BasicReturn struct{}
func (b BasicReturn) Process() string { return "Return Initiated" }

type WithLogging struct {
	Processor
}

func (w WithLogging) Process() string {
	fmt.Println("Logging return...")
	return w.Processor.Process()
}
```

**Explanation**: Wrap base logic for added responsibilities.

---

## 8. Command Pattern

**Definition**: Encapsulates a request as an object.

```go
package command

import "fmt"

type Command interface {
	Execute()
}

type RefundCommand struct{}
func (r RefundCommand) Execute() { fmt.Println("Processing Refund") }

type Invoker struct {
	commands []Command
}

func (i *Invoker) AddCommand(c Command) {
	i.commands = append(i.commands, c)
}

func (i *Invoker) ExecuteAll() {
	for _, c := range i.commands {
		c.Execute()
	}
}
```

**Explanation**: Useful for deferred execution or audit logging.

---

## 9. Prototype Pattern

**Definition**: Create new objects by copying an existing object.

```go
package prototype

type ReturnRequest struct {
	OrderID string
	Reason  string
	Items   []string
}

func (r *ReturnRequest) Clone() *ReturnRequest {
	items := make([]string, len(r.Items))
	copy(items, r.Items)
	return &ReturnRequest{
		OrderID: r.OrderID,
		Reason:  r.Reason,
		Items:   items,
	}
}
```

**Explanation**: Useful when object creation is expensive.

---

## 10. Chain of Responsibility

**Definition**: Passes requests along a chain of handlers.

```go
package chain

import "fmt"

type Handler interface {
	SetNext(Handler)
	Handle(string)
}

type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(n Handler) {
	b.next = n
}

func (b *BaseHandler) Handle(r string) {
	if b.next != nil {
		b.next.Handle(r)
	}
}

type AuthHandler struct{ BaseHandler }
func (a *AuthHandler) Handle(r string) {
	fmt.Println("Auth passed")
	a.BaseHandler.Handle(r)
}

type ValidateHandler struct{ BaseHandler }
func (v *ValidateHandler) Handle(r string) {
	fmt.Println("Validation passed")
	v.BaseHandler.Handle(r)
}
```

**Explanation**: Helps structure middleware or pre-processing chains.

---

## Summary Table

| Pattern                  | Purpose                                           | Use Case Example                  |
|--------------------------|---------------------------------------------------|-----------------------------------|
| Singleton                | One global instance                               | Config, logger                    |
| Factory                  | Abstract object creation                          | Creating return types             |
| Builder                  | Construct object step-by-step                     | ReturnRequest construction        |
| Strategy                 | Runtime algorithm switching                       | Refund calculation                |
| Adapter                  | Convert legacy interfaces                         | API adapter                       |
| Observer                 | Event notification                                | Notify on return state change     |
| Decorator                | Add responsibilities dynamically                  | Logging/metrics wrappers          |
| Command                  | Encapsulate execution logic                       | Refund queue, audit trail         |
| Prototype                | Copy object without binding to class              | Duplication of return template    |
| Chain of Responsibility  | Request pipeline with handlers                    | Middleware in API servers         |
