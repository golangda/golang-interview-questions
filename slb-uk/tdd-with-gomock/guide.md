# TDD with GoMock — a step-by-step guide for absolute beginners

This tutorial teaches **Test-Driven Development (TDD)** in Go using **GoMock**. You’ll learn the TDD loop, how to design for testability, generate mocks with `mockgen`, and write expressive, maintainable tests with real-world style examples.

---

## 0) Prerequisites

* Go 1.20+ installed
* A terminal and editor (VS Code recommended)
* Basic Go syntax knowledge (functions, interfaces, packages)
* Go modules enabled

```bash
go version
```

---

## 1) Core definitions (plain English)

**TDD (Test-Driven Development):**
A development workflow where you write a failing test **first**, then write the minimal code to pass it, and finally **refactor**. The loop is called **Red → Green → Refactor**.

**Test double:**
A generic term for stand-ins used in tests.

* **Mock:** Verifies **behavior** (which methods were called, with what arguments, how many times).
* **Stub:** Provides canned **answers** to calls.
* **Spy:** Records how it was used (calls, arguments) for later assertions.
* **Fake:** A lightweight working implementation (e.g., in-memory DB).

**GoMock:**
A mocking framework for Go that generates mocks from interfaces and lets you declare expectations: *“this method will be called with these arguments and will return those values”*.

**mockgen:**
The code generator used by GoMock to produce mock types from your interfaces.

**Dependency inversion (DIP):**
Depend on **interfaces**, not concrete implementations. This is the key to easy mocking.

---

## 2) Project setup

We’ll build a tiny **Order Service** that charges a customer and saves the order. It depends on two collaborators:

* `PaymentGateway` (external service)
* `OrderRepo` (persistence)

We’ll TDD the `PlaceOrder` behavior using GoMock.

```bash
mkdir tdd-gomock-demo && cd tdd-gomock-demo
go mod init github.com/yourname/tdd-gomock-demo
go install github.com/golang/mock/mockgen@latest
go get github.com/golang/mock/gomock@latest
```

Recommended folder layout:

```
.
├── internal/
│   ├── domain/        # business types and interfaces
│   ├── order/         # order service implementation
│   └── testutil/      # test helpers (optional)
└── go.mod
```

---

## 3) Red → Green → Refactor: the TDD loop

### Step A (RED): Write a failing test first

Create the business interfaces (tiny, focused):

```go
// internal/domain/ports.go
package domain

import "context"

type PaymentGateway interface {
    Charge(ctx context.Context, amountCents int64, currency, source string) (txID string, err error)
}

type OrderRepo interface {
    Save(ctx context.Context, o Order) error
}

type Order struct {
    ID          string
    AmountCents int64
    Currency    string
    Status      string // "pending", "paid", "failed"
    PaymentTxID string
}
```

Add the service **contract** we want to implement:

```go
// internal/order/service.go
package order

import (
    "context"

    "github.com/yourname/tdd-gomock-demo/internal/domain"
)

type Service struct {
    pay domain.PaymentGateway
    db  domain.OrderRepo
}

func NewService(pay domain.PaymentGateway, db domain.OrderRepo) *Service {
    return &Service{pay: pay, db: db}
}

// PlaceOrder charges and, if successful, persists the order.
// We'll implement this AFTER we write tests (TDD).
func (s *Service) PlaceOrder(ctx context.Context, o domain.Order, source string) (domain.Order, error) {
    // TODO (write test first)
    panic("not implemented")
}
```

Generate **mocks** for the interfaces:

```bash
mockgen -source=internal/domain/ports.go \
  -destination=internal/domain/mocks/ports_mocks.go \
  -package=mocks
```

> Tip: Add a `//go:generate` line so future devs can just run `go generate ./...`

```go
// internal/domain/ports.go
//go:generate mockgen -source=ports.go -destination=./mocks/ports_mocks.go -package=mocks
```

Now, write the **first test** for the happy path:

