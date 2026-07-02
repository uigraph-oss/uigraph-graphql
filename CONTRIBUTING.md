# Contributing

## Adding a new domain

1. Add a new `.graphqls` file under `internal/graph/schema/` with its types, inputs, and `extend type Query`/`extend type Mutation` blocks.
2. Run `go generate ./internal/graph/...` to generate the matching `<name>.resolvers.go` stub.
3. Add the REST methods and DTO structs to a new file under `internal/uigraphapi/`.
4. Add the narrow client interface for the new domain in `internal/graph/resolver.go`, and wire it into `internal/server/server.go`'s `graph.Resolver{...}` construction.
5. Implement the resolver bodies, adding conversion functions to `internal/graph/convert/` for any new model type.
6. `go build ./...` must pass before committing.

## Code style

- No comments describing *what* code does — only *why*, when non-obvious.
- Resolvers hold narrow interfaces (declared in `internal/graph/resolver.go`), never the full `*uigraphapi.Client`.
- HTTP-facing files target under ~300 lines; split by sub-domain (mirroring how `internal/graph/schema/` is split) when a file grows past that.
- Pure conversion logic belongs in `internal/graph/convert/`, not inline in resolver methods — that's what makes it unit-testable without a running server.

## Pull requests

- `go build ./...`, `golangci-lint run`, and `go test ./... -race -cover` must all pass.
- If you changed a `.graphqls` file, commit the regenerated `generated.go`/`models_gen.go`/`*.resolvers.go` — CI checks this with a `go generate && git diff --exit-code` step.
