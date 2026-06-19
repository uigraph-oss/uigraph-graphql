package uigraphapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Service struct {
	ID              string          `json:"id"`
	OrgID           string          `json:"orgId"`
	FolderID        *string         `json:"folderId,omitempty"`
	TeamID          *string         `json:"teamId,omitempty"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	Description     string          `json:"description"`
	Status          string          `json:"status"`
	Tier            string          `json:"tier"`
	Category        string          `json:"category"`
	Language        string          `json:"language"`
	GitRepoURL      *string         `json:"gitRepoUrl,omitempty"`
	JiraProjectURL  *string         `json:"jiraProjectUrl,omitempty"`
	SlackChannelURL *string         `json:"slackChannelUrl,omitempty"`
	LastCommitSha   *string         `json:"lastCommitSha,omitempty"`
	Labels          []string        `json:"labels"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CreatedBy       string          `json:"createdBy"`
	UpdatedBy       *string         `json:"updatedBy,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type ServiceStats struct {
	ServiceID     string `json:"serviceId"`
	EndpointCount int    `json:"endpointCount"`
	DiagramCount  int    `json:"diagramCount"`
	DocCount      int    `json:"docCount"`
	DBTableCount  int    `json:"dbTableCount"`
	TestCaseCount int    `json:"testCaseCount"`
}

type APIGroup struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"serviceId"`
	OrgID     string    `json:"orgId"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Label     *string   `json:"label,omitempty"`
	Protocol  string    `json:"protocol"`
	SpecKey   *string   `json:"specKey,omitempty"`
	SpecHash  *string   `json:"specHash,omitempty"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy *string   `json:"updatedBy,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type APIGroupVersion struct {
	ID            string    `json:"id"`
	APIGroupID    string    `json:"apiGroupId"`
	VersionNumber int       `json:"versionNumber"`
	Label         *string   `json:"label,omitempty"`
	SpecKey       string    `json:"specKey"`
	SpecHash      string    `json:"specHash"`
	IsAutoVersion bool      `json:"isAutoVersion"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

type ServiceDoc struct {
	ID          string    `json:"id"`
	ServiceID   string    `json:"serviceId"`
	OrgID       string    `json:"orgId"`
	FileKey     string    `json:"fileKey"`
	FileName    string    `json:"fileName"`
	FileType    string    `json:"fileType"`
	Description string    `json:"description"`
	ContentHash string    `json:"contentHash"`
	CreatedBy   string    `json:"createdBy"`
	UpdatedBy   *string   `json:"updatedBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ServiceDiagram struct {
	ServiceID string    `json:"serviceId"`
	DiagramID string    `json:"diagramId"`
	OrgID     string    `json:"orgId"`
	CreatedBy string    `json:"createdBy"`
	UpdatedBy *string   `json:"updatedBy,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Diagram   *Diagram  `json:"diagram,omitempty"`
}

type ServiceDB struct {
	ID         string          `json:"id"`
	ServiceID  string          `json:"serviceId"`
	OrgID      string          `json:"orgId"`
	DBName     string          `json:"dbName"`
	DBType     string          `json:"dbType"`
	Dialect    string          `json:"dialect"`
	SchemaJSON json.RawMessage `json:"schemaJson"`
	Source     *string         `json:"source,omitempty"`
	SourceTS   *time.Time      `json:"sourceTs,omitempty"`
	CreatedBy  string          `json:"createdBy"`
	UpdatedBy  *string         `json:"updatedBy,omitempty"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type ServiceDBVersion struct {
	ID            string          `json:"id"`
	ServiceDBID   string          `json:"serviceDbId"`
	VersionNumber int             `json:"versionNumber"`
	Label         *string         `json:"label,omitempty"`
	SchemaJSON    json.RawMessage `json:"schemaJson"`
	Source        *string         `json:"source,omitempty"`
	SourceTS      *time.Time      `json:"sourceTs,omitempty"`
	IsAutoVersion bool            `json:"isAutoVersion"`
	CreatedBy     string          `json:"createdBy"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type APIEndpoint struct {
	ID          string          `json:"id"`
	APIGroupID  string          `json:"apiGroupId"`
	ServiceID   string          `json:"serviceId"`
	OrgID       string          `json:"orgId"`
	OperationID string          `json:"operationId"`
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	Summary     string          `json:"summary"`
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	Parameters  json.RawMessage `json:"parameters"`
	RequestBody json.RawMessage `json:"requestBody"`
	Responses   json.RawMessage `json:"responses"`
	Order       float64         `json:"order"`
	CreatedBy   string          `json:"createdBy"`
	UpdatedBy   *string         `json:"updatedBy,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func (c *Client) ListServices(ctx context.Context, orgID, folderID, teamID string) ([]Service, error) {
	path := "/api/v1/orgs/" + orgID + "/services"
	q := url.Values{}
	if folderID != "" {
		q.Set("folderId", folderID)
	}
	if teamID != "" {
		q.Set("teamId", teamID)
	}
	if enc := q.Encode(); enc != "" {
		path += "?" + enc
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

func (c *Client) ListServiceStats(ctx context.Context, orgID string, serviceID *string) ([]ServiceStats, error) {
	path := "/api/v1/orgs/" + orgID + "/services/stats"
	if serviceID != nil && *serviceID != "" {
		q := url.Values{}
		q.Set("serviceId", *serviceID)
		path += "?" + q.Encode()
	}
	var out struct {
		Stats []ServiceStats `json:"stats"`
	}
	return out.Stats, c.get(ctx, path, &out)
}

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

func (c *Client) ListServiceDocs(ctx context.Context, orgID, serviceID string) ([]ServiceDoc, error) {
	var out struct {
		Docs []ServiceDoc `json:"docs"`
	}
	return out.Docs, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs", orgID, serviceID), &out)
}

func (c *Client) GetServiceDoc(ctx context.Context, orgID, serviceID, id string) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateServiceDoc(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs", orgID, serviceID), body, &out)
}

func (c *Client) UpdateServiceDoc(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteServiceDoc(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id))
}

func (c *Client) ListServiceDiagrams(ctx context.Context, orgID, serviceID string) ([]ServiceDiagram, error) {
	var out struct {
		Diagrams []ServiceDiagram `json:"diagrams"`
	}
	return out.Diagrams, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams", orgID, serviceID), &out)
}

func (c *Client) CreateServiceDiagram(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDiagram, error) {
	var out ServiceDiagram
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams", orgID, serviceID), body, &out)
}

func (c *Client) DeleteServiceDiagram(ctx context.Context, orgID, serviceID, diagramID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams/%s", orgID, serviceID, diagramID))
}

func (c *Client) ListServiceDBs(ctx context.Context, orgID, serviceID string) ([]ServiceDB, error) {
	var out struct {
		DBs []ServiceDB `json:"dbs"`
	}
	return out.DBs, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs", orgID, serviceID), &out)
}

func (c *Client) GetServiceDB(ctx context.Context, orgID, serviceID, id string) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateServiceDB(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs", orgID, serviceID), body, &out)
}

func (c *Client) UpdateServiceDB(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteServiceDB(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id))
}

func (c *Client) ListServiceDBVersions(ctx context.Context, orgID, serviceID, serviceDBID string) ([]ServiceDBVersion, error) {
	var out struct {
		Versions []ServiceDBVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions", orgID, serviceID, serviceDBID), &out)
}

func (c *Client) CreateServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID string, body map[string]interface{}) (*ServiceDBVersion, error) {
	var out ServiceDBVersion
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions", orgID, serviceID, serviceDBID), body, &out)
}

func (c *Client) RestoreServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID, versionID string) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions/%s/restore", orgID, serviceID, serviceDBID, versionID), nil, &out)
}

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
