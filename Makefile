.PHONY: gen build run test lint fmt


gen: ## generate protobufs
buf generate ./api/proto


build:
go build ./...


run: ## run local stack
docker compose up --build


test:
go test ./...


lint:
golangci-lint run ./...


fmt:
gofmt -s -w .