# syntax=docker/dockerfile:1
ARG GO_VERSION=1.24.6

# ---- Builder ----
FROM golang:${GO_VERSION} AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GO111MODULE=on

# Tools & certs (Debian uses apt-get, not apk)
RUN apt-get update \
    && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Cache deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source and build
COPY cmd/apisvc ./cmd/apisvc
COPY pkg ./pkg
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/apisvc ./cmd/apisvc

# ---- Runtime ----
FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /out/apisvc /usr/local/bin/apisvc
EXPOSE 8080
ENTRYPOINT ["apisvc"]
