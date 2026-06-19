FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /uigraph-graphql ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates && \
    addgroup -S uigraph && adduser -S uigraph -G uigraph
COPY --from=builder /uigraph-graphql /usr/local/bin/uigraph-graphql
USER uigraph
EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/uigraph-graphql"]
