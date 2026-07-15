package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func (r *queryResolver) CostSavingsSummary(ctx context.Context, orgID string, period *string, modelID *string) (*model.SavingsSummary, error) {
	s, err := r.CostSavings.GetSavingsSummary(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}
	return convert.SavingsSummaryToModel(s), nil
}

func (r *queryResolver) CostSavingsTimeseries(ctx context.Context, orgID string, period *string, modelID *string) ([]*model.DailySavings, error) {
	rows, err := r.CostSavings.GetSavingsTimeseries(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}
	return convert.DailySavingsListToModel(rows), nil
}

func (r *queryResolver) CostSavingsByTool(ctx context.Context, orgID string, period *string, modelID *string) ([]*model.ToolSavings, error) {
	rows, err := r.CostSavings.GetSavingsByTool(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}
	return convert.ToolSavingsListToModel(rows), nil
}

func (r *queryResolver) CostSavingsByClient(ctx context.Context, orgID string, period *string, modelID *string) ([]*model.ClientSavings, error) {
	rows, err := r.CostSavings.GetSavingsByClient(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}
	return convert.ClientSavingsListToModel(rows), nil
}

func (r *queryResolver) CostSavingsByModel(ctx context.Context, orgID string, period *string) ([]*model.ModelSavings, error) {
	rows, err := r.CostSavings.GetSavingsByModel(ctx, orgID, period)
	if err != nil {
		return nil, err
	}
	return convert.ModelSavingsListToModel(rows), nil
}

func (r *queryResolver) CostSavingsByUser(ctx context.Context, orgID string, period *string, modelID *string) ([]*model.UserSavings, error) {
	rows, err := r.CostSavings.GetSavingsByUser(ctx, orgID, period, modelID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.UserID != nil {
			ids = append(ids, *row.UserID)
		} else if row.ServiceAccountID != nil {
			ids = append(ids, *row.ServiceAccountID)
		}
	}

	actors := map[string]*uigraphapi.Actor{}
	if len(ids) > 0 {
		var actorErr error
		actors, actorErr = r.Resolver.Actor.ResolveActors(ctx, orgID, ids)
		if actorErr != nil {
			return nil, actorErr
		}
	}
	return convert.UserSavingsListToModel(rows, actors), nil
}
