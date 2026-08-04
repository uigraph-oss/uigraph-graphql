package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type CloudConnection struct {
	ID             string     `json:"id"`
	OrgID          string     `json:"orgId"`
	Provider       string     `json:"provider"`
	DisplayName    string     `json:"displayName"`
	Status         string     `json:"status"`
	StatusMessage  string     `json:"statusMessage"`
	LastVerifiedAt *time.Time `json:"lastVerifiedAt,omitempty"`
	CreatedBy      string     `json:"createdBy"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type ServiceCostTagRule struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"serviceId"`
	TagKey    string    `json:"tagKey"`
	TagValue  string    `json:"tagValue"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type InfraResource struct {
	ID                 string    `json:"id"`
	ExternalResourceID string    `json:"externalResourceId"`
	Name               string    `json:"name"`
	ResourceType       string    `json:"resourceType"`
	Provider           string    `json:"provider"`
	Region             string    `json:"region"`
	Environment        string    `json:"environment"`
	Status             string    `json:"status"`
	MonthlyCostUSD     float64   `json:"monthlyCostUsd"`
	Tags               []string  `json:"tags"`
	MatchedTags        []string  `json:"matchedTags"`
	LastSyncedAt       time.Time `json:"lastSyncedAt"`
}

type CostTrendPoint struct {
	Date     string  `json:"date"`
	TotalUSD float64 `json:"totalUsd"`
	AWSUSD   float64 `json:"awsUsd"`
	AzureUSD float64 `json:"azureUsd"`
	GCPUSD   float64 `json:"gcpUsd"`
}

type ServiceCostSummary struct {
	TotalMonthlyCostUSD float64 `json:"totalMonthlyCostUsd"`
	MoMChangePct        float64 `json:"momChangePct"`
	ResourceCount       int     `json:"resourceCount"`
	ProviderCount       int     `json:"providerCount"`
	TopCostDriverLabel  string  `json:"topCostDriverLabel"`
	TopCostDriverUSD    float64 `json:"topCostDriverUsd"`
}

type TestCloudConnectionResult struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func (c *Client) ListCloudConnections(ctx context.Context, orgID string) ([]CloudConnection, error) {
	var out struct {
		Connections []CloudConnection `json:"connections"`
	}
	return out.Connections, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/billing/connections", orgID), &out)
}

func (c *Client) GetCloudConnection(ctx context.Context, orgID, connectionID string) (*CloudConnection, error) {
	var out CloudConnection
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/billing/connections/%s", orgID, connectionID), &out)
}

func (c *Client) CreateCloudConnection(ctx context.Context, orgID string, body map[string]interface{}) (*CloudConnection, error) {
	var out CloudConnection
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/billing/connections", orgID), body, &out)
}

func (c *Client) DeleteCloudConnection(ctx context.Context, orgID, connectionID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/billing/connections/%s", orgID, connectionID))
}

func (c *Client) TestCloudConnection(ctx context.Context, orgID, connectionID string) (*TestCloudConnectionResult, error) {
	var out TestCloudConnectionResult
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/billing/connections/%s/test", orgID, connectionID), nil, &out)
}

// SyncCloudConnection kicks off an async sync; the REST endpoint returns as
// soon as the sync is started, not once it finishes.
func (c *Client) SyncCloudConnection(ctx context.Context, orgID, connectionID string) (bool, error) {
	var out struct {
		Started bool `json:"started"`
	}
	return out.Started, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/billing/connections/%s/sync", orgID, connectionID), nil, &out)
}

func (c *Client) GetServiceCostSummary(ctx context.Context, orgID, serviceID string) (*ServiceCostSummary, error) {
	var out ServiceCostSummary
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/costs/summary", orgID, serviceID), &out)
}

func (c *Client) ListServiceCostResources(ctx context.Context, orgID, serviceID string) ([]InfraResource, error) {
	var out struct {
		Resources []InfraResource `json:"resources"`
	}
	return out.Resources, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/costs/resources", orgID, serviceID), &out)
}

func (c *Client) GetServiceCostTrend(ctx context.Context, orgID, serviceID string, days *int) ([]CostTrendPoint, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/costs/trend", orgID, serviceID)
	if days != nil {
		path = fmt.Sprintf("%s?days=%d", path, *days)
	}
	var out struct {
		Trend []CostTrendPoint `json:"trend"`
	}
	return out.Trend, c.get(ctx, path, &out)
}

func (c *Client) ListServiceCostTagRules(ctx context.Context, orgID, serviceID string) ([]ServiceCostTagRule, error) {
	var out struct {
		TagRules []ServiceCostTagRule `json:"tagRules"`
	}
	return out.TagRules, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/costs/tag-rules", orgID, serviceID), &out)
}

func (c *Client) CreateServiceCostTagRule(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceCostTagRule, error) {
	var out ServiceCostTagRule
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/costs/tag-rules", orgID, serviceID), body, &out)
}

func (c *Client) DeleteServiceCostTagRule(ctx context.Context, orgID, serviceID, ruleID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/costs/tag-rules/%s", orgID, serviceID, ruleID))
}
