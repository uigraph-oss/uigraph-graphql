package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type Folder struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"orgId"`
	ParentID  *string    `json:"parentId,omitempty"`
	TeamID    *string    `json:"teamId,omitempty"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Order     float64    `json:"order"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (c *Client) ListFolders(ctx context.Context, orgID, folderType, parentID string) ([]Folder, error) {
	q := url.Values{}
	if folderType != "" {
		q.Set("type", folderType)
	}
	if parentID != "" {
		q.Set("parentId", parentID)
	}
	path := "/api/v1/orgs/" + orgID + "/folders"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Folders []Folder `json:"folders"`
	}
	return out.Folders, c.get(ctx, path, &out)
}

func (c *Client) GetFolder(ctx context.Context, orgID, id string) (*Folder, error) {
	var out Folder
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id), &out)
}

func (c *Client) CreateFolder(ctx context.Context, orgID string, body map[string]interface{}) (*Folder, error) {
	var out Folder
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/folders", body, &out)
}

func (c *Client) UpdateFolder(ctx context.Context, orgID, id string, body map[string]interface{}) (*Folder, error) {
	var out Folder
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id), body, &out)
}

func (c *Client) DeleteFolder(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id))
}
