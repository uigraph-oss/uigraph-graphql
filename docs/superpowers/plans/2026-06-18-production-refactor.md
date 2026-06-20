# uigraph-graphql Production Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restructure `uigraph-graphql` into a `cmd/`+`internal/` layout with domain-split files, narrow client interfaces, GraphQL-specific production hardening, a starter test suite, and OSS hygiene — with zero GraphQL schema/contract changes.

**Architecture:** Mechanical package relocation first (config, middleware, client→uigraphapi, graph→internal/graph, main→cmd/server), then content splits driven by existing domain boundaries (client structs/methods, GraphQL schema files, convert.go), then new production code (server wiring, logging, error sanitization, complexity limits), then tests, then OSS docs/CI.

**Tech Stack:** Go 1.25, gqlgen v0.17.73, `log/slog`, `net/http` stdlib only (no new web framework).

## Global Constraints

- Module path stays `github.com/uigraph/graphql` — only internal directory structure changes.
- No GraphQL schema field, type, input, or response-shape changes anywhere in this plan. Every task must leave existing valid queries returning identical results.
- Every task must end with `go build ./...` passing (and `go vet ./...` where noted) before its commit.
- No comments describing *what* code does — only non-obvious *why* comments, matching existing repo style.
- The actor/asset N+1 dataloader fix is explicitly out of scope (see spec `docs/superpowers/specs/2026-06-18-production-refactor-design.md`) — do not attempt it as part of any task below.

---

## Task 1: Relocate `config/` and `middleware/` under `internal/`

**Files:**
- Create: `internal/config/config.go` (moved from `config/config.go`, unchanged content)
- Create: `internal/middleware/auth.go` (moved from `middleware/auth.go`, unchanged content)
- Modify: `main.go` (import paths only)
- Delete: `config/config.go`, `middleware/auth.go`

**Interfaces:**
- Consumes: nothing new.
- Produces: `internal/config.Load() *Config` (same signature as today — error-returning signature lands in Task 15), `internal/middleware.Auth`, `internal/middleware.ApplyAuth`, `internal/middleware.BearerToken` (same signatures as today).

- [ ] **Step 1: Move the directories**

```bash
mkdir -p internal
git mv config internal/config
git mv middleware internal/middleware
```

- [ ] **Step 2: Update the import path in `client/client.go`**

In `client/client.go`, change:

```go
	"github.com/uigraph/graphql/middleware"
```

to:

```go
	"github.com/uigraph/graphql/internal/middleware"
```

- [ ] **Step 3: Update the import paths in `main.go`**

In `main.go`, change:

```go
	"github.com/uigraph/graphql/client"
	"github.com/uigraph/graphql/config"
	"github.com/uigraph/graphql/graph"
	"github.com/uigraph/graphql/graph/generated"
	"github.com/uigraph/graphql/middleware"
```

to:

```go
	"github.com/uigraph/graphql/client"
	"github.com/uigraph/graphql/graph"
	"github.com/uigraph/graphql/graph/generated"
	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/middleware"
```

(`client` and `graph` import paths are untouched here — they move in Tasks 2–3.)

- [ ] **Step 4: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "refactor: move config and middleware packages under internal/"
```

---

## Task 2: Rename `client/` → `internal/uigraphapi/`

**Files:**
- Create: `internal/uigraphapi/{client,actors,assets,auth,catalog,content,org,sso,types,users}.go` (moved from `client/`, package renamed)
- Modify: `main.go`, every file under `graph/` that imports `client` (`graph/resolver.go`, `graph/convert.go`, and any `*.resolvers.go` referencing `client.` types)
- Delete: `client/` directory

**Interfaces:**
- Consumes: `internal/middleware.ApplyAuth` (from Task 1).
- Produces: package `uigraphapi` at `github.com/uigraph/graphql/internal/uigraphapi`, with `Client`, `New`, `APIError`, `IsNotFound`, and every existing exported method/type, identical to today's `client` package, just renamed.

- [ ] **Step 1: Move the directory**

```bash
git mv client internal/uigraphapi
```

- [ ] **Step 2: Rename the package declaration in every moved file**

Run:

```bash
sed -i '' 's/^package client$/package uigraphapi/' internal/uigraphapi/*.go
```

(On Linux, drop the `''` after `-i`.)

- [ ] **Step 3: Update every reference across the repo**

Run:

```bash
grep -rl '"github.com/uigraph/graphql/client"' --include="*.go" . | xargs sed -i '' 's#"github.com/uigraph/graphql/client"#"github.com/uigraph/graphql/internal/uigraphapi"#g'
grep -rl '\bclient\.' --include="*.go" graph main.go | xargs sed -i '' 's/\bclient\./uigraphapi./g'
```

This second command rewrites qualified references like `client.MeResponse` → `uigraphapi.MeResponse` and `client.IsNotFound` → `uigraphapi.IsNotFound` in every `graph/*.go` file and in `main.go`. It does not touch `internal/uigraphapi/*.go` itself (those files don't qualify their own package name).

- [ ] **Step 4: Verify no stray references remain**

Run: `grep -rn '"github.com/uigraph/graphql/client"\|[^a-zA-Z]client\.' --include="*.go" . | grep -v internal/uigraphapi`
Expected: no output (empty).

- [ ] **Step 5: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: rename client package to internal/uigraphapi"
```

---

## Task 3: Move `graph/` → `internal/graph/`, update gqlgen config

**Files:**
- Create: `internal/graph/` (moved from `graph/`)
- Modify: `gqlgen.yml`, `main.go`
- Delete: `graph/` directory

**Interfaces:**
- Consumes: nothing new.
- Produces: gqlgen now reads schema from `internal/graph/schema/*.graphqls` and writes `internal/graph/generated/generated.go` + `internal/graph/model/models_gen.go`. `main.go` imports `github.com/uigraph/graphql/internal/graph` and `.../internal/graph/generated`.

- [ ] **Step 1: Move the directory**

```bash
git mv graph internal/graph
```

- [ ] **Step 2: Update `gqlgen.yml`**

Replace its contents with:

```yaml
schema:
  - internal/graph/schema/*.graphqls

exec:
  filename: internal/graph/generated/generated.go
  package: generated

model:
  filename: internal/graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: internal/graph
  package: graph
  filename_template: "{name}.resolvers.go"

autobind: []

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Time:
    model: github.com/99designs/gqlgen/graphql.Time
```

- [ ] **Step 3: Update every import path that referenced `graph/`**

Run:

```bash
grep -rl '"github.com/uigraph/graphql/graph' --include="*.go" . | xargs sed -i '' 's#"github.com/uigraph/graphql/graph#"github.com/uigraph/graphql/internal/graph#g'
```

This rewrites `github.com/uigraph/graphql/graph/model` → `.../internal/graph/model`, `github.com/uigraph/graphql/graph/generated` → `.../internal/graph/generated`, and the bare `github.com/uigraph/graphql/graph` import in `main.go`.

- [ ] **Step 4: Regenerate**

Run: `go generate ./internal/graph/...`
Expected: exits 0. `internal/graph/generated/generated.go` and `internal/graph/model/models_gen.go` are rewritten in place with the corrected import paths; no other files change.

- [ ] **Step 5: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: move graph package under internal/, update gqlgen paths"
```

---

## Task 4: Move `main.go` → `cmd/server/main.go`

**Files:**
- Create: `cmd/server/main.go` (moved from `main.go`, unchanged content)
- Delete: `main.go`

**Interfaces:**
- Consumes: nothing new (Go import paths are independent of the importing file's location).
- Produces: the binary is now built via `go build ./cmd/server` / `go run ./cmd/server`.

- [ ] **Step 1: Move the file**

```bash
mkdir -p cmd/server
git mv main.go cmd/server/main.go
```

- [ ] **Step 2: Build and run-check**

Run: `go build ./cmd/server/...`
Expected: exits 0, no output.

Run: `go vet ./...`
Expected: exits 0, no output.

- [ ] **Step 3: Update `Dockerfile` and `Dockerfile.dev` build paths**

In `Dockerfile`, change:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o /uigraph-graphql .
```

to:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o /uigraph-graphql ./cmd/server
```

In `.air.toml`, change:

```toml
  cmd = "go build -o ./tmp/main ."
```

to:

```toml
  cmd = "go build -o ./tmp/main ./cmd/server"
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "refactor: move main.go to cmd/server/main.go"
```

---

## Task 5: Drain auth/org/users/sso structs out of `types.go`

**Files:**
- Modify: `internal/uigraphapi/auth.go`, `internal/uigraphapi/org.go`, `internal/uigraphapi/users.go`, `internal/uigraphapi/sso.go`, `internal/uigraphapi/types.go`

**Interfaces:**
- Consumes: nothing new.
- Produces: `MeResponse`, `OrgSummary` now live in `auth.go`; `Org`, `Member`, `Team`, `TeamMember`, `Invitation`, `ServiceAccount`, `ServiceAccountToken`, `CreatedToken` now live in `org.go`; `User` now lives in `users.go`; `OAuthProvider`, `RoleMapping`, `LDAPConfig`, `SAMLConfig` now live in `sso.go`. No method signatures change.

- [ ] **Step 1: Add structs to `internal/uigraphapi/auth.go`**

Insert this block immediately after the `import "context"` line:

```go

type MeResponse struct {
	UserID       string `json:"userId"`
	OrgID        string `json:"orgId"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Login        string `json:"login"`
	Kind         string `json:"kind"`
	Role         string `json:"role"`
	AuthProvider string `json:"authProvider"`
	AvatarURL    string `json:"avatarUrl,omitempty"`
}

type OrgSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}
```

- [ ] **Step 2: Add structs to `internal/uigraphapi/org.go`**

Change its import block from:

```go
import (
	"context"
	"fmt"
)
```

to:

```go
import (
	"context"
	"fmt"
	"time"
)
```

Then insert this block immediately after the import block, before the `// ── Orgs ──` comment:

```go
type Org struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Member struct {
	UserID    string    `json:"userId"`
	OrgID     string    `json:"orgId"`
	Role      string    `json:"role"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Team struct {
	ID         string    `json:"id"`
	OrgID      string    `json:"orgId"`
	Name       string    `json:"name"`
	Email      string    `json:"email,omitempty"`
	ExternalID string    `json:"externalId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TeamMember struct {
	TeamID     string    `json:"teamId"`
	UserID     string    `json:"userId"`
	Permission string    `json:"permission"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Invitation struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"orgId"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Code      string     `json:"code"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type ServiceAccount struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"orgId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Role        string    `json:"role"`
	Disabled    bool      `json:"disabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ServiceAccountToken struct {
	ID               string     `json:"id"`
	ServiceAccountID string     `json:"serviceAccountId"`
	Name             string     `json:"name"`
	Prefix           string     `json:"prefix"`
	ExpiresAt        *time.Time `json:"expiresAt,omitempty"`
	LastUsedAt       *time.Time `json:"lastUsedAt,omitempty"`
	Revoked          bool       `json:"revoked"`
	CreatedAt        time.Time  `json:"createdAt"`
}

type CreatedToken struct {
	ServiceAccountToken
	Token string `json:"token"`
}

```

- [ ] **Step 3: Add struct to `internal/uigraphapi/users.go`**

Change its import line from `import "context"` to:

```go
import (
	"context"
	"time"
)
```

Then insert this block immediately after the import block, before the `// ── Users (server admin) ──` comment:

```go
type User struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Name       string     `json:"name"`
	Login      string     `json:"login"`
	Disabled   bool       `json:"disabled"`
	Role       string     `json:"role"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

```

- [ ] **Step 4: Add structs to `internal/uigraphapi/sso.go`**

Change its import line from `import "context"` to:

```go
import (
	"context"
	"time"
)
```

Then insert this block immediately after the import block, before the `// ── OAuth providers ──` comment:

```go
type OAuthProvider struct {
	ID             string    `json:"id"`
	ProviderName   string    `json:"providerName"`
	Type           string    `json:"type"`
	DisplayName    string    `json:"displayName"`
	ClientID       string    `json:"clientId"`
	ClientSecret   string    `json:"clientSecret"`
	AuthURL        string    `json:"authUrl"`
	TokenURL       string    `json:"tokenUrl"`
	UserinfoURL    string    `json:"userinfoUrl"`
	APIURL         string    `json:"apiUrl"`
	Scopes         string    `json:"scopes"`
	AllowedDomains string    `json:"allowedDomains"`
	AllowSignUp    bool      `json:"allowSignUp"`
	EmailClaim     string    `json:"emailClaim"`
	NameClaim      string    `json:"nameClaim"`
	SubClaim       string    `json:"subClaim"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type RoleMapping struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	ClaimKey       string `json:"claimKey"`
	ClaimValue     string `json:"claimValue"`
	Role           string `json:"role"`
	Scope          string `json:"scope"`
	ResourceType   string `json:"resourceType"`
	ResourceID     string `json:"resourceId"`
}

type LDAPConfig struct {
	ID                string    `json:"id"`
	Host              string    `json:"host"`
	Port              int       `json:"port"`
	UseSSL            bool      `json:"useSsl"`
	StartTLS          bool      `json:"startTls"`
	SkipTLSVerify     bool      `json:"skipTlsVerify"`
	BindDN            string    `json:"bindDn"`
	BindPassword      string    `json:"bindPassword"`
	SearchBaseDN      string    `json:"searchBaseDn"`
	SearchFilter      string    `json:"searchFilter"`
	EmailAttribute    string    `json:"emailAttribute"`
	NameAttribute     string    `json:"nameAttribute"`
	UsernameAttribute string    `json:"usernameAttribute"`
	MemberOfAttribute string    `json:"memberOfAttribute"`
	AllowSignUp       bool      `json:"allowSignUp"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type SAMLConfig struct {
	ID              string    `json:"id"`
	IDPMetadataURL  string    `json:"idpMetadataUrl"`
	IDPMetadataXML  string    `json:"idpMetadataXml"`
	IDPEntityID     string    `json:"idpEntityId"`
	IDPSsoURL       string    `json:"idpSsoUrl"`
	IDPCert         string    `json:"idpCert"`
	SPEntityID      string    `json:"spEntityId"`
	SPCert          string    `json:"spCert"`
	SPKey           string    `json:"spKey"`
	SignRequests    bool      `json:"signRequests"`
	NameIDFormat    string    `json:"nameIdFormat"`
	EmailAttribute  string    `json:"emailAttribute"`
	NameAttribute   string    `json:"nameAttribute"`
	LoginAttribute  string    `json:"loginAttribute"`
	GroupsAttribute string    `json:"groupsAttribute"`
	AllowSignUp     bool      `json:"allowSignUp"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

```

- [ ] **Step 5: Remove the now-duplicated structs from `internal/uigraphapi/types.go`**

Delete the `MeResponse`, `OrgSummary`, `Org`, `Member`, `Team`, `TeamMember`, `Invitation`, `ServiceAccount`, `ServiceAccountToken`, `CreatedToken`, `User`, `OAuthProvider`, `RoleMapping`, `LDAPConfig`, `SAMLConfig` type definitions and their section comments (`// ── Auth ──`, `// ── Users (server admin) ──`, `// ── SSO ──`, `// ── Org ──`) from `types.go`. After this step, `types.go` should contain only the `// REST DTO types...` package comment, its imports, and the `// ── Content ──` / `// ── Catalog ──` sections (handled in Tasks 6–7).

- [ ] **Step 6: Build**

Run: `go build ./...`
Expected: exits 0, no output. If `time` or `fmt` are now unused/duplicated in any file, fix the import block — `go build` will name the exact file and line.

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "refactor: move auth/org/users/sso structs out of types.go into their domain files"
```

---

## Task 6: Split `content.go` into `folder.go`, `diagram.go`, `component.go`, `uimap.go`

**Files:**
- Create: `internal/uigraphapi/folder.go`, `internal/uigraphapi/diagram.go`, `internal/uigraphapi/component.go`, `internal/uigraphapi/uimap.go`
- Modify: `internal/uigraphapi/types.go` (remove the Content section)
- Delete: `internal/uigraphapi/content.go`

**Interfaces:**
- Consumes: `Client.get/post/put/del` from `internal/uigraphapi/client.go` (Task 2).
- Produces: same exported types and methods as today (`Folder`, `Diagram`, `DiagramImage`, `DiagramVersion`, `FlowDiagramComponentField`, `FlowDiagramComponent`, `FlowComponents`, `ComponentField`, `Component`, `Components`, `UIMap`, `Frame`, `FocalPoint`, `Canvas`, `FrameGroup`, `FrameLink`, `FocalPointMeta`, and every `List*/Get*/Create*/Update*/Delete*/Sync*/Restore*` method content.go had), just redistributed across 4 files. The dead, unused `rawJSON` helper from `content.go` is dropped (confirmed zero call sites via `grep -rn "rawJSON" .` before this task).

- [ ] **Step 1: Create `internal/uigraphapi/folder.go`**

```go
package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type Folder struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"orgId"`
	ParentID  *string    `json:"parentId,omitempty"`
	TeamID    *string    `json:"teamId,omitempty"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Order     float64    `json:"order"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (c *Client) ListFolders(ctx context.Context, orgID, folderType, parentID string) ([]Folder, error) {
	q := url.Values{}
	if folderType != "" {
		q.Set("type", folderType)
	}
	if parentID != "" {
		q.Set("parentId", parentID)
	}
	path := "/api/v1/orgs/" + orgID + "/folders"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Folders []Folder `json:"folders"`
	}
	return out.Folders, c.get(ctx, path, &out)
}

func (c *Client) GetFolder(ctx context.Context, orgID, id string) (*Folder, error) {
	var out Folder
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id), &out)
}

func (c *Client) CreateFolder(ctx context.Context, orgID string, body map[string]interface{}) (*Folder, error) {
	var out Folder
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/folders", body, &out)
}

func (c *Client) UpdateFolder(ctx context.Context, orgID, id string, body map[string]interface{}) (*Folder, error) {
	var out Folder
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id), body, &out)
}

func (c *Client) DeleteFolder(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id))
}
```

- [ ] **Step 2: Create `internal/uigraphapi/diagram.go`**

```go
package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type Diagram struct {
	ID                 string     `json:"id"`
	OrgID              string     `json:"orgId"`
	FolderID           *string    `json:"folderId,omitempty"`
	TeamID             *string    `json:"teamId,omitempty"`
	Name               string     `json:"name"`
	ContentKey         string     `json:"contentKey"`
	ContentHash        string     `json:"contentHash"`
	PreviewAssetID     *string    `json:"previewAssetId,omitempty"`
	PreviewContentHash *string    `json:"previewContentHash,omitempty"`
	Source             *string    `json:"source,omitempty"`
	CreatedBy          string     `json:"createdBy"`
	UpdatedBy          *string    `json:"updatedBy,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt,omitempty"`
}

type DiagramImage struct {
	DiagramImageID string    `json:"diagramImageId"`
	DiagramID      string    `json:"diagramId"`
	OrgID          string    `json:"orgId"`
	AssetID        string    `json:"assetId"`
	FileName       *string   `json:"fileName,omitempty"`
	Order          int       `json:"order"`
	CreatedBy      string    `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
}

type DiagramVersion struct {
	ID            string    `json:"id"`
	DiagramID     string    `json:"diagramId"`
	VersionNumber int       `json:"versionNumber"`
	Label         *string   `json:"label,omitempty"`
	ContentKey    string    `json:"contentKey"`
	ContentHash   string    `json:"contentHash"`
	IsAutoVersion bool      `json:"isAutoVersion"`
	Source        *string   `json:"source,omitempty"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (c *Client) ListDiagrams(ctx context.Context, orgID, folderID string) ([]Diagram, error) {
	path := "/api/v1/orgs/" + orgID + "/diagrams"
	if folderID != "" {
		path += "?folderId=" + folderID
	}
	var out struct {
		Diagrams []Diagram `json:"diagrams"`
	}
	return out.Diagrams, c.get(ctx, path, &out)
}

func (c *Client) GetDiagram(ctx context.Context, orgID, id string) (*Diagram, error) {
	var out Diagram
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s", orgID, id), &out)
}

func (c *Client) GetDiagramContent(ctx context.Context, orgID, id string) (string, error) {
	var out struct {
		DiagramID string `json:"diagramId"`
		Content   string `json:"content"`
	}
	return out.Content, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/content", orgID, id), &out)
}

func (c *Client) CreateDiagram(ctx context.Context, orgID string, body map[string]interface{}) (*Diagram, error) {
	var out Diagram
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/diagrams", body, &out)
}

func (c *Client) UpdateDiagram(ctx context.Context, orgID, id string, body map[string]interface{}) (*Diagram, error) {
	var out Diagram
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s", orgID, id), body, &out)
}

func (c *Client) DeleteDiagram(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s", orgID, id))
}

func (c *Client) ListDiagramImages(ctx context.Context, orgID, diagramID string) ([]DiagramImage, error) {
	var out struct {
		Images []DiagramImage `json:"images"`
	}
	return out.Images, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/images", orgID, diagramID), &out)
}

func (c *Client) SyncDiagram(ctx context.Context, orgID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, "/api/v1/orgs/"+orgID+"/diagrams/sync", body, &out)
}

func (c *Client) ListDiagramVersions(ctx context.Context, orgID, diagramID string) ([]DiagramVersion, error) {
	var out struct {
		Versions []DiagramVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions", orgID, diagramID), &out)
}

func (c *Client) CreateDiagramVersion(ctx context.Context, orgID, diagramID string, body map[string]interface{}) (*DiagramVersion, error) {
	var out DiagramVersion
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions", orgID, diagramID), body, &out)
}

func (c *Client) GetDiagramVersionContent(ctx context.Context, orgID, diagramID, versionID string) (string, error) {
	var out struct {
		VersionID string `json:"versionId"`
		Content   string `json:"content"`
	}
	return out.Content, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions/%s/content", orgID, diagramID, versionID), &out)
}

func (c *Client) RestoreDiagramVersion(ctx context.Context, orgID, diagramID, versionID string) (*Diagram, error) {
	var out Diagram
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions/%s/restore", orgID, diagramID, versionID), nil, &out)
}
```

- [ ] **Step 3: Create `internal/uigraphapi/component.go`**

```go
package uigraphapi

import (
	"context"
	"fmt"
)

type FlowDiagramComponentField struct {
	FlowDiagramComponentFieldID string   `json:"flowDiagramComponentFieldId"`
	Label                       string   `json:"label"`
	Type                        string   `json:"type"`
	Required                    bool     `json:"required"`
	Readonly                    *bool    `json:"readonly,omitempty"`
	Options                     []string `json:"options,omitempty"`
	Order                       int      `json:"order"`
}

type FlowDiagramComponent struct {
	ComponentID                string                      `json:"componentId"`
	Type                       string                      `json:"type"`
	Name                       string                      `json:"name"`
	Description                string                      `json:"description"`
	Category                   string                      `json:"category"`
	Tags                       []string                    `json:"tags"`
	Slug                       string                      `json:"slug"`
	PreviewImageJpg            string                      `json:"previewImageJpg"`
	IsActive                   bool                        `json:"isActive"`
	Order                      int                         `json:"order"`
	OrganizationID             *string                     `json:"organizationId,omitempty"`
	FlowDiagramComponentFields []FlowDiagramComponentField `json:"flowDiagramComponentFields"`
}

type FlowComponents struct {
	Components       []FlowDiagramComponent `json:"components"`
	CustomComponents []FlowDiagramComponent `json:"customComponents"`
}

type ComponentField struct {
	ComponentFieldID string   `json:"componentFieldId"`
	Label            string   `json:"label"`
	Type             string   `json:"type"`
	Required         bool     `json:"required"`
	Readonly         *bool    `json:"readonly,omitempty"`
	Options          []string `json:"options,omitempty"`
	Order            int      `json:"order"`
}

type Component struct {
	ComponentID     string           `json:"componentId"`
	Type            string           `json:"type"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Category        string           `json:"category"`
	Tags            []string         `json:"tags"`
	Slug            string           `json:"slug"`
	PreviewImageJpg string           `json:"previewImageJpg"`
	IsActive        bool             `json:"isActive"`
	Order           int              `json:"order"`
	ComponentFields []ComponentField `json:"componentFields"`
}

type Components struct {
	Components       []Component `json:"components"`
	CustomComponents []Component `json:"customComponents"`
}

func (c *Client) ListFlowDiagramComponents(ctx context.Context, orgID string) (*FlowComponents, error) {
	var out FlowComponents
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/flow-diagram-components", orgID), &out)
}

func (c *Client) ListComponents(ctx context.Context, orgID string) (*Components, error) {
	var out Components
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/components", orgID), &out)
}
```

- [ ] **Step 4: Create `internal/uigraphapi/uimap.go`**

```go
package uigraphapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type UIMap struct {
	ID          string     `json:"id"`
	OrgID       string     `json:"orgId"`
	FolderID    *string    `json:"folderId,omitempty"`
	TeamID      *string    `json:"teamId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type Frame struct {
	ID                    string     `json:"id"`
	MapID                 string     `json:"mapId"`
	OrgID                 string     `json:"orgId"`
	ParentFrameID         *string    `json:"parentFrameId,omitempty"`
	Name                  string     `json:"name"`
	Description           string     `json:"description"`
	TemplateType          string     `json:"templateType"`
	ScreenshotAssetID     *string    `json:"screenshotAssetId,omitempty"`
	ScreenshotContentHash *string    `json:"screenshotContentHash,omitempty"`
	Status                string     `json:"status"`
	Order                 float64    `json:"order"`
	Source                *string    `json:"source,omitempty"`
	CreatedBy             string     `json:"createdBy"`
	UpdatedBy             *string    `json:"updatedBy,omitempty"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
	DeletedAt             *time.Time `json:"deletedAt,omitempty"`
}

