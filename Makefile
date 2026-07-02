.PHONY: run build test lint fmt fmt-check generate docker-build

run:
	go run ./cmd/server

build:
	go build -o bin/uigraph-graphql ./cmd/server

test:
	go test ./... -race -cover

fmt:
	gofmt -w .

fmt-check:
	@files=$$(gofmt -l .); \
	if [ -n "$$files" ]; then \
		echo "gofmt required on:"; \
		echo "$$files"; \
		exit 1; \
	fi

lint: fmt-check
	golangci-lint run

generate:
	go generate ./internal/graph/...

docker-build:
	docker build -t uigraph-graphql:local .
