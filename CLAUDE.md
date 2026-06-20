# uigraph-graphql — Claude Code Guidelines

Go module: `github.com/uigraph/graphql` · Go 1.25 · gqlgen v0.17 · open-source project

---

## Project layout

```
cmd/server/main.go         entry point — loads config, calls server.Run
internal/
  config/                  env config + boot-time URL validation
  middleware/              auth passthrough, request logging, CORS
  uigraphapi/              typed REST client for uigraph-api, one file per domain
  server/                  HTTP server assembly, graceful shutdown, /healthz + /readyz
  graph/
    schema/                GraphQL SDL, one file per domain
    generated/, model/     gqlgen output — never hand-edit
    convert/               pure REST-DTO → GraphQL-model mapping, unit-tested
    refs.go                resolveActor / resolveAssetURL (need Resolver, not pure)
    *.resolvers.go         resolver implementations, one file per schema file
    errors.go              GraphQL ErrorPresenter — logs + sanitizes upstream errors
```

---

## The schema-driven file-split rule

gqlgen's `follow-schema` resolver layout ties each `<name>.resolvers.go` file to its source `<name>.graphqls` file. **Never** try to manually split a `*.resolvers.go` file — the next `go generate` undoes it. To split resolvers, split the schema file instead, then run `go generate ./internal/graph/...`.

---

## Narrow interface pattern

`internal/graph/resolver.go` declares one unexported interface per domain (`authClient`, `orgClient`, `adminClient`, …). The `Resolver` struct holds these interfaces (not the concrete `*uigraphapi.Client`). All interfaces are satisfied by the same `*uigraphapi.Client` wired in `cmd/server/main.go`. This lets tests inject a 2-3 method fake instead of mocking the whole REST client.

---

## Conversion functions live in `internal/graph/convert/`

Every `*ToModel` function maps a `uigraphapi` DTO onto a `graph/model` GraphQL model. These are pure — no I/O, no context — specifically so they can be unit-tested without a running server. Add new ones in the matching domain file (e.g., `catalog.go`), not inline in resolver methods.

---

## Comments

Default: **no comments.** Well-named identifiers are self-documenting. Add a comment only when the **why** is non-obvious — a hidden constraint, a subtle invariant, a workaround for a specific bug.

---

## Forbidden patterns

- Storing a fat `*uigraphapi.Client` directly on `Resolver` — use narrow interfaces (see above).
- Hand-editing `internal/graph/generated/generated.go` or `internal/graph/model/models_gen.go` — these are gqlgen output.
- Manually splitting a `*.resolvers.go` file without first splitting its source `.graphqls` file.
- Using `log.Printf` for structured logging — use `slog` from the standard library.
- Forwarding a raw upstream error message straight to the GraphQL client without classification or sanitization.
- Using `transport.MultipartForm{}` in the GraphQL server setup — it is removed; no schema field uses file upload.
- Skipping `go build ./...` before committing.

---

## How to run

```bash
# Generate gqlgen resolvers from schema changes
go generate ./internal/graph/...

# Run tests with race detection
go test ./... -race

# Build the binary
go build ./...

# Run the server (requires API_BASE_URL env var, defaults to http://localhost:8080)
go run ./cmd/server
```
