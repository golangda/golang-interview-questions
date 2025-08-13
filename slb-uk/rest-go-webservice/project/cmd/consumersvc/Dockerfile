# syntax=docker/dockerfile:1
ARG GO_VERSION=1.24.6

FROM golang:${GO_VERSION} AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GO111MODULE=on
RUN apt-get update \
    && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY cmd/consumersvc ./cmd/consumersvc
COPY pkg ./pkg
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/consumersvc ./cmd/consumersvc

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /out/consumersvc /usr/local/bin/consumersvc
ENTRYPOINT ["consumersvc"]
