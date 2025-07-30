# Go Programming Theory Questions
# 🧠 Top 25 Go Theory Questions with Detailed Answers

## 1. What are the key features of Go?
Go is statically typed, compiled, supports concurrency (goroutines, channels), fast compilation, garbage collection, and a strong standard library.

## 2. What are goroutines?
Lightweight threads managed by Go runtime. Use `go` keyword to start a goroutine.

**Example:**
```go
func sayHello() {
    fmt.Println("Hello from goroutine")
}

func main() {
    go sayHello()
    time.Sleep(time.Second)
}
```

Lightweight threads managed by Go runtime. Use `go` keyword to start a goroutine.

## 3. What are channels?
Typed pipes that connect goroutines for communication. Use `<-` to send/receive.

**Example:**
```go
func main() {
    ch := make(chan string)
    go func() {
        ch <- "Hello from channel"
    }()
    fmt.Println(<-ch)
}
```

Typed pipes that connect goroutines for communication. Use `<-` to send/receive.

## 4. What is a buffered vs unbuffered channel?
Buffered channels block only when full/empty. Unbuffered channels block until both sender and receiver are ready.

## 5. What is the `select` statement?
Allows a goroutine to wait on multiple communication operations.

**Example:**
```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() { ch1 <- "one" }()
    go func() { ch2 <- "two" }()

    select {
    case msg1 := <-ch1:
        fmt.Println("received", msg1)
    case msg2 := <-ch2:
        fmt.Println("received", msg2)
    }
}
```

Allows a goroutine to wait on multiple communication operations.

## 6. What is a Go interface?
A collection of method signatures. Implements implicitly—no `implements` keyword.

**Example:**
```go
type Speaker interface {
    Speak() string
}

type Person struct {}

func (p Person) Speak() string {
    return "Hello from Person"
}

func greet(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    p := Person{}
    greet(p)
}
```

A collection of method signatures. Implements implicitly—no `implements` keyword.

## 7. Difference between value and pointer receivers?
Value receivers get a copy; pointer receivers can mutate the original struct.

## 8. What is `defer`?
Schedules a function call to run after the surrounding function returns. Useful for cleanup.

**Example:**
```go
func main() {
    defer fmt.Println("world")
    fmt.Println("hello")
}
// Output:
// hello
// world
```

Schedules a function call to run after the surrounding function returns. Useful for cleanup.

## 9. What is a struct in Go?
Custom composite type that groups fields. Acts like a class without inheritance.

## 10. How does garbage collection work in Go?
Go uses concurrent garbage collector (mark and sweep), automatically freeing memory no longer referenced.

## 11. Explain Go memory model.
Defines how goroutines interact with memory. Use channels or `sync/atomic` to ensure consistency.

## 12. What is `panic` and `recover`?
`panic` stops normal execution. `recover` regains control within deferred functions.

**Example:**
```go
func mayPanic() {
    panic("something went wrong")
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()
    mayPanic()
    fmt.Println("After panic")
}
```

`panic` stops normal execution. `recover` regains control within deferred functions.

## 13. How does Go handle error handling?
Explicit error returns (`error` type). No exceptions.

