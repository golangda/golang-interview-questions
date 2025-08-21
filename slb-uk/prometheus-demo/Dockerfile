# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY app/go.mod ./
RUN go mod download
COPY app/*.go ./
RUN go build -o /out/prom-demo ./main.go

# Runtime
FROM alpine:3.20
EXPOSE 2112
COPY --from=builder /out/prom-demo /usr/local/bin/prom-demo
ENTRYPOINT ["/usr/local/bin/prom-demo"]