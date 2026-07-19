package convert

import (
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
