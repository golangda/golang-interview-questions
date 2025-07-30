# 📘 Contexts in Go: A Complete Developer's Guide

---

## 📌 Table of Contents

1. [What is Context in Go?](#what-is-context-in-go)
2. [Why is Context Important?](#why-is-context-important)
3. [Context Package Overview](#context-package-overview)
4. [Types of Contexts](#types-of-contexts)
5. [Creating Contexts](#creating-contexts)
6. [Using Context in Goroutines](#using-context-in-goroutines)
7. [Context with HTTP Request](#context-with-http-request)
8. [Context with Database Operations](#context-with-database-operations)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)
11. [Interview Questions](#interview-questions)

---

## 🔍 What is Context in Go?

> `context` is a standard library package in Go used for controlling deadlines, cancellation signals, and request-scoped values across API boundaries and goroutines.

It helps manage long-running processes, timeouts, and cancellation — crucial in modern, concurrent programs (e.g., web servers, microservices, background jobs).

---

## ❓ Why is Context Important?

- **Cancel goroutines** to avoid leaks.
- **Set timeouts** and deadlines to limit execution time.
- **Pass request-scoped data** like trace IDs, auth tokens, etc.
- **Gracefully shutdown** servers and tasks on signals.

---

## 📦 Context Package Overview

```go
import "context"
```

Core types and functions:

| Function | Purpose |
|---------|---------|
| `context.Background()` | Root context (empty, never canceled) |
| `context.TODO()` | Placeholder context (for future use) |
| `context.WithCancel(ctx)` | Returns new context with cancel function |
| `context.WithTimeout(ctx, duration)` | Cancels after timeout |
| `context.WithDeadline(ctx, time)` | Cancels at exact time |
| `context.WithValue(ctx, key, value)` | Associates key-value pair with context |

---

## 🧱 Types of Contexts

### 1. `context.Background()`

Used at the top level — server startup, `main()`, root tasks.

```go
ctx := context.Background()
```

---

### 2. `context.TODO()`

Temporary placeholder — when you're unsure what context to use.

```go
ctx := context.TODO()
```

---

### 3. `context.WithCancel(parent)`

Used when you want manual control to cancel context.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

---

### 4. `context.WithTimeout(parent, timeout)`

Automatically cancels after the specified duration.

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```

---

### 5. `context.WithDeadline(parent, time.Time)`

Cancels at a specific point in time.

```go
deadline := time.Now().Add(1 * time.Minute)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

---

## 🛠️ Using Context in Goroutines

### 🔁 Problem: Goroutines leak when not terminated

#### ✅ Solution: Use context to signal termination.

```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine exiting...")
				return
			default:
				fmt.Println("Working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}
```

---

## 🌐 Context in HTTP Handlers

```go
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "Done")
	case <-ctx.Done():
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
	}
}
```

📌 Real-world: Handles cancellation if the client disconnects before the request is served.

---

## 💾 Context in Database Queries (e.g., PostgreSQL, MySQL)

```go
func queryUser(ctx context.Context, db *sql.DB, id int) (User, error) {
	query := "SELECT name, email FROM users WHERE id = ?"
	row := db.QueryRowContext(ctx, query, id)

	var user User
	err := row.Scan(&user.Name, &user.Email)
	return user, err
}
```

🧠 Tip: `QueryRowContext`, `ExecContext` are preferred over raw queries to avoid hanging DB connections.

---

## 🧭 Best Practices

✅ Always `defer cancel()`  
✅ Prefer `WithTimeout` for external calls  
✅ Pass context as the first parameter (`ctx context.Context`)  
✅ Use typed keys for `WithValue` to avoid collision

---

## ❌ Common Mistakes

| Mistake | Why It's Bad |
|--------|--------------|
| Not cancelling contexts | Goroutines keep running |
| Using `context.WithValue` for passing config/data | It's for metadata only (trace ID, etc.) |
| Ignoring `ctx.Err()` | You won't know why it failed |
| Forgetting `select { case <-ctx.Done(): ... }` in goroutines | Causes memory leaks |

---

## 🧠 Real-World Use Cases

| Scenario | Context Usage |
|---------|---------------|
| REST API server | Cancel handler if client disconnects |
| gRPC calls | Timeout-based cancellation of service calls |
| Database queries | Prevent slow DB queries from blocking |
| File uploads | Cancel if client drops midway |
| Microservices | Pass trace ID, auth token with `WithValue()` |
| Scheduled jobs | Timeout long jobs, cancel on shutdown signal |

---

## 📌 Interview Questions on Go Context

### 1. What is context in Go and why is it needed?
**Answer:** Context is used to manage deadlines, cancellation signals, and request-scoped values. It prevents goroutine leaks and supports clean shutdowns.

---

### 2. What is the difference between `context.Background()` and `context.TODO()`?
**Answer:**
- `Background()` is used in main or root-level logic.
- `TODO()` is used when you haven’t decided which context to use yet.

---

### 3. What does `ctx.Done()` do?
**Answer:** Returns a channel that's closed when the context is canceled or times out. Use it in `select` to handle cancellation.

---

### 4. How is context propagated in an HTTP server?
**Answer:** It’s accessible via `r.Context()` in handlers and is canceled when the client closes the connection or the server times out.

---

### 5. What happens if you don’t call `cancel()` in `WithCancel()` or `WithTimeout()`?
**Answer:** Resource leak! Context won’t be garbage collected, and goroutines may not exit properly.

---

### 6. How do you pass values using context?
**Answer:**

```go
ctx = context.WithValue(ctx, "traceID", "abc-123")
val := ctx.Value("traceID")
```

Use **custom types** for keys to avoid conflicts.

---

## 📘 Summary Table

| Context Function | Use Case |
|------------------|----------|
| `context.Background()` | Root context in `main()` or tests |
| `context.TODO()` | Placeholder during development |
| `context.WithCancel()` | Manual cancellation (e.g., signal interrupt) |
| `context.WithTimeout()` | Automatically cancel after a duration |
| `context.WithDeadline()` | Cancel at a specific time |
| `context.WithValue()` | Pass request-scoped metadata (not config!) |

---

## ✅ Try it on Go Playground

👉 [Basic Timeout Example](https://go.dev/play/p/gK9Crf8bfSe)