```go
// internal/order/service_test.go
package order_test

import (
    "context"
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/yourname/tdd-gomock-demo/internal/domain"
    "github.com/yourname/tdd-gomock-demo/internal/domain/mocks"
    "github.com/yourname/tdd-gomock-demo/internal/order"
)

func TestService_PlaceOrder_Success(t *testing.T) {
    t.Parallel()

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockPay := mocks.NewMockPaymentGateway(ctrl)
    mockRepo := mocks.NewMockOrderRepo(ctrl)

    svc := order.NewService(mockPay, mockRepo)

    in := domain.Order{ID: "ord_1", AmountCents: 4999, Currency: "INR", Status: "pending"}
    source := "tok_visa"

    // Expectations:
    mockPay.
        EXPECT().
        Charge(gomock.Any(), int64(4999), "INR", source).
        Return("tx_abc123", nil).
        Times(1)

    mockRepo.
        EXPECT().
        Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Order{})).
        Return(nil).
        Times(1)

    out, err := svc.PlaceOrder(context.Background(), in, source)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if out.Status != "paid" || out.PaymentTxID != "tx_abc123" {
        t.Fatalf("unexpected order: %+v", out)
    }
}
```

Run the test:

```bash
go test ./...
```

It **fails** (RED) because `PlaceOrder` isn’t implemented. Perfect.

---

### Step B (GREEN): Implement the minimum to pass

```go
// internal/order/service.go
package order

import (
    "context"
    "fmt"

    "github.com/yourname/tdd-gomock-demo/internal/domain"
)

type Service struct {
    pay domain.PaymentGateway
    db  domain.OrderRepo
}

func NewService(pay domain.PaymentGateway, db domain.OrderRepo) *Service {
    return &Service{pay: pay, db: db}
}

func (s *Service) PlaceOrder(ctx context.Context, o domain.Order, source string) (domain.Order, error) {
    txID, err := s.pay.Charge(ctx, o.AmountCents, o.Currency, source)
    if err != nil {
        o.Status = "failed"
        return o, fmt.Errorf("charge failed: %w", err)
    }
    o.Status = "paid"
    o.PaymentTxID = txID

    if err := s.db.Save(ctx, o); err != nil {
        // Real systems might refund/compensate here
        return o, fmt.Errorf("save failed: %w", err)
    }

    return o, nil
}
```

Run tests again:

```bash
go test ./...
```

They should pass (GREEN).

---

### Step C (REFACTOR): Improve design without breaking tests

* Keep interfaces slim
* Improve names and error wrapping
* Consider adding idempotency logic later
* Move constants to `const` if needed

Run tests after each refactor.

---

## 4) Growing test coverage with table-driven tests

Add more behavior: handle charge failure and save failure.

```go
// internal/order/service_test.go
func TestService_PlaceOrder(t *testing.T) {
    t.Parallel()

    cases := []struct {
        name           string
        chargeErr      error
        saveErr        error
        wantStatus     string
        wantErrSubstr  string
    }{
        {"success", nil, nil, "paid", ""},
        {"charge fails", errors.New("card declined"), nil, "failed", "charge failed"},
        {"save fails", nil, errors.New("db down"), "paid", "save failed"},
    }

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockPay := mocks.NewMockPaymentGateway(ctrl)
            mockRepo := mocks.NewMockOrderRepo(ctrl)
            svc := order.NewService(mockPay, mockRepo)

            in := domain.Order{ID: "ord_1", AmountCents: 4999, Currency: "INR", Status: "pending"}
            source := "tok_visa"

            mockPay.EXPECT().
                Charge(gomock.Any(), int64(4999), "INR", source).
                Return("tx_abc123", tc.chargeErr).
                Times(1)

            // Only expect Save if charge succeeds
            if tc.chargeErr == nil {
                mockRepo.EXPECT().
                    Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Order{})).
                    Return(tc.saveErr).
                    Times(1)
            }

            out, err := svc.PlaceOrder(context.Background(), in, source)

            if tc.wantErrSubstr == "" && err != nil {
                t.Fatalf("unexpected err: %v", err)
            }
            if tc.wantErrSubstr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErrSubstr)) {
                t.Fatalf("want err containing %q, got %v", tc.wantErrSubstr, err)
            }
            if out.Status != tc.wantStatus {
                t.Fatalf("want status %q, got %q", tc.wantStatus, out.Status)
            }
        })
    }
}
```

