package client

import (
	"context"
	"net/url"
	"strings"
)

// ResolveAssetURLs resolves the given asset ids within an org, returning a map
// from id to its presigned GET URL.
func (c *Client) ResolveAssetURLs(ctx context.Context, orgID string, ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return map[string]string{}, nil
	}
	var out struct {
		URLs map[string]string `json:"urls"`
	}
	path := "/api/v1/orgs/" + orgID + "/assets/urls?ids=" + url.QueryEscape(strings.Join(ids, ","))
	return out.URLs, c.get(ctx, path, &out)
}
