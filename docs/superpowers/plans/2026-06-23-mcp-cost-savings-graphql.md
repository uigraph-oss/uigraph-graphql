# MCP Cost-Savings GraphQL Layer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expose the new uigraph-api MCP savings aggregation endpoints (summary/timeseries/by-tool/by-model/by-user — see the `uigraph-api` repo's `2026-06-23-mcp-savings-aggregation-api.md` plan) through this GraphQL gateway, so uigraph-ui can query them via Apollo.

**Architecture:** This gateway is a thin REST-proxy with zero direct database access — every resolver delegates to a typed HTTP client in `internal/uigraphapi/`. This plan follows the exact existing `testRunsSummary` pattern: one new `.graphqls` schema file → gqlgen-generated resolver stubs → REST client methods on `*uigraphapi.Client` → a `convert` package mapping REST DTOs to gqlgen models → wiring into `Resolver`/`server.go`. The one new wrinkle is `UserSavings.displayName`: the REST API only returns raw `userId`/`serviceAccountId`, so the `CostSavingsByUser` resolver does one batched call to the existing `Actor.ResolveActors` client method (already used elsewhere for `createdByActor`/`updatedByActor` fields) to attach human-readable names, instead of resolving per-row.

**Tech Stack:** Go, gqlgen v0.17.73, stdlib `net/http`/`net/url`. No testify — plain `testing` with table-driven subtests (`t.Run`), matching `internal/graph/convert/testpack_test.go` and `internal/uigraphapi/client_test.go`.

## Global Constraints

- Module path: `github.com/uigraph/graphql`
- Regenerate gqlgen code with `go generate ./internal/graph/...` (or `make generate`) from the repo root after any schema change — this rewrites `internal/graph/generated/generated.go` and `internal/graph/model/models_gen.go`, and creates/updates `internal/graph/mcpsavings.resolvers.go` (resolver stub file named after the schema file, since `gqlgen.yml` uses `filename_template: "{name}.resolvers.go"` with `layout: follow-schema`)
- `period`/`modelId` GraphQL args are always optional (`String`, no `!`) — omitted means "use the backend's default" (`period` defaults to `7d` server-side) or "blended across all models" (`modelId` omitted)
- REST query param names are snake_case (`model_id`), GraphQL arg names are camelCase (`modelId`) — the REST client methods are responsible for this translation, exactly like `ListTestRunsSummary`
- Run from the repo root: `/Users/kranthi/workspace/go/uigraph/backend/uigraph-oss/uigraph-graphql`
- Resolver-level tests never hit a real Postgres or uigraph-api — they spin up the real gqlgen executable schema via `httptest.NewServer` with a fake client injected into `Resolver`, exactly like `internal/graph/resolver_test.go`'s `TestMeQuery`/`TestCreateFolderMutation`

---

### Task 1: REST client layer (`internal/uigraphapi`)

**Files:**
- Create: `internal/uigraphapi/mcpsavings.go`
- Create: `internal/uigraphapi/mcpsavings_test.go`

**Interfaces:**
- Produces: DTOs `SavingsSummary`, `DailySavings`, `ToolSavings`, `ModelSavings`, `UserSavings` (exact field names/JSON tags below), and methods `(c *Client) GetSavingsSummary(ctx, orgID string, period, modelID *string) (*SavingsSummary, error)`, `GetSavingsTimeseries(ctx, orgID string, period, modelID *string) ([]DailySavings, error)`, `GetSavingsByTool(ctx, orgID string, period, modelID *string) ([]ToolSavings, error)`, `GetSavingsByModel(ctx, orgID string, period *string) ([]ModelSavings, error)`, `GetSavingsByUser(ctx, orgID string, period, modelID *string) ([]UserSavings, error)` — the DTOs are consumed by Task 3 (convert), the client methods by Task 4 (resolvers)

- [ ] **Step 1: Write the failing tests**

Create `internal/uigraphapi/mcpsavings_test.go`:

