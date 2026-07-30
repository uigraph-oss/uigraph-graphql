package uigraphapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type DependencyService struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Status      *string         `json:"status,omitempty"`
	Tier        *string         `json:"tier,omitempty"`
	Category    *string         `json:"category,omitempty"`
	Language    *string         `json:"language,omitempty"`
	GitRepoURL  *string         `json:"gitRepoUrl,omitempty"`
	UpdatedAt   *string         `json:"updatedAt,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

type Dependency struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Service          *DependencyService `json:"service,omitempty"`
	Dependency       *DependencyService `json:"dependency,omitempty"`
	DependencyName   string             `json:"dependencyName"`
	Direction        string             `json:"direction"`
	Type             *string            `json:"type,omitempty"`
	Criticality      *string            `json:"criticality,omitempty"`
	Description      *string            `json:"description,omitempty"`
	APIGroupName     *string            `json:"apiGroupName,omitempty"`
	APIEndpointNames []string           `json:"apiEndpointNames,omitempty"`
	DatabaseName     *string            `json:"databaseName,omitempty"`
}

func (c *Client) ListDependencies(ctx context.Context, orgID, serviceID string, direction, criticality *string) ([]Dependency, error) {
	q := url.Values{}
	if direction != nil {
		q.Set("direction", *direction)
	}
	if criticality != nil {
		q.Set("criticality", *criticality)
	}
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/dependencies", orgID, serviceID)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Dependencies []Dependency `json:"edges"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.Dependencies, nil
}

func (c *Client) GetServiceDependencyGraph(ctx context.Context, orgID, serviceID string) ([]Dependency, error) {
	var out struct {
		Edges []Dependency `json:"edges"`
	}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dependency-graph", orgID, serviceID), &out); err != nil {
		return nil, err
	}
	return out.Edges, nil
}

func (c *Client) UpdateServiceDependencies(ctx context.Context, orgID, serviceID string, body map[string]interface{}) ([]Dependency, error) {
	if err := c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dependencies/sync", orgID, serviceID), body, nil); err != nil {
		return nil, err
	}
	return c.GetServiceDependencyGraph(ctx, orgID, serviceID)
}

func (c *Client) GetDependencyGraph(ctx context.Context, orgID string) ([]Dependency, error) {
	var out struct {
		Edges []Dependency `json:"edges"`
	}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/dependency-graph", orgID), &out); err != nil {
		return nil, err
	}
	return out.Edges, nil
}

func (c *Client) GetServiceImpact(ctx context.Context, orgID, serviceID string, direction *string, maxDepth *int) ([]Dependency, error) {
	q := url.Values{}
	if direction != nil {
		q.Set("direction", *direction)
	}
	if maxDepth != nil {
		q.Set("maxDepth", strconv.Itoa(*maxDepth))
	}
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/impact", orgID, serviceID)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Edges []Dependency `json:"edges"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.Edges, nil
}
