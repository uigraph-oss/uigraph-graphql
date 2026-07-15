package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *queryResolver) Actor(ctx context.Context, orgID string, id string) (*model.Actor, error) {
	return r.resolveActor(ctx, orgID, id)
}