```go
package uigraphapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSavingsSummary_BuildsQueryAndDecodes(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path + "?" + r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SavingsSummary{OrgID: "org-1", Period: "7d", TotalCalls: 5, CostSavedUSD: 1.23})
	}))
	defer srv.Close()

	c := New(srv.URL)
	period := "7d"
	modelID := "claude-sonnet-4-6"
	got, err := c.GetSavingsSummary(context.Background(), "org-1", &period, &modelID)
	if err != nil {
		t.Fatalf("GetSavingsSummary() error = %v", err)
	}
	if got.TotalCalls != 5 {
		t.Fatalf("TotalCalls = %d, want 5", got.TotalCalls)
	}
	wantPath := "/api/v1/orgs/org-1/mcp/savings/summary?model_id=claude-sonnet-4-6&period=7d"
	if gotPath != wantPath {
		t.Fatalf("request path = %q, want %q", gotPath, wantPath)
	}
}

func TestGetSavingsSummary_OmitsParamsWhenNil(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SavingsSummary{})
	}))
	defer srv.Close()

	c := New(srv.URL)
	if _, err := c.GetSavingsSummary(context.Background(), "org-1", nil, nil); err != nil {
		t.Fatalf("GetSavingsSummary() error = %v", err)
	}
	if gotQuery != "" {
		t.Fatalf("query = %q, want empty when period/modelId are nil", gotQuery)
	}
}

func TestGetSavingsTimeseries_UnwrapsEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"timeseries":[{"totalCalls":2,"totalTokensSaved":300}]}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.GetSavingsTimeseries(context.Background(), "org-1", nil, nil)
	if err != nil {
		t.Fatalf("GetSavingsTimeseries() error = %v", err)
	}
	if len(got) != 1 || got[0].TotalCalls != 2 || got[0].TotalTokensSaved != 300 {
		t.Fatalf("GetSavingsTimeseries() = %+v, want one row with TotalCalls=2 TotalTokensSaved=300", got)
	}
}

func TestGetSavingsByTool_UnwrapsEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"byTool":[{"toolName":"get_api_spec","totalCalls":3}]}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.GetSavingsByTool(context.Background(), "org-1", nil, nil)
	if err != nil {
		t.Fatalf("GetSavingsByTool() error = %v", err)
	}
	if len(got) != 1 || got[0].ToolName != "get_api_spec" {
		t.Fatalf("GetSavingsByTool() = %+v, want one row with ToolName=get_api_spec", got)
	}
}

func TestGetSavingsByModel_UnwrapsEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"byModel":[{"modelId":"claude-haiku-4-5","displayName":"Claude Haiku 4.5"}]}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.GetSavingsByModel(context.Background(), "org-1", nil)
	if err != nil {
		t.Fatalf("GetSavingsByModel() error = %v", err)
	}
	if len(got) != 1 || got[0].DisplayName != "Claude Haiku 4.5" {
		t.Fatalf("GetSavingsByModel() = %+v, want one row with DisplayName=Claude Haiku 4.5", got)
	}
}

func TestGetSavingsByUser_UnwrapsEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"byUser":[{"userId":"u1","totalCalls":4}]}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.GetSavingsByUser(context.Background(), "org-1", nil, nil)
	if err != nil {
		t.Fatalf("GetSavingsByUser() error = %v", err)
	}
	if len(got) != 1 || got[0].UserID == nil || *got[0].UserID != "u1" {
		t.Fatalf("GetSavingsByUser() = %+v, want one row with UserID=u1", got)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/uigraphapi/... -run TestGetSavings -v`

Expected: FAIL — compile error, `SavingsSummary`, `GetSavingsSummary`, etc. are undefined in package `uigraphapi`.

- [ ] **Step 3: Implement the DTOs and client methods**

Create `internal/uigraphapi/mcpsavings.go`:

