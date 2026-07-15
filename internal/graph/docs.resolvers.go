package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func (r *docResolver) FileURL(ctx context.Context, obj *model.Doc) (*string, error) {
	return r.resolveAssetURL(ctx, obj.OrgID, obj.FileAssetID)
}

func (r *docResolver) CreatedByActor(ctx context.Context, obj *model.Doc) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *docResolver) UpdatedByActor(ctx context.Context, obj *model.Doc) (*model.Actor, error) {
	if obj.UpdatedBy == nil {
		return nil, nil
	}
	return r.resolveActor(ctx, obj.OrgID, *obj.UpdatedBy)
}

func (r *mutationResolver) CreateDoc(ctx context.Context, orgID string, input model.CreateDocInput) (*model.Doc, error) {
	d, err := r.DocAPI.CreateDoc(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.DocToModel(d), nil
}

func (r *mutationResolver) UpdateDoc(ctx context.Context, orgID string, id string, input model.UpdateDocInput) (*model.Doc, error) {
	d, err := r.DocAPI.UpdateDoc(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.DocToModel(d), nil
}

func (r *mutationResolver) DeleteDoc(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.DocAPI.DeleteDoc(ctx, orgID, id)
}

func (r *queryResolver) Docs(ctx context.Context, orgID string, folderID *string, teamID *string, search *string, sortBy *string, sortDir *string, limit *int, offset *int) (*model.DocPage, error) {
	p := uigraphapi.ListParams{
		FolderID: derefStr(folderID),
		TeamID:   derefStr(teamID),
		Search:   derefStr(search),
		SortBy:   derefStr(sortBy),
		SortDir:  derefStr(sortDir),
		Limit:    limit,
		Offset:   offset,
	}
	docs, total, err := r.DocAPI.ListDocs(ctx, orgID, p)
	if err != nil {
		return nil, err
	}
	return &model.DocPage{Items: convert.DocsToModel(docs), TotalCount: total}, nil
}

func (r *queryResolver) Doc(ctx context.Context, orgID string, id string) (*model.Doc, error) {
	d, err := r.DocAPI.GetDoc(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.DocToModel(d), nil
}

func (r *Resolver) Doc() generated.DocResolver { return &docResolver{r} }

type docResolver struct{ *Resolver }
