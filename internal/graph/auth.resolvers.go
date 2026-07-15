package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *mutationResolver) SwitchOrg(ctx context.Context, orgID string) (bool, error) {
	return true, r.Auth.SwitchOrg(ctx, orgID)
}

func (r *mutationResolver) PrepareUserAvatarUpload(ctx context.Context) (*model.AssetUpload, error) {
	u, err := r.Auth.PrepareUserAvatarUpload(ctx)
	if err != nil {
		return nil, err
	}
	return &model.AssetUpload{AssetID: u.AssetID, UploadURL: u.UploadURL}, nil
}

func (r *mutationResolver) SetMyAvatar(ctx context.Context) (bool, error) {
	return true, r.Auth.SetMyAvatar(ctx)
}

func (r *queryResolver) Me(ctx context.Context) (*model.Me, error) {
	me, err := r.Auth.Me(ctx)
	if err != nil {
		return nil, err
	}
	return convert.MeToModel(me), nil
}

func (r *queryResolver) MyOrgs(ctx context.Context) ([]*model.OrgSummary, error) {
	orgs, err := r.Auth.MyOrgs(ctx)
	if err != nil {
		return nil, err
	}
	return convert.OrgSummariesToModel(orgs), nil
}