type FocalPoint struct {
	ID         string     `json:"id"`
	FrameID    string     `json:"frameId"`
	OrgID      string     `json:"orgId"`
	Name       string     `json:"name"`
	LocationX  float64    `json:"locationX"`
	LocationY  float64    `json:"locationY"`
	Visibility string     `json:"visibility"`
	IsActive   bool       `json:"isActive"`
	CreatedBy  string     `json:"createdBy"`
	UpdatedBy  *string    `json:"updatedBy,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

type Canvas struct {
	MapID          string          `json:"mapId"`
	OrgID          string          `json:"orgId"`
	Zoom           float64         `json:"zoom"`
	NavigationX    float64         `json:"navigationX"`
	NavigationY    float64         `json:"navigationY"`
	FramePositions json.RawMessage `json:"framePositions"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type FrameGroup struct {
	ID          string     `json:"id"`
	FrameID     string     `json:"frameId"`
	OrgID       string     `json:"orgId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	LocationX   float64    `json:"locationX"`
	LocationY   float64    `json:"locationY"`
	Width       float64    `json:"width"`
	Height      float64    `json:"height"`
	Order       float64    `json:"order"`
	IsActive    bool       `json:"isActive"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type FrameLink struct {
	ID            string     `json:"id"`
	FrameID       string     `json:"frameId"`
	OrgID         string     `json:"orgId"`
	Kind          string     `json:"kind"`
	TargetFrameID *string    `json:"targetFrameId,omitempty"`
	TargetMapID   *string    `json:"targetMapId,omitempty"`
	Label         string     `json:"label"`
	LocationX     float64    `json:"locationX"`
	LocationY     float64    `json:"locationY"`
	IsActive      bool       `json:"isActive"`
	CreatedBy     string     `json:"createdBy"`
	UpdatedBy     *string    `json:"updatedBy,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

type FocalPointMeta struct {
	ID                   string          `json:"id"`
	FocalPointID         string          `json:"focalPointId"`
	OrgID                string          `json:"orgId"`
	FrameID              string          `json:"frameId"`
	ComponentID          string          `json:"componentId"`
	ComponentLinkID      *string         `json:"componentLinkId,omitempty"`
	ComponentImages      json.RawMessage `json:"componentImages"`
	ComponentFlowDiagram *string         `json:"componentFlowDiagram,omitempty"`
	ComponentModalFields json.RawMessage `json:"componentModalFields"`
	CreatedBy            string          `json:"createdBy"`
	UpdatedBy            *string         `json:"updatedBy,omitempty"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
	DeletedAt            *time.Time      `json:"deletedAt,omitempty"`
}

func (c *Client) ListMaps(ctx context.Context, orgID, folderID string) ([]UIMap, error) {
	path := "/api/v1/orgs/" + orgID + "/maps"
	if folderID != "" {
		path += "?folderId=" + folderID
	}
	var out struct {
		Maps []UIMap `json:"maps"`
	}
	return out.Maps, c.get(ctx, path, &out)
}

func (c *Client) GetMap(ctx context.Context, orgID, id string) (*UIMap, error) {
	var out UIMap
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s", orgID, id), &out)
}

func (c *Client) CreateMap(ctx context.Context, orgID string, body map[string]interface{}) (*UIMap, error) {
	var out UIMap
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/maps", body, &out)
}

func (c *Client) UpdateMap(ctx context.Context, orgID, id string, body map[string]interface{}) (*UIMap, error) {
	var out UIMap
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s", orgID, id), body, &out)
}

func (c *Client) DeleteMap(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s", orgID, id))
}

func (c *Client) ListFrames(ctx context.Context, orgID, mapID string) ([]Frame, error) {
	var out struct {
		Frames []Frame `json:"frames"`
	}
	return out.Frames, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames", orgID, mapID), &out)
}

func (c *Client) GetFrame(ctx context.Context, orgID, mapID, id string) (*Frame, error) {
	var out Frame
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s", orgID, mapID, id), &out)
}

func (c *Client) GetFrameByID(ctx context.Context, orgID, id string) (*Frame, error) {
	var out Frame
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/frames/%s", orgID, id), &out)
}

func (c *Client) CreateFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*Frame, error) {
	var out Frame
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames", orgID, mapID), body, &out)
}

func (c *Client) UpdateFrame(ctx context.Context, orgID, mapID, id string, body map[string]interface{}) (*Frame, error) {
	var out Frame
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s", orgID, mapID, id), body, &out)
}

func (c *Client) DeleteFrame(ctx context.Context, orgID, mapID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s", orgID, mapID, id))
}

func (c *Client) SyncFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/sync", orgID, mapID), body, &out)
}

func (c *Client) ListFocalPoints(ctx context.Context, orgID, mapID, frameID string) ([]FocalPoint, error) {
	var out struct {
		FocalPoints []FocalPoint `json:"focalPoints"`
	}
	return out.FocalPoints, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points", orgID, mapID, frameID), &out)
}

func (c *Client) GetFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) (*FocalPoint, error) {
	var out FocalPoint
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s", orgID, mapID, frameID, id), &out)
}

func (c *Client) CreateFocalPoint(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*FocalPoint, error) {
	var out FocalPoint
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points", orgID, mapID, frameID), body, &out)
}

func (c *Client) UpdateFocalPoint(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*FocalPoint, error) {
	var out FocalPoint
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s", orgID, mapID, frameID, id), body, &out)
}

func (c *Client) DeleteFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s", orgID, mapID, frameID, id))
}

func (c *Client) GetCanvas(ctx context.Context, orgID, mapID string) (*Canvas, error) {
	var out Canvas
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/canvas", orgID, mapID), &out)
}

func (c *Client) UpsertCanvas(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*Canvas, error) {
	var out Canvas
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/canvas", orgID, mapID), body, &out)
}

func (c *Client) ListFrameGroups(ctx context.Context, orgID, mapID, frameID string) ([]FrameGroup, error) {
	var out struct {
		Groups []FrameGroup `json:"groups"`
	}
	return out.Groups, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups", orgID, mapID, frameID), &out)
}

func (c *Client) CreateFrameGroup(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*FrameGroup, error) {
	var out FrameGroup
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups", orgID, mapID, frameID), body, &out)
}

func (c *Client) UpdateFrameGroup(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*FrameGroup, error) {
	var out FrameGroup
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups/%s", orgID, mapID, frameID, id), body, &out)
}

func (c *Client) DeleteFrameGroup(ctx context.Context, orgID, mapID, frameID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups/%s", orgID, mapID, frameID, id))
}

func (c *Client) ListFrameLinks(ctx context.Context, orgID, mapID, frameID string) ([]FrameLink, error) {
	var out struct {
		Links []FrameLink `json:"links"`
	}
	return out.Links, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links", orgID, mapID, frameID), &out)
}

func (c *Client) CreateFrameLink(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*FrameLink, error) {
	var out FrameLink
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links", orgID, mapID, frameID), body, &out)
}

func (c *Client) UpdateFrameLink(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*FrameLink, error) {
	var out FrameLink
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links/%s", orgID, mapID, frameID, id), body, &out)
}

func (c *Client) DeleteFrameLink(ctx context.Context, orgID, mapID, frameID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links/%s", orgID, mapID, frameID, id))
}

func (c *Client) ListFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string) ([]FocalPointMeta, error) {
	var out struct {
		Meta []FocalPointMeta `json:"meta"`
	}
	return out.Meta, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta", orgID, mapID, frameID, fpID), &out)
}

func (c *Client) CreateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string, body map[string]interface{}) (*FocalPointMeta, error) {
	var out FocalPointMeta
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta", orgID, mapID, frameID, fpID), body, &out)
}

func (c *Client) UpdateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string, body map[string]interface{}) (*FocalPointMeta, error) {
	var out FocalPointMeta
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta/%s", orgID, mapID, frameID, fpID, id), body, &out)
}

func (c *Client) DeleteFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta/%s", orgID, mapID, frameID, fpID, id))
}
```

- [ ] **Step 5: Delete `content.go` and the moved structs from `types.go`**

```bash
rm internal/uigraphapi/content.go
```

From `internal/uigraphapi/types.go`, delete the `Folder`, `Diagram`, `FlowDiagramComponentField`, `FlowDiagramComponent`, `FlowComponents`, `ComponentField`, `Component`, `Components`, `DiagramImage`, `DiagramVersion`, `UIMap`, `Frame`, `FocalPoint`, `Canvas`, `FrameGroup`, `FrameLink`, `FocalPointMeta` type definitions and the `// ── Content ──` section comment. After this step, `types.go` should contain only the `// ── Catalog ──` section (handled in Task 7).

