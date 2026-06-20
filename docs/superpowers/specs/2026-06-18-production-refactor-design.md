# uigraph-graphql: Production-Ready Refactor

**Date:** 2026-06-18
**Status:** Approved — pending implementation
**Goal:** Restructure this GraphQL BFF to industry-standard Go/GraphQL conventions, close real production gaps, and make the repo genuinely open-source-ready — without changing the GraphQL schema/contract.

---

## Context

`uigraph-graphql` is a thin GraphQL gateway (gqlgen) in front of the `uigraph-api` REST backend. It has no database and no business logic of its own beyond translating between REST DTOs and GraphQL models. The sibling repo `uigraph-api` went through an analogous folder-structure refactor (see its `docs/superpowers/specs/2026-06-18-folder-structure-design.md`); this document defines the equivalent target for this repo, adapted to what a GraphQL/gqlgen service actually needs (it has no store/postgres/cache layers — those don't apply here).

Current pain points:

1. No `cmd/`/`internal/` split — everything lives in flat top-level packages (`client/`, `config/`, `middleware/`, `graph/`), importable by anyone even though nothing should be imported outside this binary.
2. `graph/convert.go` (957 lines) and `client/types.go` (747 lines) are monolithic, mixing every domain's conversion/DTO code in one file each.
3. `graph/catalog.resolvers.go` (499 lines) and `graph/content.resolvers.go` (535 lines) mix multiple unrelated domains because their source schema files (`catalog.graphqls`, `content.graphqls`) do. gqlgen's `follow-schema` layout ties resolver file boundaries 1:1 to schema file names, so the resolver files can't be split cleanly without splitting the schema files first.
4. Zero tests, zero CI, no lint config, no README/LICENSE/CONTRIBUTING.
5. Real production gaps: no graceful shutdown, no HTTP server timeouts, `log.Printf` instead of structured logging, raw upstream error bodies forwarded verbatim as GraphQL error messages, no query complexity/depth limit, container runs as root.
6. A known N+1 call pattern (`resolveActor`/`resolveAssetURL` call a *batch* upstream endpoint one id at a time, once per row) — confirmed present despite the misleading "Cache actor and assets" commit title; no caching or batching actually exists.

---

## Design Principles

1. **`internal/` for everything not meant to be imported externally.** This is a single binary, not a library — matches `uigraph-api`'s convention and Go's own guidance.
2. **Split schema files, not resolver files.** gqlgen regenerates resolver stubs per schema file; splitting `.graphqls` files is the only regeneration-safe way to get smaller resolver files.
3. **Narrow client interfaces, not the fat concrete client.** The `Resolver` struct declares only the methods each domain needs, all satisfied by the same `*uigraphapi.Client` — enables tests without restructuring the client.
4. **No schema/contract changes.** Every change here is internal. `uigraph-ui` and any other consumer sees identical responses for valid queries.
5. **Errors are logged before sanitizing.** Every internal failure is logged server-side with `slog`; only safe, classified messages reach the GraphQL client.

---

## Target Directory Tree

```
uigraph-graphql/
├── cmd/
│   └── server/
│       └── main.go                 # entry point only: load config, build server, run
│
├── internal/
│   ├── config/
│   │   └── config.go               # unchanged content + boot-time validation
│   │
│   ├── middleware/
│   │   ├── auth.go                 # unchanged — header/cookie passthrough
│   │   └── logging.go              # NEW — request ID + structured access log
│   │
│   ├── uigraphapi/                 # renamed from client/ — REST client for uigraph-api
│   │   ├── client.go               # Client, do/get/post/put/del, APIError
│   │   ├── auth.go                 # Me, MyOrgs, SwitchOrg + MeResponse, OrgSummary
│   │   ├── org.go                  # org/member/team/invitation/service-account + structs
│   │   ├── users.go                # server-admin user methods + User struct
│   │   ├── sso.go                  # oauth/role-mapping/ldap/saml + structs
│   │   ├── folder.go               # folder methods + Folder struct
│   │   ├── diagram.go              # diagram/version/image methods + structs
│   │   ├── uimap.go                # map/frame/focalpoint/canvas/group/link + structs
│   │   ├── component.go            # flow-component & component palette + structs
│   │   ├── catalog.go              # service/api-group/doc/db/endpoint + structs
│   │   ├── testpack.go             # test pack/case/run/result + structs
│   │   ├── actors.go               # unchanged
│   │   └── assets.go               # unchanged
│   │
│   ├── server/
│   │   └── server.go               # NEW — *http.Server, mux wiring, graceful shutdown
│   │
│   └── graph/
│       ├── schema/
│       │   ├── schema.graphqls
│       │   ├── directives.graphqls
│       │   ├── actor.graphqls
│       │   ├── auth.graphqls
│       │   ├── admin.graphqls
│       │   ├── org.graphqls
│       │   ├── folder.graphqls     # SPLIT from content.graphqls
│       │   ├── diagram.graphqls    # SPLIT from content.graphqls
│       │   ├── uimap.graphqls      # SPLIT from content.graphqls
│       │   ├── catalog.graphqls    # trimmed: services/api-groups/docs/dbs/endpoints
│       │   └── testpack.graphqls   # SPLIT from catalog.graphqls
│       │
│       ├── generated/generated.go  # gqlgen output, unchanged mechanism
│       ├── model/models_gen.go     # gqlgen output, unchanged mechanism
│       ├── resolver.go             # Resolver struct — narrow client interfaces
│       ├── refs.go                 # resolveActor/resolveAssetURL (need r.Client)
│       ├── helpers.go              # unchanged map/json helpers
│       ├── auth.resolvers.go
│       ├── admin.resolvers.go
│       ├── org.resolvers.go
│       ├── folder.resolvers.go     # NEW after schema split
│       ├── diagram.resolvers.go    # NEW after schema split
│       ├── uimap.resolvers.go      # NEW after schema split
│       ├── catalog.resolvers.go    # shrinks
│       ├── testpack.resolvers.go   # NEW after schema split
│       │
│       └── convert/                # NEW — pure DTO→GraphQL-model mapping, unit-testable
│           ├── auth.go
│           ├── org.go
│           ├── admin.go
│           ├── folder.go
│           ├── diagram.go
│           ├── uimap.go
│           ├── catalog.go
│           └── testpack.go
│
├── docs/
│   └── superpowers/{specs,plans}/  # design history
│
├── .github/workflows/ci.yml        # NEW
├── .golangci.yml                   # NEW
├── Makefile                        # NEW
├── README.md                       # NEW
├── CONTRIBUTING.md                 # NEW
├── CLAUDE.md                       # NEW — contributor/AI conventions
├── gqlgen.yml                      # updated paths only
├── Dockerfile / Dockerfile.dev      # minor hardening (non-root user)
└── go.mod
```

---

## Schema-Driven Resolver Split

`content.graphqls` currently mixes four unrelated domains (Folder, Diagram, UIMap/Frame/FocalPoint/Canvas/FrameGroup/FrameLink, FocalPointMeta) in one 476-line file, which is why `content.resolvers.go` is 535 lines. `catalog.graphqls` mixes the Service/APIGroup/Endpoint/Doc/DB domain with the unrelated TestPack/TestCase/TestRun/TestRunResult domain in one 722-line file.

gqlgen's `follow-schema` resolver layout regenerates `<schema-file-name>.resolvers.go` per schema file on every `go generate` run. Manually splitting a resolver file without splitting its source schema file gets silently undone on next generation. The fix is to split the schema files themselves:

- `content.graphqls` → `folder.graphqls`, `diagram.graphqls`, `uimap.graphqls`
- `catalog.graphqls` → `catalog.graphqls` (trimmed), `testpack.graphqls`

gqlgen merges `extend type Query { ... }` / `extend type Mutation { ... }` blocks across files automatically (the existing files already rely on this), so this split is mechanical and safe. It also happens to mirror `uigraph-api`'s own domain boundaries (`folder`, `diagram`, `uimap`, `catalog`) almost exactly.

`admin.graphqls` (184 lines → 132-line resolver) and `org.graphqls` (165 lines → 235-line resolver) stay as single files — both are under the ~300-line target and each represents one coherent domain from the client's perspective (org administration; server admin + SSO).

---

## Narrow Client Interfaces

Today `Resolver` holds one field: `Client *client.Client` (the entire REST client, every domain). Per the same principle `uigraph-api` applies to its handlers, `internal/graph/resolver.go` instead declares one narrow interface per domain:

```go
type authClient interface {
    Me(ctx context.Context) (*uigraphapi.MeResponse, error)
    MyOrgs(ctx context.Context) ([]uigraphapi.OrgSummary, error)
    SwitchOrg(ctx context.Context, orgID string) error
}

type Resolver struct {
    Auth    authClient
    Org     orgClient
    Admin   adminClient
    Folder  folderClient
    Diagram diagramClient
    UIMap   uimapClient
    Catalog catalogClient
    TestPack testPackClient
    Actor   actorClient   // ResolveActors, ResolveAssetURLs — used by refs.go
}
```

All satisfied automatically by the same concrete `*uigraphapi.Client` in `main.go` — zero behavior change, but a test can now inject a 3-method fake instead of mocking the entire REST client.

---

## Production Hardening

| Area | Current | Change |
|---|---|---|
| Shutdown | `http.ListenAndServe`, no signal handling | `signal.NotifyContext` + `srv.Shutdown(ctx)` with deadline, in `internal/server` |
| Server timeouts | None (zero-value `http.Server`) | `ReadHeaderTimeout`/`ReadTimeout`/`WriteTimeout`/`IdleTimeout`, env-configurable with sane defaults |
| Logging | `log.Printf`, no access log, no request id | `log/slog` everywhere; `internal/middleware/logging.go` adds request id + structured access log |
| Error exposure | Raw `client.APIError.Body` becomes the literal GraphQL error message | gqlgen `SetErrorPresenter`: log real error via `slog`, return sanitized/classified message to client |
| Abuse resistance | No complexity/depth limit | gqlgen `extension.FixedComplexityLimit` + depth limit — defense-in-depth independent of the deferred dataloader fix |
| Config | No validation | Validate `APIBaseURL` parses as a URL at boot; fail fast |
| Health | `/healthz` always 200 | Add `/readyz` checking upstream `uigraph-api` reachability (liveness vs. readiness) |
| CORS | None | Optional, env-gated allow-list (`ALLOWED_ORIGINS`), default disabled (prod is same-origin behind Caddy per `uigraph-deploy`) |
| Container | Runs as root | Add unprivileged `USER` in `Dockerfile` |
| Dead code | `transport.MultipartForm{}` registered, no schema field uses file upload | Remove |

---

## Testing Strategy

Zero tests exist today. Priority order:

1. **`internal/graph/convert/`** — pure functions, table-driven tests, especially optional-pointer edge cases. Highest value/cost ratio; this is most of the current `convert.go`.
2. **`internal/uigraphapi/`** — `httptest.Server`-backed tests verifying request construction and response decoding, plus `APIError`/`IsNotFound`.
3. **`internal/middleware/`** — context propagation tests for `Auth`/`ApplyAuth`/`BearerToken`.
4. **A handful of resolver-level integration tests** — using the narrow interfaces, inject a hand-rolled fake and exercise 2-3 representative resolvers through the real generated executable schema.

CI (`.github/workflows/ci.yml`): `go build ./...`, `golangci-lint run`, `go test ./... -race -cover`, and a generated-code drift check (`go generate ./... && git diff --exit-code`).

---

## OSS Deliverables

- **README.md** — what this service is, architecture summary, local dev (`air` via `Dockerfile.dev` or `make run`), env vars table, `/playground` usage, schema-change workflow.
- **CONTRIBUTING.md** — new-domain checklist (schema file → `go generate` → resolver → convert mapper → client method), style notes, PR expectations.
- **CLAUDE.md** — package layout, narrow-interface pattern, schema-driven split rule, comment policy, forbidden patterns — mirrors `uigraph-api`'s `CLAUDE.md`.
- **.golangci.yml** — govet, staticcheck, errcheck, unused, gofmt/goimports, revive.
- **Makefile** — `run`, `build`, `test`, `lint`, `generate`, `docker-build`.
- **.github/workflows/ci.yml** — as above.
- **LICENSE** — deferred; add immediately before the repo goes public.

---

## Out of Scope (explicit follow-ups)

- **Actor/asset dataloader.** The N+1 pattern in `resolveActor`/`resolveAssetURL` is real and confirmed (each row triggers its own call to a batch endpoint), but fixing it is deferred to a dedicated follow-up: introduce a per-request dataloader (e.g. `vikstrous/dataloadgen`) that batches all such calls within one GraphQL request into a single upstream call.
- GraphQL subscriptions, persisted queries/APQ, rate limiting — not currently used/needed.
- Any change to `uigraph-deploy`'s Caddyfile (it doesn't route `/graphql` yet — a deploy-repo concern).

---

## Success Criteria

- [ ] `cmd/server/main.go` is the only file outside `internal/`.
- [ ] No resolver, convert, or client file exceeds ~300 lines.
- [ ] `Resolver` holds narrow per-domain interfaces, not a single fat client.
- [ ] `content.graphqls` and the TestPack portion of `catalog.graphqls` are split into focused schema files; `go generate` reproduces the same resolver file boundaries on every run.
- [ ] Every GraphQL-facing error path logs the real error via `slog` before returning a sanitized message.
- [ ] Graceful shutdown, server timeouts, and `/readyz` are in place.
- [ ] `internal/graph/convert/` and `internal/uigraphapi/` have meaningful test coverage; CI runs build, lint, test, and a generated-code drift check.
- [ ] README, CONTRIBUTING, CLAUDE.md, Makefile, `.golangci.yml`, and CI workflow exist.
- [ ] `go build ./...` passes; the GraphQL schema and all responses are unchanged for existing valid queries.