---

## 5) GoMock essentials you’ll use every day

### a) Argument matchers

* `gomock.Any()` — match any value
* `gomock.Eq(v)` — exact equality
* `gomock.Nil()` / `gomock.Not(gomock.Nil())`
* `gomock.AssignableToTypeOf(T{})` — any value of that type

**Custom matcher** (e.g., ensure amount > 0):

```go
type positiveAmount int64

func (positiveAmount) Matches(x interface{}) bool {
    v, ok := x.(int64)
    return ok && v > 0
}
func (positiveAmount) String() string { return "amount > 0" }

// usage:
mockPay.EXPECT().
  Charge(gomock.Any(), positiveAmount(0), "INR", gomock.Any()).
  Return("tx", nil)
```

### b) Side effects and computed returns

```go
mockRepo.EXPECT().
  Save(gomock.Any(), gomock.AssignableToTypeOf(domain.Order{})).
  DoAndReturn(func(_ context.Context, o domain.Order) error {
      if o.Status != "paid" || o.PaymentTxID == "" {
          return errors.New("invalid order")
      }
      return nil
  })
```

### c) Call order

```go
gomock.InOrder(
    mockPay.EXPECT().Charge(gomock.Any(), int64(4999), "INR", "tok_visa").Return("tx", nil),
    mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil),
)
```

### d) Call counts

* `Times(n)`
* `AnyTimes()`
* `MinTimes(n)` / `MaxTimes(n)`

### e) Lifecycle

* Always create a controller: `ctrl := gomock.NewController(t)`
* Ensure `ctrl.Finish()` runs (`defer` or `t.Cleanup`)

---

## 6) Real-world scenarios (with patterns & examples)

### Scenario 1 — Mocking an external HTTP API

Instead of mocking `*http.Client` (not an interface), mock its **transport**:

```go
// internal/domain/http_client.go
package domain

import "net/http"

type RoundTripper interface {
    RoundTrip(*http.Request) (*http.Response, error)
}
```

Wire your client like:

```go
client := &http.Client{Transport: myTransport}
```

Generate a mock for `RoundTripper` and assert the correct HTTP request is made (URL, headers, body), returning a crafted `*http.Response`. Pattern: **Adapter + Mock Transport**.

### Scenario 2 — Messaging/Kafka publisher

Wrap the library in your own interface:

```go
// internal/domain/publisher.go
package domain

import "context"

type Publisher interface {
    Publish(ctx context.Context, topic string, key, value []byte) error
}
```

Mock `Publisher` in tests to verify **at-least once** logic, retries, and backoff without involving a real broker. Pattern: **Port/Adapter**.

### Scenario 3 — Time, randomness, and IDs

Inject a clock and ID generator:

```go
type Clock interface{ Now() time.Time }
type IDGen interface{ New() string }
```

Mocks let you assert time-dependent behavior and stable IDs. Pattern: **Deterministic boundaries**.

### Scenario 4 — Database repositories

Keep `OrderRepo` **interface-shaped**. In unit tests, mock it to validate call sequences and domain rules. Keep DB integration tests separate (a small set) for migrations/SQL. Pattern: **Unit vs. integration split**.

### Scenario 5 — Orchestrations with SAGA / compensation

Use `gomock.InOrder(...)` to enforce **charge → reserve → save** sequences. Add tests for **compensation** paths (e.g., refund if `Save` fails). Pattern: **Workflow verification**.

---

## 7) Going a level deeper: designing for testability

* **Narrow interfaces** (one or two methods) keep your mocks simple.
* **Constructor injection** (`NewService(dep1, dep2)`) makes dependencies explicit.
* **Functional options** can inject optional deps (logger, metrics).
* **Avoid global state**; inject everything.
* **Return rich domain errors** (wrap with `%w`) so tests can assert error categories.

---

## 8) Developer workflow tips

### a) Automate mock generation

Add to `Makefile`:

