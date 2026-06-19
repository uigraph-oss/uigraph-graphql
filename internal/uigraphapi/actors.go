package uigraphapi

import (
	"context"
	"net/url"
	"strings"
)

// Actor is public identity info for a created_by / updated_by / deleted_by id,
// resolved by uigraph-api to either a user or a service account.
type Actor struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	Disabled  bool   `json:"disabled"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// ResolveActors resolves the given ids within an org, returning a map from id
// to its actor (nil for ids that match no user or service account).
func (c *Client) ResolveActors(ctx context.Context, orgID string, ids []string) (map[string]*Actor, error) {
	if len(ids) == 0 {
		return map[string]*Actor{}, nil
	}
	var out struct {
		Actors map[string]*Actor `json:"actors"`
	}
	path := "/api/v1/orgs/" + orgID + "/actors?ids=" + url.QueryEscape(strings.Join(ids, ","))
	return out.Actors, c.get(ctx, path, &out)
}
