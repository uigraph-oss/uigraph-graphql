package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type Comment struct {
	ID              string    `json:"id"`
	OrgID           string    `json:"orgId"`
	ResourceID      string    `json:"resourceId"`
	ParentCommentID *string   `json:"parentCommentId,omitempty"`
	Text            string    `json:"text"`
	CreatedBy       string    `json:"createdBy"`
	UpdatedBy       *string   `json:"updatedBy,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (c *Client) ListComments(ctx context.Context, orgID, resourceID string) ([]Comment, error) {
	var out struct {
		Comments []Comment `json:"comments"`
	}
	path := fmt.Sprintf("/api/v1/orgs/%s/comments?resourceId=%s", orgID, url.QueryEscape(resourceID))
	return out.Comments, c.get(ctx, path, &out)
}

func (c *Client) CreateComment(ctx context.Context, orgID string, body map[string]interface{}) (*Comment, error) {
	var out Comment
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/comments", orgID), body, &out)
}

func (c *Client) UpdateComment(ctx context.Context, orgID, id string, body map[string]interface{}) (*Comment, error) {
	var out Comment
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/comments/%s", orgID, id), body, &out)
}

func (c *Client) DeleteComment(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/comments/%s", orgID, id))
}
