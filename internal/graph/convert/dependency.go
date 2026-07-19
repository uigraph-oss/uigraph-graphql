package convert

import (
	"encoding/json"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func DependencyServiceToModel(service *uigraphapi.DependencyService) *model.DependencyService {
	if service == nil {
		return nil
	}
	return &model.DependencyService{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		Status:      service.Status,
		Tier:        service.Tier,
		Category:    service.Category,
		Language:    service.Language,
		GitRepoURL:  service.GitRepoURL,
		UpdatedAt:   service.UpdatedAt,
		Metadata:    rawJSON(service.Metadata),
	}
}

func DependencyToModel(dependency uigraphapi.Dependency) *model.Dependency {
	return &model.Dependency{
		ID:               dependency.ID,
		Name:             dependency.Name,
		ConsumerService:  DependencyServiceToModel(&dependency.ConsumerService),
		ProviderService:  DependencyServiceToModel(dependency.ProviderService),
		ProviderName:     dependency.ProviderName,
		Type:             dependency.Type,
		Criticality:      dependency.Criticality,
		Description:      dependency.Description,
		APIGroupName:     dependency.APIGroupName,
		APIEndpointNames: dependency.APIEndpointNames,
		DatabaseName:     dependency.DatabaseName,
		Direction:        dependency.Direction,
	}
}

func DependenciesToModel(dependencies []uigraphapi.Dependency) []*model.Dependency {
	out := make([]*model.Dependency, len(dependencies))
	for i, dependency := range dependencies {
		out[i] = DependencyToModel(dependency)
	}
	return out
}

func rawJSON(raw json.RawMessage) interface{} {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var value interface{}
	if json.Unmarshal(raw, &value) != nil {
		return nil
	}
	return value
}
