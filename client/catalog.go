package client

import (
	"context"
	"fmt"
)

// ── Services ──────────────────────────────────────────────────────────────────

func (c *Client) ListServices(ctx context.Context, orgID, folderID string) ([]Service, error) {
	path := "/api/v1/orgs/" + orgID + "/services"
	if folderID != "" {
		path += "?folderId=" + folderID
	}
	var out struct {
		Services []Service `json:"services"`
	}
	return out.Services, c.get(ctx, path, &out)
}

func (c *Client) GetService(ctx context.Context, orgID, id string) (*Service, error) {
	var out Service
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id), &out)
}

func (c *Client) CreateService(ctx context.Context, orgID string, body map[string]interface{}) (*Service, error) {
	var out Service
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/services", body, &out)
}

func (c *Client) UpdateService(ctx context.Context, orgID, id string, body map[string]interface{}) (*Service, error) {
	var out Service
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id), body, &out)
}

func (c *Client) DeleteService(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id))
}

// ── API Groups ────────────────────────────────────────────────────────────────

func (c *Client) ListAPIGroups(ctx context.Context, orgID, serviceID string) ([]APIGroup, error) {
	var out struct {
		APIGroups []APIGroup `json:"apiGroups"`
	}
	return out.APIGroups, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups", orgID, serviceID), &out)
}

func (c *Client) GetAPIGroup(ctx context.Context, orgID, serviceID, id string) (*APIGroup, error) {
	var out APIGroup
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*APIGroup, error) {
	var out APIGroup
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups", orgID, serviceID), body, &out)
}

func (c *Client) UpdateAPIGroup(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*APIGroup, error) {
	var out APIGroup
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteAPIGroup(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id))
}

func (c *Client) SyncAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/sync", orgID, serviceID), body, &out)
}

func (c *Client) ListAPIGroupVersions(ctx context.Context, orgID, serviceID, apiGroupID string) ([]APIGroupVersion, error) {
	var out struct {
		Versions []APIGroupVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/versions", orgID, serviceID, apiGroupID), &out)
}

// ── API Endpoints ─────────────────────────────────────────────────────────────

func (c *Client) ListAPIEndpoints(ctx context.Context, orgID, serviceID, apiGroupID string) ([]APIEndpoint, error) {
	var out struct {
		Endpoints []APIEndpoint `json:"endpoints"`
	}
	return out.Endpoints, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints", orgID, serviceID, apiGroupID), &out)
}

func (c *Client) GetAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id), &out)
}

func (c *Client) CreateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID string, body map[string]interface{}) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints", orgID, serviceID, apiGroupID), body, &out)
}

func (c *Client) UpdateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string, body map[string]interface{}) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id), body, &out)
}

func (c *Client) DeleteAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id))
}