```go
package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type SavingsSummary struct {
	OrgID             string  `json:"orgId"`
	Period            string  `json:"period"`
	ModelID           string  `json:"modelId"`
	TotalCalls        int     `json:"totalCalls"`
	TotalTokensServed int     `json:"totalTokensServed"`
	TotalTokensSaved  int     `json:"totalTokensSaved"`
	CostServedUSD     float64 `json:"costServedUsd"`
	CostRawUSD        float64 `json:"costRawUsd"`
	CostSavedUSD      float64 `json:"costSavedUsd"`
	UniqueUsersCount  int     `json:"uniqueUsersCount"`
}

type DailySavings struct {
	Date              time.Time `json:"date"`
	TotalCalls        int       `json:"totalCalls"`
	TotalTokensServed int       `json:"totalTokensServed"`
	TotalTokensSaved  int       `json:"totalTokensSaved"`
	CostServedUSD     float64   `json:"costServedUsd"`
	CostRawUSD        float64   `json:"costRawUsd"`
	CostSavedUSD      float64   `json:"costSavedUsd"`
}

type ToolSavings struct {
	ToolName     string  `json:"toolName"`
	TotalCalls   int     `json:"totalCalls"`
	TokensSaved  int     `json:"tokensSaved"`
	CostSavedUSD float64 `json:"costSavedUsd"`
}

type ModelSavings struct {
	ModelID      string  `json:"modelId"`
	DisplayName  string  `json:"displayName"`
	Provider     string  `json:"provider"`
	TotalCalls   int     `json:"totalCalls"`
	TokensSaved  int     `json:"tokensSaved"`
	CostSavedUSD float64 `json:"costSavedUsd"`
}

type UserSavings struct {
	UserID           *string `json:"userId,omitempty"`
	ServiceAccountID *string `json:"serviceAccountId,omitempty"`
	TotalCalls       int     `json:"totalCalls"`
	TokensSaved      int     `json:"tokensSaved"`
	CostSavedUSD     float64 `json:"costSavedUsd"`
}

func savingsQuery(period, modelID *string) url.Values {
	q := url.Values{}
	if period != nil && *period != "" {
		q.Set("period", *period)
	}
	if modelID != nil && *modelID != "" {
		q.Set("model_id", *modelID)
	}
	return q
}

func withQuery(path string, q url.Values) string {
	if len(q) > 0 {
		return path + "?" + q.Encode()
	}
	return path
}

func (c *Client) GetSavingsSummary(ctx context.Context, orgID string, period, modelID *string) (*SavingsSummary, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/summary", orgID), savingsQuery(period, modelID))
	var out SavingsSummary
	return &out, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsTimeseries(ctx context.Context, orgID string, period, modelID *string) ([]DailySavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/timeseries", orgID), savingsQuery(period, modelID))
	var out struct {
		Timeseries []DailySavings `json:"timeseries"`
	}
	return out.Timeseries, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsByTool(ctx context.Context, orgID string, period, modelID *string) ([]ToolSavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/by-tool", orgID), savingsQuery(period, modelID))
	var out struct {
		ByTool []ToolSavings `json:"byTool"`
	}
	return out.ByTool, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsByModel(ctx context.Context, orgID string, period *string) ([]ModelSavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/by-model", orgID), savingsQuery(period, nil))
	var out struct {
		ByModel []ModelSavings `json:"byModel"`
	}
	return out.ByModel, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsByUser(ctx context.Context, orgID string, period, modelID *string) ([]UserSavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/by-user", orgID), savingsQuery(period, modelID))
	var out struct {
		ByUser []UserSavings `json:"byUser"`
	}
	return out.ByUser, c.get(ctx, path, &out)
}
```

(`savingsQuery`/`withQuery` are small shared helpers factored out since all five methods build the same `period`/`model_id` query string — avoids repeating the same four `if` lines five times.)

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/uigraphapi/... -run TestGetSavings -v`

Expected: `PASS` for all 6 tests.

- [ ] **Step 5: Run the package's full test suite to confirm no regressions**

Run: `go test ./internal/uigraphapi/...`

Expected: `ok`.

- [ ] **Step 6: Commit**

```bash
git add internal/uigraphapi/mcpsavings.go internal/uigraphapi/mcpsavings_test.go
git commit -m "feat: add MCP savings REST client methods"
```

---

### Task 2: GraphQL schema + codegen

**Files:**
- Create: `internal/graph/schema/mcpsavings.graphqls`
- Generated (do not hand-edit): `internal/graph/generated/generated.go`, `internal/graph/model/models_gen.go`, `internal/graph/mcpsavings.resolvers.go`

**Interfaces:**
- Produces: gqlgen models `model.SavingsSummary`, `model.DailySavings`, `model.ToolSavings`, `model.ModelSavings`, `model.UserSavings` — consumed by Task 3 (convert) and Task 4 (resolvers) — and five stub resolver methods on `*queryResolver` (bodies are gqlgen's default `panic(fmt.Errorf("not implemented..."))` until Task 4 fills them in)

- [ ] **Step 1: Write the schema file**

Create `internal/graph/schema/mcpsavings.graphqls`:

```graphql
extend type Query {
    costSavingsSummary(orgId: ID!, period: String, modelId: String): SavingsSummary!
    costSavingsTimeseries(orgId: ID!, period: String, modelId: String): [DailySavings!]!
    costSavingsByTool(orgId: ID!, period: String, modelId: String): [ToolSavings!]!
    costSavingsByModel(orgId: ID!, period: String): [ModelSavings!]!
    costSavingsByUser(orgId: ID!, period: String, modelId: String): [UserSavings!]!
}

type SavingsSummary {
    orgId:             ID!
    period:            String!
    modelId:           String
    totalCalls:        Int!
    totalTokensServed: Int!
    totalTokensSaved:  Int!
    costServedUsd:     Float!
    costRawUsd:        Float!
    costSavedUsd:      Float!
    uniqueUsersCount:  Int!
}

type DailySavings {
    date:              Time!
    totalCalls:        Int!
    totalTokensServed: Int!
    totalTokensSaved:  Int!
    costServedUsd:     Float!
    costRawUsd:        Float!
    costSavedUsd:      Float!
}

type ToolSavings {
    toolName:     String!
    totalCalls:   Int!
    tokensSaved:  Int!
    costSavedUsd: Float!
}

