package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *mutationResolver) CreateCustomComponent(ctx context.Context, orgID string, input model.CustomComponentInput) (*model.Component, error) {
	c, err := r.Component.CreateCustomComponent(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ComponentToModel(*c), nil
}

func (r *mutationResolver) UpdateCustomComponent(ctx context.Context, orgID string, id string, input model.CustomComponentInput) (*model.Component, error) {
	c, err := r.Component.UpdateCustomComponent(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ComponentToModel(*c), nil
}

func (r *mutationResolver) DeleteCustomComponent(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.Component.DeleteCustomComponent(ctx, orgID, id)
}

func (r *queryResolver) FlowDiagramComponents(ctx context.Context, orgID string) (*model.FlowDiagramComponents, error) {
	res, err := r.Component.ListFlowDiagramComponents(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return &model.FlowDiagramComponents{
		Components:       convert.FlowComponentsToModel(res.Components),
		CustomComponents: convert.FlowComponentsToModel(res.CustomComponents),
	}, nil
}

func (r *queryResolver) Components(ctx context.Context, orgID string) (*model.Components, error) {
	res, err := r.Component.ListComponents(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return &model.Components{
		Components:       convert.ComponentsToModel(res.Components),
		CustomComponents: convert.ComponentsToModel(res.CustomComponents),
	}, nil
}
