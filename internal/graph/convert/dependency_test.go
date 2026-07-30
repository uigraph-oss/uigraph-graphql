package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestDependencyToModel(t *testing.T) {
	apiGroupName := "v1"
	dependency := DependencyToModel(uigraphapi.Dependency{
		ID:               "dependency-1",
		Name:             "Payments",
		Service:          &uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"},
		DependencyName:   "External payments",
		Direction:        "downstream",
		APIGroupName:     &apiGroupName,
		APIEndpointNames: []string{"createPayment"},
	})
	if dependency.Service == nil || dependency.Service.ID != "service-1" {
		t.Errorf("Service = %v, want id service-1", dependency.Service)
	}
	if dependency.Dependency != nil {
		t.Errorf("Dependency = %v, want nil", dependency.Dependency)
	}
	if dependency.DependencyName != "External payments" {
		t.Errorf("DependencyName = %q, want External payments", dependency.DependencyName)
	}
	if dependency.Direction != "downstream" {
		t.Errorf("Direction = %q, want down", dependency.Direction)
	}
	if dependency.APIGroupName == nil || *dependency.APIGroupName != apiGroupName {
		t.Errorf("APIGroupName = %v, want %q", dependency.APIGroupName, apiGroupName)
	}
	if len(dependency.APIEndpointNames) != 1 || dependency.APIEndpointNames[0] != "createPayment" {
		t.Errorf("APIEndpointNames = %#v, want [createPayment]", dependency.APIEndpointNames)
	}
}

func TestDependencyToModel_directionIsReadVerbatim(t *testing.T) {
	for _, direction := range []string{"upstream", "downstream"} {
		dependency := DependencyToModel(uigraphapi.Dependency{
			ID:             "dependency-1",
			Name:           "Peer",
			Service:        &uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"},
			Dependency:     &uigraphapi.DependencyService{ID: "service-2", Name: "Payments"},
			DependencyName: "Payments",
			Direction:      direction,
		})
		if dependency.Direction != direction {
			t.Errorf("Direction = %q, want %q", dependency.Direction, direction)
		}
	}
}
