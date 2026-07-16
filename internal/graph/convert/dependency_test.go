package convert

import (
	"encoding/json"
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestDependencyToModel(t *testing.T) {
	providerName := "External payments"
	apiGroupName := "v1"
	dependency := DependencyToModel(uigraphapi.Dependency{
		ID:               "dependency-1",
		Name:             "Payments",
		ConsumerService:  uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"},
		ProviderName:     &providerName,
		APIGroupName:     &apiGroupName,
		APIEndpointNames: []string{"createPayment"},
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
	if dependency.APIGroupName == nil || *dependency.APIGroupName != apiGroupName {
		t.Errorf("APIGroupName = %v, want %q", dependency.APIGroupName, apiGroupName)
	}
	if len(dependency.APIEndpointNames) != 1 || dependency.APIEndpointNames[0] != "createPayment" {
		t.Errorf("APIEndpointNames = %#v, want [createPayment]", dependency.APIEndpointNames)
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