type ModelSavings {
    modelId:      String!
    displayName:  String!
    provider:     String!
    totalCalls:   Int!
    tokensSaved:  Int!
    costSavedUsd: Float!
}

type UserSavings {
    userId:           ID
    serviceAccountId: ID
    displayName:      String!
    totalCalls:       Int!
    tokensSaved:      Int!
    costSavedUsd:     Float!
}
```

- [ ] **Step 2: Regenerate gqlgen code**

Run: `make generate` (equivalently `go generate ./internal/graph/...`)

Expected: no errors. Confirm `internal/graph/mcpsavings.resolvers.go` now exists and confirm `internal/graph/model/models_gen.go` now contains a `SavingsSummary` struct, a `DailySavings` struct, a `ToolSavings` struct, a `ModelSavings` struct, and a `UserSavings` struct with fields matching the schema above (gqlgen's default Go field naming: `OrgID`, `Period`, `ModelID *string`, `TotalCalls`, etc.).

- [ ] **Step 3: Confirm the project builds with the stub resolvers**

Run: `go build ./...`

Expected: succeeds. (The generated resolver stub bodies `panic(fmt.Errorf("not implemented: ..."))` are fine at this stage — `go build` only checks that signatures type-check, not that anything is implemented.)

- [ ] **Step 4: Commit**

```bash
git add internal/graph/schema/mcpsavings.graphqls internal/graph/generated/generated.go internal/graph/model/models_gen.go internal/graph/mcpsavings.resolvers.go
git commit -m "feat: add MCP savings GraphQL schema and generated scaffolding"
```

---

### Task 3: Convert layer (`internal/graph/convert`)

**Files:**
- Create: `internal/graph/convert/mcpsavings.go`
- Create: `internal/graph/convert/mcpsavings_test.go`

**Interfaces:**
- Consumes: `uigraphapi.SavingsSummary`, `DailySavings`, `ToolSavings`, `ModelSavings`, `UserSavings`, and `uigraphapi.Actor` (existing type, has a `Name string` field — see `internal/graph/refs.go`'s `resolveActor`) from Task 1, and `model.SavingsSummary`/`model.DailySavings`/etc. from Task 2
- Produces: `SavingsSummaryToModel(s *uigraphapi.SavingsSummary) *model.SavingsSummary`, `DailySavingsListToModel([]uigraphapi.DailySavings) []*model.DailySavings`, `ToolSavingsListToModel(...)`, `ModelSavingsListToModel(...)`, `UserSavingsListToModel(rows []uigraphapi.UserSavings, actors map[string]*uigraphapi.Actor) []*model.UserSavings` — all consumed by Task 4

- [ ] **Step 1: Write the failing tests**

Create `internal/graph/convert/mcpsavings_test.go`:

```go
package convert

import (
	"testing"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestSavingsSummaryToModel(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		if got := SavingsSummaryToModel(nil); got != nil {
			t.Fatalf("SavingsSummaryToModel(nil) = %v, want nil", got)
		}
	})

	t.Run("empty ModelID becomes nil (blended)", func(t *testing.T) {
		got := SavingsSummaryToModel(&uigraphapi.SavingsSummary{OrgID: "o1", Period: "7d", ModelID: "", TotalCalls: 3})
		if got.ModelID != nil {
			t.Errorf("ModelID = %v, want nil for blended summary", got.ModelID)
		}
		if got.TotalCalls != 3 {
			t.Errorf("TotalCalls = %d, want 3", got.TotalCalls)
		}
	})

	t.Run("non-empty ModelID is preserved as pointer", func(t *testing.T) {
		got := SavingsSummaryToModel(&uigraphapi.SavingsSummary{ModelID: "claude-sonnet-4-6"})
		if got.ModelID == nil || *got.ModelID != "claude-sonnet-4-6" {
			t.Errorf("ModelID = %v, want pointer to claude-sonnet-4-6", got.ModelID)
		}
	})
}

func TestDailySavingsListToModel(t *testing.T) {
	t.Run("maps each row in order", func(t *testing.T) {
		now := time.Now()
		got := DailySavingsListToModel([]uigraphapi.DailySavings{
			{Date: now, TotalCalls: 1, TotalTokensSaved: 10},
			{Date: now.AddDate(0, 0, 1), TotalCalls: 2, TotalTokensSaved: 20},
		})
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
		if got[0].TotalCalls != 1 || got[1].TotalCalls != 2 {
			t.Errorf("TotalCalls = [%d, %d], want [1, 2]", got[0].TotalCalls, got[1].TotalCalls)
		}
	})

	t.Run("empty input returns empty slice, not nil", func(t *testing.T) {
		got := DailySavingsListToModel(nil)
		if got == nil {
			t.Fatal("got nil, want empty slice")
		}
		if len(got) != 0 {
			t.Errorf("len = %d, want 0", len(got))
		}
	})
}

