package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type SavedQueryFolder struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"orgId"`
	ServiceDBID string    `json:"serviceDbId"`
	Scope       string    `json:"scope"`
	OwnerUserID *string   `json:"ownerUserId,omitempty"`
	TeamID      *string   `json:"teamId,omitempty"`
	Name        string    `json:"name"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type SavedQuery struct {
	ID                  string    `json:"id"`
	OrgID               string    `json:"orgId"`
	ServiceDBID         string    `json:"serviceDbId"`
	FolderID            *string   `json:"folderId,omitempty"`
	Scope               string    `json:"scope"`
	OwnerUserID         *string   `json:"ownerUserId,omitempty"`
	TeamID              *string   `json:"teamId,omitempty"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	QueryText           string    `json:"queryText"`
	Tags                []string  `json:"tags"`
	Source              *string   `json:"source,omitempty"`
	CreatedBy           string    `json:"createdBy"`
	UpdatedBy           *string   `json:"updatedBy,omitempty"`
	CreatedByCommitHash *string   `json:"createdByCommitHash,omitempty"`
	UpdatedByCommitHash *string   `json:"updatedByCommitHash,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

func savedQueriesBasePath(orgID, serviceID, dbID string) string {
	return fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, dbID)
}

func (c *Client) ListSavedQueryFolders(ctx context.Context, orgID, serviceID, dbID, scope string) ([]SavedQueryFolder, error) {
	var out struct {
		Folders []SavedQueryFolder `json:"folders"`
	}
	path := fmt.Sprintf("%s/query-folders?scope=%s", savedQueriesBasePath(orgID, serviceID, dbID), scope)
	return out.Folders, c.get(ctx, path, &out)
}

func (c *Client) CreateSavedQueryFolder(ctx context.Context, orgID, serviceID, dbID string, body map[string]interface{}) (*SavedQueryFolder, error) {
	var out SavedQueryFolder
	return &out, c.post(ctx, savedQueriesBasePath(orgID, serviceID, dbID)+"/query-folders", body, &out)
}

func (c *Client) DeleteSavedQueryFolder(ctx context.Context, orgID, serviceID, dbID, id string) error {
	return c.del(ctx, fmt.Sprintf("%s/query-folders/%s", savedQueriesBasePath(orgID, serviceID, dbID), id))
}

func (c *Client) ListSavedQueries(ctx context.Context, orgID, serviceID, dbID, scope string) ([]SavedQuery, error) {
	var out struct {
		Queries []SavedQuery `json:"queries"`
	}
	path := fmt.Sprintf("%s/queries?scope=%s", savedQueriesBasePath(orgID, serviceID, dbID), scope)
	return out.Queries, c.get(ctx, path, &out)
}

func (c *Client) CreateSavedQuery(ctx context.Context, orgID, serviceID, dbID string, body map[string]interface{}) (*SavedQuery, error) {
	var out SavedQuery
	return &out, c.post(ctx, savedQueriesBasePath(orgID, serviceID, dbID)+"/queries", body, &out)
}

func (c *Client) UpdateSavedQuery(ctx context.Context, orgID, serviceID, dbID, id string, body map[string]interface{}) (*SavedQuery, error) {
	var out SavedQuery
	return &out, c.put(ctx, fmt.Sprintf("%s/queries/%s", savedQueriesBasePath(orgID, serviceID, dbID), id), body, &out)
}

func (c *Client) DeleteSavedQuery(ctx context.Context, orgID, serviceID, dbID, id string) error {
	return c.del(ctx, fmt.Sprintf("%s/queries/%s", savedQueriesBasePath(orgID, serviceID, dbID), id))
}
