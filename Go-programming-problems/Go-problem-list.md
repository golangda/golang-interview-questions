# 25 Go Interview Problems with Hints and Playground Links

## ⯿️ Beginner (Conceptual & Syntax)

| # | Problem | Concept | Hint | Link |
|--|---------|---------|------|------|
| 1 | What’s the difference between value and pointer receivers in Go? | Method receivers | Value makes a copy, pointer modifies the original | [Playground](https://go.dev/play/p/NT8toN7yGHE) |
| 2 | Write a function to reverse a string in Go. | Strings, runes | Use `[]rune` for Unicode support | [Playground](https://go.dev/play/p/XoErBFL7NQ7) |
| 3 | Explain the difference between `make` and `new`. | Memory allocation | `make` is for slices/maps/channels, `new` allocates zeroed memory | [Playground](https://go.dev/play/p/Q9IcsAOdvTk) |
| 4 | What does `defer` do? Show an example with `panic`. | Defer, panic recovery | Deferred calls execute in LIFO just before function exits | [Playground](https://go.dev/play/p/cAP_41H__Mo) |
| 5 | Write a simple interface and struct that implements it. | Interfaces | Define a `Shape` interface and a `Circle` struct | [Playground](https://go.dev/play/p/ErXrGUOfkUu) |

---

## 🟡 Intermediate (Concurrency, Errors, Structs)

| # | Problem | Concept | Hint | Link |
|--|---------|---------|------|------|
| 6 | Write a function that runs two goroutines and waits for both to finish. | Goroutines, `sync.WaitGroup` | Use `Add`, `Done`, `Wait` | [Playground](https://go.dev/play/p/vG4OHnyB3nW) |
| 7 | Implement a rate limiter using Go channels. | Concurrency | Use `time.Tick` or a buffered channel | [Playground](https://go.dev/play/p/xgeChFGok1U) |
| 8 | Write a JSON decoder for a `User` struct. | JSON unmarshalling | Use `json.Unmarshal([]byte, &obj)` | [Playground](https://go.dev/play/p/BhLKWBk1mFD) |
| 9 | What is `select` in Go? Write a select statement with 2 channels. | `select`, channels | Use timeouts with `select` and `time.After` | [Playground](https://go.dev/play/p/1N3hKglHn2I) |
| 10 | What is the difference between buffered and unbuffered channels? Show an example. | Channels | Buffered channels do not block immediately on send | [Playground](https://go.dev/play/p/l4aLkDbp9LR) |

---

## 🟠 Applied Go (Tools, Testing, Practices)

| # | Problem | Concept | Hint | Link |
|--|---------|---------|------|------|
| 11 | Write a table-driven unit test for a `Sum` function. | Testing | Use `t.Run()` for subtests | [Playground](https://go.dev/play/p/ikJDUDuHDd0) |
| 12 | How would you structure a Go project with multiple packages? | Code organization | Use `cmd/`, `pkg/`, `internal/` |
| 13 | Write a Go program that reads environment variables and falls back to defaults. | `os.Getenv`, fallback | Use `os.LookupEnv` |
| 14 | Explain the role of `go.mod` and `go.sum`. | Modules | `go.mod` declares, `go.sum` secures dependencies |
| 15 | What is `interface{}` in Go? When should you avoid it? | Empty interface | Prefer typed interfaces; use when flexibility is required |

---

## 🔵 Advanced (Design, Internals, Optimization)

| # | Problem | Concept | Hint | Link |
|--|---------|---------|------|------|
| 16 | Implement a thread-safe in-memory cache. | `sync.Map`, `sync.RWMutex` | Prefer `sync.RWMutex` for read-heavy cache | [Playground](https://go.dev/play/p/ZVrwaNjJiqz) |
| 17 | How would you detect goroutine leaks? | Concurrency, debugging | Use `pprof`, bounded channel, context cancellation |
| 18 | What is `context.Context`? Show usage in an API handler. | Context, cancellation | Use `context.WithTimeout()` | [Playground](https://go.dev/play/p/fOr9axGvHeD) |
| 19 | Compare Go’s garbage collector to manual memory management. | Runtime | Go is GC-based, but you manage lifecycles via design |
| 20 | Build a simple middleware for logging API requests in Go. | HTTP handler chaining | Wrap `http.HandlerFunc` | [Playground](https://go.dev/play/p/kEDw-3Zdk_w) |

---

## 🔴 Expert / Production Scenarios

| # | Problem | Concept | Hint | Link |
|--|---------|---------|------|------|
| 21 | How do you benchmark Go code? | `testing.B` | Use `go test -bench .` | [Playground](https://go.dev/play/p/LXuwE-wRU1U) |
| 22 | Write a Go program to gracefully shut down an HTTP server. | `context.WithCancel`, `os.Signal` | Use `signal.NotifyContext` | [Playground](https://go.dev/play/p/fRC3FYqxHRg) |
| 23 | How does Go prevent data races? | Race detector | Use `go run -race` and `sync` primitives |
| 24 | Implement a worker pool with fixed concurrency. | Channels + goroutines | Use job queue + `WaitGroup` | [Playground](https://go.dev/play/p/pEr4p3U40Tz) |
| 25 | How to integrate Prometheus metrics into a Go service? | Observability | Use `promhttp.Handler()` | [Guide](https://prometheus.io/docs/guides/go-application/) |

---