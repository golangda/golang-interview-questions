.PHONY: generate test tidy

generate:
	go generate ./...

test:
	go test ./... -race -coverprofile=cover.out
	@go tool cover -func=cover.out | tail -n +2 || true

tidy:
	go mod tidy
