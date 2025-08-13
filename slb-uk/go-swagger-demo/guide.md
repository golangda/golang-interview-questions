# Go Swagger Demo Guide

This guide explains how to set up, run, and explore the Swagger 2.0 documentation for the Go web service included in this repository.

---

## 1. Prerequisites

Ensure you have the following installed on your machine:

- [Go](https://go.dev/dl/) v1.22+
- Git
- Swagger CLI tool (`swag`)

Install the `swag` CLI tool globally:

```bash
go install github.com/swaggo/swag/cmd/swag@v1
```

Make sure `$GOPATH/bin` is in your `PATH` so you can run `swag` directly.

---

## 2. Install Dependencies

In the project directory, run:

```bash
go mod tidy
```

This will download and install all required dependencies listed in `go.mod`.

---

## 3. Generate Swagger Documentation

Run:

```bash
swag init
```

This command scans for `// @` comments in your Go source files and generates the Swagger 2.0 specification (`docs/swagger.json` and `docs/docs.go`).

---

## 4. Run the Service

Start the API service with:

```bash
go run .
```

The service will run on:

```
http://localhost:8080
```

Swagger UI will be available at:

```
http://localhost:8080/swagger/index.html
```

---

## 5. Available Endpoints

### **Misc**
- **GET** `/v1/hello` — Returns a welcome message.

### **Messages**
- **GET** `/v1/messages` — Lists all messages.
- **GET** `/v1/message/{id}` — Fetches a single message by ID.
- **POST** `/v1/message` — Creates a new message.
- **PUT** `/v1/message/{id}` — Updates an existing message by ID.
- **DELETE** `/v1/message/{id}` — Deletes a message by ID.

---

## 6. Testing with `curl`

List messages:
```bash
curl -s http://localhost:8080/v1/messages | jq
```

Create a message:
```bash
curl -s -X POST http://localhost:8080/v1/message   -H 'Content-Type: application/json'   -d '{"message":"bonjour"}' | jq
```

Update a message:
```bash
curl -s -X PUT http://localhost:8080/v1/message/1   -H 'Content-Type: application/json'   -d '{"message":"updated text"}' | jq
```

Delete a message:
```bash
curl -s -X DELETE http://localhost:8080/v1/message/2
```

---

## 7. Customizing the API Docs

- Update the package-level Swagger annotations in `main.go` to change metadata (title, description, version).
- Add or modify endpoint annotations directly above handler functions to reflect request/response shapes.

Run `swag init` again after any changes to update the docs.

---

## 8. Troubleshooting

- **Swagger UI not loading**: Ensure the `docs` package is imported in `main.go` as `_ "example.com/go-swagger-demo/docs"`.
- **No operations detected**: Check that annotation comments (`// @...`) are immediately above the related function.
- **Wrong host or base path**: Adjust `@host` and `@BasePath` in the package-level annotations.

---

Enjoy building and documenting your APIs!
