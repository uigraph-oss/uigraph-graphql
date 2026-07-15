package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/model"
)

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
