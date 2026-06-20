.PHONY: run build test lint generate docker-build

run:
	go run ./cmd/server

build:
	go build -o bin/uigraph-graphql ./cmd/server

test:
	go test ./... -race -cover

lint:
	golangci-lint run

generate:
	go generate ./internal/graph/...

docker-build:
	docker build -t uigraph-graphql:local .
