package graph

import (
	"context"
	"strings"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *mutationResolver) CreateSavedQueryFolder(ctx context.Context, orgID string, serviceID string, serviceDbID string, input model.CreateSavedQueryFolderInput) (*model.SavedQueryFolder, error) {
	body := convert.ToMap(input)
	body["scope"] = strings.ToLower(input.Scope.String())
	f, err := r.Catalog.CreateSavedQueryFolder(ctx, orgID, serviceID, serviceDbID, body)
	if err != nil {
		return nil, err
	}
	return convert.SavedQueryFolderToModel(f), nil
}

func (r *mutationResolver) DeleteSavedQueryFolder(ctx context.Context, orgID string, serviceID string, serviceDbID string, id string) (bool, error) {
	return true, r.Catalog.DeleteSavedQueryFolder(ctx, orgID, serviceID, serviceDbID, id)
}

func (r *mutationResolver) CreateSavedQuery(ctx context.Context, orgID string, serviceID string, serviceDbID string, input model.CreateSavedQueryInput) (*model.SavedQuery, error) {
	body := convert.ToMap(input)
	body["scope"] = strings.ToLower(input.Scope.String())
	q, err := r.Catalog.CreateSavedQuery(ctx, orgID, serviceID, serviceDbID, body)
	if err != nil {
		return nil, err
	}
	return convert.SavedQueryToModel(q), nil
}

func (r *mutationResolver) UpdateSavedQuery(ctx context.Context, orgID string, serviceID string, serviceDbID string, id string, input model.UpdateSavedQueryInput) (*model.SavedQuery, error) {
	q, err := r.Catalog.UpdateSavedQuery(ctx, orgID, serviceID, serviceDbID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.SavedQueryToModel(q), nil
}

func (r *mutationResolver) DeleteSavedQuery(ctx context.Context, orgID string, serviceID string, serviceDbID string, id string) (bool, error) {
	return true, r.Catalog.DeleteSavedQuery(ctx, orgID, serviceID, serviceDbID, id)
}

func (r *queryResolver) SavedQueryFolders(ctx context.Context, orgID string, serviceID string, serviceDbID string, scope model.SavedQueryScope) ([]*model.SavedQueryFolder, error) {
	folders, err := r.Catalog.ListSavedQueryFolders(ctx, orgID, serviceID, serviceDbID, strings.ToLower(scope.String()))
	if err != nil {
		return nil, err
	}
	return convert.SavedQueryFoldersToModel(folders), nil
}

func (r *queryResolver) SavedQueries(ctx context.Context, orgID string, serviceID string, serviceDbID string, scope model.SavedQueryScope) ([]*model.SavedQuery, error) {
	queries, err := r.Catalog.ListSavedQueries(ctx, orgID, serviceID, serviceDbID, strings.ToLower(scope.String()))
	if err != nil {
		return nil, err
	}
	return convert.SavedQueriesToModel(queries), nil
}

func (r *savedQueryResolver) CreatedByActor(ctx context.Context, obj *model.SavedQuery) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *savedQueryResolver) UpdatedByActor(ctx context.Context, obj *model.SavedQuery) (*model.Actor, error) {
	if obj.UpdatedBy == nil {
		return nil, nil
	}
	return r.resolveActor(ctx, obj.OrgID, *obj.UpdatedBy)
}

func (r *Resolver) SavedQuery() generated.SavedQueryResolver { return &savedQueryResolver{r} }

type savedQueryResolver struct{ *Resolver }
