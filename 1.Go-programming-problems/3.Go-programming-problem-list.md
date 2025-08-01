# Top 10 Go Programming Problems for BigCommerce Interview

These problems are curated for the BigCommerce interview process with a focus on Go programming, REST APIs, system design, concurrency, and clean coding practices.

---

## 🚀 1. Design a URL Shortener Service (like bit.ly)
**Problem:**  
Build a Go-based in-memory URL shortener with APIs to shorten and expand URLs.

**Solution Hint:**  
Use a `map[string]string` to store short-to-long mappings. Generate short keys using base62 encoding of an incrementing counter or `crypto/rand`. Expose REST endpoints using `net/http` or `gin`.

---

## 🧵 2. Write a Safe Concurrent Counter
**Problem:**  
Design a counter that can be safely incremented by multiple goroutines.

**Solution Hint:**  
Use `sync.Mutex` or `sync/atomic` package for thread safety. Wrap the counter in a struct with lock methods.

---

## 📦 3. Implement a Rate Limiter
**Problem:**  
Write a rate limiter that allows only N requests per second per user.

**Solution Hint:**  
Use a `map[string]*rate.Limiter` from `golang.org/x/time/rate` package. For custom logic, use a token bucket pattern with `time.Ticker`.

---

## 🧹 4. Graceful Shutdown of an HTTP Server
**Problem:**  
Create an HTTP server that shuts down gracefully when it receives a SIGINT or SIGTERM signal.

**Solution Hint:**  
Use `http.Server` with `Shutdown(ctx)` and handle OS signals using the `os/signal` package and `context.WithTimeout`.

---

## 🧪 5. Implement a Middleware Chain
**Problem:**  
Create your own middleware functions in a mini web framework.

**Solution Hint:**  
Use handler functions that accept `http.Handler` and return `http.Handler`. Chain them together for logging, authentication, etc.

---

## 🔁 6. Build a Retry Mechanism
**Problem:**  
Implement a generic retry function that retries an operation N times with exponential backoff.

**Solution Hint:**  
Use function types like `func() error` and loops with `time.Sleep(time.Duration(math.Pow(2, i)) * time.Millisecond)`.

---

## 🧠 7. Parse and Validate JSON Input
**Problem:**  
Create a REST endpoint that accepts JSON, validates required fields, and returns appropriate responses.

**Solution Hint:**  
Use `encoding/json` to decode into a struct, then validate fields manually or with validation libraries like `go-playground/validator`.

---

## 🧊 8. Prevent Goroutine Leaks
**Problem:**  
Write code that launches a goroutine for processing messages but stops gracefully when the app shuts down.

**Solution Hint:**  
Use `context.Context` to pass cancellation signals to the goroutine. Always select on `ctx.Done()` and the message channel.

---

## 🕵️ 9. Implement Custom JSON Marshalling
**Problem:**  
Modify how a struct gets marshaled to JSON (e.g., omit empty fields, format timestamps).

**Solution Hint:**  
Implement the `MarshalJSON()` method for the struct. You can use aliases or build a custom `map[string]interface{}`.

---

## 🧬 10. Build a CRUD API with MySQL
**Problem:**  
Build a basic REST API to Create, Read, Update, Delete `Product` entities stored in MySQL.

**Solution Hint:**  
Use `database/sql` with `mysql` driver. Structure code with a repository layer. Use `gorm` for faster setup if allowed.

---