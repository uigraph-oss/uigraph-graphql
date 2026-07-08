package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type SavingsSummary struct {
	OrgID             string  `json:"orgId"`
	Period            string  `json:"period"`
	ModelID           string  `json:"modelId"`
	TotalCalls        int     `json:"totalCalls"`
	TotalTokensServed int     `json:"totalTokensServed"`
	TotalTokensSaved  int     `json:"totalTokensSaved"`
	CostServedUSD     float64 `json:"costServedUsd"`
	CostRawUSD        float64 `json:"costRawUsd"`
	CostSavedUSD      float64 `json:"costSavedUsd"`
	UniqueUsersCount  int     `json:"uniqueUsersCount"`
}

type DailySavings struct {
	Date              time.Time `json:"date"`
	TotalCalls        int       `json:"totalCalls"`
	TotalTokensServed int       `json:"totalTokensServed"`
	TotalTokensSaved  int       `json:"totalTokensSaved"`
	CostServedUSD     float64   `json:"costServedUsd"`
	CostRawUSD        float64   `json:"costRawUsd"`
	CostSavedUSD      float64   `json:"costSavedUsd"`
}

type ToolSavings struct {
	ToolName     string  `json:"toolName"`
	TotalCalls   int     `json:"totalCalls"`
	TokensSaved  int     `json:"tokensSaved"`
	CostSavedUSD float64 `json:"costSavedUsd"`
}

type ClientSavings struct {
	ClientName   string  `json:"clientName"`
	TotalCalls   int     `json:"totalCalls"`
	TokensSaved  int     `json:"tokensSaved"`
	CostSavedUSD float64 `json:"costSavedUsd"`
}

type ModelSavings struct {
	ModelID      string  `json:"modelId"`
	DisplayName  string  `json:"displayName"`
	Provider     string  `json:"provider"`
	TotalCalls   int     `json:"totalCalls"`
	TokensSaved  int     `json:"tokensSaved"`
	CostRawUSD   float64 `json:"costRawUsd"`
	CostSavedUSD float64 `json:"costSavedUsd"`
}

type UserSavings struct {
	UserID           *string `json:"userId,omitempty"`
	ServiceAccountID *string `json:"serviceAccountId,omitempty"`
	TotalCalls       int     `json:"totalCalls"`
	TokensSaved      int     `json:"tokensSaved"`
	CostSavedUSD     float64 `json:"costSavedUsd"`
}

func savingsQuery(period, modelID *string) url.Values {
	q := url.Values{}
	if period != nil && *period != "" {
		q.Set("period", *period)
	}
	if modelID != nil && *modelID != "" {
		q.Set("model_id", *modelID)
	}
	return q
}

func withQuery(path string, q url.Values) string {
	if len(q) > 0 {
		return path + "?" + q.Encode()
	}
	return path
}

func (c *Client) GetSavingsSummary(ctx context.Context, orgID string, period, modelID *string) (*SavingsSummary, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/summary", orgID), savingsQuery(period, modelID))
	var out SavingsSummary
	return &out, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsTimeseries(ctx context.Context, orgID string, period, modelID *string) ([]DailySavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/timeseries", orgID), savingsQuery(period, modelID))
	var out struct {
		Timeseries []DailySavings `json:"timeseries"`
	}
	return out.Timeseries, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsByTool(ctx context.Context, orgID string, period, modelID *string) ([]ToolSavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/by-tool", orgID), savingsQuery(period, modelID))
	var out struct {
		ByTool []ToolSavings `json:"byTool"`
	}
	return out.ByTool, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsByClient(ctx context.Context, orgID string, period, modelID *string) ([]ClientSavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/by-client", orgID), savingsQuery(period, modelID))
	var out struct {
		ByClient []ClientSavings `json:"byClient"`
	}
	return out.ByClient, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsByModel(ctx context.Context, orgID string, period *string) ([]ModelSavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/by-model", orgID), savingsQuery(period, nil))
	var out struct {
		ByModel []ModelSavings `json:"byModel"`
	}
	return out.ByModel, c.get(ctx, path, &out)
}

func (c *Client) GetSavingsByUser(ctx context.Context, orgID string, period, modelID *string) ([]UserSavings, error) {
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/mcp/savings/by-user", orgID), savingsQuery(period, modelID))
	var out struct {
		ByUser []UserSavings `json:"byUser"`
	}
	return out.ByUser, c.get(ctx, path, &out)
}
