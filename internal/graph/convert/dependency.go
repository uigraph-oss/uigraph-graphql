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
		OnboardingStatus: dependency.OnboardingStatus,
		Type:             dependency.Type,
		Criticality:      dependency.Criticality,
		Description:      dependency.Description,
		API:              rawJSON(dependency.API),
		Operations:       rawJSON(dependency.Operations),
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

func DependencyGraphToModel(graph *uigraphapi.DependencyGraph) *model.DependencyGraph {
	nodes := make([]*model.DependencyGraphNode, len(graph.Nodes))
	for i, node := range graph.Nodes {
		nodes[i] = &model.DependencyGraphNode{
			ID:               node.ID,
			Name:             node.Name,
			Type:             node.Type,
			Service:          DependencyServiceToModel(node.Service),
			OnboardingStatus: node.OnboardingStatus,
			Depth:            node.Depth,
			Metadata:         rawJSON(node.Metadata),
		}
	}
	edges := make([]*model.DependencyGraphEdge, len(graph.Edges))
	for i, edge := range graph.Edges {
		source := edge.Source
		if source == "" {
			source = edge.SourceServiceID
		}
		if source == "" && edge.Consumer != nil {
			source = edge.Consumer.ID
		}
		target := edge.Target
		if target == "" && edge.Provider != nil {
			target = edge.Provider.ID
		}
		if target == "" {
			target = "ghost:" + edge.ProviderServiceName
		}
		edges[i] = &model.DependencyGraphEdge{
			ID:           edge.ID,
			Source:       source,
			Target:       target,
			DependencyID: edge.DependencyID,
			Type:         edge.Type,
			Criticality:  edge.Criticality,
			Direction:    edge.Direction,
			Depth:        edge.Depth,
			Operations:   rawJSON(edge.Operations),
			Metadata:     rawJSON(edge.Metadata),
		}
	}
	return &model.DependencyGraph{Nodes: nodes, Edges: edges}
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
