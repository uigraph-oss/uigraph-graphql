# uigraph-graphql — Claude Code Guidelines

Go module: `github.com/uigraph/graphql` · Go 1.25 · gqlgen v0.17 · open-source project

---

## Project layout

```
main.go                entry point (HTTP server, middleware wiring)
config/                env config + validation
middleware/            auth passthrough, request logging, CORS
client/                typed REST client for uigraph-api, one file per domain
graph/
  schema/              GraphQL SDL, one file per domain
  generated/, model/   gqlgen output — never hand-edit
  convert.go           pure REST-DTO -> GraphQL-model mapping, unit-tested
  *.resolvers.go       resolver implementations, one file per schema file
```

---

## The schema-driven file-split rule

gqlgen's `follow-schema` resolver layout ties each `<name>.resolvers.go` file to its source `<name>.graphqls` file. **Never** try to manually split a `*.resolvers.go` file — the next `go generate` undoes it. To split resolvers, split the schema file instead, then run `go generate ./graph/...`.

---

## Narrow interface pattern

`graph/resolver.go` holds a single `*client.Client` field. When refactoring toward narrow interfaces, declare one interface per domain (e.g., `authClient`, `orgClient`, `catalogClient`, …), each listing only the `*client.Client` methods that domain's resolvers actually call. `Resolver` will hold these interfaces instead of the concrete client — this is what lets tests inject a 2-3 method fake instead of mocking the whole REST client.

---

## Conversion functions live in `graph/convert.go`

Every `*ToModel` function maps a `client` DTO onto a `graph/model` GraphQL model. These are pure — no I/O, no context — specifically so they can be unit-tested without a running server. Add new ones there, not inline in resolver methods.

---

## Comments

Default: **no comments.** Well-named identifiers are self-documenting. Add a comment only when the **why** is non-obvious — a hidden constraint, a subtle invariant, a workaround for a specific bug.

---

## Forbidden patterns

- Storing a fat `*client.Client` directly on `Resolver` without narrow interfaces (eventual target; see Narrow interface pattern above).
- Hand-editing `graph/generated/generated.go` or `graph/model/models_gen.go` — these are gqlgen output.
- Manually splitting a `*.resolvers.go` file without first splitting its source `.graphqls` file.
- Using `log.Printf` for structured logging — use `slog` from the standard library.
- Forwarding a raw upstream error message straight to the GraphQL client without classification or sanitization.
- Using `transport.MultipartForm{}` in the GraphQL server setup without explicit validation.
- Skipping `go build ./...` before committing.

---

## How to run

```bash
# Generate gqlgen resolvers from schema changes
go generate ./graph/...

# Run tests with race detection
go test ./... -race

# Build the binary
go build ./...

# Run the server (requires UIGRAPH_API_BASE_URL env var)
./graphql
```

---
