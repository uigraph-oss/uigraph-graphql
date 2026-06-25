package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type Doc struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"orgId"`
	FolderID    *string   `json:"folderId,omitempty"`
	TeamID      *string   `json:"teamId,omitempty"`
	FileAssetID string    `json:"fileAssetId"`
	FileName    string    `json:"fileName"`
	FileType    string    `json:"fileType"`
	Description string    `json:"description"`
	ContentHash string    `json:"contentHash"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedBy   *string   `json:"updatedBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (c *Client) ListDocs(ctx context.Context, orgID string, p ListParams) ([]Doc, int, error) {
	path := "/api/v1/orgs/" + orgID + "/docs" + listQuery(p)
	var out struct {
		Docs  []Doc `json:"docs"`
		Total int   `json:"total"`
	}
	return out.Docs, out.Total, c.get(ctx, path, &out)
}

func (c *Client) GetDoc(ctx context.Context, orgID, id string) (*Doc, error) {
	var out Doc
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/docs/%s", orgID, id), &out)
}

func (c *Client) CreateDoc(ctx context.Context, orgID string, body map[string]interface{}) (*Doc, error) {
	var out Doc
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/docs", body, &out)
}

func (c *Client) UpdateDoc(ctx context.Context, orgID, id string, body map[string]interface{}) (*Doc, error) {
	var out Doc
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/docs/%s", orgID, id), body, &out)
}

func (c *Client) DeleteDoc(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/docs/%s", orgID, id))
}
