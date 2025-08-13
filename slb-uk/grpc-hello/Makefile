PROTO=api/hello.proto

.PHONY: tools gen tidy run-server run-client all
tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

gen:
	protoc --go_out=. --go-grpc_out=. $(PROTO)

tidy:
	go mod tidy

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client

all: tidy tools gen run-server