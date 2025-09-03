.PHONY: run build lint test


run:
go run ./cmd/gateway


build:
go build -o bin/gateway ./cmd/gateway


test:
go test ./...


lint:
golangci-lint run ./...