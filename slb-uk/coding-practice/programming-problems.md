# Top 10 Golang Concurrency Patterns (Interview‑Ready)

Includes your must‑haves: **Worker Pool, Rate Limiter, Fan‑In/Fan‑Out, Context Binding, LRU Caching**.

---

## 1) Worker Pool (bounded concurrency)

**Idea:** N workers read tasks from a channel; limits parallelism and evens load.
**Use when:** Many similar jobs (I/O or CPU) to process.

```go
var wg sync.WaitGroup
for i := 0; i < N; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range tasks { results <- do(t) }
	}()
}
go func() { wg.Wait(); close(results) }()
```

**Pitfalls:** Close `results` only after workers finish; apply backpressure to avoid blocking producers/consumers.

---

## 2) Fan‑Out

**Idea:** Launch multiple goroutines to process the same input stream in parallel.
**Use when:** Parallelize independent work on each item.

```go
for i := 0; i < N; i++ {
	go func() { for v := range in { out <- f(v) } }()
}
```

**Pitfalls:** Coordinate output closing (see fan‑in); preserve ordering if required.

---

## 3) Fan‑In

**Idea:** Merge multiple channels into one.
**Use when:** Aggregate results from many producers.

```go
func fanIn[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(chs))
	for _, ch := range chs {
		go func(c <-chan T) { defer wg.Done(); for v := range c { out <- v } }(ch)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
```

**Pitfalls:** Never double‑close; propagate cancellation upstream.

---

## 4) Pipeline

**Idea:** Staged processing with goroutines connected by channels.
**Use when:** Clear stages (gen → transform → filter → aggregate) with streaming.

```go
stage2 := func(in <-chan A) <-chan B {
	out := make(chan B)
	go func(){
		defer close(out)
		for a := range in { out <- g(a) }
	}()
	return out
}
```

**Pitfalls:** Always add cancellation paths (see context binding) to avoid goroutine leaks.

---

## 5) Context Binding & Cancellation Propagation

**Idea:** Pass `context.Context`; each `select` listens on `<-ctx.Done()`.

```go
func work(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-doSomething():
		return nil
	}
}
```

**Pitfalls:** Forgetting to pass/observe `ctx`; not calling `cancel()` when you own it.

---

## 6) Rate Limiter (ticker or token bucket)

**Idea:** Throttle ops to X/sec; smooth spikes.
**Use when:** External APIs, DB load, fair use.

```go
// With ticker
tick := time.NewTicker(100 * time.Millisecond)
defer tick.Stop()
for req := range in {
	<-tick.C
	handle(req)
}

// With token bucket (preferred)
lim := rate.NewLimiter(10, 20) // 10 rps, burst 20
_ = lim.Wait(ctx)              // block until a token
```

**Pitfalls:** Leaking tickers/timers; using `time.After` in tight loops.

---

## 7) Concurrent LRU Cache

**Idea:** LRU map guarded by a lock or “actor” goroutine; evicts least‑recently‑used on cap.
**Use when:** Hot, repeated lookups with limited memory.

**Mutex version (simple):**

```go
type LRU struct{ mu sync.RWMutex; m map[string]*list.Element; l *list.List /* ... */ }
func (c *LRU) Get(k string) (v any, ok bool) { c.mu.RLock(); /* ... */ c.mu.RUnlock(); return }
func (c *LRU) Set(k string, v any) { c.mu.Lock(); /* move/evict */ c.mu.Unlock() }
```

**Actor version (no locks, single owner):**

```go
type req struct{ get string; setK string; setV any; reply chan any }
func runLRU(in <-chan req) { /* single goroutine owns map+list */ }
```

**Pitfalls:** Lock contention → consider sharding; keep ops O(1); avoid eviction races.

---

## 8) Semaphore (limit in‑flight tasks)

**Idea:** Use a buffered channel or `x/sync/semaphore` to bound concurrent sections.
**Use when:** Only K goroutines inside a critical region/expensive call.

```go
sem := make(chan struct{}, K)
do := func() {
	sem <- struct{}{}       // acquire
	defer func(){ <-sem }() // release
	// critical / limited work
}
```

**Pitfalls:** Always release (use `defer`); avoid deadlocks on panics/errors.

---

## 9) `errgroup` + Context (scatter‑gather with fail‑fast)

**Idea:** Run M tasks in parallel; cancel all on first error.
**Use when:** Parallel I/O calls, fan‑out requests.

```go
g, ctx := errgroup.WithContext(ctx)
for _, u := range urls {
	u := u
	g.Go(func() error { return fetch(ctx, u) })
}
if err := g.Wait(); err != nil { /* handle */ }
```

**Pitfalls:** Work must respect `ctx`; avoid capturing loop vars incorrectly.

---

## 10) Graceful Shutdown (signals + WaitGroup)

**Idea:** On `SIGINT/SIGTERM`, cancel context, stop intake, drain work, wait.
**Use when:** Services/daemons with goroutines.

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer cancel()
go serve(ctx)
<-ctx.Done()     // signal received
wg.Wait()        // wait for workers to finish
```

**Pitfalls:** Forgetting to stop accepting new work; not closing channels; missing time‑bounded shutdown.
