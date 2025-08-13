# gRPC in Go — Beginner Guide (Ready-to-Run)

This project gives you a minimal **gRPC** service in **Go** with:
- Unary RPC: `SayHello`
- Server-streaming RPC: `GreetManyTimes`
- Logging + optional token auth via interceptors
- Makefile targets to generate protobuf code and run

## 1) Prerequisites

- **Go 1.21+**
- **Protocol Buffers compiler (`protoc`)**
  - macOS: `brew install protobuf`
  - Ubuntu/Debian: `sudo apt-get install -y protobuf-compiler`
  - Windows (choco): `choco install protoc`

## 2) One-time setup

```bash
make tidy        # pulls go deps
make tools       # installs protoc plugins into $GOPATH/bin
```

Make sure `$GOPATH/bin` is on your `PATH`:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 3) Generate gRPC code

```bash
make gen
```
This creates `api/hellopb/hello.pb.go` and `api/hellopb/hello_grpc.pb.go`.

> Tip: If you ever change `api/hello.proto`, re-run `make gen`.

## 4) Run the server

```bash
make run-server
```
Optional environment variables:
- `GRPC_ADDR` — listen address (default `:50051`)
- `GREETER_TOKEN` — if set, enables simple bearer-token auth (e.g., `s3cr3t`).

## 5) Run the client (in a new terminal)

```bash
make run-client
```
Optional environment variables:
- `GRPC_ADDR` — server address (default `localhost:50051`)
- `GREETER_TOKEN` — must match the server token if auth enabled.

Expected output:
```
Unary: Hello, Rahul! 👋
Stream:
  [1/5] Hello, Rahul!
  [2/5] Hello, Rahul!
  [3/5] Hello, Rahul!
  [4/5] Hello, Rahul!
  [5/5] Hello, Rahul!
```

## 6) How it works (high level)

- API contract lives in **`api/hello.proto`**
- `make gen` uses `protoc` + Go plugins to generate Go code in **`api/hellopb`**
- Server code implements the generated `GreeterServer` interface
- Client uses the generated `GreeterClient` to call methods
- Interceptors add logging and optional metadata-based auth
- Deadlines/cancellation handled via `context.Context`

## 7) Common fixes

- **`protoc-gen-go: program not found`** — run `make tools` and ensure `$GOPATH/bin` is in `PATH`.
- **`protoc: command not found`** — install protoc (see prerequisites).
- **`WithInsecure deprecated`** — we use `credentials/insecure` for local dev; use TLS in production.

## 8) Next steps

- Add client-streaming and bidirectional-streaming RPCs
- Add TLS/mTLS (`credentials.NewServerTLSFromFile` / `NewClientTLSFromFile`)
- Add OpenTelemetry for tracing + metrics
- Containerize (Docker) and deploy to Kubernetes with health checks