package uigraphapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDependencyClientRequests(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/orgs/org-1/services/service-1/dependencies":
			if r.URL.Query().Get("direction") != "outbound" || r.URL.Query().Get("criticality") != "high" {
				t.Errorf("query = %q, want direction=outbound and criticality=high", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"edges": []Dependency{{ID: "dependency-1", Name: "Payments"}}})
		case "/api/v1/orgs/org-1/services/service-1/dependency-graph", "/api/v1/orgs/org-1/dependency-graph":
			_ = json.NewEncoder(w).Encode(DependencyGraph{Nodes: []DependencyGraphNode{{ID: "service-1", Name: "Checkout"}}})
		case "/api/v1/orgs/org-1/services/service-1/impact":
			if r.URL.Query().Get("direction") != "inbound" || r.URL.Query().Get("maxDepth") != "3" {
				t.Errorf("query = %q, want direction=inbound and maxDepth=3", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(DependencyGraph{Edges: []DependencyGraphEdge{{ID: "edge-1", Source: "a", Target: "b"}}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := New(srv.URL)
	direction := "outbound"
	criticality := "high"
	dependencies, err := c.ListDependencies(context.Background(), "org-1", "service-1", &direction, &criticality)
	if err != nil {
		t.Fatalf("ListDependencies() error = %v", err)
	}
	if len(dependencies) != 1 || dependencies[0].ID != "dependency-1" {
		t.Fatalf("ListDependencies() = %+v, want dependency-1", dependencies)
	}
	if _, err := c.GetServiceDependencyGraph(context.Background(), "org-1", "service-1"); err != nil {
		t.Fatalf("GetServiceDependencyGraph() error = %v", err)
	}
	if _, err := c.GetDependencyGraph(context.Background(), "org-1"); err != nil {
		t.Fatalf("GetDependencyGraph() error = %v", err)
	}
	inbound := "inbound"
	maxDepth := 3
	impact, err := c.GetServiceImpact(context.Background(), "org-1", "service-1", &inbound, &maxDepth)
	if err != nil {
		t.Fatalf("GetServiceImpact() error = %v", err)
	}
	if len(impact.Edges) != 1 || impact.Edges[0].ID != "edge-1" {
		t.Fatalf("GetServiceImpact() = %+v, want edge-1", impact)
	}
}