func TestUserSavingsToModel(t *testing.T) {
	t.Run("resolves display name from actor map for a user", func(t *testing.T) {
		uid := "u1"
		actors := map[string]*uigraphapi.Actor{"u1": {ID: "u1", Name: "Ada Lovelace"}}
		got := UserSavingsToModel(uigraphapi.UserSavings{UserID: &uid, TotalCalls: 5}, actors)
		if got.DisplayName != "Ada Lovelace" {
			t.Errorf("DisplayName = %q, want Ada Lovelace", got.DisplayName)
		}
	})

	t.Run("falls back to Service Account when actor not found", func(t *testing.T) {
		said := "sa1"
		got := UserSavingsToModel(uigraphapi.UserSavings{ServiceAccountID: &said}, map[string]*uigraphapi.Actor{})
		if got.DisplayName != "Service Account" {
			t.Errorf("DisplayName = %q, want Service Account", got.DisplayName)
		}
	})

	t.Run("falls back to Unknown User when neither id nor actor resolves", func(t *testing.T) {
		got := UserSavingsToModel(uigraphapi.UserSavings{}, map[string]*uigraphapi.Actor{})
		if got.DisplayName != "Unknown User" {
			t.Errorf("DisplayName = %q, want Unknown User", got.DisplayName)
		}
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/graph/convert/... -run "TestSavingsSummaryToModel|TestDailySavingsListToModel|TestUserSavingsToModel" -v`

Expected: FAIL — compile error, `SavingsSummaryToModel` etc. are undefined in package `convert` (the `model.SavingsSummary` etc. types referenced in the test already exist from Task 2's codegen).

- [ ] **Step 3: Implement the convert functions**

Create `internal/graph/convert/mcpsavings.go`:

```go
package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func SavingsSummaryToModel(s *uigraphapi.SavingsSummary) *model.SavingsSummary {
	if s == nil {
		return nil
	}
	var modelID *string
	if s.ModelID != "" {
		modelID = &s.ModelID
	}
	return &model.SavingsSummary{
		OrgID:             s.OrgID,
		Period:            s.Period,
		ModelID:           modelID,
		TotalCalls:        s.TotalCalls,
		TotalTokensServed: s.TotalTokensServed,
		TotalTokensSaved:  s.TotalTokensSaved,
		CostServedUSD:     s.CostServedUSD,
		CostRawUSD:        s.CostRawUSD,
		CostSavedUSD:      s.CostSavedUSD,
		UniqueUsersCount:  s.UniqueUsersCount,
	}
}

func DailySavingsToModel(d uigraphapi.DailySavings) *model.DailySavings {
	return &model.DailySavings{
		Date:              d.Date,
		TotalCalls:        d.TotalCalls,
		TotalTokensServed: d.TotalTokensServed,
		TotalTokensSaved:  d.TotalTokensSaved,
		CostServedUSD:     d.CostServedUSD,
		CostRawUSD:        d.CostRawUSD,
		CostSavedUSD:      d.CostSavedUSD,
	}
}

func DailySavingsListToModel(rows []uigraphapi.DailySavings) []*model.DailySavings {
	out := make([]*model.DailySavings, len(rows))
	for i, row := range rows {
		out[i] = DailySavingsToModel(row)
	}
	return out
}

func ToolSavingsToModel(s uigraphapi.ToolSavings) *model.ToolSavings {
	return &model.ToolSavings{
		ToolName:     s.ToolName,
		TotalCalls:   s.TotalCalls,
		TokensSaved:  s.TokensSaved,
		CostSavedUSD: s.CostSavedUSD,
	}
}

func ToolSavingsListToModel(rows []uigraphapi.ToolSavings) []*model.ToolSavings {
	out := make([]*model.ToolSavings, len(rows))
	for i, row := range rows {
		out[i] = ToolSavingsToModel(row)
	}
	return out
}

func ModelSavingsToModel(s uigraphapi.ModelSavings) *model.ModelSavings {
	return &model.ModelSavings{
		ModelID:      s.ModelID,
		DisplayName:  s.DisplayName,
		Provider:     s.Provider,
		TotalCalls:   s.TotalCalls,
		TokensSaved:  s.TokensSaved,
		CostSavedUSD: s.CostSavedUSD,
	}
}

func ModelSavingsListToModel(rows []uigraphapi.ModelSavings) []*model.ModelSavings {
	out := make([]*model.ModelSavings, len(rows))
	for i, row := range rows {
		out[i] = ModelSavingsToModel(row)
	}
	return out
}

// UserSavingsToModel resolves DisplayName from actors (keyed by user ID or
// service account ID), falling back to "Service Account" or "Unknown User"
// when no actor was resolved for that ID.
func UserSavingsToModel(s uigraphapi.UserSavings, actors map[string]*uigraphapi.Actor) *model.UserSavings {
	id := ""
	if s.UserID != nil {
		id = *s.UserID
	} else if s.ServiceAccountID != nil {
		id = *s.ServiceAccountID
	}
	displayName := "Unknown User"
	if a := actors[id]; a != nil {
		displayName = a.Name
	} else if s.ServiceAccountID != nil {
		displayName = "Service Account"
	}
	return &model.UserSavings{
		UserID:           s.UserID,
		ServiceAccountID: s.ServiceAccountID,
		DisplayName:      displayName,
		TotalCalls:       s.TotalCalls,
		TokensSaved:      s.TokensSaved,
		CostSavedUSD:     s.CostSavedUSD,
	}
}

func UserSavingsListToModel(rows []uigraphapi.UserSavings, actors map[string]*uigraphapi.Actor) []*model.UserSavings {
	out := make([]*model.UserSavings, len(rows))
	for i, row := range rows {
		out[i] = UserSavingsToModel(row, actors)
	}
	return out
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/graph/convert/... -run "TestSavingsSummaryToModel|TestDailySavingsListToModel|TestUserSavingsToModel" -v`

Expected: `PASS` for all subtests.

- [ ] **Step 5: Commit**

```bash
git add internal/graph/convert/mcpsavings.go internal/graph/convert/mcpsavings_test.go
git commit -m "feat: add MCP savings REST-to-GraphQL convert functions"
```

---

### Task 4: Resolver implementation + wiring

**Files:**
- Modify: `internal/graph/mcpsavings.resolvers.go` (replace the 5 generated panic stubs with real implementations)
- Modify: `internal/graph/resolver.go` (add `costSavingsClient` interface + `CostSavings` field on `Resolver`)
- Modify: `internal/server/server.go` (wire `CostSavings: c,` into the `Resolver` literal)
- Create: `internal/graph/mcpsavings_test.go`

**Interfaces:**
- Consumes: Task 1's `uigraphapi.Client` methods, Task 3's `convert` functions, the existing `actorClient` interface's `ResolveActors(ctx, orgID string, ids []string) (map[string]*uigraphapi.Actor, error)` (already wired as `Resolver.Actor`, see `internal/graph/refs.go`)

- [ ] **Step 1: Write the failing resolver tests**

Create `internal/graph/mcpsavings_test.go`:

```go
package graph_test

import (
	"context"
	"testing"

	"github.com/uigraph/graphql/internal/graph"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

type fakeCostSavingsClient struct {
	summary    *uigraphapi.SavingsSummary
	timeseries []uigraphapi.DailySavings
	byTool     []uigraphapi.ToolSavings
	byModel    []uigraphapi.ModelSavings
	byUser     []uigraphapi.UserSavings
}

func (f *fakeCostSavingsClient) GetSavingsSummary(_ context.Context, _ string, _, _ *string) (*uigraphapi.SavingsSummary, error) {
	return f.summary, nil
}
func (f *fakeCostSavingsClient) GetSavingsTimeseries(_ context.Context, _ string, _, _ *string) ([]uigraphapi.DailySavings, error) {
	return f.timeseries, nil
}
func (f *fakeCostSavingsClient) GetSavingsByTool(_ context.Context, _ string, _, _ *string) ([]uigraphapi.ToolSavings, error) {
	return f.byTool, nil
}
func (f *fakeCostSavingsClient) GetSavingsByModel(_ context.Context, _ string, _ *string) ([]uigraphapi.ModelSavings, error) {
	return f.byModel, nil
}
func (f *fakeCostSavingsClient) GetSavingsByUser(_ context.Context, _ string, _, _ *string) ([]uigraphapi.UserSavings, error) {
	return f.byUser, nil
}

type fakeActorClient struct {
	actors map[string]*uigraphapi.Actor
}

func (f *fakeActorClient) ResolveActors(_ context.Context, _ string, ids []string) (map[string]*uigraphapi.Actor, error) {
	out := map[string]*uigraphapi.Actor{}
	for _, id := range ids {
		if a, ok := f.actors[id]; ok {
			out[id] = a
		}
	}
	return out, nil
}
func (f *fakeActorClient) ResolveAssetURLs(_ context.Context, _ string, _ []string) (map[string]string, error) {
	return nil, nil
}
func (f *fakeActorClient) CreateAssetUpload(_ context.Context, _ string) (*uigraphapi.AssetUpload, error) {
	return nil, nil
}

func TestCostSavingsSummaryQuery(t *testing.T) {
	resolver := &graph.Resolver{
		CostSavings: &fakeCostSavingsClient{
			summary: &uigraphapi.SavingsSummary{OrgID: "org-1", Period: "7d", TotalCalls: 9, CostSavedUSD: 4.5},
		},
	}
	srv := newTestServer(resolver)
	defer srv.Close()

	data := doGraphQL(t, srv, `{ costSavingsSummary(orgId: "org-1") { orgId period totalCalls costSavedUsd modelId } }`)
	summary, ok := data["costSavingsSummary"].(map[string]interface{})
	if !ok {
		t.Fatalf("costSavingsSummary field missing or wrong type: %+v", data)
	}
	if summary["totalCalls"] != float64(9) {
		t.Errorf("totalCalls: got %v, want 9", summary["totalCalls"])
	}
	if summary["modelId"] != nil {
		t.Errorf("modelId: got %v, want nil for blended summary", summary["modelId"])
	}
}

func TestCostSavingsByUserQuery_ResolvesDisplayNames(t *testing.T) {
	uid := "u1"
	resolver := &graph.Resolver{
		CostSavings: &fakeCostSavingsClient{
			byUser: []uigraphapi.UserSavings{{UserID: &uid, TotalCalls: 3, TokensSaved: 100, CostSavedUSD: 1.0}},
		},
		Actor: &fakeActorClient{actors: map[string]*uigraphapi.Actor{"u1": {ID: "u1", Name: "Ada Lovelace"}}},
	}
	srv := newTestServer(resolver)
	defer srv.Close()

	data := doGraphQL(t, srv, `{ costSavingsByUser(orgId: "org-1") { userId displayName totalCalls } }`)
	rows, ok := data["costSavingsByUser"].([]interface{})
	if !ok || len(rows) != 1 {
		t.Fatalf("costSavingsByUser field missing or wrong shape: %+v", data)
	}
	row := rows[0].(map[string]interface{})
	if row["displayName"] != "Ada Lovelace" {
		t.Errorf("displayName: got %v, want Ada Lovelace", row["displayName"])
	}
}

func TestCostSavingsByUserQuery_FallsBackForServiceAccount(t *testing.T) {
	said := "sa1"
	resolver := &graph.Resolver{
		CostSavings: &fakeCostSavingsClient{
			byUser: []uigraphapi.UserSavings{{ServiceAccountID: &said, TotalCalls: 7}},
		},
		Actor: &fakeActorClient{actors: map[string]*uigraphapi.Actor{}},
	}
	srv := newTestServer(resolver)
	defer srv.Close()

	data := doGraphQL(t, srv, `{ costSavingsByUser(orgId: "org-1") { serviceAccountId displayName } }`)
	rows := data["costSavingsByUser"].([]interface{})
	row := rows[0].(map[string]interface{})
	if row["displayName"] != "Service Account" {
		t.Errorf("displayName: got %v, want Service Account", row["displayName"])
	}
}
```

(`newTestServer`/`doGraphQL` are already defined in `internal/graph/resolver_test.go`, same `graph_test` package — no need to redeclare them.)

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/graph/... -run TestCostSavings -v`

Expected: FAIL — compile error, `graph.Resolver` has no field `CostSavings`, and `fakeCostSavingsClient` doesn't satisfy any known interface yet.

- [ ] **Step 3: Add the `costSavingsClient` interface and `Resolver` field**

In `internal/graph/resolver.go`, add the interface (alongside the other domain interfaces like `testPackClient`):

```go
type costSavingsClient interface {
	GetSavingsSummary(ctx context.Context, orgID string, period, modelID *string) (*uigraphapi.SavingsSummary, error)
	GetSavingsTimeseries(ctx context.Context, orgID string, period, modelID *string) ([]uigraphapi.DailySavings, error)
	GetSavingsByTool(ctx context.Context, orgID string, period, modelID *string) ([]uigraphapi.ToolSavings, error)
	GetSavingsByModel(ctx context.Context, orgID string, period *string) ([]uigraphapi.ModelSavings, error)
	GetSavingsByUser(ctx context.Context, orgID string, period, modelID *string) ([]uigraphapi.UserSavings, error)
}
```

And add a field to the `Resolver` struct:

```go
type Resolver struct {
	Auth        authClient
	OrgAPI      orgClient
	Admin       adminClient
	FolderAPI   folderClient
	DiagramAPI  diagramClient
	DocAPI      docsClient
	Component   componentClient
	UIMapAPI    uimapClient
	Catalog     catalogClient
	TestPack    testPackClient
	Actor       actorClient
	CommentAPI  commentClient
	CostSavings costSavingsClient
}
```

- [ ] **Step 4: Implement the resolvers**

Replace the 5 generated stub method bodies in `internal/graph/mcpsavings.resolvers.go` (keep the generated file's header comment and imports block as-is, just replace the panicking bodies):

```go
package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

// CostSavingsSummary is the resolver for the costSavingsSummary field.
func (r *queryResolver) CostSavingsSummary(ctx context.Context, orgID string, period *string, modelID *string) (*model.SavingsSummary, error) {
	s, err := r.CostSavings.GetSavingsSummary(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}
	return convert.SavingsSummaryToModel(s), nil
}

// CostSavingsTimeseries is the resolver for the costSavingsTimeseries field.
func (r *queryResolver) CostSavingsTimeseries(ctx context.Context, orgID string, period *string, modelID *string) ([]*model.DailySavings, error) {
	rows, err := r.CostSavings.GetSavingsTimeseries(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}
	return convert.DailySavingsListToModel(rows), nil
}

// CostSavingsByTool is the resolver for the costSavingsByTool field.
func (r *queryResolver) CostSavingsByTool(ctx context.Context, orgID string, period *string, modelID *string) ([]*model.ToolSavings, error) {
	rows, err := r.CostSavings.GetSavingsByTool(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}
	return convert.ToolSavingsListToModel(rows), nil
}

// CostSavingsByModel is the resolver for the costSavingsByModel field.
func (r *queryResolver) CostSavingsByModel(ctx context.Context, orgID string, period *string) ([]*model.ModelSavings, error) {
	rows, err := r.CostSavings.GetSavingsByModel(ctx, orgID, period)
	if err != nil {
		return nil, err
	}
	return convert.ModelSavingsListToModel(rows), nil
}

// CostSavingsByUser is the resolver for the costSavingsByUser field.
func (r *queryResolver) CostSavingsByUser(ctx context.Context, orgID string, period *string, modelID *string) ([]*model.UserSavings, error) {
	rows, err := r.CostSavings.GetSavingsByUser(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.UserID != nil {
			ids = append(ids, *row.UserID)
		} else if row.ServiceAccountID != nil {
			ids = append(ids, *row.ServiceAccountID)
		}
	}

	actors := map[string]*uigraphapi.Actor{}
	if len(ids) > 0 {
		var actorErr error
		actors, actorErr = r.Actor.ResolveActors(ctx, orgID, ids)
		if actorErr != nil {
			return nil, actorErr
		}
	}
	return convert.UserSavingsListToModel(rows, actors), nil
}
```

- [ ] **Step 5: Wire the client into `server.go`**

In `internal/server/server.go`, add `CostSavings: c,` to the `&graph.Resolver{...}` literal:

```go
resolver := &graph.Resolver{
	Auth:        c,
	OrgAPI:      c,
	Admin:       c,
	FolderAPI:   c,
	DiagramAPI:  c,
	DocAPI:      c,
	Component:   c,
	UIMapAPI:    c,
	Catalog:     c,
	TestPack:    c,
	Actor:       c,
	CommentAPI:  c,
	CostSavings: c,
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./internal/graph/... -run TestCostSavings -v`

Expected: `PASS` for all 3 tests.

- [ ] **Step 7: Run the full test suite and build**

Run: `go build ./... && go test ./...`

Expected: build succeeds, all packages `ok`.

- [ ] **Step 8: Commit**

```bash
git add internal/graph/mcpsavings.resolvers.go internal/graph/resolver.go internal/server/server.go internal/graph/mcpsavings_test.go
git commit -m "feat: implement MCP savings GraphQL resolvers"
```

---

## Summary of new GraphQL surface (for the uigraph-ui plan)

```graphql
costSavingsSummary(orgId: ID!, period: String, modelId: String): SavingsSummary!
costSavingsTimeseries(orgId: ID!, period: String, modelId: String): [DailySavings!]!
costSavingsByTool(orgId: ID!, period: String, modelId: String): [ToolSavings!]!
costSavingsByModel(orgId: ID!, period: String): [ModelSavings!]!
costSavingsByUser(orgId: ID!, period: String, modelId: String): [UserSavings!]!
```

`SavingsSummary` fields: `orgId period modelId totalCalls totalTokensServed totalTokensSaved costServedUsd costRawUsd costSavedUsd uniqueUsersCount`
`DailySavings` fields: `date totalCalls totalTokensServed totalTokensSaved costServedUsd costRawUsd costSavedUsd`
`ToolSavings` fields: `toolName totalCalls tokensSaved costSavedUsd`
`ModelSavings` fields: `modelId displayName provider totalCalls tokensSaved costSavedUsd`
`UserSavings` fields: `userId serviceAccountId displayName totalCalls tokensSaved costSavedUsd`
