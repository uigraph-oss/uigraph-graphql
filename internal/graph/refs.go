package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/model"
)

// resolveActor resolves a single created_by / updated_by id within an org to
// its public actor info, returning nil when id is empty or matches no actor.
func (r *Resolver) resolveActor(ctx context.Context, orgID, id string) (*model.Actor, error) {
	if id == "" {
		return nil, nil
	}
	actors, err := r.Actor.ResolveActors(ctx, orgID, []string{id})
	if err != nil {
		return nil, err
	}
	a := actors[id]
	if a == nil {
		return nil, nil
	}
	m := &model.Actor{ID: a.ID, Type: a.Type, Name: a.Name, Disabled: a.Disabled}
	if a.Email != "" {
		m.Email = &a.Email
	}
	if a.AvatarURL != "" {
		m.AvatarURL = &a.AvatarURL
	}
	return m, nil
}

// resolveAssetURL resolves a single asset id within an org to a presigned GET
// URL, returning nil when id is empty or no url is produced.
func (r *Resolver) resolveAssetURL(ctx context.Context, orgID, assetID string) (*string, error) {
	if assetID == "" {
		return nil, nil
	}
	urls, err := r.Actor.ResolveAssetURLs(ctx, orgID, []string{assetID})
	if err != nil {
		return nil, err
	}
	u, ok := urls[assetID]
	if !ok || u == "" {
		return nil, nil
	}
	return &u, nil
}