**Example:**
```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

Explicit error returns (`error` type). No exceptions.

## 14. What is the zero value in Go?
Default value for variables: `0`, `false`, `""`, `nil`, etc.

## 15. What is a slice in Go?
Dynamic array. Backed by an underlying array. `len` and `cap` control view.

**Example:**
```go
func main() {
    arr := []int{1, 2, 3}
    fmt.Println("Length:", len(arr))
    fmt.Println("Capacity:", cap(arr))
    arr = append(arr, 4)
    fmt.Println("Appended Slice:", arr)
}
```

Dynamic array. Backed by an underlying array. `len` and `cap` control view.

## 16. What are Go maps?
Unordered key-value data structures. `map[keyType]valueType`

**Example:**
```go
func main() {
    m := make(map[string]int)
    m["one"] = 1
    m["two"] = 2
    fmt.Println(m)
    delete(m, "one")
    fmt.Println("After delete:", m)
}
```

Unordered key-value data structures. `map[keyType]valueType`

## 17. How are packages organized in Go?
One package per directory. Entry point is `main` package with `main()` function.

## 18. What is the `init()` function?
Runs automatically before `main()` for package initialization.

## 19. How to handle concurrency safely?
Use channels, `sync.Mutex`, or atomic operations.

## 20. What is `context.Context` used for?
Carry deadlines, cancellations, and other request-scoped values.

**Example:**
```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    ch := make(chan string)
    go func() {
        time.Sleep(1 * time.Second)
        ch <- "done"
    }()

    select {
    case res := <-ch:
        fmt.Println("Result:", res)
    case <-ctx.Done():
        fmt.Println("Timeout:", ctx.Err())
    }
}
```

### ➕ 10 Follow-up Questions on `context.Context`

---

### 🧩 Real-World Use Cases of `context.Context`

**1. Context in REST APIs**
Used to enforce request timeouts and cancellations.
```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    select {
    case <-time.After(3 * time.Second):
        fmt.Fprintln(w, "Processed")
    case <-ctx.Done():
        http.Error(w, "Request cancelled", http.StatusRequestTimeout)
    }
}
```

**2. Context in HTTP Client Calls**
Used to cancel outbound API calls if deadline exceeded.
```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.example.com/data", nil)
resp, err := http.DefaultClient.Do(req)
```

**3. Context in Goroutines**
Used to cancel background tasks.
```go
func startWorker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker stopped")
            return
        default:
            fmt.Println("Working...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

**4. Context with Databases (e.g., SQL)**
Used to cancel long-running queries.
```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

rows, err := db.QueryContext(ctx, "SELECT * FROM orders")
```

**5. Context in Kafka Consumers**
Used to shutdown consumers gracefully.
```go
func consume(ctx context.Context, ch <-chan string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Stopped consumer")
            return
        case msg := <-ch:
            fmt.Println("Received:", msg)
        }
    }
}
```

**6. Context Propagation in Microservices**
Pass context between services to maintain traceability, deadlines, and correlation IDs.

**7. Context in Worker Pools**
Gracefully stop all workers by broadcasting cancel.

**8. Context in File Uploads/Downloads**
Used to cancel if the client disconnects or timeout is reached.

**9. Context in gRPC Services**
gRPC natively uses context for deadline, metadata, and cancellation.

**10. Context in Cron Jobs or Background Schedulers**
Allows scheduler to cancel job run if needed.

---

### 🔁 Continued: 10 Follow-up Questions on `context.Context`

**1. How do you create a cancellable context?**  
Use `context.WithCancel(parent)` to create a cancellable context. `cancel()` should be called to release resources.

**2. How do you pass context to a function?**  
By adding `ctx context.Context` as the first parameter in your function signature.

**3. When should you cancel a context?**  
Always call the `cancel` function returned by `WithCancel`, `WithTimeout`, or `WithDeadline` to release resources.

**4. What is the difference between `WithTimeout` and `WithDeadline`?**  
`WithTimeout` sets a timeout relative to now. `WithDeadline` sets a specific end time.

**5. How do you detect if a context is done?**  
Use `select { case <-ctx.Done(): ... }`.

**6. What does `ctx.Err()` return?**  
Returns `context.Canceled` or `context.DeadlineExceeded` when done.

**7. How is context propagated in Go applications?**  
Always pass it down through function calls to ensure cancellation and deadlines are respected.

**8. How do you attach values to context?**  
Use `context.WithValue(parent, key, value)`. Avoid using it for passing business logic data.

**9. Why is context passed as the first parameter?**  
To enforce consistent use and ensure that cancellation signals and deadlines are handled early.

**10. What are best practices with `context.Context`?**  
- Always cancel context
- Don’t store in struct fields
- Don’t pass nil context
- Use `context.Background()` or `context.TODO()` as base

Carry deadlines, cancellations, and other request-scoped values.

**Example:**
```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    ch := make(chan string)
    go func() {
        time.Sleep(1 * time.Second)
        ch <- "done"
    }()

    select {
    case res := <-ch:
        fmt.Println("Result:", res)
    case <-ctx.Done():
        fmt.Println("Timeout:", ctx.Err())
    }
}
```

Carry deadlines, cancellations, and other request-scoped values.

## 21. How is Go different from Java/C++?
Simpler syntax, no classes, interfaces are implicit, goroutines over OS threads.

## 22. What is `iota`?
Auto-incrementing identifier used in const declarations.

## 23. What are Go modules?
Dependency management system. Defined in `go.mod` and `go.sum`.

## 24. What is embedding in Go?
Structs can embed other structs/interfaces to inherit behavior without inheritance.

## 25. What tools are used for testing in Go?
`testing` package, `go test`, `go bench`, plus mocks and coverage via `go tool cover`.

---

Would you like to add code snippets, Playground links, or quiz questions for each theory topic?
