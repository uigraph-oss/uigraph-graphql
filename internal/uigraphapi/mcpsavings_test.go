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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
