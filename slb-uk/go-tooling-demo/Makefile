# Go Tooling Starter Makefile
APP_NAME ?= app
PKGS := ./...
GOLANGCI := $(GOPATH)/bin/golangci-lint

.PHONY: all fmt vet lint test race cover coverhtml build run tidy deps generate tools clean

all: fmt vet lint test build

fmt:
	@echo "==> go fmt"
	@go fmt $(PKGS)

vet:
	@echo "==> go vet"
	@go vet $(PKGS)

lint: tools
	@echo "==> golangci-lint"
	@$(GOLANGCI) run || true

test:
	@echo "==> go test"
	@go test $(PKGS) -cover

race:
	@echo "==> go test -race"
	@go test $(PKGS) -race

cover:
	@echo "==> coverage profile"
	@go test $(PKGS) -coverprofile=cover.out
	@go tool cover -func=cover.out

coverhtml: cover
	@echo "==> HTML coverage report at cover.html"
	@go tool cover -html=cover.out -o cover.html

build:
	@echo "==> building $(APP_NAME)"
	@go build -o $(APP_NAME) .

run: build
	@./$(APP_NAME)

tidy:
	@echo "==> go mod tidy"
	@go mod tidy

deps:
	@echo "==> listing modules"
	@go list -m all

generate:
	@echo "==> go generate"
	@go generate $(PKGS)

tools:
	@echo "==> installing tools"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

clean:
	@echo "==> clean"
	@rm -f $(APP_NAME) cover.out cover.html
