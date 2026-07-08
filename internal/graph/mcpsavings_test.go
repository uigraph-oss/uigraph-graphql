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
	byClient   []uigraphapi.ClientSavings
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
func (f *fakeCostSavingsClient) GetSavingsByClient(_ context.Context, _ string, _, _ *string) ([]uigraphapi.ClientSavings, error) {
	return f.byClient, nil
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
