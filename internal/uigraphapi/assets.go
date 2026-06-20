package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type AssetUpload struct {
	AssetID   string `json:"assetId"`
	UploadURL string `json:"uploadUrl"`
}

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

// CreateAssetUpload allocates a new asset id and returns a presigned PUT URL.
func (c *Client) CreateAssetUpload(ctx context.Context, orgID string) (*AssetUpload, error) {
	var out AssetUpload
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/assets", orgID), map[string]interface{}{}, &out)
}
