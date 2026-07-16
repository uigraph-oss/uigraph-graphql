package convert

import (
	"encoding/json"
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestDependencyToModel(t *testing.T) {
	providerName := "External payments"
	dependency := DependencyToModel(uigraphapi.Dependency{
		ID:              "dependency-1",
		Name:            "Payments",
		ConsumerService: uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"},
		ProviderName:    &providerName,
		API:             json.RawMessage(`{"protocol":"http"}`),
		Operations:      json.RawMessage(`[{"method":"POST"}]`),
	})
	if dependency.ConsumerService.ID != "service-1" {
		t.Errorf("ConsumerService.ID = %q, want service-1", dependency.ConsumerService.ID)
	}
	if dependency.ProviderService != nil {
		t.Errorf("ProviderService = %v, want nil", dependency.ProviderService)
	}
	if dependency.ProviderName == nil || *dependency.ProviderName != providerName {
		t.Errorf("ProviderName = %v, want %q", dependency.ProviderName, providerName)
	}
	api, ok := dependency.API.(map[string]interface{})
	if !ok || api["protocol"] != "http" {
		t.Errorf("API = %#v, want protocol http", dependency.API)
	}
	operations, ok := dependency.Operations.([]interface{})
	if !ok || len(operations) != 1 {
		t.Errorf("Operations = %#v, want one operation", dependency.Operations)
	}
}

func TestDependencyGraphToModel(t *testing.T) {
	depth := 2
	graph := DependencyGraphToModel(&uigraphapi.DependencyGraph{
		Nodes: []uigraphapi.DependencyGraphNode{{
			ID: "service-1", Name: "Checkout", Depth: &depth,
			Service:  &uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"},
			Metadata: json.RawMessage(`{"team":"core"}`),
		}},
		Edges: []uigraphapi.DependencyGraphEdge{{ID: "edge-1", Source: "service-1", Target: "service-2"}},
	})
	if graph.Nodes[0].Depth == nil || *graph.Nodes[0].Depth != depth {
		t.Errorf("Nodes[0].Depth = %v, want %d", graph.Nodes[0].Depth, depth)
	}
	if graph.Nodes[0].Service == nil || graph.Nodes[0].Service.Name != "Checkout" {
		t.Errorf("Nodes[0].Service = %+v, want Checkout", graph.Nodes[0].Service)
	}
	if graph.Edges[0].Source != "service-1" || graph.Edges[0].Target != "service-2" {
		t.Errorf("Edges[0] = %+v, want service-1 -> service-2", graph.Edges[0])
	}
}