- [ ] **Step 6: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "refactor: split content.go into folder/diagram/component/uimap files"
```

---

## Task 7: Split `catalog.go` into `catalog.go` (trimmed) and `testpack.go`; delete `types.go`

**Files:**
- Modify: `internal/uigraphapi/catalog.go` (rewrite — trimmed to Service/APIGroup/Doc/Diagram-link/DB/Endpoint only)
- Create: `internal/uigraphapi/testpack.go`
- Delete: `internal/uigraphapi/types.go`

**Interfaces:**
- Consumes: `Client.get/post/put/del`, and `Diagram` type from `internal/uigraphapi/diagram.go` (Task 6) — `ServiceDiagram.Diagram *Diagram`.
- Produces: same exported types/methods as today, just redistributed: `Service`, `ServiceStats`, `APIGroup`, `APIGroupVersion`, `ServiceDoc`, `ServiceDiagram`, `ServiceDB`, `ServiceDBVersion`, `APIEndpoint` + their methods stay in `catalog.go`; `KeyValue`, `Assertion`, `AuthConfig`, `TestCaseStep`, `ManualTestCase`, `APITestCase`, `GraphQLTestCase`, `DatabaseTestCase`, `GRPCTestCase`, `TestPack`, `TestCase`, `TestRun`, `TestRunSummary`, `TestRunResult` + their methods move to `testpack.go`.

- [ ] **Step 1: Rewrite `internal/uigraphapi/catalog.go`**

```go
package uigraphapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Service struct {
	ID              string          `json:"id"`
	OrgID           string          `json:"orgId"`
	FolderID        *string         `json:"folderId,omitempty"`
	TeamID          *string         `json:"teamId,omitempty"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	Description     string          `json:"description"`
	Status          string          `json:"status"`
	Tier            string          `json:"tier"`
	Category        string          `json:"category"`
	Language        string          `json:"language"`
	GitRepoURL      *string         `json:"gitRepoUrl,omitempty"`
	JiraProjectURL  *string         `json:"jiraProjectUrl,omitempty"`
	SlackChannelURL *string         `json:"slackChannelUrl,omitempty"`
	LastCommitSha   *string         `json:"lastCommitSha,omitempty"`
	Labels          []string        `json:"labels"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CreatedBy       string          `json:"createdBy"`
	UpdatedBy       *string         `json:"updatedBy,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type ServiceStats struct {
	ServiceID     string `json:"serviceId"`
	EndpointCount int    `json:"endpointCount"`
	DiagramCount  int    `json:"diagramCount"`
	DocCount      int    `json:"docCount"`
	DBTableCount  int    `json:"dbTableCount"`
	TestCaseCount int    `json:"testCaseCount"`
}

type APIGroup struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"serviceId"`
	OrgID     string    `json:"orgId"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Label     *string   `json:"label,omitempty"`
	Protocol  string    `json:"protocol"`
	SpecKey   *string   `json:"specKey,omitempty"`
	SpecHash  *string   `json:"specHash,omitempty"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy *string   `json:"updatedBy,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type APIGroupVersion struct {
	ID            string    `json:"id"`
	APIGroupID    string    `json:"apiGroupId"`
	VersionNumber int       `json:"versionNumber"`
	Label         *string   `json:"label,omitempty"`
	SpecKey       string    `json:"specKey"`
	SpecHash      string    `json:"specHash"`
	IsAutoVersion bool      `json:"isAutoVersion"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

type ServiceDoc struct {
	ID          string    `json:"id"`
	ServiceID   string    `json:"serviceId"`
	OrgID       string    `json:"orgId"`
	FileKey     string    `json:"fileKey"`
	FileName    string    `json:"fileName"`
	FileType    string    `json:"fileType"`
	Description string    `json:"description"`
	ContentHash string    `json:"contentHash"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedBy   *string   `json:"updatedBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ServiceDiagram struct {
	ServiceID string    `json:"serviceId"`
	DiagramID string    `json:"diagramId"`
	OrgID     string    `json:"orgId"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy *string   `json:"updatedBy,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Diagram   *Diagram  `json:"diagram,omitempty"`
}

type ServiceDB struct {
	ID         string          `json:"id"`
	ServiceID  string          `json:"serviceId"`
	OrgID      string          `json:"orgId"`
	DBName     string          `json:"dbName"`
	DBType     string          `json:"dbType"`
	Dialect    string          `json:"dialect"`
	SchemaJSON json.RawMessage `json:"schemaJson"`
	Source     *string         `json:"source,omitempty"`
	SourceTS   *time.Time      `json:"sourceTs,omitempty"`
	CreatedBy  string          `json:"createdBy"`
	UpdatedBy  *string         `json:"updatedBy,omitempty"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type ServiceDBVersion struct {
	ID            string          `json:"id"`
	ServiceDBID   string          `json:"serviceDbId"`
	VersionNumber int             `json:"versionNumber"`
	Label         *string         `json:"label,omitempty"`
	SchemaJSON    json.RawMessage `json:"schemaJson"`
	Source        *string         `json:"source,omitempty"`
	SourceTS      *time.Time      `json:"sourceTs,omitempty"`
	IsAutoVersion bool            `json:"isAutoVersion"`
	CreatedBy     string          `json:"createdBy"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type APIEndpoint struct {
	ID          string          `json:"id"`
	APIGroupID  string          `json:"apiGroupId"`
	ServiceID   string          `json:"serviceId"`
	OrgID       string          `json:"orgId"`
	OperationID string          `json:"operationId"`
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	Summary     string          `json:"summary"`
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	Parameters  json.RawMessage `json:"parameters"`
	RequestBody json.RawMessage `json:"requestBody"`
	Responses   json.RawMessage `json:"responses"`
	Order       float64         `json:"order"`
	CreatedBy   string          `json:"createdBy"`
	UpdatedBy   *string         `json:"updatedBy,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func (c *Client) ListServices(ctx context.Context, orgID, folderID, teamID string) ([]Service, error) {
	path := "/api/v1/orgs/" + orgID + "/services"
	q := url.Values{}
	if folderID != "" {
		q.Set("folderId", folderID)
	}
	if teamID != "" {
		q.Set("teamId", teamID)
	}
	if enc := q.Encode(); enc != "" {
		path += "?" + enc
	}
	var out struct {
		Services []Service `json:"services"`
	}
	return out.Services, c.get(ctx, path, &out)
}

func (c *Client) GetService(ctx context.Context, orgID, id string) (*Service, error) {
	var out Service
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id), &out)
}

func (c *Client) CreateService(ctx context.Context, orgID string, body map[string]interface{}) (*Service, error) {
	var out Service
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/services", body, &out)
}

func (c *Client) UpdateService(ctx context.Context, orgID, id string, body map[string]interface{}) (*Service, error) {
	var out Service
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id), body, &out)
}

func (c *Client) DeleteService(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id))
}

func (c *Client) ListServiceStats(ctx context.Context, orgID string, serviceID *string) ([]ServiceStats, error) {
	path := "/api/v1/orgs/" + orgID + "/services/stats"
	if serviceID != nil && *serviceID != "" {
		q := url.Values{}
		q.Set("serviceId", *serviceID)
		path += "?" + q.Encode()
	}
	var out struct {
		Stats []ServiceStats `json:"stats"`
	}
	return out.Stats, c.get(ctx, path, &out)
}

func (c *Client) ListAPIGroups(ctx context.Context, orgID, serviceID string) ([]APIGroup, error) {
	var out struct {
		APIGroups []APIGroup `json:"apiGroups"`
	}
	return out.APIGroups, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups", orgID, serviceID), &out)
}

func (c *Client) GetAPIGroup(ctx context.Context, orgID, serviceID, id string) (*APIGroup, error) {
	var out APIGroup
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*APIGroup, error) {
	var out APIGroup
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups", orgID, serviceID), body, &out)
}

func (c *Client) UpdateAPIGroup(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*APIGroup, error) {
	var out APIGroup
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteAPIGroup(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id))
}

func (c *Client) SyncAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/sync", orgID, serviceID), body, &out)
}

func (c *Client) ListAPIGroupVersions(ctx context.Context, orgID, serviceID, apiGroupID string) ([]APIGroupVersion, error) {
	var out struct {
		Versions []APIGroupVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/versions", orgID, serviceID, apiGroupID), &out)
}

func (c *Client) ListServiceDocs(ctx context.Context, orgID, serviceID string) ([]ServiceDoc, error) {
	var out struct {
		Docs []ServiceDoc `json:"docs"`
	}
	return out.Docs, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs", orgID, serviceID), &out)
}

func (c *Client) GetServiceDoc(ctx context.Context, orgID, serviceID, id string) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateServiceDoc(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs", orgID, serviceID), body, &out)
}

func (c *Client) UpdateServiceDoc(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteServiceDoc(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id))
}

func (c *Client) ListServiceDiagrams(ctx context.Context, orgID, serviceID string) ([]ServiceDiagram, error) {
	var out struct {
		Diagrams []ServiceDiagram `json:"diagrams"`
	}
	return out.Diagrams, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams", orgID, serviceID), &out)
}

func (c *Client) CreateServiceDiagram(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDiagram, error) {
	var out ServiceDiagram
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams", orgID, serviceID), body, &out)
}

func (c *Client) DeleteServiceDiagram(ctx context.Context, orgID, serviceID, diagramID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams/%s", orgID, serviceID, diagramID))
}

func (c *Client) ListServiceDBs(ctx context.Context, orgID, serviceID string) ([]ServiceDB, error) {
	var out struct {
		DBs []ServiceDB `json:"dbs"`
	}
	return out.DBs, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs", orgID, serviceID), &out)
}

func (c *Client) GetServiceDB(ctx context.Context, orgID, serviceID, id string) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateServiceDB(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs", orgID, serviceID), body, &out)
}

func (c *Client) UpdateServiceDB(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteServiceDB(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id))
}

func (c *Client) ListServiceDBVersions(ctx context.Context, orgID, serviceID, serviceDBID string) ([]ServiceDBVersion, error) {
	var out struct {
		Versions []ServiceDBVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions", orgID, serviceID, serviceDBID), &out)
}

func (c *Client) CreateServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID string, body map[string]interface{}) (*ServiceDBVersion, error) {
	var out ServiceDBVersion
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions", orgID, serviceID, serviceDBID), body, &out)
}

func (c *Client) RestoreServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID, versionID string) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions/%s/restore", orgID, serviceID, serviceDBID, versionID), nil, &out)
}

func (c *Client) ListAPIEndpoints(ctx context.Context, orgID, serviceID, apiGroupID string) ([]APIEndpoint, error) {
	var out struct {
		Endpoints []APIEndpoint `json:"endpoints"`
	}
	return out.Endpoints, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints", orgID, serviceID, apiGroupID), &out)
}

func (c *Client) GetAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id), &out)
}

func (c *Client) CreateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID string, body map[string]interface{}) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints", orgID, serviceID, apiGroupID), body, &out)
}

func (c *Client) UpdateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string, body map[string]interface{}) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id), body, &out)
}

func (c *Client) DeleteAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id))
}
```

- [ ] **Step 2: Create `internal/uigraphapi/testpack.go`**

```go
package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Assertion struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AuthConfig struct {
	Type          string  `json:"type"`
	BearerToken   *string `json:"bearerToken,omitempty"`
	APIKeyHeader  *string `json:"apiKeyHeader,omitempty"`
	APIKeyValue   *string `json:"apiKeyValue,omitempty"`
	BasicUsername *string `json:"basicUsername,omitempty"`
	BasicPassword *string `json:"basicPassword,omitempty"`
}

type TestCaseStep struct {
	Order          int    `json:"order"`
	Action         string `json:"action"`
	ExpectedResult string `json:"expectedResult"`
}

type ManualTestCase struct {
	Preconditions   *string        `json:"preconditions,omitempty"`
	TestData        *string        `json:"testData,omitempty"`
	Steps           []TestCaseStep `json:"steps,omitempty"`
	ExpectedOutcome *string        `json:"expectedOutcome,omitempty"`
	Postconditions  *string        `json:"postconditions,omitempty"`
}

type APITestCase struct {
	HTTPMethod         string      `json:"httpMethod"`
	APISpecID          *string     `json:"apiSpecId,omitempty"`
	OperationID        *string     `json:"operationId,omitempty"`
	Auth               *AuthConfig `json:"auth,omitempty"`
	RequestHeaders     []KeyValue  `json:"requestHeaders,omitempty"`
	QueryParams        []KeyValue  `json:"queryParams,omitempty"`
	RequestBody        *string     `json:"requestBody,omitempty"`
	ExpectedStatusCode *int        `json:"expectedStatusCode,omitempty"`
	MaxResponseTimeMs  *int        `json:"maxResponseTimeMs,omitempty"`
	ResponseBody       *string     `json:"responseBody,omitempty"`
	Assertions         []Assertion `json:"assertions,omitempty"`
}

type GraphQLTestCase struct {
	OperationType string      `json:"operationType"`
	OperationName *string     `json:"operationName,omitempty"`
	Query         string      `json:"query"`
	Variables     *string     `json:"variables,omitempty"`
	ResponseBody  *string     `json:"responseBody,omitempty"`
	Assertions    []Assertion `json:"assertions,omitempty"`
	ExpectError   bool        `json:"expectError"`
}

type DatabaseTestCase struct {
	Dialect       string      `json:"dialect"`
	SchemaID      *string     `json:"schemaId,omitempty"`
	Query         string      `json:"query"`
	Assertions    []Assertion `json:"assertions,omitempty"`
	SetupQuery    *string     `json:"setupQuery,omitempty"`
	TeardownQuery *string     `json:"teardownQuery,omitempty"`
}

type GRPCTestCase struct {
	ServiceName    string      `json:"serviceName"`
	MethodName     string      `json:"methodName"`
	CallMode       string      `json:"callMode"`
	ProtoFileID    *string     `json:"protoFileId,omitempty"`
	ServerAddress  *string     `json:"serverAddress,omitempty"`
	RequestMessage *string     `json:"requestMessage,omitempty"`
	Metadata       []KeyValue  `json:"metadata,omitempty"`
	ExpectedStatus string      `json:"expectedStatus"`
	DeadlineMs     *int        `json:"deadlineMs,omitempty"`
	ResponseBody   *string     `json:"responseBody,omitempty"`
	Assertions     []Assertion `json:"assertions,omitempty"`
	UseTLS         bool        `json:"useTLS"`
	ExpectError    bool        `json:"expectError"`
}

type TestPack struct {
	TestPackID string     `json:"testPackId"`
	ServiceID  string     `json:"serviceId"`
	OrgID      string     `json:"orgId"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	CreatedBy  string     `json:"createdBy"`
	UpdatedBy  *string    `json:"updatedBy,omitempty"`
	DeletedBy  *string    `json:"deletedBy,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

type TestCase struct {
	TestCaseID            string            `json:"testCaseId"`
	TestPackID            string            `json:"testPackId"`
	ServiceID             string            `json:"serviceId"`
	OrgID                 string            `json:"orgId"`
	Title                 string            `json:"title"`
	Order                 float64           `json:"order"`
	Type                  string            `json:"type"`
	Description           *string           `json:"description,omitempty"`
	Priority              *string           `json:"priority,omitempty"`
	Labels                []string          `json:"labels,omitempty"`
	LinkedTicket          *string           `json:"linkedTicket,omitempty"`
	EstimatedDurationMins *int              `json:"estimatedDurationMins,omitempty"`
	TestOwner             *string           `json:"testOwner,omitempty"`
	LinkedMapNodeID       *string           `json:"linkedMapNodeId,omitempty"`
	IsCritical            bool              `json:"isCritical"`
	EvidenceRequired      bool              `json:"evidenceRequired"`
	Manual                *ManualTestCase   `json:"manual,omitempty"`
	API                   *APITestCase      `json:"api,omitempty"`
	GraphQL               *GraphQLTestCase  `json:"graphql,omitempty"`
	Database              *DatabaseTestCase `json:"database,omitempty"`
	GRPC                  *GRPCTestCase     `json:"grpc,omitempty"`
	Status                string            `json:"status"`
	Version               int               `json:"version"`
	BaselineRunResultID   *string           `json:"baselineRunResultId,omitempty"`
	Dependencies          []string          `json:"dependencies,omitempty"`
	CreatedBy             string            `json:"createdBy"`
	UpdatedBy             *string           `json:"updatedBy,omitempty"`
	DeletedBy             *string           `json:"deletedBy,omitempty"`
	CreatedAt             time.Time         `json:"createdAt"`
	UpdatedAt             time.Time         `json:"updatedAt"`
	DeletedAt             *time.Time        `json:"deletedAt,omitempty"`
}

type TestRun struct {
	TestRunID     string     `json:"testRunId"`
	TestPackID    string     `json:"testPackId"`
	ServiceID     string     `json:"serviceId"`
	OrgID         string     `json:"orgId"`
	Environment   string     `json:"environment"`
	ReleaseLabel  *string    `json:"releaseLabel,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	Status        string     `json:"status"`
	StartedBy     *string    `json:"startedBy,omitempty"`
	ExecutedBy    string     `json:"executedBy"`
	ExecutedAt    time.Time  `json:"executedAt"`
	OverallStatus string     `json:"overallStatus"`
}

type TestRunSummary struct {
	TestRunID     string     `json:"testRunId"`
	TestPackID    string     `json:"testPackId"`
	ServiceID     string     `json:"serviceId"`
	Environment   string     `json:"environment"`
	ReleaseLabel  *string    `json:"releaseLabel,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	Status        string     `json:"status"`
	StartedBy     *string    `json:"startedBy,omitempty"`
	ExecutedBy    string     `json:"executedBy"`
	ExecutedAt    time.Time  `json:"executedAt"`
	OverallStatus string     `json:"overallStatus"`
	PassedCount   int        `json:"passedCount"`
	FailedCount   int        `json:"failedCount"`
	SkippedCount  int        `json:"skippedCount"`
	BlockedCount  int        `json:"blockedCount"`
}

type TestRunResult struct {
	TestRunResultID string    `json:"testRunResultId"`
	TestRunID       string    `json:"testRunId"`
	TestCaseID      string    `json:"testCaseId"`
	ServiceID       string    `json:"serviceId"`
	OrgID           string    `json:"orgId"`
	Status          string    `json:"status"`
	BlockedReason   *string   `json:"blockedReason,omitempty"`
	ResponseStatus  *int      `json:"responseStatus,omitempty"`
	ResponseBody    *string   `json:"responseBody,omitempty"`
	ResponseTimeMs  *int64    `json:"responseTimeMs,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	ScreenshotURLs  []string  `json:"screenshotUrls,omitempty"`
	ExecutedAt      time.Time `json:"executedAt"`
	ExecutedBy      string    `json:"executedBy"`
}

func (c *Client) ListTestPacks(ctx context.Context, orgID, serviceID string) ([]TestPack, error) {
	var out struct {
		TestPacks []TestPack `json:"testPacks"`
	}
	return out.TestPacks, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-packs", orgID, serviceID), &out)
}

func (c *Client) CreateTestPack(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestPack, error) {
	var out TestPack
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-pack", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestPack(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestPack, error) {
	var out TestPack
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-pack/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteTestPack(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-pack/%s", orgID, serviceID, id))
}

func (c *Client) ListTestCases(ctx context.Context, orgID, serviceID string, testPackID *string) ([]TestCase, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-cases", orgID, serviceID)
	if testPackID != nil && *testPackID != "" {
		q := url.Values{}
		q.Set("testPackId", *testPackID)
		path += "?" + q.Encode()
	}
	var out struct {
		TestCases []TestCase `json:"testCases"`
	}
	return out.TestCases, c.get(ctx, path, &out)
}

func (c *Client) CreateTestCase(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestCase, error) {
	var out TestCase
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-case", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestCase(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestCase, error) {
	var out TestCase
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-case/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteTestCase(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-case/%s", orgID, serviceID, id))
}

func (c *Client) GetTestRun(ctx context.Context, orgID, serviceID, id string) (*TestRun, error) {
	var out TestRun
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run/%s", orgID, serviceID, id), &out)
}

func (c *Client) ListTestRuns(ctx context.Context, orgID, serviceID string, testPackID *string) ([]TestRun, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-runs", orgID, serviceID)
	if testPackID != nil && *testPackID != "" {
		q := url.Values{}
		q.Set("testPackId", *testPackID)
		path += "?" + q.Encode()
	}
	var out struct {
		TestRuns []TestRun `json:"testRuns"`
	}
	return out.TestRuns, c.get(ctx, path, &out)
}

func (c *Client) ListTestRunsSummary(
	ctx context.Context,
	orgID, serviceID string,
	testPackID *string,
	environment *string,
	status *string,
	executedBy *string,
	fromDate *time.Time,
	toDate *time.Time,
) ([]TestRunSummary, error) {
	q := url.Values{}
	if testPackID != nil && *testPackID != "" {
		q.Set("testPackId", *testPackID)
	}
	if environment != nil && *environment != "" {
		q.Set("environment", *environment)
	}
	if status != nil && *status != "" {
		q.Set("status", *status)
	}
	if executedBy != nil && *executedBy != "" {
		q.Set("executedBy", *executedBy)
	}
	if fromDate != nil {
		q.Set("fromDate", fromDate.UTC().Format(time.RFC3339))
	}
	if toDate != nil {
		q.Set("toDate", toDate.UTC().Format(time.RFC3339))
	}
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-runs-summary", orgID, serviceID)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		TestRunsSummary []TestRunSummary `json:"testRunsSummary"`
	}
	return out.TestRunsSummary, c.get(ctx, path, &out)
}

func (c *Client) CreateTestRun(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestRun, error) {
	var out TestRun
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestRun(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestRun, error) {
	var out TestRun
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) ListTestRunResults(ctx context.Context, orgID, serviceID, testRunID string) ([]TestRunResult, error) {
	q := url.Values{}
	q.Set("testRunId", testRunID)
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run-results?%s", orgID, serviceID, q.Encode())
	var out struct {
		TestRunResults []TestRunResult `json:"testRunResults"`
	}
	return out.TestRunResults, c.get(ctx, path, &out)
}

func (c *Client) CreateTestRunResult(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestRunResult, error) {
	var out TestRunResult
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run-result", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestRunResult(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestRunResult, error) {
	var out TestRunResult
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run-result/%s", orgID, serviceID, id), body, &out)
}
```

- [ ] **Step 3: Delete `types.go`**

```bash
rm internal/uigraphapi/types.go
```

By this point every struct and method that was in `types.go` has a new home. Confirm with:

Run: `git status --short internal/uigraphapi/`
Expected: `types.go` shows as deleted (`D`); `catalog.go` shows as modified; `testpack.go` shows as new.

- [ ] **Step 4: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "refactor: split catalog.go into catalog/testpack files, remove types.go"
```

---

## Task 8: Split GraphQL schema files (`content.graphqls`, `catalog.graphqls`)

**Files:**
- Create: `internal/graph/schema/folder.graphqls`, `internal/graph/schema/diagram.graphqls`, `internal/graph/schema/component.graphqls`, `internal/graph/schema/uimap.graphqls`, `internal/graph/schema/testpack.graphqls`
- Modify: `internal/graph/schema/catalog.graphqls` (trimmed)
- Delete: `internal/graph/schema/content.graphqls`
- Regenerated by `go generate`: `internal/graph/generated/generated.go`, `internal/graph/model/models_gen.go`, `internal/graph/folder.resolvers.go`, `internal/graph/diagram.resolvers.go`, `internal/graph/component.resolvers.go`, `internal/graph/uimap.resolvers.go`, `internal/graph/testpack.resolvers.go`, and the now-shrunk `internal/graph/catalog.resolvers.go`

**Interfaces:**
- Consumes: nothing new — this is a schema-file reorganization with the exact same SDL content as today, just redistributed.
- Produces: gqlgen's `follow-schema` layout now generates one resolver file per new schema file. No GraphQL type, field, or query/mutation signature changes.

- [ ] **Step 1: Create `internal/graph/schema/folder.graphqls`**

```graphql
extend type Query {
    folders(orgId: ID!, type: String, parentId: ID):          [Folder!]!
    folder(orgId: ID!, id: ID!):                              Folder!
}

extend type Mutation {
    createFolder(orgId: ID!, input: CreateFolderInput!):              Folder!
    updateFolder(orgId: ID!, id: ID!, input: UpdateFolderInput!):     Folder!
    deleteFolder(orgId: ID!, id: ID!):                                Boolean!
}

type Folder {
    id:        ID!
    orgId:     ID!
    parentId:  ID
    teamId:    ID
    type:      String!
    name:      String!
    order:     Float!
    createdBy: ID!
    createdAt: Time!
    updatedAt: Time!
}

input CreateFolderInput {
    name:     String!
    type:     String!
    parentId: ID
    teamId:   ID
    order:    Float
}

input UpdateFolderInput {
    name:     String
    parentId: ID
    teamId:   ID
    order:    Float
}
```

- [ ] **Step 2: Create `internal/graph/schema/diagram.graphqls`**

```graphql
extend type Query {
    diagrams(orgId: ID!, folderId: ID):                       [Diagram!]!
    diagram(orgId: ID!, id: ID!):                             Diagram!
    diagramContent(orgId: ID!, id: ID!):                      DiagramContent!
    diagramVersions(orgId: ID!, diagramId: ID!):              [DiagramVersion!]!
    diagramVersionContent(orgId: ID!, diagramId: ID!, versionId: ID!): DiagramContent!
    diagramImages(orgId: ID!, diagramId: ID!):               [DiagramImage!]!
}

extend type Mutation {
    createDiagram(orgId: ID!, input: CreateDiagramInput!):            Diagram!
    updateDiagram(orgId: ID!, id: ID!, input: UpdateDiagramInput!):   Diagram!
    deleteDiagram(orgId: ID!, id: ID!):                               Boolean!
    syncDiagram(orgId: ID!, input: SyncDiagramInput!):                SyncDiagramResult!
    createDiagramVersion(orgId: ID!, diagramId: ID!, label: String):  DiagramVersion!
    restoreDiagramVersion(orgId: ID!, diagramId: ID!, versionId: ID!): Diagram!
}

type Diagram {
    id:                 ID!
    orgId:              ID!
    folderId:           ID
    teamId:             ID
    name:               String!
    contentKey:         String!
    contentHash:        String!
    previewAssetId:     String
    previewImageUrl:    String @goField(forceResolver: true)
    previewContentHash: String
    source:             String
    createdBy:          ID!
    updatedBy:          ID
    createdByActor:     Actor @goField(forceResolver: true)
    updatedByActor:     Actor @goField(forceResolver: true)
    createdAt:          Time!
    updatedAt:          Time!
}

type DiagramImage {
    diagramImageId: String!
    diagramId:      String!
    orgId:          ID!
    assetId:        String!
    imageUrl:       String @goField(forceResolver: true)
    fileName:       String
    order:          Int!
    createdBy:      ID!
    createdAt:      Time!
}

type DiagramContent {
    diagramId: ID!
    content:   String!
}

type DiagramVersion {
    id:            ID!
    orgId:         ID!
    diagramId:     ID!
    versionNumber: Int!
    label:         String
    contentKey:    String!
    contentHash:   String!
    isAutoVersion: Boolean!
    source:        String
    createdBy:     ID!
    createdByActor: Actor @goField(forceResolver: true)
    createdAt:     Time!
}

type SyncDiagramResult {
    diagramId:      ID!
    versionCreated: Boolean!
    versionId:      ID
}

input CreateDiagramInput {
    name:     String!
    content:  String!
    folderId: ID
    teamId:   ID
    source:   String
}

input UpdateDiagramInput {
    name:     String
    content:  String
    folderId: ID
    teamId:   ID
    source:   String
}

input SyncDiagramInput {
    diagramId: ID
    name:      String!
    content:   String!
    folderId:  ID
    teamId:    ID
    source:    String
}
```

- [ ] **Step 3: Create `internal/graph/schema/component.graphqls`**

```graphql
extend type Query {
    flowDiagramComponents(orgId: ID!):                        FlowDiagramComponents!
    components(orgId: ID!):                                   Components!
}

type FlowDiagramComponentField {
    flowDiagramComponentFieldId: String!
    label:                       String!
    type:                        String!
    required:                    Boolean!
    readonly:                    Boolean
    options:                     [String!]
    order:                       Int!
}

type FlowDiagramComponent {
    componentId:                String!
    type:                       String!
    name:                       String!
    description:                String!
    category:                   String!
    tags:                       [String!]!
    slug:                       String!
    previewImageJpg:            String!
    isActive:                   Boolean!
    order:                      Int!
    organizationId:             String
    flowDiagramComponentFields: [FlowDiagramComponentField!]!
}

type FlowDiagramComponents {
    components:       [FlowDiagramComponent!]!
    customComponents: [FlowDiagramComponent!]!
}

type ComponentField {
    componentFieldId: String!
    label:            String!
    type:             String!
    required:         Boolean!
    readonly:         Boolean
    options:          [String!]
    order:            Int!
}

type Component {
    componentId:     String!
    type:            String!
    name:            String!
    description:     String!
    category:        String!
    tags:            [String!]!
    slug:            String!
    previewImageJpg: String!
    isActive:        Boolean!
    order:           Int!
    componentFields: [ComponentField!]!
}

type Components {
    components:       [Component!]!
    customComponents: [Component!]!
}
```

- [ ] **Step 4: Create `internal/graph/schema/uimap.graphqls`**

```graphql
extend type Query {
    maps(orgId: ID!, folderId: ID):                           [UIMap!]!
    map(orgId: ID!, id: ID!):                                 UIMap!
    frames(orgId: ID!, mapId: ID!):                           [Frame!]!
    frame(orgId: ID!, mapId: ID!, id: ID!):                   Frame!
    frameById(orgId: ID!, id: ID!):                           Frame!
    focalPoints(orgId: ID!, mapId: ID!, frameId: ID!):        [FocalPoint!]!
    canvas(orgId: ID!, mapId: ID!):                           Canvas!
    frameGroups(orgId: ID!, mapId: ID!, frameId: ID!):        [FrameGroup!]!
    frameLinks(orgId: ID!, mapId: ID!, frameId: ID!):         [FrameLink!]!
    focalPointMeta(orgId: ID!, mapId: ID!, frameId: ID!, focalPointId: ID!): [FocalPointMeta!]!
}

extend type Mutation {
    createMap(orgId: ID!, input: CreateMapInput!):                    UIMap!
    updateMap(orgId: ID!, id: ID!, input: UpdateMapInput!):           UIMap!
    deleteMap(orgId: ID!, id: ID!):                                   Boolean!

    createFrame(orgId: ID!, mapId: ID!, input: CreateFrameInput!):    Frame!
    updateFrame(orgId: ID!, mapId: ID!, id: ID!, input: UpdateFrameInput!): Frame!
    deleteFrame(orgId: ID!, mapId: ID!, id: ID!):                     Boolean!
    syncFrame(orgId: ID!, mapId: ID!, input: SyncFrameInput!):        SyncFrameResult!

    createFocalPoint(orgId: ID!, mapId: ID!, frameId: ID!, input: CreateFocalPointInput!): FocalPoint!
    updateFocalPoint(orgId: ID!, mapId: ID!, frameId: ID!, id: ID!, input: UpdateFocalPointInput!): FocalPoint!
    deleteFocalPoint(orgId: ID!, mapId: ID!, frameId: ID!, id: ID!):  Boolean!

    upsertCanvas(orgId: ID!, mapId: ID!, input: UpsertCanvasInput!):  Canvas!

    createFrameGroup(orgId: ID!, mapId: ID!, frameId: ID!, input: CreateFrameGroupInput!): FrameGroup!
    updateFrameGroup(orgId: ID!, mapId: ID!, frameId: ID!, id: ID!, input: UpdateFrameGroupInput!): FrameGroup!
    deleteFrameGroup(orgId: ID!, mapId: ID!, frameId: ID!, id: ID!): Boolean!

    createFrameLink(orgId: ID!, mapId: ID!, frameId: ID!, input: CreateFrameLinkInput!): FrameLink!
    updateFrameLink(orgId: ID!, mapId: ID!, frameId: ID!, id: ID!, input: UpdateFrameLinkInput!): FrameLink!
    deleteFrameLink(orgId: ID!, mapId: ID!, frameId: ID!, id: ID!): Boolean!

    createFocalPointMeta(orgId: ID!, mapId: ID!, frameId: ID!, focalPointId: ID!, input: CreateFocalPointMetaInput!): FocalPointMeta!
    updateFocalPointMeta(orgId: ID!, mapId: ID!, frameId: ID!, focalPointId: ID!, id: ID!, input: UpdateFocalPointMetaInput!): FocalPointMeta!
    deleteFocalPointMeta(orgId: ID!, mapId: ID!, frameId: ID!, focalPointId: ID!, id: ID!): Boolean!
}

type UIMap {
    id:          ID!
    orgId:       ID!
    folderId:    ID
    teamId:      ID
    name:        String!
    description: String!
    status:      String!
    createdBy:   ID!
    updatedBy:   ID
    createdAt:   Time!
    updatedAt:   Time!
}

type Frame {
    id:                   ID!
    mapId:                ID!
    orgId:                ID!
    parentFrameId:        ID
    name:                 String!
    description:          String!
    templateType:         String!
    screenshotAssetId:    String
    screenshotImageUrl:   String @goField(forceResolver: true)
    screenshotContentHash: String
    status:               String!
    order:                Float!
    source:               String
    createdBy:            ID!
    updatedBy:            ID
    createdByActor:       Actor @goField(forceResolver: true)
    updatedByActor:       Actor @goField(forceResolver: true)
    createdAt:            Time!
    updatedAt:            Time!
}

type SyncFrameResult {
    frameId:        ID!
    versionCreated: Boolean!
}

type FocalPoint {
    id:         ID!
    frameId:    ID!
    orgId:      ID!
    name:       String!
    locationX:  Float!
    locationY:  Float!
    visibility: String!
    isActive:   Boolean!
    createdBy:  ID!
    updatedBy:  ID
    createdAt:  Time!
    updatedAt:  Time!
}

type Canvas {
    mapId:           ID!
    orgId:           ID!
    zoom:            Float!
    navigationX:     Float!
    navigationY:     Float!
    framePositions:  String!
    updatedAt:       Time!
}

type FrameGroup {
    id:          ID!
    frameId:     ID!
    orgId:       ID!
    name:        String!
    description: String!
    locationX:   Float!
    locationY:   Float!
    width:       Float!
    height:      Float!
    order:       Float!
    isActive:    Boolean!
    createdBy:   ID!
    updatedBy:   ID
    createdAt:   Time!
    updatedAt:   Time!
}

type FrameLink {
    id:            ID!
    frameId:       ID!
    orgId:         ID!
    kind:          String!
    targetFrameId: ID
    targetMapId:   ID
    label:         String!
    locationX:     Float!
    locationY:     Float!
    isActive:      Boolean!
    createdBy:     ID!
    updatedBy:     ID
    createdAt:     Time!
    updatedAt:     Time!
}

type FocalPointMeta {
    id:                   ID!
    focalPointId:         ID!
    orgId:                ID!
    frameId:              ID!
    componentId:          String!
    componentLinkId:      String
    componentImages:      String!
    componentFlowDiagram: String
    componentModalFields: String!
    createdBy:            ID!
    updatedBy:            ID
    createdAt:            Time!
    updatedAt:            Time!
}

input CreateMapInput {
    name:        String!
    description: String
    folderId:    ID
    teamId:      ID
}

input UpdateMapInput {
    name:        String
    description: String
    status:      String
    folderId:    ID
    teamId:      ID
}

input CreateFrameInput {
    name:          String!
    description:   String
    templateType:  String!
    parentFrameId: ID
    order:         Float
    screenshot:    String
}

input UpdateFrameInput {
    name:         String
    description:  String
    templateType: String
    status:       String
    order:        Float
    screenshot:   String
}

input SyncFrameInput {
    frameId:      ID
    name:         String!
    templateType: String!
    description:  String
    screenshot:   String!
    source:       String
}

input CreateFocalPointInput {
    name:       String!
    locationX:  Float!
    locationY:  Float!
    visibility: String
    isActive:   Boolean
}

input UpdateFocalPointInput {
    name:       String
    locationX:  Float
    locationY:  Float
    visibility: String
    isActive:   Boolean
}

input UpsertCanvasInput {
    zoom:           Float
    navigationX:    Float
    navigationY:    Float
    framePositions: String
}

input CreateFrameGroupInput {
    name:        String!
    description: String
    locationX:   Float
    locationY:   Float
    width:       Float
    height:      Float
    order:       Float
    isActive:    Boolean
}

input UpdateFrameGroupInput {
    name:        String
    description: String
    locationX:   Float
    locationY:   Float
    width:       Float
    height:      Float
    order:       Float
    isActive:    Boolean
}

input CreateFrameLinkInput {
    kind:          String!
    targetFrameId: ID
    targetMapId:   ID
    label:         String
    locationX:     Float
    locationY:     Float
    isActive:      Boolean
}

input UpdateFrameLinkInput {
    kind:          String
    targetFrameId: ID
    targetMapId:   ID
    label:         String
    locationX:     Float
    locationY:     Float
    isActive:      Boolean
}

input CreateFocalPointMetaInput {
    componentId:          String!
    componentLinkId:      String
    componentImages:      String
    componentFlowDiagram: String
    componentModalFields: String
}

input UpdateFocalPointMetaInput {
    componentId:          String
    componentLinkId:      String
    componentImages:      String
    componentFlowDiagram: String
    componentModalFields: String
}
```

- [ ] **Step 5: Delete `content.graphqls`**

```bash
rm internal/graph/schema/content.graphqls
```

- [ ] **Step 6: Build and regenerate — checkpoint before touching `catalog.graphqls`**

Run: `go generate ./internal/graph/...`
Expected: exits 0. New files appear: `internal/graph/folder.resolvers.go`, `diagram.resolvers.go`, `component.resolvers.go`, `uimap.resolvers.go`. `internal/graph/content.resolvers.go` is deleted by gqlgen (now empty — gqlgen removes resolver files with no remaining fields) or left with a header comment and no functions; if the latter, delete it manually: `rm -f internal/graph/content.resolvers.go` only if `git diff internal/graph/content.resolvers.go` shows no function bodies remain.

Run: `go build ./...`
Expected: exits 0, no output. If it fails with "missing method" errors, gqlgen did not find a matching hand-written implementation for some field — check the error message for the exact resolver type/method name and confirm it exists verbatim (same signature) somewhere under `internal/graph/`.

- [ ] **Step 7: Rewrite `internal/graph/schema/catalog.graphqls` (trimmed)**

Replace its contents with:

```graphql
extend type Query {
    services(orgId: ID!, folderId: ID, teamId: ID):                              [Service!]!
    service(orgId: ID!, id: ID!):                                    Service!
    apiGroups(orgId: ID!, serviceId: ID!):                           [APIGroup!]!
    apiGroup(orgId: ID!, serviceId: ID!, id: ID!):                   APIGroup!
    apiGroupVersions(orgId: ID!, serviceId: ID!, apiGroupId: ID!):   [APIGroupVersion!]!
    serviceDocs(orgId: ID!, serviceId: ID!):                         [ServiceDoc!]!
    serviceDoc(orgId: ID!, serviceId: ID!, id: ID!):                 ServiceDoc!
    serviceDiagrams(orgId: ID!, serviceId: ID!):                     [ServiceDiagram!]!
    serviceDBs(orgId: ID!, serviceId: ID!):                          [ServiceDB!]!
    serviceDB(orgId: ID!, serviceId: ID!, id: ID!):                  ServiceDB!
    serviceDBVersions(orgId: ID!, serviceId: ID!, serviceDbId: ID!): [ServiceDBVersion!]!
    apiEndpoints(orgId: ID!, serviceId: ID!, apiGroupId: ID!):       [APIEndpoint!]!
    apiEndpoint(orgId: ID!, serviceId: ID!, apiGroupId: ID!, id: ID!): APIEndpoint!
    serviceStats(orgId: ID!, serviceId: ID):                         [ServiceStats!]!
}

extend type Mutation {
    createService(orgId: ID!, input: CreateServiceInput!):                         Service!
    updateService(orgId: ID!, id: ID!, input: UpdateServiceInput!):                Service!
    deleteService(orgId: ID!, id: ID!):                                            Boolean!

    createAPIGroup(orgId: ID!, serviceId: ID!, input: CreateAPIGroupInput!):                        APIGroup!
    updateAPIGroup(orgId: ID!, serviceId: ID!, id: ID!, input: UpdateAPIGroupInput!):               APIGroup!
    deleteAPIGroup(orgId: ID!, serviceId: ID!, id: ID!):                                            Boolean!
    syncAPIGroup(orgId: ID!, serviceId: ID!, input: SyncAPIGroupInput!):                            SyncAPIGroupResult!
    createServiceDoc(orgId: ID!, serviceId: ID!, input: CreateServiceDocInput!):                    ServiceDoc!
    updateServiceDoc(orgId: ID!, serviceId: ID!, id: ID!, input: UpdateServiceDocInput!):           ServiceDoc!
    deleteServiceDoc(orgId: ID!, serviceId: ID!, id: ID!):                                          Boolean!
    createServiceDiagram(orgId: ID!, serviceId: ID!, input: CreateServiceDiagramInput!):            ServiceDiagram!
    deleteServiceDiagram(orgId: ID!, serviceId: ID!, diagramId: ID!):                               Boolean!
    createServiceDB(orgId: ID!, serviceId: ID!, input: CreateServiceDBInput!):                      ServiceDB!
    updateServiceDB(orgId: ID!, serviceId: ID!, id: ID!, input: UpdateServiceDBInput!):             ServiceDB!
    deleteServiceDB(orgId: ID!, serviceId: ID!, id: ID!):                                           Boolean!
    createServiceDBVersion(orgId: ID!, serviceId: ID!, serviceDbId: ID!, input: CreateServiceDBVersionInput!): ServiceDBVersion!
    restoreServiceDBVersion(orgId: ID!, serviceId: ID!, serviceDbId: ID!, versionId: ID!):          ServiceDB!

    createAPIEndpoint(orgId: ID!, serviceId: ID!, apiGroupId: ID!, input: CreateAPIEndpointInput!):             APIEndpoint!
    updateAPIEndpoint(orgId: ID!, serviceId: ID!, apiGroupId: ID!, id: ID!, input: UpdateAPIEndpointInput!):    APIEndpoint!
    deleteAPIEndpoint(orgId: ID!, serviceId: ID!, apiGroupId: ID!, id: ID!):                                    Boolean!
}

type Service {
    id:               ID!
    orgId:            ID!
    folderId:         ID
    teamId:           ID
    name:             String!
    slug:             String!
    description:      String!
    status:           String!
    tier:             String!
    category:         String!
    language:         String!
    gitRepoUrl:       String
    jiraProjectUrl:   String
    slackChannelUrl:  String
    lastCommitSha:    String
    labels:           [String!]!
    metadata:         String!
    createdBy:        ID!
    updatedBy:        ID
    createdByActor:   Actor @goField(forceResolver: true)
    updatedByActor:   Actor @goField(forceResolver: true)
    createdAt:        Time!
    updatedAt:        Time!
}

type ServiceStats {
    serviceId:     ID!
    endpointCount: Int!
    diagramCount:  Int!
    docCount:      Int!
    dbTableCount:  Int!
    testCaseCount: Int!
}

type APIGroup {
    id:        ID!
    serviceId: ID!
    orgId:     ID!
    name:      String!
    version:   String!
    label:     String
    protocol:  String!
    specKey:   String
    specHash:  String
    createdBy: ID!
    updatedBy: ID
    createdAt: Time!
    updatedAt: Time!
}

type APIGroupVersion {
    id:            ID!
    orgId:         ID!
    apiGroupId:    ID!
    versionNumber: Int!
    label:         String
    specKey:       String!
    specHash:      String!
    isAutoVersion: Boolean!
    createdBy:     ID!
    createdByActor: Actor @goField(forceResolver: true)
    createdAt:     Time!
}

type ServiceDoc {
    id:          ID!
    serviceId:   ID!
    orgId:       ID!
    fileKey:     String!
    fileName:    String!
    fileType:    String!
    description: String!
    contentHash: String!
    createdBy:   ID!
    updatedBy:   ID
    createdAt:   Time!
    updatedAt:   Time!
}

type ServiceDiagram {
    serviceId: ID!
    diagramId: ID!
    orgId: ID!
    createdBy: ID!
    updatedBy: ID
    createdAt: Time!
    updatedAt: Time!
    diagram: Diagram
}

type ServiceDB {
    id:         ID!
    serviceId:  ID!
    orgId:      ID!
    dbName:     String!
    dbType:     String!
    dialect:    String!
    schemaJson: String!
    source:     String
    sourceTs:   Time
    createdBy:  ID!
    updatedBy:  ID
    createdByActor: Actor @goField(forceResolver: true)
    updatedByActor: Actor @goField(forceResolver: true)
    createdAt:  Time!
    updatedAt:  Time!
}

type ServiceDBVersion {
    id:            ID!
    orgId:         ID!
    serviceDbId:   ID!
    versionNumber: Int!
    label:         String
    schemaJson:    String!
    source:        String
    sourceTs:      Time
    isAutoVersion: Boolean!
    createdBy:     ID!
    createdByActor: Actor @goField(forceResolver: true)
    createdAt:     Time!
}

type APIEndpoint {
    id:          ID!
    apiGroupId:  ID!
    serviceId:   ID!
    orgId:       ID!
    operationId: String!
    method:      String!
    path:        String!
    summary:     String!
    description: String!
    tags:        [String!]!
    parameters:  String!
    requestBody: String!
    responses:   String!
    order:       Float!
    createdBy:   ID!
    updatedBy:   ID
    createdAt:   Time!
    updatedAt:   Time!
}

type SyncAPIGroupResult {
    apiGroupId:     ID!
    versionCreated: Boolean!
}

input CreateServiceInput {
    name:             String!
    slug:             String
    description:      String
    status:           String
    tier:             String
    category:         String
    language:         String
    folderId:         ID
    teamId:           ID
    gitRepoUrl:       String
    jiraProjectUrl:   String
    slackChannelUrl:  String
    labels:           [String!]
    metadata:         String
}

input UpdateServiceInput {
    name:             String
    slug:             String
    description:      String
    status:           String
    tier:             String
    category:         String
    language:         String
    folderId:         ID
    teamId:           ID
    gitRepoUrl:       String
    jiraProjectUrl:   String
    slackChannelUrl:  String
    lastCommitSha:    String
    labels:           [String!]
    metadata:         String
}

input CreateAPIGroupInput {
    name:     String!
    version:  String
    label:    String
    protocol: String
    spec:     String
}

input UpdateAPIGroupInput {
    name:     String
    version:  String
    label:    String
    protocol: String
    spec:     String
}

input SyncAPIGroupInput {
    apiGroupId: ID
    name:       String!
    version:    String
    protocol:   String
    spec:       String!
}

input CreateServiceDocInput {
    fileName:      String!
    fileType:      String
    description:   String
    contentBase64: String!
}

input UpdateServiceDocInput {
    fileName:      String
    fileType:      String
    description:   String
    contentBase64: String
}

input CreateServiceDiagramInput {
    diagramId: ID
    name: String
    content: String
    folderId: ID
    teamId: ID
    source: String
}

input CreateServiceDBInput {
    dbName:     String!
    dbType:     String
    dialect:    String
    schemaJson: String
    source:     String
    sourceTs:   Time
}

input UpdateServiceDBInput {
    dbName:     String
    dbType:     String
    dialect:    String
    schemaJson: String
    source:     String
    sourceTs:   Time
}

input CreateServiceDBVersionInput {
    label:         String
    isAutoVersion: Boolean
    dbName:        String
    dbType:        String
    dialect:       String
    schemaJson:    String
    source:        String
    sourceTs:      Time
}

input CreateAPIEndpointInput {
    operationId: String
    method:      String!
    path:        String!
    summary:     String
    description: String
    tags:        [String!]
    parameters:  String
    requestBody: String
    responses:   String
    order:       Float
}

input UpdateAPIEndpointInput {
    operationId: String
    method:      String
    path:        String
    summary:     String
    description: String
    tags:        [String!]
    parameters:  String
    requestBody: String
    responses:   String
    order:       Float
}
```

- [ ] **Step 8: Create `internal/graph/schema/testpack.graphqls`**

```graphql
extend type Query {
    testPacks(orgId: ID!, serviceId: ID!):                           [TestPack!]!
    testCases(orgId: ID!, serviceId: ID!, testPackId: ID):           [TestCase!]!
    testRun(orgId: ID!, serviceId: ID!, id: ID!):                    TestRun!
    testRuns(orgId: ID!, serviceId: ID!, testPackId: ID):            [TestRun!]!
    testRunsSummary(
        orgId: ID!,
        serviceId: ID!,
        testPackId: ID,
        environment: String,
        status: String,
        executedBy: ID,
        fromDate: Time,
        toDate: Time
    ): [TestRunSummary!]!
    testRunResults(orgId: ID!, serviceId: ID!, testRunId: ID!):      [TestRunResult!]!
}

extend type Mutation {
    createTestPack(orgId: ID!, serviceId: ID!, input: CreateTestPackInput!):                                    TestPack!
    updateTestPack(orgId: ID!, serviceId: ID!, id: ID!, input: UpdateTestPackInput!):                           TestPack!
    deleteTestPack(orgId: ID!, serviceId: ID!, id: ID!):                                                         Boolean!
    createTestCase(orgId: ID!, serviceId: ID!, input: CreateTestCaseInput!):                                     TestCase!
    updateTestCase(orgId: ID!, serviceId: ID!, id: ID!, input: UpdateTestCaseInput!):                           TestCase!
    deleteTestCase(orgId: ID!, serviceId: ID!, id: ID!):                                                         Boolean!
    createTestRun(orgId: ID!, serviceId: ID!, input: CreateTestRunInput!):                                       TestRun!
    updateTestRun(orgId: ID!, serviceId: ID!, id: ID!, input: UpdateTestRunInput!):                             TestRun!
    createTestRunResult(orgId: ID!, serviceId: ID!, input: CreateTestRunResultInput!):                          TestRunResult!
    updateTestRunResult(orgId: ID!, serviceId: ID!, id: ID!, input: UpdateTestRunResultInput!):                TestRunResult!
}

type KeyValue {
    key:   String!
    value: String!
}

type Assertion {
    field: String!
    type:  String!
    value: String!
}

type AuthConfig {
    type:          String!
    bearerToken:   String
    apiKeyHeader:  String
    apiKeyValue:   String
    basicUsername: String
    basicPassword: String
}

type TestCaseStep {
    order:          Int!
    action:         String!
    expectedResult: String!
}

type ManualTestCase {
    preconditions:   String
    testData:        String
    steps:           [TestCaseStep!]
    expectedOutcome: String
    postconditions:  String
}

type APITestCase {
    httpMethod:         String!
    apiSpecId:          String
    operationId:        String
    auth:               AuthConfig
    requestHeaders:     [KeyValue!]
    queryParams:        [KeyValue!]
    requestBody:        String
    expectedStatusCode: Int
    maxResponseTimeMs:  Int
    responseBody:       String
    assertions:         [Assertion!]
}

type GraphQLTestCase {
    operationType: String!
    operationName: String
    query:         String!
    variables:     String
    responseBody:  String
    assertions:    [Assertion!]
    expectError:   Boolean!
}

type DatabaseTestCase {
    dialect:       String!
    schemaId:      String
    query:         String!
    assertions:    [Assertion!]
    setupQuery:    String
    teardownQuery: String
}

type GRPCTestCase {
    serviceName:    String!
    methodName:     String!
    callMode:       String!
    protoFileId:    String
    serverAddress:  String
    requestMessage: String
    metadata:       [KeyValue!]
    expectedStatus: String!
    deadlineMs:     Int
    responseBody:   String
    assertions:     [Assertion!]
    useTLS:         Boolean!
    expectError:    Boolean!
}

type TestPack {
    testPackId: ID!
    serviceId:  ID!
    orgId:      ID!
    name:       String!
    type:       String!
    createdBy:  ID!
    updatedBy:  ID
    deletedBy:  ID
    createdAt:  Time!
    updatedAt:  Time!
    deletedAt:  Time
}

type TestCase {
    testCaseId:             ID!
    testPackId:             ID!
    serviceId:              ID!
    orgId:                  ID!
    title:                  String!
    order:                  Float!
    type:                   String!
    description:            String
    priority:               String
    labels:                 [String!]!
    linkedTicket:           String
    estimatedDurationMins:  Int
    testOwner:              String
    linkedMapNodeId:        ID
    isCritical:             Boolean!
    evidenceRequired:       Boolean!
    manual:                 ManualTestCase
    api:                    APITestCase
    graphql:                GraphQLTestCase
    database:               DatabaseTestCase
    grpc:                   GRPCTestCase
    status:                 String!
    version:                Int!
    baselineRunResultId:    ID
    dependencies:           [ID!]!
    createdBy:              ID!
    updatedBy:              ID
    deletedBy:              ID
    createdAt:              Time!
    updatedAt:              Time!
    deletedAt:              Time
}

type TestRun {
    testRunId:     ID!
    testPackId:    ID!
    serviceId:     ID!
    orgId:         ID!
    environment:   String!
    releaseLabel:  String
    startedAt:     Time
    completedAt:   Time
    status:        String!
    startedBy:     ID
    executedBy:    ID!
    executedAt:    Time!
    overallStatus: String!
}

type TestRunSummary {
    testRunId:     ID!
    testPackId:    ID!
    serviceId:     ID!
    environment:   String!
    releaseLabel:  String
    startedAt:     Time
    completedAt:   Time
    status:        String!
    startedBy:     ID
    executedBy:    ID!
    executedAt:    Time!
    overallStatus: String!
    passedCount:   Int!
    failedCount:   Int!
    skippedCount:  Int!
    blockedCount:  Int!
}

type TestRunResult {
    testRunResultId: ID!
    testRunId:       ID!
    testCaseId:      ID!
    serviceId:       ID!
    orgId:           ID!
    status:          String!
    blockedReason:   String
    responseStatus:  Int
    responseBody:    String
    responseTimeMs:  Int
    notes:           String
    screenshotUrls:  [String!]!
    executedAt:      Time!
    executedBy:      ID!
}

input KeyValueInput {
    key:   String!
    value: String!
}

input AssertionInput {
    field: String!
    type:  String!
    value: String!
}

input AuthConfigInput {
    type:          String!
    bearerToken:   String
    apiKeyHeader:  String
    apiKeyValue:   String
    basicUsername: String
    basicPassword: String
}

input TestCaseStepInput {
    order:          Int!
    action:         String!
    expectedResult: String!
}

input ManualTestCaseInput {
    preconditions:   String
    testData:        String
    steps:           [TestCaseStepInput!]
    expectedOutcome: String
    postconditions:  String
}

input APITestCaseInput {
    httpMethod:         String!
    apiSpecId:          String
    operationId:        String
    auth:               AuthConfigInput
    requestHeaders:     [KeyValueInput!]
    queryParams:        [KeyValueInput!]
    requestBody:        String
    expectedStatusCode: Int
    maxResponseTimeMs:  Int
    responseBody:       String
    assertions:         [AssertionInput!]
}

input GraphQLTestCaseInput {
    operationType: String!
    operationName: String
    query:         String!
    variables:     String
    responseBody:  String
    assertions:    [AssertionInput!]
    expectError:   Boolean!
}

input DatabaseTestCaseInput {
    dialect:       String!
    schemaId:      String
    query:         String!
    assertions:    [AssertionInput!]
    setupQuery:    String
    teardownQuery: String
}

input GRPCTestCaseInput {
    serviceName:    String!
    methodName:     String!
    callMode:       String!
    protoFileId:    String
    serverAddress:  String
    requestMessage: String
    metadata:       [KeyValueInput!]
    expectedStatus: String!
    deadlineMs:     Int
    responseBody:   String
    assertions:     [AssertionInput!]
    useTLS:         Boolean!
    expectError:    Boolean!
}

input CreateTestPackInput {
    name: String!
    type: String
}

input UpdateTestPackInput {
    name: String
    type: String
}

input CreateTestCaseInput {
    testPackId:            ID!
    title:                 String!
    order:                 Float
    type:                  String!
    description:           String
    priority:              String
    labels:                [String!]
    linkedTicket:          String
    estimatedDurationMins: Int
    testOwner:             String
    linkedMapNodeId:       ID
    isCritical:            Boolean
    evidenceRequired:      Boolean
    manual:                ManualTestCaseInput
    api:                   APITestCaseInput
    graphql:               GraphQLTestCaseInput
    database:              DatabaseTestCaseInput
    grpc:                  GRPCTestCaseInput
    status:                String
    version:               Int
    baselineRunResultId:   ID
    dependencies:          [ID!]
}

input UpdateTestCaseInput {
    testPackId:            ID
    title:                 String
    order:                 Float
    type:                  String
    description:           String
    priority:              String
    labels:                [String!]
    linkedTicket:          String
    estimatedDurationMins: Int
    testOwner:             String
    linkedMapNodeId:       ID
    isCritical:            Boolean
    evidenceRequired:      Boolean
    manual:                ManualTestCaseInput
    api:                   APITestCaseInput
    graphql:               GraphQLTestCaseInput
    database:              DatabaseTestCaseInput
    grpc:                  GRPCTestCaseInput
    status:                String
    version:               Int
    baselineRunResultId:   ID
    dependencies:          [ID!]
}

input CreateTestRunInput {
    testPackId:    ID!
    environment:   String!
    releaseLabel:  String
    startedAt:     Time
    completedAt:   Time
    status:        String
    startedBy:     ID
    overallStatus: String
}

input UpdateTestRunInput {
    overallStatus: String
    completedAt:   Time
    status:        String
}

input CreateTestRunResultInput {
    testRunId:      ID!
    testCaseId:     ID!
    status:         String!
    blockedReason:  String
    responseStatus: Int
    responseBody:   String
    responseTimeMs: Int
    notes:          String
    screenshotUrls: [String!]
}

input UpdateTestRunResultInput {
    status:         String
    blockedReason:  String
    responseStatus: Int
    responseBody:   String
    responseTimeMs: Int
    notes:          String
    screenshotUrls: [String!]
}
```

- [ ] **Step 9: Regenerate and build**

Run: `go generate ./internal/graph/...`
Expected: exits 0. `internal/graph/testpack.resolvers.go` is created; `internal/graph/catalog.resolvers.go` shrinks to only the Service/APIGroup/Doc/Diagram-link/DB/Endpoint resolvers.

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 10: Commit**

```bash
git add -A
git commit -m "refactor: split content.graphqls and catalog.graphqls into per-domain schema files"
```

---

## Task 9: Extract `convert.go` into `internal/graph/convert/` package

**Files:**
- Create: `internal/graph/convert/helpers.go`, `auth.go`, `org.go`, `admin.go`, `folder.go`, `diagram.go`, `component.go`, `uimap.go`, `catalog.go`, `testpack.go`
- Create: `internal/graph/refs.go` (resolveActor/resolveAssetURL, moved out of `convert.go` unchanged — these need `r.Client` and can't be pure functions)
- Modify: every `internal/graph/*.resolvers.go` file (call-site renames + import)
- Delete: `internal/graph/convert.go`, `internal/graph/helpers.go`

**Interfaces:**
- Consumes: `internal/graph/model` (gqlgen-generated), `internal/uigraphapi` types (Task 2/5/6/7).
- Produces: package `internal/graph/convert` exporting one `XxxToModel` function per type (full list in Step 9 below), plus `convert.ToMap`, `convert.RawStr`, `convert.RawArrStr`, `convert.FocalPointMetaBody`, `convert.StrFromMap`, `convert.BoolFromMap`, `convert.OptStrFromMap`, `convert.UnmarshalJSONString`. `internal/graph` keeps `Resolver.resolveActor` / `Resolver.resolveAssetURL` as before, just relocated to `refs.go`.

- [ ] **Step 1: Create `internal/graph/convert/helpers.go`**

```go
// Package convert maps internal/uigraphapi REST DTOs onto internal/graph/model
// GraphQL models. Every function here is pure — no I/O, no context — which is
// what makes this package unit-testable without a running server.
package convert

import "encoding/json"

func StrFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func BoolFromMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func OptStrFromMap(m map[string]interface{}, key string) *string {
	if v, ok := m[key].(string); ok && v != "" {
		return &v
	}
	return nil
}

func UnmarshalJSONString(s string, out interface{}) error {
	return json.Unmarshal([]byte(s), out)
}

// ToMap JSON-round-trips a struct into map[string]interface{}.
// This correctly handles optional fields: nil pointer fields are omitted
// from the resulting map (because of omitempty in the input struct tags).
func ToMap(v interface{}) map[string]interface{} {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

// RawStr returns the JSON string of a raw message, defaulting to "{}".
func RawStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "{}"
	}
	return string(b)
}

func RawArrStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "[]"
	}
	return string(b)
}
```

- [ ] **Step 2: Create `internal/graph/convert/auth.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func MeToModel(m *uigraphapi.MeResponse) *model.Me {
	me := &model.Me{
		UserID: m.UserID, OrgID: m.OrgID,
		Email: m.Email, Name: m.Name, Login: m.Login,
		Kind: m.Kind, Role: m.Role, AuthProvider: m.AuthProvider,
	}
	if m.AvatarURL != "" {
		me.AvatarURL = &m.AvatarURL
	}
	return me
}

func OrgSummaryToModel(o uigraphapi.OrgSummary) *model.OrgSummary {
	return &model.OrgSummary{ID: o.ID, Name: o.Name, Slug: o.Slug, Role: o.Role, Active: o.Active}
}

func OrgSummariesToModel(orgs []uigraphapi.OrgSummary) []*model.OrgSummary {
	out := make([]*model.OrgSummary, len(orgs))
	for i, o := range orgs {
		out[i] = OrgSummaryToModel(o)
	}
	return out
}
```

- [ ] **Step 3: Create `internal/graph/convert/org.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func OrgToModel(o *uigraphapi.Org) *model.Org {
	return &model.Org{ID: o.ID, Name: o.Name, Slug: o.Slug, Disabled: o.Disabled, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt}
}

func MemberToModel(m uigraphapi.Member) *model.Member {
	return &model.Member{UserID: m.UserID, OrgID: m.OrgID, Role: m.Role, Source: m.Source, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func TeamToModel(t *uigraphapi.Team) *model.Team {
	m := &model.Team{ID: t.ID, OrgID: t.OrgID, Name: t.Name, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
	if t.Email != "" {
		m.Email = &t.Email
	}
	if t.ExternalID != "" {
		m.ExternalID = &t.ExternalID
	}
	return m
}

func TeamMemberToModel(m uigraphapi.TeamMember) *model.TeamMember {
	return &model.TeamMember{TeamID: m.TeamID, UserID: m.UserID, Permission: m.Permission, CreatedAt: m.CreatedAt}
}

func InvitationToModel(i uigraphapi.Invitation) *model.Invitation {
	return &model.Invitation{
		ID: i.ID, OrgID: i.OrgID, Email: i.Email, Role: i.Role,
		Code: i.Code, CreatedBy: i.CreatedBy, CreatedAt: i.CreatedAt, ExpiresAt: i.ExpiresAt,
	}
}

func ServiceAccountToModel(sa uigraphapi.ServiceAccount) *model.ServiceAccount {
	return &model.ServiceAccount{
		ID: sa.ID, OrgID: sa.OrgID, Name: sa.Name, Description: sa.Description,
		Role: sa.Role, Disabled: sa.Disabled, CreatedAt: sa.CreatedAt, UpdatedAt: sa.UpdatedAt,
	}
}

func SATokenToModel(t uigraphapi.ServiceAccountToken) *model.ServiceAccountToken {
	return &model.ServiceAccountToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name, Prefix: t.Prefix,
		ExpiresAt: t.ExpiresAt, LastUsedAt: t.LastUsedAt, Revoked: t.Revoked, CreatedAt: t.CreatedAt,
	}
}

func CreatedTokenToModel(t *uigraphapi.CreatedToken) *model.CreatedToken {
	return &model.CreatedToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name,
		Prefix: t.Prefix, Token: t.Token, CreatedAt: t.CreatedAt,
	}
}

func OrgsToModel(orgs []uigraphapi.Org) []*model.Org {
	out := make([]*model.Org, len(orgs))
	for i := range orgs {
		out[i] = OrgToModel(&orgs[i])
	}
	return out
}

func MembersToModel(members []uigraphapi.Member) []*model.Member {
	out := make([]*model.Member, len(members))
	for i, m := range members {
		out[i] = MemberToModel(m)
	}
	return out
}

func TeamsToModel(teams []uigraphapi.Team) []*model.Team {
	out := make([]*model.Team, len(teams))
	for i := range teams {
		out[i] = TeamToModel(&teams[i])
	}
	return out
}

func TeamMembersToModel(members []uigraphapi.TeamMember) []*model.TeamMember {
	out := make([]*model.TeamMember, len(members))
	for i, m := range members {
		out[i] = TeamMemberToModel(m)
	}
	return out
}

func InvitationsToModel(invs []uigraphapi.Invitation) []*model.Invitation {
	out := make([]*model.Invitation, len(invs))
	for i, inv := range invs {
		out[i] = InvitationToModel(inv)
	}
	return out
}

func ServiceAccountsToModel(sas []uigraphapi.ServiceAccount) []*model.ServiceAccount {
	out := make([]*model.ServiceAccount, len(sas))
	for i, sa := range sas {
		out[i] = ServiceAccountToModel(sa)
	}
	return out
}

func SATokensToModel(tokens []uigraphapi.ServiceAccountToken) []*model.ServiceAccountToken {
	out := make([]*model.ServiceAccountToken, len(tokens))
	for i, t := range tokens {
		out[i] = SATokenToModel(t)
	}
	return out
}
```

- [ ] **Step 4: Create `internal/graph/convert/admin.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func UserToModel(u *uigraphapi.User) *model.User {
	return &model.User{
		ID: u.ID, Email: u.Email, Name: u.Name, Login: u.Login,
		Disabled: u.Disabled, Role: u.Role, LastSeenAt: u.LastSeenAt,
		CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}

func UsersToModel(users []uigraphapi.User) []*model.User {
	out := make([]*model.User, len(users))
	for i := range users {
		out[i] = UserToModel(&users[i])
	}
	return out
}

func OAuthProviderToModel(p uigraphapi.OAuthProvider) *model.OAuthProvider {
	return &model.OAuthProvider{
		ID: p.ID, ProviderName: p.ProviderName, Type: p.Type, DisplayName: p.DisplayName,
		ClientID: p.ClientID, ClientSecret: p.ClientSecret,
		AuthURL: p.AuthURL, TokenURL: p.TokenURL, UserinfoURL: p.UserinfoURL, APIURL: p.APIURL,
		Scopes: p.Scopes, AllowedDomains: p.AllowedDomains, AllowSignUp: p.AllowSignUp,
		EmailClaim: p.EmailClaim, NameClaim: p.NameClaim, SubClaim: p.SubClaim,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func OAuthProvidersToModel(providers []uigraphapi.OAuthProvider) []*model.OAuthProvider {
	out := make([]*model.OAuthProvider, len(providers))
	for i := range providers {
		out[i] = OAuthProviderToModel(providers[i])
	}
	return out
}

func RoleMappingToModel(m uigraphapi.RoleMapping) *model.RoleMapping {
	return &model.RoleMapping{
		ID: m.ID, OrganizationID: m.OrganizationID,
		ClaimKey: m.ClaimKey, ClaimValue: m.ClaimValue, Role: m.Role, Scope: m.Scope,
		ResourceType: m.ResourceType, ResourceID: m.ResourceID,
	}
}

func RoleMappingsToModel(mappings []uigraphapi.RoleMapping) []*model.RoleMapping {
	out := make([]*model.RoleMapping, len(mappings))
	for i := range mappings {
		out[i] = RoleMappingToModel(mappings[i])
	}
	return out
}

func LDAPToModel(l *uigraphapi.LDAPConfig) *model.LDAPConfig {
	return &model.LDAPConfig{
		ID: l.ID, Host: l.Host, Port: l.Port,
		UseSsl: l.UseSSL, StartTLS: l.StartTLS, SkipTLSVerify: l.SkipTLSVerify,
		BindDn: l.BindDN, BindPassword: l.BindPassword,
		SearchBaseDn: l.SearchBaseDN, SearchFilter: l.SearchFilter,
		EmailAttribute: l.EmailAttribute, NameAttribute: l.NameAttribute,
		UsernameAttribute: l.UsernameAttribute, MemberOfAttribute: l.MemberOfAttribute,
		AllowSignUp: l.AllowSignUp, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func SAMLToModel(s *uigraphapi.SAMLConfig) *model.SAMLConfig {
	return &model.SAMLConfig{
		ID: s.ID, IdpMetadataURL: s.IDPMetadataURL, IdpMetadataXML: s.IDPMetadataXML,
		IdpEntityID: s.IDPEntityID, IdpSsoURL: s.IDPSsoURL, IdpCert: s.IDPCert,
		SpEntityID: s.SPEntityID, SpCert: s.SPCert, SpKey: s.SPKey,
		SignRequests: s.SignRequests, NameIDFormat: s.NameIDFormat,
		EmailAttribute: s.EmailAttribute, NameAttribute: s.NameAttribute, LoginAttribute: s.LoginAttribute,
		GroupsAttribute: s.GroupsAttribute, AllowSignUp: s.AllowSignUp,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}
```

- [ ] **Step 5: Create `internal/graph/convert/folder.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func FolderToModel(f *uigraphapi.Folder) *model.Folder {
	return &model.Folder{
		ID: f.ID, OrgID: f.OrgID, ParentID: f.ParentID, TeamID: f.TeamID, Type: f.Type,
		Name: f.Name, Order: f.Order, CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func FoldersToModel(folders []uigraphapi.Folder) []*model.Folder {
	out := make([]*model.Folder, len(folders))
	for i := range folders {
		out[i] = FolderToModel(&folders[i])
	}
	return out
}
```

- [ ] **Step 6: Create `internal/graph/convert/diagram.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func DiagramToModel(d *uigraphapi.Diagram) *model.Diagram {
	return &model.Diagram{
		ID: d.ID, OrgID: d.OrgID, FolderID: d.FolderID, TeamID: d.TeamID,
		Name: d.Name, ContentKey: d.ContentKey, ContentHash: d.ContentHash,
		PreviewAssetID: d.PreviewAssetID, PreviewContentHash: d.PreviewContentHash,
		Source: d.Source, CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func DiagramVersionToModel(orgID string, v uigraphapi.DiagramVersion) *model.DiagramVersion {
	return &model.DiagramVersion{
		ID: v.ID, OrgID: orgID, DiagramID: v.DiagramID, VersionNumber: v.VersionNumber,
		Label: v.Label, ContentKey: v.ContentKey, ContentHash: v.ContentHash,
		IsAutoVersion: v.IsAutoVersion, Source: v.Source, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func DiagramImageToModel(img uigraphapi.DiagramImage) *model.DiagramImage {
	return &model.DiagramImage{
		DiagramImageID: img.DiagramImageID, DiagramID: img.DiagramID,
		OrgID: img.OrgID, AssetID: img.AssetID, FileName: img.FileName,
		Order: img.Order, CreatedBy: img.CreatedBy, CreatedAt: img.CreatedAt,
	}
}

func DiagramsToModel(diagrams []uigraphapi.Diagram) []*model.Diagram {
	out := make([]*model.Diagram, len(diagrams))
	for i := range diagrams {
		out[i] = DiagramToModel(&diagrams[i])
	}
	return out
}

func DiagramVersionsToModel(orgID string, versions []uigraphapi.DiagramVersion) []*model.DiagramVersion {
	out := make([]*model.DiagramVersion, len(versions))
	for i, v := range versions {
		out[i] = DiagramVersionToModel(orgID, v)
	}
	return out
}

func DiagramImagesToModel(images []uigraphapi.DiagramImage) []*model.DiagramImage {
	out := make([]*model.DiagramImage, len(images))
	for i, img := range images {
		out[i] = DiagramImageToModel(img)
	}
	return out
}
```

- [ ] **Step 7: Create `internal/graph/convert/component.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func FlowComponentFieldToModel(f uigraphapi.FlowDiagramComponentField) *model.FlowDiagramComponentField {
	return &model.FlowDiagramComponentField{
		FlowDiagramComponentFieldID: f.FlowDiagramComponentFieldID,
		Label:                       f.Label,
		Type:                        f.Type,
		Required:                    f.Required,
		Readonly:                    f.Readonly,
		Options:                     f.Options,
		Order:                       f.Order,
	}
}

func FlowComponentToModel(c uigraphapi.FlowDiagramComponent) *model.FlowDiagramComponent {
	fields := make([]*model.FlowDiagramComponentField, len(c.FlowDiagramComponentFields))
	for i, f := range c.FlowDiagramComponentFields {
		fields[i] = FlowComponentFieldToModel(f)
	}
	return &model.FlowDiagramComponent{
		ComponentID: c.ComponentID, Type: c.Type, Name: c.Name,
		Description: c.Description, Category: c.Category, Tags: c.Tags,
		Slug: c.Slug, PreviewImageJpg: c.PreviewImageJpg, IsActive: c.IsActive,
		Order: c.Order, OrganizationID: c.OrganizationID,
		FlowDiagramComponentFields: fields,
	}
}

func FlowComponentsToModel(components []uigraphapi.FlowDiagramComponent) []*model.FlowDiagramComponent {
	out := make([]*model.FlowDiagramComponent, len(components))
	for i, c := range components {
		out[i] = FlowComponentToModel(c)
	}
	return out
}

func ComponentFieldToModel(f uigraphapi.ComponentField) *model.ComponentField {
	return &model.ComponentField{
		ComponentFieldID: f.ComponentFieldID,
		Label:            f.Label,
		Type:             f.Type,
		Required:         f.Required,
		Readonly:         f.Readonly,
		Options:          f.Options,
		Order:            f.Order,
	}
}

func ComponentToModel(c uigraphapi.Component) *model.Component {
	fields := make([]*model.ComponentField, len(c.ComponentFields))
	for i, f := range c.ComponentFields {
		fields[i] = ComponentFieldToModel(f)
	}
	return &model.Component{
		ComponentID: c.ComponentID, Type: c.Type, Name: c.Name,
		Description: c.Description, Category: c.Category, Tags: c.Tags,
		Slug: c.Slug, PreviewImageJpg: c.PreviewImageJpg, IsActive: c.IsActive,
		Order: c.Order, ComponentFields: fields,
	}
}

func ComponentsToModel(components []uigraphapi.Component) []*model.Component {
	out := make([]*model.Component, len(components))
	for i, c := range components {
		out[i] = ComponentToModel(c)
	}
	return out
}
```

- [ ] **Step 8: Create `internal/graph/convert/uimap.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func UIMapToModel(m *uigraphapi.UIMap) *model.UIMap {
	return &model.UIMap{
		ID: m.ID, OrgID: m.OrgID, FolderID: m.FolderID, TeamID: m.TeamID,
		Name: m.Name, Description: m.Description, Status: m.Status,
		CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func FrameToModel(f *uigraphapi.Frame) *model.Frame {
	return &model.Frame{
		ID: f.ID, MapID: f.MapID, OrgID: f.OrgID, ParentFrameID: f.ParentFrameID,
		Name: f.Name, Description: f.Description, TemplateType: f.TemplateType,
		ScreenshotAssetID: f.ScreenshotAssetID, ScreenshotContentHash: f.ScreenshotContentHash,
		Status: f.Status, Order: f.Order, Source: f.Source,
		CreatedBy: f.CreatedBy, UpdatedBy: f.UpdatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func FocalPointToModel(fp *uigraphapi.FocalPoint) *model.FocalPoint {
	return &model.FocalPoint{
		ID: fp.ID, FrameID: fp.FrameID, OrgID: fp.OrgID,
		Name: fp.Name, LocationX: fp.LocationX, LocationY: fp.LocationY,
		Visibility: fp.Visibility, IsActive: fp.IsActive,
		CreatedBy: fp.CreatedBy, UpdatedBy: fp.UpdatedBy, CreatedAt: fp.CreatedAt, UpdatedAt: fp.UpdatedAt,
	}
}

func CanvasToModel(c *uigraphapi.Canvas) *model.Canvas {
	return &model.Canvas{
		MapID: c.MapID, OrgID: c.OrgID,
		Zoom: c.Zoom, NavigationX: c.NavigationX, NavigationY: c.NavigationY,
		FramePositions: RawStr(c.FramePositions),
		UpdatedAt:      c.UpdatedAt,
	}
}

func FrameGroupToModel(g *uigraphapi.FrameGroup) *model.FrameGroup {
	return &model.FrameGroup{
		ID: g.ID, FrameID: g.FrameID, OrgID: g.OrgID,
		Name: g.Name, Description: g.Description,
		LocationX: g.LocationX, LocationY: g.LocationY,
		Width: g.Width, Height: g.Height, Order: g.Order, IsActive: g.IsActive,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy,
		CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func FrameGroupsToModel(gs []uigraphapi.FrameGroup) []*model.FrameGroup {
	out := make([]*model.FrameGroup, len(gs))
	for i := range gs {
		out[i] = FrameGroupToModel(&gs[i])
	}
	return out
}

func FrameLinkToModel(l *uigraphapi.FrameLink) *model.FrameLink {
	return &model.FrameLink{
		ID: l.ID, FrameID: l.FrameID, OrgID: l.OrgID, Kind: l.Kind,
		TargetFrameID: l.TargetFrameID, TargetMapID: l.TargetMapID,
		Label: l.Label, LocationX: l.LocationX, LocationY: l.LocationY, IsActive: l.IsActive,
		CreatedBy: l.CreatedBy, UpdatedBy: l.UpdatedBy,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func FrameLinksToModel(ls []uigraphapi.FrameLink) []*model.FrameLink {
	out := make([]*model.FrameLink, len(ls))
	for i := range ls {
		out[i] = FrameLinkToModel(&ls[i])
	}
	return out
}

func FocalPointMetaToModel(m *uigraphapi.FocalPointMeta) *model.FocalPointMeta {
	return &model.FocalPointMeta{
		ID: m.ID, FocalPointID: m.FocalPointID, OrgID: m.OrgID, FrameID: m.FrameID,
		ComponentID: m.ComponentID, ComponentLinkID: m.ComponentLinkID,
		ComponentImages:      RawArrStr(m.ComponentImages),
		ComponentFlowDiagram: m.ComponentFlowDiagram,
		ComponentModalFields: RawArrStr(m.ComponentModalFields),
		CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func FocalPointMetasToModel(ms []uigraphapi.FocalPointMeta) []*model.FocalPointMeta {
	out := make([]*model.FocalPointMeta, len(ms))
	for i := range ms {
		out[i] = FocalPointMetaToModel(&ms[i])
	}
	return out
}

func FocalPointMetaBody(body map[string]interface{}) map[string]interface{} {
	for _, key := range []string{"componentImages", "componentModalFields"} {
		if s, ok := body[key].(string); ok {
			var raw interface{}
			if err := UnmarshalJSONString(s, &raw); err == nil {
				body[key] = raw
			}
		}
	}
	return body
}

func UIMapsToModel(maps []uigraphapi.UIMap) []*model.UIMap {
	out := make([]*model.UIMap, len(maps))
	for i := range maps {
		out[i] = UIMapToModel(&maps[i])
	}
	return out
}

func FramesToModel(frames []uigraphapi.Frame) []*model.Frame {
	out := make([]*model.Frame, len(frames))
	for i := range frames {
		out[i] = FrameToModel(&frames[i])
	}
	return out
}

func FocalPointsToModel(fps []uigraphapi.FocalPoint) []*model.FocalPoint {
	out := make([]*model.FocalPoint, len(fps))
	for i := range fps {
		out[i] = FocalPointToModel(&fps[i])
	}
	return out
}
```

- [ ] **Step 9: Create `internal/graph/convert/catalog.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func ServiceToModel(s *uigraphapi.Service) *model.Service {
	return &model.Service{
		ID: s.ID, OrgID: s.OrgID, FolderID: s.FolderID, TeamID: s.TeamID,
		Name: s.Name, Slug: s.Slug, Description: s.Description,
		Status: s.Status, Tier: s.Tier, Category: s.Category, Language: s.Language,
		GitRepoURL: s.GitRepoURL, JiraProjectURL: s.JiraProjectURL,
		SlackChannelURL: s.SlackChannelURL, LastCommitSha: s.LastCommitSha,
		Labels:    s.Labels,
		Metadata:  RawStr(s.Metadata),
		CreatedBy: s.CreatedBy, UpdatedBy: s.UpdatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func ServiceStatsToModel(s uigraphapi.ServiceStats) *model.ServiceStats {
	return &model.ServiceStats{
		ServiceID:     s.ServiceID,
		EndpointCount: s.EndpointCount,
		DiagramCount:  s.DiagramCount,
		DocCount:      s.DocCount,
		DbTableCount:  s.DBTableCount,
		TestCaseCount: s.TestCaseCount,
	}
}

func APIGroupToModel(g *uigraphapi.APIGroup) *model.APIGroup {
	return &model.APIGroup{
		ID: g.ID, ServiceID: g.ServiceID, OrgID: g.OrgID,
		Name: g.Name, Version: g.Version, Label: g.Label, Protocol: g.Protocol,
		SpecKey: g.SpecKey, SpecHash: g.SpecHash,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func APIGroupVersionToModel(orgID string, v uigraphapi.APIGroupVersion) *model.APIGroupVersion {
	return &model.APIGroupVersion{
		ID: v.ID, OrgID: orgID, APIGroupID: v.APIGroupID, VersionNumber: v.VersionNumber,
		Label: v.Label, SpecKey: v.SpecKey, SpecHash: v.SpecHash,
		IsAutoVersion: v.IsAutoVersion, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func ServiceDocToModel(d *uigraphapi.ServiceDoc) *model.ServiceDoc {
	return &model.ServiceDoc{
		ID: d.ID, ServiceID: d.ServiceID, OrgID: d.OrgID,
		FileKey: d.FileKey, FileName: d.FileName, FileType: d.FileType,
		Description: d.Description, ContentHash: d.ContentHash,
		CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func ServiceDiagramToModel(d *uigraphapi.ServiceDiagram) *model.ServiceDiagram {
	out := &model.ServiceDiagram{
		ServiceID: d.ServiceID,
		DiagramID: d.DiagramID,
		OrgID:     d.OrgID,
		CreatedBy: d.CreatedBy,
		UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
	if d.Diagram != nil {
		out.Diagram = DiagramToModel(d.Diagram)
	}
	return out
}

func ServiceDBToModel(d *uigraphapi.ServiceDB) *model.ServiceDb {
	return &model.ServiceDb{
		ID: d.ID, ServiceID: d.ServiceID, OrgID: d.OrgID,
		DbName: d.DBName, DbType: d.DBType, Dialect: d.Dialect,
		SchemaJSON: RawStr(d.SchemaJSON),
		Source:     d.Source, SourceTs: d.SourceTS,
		CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func ServiceDBVersionToModel(orgID string, v uigraphapi.ServiceDBVersion) *model.ServiceDBVersion {
	return &model.ServiceDBVersion{
		ID: v.ID, OrgID: orgID, ServiceDbID: v.ServiceDBID, VersionNumber: v.VersionNumber,
		Label: v.Label, SchemaJSON: RawStr(v.SchemaJSON),
		Source: v.Source, SourceTs: v.SourceTS,
		IsAutoVersion: v.IsAutoVersion, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func APIEndpointToModel(e *uigraphapi.APIEndpoint) *model.APIEndpoint {
	return &model.APIEndpoint{
		ID: e.ID, APIGroupID: e.APIGroupID, ServiceID: e.ServiceID, OrgID: e.OrgID,
		OperationID: e.OperationID, Method: e.Method, Path: e.Path,
		Summary: e.Summary, Description: e.Description, Tags: e.Tags,
		Parameters:  RawArrStr(e.Parameters),
		RequestBody: RawStr(e.RequestBody),
		Responses:   RawStr(e.Responses),
		Order:       e.Order,
		CreatedBy:   e.CreatedBy, UpdatedBy: e.UpdatedBy, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

func ServicesToModel(services []uigraphapi.Service) []*model.Service {
	out := make([]*model.Service, len(services))
	for i := range services {
		out[i] = ServiceToModel(&services[i])
	}
	return out
}

func ServiceStatsListToModel(stats []uigraphapi.ServiceStats) []*model.ServiceStats {
	out := make([]*model.ServiceStats, len(stats))
	for i, s := range stats {
		out[i] = ServiceStatsToModel(s)
	}
	return out
}

func APIGroupsToModel(groups []uigraphapi.APIGroup) []*model.APIGroup {
	out := make([]*model.APIGroup, len(groups))
	for i := range groups {
		out[i] = APIGroupToModel(&groups[i])
	}
	return out
}

func APIGroupVersionsToModel(orgID string, versions []uigraphapi.APIGroupVersion) []*model.APIGroupVersion {
	out := make([]*model.APIGroupVersion, len(versions))
	for i, v := range versions {
		out[i] = APIGroupVersionToModel(orgID, v)
	}
	return out
}

func ServiceDocsToModel(docs []uigraphapi.ServiceDoc) []*model.ServiceDoc {
	out := make([]*model.ServiceDoc, len(docs))
	for i := range docs {
		out[i] = ServiceDocToModel(&docs[i])
	}
	return out
}

func ServiceDiagramsToModel(diagrams []uigraphapi.ServiceDiagram) []*model.ServiceDiagram {
	out := make([]*model.ServiceDiagram, len(diagrams))
	for i := range diagrams {
		out[i] = ServiceDiagramToModel(&diagrams[i])
	}
	return out
}

func ServiceDBsToModel(dbs []uigraphapi.ServiceDB) []*model.ServiceDb {
	out := make([]*model.ServiceDb, len(dbs))
	for i := range dbs {
		out[i] = ServiceDBToModel(&dbs[i])
	}
	return out
}

func ServiceDBVersionsToModel(orgID string, versions []uigraphapi.ServiceDBVersion) []*model.ServiceDBVersion {
	out := make([]*model.ServiceDBVersion, len(versions))
	for i, v := range versions {
		out[i] = ServiceDBVersionToModel(orgID, v)
	}
	return out
}

func APIEndpointsToModel(endpoints []uigraphapi.APIEndpoint) []*model.APIEndpoint {
	out := make([]*model.APIEndpoint, len(endpoints))
	for i := range endpoints {
		out[i] = APIEndpointToModel(&endpoints[i])
	}
	return out
}
```

- [ ] **Step 10: Create `internal/graph/convert/testpack.go`**

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func KeyValueToModel(v uigraphapi.KeyValue) *model.KeyValue {
	return &model.KeyValue{Key: v.Key, Value: v.Value}
}

func AssertionToModel(a uigraphapi.Assertion) *model.Assertion {
	return &model.Assertion{Field: a.Field, Type: a.Type, Value: a.Value}
}

func AuthConfigToModel(a *uigraphapi.AuthConfig) *model.AuthConfig {
	if a == nil {
		return nil
	}
	return &model.AuthConfig{
		Type:          a.Type,
		BearerToken:   a.BearerToken,
		APIKeyHeader:  a.APIKeyHeader,
		APIKeyValue:   a.APIKeyValue,
		BasicUsername: a.BasicUsername,
		BasicPassword: a.BasicPassword,
	}
}

func TestCaseStepToModel(s uigraphapi.TestCaseStep) *model.TestCaseStep {
	return &model.TestCaseStep{Order: s.Order, Action: s.Action, ExpectedResult: s.ExpectedResult}
}

func ManualTestCaseToModel(m *uigraphapi.ManualTestCase) *model.ManualTestCase {
	if m == nil {
		return nil
	}
	steps := make([]*model.TestCaseStep, len(m.Steps))
	for i, s := range m.Steps {
		steps[i] = TestCaseStepToModel(s)
	}
	return &model.ManualTestCase{
		Preconditions:   m.Preconditions,
		TestData:        m.TestData,
		Steps:           steps,
		ExpectedOutcome: m.ExpectedOutcome,
		Postconditions:  m.Postconditions,
	}
}

func APITestCaseToModel(a *uigraphapi.APITestCase) *model.APITestCase {
	if a == nil {
		return nil
	}
	headers := make([]*model.KeyValue, len(a.RequestHeaders))
	for i, v := range a.RequestHeaders {
		headers[i] = KeyValueToModel(v)
	}
	params := make([]*model.KeyValue, len(a.QueryParams))
	for i, v := range a.QueryParams {
		params[i] = KeyValueToModel(v)
	}
	assertions := make([]*model.Assertion, len(a.Assertions))
	for i, v := range a.Assertions {
		assertions[i] = AssertionToModel(v)
	}
	return &model.APITestCase{
		HTTPMethod:         a.HTTPMethod,
		APISpecID:          a.APISpecID,
		OperationID:        a.OperationID,
		Auth:               AuthConfigToModel(a.Auth),
		RequestHeaders:     headers,
		QueryParams:        params,
		RequestBody:        a.RequestBody,
		ExpectedStatusCode: a.ExpectedStatusCode,
		MaxResponseTimeMs:  a.MaxResponseTimeMs,
		ResponseBody:       a.ResponseBody,
		Assertions:         assertions,
	}
}

func GraphQLTestCaseToModel(g *uigraphapi.GraphQLTestCase) *model.GraphQLTestCase {
	if g == nil {
		return nil
	}
	assertions := make([]*model.Assertion, len(g.Assertions))
	for i, v := range g.Assertions {
		assertions[i] = AssertionToModel(v)
	}
	return &model.GraphQLTestCase{
		OperationType: g.OperationType,
		OperationName: g.OperationName,
		Query:         g.Query,
		Variables:     g.Variables,
		ResponseBody:  g.ResponseBody,
		Assertions:    assertions,
		ExpectError:   g.ExpectError,
	}
}

func DatabaseTestCaseToModel(d *uigraphapi.DatabaseTestCase) *model.DatabaseTestCase {
	if d == nil {
		return nil
	}
	assertions := make([]*model.Assertion, len(d.Assertions))
	for i, v := range d.Assertions {
		assertions[i] = AssertionToModel(v)
	}
	return &model.DatabaseTestCase{
		Dialect:       d.Dialect,
		SchemaID:      d.SchemaID,
		Query:         d.Query,
		Assertions:    assertions,
		SetupQuery:    d.SetupQuery,
		TeardownQuery: d.TeardownQuery,
	}
}

func GRPCTestCaseToModel(g *uigraphapi.GRPCTestCase) *model.GRPCTestCase {
	if g == nil {
		return nil
	}
	metadata := make([]*model.KeyValue, len(g.Metadata))
	for i, v := range g.Metadata {
		metadata[i] = KeyValueToModel(v)
	}
	assertions := make([]*model.Assertion, len(g.Assertions))
	for i, v := range g.Assertions {
		assertions[i] = AssertionToModel(v)
	}
	return &model.GRPCTestCase{
		ServiceName:    g.ServiceName,
		MethodName:     g.MethodName,
		CallMode:       g.CallMode,
		ProtoFileID:    g.ProtoFileID,
		ServerAddress:  g.ServerAddress,
		RequestMessage: g.RequestMessage,
		Metadata:       metadata,
		ExpectedStatus: g.ExpectedStatus,
		DeadlineMs:     g.DeadlineMs,
		ResponseBody:   g.ResponseBody,
		Assertions:     assertions,
		UseTLS:         g.UseTLS,
		ExpectError:    g.ExpectError,
	}
}

func TestPackToModel(p *uigraphapi.TestPack) *model.TestPack {
	return &model.TestPack{
		TestPackID: p.TestPackID, ServiceID: p.ServiceID, OrgID: p.OrgID,
		Name: p.Name, Type: p.Type,
		CreatedBy: p.CreatedBy, UpdatedBy: p.UpdatedBy, DeletedBy: p.DeletedBy,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, DeletedAt: p.DeletedAt,
	}
}

func TestCaseToModel(tc *uigraphapi.TestCase) *model.TestCase {
	return &model.TestCase{
		TestCaseID: tc.TestCaseID, TestPackID: tc.TestPackID, ServiceID: tc.ServiceID, OrgID: tc.OrgID,
		Title: tc.Title, Order: tc.Order, Type: tc.Type, Description: tc.Description, Priority: tc.Priority,
		Labels: tc.Labels, LinkedTicket: tc.LinkedTicket, EstimatedDurationMins: tc.EstimatedDurationMins,
		TestOwner: tc.TestOwner, LinkedMapNodeID: tc.LinkedMapNodeID, IsCritical: tc.IsCritical, EvidenceRequired: tc.EvidenceRequired,
		Manual: ManualTestCaseToModel(tc.Manual), API: APITestCaseToModel(tc.API),
		Graphql: GraphQLTestCaseToModel(tc.GraphQL), Database: DatabaseTestCaseToModel(tc.Database), Grpc: GRPCTestCaseToModel(tc.GRPC),
		Status: tc.Status, Version: tc.Version, BaselineRunResultID: tc.BaselineRunResultID, Dependencies: tc.Dependencies,
		CreatedBy: tc.CreatedBy, UpdatedBy: tc.UpdatedBy, DeletedBy: tc.DeletedBy, CreatedAt: tc.CreatedAt, UpdatedAt: tc.UpdatedAt, DeletedAt: tc.DeletedAt,
	}
}

func TestRunToModel(tr *uigraphapi.TestRun) *model.TestRun {
	return &model.TestRun{
		TestRunID: tr.TestRunID, TestPackID: tr.TestPackID, ServiceID: tr.ServiceID, OrgID: tr.OrgID,
		Environment: tr.Environment, ReleaseLabel: tr.ReleaseLabel, StartedAt: tr.StartedAt, CompletedAt: tr.CompletedAt,
		Status: tr.Status, StartedBy: tr.StartedBy, ExecutedBy: tr.ExecutedBy, ExecutedAt: tr.ExecutedAt, OverallStatus: tr.OverallStatus,
	}
}

func TestRunSummaryToModel(s uigraphapi.TestRunSummary) *model.TestRunSummary {
	return &model.TestRunSummary{
		TestRunID: s.TestRunID, TestPackID: s.TestPackID, ServiceID: s.ServiceID,
		Environment: s.Environment, ReleaseLabel: s.ReleaseLabel, StartedAt: s.StartedAt, CompletedAt: s.CompletedAt,
		Status: s.Status, StartedBy: s.StartedBy, ExecutedBy: s.ExecutedBy, ExecutedAt: s.ExecutedAt, OverallStatus: s.OverallStatus,
		PassedCount: s.PassedCount, FailedCount: s.FailedCount, SkippedCount: s.SkippedCount, BlockedCount: s.BlockedCount,
	}
}

func TestRunResultToModel(rr *uigraphapi.TestRunResult) *model.TestRunResult {
	var responseTimeMs *int
	if rr.ResponseTimeMs != nil {
		v := int(*rr.ResponseTimeMs)
		responseTimeMs = &v
	}
	return &model.TestRunResult{
		TestRunResultID: rr.TestRunResultID, TestRunID: rr.TestRunID, TestCaseID: rr.TestCaseID,
		ServiceID: rr.ServiceID, OrgID: rr.OrgID, Status: rr.Status, BlockedReason: rr.BlockedReason,
		ResponseStatus: rr.ResponseStatus, ResponseBody: rr.ResponseBody, ResponseTimeMs: responseTimeMs,
		Notes: rr.Notes, ScreenshotUrls: rr.ScreenshotURLs, ExecutedAt: rr.ExecutedAt, ExecutedBy: rr.ExecutedBy,
	}
}

func TestPacksToModel(packs []uigraphapi.TestPack) []*model.TestPack {
	out := make([]*model.TestPack, len(packs))
	for i := range packs {
		out[i] = TestPackToModel(&packs[i])
	}
	return out
}

func TestCasesToModel(cases []uigraphapi.TestCase) []*model.TestCase {
	out := make([]*model.TestCase, len(cases))
	for i := range cases {
		out[i] = TestCaseToModel(&cases[i])
	}
	return out
}

func TestRunsToModel(runs []uigraphapi.TestRun) []*model.TestRun {
	out := make([]*model.TestRun, len(runs))
	for i := range runs {
		out[i] = TestRunToModel(&runs[i])
	}
	return out
}

func TestRunSummariesToModel(summaries []uigraphapi.TestRunSummary) []*model.TestRunSummary {
	out := make([]*model.TestRunSummary, len(summaries))
	for i, s := range summaries {
		out[i] = TestRunSummaryToModel(s)
	}
	return out
}

func TestRunResultsToModel(results []uigraphapi.TestRunResult) []*model.TestRunResult {
	out := make([]*model.TestRunResult, len(results))
	for i := range results {
		out[i] = TestRunResultToModel(&results[i])
	}
	return out
}
```

- [ ] **Step 11: Create `internal/graph/refs.go`**

```go
package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/model"
)

// resolveActor resolves a single created_by / updated_by id within an org to
// its public actor info, returning nil when id is empty or matches no actor.
func (r *Resolver) resolveActor(ctx context.Context, orgID, id string) (*model.Actor, error) {
	if id == "" {
		return nil, nil
	}
	actors, err := r.Client.ResolveActors(ctx, orgID, []string{id})
	if err != nil {
		return nil, err
	}
	a := actors[id]
	if a == nil {
		return nil, nil
	}
	m := &model.Actor{ID: a.ID, Type: a.Type, Name: a.Name, Disabled: a.Disabled}
	if a.Email != "" {
		m.Email = &a.Email
	}
	if a.AvatarURL != "" {
		m.AvatarURL = &a.AvatarURL
	}
	return m, nil
}

// resolveAssetURL resolves a single asset id within an org to a presigned GET
// URL, returning nil when id is empty or no url is produced.
func (r *Resolver) resolveAssetURL(ctx context.Context, orgID, assetID string) (*string, error) {
	if assetID == "" {
		return nil, nil
	}
	urls, err := r.Client.ResolveAssetURLs(ctx, orgID, []string{assetID})
	if err != nil {
		return nil, err
	}
	u, ok := urls[assetID]
	if !ok || u == "" {
		return nil, nil
	}
	return &u, nil
}
```

- [ ] **Step 12: Delete the old `convert.go` and `helpers.go`**

```bash
rm internal/graph/convert.go internal/graph/helpers.go
```

- [ ] **Step 13: Rename every call site across `internal/graph/*.resolvers.go`**

Run this script from the repo root. It rewrites every old unqualified call to its new `convert.`-qualified name, across every resolver file (the glob does not recurse into `generated/`, `model/`, or `convert/`):

```bash
declare -a RENAMES=(
  "toMap(:convert.ToMap("
  "rawStr(:convert.RawStr("
  "rawArrStr(:convert.RawArrStr("
  "strFromMap(:convert.StrFromMap("
  "boolFromMap(:convert.BoolFromMap("
  "optStrFromMap(:convert.OptStrFromMap("
  "unmarshalJSONString(:convert.UnmarshalJSONString("
  "meToModel(:convert.MeToModel("
  "orgSummaryToModel(:convert.OrgSummaryToModel("
  "orgSummariesToModel(:convert.OrgSummariesToModel("
  "orgToModel(:convert.OrgToModel("
  "memberToModel(:convert.MemberToModel("
  "teamToModel(:convert.TeamToModel("
  "teamMemberToModel(:convert.TeamMemberToModel("
  "invitationToModel(:convert.InvitationToModel("
  "serviceAccountToModel(:convert.ServiceAccountToModel("
  "saTokenToModel(:convert.SATokenToModel("
  "createdTokenToModel(:convert.CreatedTokenToModel("
  "orgsToModel(:convert.OrgsToModel("
  "membersToModel(:convert.MembersToModel("
  "teamsToModel(:convert.TeamsToModel("
  "teamMembersToModel(:convert.TeamMembersToModel("
  "invitationsToModel(:convert.InvitationsToModel("
  "serviceAccountsToModel(:convert.ServiceAccountsToModel("
  "saTokensToModel(:convert.SATokensToModel("
  "folderToModel(:convert.FolderToModel("
  "foldersToModel(:convert.FoldersToModel("
  "diagramToModel(:convert.DiagramToModel("
  "diagramsToModel(:convert.DiagramsToModel("
  "diagramVersionToModel(:convert.DiagramVersionToModel("
  "diagramVersionsToModel(:convert.DiagramVersionsToModel("
  "diagramImageToModel(:convert.DiagramImageToModel("
  "diagramImagesToModel(:convert.DiagramImagesToModel("
  "flowComponentFieldToModel(:convert.FlowComponentFieldToModel("
  "flowComponentToModel(:convert.FlowComponentToModel("
  "flowComponentsToModel(:convert.FlowComponentsToModel("
  "componentFieldToModel(:convert.ComponentFieldToModel("
  "componentToModel(:convert.ComponentToModel("
  "componentsToModel(:convert.ComponentsToModel("
  "uimapToModel(:convert.UIMapToModel("
  "uimapsToModel(:convert.UIMapsToModel("
  "frameToModel(:convert.FrameToModel("
  "framesToModel(:convert.FramesToModel("
  "focalPointToModel(:convert.FocalPointToModel("
  "focalPointsToModel(:convert.FocalPointsToModel("
  "canvasToModel(:convert.CanvasToModel("
  "frameGroupToModel(:convert.FrameGroupToModel("
  "frameGroupsToModel(:convert.FrameGroupsToModel("
  "frameLinkToModel(:convert.FrameLinkToModel("
  "frameLinksToModel(:convert.FrameLinksToModel("
  "focalPointMetaToModel(:convert.FocalPointMetaToModel("
  "focalPointMetasToModel(:convert.FocalPointMetasToModel("
  "focalPointMetaBody(:convert.FocalPointMetaBody("
  "serviceToModel(:convert.ServiceToModel("
  "servicesToModel(:convert.ServicesToModel("
  "serviceStatsToModel(:convert.ServiceStatsToModel("
  "serviceStatsListToModel(:convert.ServiceStatsListToModel("
  "apiGroupToModel(:convert.APIGroupToModel("
  "apiGroupsToModel(:convert.APIGroupsToModel("
  "apiGroupVersionToModel(:convert.APIGroupVersionToModel("
  "apiGroupVersionsToModel(:convert.APIGroupVersionsToModel("
  "serviceDocToModel(:convert.ServiceDocToModel("
  "serviceDocsToModel(:convert.ServiceDocsToModel("
  "serviceDiagramToModel(:convert.ServiceDiagramToModel("
  "serviceDiagramsToModel(:convert.ServiceDiagramsToModel("
  "serviceDBToModel(:convert.ServiceDBToModel("
  "serviceDBsToModel(:convert.ServiceDBsToModel("
  "serviceDBVersionToModel(:convert.ServiceDBVersionToModel("
  "serviceDBVersionsToModel(:convert.ServiceDBVersionsToModel("
  "apiEndpointToModel(:convert.APIEndpointToModel("
  "apiEndpointsToModel(:convert.APIEndpointsToModel("
  "keyValueToModel(:convert.KeyValueToModel("
  "assertionToModel(:convert.AssertionToModel("
  "authConfigToModel(:convert.AuthConfigToModel("
  "testCaseStepToModel(:convert.TestCaseStepToModel("
  "manualTestCaseToModel(:convert.ManualTestCaseToModel("
  "apiTestCaseToModel(:convert.APITestCaseToModel("
  "graphQLTestCaseToModel(:convert.GraphQLTestCaseToModel("
  "databaseTestCaseToModel(:convert.DatabaseTestCaseToModel("
  "grpcTestCaseToModel(:convert.GRPCTestCaseToModel("
  "testPackToModel(:convert.TestPackToModel("
  "testPacksToModel(:convert.TestPacksToModel("
  "testCaseToModel(:convert.TestCaseToModel("
  "testCasesToModel(:convert.TestCasesToModel("
  "testRunToModel(:convert.TestRunToModel("
  "testRunsToModel(:convert.TestRunsToModel("
  "testRunSummaryToModel(:convert.TestRunSummaryToModel("
  "testRunSummariesToModel(:convert.TestRunSummariesToModel("
  "testRunResultToModel(:convert.TestRunResultToModel("
  "testRunResultsToModel(:convert.TestRunResultsToModel("
  "userToModel(:convert.UserToModel("
  "usersToModel(:convert.UsersToModel("
  "oauthProviderToModel(:convert.OAuthProviderToModel("
  "oauthProvidersToModel(:convert.OAuthProvidersToModel("
  "roleMappingToModel(:convert.RoleMappingToModel("
  "roleMappingsToModel(:convert.RoleMappingsToModel("
  "ldapToModel(:convert.LDAPToModel("
  "samlToModel(:convert.SAMLToModel("
)

for f in internal/graph/*.resolvers.go; do
  for pair in "${RENAMES[@]}"; do
    old="${pair%%:*}"
    new="${pair#*:}"
    sed -i '' "s/${old}/${new}/g" "$f"
  done
done
```

(On Linux, drop the `''` after `-i` in the `sed` line. List order matters only in that every entry is processed for every file — since each pattern includes the function's exact name plus `(`, and no old name is a substring of another old name's match text, ordering does not cause double-rewrites.)

- [ ] **Step 14: Add the `convert` import to resolver files that use it**

Add this import line (alongside the existing `"github.com/uigraph/graphql/internal/graph/model"` import) to each of: `auth.resolvers.go`, `admin.resolvers.go`, `org.resolvers.go`, `folder.resolvers.go`, `diagram.resolvers.go`, `component.resolvers.go`, `uimap.resolvers.go`, `catalog.resolvers.go`, `testpack.resolvers.go`:

```go
	"github.com/uigraph/graphql/internal/graph/convert"
```

- [ ] **Step 15: Build and fix remaining import errors**

Run: `go build ./...`

If a file reports `"convert" imported and not used`, remove the import line you just added to that specific file (it had no `convert.X` call after the rename). If a file reports `undefined: convert`, you missed adding the import in Step 14 — add it. Repeat until `go build ./...` exits 0 with no output.

- [ ] **Step 16: Commit**

```bash
git add -A
git commit -m "refactor: extract convert.go into internal/graph/convert package"
```

---

## Task 10: Narrow `Resolver`'s client dependency into per-domain interfaces

**Files:**
- Modify: `internal/graph/resolver.go`, `internal/graph/refs.go`, every `internal/graph/*.resolvers.go`, `cmd/server/main.go`

**Interfaces:**
- Consumes: every exported method on `*uigraphapi.Client` (Tasks 2/5/6/7).
- Produces: `Resolver` struct with 10 narrow interface fields (`Auth`, `Org`, `Admin`, `Folder`, `Diagram`, `Component`, `UIMap`, `Catalog`, `TestPack`, `Actor`) instead of one `Client *uigraphapi.Client` field. `cmd/server/main.go` wires the same single `*uigraphapi.Client` instance into all 10 fields — `*uigraphapi.Client` satisfies every interface automatically since it already has all these methods.

- [ ] **Step 1: Rewrite `internal/graph/resolver.go`**

```go
package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

type authClient interface {
	Me(ctx context.Context) (*uigraphapi.MeResponse, error)
	MyOrgs(ctx context.Context) ([]uigraphapi.OrgSummary, error)
	SwitchOrg(ctx context.Context, orgID string) error
}

type orgClient interface {
	ListOrgs(ctx context.Context) ([]uigraphapi.Org, error)
	GetOrg(ctx context.Context, id string) (*uigraphapi.Org, error)
	CreateOrg(ctx context.Context, body map[string]interface{}) (*uigraphapi.Org, error)
	UpdateOrg(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.Org, error)
	DeleteOrg(ctx context.Context, id string) error
	ListMembers(ctx context.Context, orgID string) ([]uigraphapi.Member, error)
	AddMember(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Member, error)
	UpdateMemberRole(ctx context.Context, orgID, userID string, body map[string]interface{}) (*uigraphapi.Member, error)
	RemoveMember(ctx context.Context, orgID, userID string) error
	ListTeams(ctx context.Context, orgID string) ([]uigraphapi.Team, error)
	GetTeam(ctx context.Context, orgID, teamID string) (*uigraphapi.Team, error)
	CreateTeam(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Team, error)
	UpdateTeam(ctx context.Context, orgID, teamID string, body map[string]interface{}) (*uigraphapi.Team, error)
	DeleteTeam(ctx context.Context, orgID, teamID string) error
	ListTeamMembers(ctx context.Context, orgID, teamID string) ([]uigraphapi.TeamMember, error)
	AddTeamMember(ctx context.Context, orgID, teamID string, body map[string]interface{}) error
	RemoveTeamMember(ctx context.Context, orgID, teamID, userID string) error
	ListInvitations(ctx context.Context, orgID string) ([]uigraphapi.Invitation, error)
	CreateInvitation(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Invitation, error)
	RevokeInvitation(ctx context.Context, orgID, invitationID string) error
	ListServiceAccounts(ctx context.Context, orgID string) ([]uigraphapi.ServiceAccount, error)
	GetServiceAccount(ctx context.Context, orgID, id string) (*uigraphapi.ServiceAccount, error)
	CreateServiceAccount(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.ServiceAccount, error)
	UpdateServiceAccount(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, orgID, id string) error
	ListServiceAccountTokens(ctx context.Context, orgID, saID string) ([]uigraphapi.ServiceAccountToken, error)
	CreateServiceAccountToken(ctx context.Context, orgID, saID string, body map[string]interface{}) (*uigraphapi.CreatedToken, error)
	RevokeServiceAccountToken(ctx context.Context, orgID, saID, tokenID string) error
}

type adminClient interface {
	ListUsers(ctx context.Context) ([]uigraphapi.User, error)
	GetUser(ctx context.Context, id string) (*uigraphapi.User, error)
	CreateUser(ctx context.Context, body map[string]interface{}) (*uigraphapi.User, error)
	UpdateUser(ctx context.Context, id string, body map[string]interface{}) (*uigraphapi.User, error)
	DisableUser(ctx context.Context, id string) error
	ListOAuthProviders(ctx context.Context) ([]uigraphapi.OAuthProvider, error)
	UpsertOAuthProvider(ctx context.Context, provider string, body map[string]interface{}) error
	DeleteOAuthProvider(ctx context.Context, provider string) error
	ListRoleMappings(ctx context.Context) ([]uigraphapi.RoleMapping, error)
	CreateRoleMapping(ctx context.Context, body map[string]interface{}) error
	DeleteRoleMapping(ctx context.Context, id string) error
	GetLDAP(ctx context.Context) (*uigraphapi.LDAPConfig, error)
	UpsertLDAP(ctx context.Context, body map[string]interface{}) error
	DeleteLDAP(ctx context.Context) error
	GetSAML(ctx context.Context) (*uigraphapi.SAMLConfig, error)
	UpsertSAML(ctx context.Context, body map[string]interface{}) error
}

type folderClient interface {
	ListFolders(ctx context.Context, orgID, folderType, parentID string) ([]uigraphapi.Folder, error)
	GetFolder(ctx context.Context, orgID, id string) (*uigraphapi.Folder, error)
	CreateFolder(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Folder, error)
	UpdateFolder(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Folder, error)
	DeleteFolder(ctx context.Context, orgID, id string) error
}

type diagramClient interface {
	ListDiagrams(ctx context.Context, orgID, folderID string) ([]uigraphapi.Diagram, error)
	GetDiagram(ctx context.Context, orgID, id string) (*uigraphapi.Diagram, error)
	GetDiagramContent(ctx context.Context, orgID, id string) (string, error)
	CreateDiagram(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Diagram, error)
	UpdateDiagram(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Diagram, error)
	DeleteDiagram(ctx context.Context, orgID, id string) error
	ListDiagramImages(ctx context.Context, orgID, diagramID string) ([]uigraphapi.DiagramImage, error)
	SyncDiagram(ctx context.Context, orgID string, body map[string]interface{}) (map[string]interface{}, error)
	ListDiagramVersions(ctx context.Context, orgID, diagramID string) ([]uigraphapi.DiagramVersion, error)
	CreateDiagramVersion(ctx context.Context, orgID, diagramID string, body map[string]interface{}) (*uigraphapi.DiagramVersion, error)
	GetDiagramVersionContent(ctx context.Context, orgID, diagramID, versionID string) (string, error)
	RestoreDiagramVersion(ctx context.Context, orgID, diagramID, versionID string) (*uigraphapi.Diagram, error)
}

type componentClient interface {
	ListFlowDiagramComponents(ctx context.Context, orgID string) (*uigraphapi.FlowComponents, error)
	ListComponents(ctx context.Context, orgID string) (*uigraphapi.Components, error)
}

type uimapClient interface {
	ListMaps(ctx context.Context, orgID, folderID string) ([]uigraphapi.UIMap, error)
	GetMap(ctx context.Context, orgID, id string) (*uigraphapi.UIMap, error)
	CreateMap(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.UIMap, error)
	UpdateMap(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.UIMap, error)
	DeleteMap(ctx context.Context, orgID, id string) error
	ListFrames(ctx context.Context, orgID, mapID string) ([]uigraphapi.Frame, error)
	GetFrame(ctx context.Context, orgID, mapID, id string) (*uigraphapi.Frame, error)
	GetFrameByID(ctx context.Context, orgID, id string) (*uigraphapi.Frame, error)
	CreateFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*uigraphapi.Frame, error)
	UpdateFrame(ctx context.Context, orgID, mapID, id string, body map[string]interface{}) (*uigraphapi.Frame, error)
	DeleteFrame(ctx context.Context, orgID, mapID, id string) error
	SyncFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (map[string]interface{}, error)
	ListFocalPoints(ctx context.Context, orgID, mapID, frameID string) ([]uigraphapi.FocalPoint, error)
	GetFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) (*uigraphapi.FocalPoint, error)
	CreateFocalPoint(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*uigraphapi.FocalPoint, error)
	UpdateFocalPoint(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*uigraphapi.FocalPoint, error)
	DeleteFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) error
	GetCanvas(ctx context.Context, orgID, mapID string) (*uigraphapi.Canvas, error)
	UpsertCanvas(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*uigraphapi.Canvas, error)
	ListFrameGroups(ctx context.Context, orgID, mapID, frameID string) ([]uigraphapi.FrameGroup, error)
	CreateFrameGroup(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*uigraphapi.FrameGroup, error)
	UpdateFrameGroup(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*uigraphapi.FrameGroup, error)
	DeleteFrameGroup(ctx context.Context, orgID, mapID, frameID, id string) error
	ListFrameLinks(ctx context.Context, orgID, mapID, frameID string) ([]uigraphapi.FrameLink, error)
	CreateFrameLink(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*uigraphapi.FrameLink, error)
	UpdateFrameLink(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*uigraphapi.FrameLink, error)
	DeleteFrameLink(ctx context.Context, orgID, mapID, frameID, id string) error
	ListFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string) ([]uigraphapi.FocalPointMeta, error)
	CreateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string, body map[string]interface{}) (*uigraphapi.FocalPointMeta, error)
	UpdateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string, body map[string]interface{}) (*uigraphapi.FocalPointMeta, error)
	DeleteFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string) error
}

type catalogClient interface {
	ListServices(ctx context.Context, orgID, folderID, teamID string) ([]uigraphapi.Service, error)
	GetService(ctx context.Context, orgID, id string) (*uigraphapi.Service, error)
	CreateService(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Service, error)
	UpdateService(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Service, error)
	DeleteService(ctx context.Context, orgID, id string) error
	ListServiceStats(ctx context.Context, orgID string, serviceID *string) ([]uigraphapi.ServiceStats, error)
	ListAPIGroups(ctx context.Context, orgID, serviceID string) ([]uigraphapi.APIGroup, error)
	GetAPIGroup(ctx context.Context, orgID, serviceID, id string) (*uigraphapi.APIGroup, error)
	CreateAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.APIGroup, error)
	UpdateAPIGroup(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.APIGroup, error)
	DeleteAPIGroup(ctx context.Context, orgID, serviceID, id string) error
	SyncAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (map[string]interface{}, error)
	ListAPIGroupVersions(ctx context.Context, orgID, serviceID, apiGroupID string) ([]uigraphapi.APIGroupVersion, error)
	ListServiceDocs(ctx context.Context, orgID, serviceID string) ([]uigraphapi.ServiceDoc, error)
	GetServiceDoc(ctx context.Context, orgID, serviceID, id string) (*uigraphapi.ServiceDoc, error)
	CreateServiceDoc(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.ServiceDoc, error)
	UpdateServiceDoc(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.ServiceDoc, error)
	DeleteServiceDoc(ctx context.Context, orgID, serviceID, id string) error
	ListServiceDiagrams(ctx context.Context, orgID, serviceID string) ([]uigraphapi.ServiceDiagram, error)
	CreateServiceDiagram(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.ServiceDiagram, error)
	DeleteServiceDiagram(ctx context.Context, orgID, serviceID, diagramID string) error
	ListServiceDBs(ctx context.Context, orgID, serviceID string) ([]uigraphapi.ServiceDB, error)
	GetServiceDB(ctx context.Context, orgID, serviceID, id string) (*uigraphapi.ServiceDB, error)
	CreateServiceDB(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.ServiceDB, error)
	UpdateServiceDB(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.ServiceDB, error)
	DeleteServiceDB(ctx context.Context, orgID, serviceID, id string) error
	ListServiceDBVersions(ctx context.Context, orgID, serviceID, serviceDBID string) ([]uigraphapi.ServiceDBVersion, error)
	CreateServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID string, body map[string]interface{}) (*uigraphapi.ServiceDBVersion, error)
	RestoreServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID, versionID string) (*uigraphapi.ServiceDB, error)
	ListAPIEndpoints(ctx context.Context, orgID, serviceID, apiGroupID string) ([]uigraphapi.APIEndpoint, error)
	GetAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) (*uigraphapi.APIEndpoint, error)
	CreateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID string, body map[string]interface{}) (*uigraphapi.APIEndpoint, error)
	UpdateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string, body map[string]interface{}) (*uigraphapi.APIEndpoint, error)
	DeleteAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) error
}

type testPackClient interface {
	ListTestPacks(ctx context.Context, orgID, serviceID string) ([]uigraphapi.TestPack, error)
	CreateTestPack(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestPack, error)
	UpdateTestPack(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestPack, error)
	DeleteTestPack(ctx context.Context, orgID, serviceID, id string) error
	ListTestCases(ctx context.Context, orgID, serviceID string, testPackID *string) ([]uigraphapi.TestCase, error)
	CreateTestCase(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestCase, error)
	UpdateTestCase(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestCase, error)
	DeleteTestCase(ctx context.Context, orgID, serviceID, id string) error
	GetTestRun(ctx context.Context, orgID, serviceID, id string) (*uigraphapi.TestRun, error)
	ListTestRuns(ctx context.Context, orgID, serviceID string, testPackID *string) ([]uigraphapi.TestRun, error)
	ListTestRunsSummary(ctx context.Context, orgID, serviceID string, testPackID, environment, status, executedBy *string, fromDate, toDate *time.Time) ([]uigraphapi.TestRunSummary, error)
	CreateTestRun(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestRun, error)
	UpdateTestRun(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestRun, error)
	ListTestRunResults(ctx context.Context, orgID, serviceID, testRunID string) ([]uigraphapi.TestRunResult, error)
	CreateTestRunResult(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*uigraphapi.TestRunResult, error)
	UpdateTestRunResult(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*uigraphapi.TestRunResult, error)
}

type actorClient interface {
	ResolveActors(ctx context.Context, orgID string, ids []string) (map[string]*uigraphapi.Actor, error)
	ResolveAssetURLs(ctx context.Context, orgID string, ids []string) (map[string]string, error)
}

// Resolver is the root dependency-injection struct for all resolvers. Each
// field is the minimal interface its domain's resolvers need — not the full
// *uigraphapi.Client — so tests can inject a narrow fake instead of mocking
// every REST method.
type Resolver struct {
	Auth      authClient
	Org       orgClient
	Admin     adminClient
	Folder    folderClient
	Diagram   diagramClient
	Component componentClient
	UIMap     uimapClient
	Catalog   catalogClient
	TestPack  testPackClient
	Actor     actorClient
}
```

`ListTestRunsSummary` needs `time.Time`, so add `"time"` to the import block:

```go
import (
	"context"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)
```

- [ ] **Step 2: Update `internal/graph/refs.go`**

Change both `r.Client.ResolveActors(...)` and `r.Client.ResolveAssetURLs(...)` to `r.Actor.ResolveActors(...)` and `r.Actor.ResolveAssetURLs(...)`.

- [ ] **Step 3: Rename every `r.Client.<Method>(` call site**

Run this script from the repo root:

```bash
declare -a CLIENT_RENAMES=(
  "r.Client.Me(:r.Auth.Me("
  "r.Client.MyOrgs(:r.Auth.MyOrgs("
  "r.Client.SwitchOrg(:r.Auth.SwitchOrg("
  "r.Client.ListOrgs(:r.Org.ListOrgs("
  "r.Client.GetOrg(:r.Org.GetOrg("
  "r.Client.CreateOrg(:r.Org.CreateOrg("
  "r.Client.UpdateOrg(:r.Org.UpdateOrg("
  "r.Client.DeleteOrg(:r.Org.DeleteOrg("
  "r.Client.ListMembers(:r.Org.ListMembers("
  "r.Client.AddMember(:r.Org.AddMember("
  "r.Client.UpdateMemberRole(:r.Org.UpdateMemberRole("
  "r.Client.RemoveMember(:r.Org.RemoveMember("
  "r.Client.ListTeams(:r.Org.ListTeams("
  "r.Client.GetTeam(:r.Org.GetTeam("
  "r.Client.CreateTeam(:r.Org.CreateTeam("
  "r.Client.UpdateTeam(:r.Org.UpdateTeam("
  "r.Client.DeleteTeam(:r.Org.DeleteTeam("
  "r.Client.ListTeamMembers(:r.Org.ListTeamMembers("
  "r.Client.AddTeamMember(:r.Org.AddTeamMember("
  "r.Client.RemoveTeamMember(:r.Org.RemoveTeamMember("
  "r.Client.ListInvitations(:r.Org.ListInvitations("
  "r.Client.CreateInvitation(:r.Org.CreateInvitation("
  "r.Client.RevokeInvitation(:r.Org.RevokeInvitation("
  "r.Client.ListServiceAccounts(:r.Org.ListServiceAccounts("
  "r.Client.GetServiceAccount(:r.Org.GetServiceAccount("
  "r.Client.CreateServiceAccount(:r.Org.CreateServiceAccount("
  "r.Client.UpdateServiceAccount(:r.Org.UpdateServiceAccount("
  "r.Client.DeleteServiceAccount(:r.Org.DeleteServiceAccount("
  "r.Client.ListServiceAccountTokens(:r.Org.ListServiceAccountTokens("
  "r.Client.CreateServiceAccountToken(:r.Org.CreateServiceAccountToken("
  "r.Client.RevokeServiceAccountToken(:r.Org.RevokeServiceAccountToken("
  "r.Client.ListUsers(:r.Admin.ListUsers("
  "r.Client.GetUser(:r.Admin.GetUser("
  "r.Client.CreateUser(:r.Admin.CreateUser("
  "r.Client.UpdateUser(:r.Admin.UpdateUser("
  "r.Client.DisableUser(:r.Admin.DisableUser("
  "r.Client.ListOAuthProviders(:r.Admin.ListOAuthProviders("
  "r.Client.UpsertOAuthProvider(:r.Admin.UpsertOAuthProvider("
  "r.Client.DeleteOAuthProvider(:r.Admin.DeleteOAuthProvider("
  "r.Client.ListRoleMappings(:r.Admin.ListRoleMappings("
  "r.Client.CreateRoleMapping(:r.Admin.CreateRoleMapping("
  "r.Client.DeleteRoleMapping(:r.Admin.DeleteRoleMapping("
  "r.Client.GetLDAP(:r.Admin.GetLDAP("
  "r.Client.UpsertLDAP(:r.Admin.UpsertLDAP("
  "r.Client.DeleteLDAP(:r.Admin.DeleteLDAP("
  "r.Client.GetSAML(:r.Admin.GetSAML("
  "r.Client.UpsertSAML(:r.Admin.UpsertSAML("
  "r.Client.ListFolders(:r.Folder.ListFolders("
  "r.Client.GetFolder(:r.Folder.GetFolder("
  "r.Client.CreateFolder(:r.Folder.CreateFolder("
  "r.Client.UpdateFolder(:r.Folder.UpdateFolder("
  "r.Client.DeleteFolder(:r.Folder.DeleteFolder("
  "r.Client.ListDiagrams(:r.Diagram.ListDiagrams("
  "r.Client.GetDiagram(:r.Diagram.GetDiagram("
  "r.Client.GetDiagramContent(:r.Diagram.GetDiagramContent("
  "r.Client.CreateDiagram(:r.Diagram.CreateDiagram("
  "r.Client.UpdateDiagram(:r.Diagram.UpdateDiagram("
  "r.Client.DeleteDiagram(:r.Diagram.DeleteDiagram("
  "r.Client.ListDiagramImages(:r.Diagram.ListDiagramImages("
  "r.Client.SyncDiagram(:r.Diagram.SyncDiagram("
  "r.Client.ListDiagramVersions(:r.Diagram.ListDiagramVersions("
  "r.Client.CreateDiagramVersion(:r.Diagram.CreateDiagramVersion("
  "r.Client.GetDiagramVersionContent(:r.Diagram.GetDiagramVersionContent("
  "r.Client.RestoreDiagramVersion(:r.Diagram.RestoreDiagramVersion("
  "r.Client.ListFlowDiagramComponents(:r.Component.ListFlowDiagramComponents("
  "r.Client.ListComponents(:r.Component.ListComponents("
  "r.Client.ListMaps(:r.UIMap.ListMaps("
  "r.Client.GetMap(:r.UIMap.GetMap("
  "r.Client.CreateMap(:r.UIMap.CreateMap("
  "r.Client.UpdateMap(:r.UIMap.UpdateMap("
  "r.Client.DeleteMap(:r.UIMap.DeleteMap("
  "r.Client.ListFrames(:r.UIMap.ListFrames("
  "r.Client.GetFrame(:r.UIMap.GetFrame("
  "r.Client.GetFrameByID(:r.UIMap.GetFrameByID("
  "r.Client.CreateFrame(:r.UIMap.CreateFrame("
  "r.Client.UpdateFrame(:r.UIMap.UpdateFrame("
  "r.Client.DeleteFrame(:r.UIMap.DeleteFrame("
  "r.Client.SyncFrame(:r.UIMap.SyncFrame("
  "r.Client.ListFocalPoints(:r.UIMap.ListFocalPoints("
  "r.Client.GetFocalPoint(:r.UIMap.GetFocalPoint("
  "r.Client.CreateFocalPoint(:r.UIMap.CreateFocalPoint("
  "r.Client.UpdateFocalPoint(:r.UIMap.UpdateFocalPoint("
  "r.Client.DeleteFocalPoint(:r.UIMap.DeleteFocalPoint("
  "r.Client.GetCanvas(:r.UIMap.GetCanvas("
  "r.Client.UpsertCanvas(:r.UIMap.UpsertCanvas("
  "r.Client.ListFrameGroups(:r.UIMap.ListFrameGroups("
  "r.Client.CreateFrameGroup(:r.UIMap.CreateFrameGroup("
  "r.Client.UpdateFrameGroup(:r.UIMap.UpdateFrameGroup("
  "r.Client.DeleteFrameGroup(:r.UIMap.DeleteFrameGroup("
  "r.Client.ListFrameLinks(:r.UIMap.ListFrameLinks("
  "r.Client.CreateFrameLink(:r.UIMap.CreateFrameLink("
  "r.Client.UpdateFrameLink(:r.UIMap.UpdateFrameLink("
  "r.Client.DeleteFrameLink(:r.UIMap.DeleteFrameLink("
  "r.Client.ListFocalPointMeta(:r.UIMap.ListFocalPointMeta("
  "r.Client.CreateFocalPointMeta(:r.UIMap.CreateFocalPointMeta("
  "r.Client.UpdateFocalPointMeta(:r.UIMap.UpdateFocalPointMeta("
  "r.Client.DeleteFocalPointMeta(:r.UIMap.DeleteFocalPointMeta("
  "r.Client.ListServices(:r.Catalog.ListServices("
  "r.Client.GetService(:r.Catalog.GetService("
  "r.Client.CreateService(:r.Catalog.CreateService("
  "r.Client.UpdateService(:r.Catalog.UpdateService("
  "r.Client.DeleteService(:r.Catalog.DeleteService("
  "r.Client.ListServiceStats(:r.Catalog.ListServiceStats("
  "r.Client.ListAPIGroups(:r.Catalog.ListAPIGroups("
  "r.Client.GetAPIGroup(:r.Catalog.GetAPIGroup("
  "r.Client.CreateAPIGroup(:r.Catalog.CreateAPIGroup("
  "r.Client.UpdateAPIGroup(:r.Catalog.UpdateAPIGroup("
  "r.Client.DeleteAPIGroup(:r.Catalog.DeleteAPIGroup("
  "r.Client.SyncAPIGroup(:r.Catalog.SyncAPIGroup("
  "r.Client.ListAPIGroupVersions(:r.Catalog.ListAPIGroupVersions("
  "r.Client.ListServiceDocs(:r.Catalog.ListServiceDocs("
  "r.Client.GetServiceDoc(:r.Catalog.GetServiceDoc("
  "r.Client.CreateServiceDoc(:r.Catalog.CreateServiceDoc("
  "r.Client.UpdateServiceDoc(:r.Catalog.UpdateServiceDoc("
  "r.Client.DeleteServiceDoc(:r.Catalog.DeleteServiceDoc("
  "r.Client.ListServiceDiagrams(:r.Catalog.ListServiceDiagrams("
  "r.Client.CreateServiceDiagram(:r.Catalog.CreateServiceDiagram("
  "r.Client.DeleteServiceDiagram(:r.Catalog.DeleteServiceDiagram("
  "r.Client.ListServiceDBs(:r.Catalog.ListServiceDBs("
  "r.Client.GetServiceDB(:r.Catalog.GetServiceDB("
  "r.Client.CreateServiceDB(:r.Catalog.CreateServiceDB("
  "r.Client.UpdateServiceDB(:r.Catalog.UpdateServiceDB("
  "r.Client.DeleteServiceDB(:r.Catalog.DeleteServiceDB("
  "r.Client.ListServiceDBVersions(:r.Catalog.ListServiceDBVersions("
  "r.Client.CreateServiceDBVersion(:r.Catalog.CreateServiceDBVersion("
  "r.Client.RestoreServiceDBVersion(:r.Catalog.RestoreServiceDBVersion("
  "r.Client.ListAPIEndpoints(:r.Catalog.ListAPIEndpoints("
  "r.Client.GetAPIEndpoint(:r.Catalog.GetAPIEndpoint("
  "r.Client.CreateAPIEndpoint(:r.Catalog.CreateAPIEndpoint("
  "r.Client.UpdateAPIEndpoint(:r.Catalog.UpdateAPIEndpoint("
  "r.Client.DeleteAPIEndpoint(:r.Catalog.DeleteAPIEndpoint("
  "r.Client.ListTestPacks(:r.TestPack.ListTestPacks("
  "r.Client.CreateTestPack(:r.TestPack.CreateTestPack("
  "r.Client.UpdateTestPack(:r.TestPack.UpdateTestPack("
  "r.Client.DeleteTestPack(:r.TestPack.DeleteTestPack("
  "r.Client.ListTestCases(:r.TestPack.ListTestCases("
  "r.Client.CreateTestCase(:r.TestPack.CreateTestCase("
  "r.Client.UpdateTestCase(:r.TestPack.UpdateTestCase("
  "r.Client.DeleteTestCase(:r.TestPack.DeleteTestCase("
  "r.Client.GetTestRun(:r.TestPack.GetTestRun("
  "r.Client.ListTestRuns(:r.TestPack.ListTestRuns("
  "r.Client.ListTestRunsSummary(:r.TestPack.ListTestRunsSummary("
  "r.Client.CreateTestRun(:r.TestPack.CreateTestRun("
  "r.Client.UpdateTestRun(:r.TestPack.UpdateTestRun("
  "r.Client.ListTestRunResults(:r.TestPack.ListTestRunResults("
  "r.Client.CreateTestRunResult(:r.TestPack.CreateTestRunResult("
  "r.Client.UpdateTestRunResult(:r.TestPack.UpdateTestRunResult("
)

for f in internal/graph/*.resolvers.go; do
  for pair in "${CLIENT_RENAMES[@]}"; do
    old="${pair%%:*}"
    new="${pair#*:}"
    sed -i '' "s/${old//./\\.}/${new}/g" "$f"
  done
done
```

(On Linux, drop the `''` after `-i`. Longer method names like `ListServiceAccountTokens` are listed before nothing that could prefix-collide — each pattern includes the trailing `(`, and no method name is a prefix of another method name in this list followed immediately by `(`, so single-pass replacement is safe regardless of order.)

- [ ] **Step 4: Update `cmd/server/main.go` wiring**

Change:

```go
	resolver := &graph.Resolver{Client: c}
```

to:

```go
	resolver := &graph.Resolver{
		Auth:      c,
		Org:       c,
		Admin:     c,
		Folder:    c,
		Diagram:   c,
		Component: c,
		UIMap:     c,
		Catalog:   c,
		TestPack:  c,
		Actor:     c,
	}
```

- [ ] **Step 5: Build**

Run: `go build ./...`
Expected: exits 0, no output. A failure here means `*uigraphapi.Client` is missing a method one of the new interfaces declares, or a method signature was transcribed incorrectly — the compiler error names the exact interface and missing/mismatched method.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: narrow Resolver's client dependency into per-domain interfaces"
```

---

## Task 11: Extract `internal/server/server.go` with graceful shutdown and HTTP timeouts

**Files:**
- Create: `internal/server/server.go`
- Modify: `cmd/server/main.go` (shrinks to config load + `server.Run`)

**Interfaces:**
- Consumes: `config.Config` (Task 1), `graph.Resolver` + its 10 interface fields (Task 10), `uigraphapi.New` (Task 2).
- Produces: `server.Run(cfg *config.Config) error` — builds the client, resolver, executable schema, mux, and `*http.Server`; blocks until `SIGINT`/`SIGTERM` triggers a graceful `Shutdown` (10s deadline) or `ListenAndServe` returns a non-`http.ErrServerClosed` error.

- [ ] **Step 1: Create `internal/server/server.go`**

```go
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/graph"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/middleware"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 60 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// Run builds the HTTP server from cfg and blocks until it shuts down, either
// because ListenAndServe failed or the process received SIGINT/SIGTERM.
func Run(cfg *config.Config) error {
	c := uigraphapi.New(cfg.APIBaseURL)

	resolver := &graph.Resolver{
		Auth:      c,
		Org:       c,
		Admin:     c,
		Folder:    c,
		Diagram:   c,
		Component: c,
		UIMap:     c,
		Catalog:   c,
		TestPack:  c,
		Actor:     c,
	}
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	gqlSrv := newGraphQLServer(schema, cfg.Env)

	mux := http.NewServeMux()
	mux.Handle("POST /graphql", middleware.Auth(gqlSrv))
	mux.Handle("GET /graphql", middleware.Auth(gqlSrv))

	if cfg.Env != "prod" {
		mux.Handle("GET /playground", playground.Handler("UIGraph GraphQL", "/graphql"))
		slog.Info("playground enabled", "url", "http://localhost:"+cfg.Port+"/playground")
	}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("uigraph-graphql listening", "port", cfg.Port, "upstream", cfg.APIBaseURL)
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	}
}

func newGraphQLServer(schema graphql.ExecutableSchema, env string) *handler.Server {
	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	if env != "prod" {
		srv.Use(extension.Introspection{})
	}
	return srv
}
```

- [ ] **Step 2: Rewrite `cmd/server/main.go`**

```go
package main

import (
	"log/slog"
	"os"

	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/server"
)

func main() {
	cfg := config.Load()

	if err := server.Run(cfg); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 4: Manual smoke test**

Run: `go run ./cmd/server &` then `curl -s localhost:8090/healthz; kill %1`
Expected: `curl` prints `ok`; the process exits cleanly when killed (no panic, no leaked-goroutine warnings from `-race` if you ran `go run -race`).

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "refactor: extract internal/server with graceful shutdown and HTTP timeouts"
```

---

## Task 12: Add request-ID + structured access-log middleware

**Files:**
- Create: `internal/middleware/logging.go`
- Modify: `internal/server/server.go` (wrap the mux with `middleware.Logging`)
- Modify: `go.mod`, `go.sum` (promote `github.com/google/uuid` from indirect to direct)

**Interfaces:**
- Consumes: nothing new.
- Produces: `middleware.Logging(next http.Handler) http.Handler`, `middleware.RequestID(ctx context.Context) string`.

- [ ] **Step 1: Create `internal/middleware/logging.go`**

```go
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type requestIDKey struct{}

// Logging wraps next with a per-request access log line and a request id
// propagated through the context so downstream code can correlate logs.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(ctx))

		slog.InfoContext(ctx, "http request",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// RequestID returns the request id stored in ctx by Logging, or "" if absent.
func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey{}).(string)
	return v
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
```

- [ ] **Step 2: Wire it into `internal/server/server.go`**

Change:

```go
	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
```

to:

```go
	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           middleware.Logging(mux),
```

- [ ] **Step 3: Promote `google/uuid` to a direct dependency**

Run: `go mod tidy`
Expected: `go.mod` now lists `github.com/google/uuid v1.6.0` under the direct `require` block (no `// indirect` suffix); `go.sum` is unchanged or only reformatted.

- [ ] **Step 4: Build and smoke-test**

Run: `go build ./...`
Expected: exits 0, no output.

Run: `go run ./cmd/server & sleep 1 && curl -s localhost:8090/healthz >/dev/null; kill %1`
Expected: a JSON-ish text log line is printed to stdout containing `"msg":"http request"`, `"path":"/healthz"`, `"status":200`.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add request-id and structured access-log middleware"
```

---

## Task 13: Sanitize GraphQL errors before they reach the client

**Files:**
- Create: `internal/graph/errors.go`
- Modify: `internal/server/server.go` (call `SetErrorPresenter`)

**Interfaces:**
- Consumes: `uigraphapi.APIError` (Task 2).
- Produces: `graph.ErrorPresenter(ctx context.Context, err error) *gqlerror.Error`, matching gqlgen's `graphql.ErrorPresenterFunc` signature.

- [ ] **Step 1: Create `internal/graph/errors.go`**

```go
package graph

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

// ErrorPresenter logs every resolver error server-side, then returns a
// sanitized message to the GraphQL client so upstream REST error bodies
// (which may contain internal details) are never forwarded verbatim.
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	gqlErr := graphql.DefaultErrorPresenter(ctx, err)

	var apiErr *uigraphapi.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Status {
		case http.StatusNotFound:
			gqlErr.Message = "not found"
			return gqlErr
		case http.StatusUnauthorized:
			gqlErr.Message = "unauthorized"
			return gqlErr
		case http.StatusForbidden:
			gqlErr.Message = "forbidden"
			return gqlErr
		case http.StatusBadRequest, http.StatusConflict, http.StatusUnprocessableEntity:
			gqlErr.Message = "invalid request"
			return gqlErr
		}
	}

	slog.ErrorContext(ctx, "graphql resolver error", "err", err, "path", graphql.GetPath(ctx).String())
	gqlErr.Message = "internal error"
	return gqlErr
}
```

- [ ] **Step 2: Wire it into `internal/server/server.go`**

In `newGraphQLServer`, change:

```go
func newGraphQLServer(schema graphql.ExecutableSchema, env string) *handler.Server {
	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	if env != "prod" {
		srv.Use(extension.Introspection{})
	}
	return srv
}
```

to:

```go
func newGraphQLServer(schema graphql.ExecutableSchema, env string) *handler.Server {
	srv := handler.New(schema)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetErrorPresenter(graph.ErrorPresenter)
	if env != "prod" {
		srv.Use(extension.Introspection{})
	}
	return srv
}
```

- [ ] **Step 3: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: sanitize GraphQL error messages, log real errors server-side"
```

---

## Task 14: Add a GraphQL query complexity limit

**Files:**
- Modify: `internal/server/server.go`

**Interfaces:**
- Consumes: `extension.FixedComplexityLimit` (gqlgen, already a dependency).
- Produces: queries exceeding 1000 complexity points are rejected with `COMPLEXITY_LIMIT_EXCEEDED` before execution. Default per-field cost is 1, so this rejects pathologically wide/deep queries while leaving normal UI queries (tens of fields) untouched.

- [ ] **Step 1: Add the limit in `newGraphQLServer`**

Change:

```go
	srv.SetErrorPresenter(graph.ErrorPresenter)
	if env != "prod" {
		srv.Use(extension.Introspection{})
	}
	return srv
```

to:

```go
	srv.SetErrorPresenter(graph.ErrorPresenter)
	srv.Use(extension.FixedComplexityLimit(1000))
	if env != "prod" {
		srv.Use(extension.Introspection{})
	}
	return srv
```

- [ ] **Step 2: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 3: Manual smoke test**

Run: `go run ./cmd/server &` then, from the GraphQL Playground at `http://localhost:8090/playground`, run a normal query like `{ myOrgs { id name } }`.
Expected: succeeds normally — confirms the limit doesn't interfere with ordinary queries. Stop the server with `kill %1`.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: add a fixed GraphQL query complexity limit"
```

---

## Task 15: Validate configuration at boot

**Files:**
- Modify: `internal/config/config.go`, `cmd/server/main.go`

**Interfaces:**
- Consumes: nothing new.
- Produces: `config.Load() (*Config, error)` — was `config.Load() *Config`. Returns an error when `API_BASE_URL` doesn't parse as an absolute URL, so a misconfigured deployment fails at startup instead of on the first GraphQL request.

- [ ] **Step 1: Rewrite `internal/config/config.go`**

```go
package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	APIBaseURL string // uigraph-api base URL, e.g. http://uigraph-api:8080
	Port       string // HTTP listen port for this server
	Env        string // local | dev | prod
}

func Load() (*Config, error) {
	cfg := &Config{
		APIBaseURL: getenv("API_BASE_URL", "http://localhost:8080"),
		Port:       getenv("PORT", "8090"),
		Env:        getenv("ENV", "local"),
	}

	u, err := url.Parse(cfg.APIBaseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("API_BASE_URL %q is not a valid absolute URL", cfg.APIBaseURL)
	}

	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
```

- [ ] **Step 2: Update `cmd/server/main.go`**

```go
package main

import (
	"log/slog"
	"os"

	"github.com/uigraph/graphql/internal/config"
	"github.com/uigraph/graphql/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "err", err)
		os.Exit(1)
	}

	if err := server.Run(cfg); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 4: Manual smoke test**

Run: `API_BASE_URL="not-a-url" go run ./cmd/server`
Expected: prints a structured error log line containing `invalid configuration` and `API_BASE_URL "not-a-url" is not a valid absolute URL`, then exits with status 1 (check with `echo $?`).

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: validate API_BASE_URL at startup instead of failing on first request"
```

---

## Task 16: Add a `/readyz` endpoint that checks upstream reachability

**Files:**
- Modify: `internal/uigraphapi/client.go` (add `Ping`), `internal/server/server.go` (add route + handler)

**Interfaces:**
- Consumes: `Client.get` (Task 2).
- Produces: `(*Client) Ping(ctx context.Context) error` — `GET /healthz` against the configured `uigraph-api` base URL. `GET /readyz` on this server returns 200 if `Ping` succeeds within 3s, 503 otherwise.

- [ ] **Step 1: Add `Ping` to `internal/uigraphapi/client.go`**

Add this method after `New`:

```go
// Ping checks that the configured uigraph-api backend is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.get(ctx, "/healthz", nil)
}
```

- [ ] **Step 2: Add the route and handler to `internal/server/server.go`**

Change:

```go
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
```

to:

```go
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /readyz", readyzHandler(c))
```

Then add this function at the end of the file:

```go
func readyzHandler(c *uigraphapi.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if err := c.Ping(ctx); err != nil {
			slog.WarnContext(ctx, "readiness check failed", "err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
```

- [ ] **Step 3: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 4: Manual smoke test**

Run: `go run ./cmd/server &` then `curl -s -o /dev/null -w "%{http_code}\n" localhost:8090/readyz; kill %1`
Expected: prints `200` if a reachable `uigraph-api` is running at the configured `API_BASE_URL` (default `http://localhost:8080`), or `503` if not — either is correct behavior to observe, the point is it's not always `200` regardless of upstream state.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add /readyz endpoint checking uigraph-api reachability"
```

---

## Task 17: Optional, env-gated CORS middleware

**Files:**
- Create: `internal/middleware/cors.go`
- Modify: `internal/config/config.go` (add `AllowedOrigins`), `internal/server/server.go` (wire it)

**Interfaces:**
- Consumes: nothing new.
- Produces: `middleware.CORS(allowedOrigins []string, next http.Handler) http.Handler` — returns `next` unmodified when `allowedOrigins` is empty (the default; production is same-origin behind Caddy per `uigraph-deploy`).

- [ ] **Step 1: Create `internal/middleware/cors.go`**

```go
package middleware

import "net/http"

// CORS sets Access-Control-* headers for the given allowed origins. If
// allowedOrigins is empty, it returns next unmodified — CORS is opt-in
// because production traffic is same-origin behind a reverse proxy.
func CORS(allowedOrigins []string, next http.Handler) http.Handler {
	if len(allowedOrigins) == 0 {
		return next
	}
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

- [ ] **Step 2: Add `AllowedOrigins` to `internal/config/config.go`**

Change the `Config` struct to:

```go
type Config struct {
	APIBaseURL     string   // uigraph-api base URL, e.g. http://uigraph-api:8080
	Port           string   // HTTP listen port for this server
	Env            string   // local | dev | prod
	AllowedOrigins []string // CORS allow-list; empty disables CORS handling entirely
}
```

Add `"strings"` to the import block, and change the body of `Load` to:

```go
	cfg := &Config{
		APIBaseURL: getenv("API_BASE_URL", "http://localhost:8080"),
		Port:       getenv("PORT", "8090"),
		Env:        getenv("ENV", "local"),
	}
	if v := getenv("ALLOWED_ORIGINS", ""); v != "" {
		cfg.AllowedOrigins = strings.Split(v, ",")
	}
```

(keep the existing `url.Parse` validation block right after this, unchanged.)

- [ ] **Step 3: Wire it into `internal/server/server.go`**

Change:

```go
	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           middleware.Logging(mux),
```

to:

```go
	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           middleware.Logging(middleware.CORS(cfg.AllowedOrigins, mux)),
```

- [ ] **Step 4: Build**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add optional env-gated CORS middleware"
```

---

## Task 18: Run the container as a non-root user

**Files:**
- Modify: `Dockerfile`

**Interfaces:**
- Consumes: nothing.
- Produces: the production image runs as an unprivileged `uigraph` user instead of root.

- [ ] **Step 1: Update the final stage of `Dockerfile`**

Change:

```dockerfile
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /uigraph-graphql /usr/local/bin/uigraph-graphql
EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/uigraph-graphql"]
```

to:

```dockerfile
FROM alpine:3.20
RUN apk add --no-cache ca-certificates && \
    addgroup -S uigraph && adduser -S uigraph -G uigraph
COPY --from=builder /uigraph-graphql /usr/local/bin/uigraph-graphql
USER uigraph
EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/uigraph-graphql"]
```

- [ ] **Step 2: Verify the image builds and runs**

Run: `docker build -t uigraph-graphql:test .`
Expected: build succeeds.

Run: `docker run --rm -p 8090:8090 uigraph-graphql:test &` then `curl -s localhost:8090/healthz; docker stop $(docker ps -q --filter ancestor=uigraph-graphql:test)`
Expected: `curl` prints `ok`.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "fix: run container as non-root user"
```

---

## Task 19: Unit tests for `internal/graph/convert` — auth, org, admin, folder

**Files:**
- Create: `internal/graph/convert/auth_test.go`, `org_test.go`, `admin_test.go`, `folder_test.go`

**Interfaces:**
- Consumes: `convert.MeToModel`, `convert.OrgSummariesToModel`, `convert.TeamToModel`, `convert.InvitationToModel`, `convert.UserToModel`, `convert.FolderToModel` (Task 9).
- Produces: nothing new — pure test coverage.

These tests focus on functions with actual branching (optional-pointer defaulting, nil handling) rather than every trivial 1:1 field-copy mapper — a pure struct literal copy has no behavior for a test to catch beyond what the compiler already guarantees.

- [ ] **Step 1: Create `internal/graph/convert/auth_test.go`**

```go
package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestMeToModel(t *testing.T) {
	withAvatar := MeToModel(&uigraphapi.MeResponse{
		UserID: "u1", OrgID: "o1", Email: "a@b.com", Name: "Ann", Login: "ann",
		Kind: "user", Role: "admin", AuthProvider: "local", AvatarURL: "https://x/a.png",
	})
	if withAvatar.AvatarURL == nil || *withAvatar.AvatarURL != "https://x/a.png" {
		t.Fatalf("AvatarURL = %v, want pointer to https://x/a.png", withAvatar.AvatarURL)
	}
	if withAvatar.UserID != "u1" || withAvatar.Role != "admin" {
		t.Fatalf("unexpected fields: %+v", withAvatar)
	}

	withoutAvatar := MeToModel(&uigraphapi.MeResponse{UserID: "u2"})
	if withoutAvatar.AvatarURL != nil {
		t.Fatalf("AvatarURL = %v, want nil for empty AvatarURL", *withoutAvatar.AvatarURL)
	}
}

func TestOrgSummariesToModel(t *testing.T) {
	in := []uigraphapi.OrgSummary{
		{ID: "1", Name: "A", Active: true},
		{ID: "2", Name: "B", Active: false},
	}
	out := OrgSummariesToModel(in)
	if len(out) != 2 {
		t.Fatalf("len = %d, want 2", len(out))
	}
	if out[0].ID != "1" || out[1].Active != false {
		t.Fatalf("unexpected output: %+v", out)
	}
}
```

- [ ] **Step 2: Create `internal/graph/convert/org_test.go`**

```go
package convert

import (
	"testing"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestTeamToModel(t *testing.T) {
	withExtras := TeamToModel(&uigraphapi.Team{
		ID: "t1", OrgID: "o1", Name: "Platform", Email: "platform@x.com", ExternalID: "ext-1",
		CreatedAt: time.Unix(0, 0), UpdatedAt: time.Unix(0, 0),
	})
	if withExtras.Email == nil || *withExtras.Email != "platform@x.com" {
		t.Fatalf("Email = %v, want pointer to platform@x.com", withExtras.Email)
	}
	if withExtras.ExternalID == nil || *withExtras.ExternalID != "ext-1" {
		t.Fatalf("ExternalID = %v, want pointer to ext-1", withExtras.ExternalID)
	}

	bare := TeamToModel(&uigraphapi.Team{ID: "t2", OrgID: "o1", Name: "Bare"})
	if bare.Email != nil || bare.ExternalID != nil {
		t.Fatalf("expected nil Email/ExternalID for empty input, got %+v", bare)
	}
}

func TestInvitationToModel(t *testing.T) {
	expires := time.Now().Add(24 * time.Hour)
	withExpiry := InvitationToModel(uigraphapi.Invitation{ID: "i1", ExpiresAt: &expires})
	if withExpiry.ExpiresAt == nil || !withExpiry.ExpiresAt.Equal(expires) {
		t.Fatalf("ExpiresAt = %v, want %v", withExpiry.ExpiresAt, expires)
	}

	noExpiry := InvitationToModel(uigraphapi.Invitation{ID: "i2"})
	if noExpiry.ExpiresAt != nil {
		t.Fatalf("ExpiresAt = %v, want nil", noExpiry.ExpiresAt)
	}
}
```

- [ ] **Step 3: Create `internal/graph/convert/admin_test.go`**

```go
package convert

import (
	"testing"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestUserToModel(t *testing.T) {
	lastSeen := time.Now()
	active := UserToModel(&uigraphapi.User{ID: "u1", Email: "a@b.com", LastSeenAt: &lastSeen})
	if active.LastSeenAt == nil || !active.LastSeenAt.Equal(lastSeen) {
		t.Fatalf("LastSeenAt = %v, want %v", active.LastSeenAt, lastSeen)
	}

	neverSeen := UserToModel(&uigraphapi.User{ID: "u2"})
	if neverSeen.LastSeenAt != nil {
		t.Fatalf("LastSeenAt = %v, want nil", neverSeen.LastSeenAt)
	}
}
```

- [ ] **Step 4: Create `internal/graph/convert/folder_test.go`**

```go
package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestFolderToModel(t *testing.T) {
	parentID := "parent-1"
	in := &uigraphapi.Folder{
		ID: "f1", OrgID: "o1", ParentID: &parentID, Type: "diagrams",
		Name: "My Folder", Order: 1.5, CreatedBy: "u1",
	}
	out := FolderToModel(in)
	if out.ID != "f1" || out.Name != "My Folder" || out.Order != 1.5 {
		t.Fatalf("unexpected fields: %+v", out)
	}
	if out.ParentID == nil || *out.ParentID != parentID {
		t.Fatalf("ParentID = %v, want pointer to %q", out.ParentID, parentID)
	}
	if out.TeamID != nil {
		t.Fatalf("TeamID = %v, want nil", out.TeamID)
	}
}
```

- [ ] **Step 5: Run the tests**

Run: `go test ./internal/graph/convert/...`
Expected: `ok  github.com/uigraph/graphql/internal/graph/convert`, all tests pass.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "test: add convert package tests for auth/org/admin/folder"
```

---

## Task 20: Unit tests for `internal/graph/convert` — diagram, uimap, catalog, testpack

**Files:**
- Create: `internal/graph/convert/diagram_test.go`, `uimap_test.go`, `catalog_test.go`, `testpack_test.go`

**Interfaces:**
- Consumes: `convert.DiagramVersionToModel`, `convert.DiagramToModel`, `convert.CanvasToModel`, `convert.FocalPointMetaToModel`, `convert.FocalPointMetaBody`, `convert.ServiceToModel`, `convert.ServiceDiagramToModel`, `convert.AuthConfigToModel`, `convert.ManualTestCaseToModel`, `convert.TestRunResultToModel` (Task 9).
- Produces: nothing new — pure test coverage.

- [ ] **Step 1: Create `internal/graph/convert/diagram_test.go`**

```go
package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestDiagramVersionToModel(t *testing.T) {
	out := DiagramVersionToModel("org-1", uigraphapi.DiagramVersion{ID: "v1", DiagramID: "d1", VersionNumber: 2})
	if out.OrgID != "org-1" {
		t.Fatalf("OrgID = %q, want %q (the orgID parameter, since DiagramVersion carries no OrgID field itself)", out.OrgID, "org-1")
	}
	if out.ID != "v1" || out.VersionNumber != 2 {
		t.Fatalf("unexpected fields: %+v", out)
	}
}

func TestDiagramToModel(t *testing.T) {
	previewID := "asset-1"
	out := DiagramToModel(&uigraphapi.Diagram{ID: "d1", OrgID: "o1", Name: "Checkout", PreviewAssetID: &previewID})
	if out.PreviewAssetID == nil || *out.PreviewAssetID != previewID {
		t.Fatalf("PreviewAssetID = %v, want pointer to %q", out.PreviewAssetID, previewID)
	}
	if out.UpdatedBy != nil {
		t.Fatalf("UpdatedBy = %v, want nil", out.UpdatedBy)
	}
}
```

- [ ] **Step 2: Create `internal/graph/convert/uimap_test.go`**

```go
package convert

import (
	"encoding/json"
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestCanvasToModel(t *testing.T) {
	empty := CanvasToModel(&uigraphapi.Canvas{MapID: "m1", OrgID: "o1"})
	if empty.FramePositions != "{}" {
		t.Fatalf("FramePositions = %q, want %q for empty RawMessage", empty.FramePositions, "{}")
	}

	populated := CanvasToModel(&uigraphapi.Canvas{
		MapID: "m1", OrgID: "o1",
		FramePositions: json.RawMessage(`{"frame-1":{"x":10,"y":20}}`),
	})
	if populated.FramePositions != `{"frame-1":{"x":10,"y":20}}` {
		t.Fatalf("FramePositions = %q, want passthrough of input JSON", populated.FramePositions)
	}
}

func TestFocalPointMetaToModel(t *testing.T) {
	out := FocalPointMetaToModel(&uigraphapi.FocalPointMeta{
		ID: "fpm1", ComponentImages: json.RawMessage(`["a.png"]`),
	})
	if out.ComponentImages != `["a.png"]` {
		t.Fatalf("ComponentImages = %q, want passthrough", out.ComponentImages)
	}

	empty := FocalPointMetaToModel(&uigraphapi.FocalPointMeta{ID: "fpm2"})
	if empty.ComponentImages != "[]" {
		t.Fatalf("ComponentImages = %q, want %q for empty RawMessage", empty.ComponentImages, "[]")
	}
}

func TestFocalPointMetaBody(t *testing.T) {
	body := map[string]interface{}{
		"componentImages":      `["a.png","b.png"]`,
		"componentModalFields": `{"foo":"bar"}`,
		"componentId":          "c1",
	}
	out := FocalPointMetaBody(body)

	images, ok := out["componentImages"].([]interface{})
	if !ok || len(images) != 2 {
		t.Fatalf("componentImages = %#v, want a 2-element slice decoded from JSON", out["componentImages"])
	}
	if out["componentId"] != "c1" {
		t.Fatalf("componentId = %v, want unchanged passthrough", out["componentId"])
	}
}
```

- [ ] **Step 3: Create `internal/graph/convert/catalog_test.go`**

```go
package convert

import (
	"encoding/json"
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestServiceToModel(t *testing.T) {
	withMetadata := ServiceToModel(&uigraphapi.Service{ID: "s1", Metadata: json.RawMessage(`{"team":"core"}`)})
	if withMetadata.Metadata != `{"team":"core"}` {
		t.Fatalf("Metadata = %q, want passthrough", withMetadata.Metadata)
	}

	noMetadata := ServiceToModel(&uigraphapi.Service{ID: "s2"})
	if noMetadata.Metadata != "{}" {
		t.Fatalf("Metadata = %q, want %q for empty RawMessage", noMetadata.Metadata, "{}")
	}
}

func TestServiceDiagramToModel(t *testing.T) {
	withDiagram := ServiceDiagramToModel(&uigraphapi.ServiceDiagram{
		ServiceID: "svc1", DiagramID: "d1",
		Diagram: &uigraphapi.Diagram{ID: "d1", Name: "Checkout"},
	})
	if withDiagram.Diagram == nil || withDiagram.Diagram.Name != "Checkout" {
		t.Fatalf("Diagram = %v, want a converted Diagram named Checkout", withDiagram.Diagram)
	}

	withoutDiagram := ServiceDiagramToModel(&uigraphapi.ServiceDiagram{ServiceID: "svc1", DiagramID: "d2"})
	if withoutDiagram.Diagram != nil {
		t.Fatalf("Diagram = %v, want nil when source Diagram pointer is nil", withoutDiagram.Diagram)
	}
}
```

- [ ] **Step 4: Create `internal/graph/convert/testpack_test.go`**

```go
package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestAuthConfigToModel(t *testing.T) {
	if got := AuthConfigToModel(nil); got != nil {
		t.Fatalf("AuthConfigToModel(nil) = %v, want nil", got)
	}

	token := "secret"
	got := AuthConfigToModel(&uigraphapi.AuthConfig{Type: "bearer", BearerToken: &token})
	if got.Type != "bearer" || got.BearerToken == nil || *got.BearerToken != token {
		t.Fatalf("unexpected output: %+v", got)
	}
}

func TestManualTestCaseToModel(t *testing.T) {
	if got := ManualTestCaseToModel(nil); got != nil {
		t.Fatalf("ManualTestCaseToModel(nil) = %v, want nil", got)
	}

	in := &uigraphapi.ManualTestCase{
		Steps: []uigraphapi.TestCaseStep{
			{Order: 1, Action: "open page", ExpectedResult: "page loads"},
		},
	}
	got := ManualTestCaseToModel(in)
	if len(got.Steps) != 1 || got.Steps[0].Action != "open page" {
		t.Fatalf("unexpected steps: %+v", got.Steps)
	}
}

func TestTestRunResultToModel(t *testing.T) {
	var ms int64 = 1500
	got := TestRunResultToModel(&uigraphapi.TestRunResult{TestRunResultID: "r1", ResponseTimeMs: &ms})
	if got.ResponseTimeMs == nil || *got.ResponseTimeMs != 1500 {
		t.Fatalf("ResponseTimeMs = %v, want pointer to 1500 (converted from *int64 to *int)", got.ResponseTimeMs)
	}

	noTiming := TestRunResultToModel(&uigraphapi.TestRunResult{TestRunResultID: "r2"})
	if noTiming.ResponseTimeMs != nil {
		t.Fatalf("ResponseTimeMs = %v, want nil", noTiming.ResponseTimeMs)
	}
}
```

- [ ] **Step 5: Run the tests**

Run: `go test ./internal/graph/convert/...`
Expected: `ok  github.com/uigraph/graphql/internal/graph/convert`, all tests pass.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "test: add convert package tests for diagram/uimap/catalog/testpack"
```

---

## Task 21: `httptest`-based tests for `internal/uigraphapi`

**Files:**
- Create: `internal/uigraphapi/client_test.go`

**Interfaces:**
- Consumes: `Client.New`, `Client.GetOrg`, `Client.CreateOrg`, `APIError`, `IsNotFound` (Task 2).
- Produces: nothing new — pure test coverage of the shared `do`/`get`/`post` machinery in `client.go`, exercised through two representative domain methods.

- [ ] **Step 1: Create `internal/uigraphapi/client_test.go`**

```go
package uigraphapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet_DecodesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/orgs/org-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Org{ID: "org-1", Name: "Acme"})
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.GetOrg(context.Background(), "org-1")
	if err != nil {
		t.Fatalf("GetOrg() error = %v", err)
	}
	if got.ID != "org-1" || got.Name != "Acme" {
		t.Fatalf("GetOrg() = %+v, want ID=org-1 Name=Acme", got)
	}
}

func TestClientGet_404ReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.GetOrg(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected an error for a 404 response, got nil")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(err) = false, want true for err = %v", err)
	}
}

func TestClientGet_500IsNotIsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.GetOrg(context.Background(), "x")
	if err == nil {
		t.Fatal("expected an error for a 500 response, got nil")
	}
	if IsNotFound(err) {
		t.Fatal("IsNotFound(err) = true, want false for a 500 response")
	}
}

func TestClientPost_SendsBody(t *testing.T) {
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/orgs" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Org{ID: "org-2", Name: gotBody["name"].(string)})
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.CreateOrg(context.Background(), map[string]interface{}{"name": "Globex"})
	if err != nil {
		t.Fatalf("CreateOrg() error = %v", err)
	}
	if gotBody["name"] != "Globex" {
		t.Fatalf("server received body %v, want name=Globex", gotBody)
	}
	if got.Name != "Globex" {
		t.Fatalf("CreateOrg() = %+v, want Name=Globex", got)
	}
}
```

- [ ] **Step 2: Run the tests**

Run: `go test ./internal/uigraphapi/...`
Expected: `ok  github.com/uigraph/graphql/internal/uigraphapi`, all tests pass.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "test: add httptest-based tests for internal/uigraphapi client"
```

---

## Task 22: Tests for `internal/middleware`

**Files:**
- Create: `internal/middleware/auth_test.go`, `internal/middleware/logging_test.go`

**Interfaces:**
- Consumes: `middleware.Auth`, `middleware.ApplyAuth`, `middleware.BearerToken` (Task 1), `middleware.Logging`, `middleware.RequestID` (Task 12).
- Produces: nothing new — pure test coverage.

- [ ] **Step 1: Create `internal/middleware/auth_test.go`**

```go
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth_PropagatesHeadersToOutgoingRequest(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	inbound.Header.Set("Authorization", "Bearer abc123")
	inbound.Header.Set("Cookie", "session=xyz")

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	outbound, _ := http.NewRequest(http.MethodGet, "http://upstream/api/v1/auth/me", nil)
	ApplyAuth(capturedCtx, outbound)

	if got := outbound.Header.Get("Authorization"); got != "Bearer abc123" {
		t.Fatalf("Authorization = %q, want %q", got, "Bearer abc123")
	}
	if got := outbound.Header.Get("Cookie"); got != "session=xyz" {
		t.Fatalf("Cookie = %q, want %q", got, "session=xyz")
	}
}

func TestAuth_NoHeadersMeansNothingPropagated(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	outbound, _ := http.NewRequest(http.MethodGet, "http://upstream/api/v1/auth/me", nil)
	ApplyAuth(capturedCtx, outbound)

	if got := outbound.Header.Get("Authorization"); got != "" {
		t.Fatalf("Authorization = %q, want empty", got)
	}
}

func TestBearerToken(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	inbound.Header.Set("Authorization", "Bearer abc123")

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	if got := BearerToken(capturedCtx); got != "abc123" {
		t.Fatalf("BearerToken() = %q, want %q", got, "abc123")
	}
}

func TestBearerToken_NonBearerAuthReturnsEmpty(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	inbound.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	if got := BearerToken(capturedCtx); got != "" {
		t.Fatalf("BearerToken() = %q, want empty for non-Bearer Authorization header", got)
	}
}
```

- [ ] **Step 2: Create `internal/middleware/logging_test.go`**

```go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging_SetsRequestIDInContext(t *testing.T) {
	var gotID string
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = RequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if gotID == "" {
		t.Fatal("RequestID(ctx) is empty, want a generated request id")
	}
}

func TestLogging_RecordsResponseStatus(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusTeapot {
		t.Fatalf("recorder status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}
```

- [ ] **Step 3: Run the tests**

Run: `go test ./internal/middleware/...`
Expected: `ok  github.com/uigraph/graphql/internal/middleware`, all tests pass.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "test: add tests for internal/middleware auth and logging"
```

---

## Task 23: Resolver-level integration tests through the real executable schema

**Files:**
- Create: `internal/graph/resolver_test.go`

**Interfaces:**
- Consumes: `graph.Resolver`, the `authClient`/`folderClient` interfaces (Task 10), `generated.NewExecutableSchema` (gqlgen-generated).
- Produces: nothing new — proves the GraphQL wiring works end-to-end (HTTP → gqlgen → resolver → fake client → response) without a real `uigraph-api`.

- [ ] **Step 1: Create `internal/graph/resolver_test.go`**

```go
package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/uigraph/graphql/internal/graph"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

type fakeAuthClient struct {
	me *uigraphapi.MeResponse
}

func (f *fakeAuthClient) Me(ctx context.Context) (*uigraphapi.MeResponse, error) { return f.me, nil }
func (f *fakeAuthClient) MyOrgs(ctx context.Context) ([]uigraphapi.OrgSummary, error) {
	return nil, nil
}
func (f *fakeAuthClient) SwitchOrg(ctx context.Context, orgID string) error { return nil }

type fakeFolderClient struct {
	created *uigraphapi.Folder
}

func (f *fakeFolderClient) ListFolders(ctx context.Context, orgID, folderType, parentID string) ([]uigraphapi.Folder, error) {
	return nil, nil
}
func (f *fakeFolderClient) GetFolder(ctx context.Context, orgID, id string) (*uigraphapi.Folder, error) {
	return nil, nil
}
func (f *fakeFolderClient) CreateFolder(ctx context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Folder, error) {
	f.created = &uigraphapi.Folder{ID: "folder-1", OrgID: orgID, Name: body["name"].(string), Type: body["type"].(string)}
	return f.created, nil
}
func (f *fakeFolderClient) UpdateFolder(ctx context.Context, orgID, id string, body map[string]interface{}) (*uigraphapi.Folder, error) {
	return nil, nil
}
func (f *fakeFolderClient) DeleteFolder(ctx context.Context, orgID, id string) error { return nil }

func newTestServer(resolver *graph.Resolver) *httptest.Server {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	srv := handler.New(schema)
	srv.AddTransport(transport.POST{})
	return httptest.NewServer(srv)
}

func doGraphQL(t *testing.T, srv *httptest.Server, query string) map[string]interface{} {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /graphql: %v", err)
	}
	defer resp.Body.Close()

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if errs, ok := out["errors"]; ok {
		t.Fatalf("graphql errors: %v", errs)
	}
	return out["data"].(map[string]interface{})
}

func TestMeQuery(t *testing.T) {
	resolver := &graph.Resolver{
		Auth: &fakeAuthClient{me: &uigraphapi.MeResponse{UserID: "u1", Email: "a@b.com", Name: "Ann", Role: "admin"}},
	}
	srv := newTestServer(resolver)
	defer srv.Close()

	data := doGraphQL(t, srv, `{ me { userId email name role } }`)
	me := data["me"].(map[string]interface{})
	if me["userId"] != "u1" || me["email"] != "a@b.com" {
		t.Fatalf("unexpected me: %+v", me)
	}
}

func TestCreateFolderMutation(t *testing.T) {
	folders := &fakeFolderClient{}
	resolver := &graph.Resolver{Folder: folders}
	srv := newTestServer(resolver)
	defer srv.Close()

	data := doGraphQL(t, srv, `mutation { createFolder(orgId: "org-1", input: { name: "Diagrams", type: "diagram" }) { id name type } }`)
	created := data["createFolder"].(map[string]interface{})
	if created["name"] != "Diagrams" || created["type"] != "diagram" {
		t.Fatalf("unexpected createFolder result: %+v", created)
	}
	if folders.created == nil || folders.created.OrgID != "org-1" {
		t.Fatalf("fake folder client did not record the call correctly: %+v", folders.created)
	}
}
```

This file is `package graph_test` (external test package, same directory) so it exercises `graph.Resolver` exactly as `cmd/server/main.go` does, including the unexported `authClient`/`folderClient` interfaces — Go satisfies those structurally since every method they require is exported, so `*fakeAuthClient`/`*fakeFolderClient` from another package can implement them.

- [ ] **Step 2: Run the tests**

Run: `go test ./internal/graph/...`
Expected: `ok  github.com/uigraph/graphql/internal/graph` and `ok  .../internal/graph/convert`, all tests pass.

- [ ] **Step 3: Run the full suite**

Run: `go test ./... -race -cover`
Expected: every package reports `ok`; no `FAIL`.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "test: add resolver-level integration tests through the real executable schema"
```

---

## Task 24: Add `README.md`

**Files:**
- Create: `README.md`

**Interfaces:**
- Consumes: nothing.
- Produces: nothing — documentation only.

- [ ] **Step 1: Create `README.md`**

```markdown
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
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "docs: add README"
```

---

## Task 25: Add `CONTRIBUTING.md`

**Files:**
- Create: `CONTRIBUTING.md`

**Interfaces:**
- Consumes: nothing.
- Produces: nothing — documentation only.

- [ ] **Step 1: Create `CONTRIBUTING.md`**

```markdown
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
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "docs: add CONTRIBUTING"
```

---

## Task 26: Add `CLAUDE.md`

**Files:**
- Create: `CLAUDE.md`

**Interfaces:**
- Consumes: nothing.
- Produces: nothing — documentation only.

- [ ] **Step 1: Create `CLAUDE.md`**

```markdown
# uigraph-graphql — Claude Code Guidelines

Go module: `github.com/uigraph/graphql` · Go 1.25 · gqlgen v0.17 · open-source project

---

## Project layout

```
cmd/server/        entry point
internal/
  config/           env config + validation
  middleware/       auth passthrough, request logging, CORS
  uigraphapi/       typed REST client for uigraph-api, one file per domain
  server/           HTTP server wiring, graceful shutdown
  graph/
    schema/         GraphQL SDL, one file per domain
    generated/, model/  gqlgen output — never hand-edit
    convert/        pure REST-DTO -> GraphQL-model mapping, unit-tested
    *.resolvers.go  resolver implementations, one file per schema file
```

---

## The schema-driven file-split rule

gqlgen's `follow-schema` resolver layout ties each `<name>.resolvers.go` file to its source `<name>.graphqls` file. **Never** try to manually split a `*.resolvers.go` file — the next `go generate` undoes it. To split resolvers, split the schema file instead, then run `go generate ./internal/graph/...`.

---

## Narrow interface pattern

`internal/graph/resolver.go` declares one interface per domain (`authClient`, `orgClient`, `folderClient`, …), each listing only the `*uigraphapi.Client` methods that domain's resolvers actually call. `Resolver` holds these interfaces, not the concrete client — this is what lets tests inject a 2-3 method fake instead of mocking the whole REST client. When adding a domain, add its interface here and wire it in `internal/server/server.go`.

---

## Conversion functions live in `internal/graph/convert/`

Every `*ToModel` function maps an `internal/uigraphapi` DTO onto an `internal/graph/model` GraphQL model. These are pure — no I/O, no context — specifically so they can be unit-tested without a running server. Add new ones there, not inline in resolver methods.

---

## Comments

Default: **no comments.** Well-named identifiers are self-documenting. Add a comment only when the **why** is non-obvious — a hidden constraint, a subtle invariant, a workaround for a specific bug.

---

## Forbidden patterns

- `Resolver.Client *uigraphapi.Client` as a single fat field — use the narrow per-domain interfaces.
- Hand-editing `internal/graph/generated/generated.go` or `internal/graph/model/models_gen.go` — these are gqlgen output.
- Manually splitting a `*.resolvers.go` file without first splitting its source `.graphqls` file.
- Forwarding a raw upstream `uigraphapi.APIError` message straight to the GraphQL client — go through `graph.ErrorPresenter` (`internal/graph/errors.go`), which classifies and sanitizes it.
- Skipping `go build ./...` before committing.
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "docs: add CLAUDE.md contributor/AI guidelines"
```

---

## Task 27: Add `.golangci.yml` and `Makefile`

**Files:**
- Create: `.golangci.yml`, `Makefile`

**Interfaces:**
- Consumes: nothing.
- Produces: `make run`, `make build`, `make test`, `make lint`, `make generate`, `make docker-build`; a `golangci-lint run` config.

- [ ] **Step 1: Create `.golangci.yml`**

```yaml
run:
  timeout: 5m

linters:
  disable-all: true
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - ineffassign
    - gofmt
    - goimports

issues:
  exclude-dirs:
    - internal/graph/generated
    - internal/graph/model
```

- [ ] **Step 2: Create `Makefile`**

```makefile
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
```

- [ ] **Step 3: Verify**

Run: `make build`
Expected: produces `bin/uigraph-graphql`, exits 0.

Run: `make test`
Expected: exits 0, all packages `ok`.

If `golangci-lint` is installed locally, run: `make lint`
Expected: exits 0, no issues reported. (If `golangci-lint` isn't installed locally, this is fine to skip — CI will run it in Task 28.)

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "build: add .golangci.yml and Makefile"
```

---

## Task 28: Add CI workflow

**Files:**
- Create: `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: `make build`/`test`/`generate` targets (Task 27).
- Produces: a GitHub Actions workflow running on every push to `main` and every PR: build, lint, test, and a generated-code drift check.

- [ ] **Step 1: Create `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - run: go build ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - run: go test ./... -race -cover

  generate-drift:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - run: go generate ./internal/graph/...
      - run: git diff --exit-code
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "ci: add GitHub Actions workflow for build/lint/test/generate-drift"
```

---

## Task 29: Final verification against the design spec's success criteria

**Files:** none — verification only.

**Interfaces:** none.

- [ ] **Step 1: Run the full local pipeline**

```bash
go build ./...
go vet ./...
go test ./... -race -cover
go generate ./internal/graph/... && git diff --exit-code
```

Expected: every command exits 0 with no diff from the drift check.

- [ ] **Step 2: Walk the spec's success criteria**

Open `docs/superpowers/specs/2026-06-18-production-refactor-design.md` and confirm each box:

- `cmd/server/main.go` is the only file outside `internal/` — check with `find . -maxdepth 1 -name "*.go"` (expect no output).
- No resolver, convert, or client file exceeds ~300 lines — check with `wc -l internal/graph/*.go internal/graph/convert/*.go internal/uigraphapi/*.go | sort -n | tail -5`.
- `Resolver` holds narrow per-domain interfaces — confirmed by Task 10's `resolver.go`.
- `content.graphqls`/the TestPack portion of `catalog.graphqls` are split into focused schema files, regeneration-stable — confirmed by Task 8.
- Every GraphQL-facing error path logs via `slog` before sanitizing — confirmed by Task 13's `ErrorPresenter`.
- Graceful shutdown, server timeouts, `/readyz` are in place — confirmed by Tasks 11/16.
- `internal/graph/convert/` and `internal/uigraphapi/` have test coverage; CI runs build/lint/test/drift-check — confirmed by Tasks 19–22, 28.
- README, CONTRIBUTING, CLAUDE.md, Makefile, `.golangci.yml`, CI workflow exist — confirmed by Tasks 24–28.

- [ ] **Step 3: Manual end-to-end smoke test**

Run: `go run ./cmd/server &`, then from another terminal open `http://localhost:8090/playground` and run a real query against a running `uigraph-api` instance, e.g. `{ myOrgs { id name } }` (with a valid session cookie/Authorization header attached via the Playground's HTTP headers panel).
Expected: the response is identical in shape to what the pre-refactor server returned for the same query — confirming the "no schema/contract changes" constraint held throughout.

Stop the server: `kill %1`

No commit for this task — it's a verification gate, not a code change. If any step fails, fix the regression and re-run from Step 1 before considering the refactor complete.

