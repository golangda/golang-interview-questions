# Go Tooling — Beginner’s Printable Guide

Hands‑on guide to Go tooling, with copy‑paste code and a ready Makefile + pre‑commit.

## 1) Init a module
```bash
mkdir go-tooling-demo && cd go-tooling-demo
go mod init example.com/go-tooling-demo
```
## 2) Build & Run
`go run .` (iterate quickly) • `go build -o app && ./app` (ship/build)
### main.go
```go
package main
import "fmt"
func main(){ fmt.Println("Hello, Go tooling!") }
```
## 3) Format & Vet
```bash
go fmt ./...
go vet ./...
```
## 4) Tests, Coverage, Benchmarks
```bash
go test ./... -cover
go test ./... -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
go test -bench=. ./word
```
## 5) Dependencies
```bash
go get github.com/google/uuid@latest
go mod tidy
go list ./...           # packages
go list -m all          # modules
```
## 6) Docs
```bash
go doc strings.Fields
go doc -all fmt
```
## 7) Race Detector
```bash
go test -race ./...
```
## 8) Profiling (pprof)
```bash
go test ./word -run TestRepeatCount -cpuprofile=cpu.out
go tool pprof cpu.out
```
## 9) Code Generation
Use `//go:generate` + `go generate ./...` to standardize generation.

## 10) Linting (golangci-lint)
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
$(go env GOPATH)/bin/golangci-lint run
```

## 11) Debugging (Delve)
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
$(go env GOPATH)/bin/dlv debug
```

## 12) Workspaces & Env
```bash
go env GOPATH GOMODCACHE GOROOT GOOS GOARCH
go work init ./app ./lib
```

## Real‑world scenarios
- PR broke quality → `go fmt`, `go vet`, `golangci-lint run`, `go test -race -cover`
- Perf regression → add bench + pprof
- Flaky under load → `-race`, targeted tests, `dlv`
- Dep mess → `go get ...`, `go mod tidy`, `go list -m all`

---
### Daily commands
`go run .` • `go build -o app` • `go fmt ./...` • `go vet ./...` • `golangci-lint run` • `go test ./...` • `go mod tidy` • `go tool cover -html=cover.out -o cover.html`

## pprof Quick‑Start

There are two common ways to use pprof: from tests and from a running HTTP server.

### A) From tests (no server code)
CPU profile:
```bash
go test ./word -run TestRepeatCount -cpuprofile=cpu.out
go tool pprof -http=:0 cpu.out
```
Heap profile:
```bash
go test ./word -run TestRepeatCount -memprofile=heap.out -benchmem
go tool pprof -http=:0 heap.out
```

### B) From a running service (HTTP pprof)
We included a small server at `pprof/server.go` that exposes pprof on `localhost:6060` and a `/work` endpoint to generate CPU load.

Run it:
```bash
go run ./pprof/server.go
```
Hit the app to create some load in another terminal:
```bash
curl http://localhost:6060/work
```
Capture a 30‑second CPU profile (in a third terminal):
```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```
Other useful endpoints:
- Heap (current live allocations): `http://localhost:6060/debug/pprof/heap`
- Goroutines snapshot:          `http://localhost:6060/debug/pprof/goroutine?debug=2`
- Mutex/blocking profiles (enabled in code): `.../mutex`, `.../block`

Inside `pprof`:
```
top            # hottest symbols
top -cum       # cumulative time (including callees)
list FuncName  # annotated source
web            # call graph (needs graphviz), or use -http=:0 on startup
```
