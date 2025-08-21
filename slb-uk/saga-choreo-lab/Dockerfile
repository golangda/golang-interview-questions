# Multi-stage build for all services (set build-arg CMD=<service>)
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /src

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY cmd ./cmd
COPY pkg ./pkg

ARG CMD
ENV CGO_ENABLED=0 GOOS=linux
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o /out/app ./cmd/${CMD}

FROM alpine:3.20
RUN adduser -D -H app && apk add --no-cache ca-certificates && update-ca-certificates
USER app
COPY --from=builder /out/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