```makefile
generate:
	go generate ./...

test:
	go test ./... -race -coverprofile=cover.out
	@go tool cover -func=cover.out | tail -n +2 || true
```

### b) Keep tests readable

* Use **Given/When/Then** comments
* Use **table-driven** tests for variants
* Prefer **behavior** assertions over implementation details
* Name test cases with the **reason** they fail/succeed

### c) Debugging common GoMock errors

* **“no matching calls”** → argument mismatch; add matchers or log args in `DoAndReturn`.
* **“expected call … not satisfied”** → you didn’t call it; check branches or `Times`.
* **“unexpected call”** → overspecified; loosen with `AnyTimes` or adjust flow.
* **Mocks not generated** → run `go generate ./...` or `mockgen …` manually; check import paths.

---

## 9) Extra example: validating input and retries

Add retries for transient payment errors and validate amounts.

```go
// internal/order/service.go (excerpt)
func (s *Service) PlaceOrder(ctx context.Context, o domain.Order, source string) (domain.Order, error) {
    if o.AmountCents <= 0 {
        return o, fmt.Errorf("invalid amount")
    }
    var txID string
    var err error
    for i := 0; i < 3; i++ {
        txID, err = s.pay.Charge(ctx, o.AmountCents, o.Currency, source)
        if err == nil {
            break
        }
        if !isTransient(err) {
            break
        }
    }
    if err != nil {
        o.Status = "failed"
        return o, fmt.Errorf("charge failed: %w", err)
    }
    o.Status = "paid"
    o.PaymentTxID = txID
    if err := s.db.Save(ctx, o); err != nil {
        return o, fmt.Errorf("save failed: %w", err)
    }
    return o, nil
}
```

Test the retry with `Times(3)`:

```go
mockPay.EXPECT().
  Charge(gomock.Any(), int64(4999), "INR", source).
  Return("", transientErr).
  Times(2) // first two transient

mockPay.EXPECT().
  Charge(gomock.Any(), int64(4999), "INR", source).
  Return("tx_final", nil).
  Times(1)
```

---

## 10) End-to-end checklist for beginners

1. **Model the behavior** you want (write tests first).
2. **Define small interfaces** around external boundaries.
3. `mockgen` your interfaces (`//go:generate` is your friend).
4. In tests, **set expectations** with matchers and return values.
5. Implement the **minimal code** to pass (Green).
6. **Refactor** with confidence; tests keep you safe.
7. Use **table-driven tests** and **InOrder** for workflows.
8. Keep **unit tests hermetic** (no network/DB); do a few focused integration tests separately.
9. Run `go test -race -cover` regularly.
10. Automate in a **Makefile** or CI.

---

## 11) FAQ (quick wins)

* **Should I mock standard library types?**
  Only if there’s an interface (e.g., `io.Reader`). Otherwise, create your **own** small interface and adapt.

* **When not to mock?**
  Pure functions don’t need mocks. Prefer **fakes** for complex protocols you own.

* **Are mocks “bad”?**
  Mocks are great for **behavioral** verification at the system boundary. Use judiciously; avoid over-specifying internal details.

---

## 12) Complete minimal working example (copy-paste ready)

**Files you’ll have:**

* `internal/domain/ports.go` (interfaces & `//go:generate`)
* `internal/domain/mocks/ports_mocks.go` (generated)
* `internal/order/service.go` (implementation)
* `internal/order/service_test.go` (tests)

**Commands:**

```bash
go mod init github.com/yourname/tdd-gomock-demo
go install github.com/golang/mock/mockgen@latest
go get github.com/golang/mock/gomock@latest
go generate ./...
go test ./... -race -cover
```

---

## 13) What you can build next (real-world practice ideas)

* **User signup workflow:** Email verification (mock EmailSender), user repo (mock), rate-limit (mock Clock).
* **Payment + inventory reservation:** Use `gomock.InOrder` for multi-step workflow; add compensation tests.
* **Webhook receiver:** Mock downstream publisher; test retries and DLQ (dead-letter queue) publishing.
* **HTTP client wrapper:** Mock `RoundTripper` to test request formation, auth headers, and backoff.
