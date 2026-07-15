package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *mutationResolver) CreateAssetUpload(ctx context.Context, orgID string) (*model.AssetUpload, error) {
	u, err := r.Actor.CreateAssetUpload(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return &model.AssetUpload{AssetID: u.AssetID, UploadURL: u.UploadURL}, nil
}

func (r *queryResolver) AssetURL(ctx context.Context, orgID string, assetID string) (*string, error) {
	return r.resolveAssetURL(ctx, orgID, assetID)
}

func (r *queryResolver) AssetUrls(ctx context.Context, orgID string, assetIds []string) ([]*model.AssetURL, error) {
	urlsByID, err := r.Resolver.Actor.ResolveAssetURLs(ctx, orgID, assetIds)
	if err != nil {
		return nil, err
	}

	out := make([]*model.AssetURL, 0, len(assetIds))
	for _, id := range assetIds {
		if u, ok := urlsByID[id]; ok && u != "" {
			out = append(out, &model.AssetURL{AssetID: id, URL: u})
		}
	}
	return out, nil
}
