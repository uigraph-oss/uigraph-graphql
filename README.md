# uigraph-graphql

A GraphQL BFF (backend-for-frontend) in front of [`uigraph-api`](https://github.com/uigraph-oss/uigraph-api). It translates GraphQL queries and mutations into REST calls, and REST DTOs into GraphQL models. It has no database and no business logic of its own beyond that translation.

## Architecture

```
cmd/server/main.go          entry point: load config, run the server
internal/
  config/                   env-based configuration + validation
  middleware/               auth header/cookie passthrough, request logging, CORS
  uigraphapi/                typed REST client for uigraph-api, one file per domain
  server/                   HTTP server wiring, graceful shutdown
  graph/
    schema/                 GraphQL SDL, one file per domain
    generated/, model/      gqlgen-generated code — do not edit by hand
    convert/                pure REST-DTO -> GraphQL-model mapping functions
    *.resolvers.go          gqlgen resolver implementations, one file per schema file
```

## Local development

```bash
go run ./cmd/server
```

The GraphQL Playground is available at `http://localhost:8090/playground` whenever `ENV` is not `prod`.

For hot reload during development, see `Dockerfile.dev` and `.air.toml` (uses [air](https://github.com/air-verse/air)).

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `API_BASE_URL` | `http://localhost:8080` | Base URL of the `uigraph-api` backend this server proxies to |
| `PORT` | `8090` | HTTP listen port |
| `ENV` | `local` | `local` \| `dev` \| `prod` — controls whether the Playground and GraphQL introspection are enabled |
| `ALLOWED_ORIGINS` | (empty) | Comma-separated CORS allow-list. Empty disables CORS handling entirely — the default, since production deployments are same-origin behind a reverse proxy |

## Changing the schema

1. Edit the relevant `.graphqls` file under `internal/graph/schema/` (or add a new file for a new domain).
2. Run `go generate ./internal/graph/...` — this regenerates `internal/graph/generated/generated.go` and `internal/graph/model/models_gen.go`, and creates or updates the matching `<name>.resolvers.go` stub.
3. Implement the resolver body. Add a converter in `internal/graph/convert/` for any new model type, and the corresponding REST method to `internal/uigraphapi/` if it doesn't exist yet.
4. Commit the regenerated files alongside your change — CI fails if they're out of sync with the schema.

## Testing

```bash
go test ./... -race -cover
```

## Known limitations

`resolveActor`/`resolveAssetURL` (in `internal/graph/refs.go`) resolve one id at a time per row instead of batching every such call within a single GraphQL request. A per-request dataloader is a planned follow-up — see `docs/superpowers/specs/2026-06-18-production-refactor-design.md`.
