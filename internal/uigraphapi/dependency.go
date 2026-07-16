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
	ConsumerService  DependencyService  `json:"consumer"`
	ProviderService  *DependencyService `json:"provider,omitempty"`
	ProviderName     *string            `json:"providerName,omitempty"`
	Type             *string            `json:"type,omitempty"`
	Criticality      *string            `json:"criticality,omitempty"`
	Description      *string            `json:"description,omitempty"`
	APIGroupName     *string            `json:"apiGroupName,omitempty"`
	APIEndpointNames []string           `json:"apiEndpointNames,omitempty"`
	DatabaseName     *string            `json:"databaseName,omitempty"`
	Direction        *string            `json:"direction,omitempty"`
}

type DependencyGraphNode struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Type             *string            `json:"type,omitempty"`
	Service          *DependencyService `json:"service,omitempty"`
	Depth            *int               `json:"depth,omitempty"`
	Metadata         json.RawMessage    `json:"metadata,omitempty"`
}

type DependencyGraphEdge struct {
	ID                  string             `json:"id"`
	Source              string             `json:"source,omitempty"`
	Target              string             `json:"target,omitempty"`
	SourceServiceID     string             `json:"sourceServiceId,omitempty"`
	ProviderServiceName string             `json:"providerName,omitempty"`
	Consumer            *DependencyService `json:"consumer,omitempty"`
	Provider            *DependencyService `json:"provider,omitempty"`
	DependencyID        *string            `json:"dependencyId,omitempty"`
	Type                *string            `json:"type,omitempty"`
	Criticality         *string            `json:"criticality,omitempty"`
	Direction           *string            `json:"direction,omitempty"`
	Depth               *int               `json:"depth,omitempty"`
	APIGroupName        *string            `json:"apiGroupName,omitempty"`
	APIEndpointNames    []string           `json:"apiEndpointNames,omitempty"`
	DatabaseName        *string            `json:"databaseName,omitempty"`
	Metadata            json.RawMessage    `json:"metadata,omitempty"`
}

type DependencyGraph struct {
	Nodes []DependencyGraphNode `json:"nodes"`
	Edges []DependencyGraphEdge `json:"edges"`
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
	return out.Dependencies, c.get(ctx, path, &out)
}

func (c *Client) GetServiceDependencyGraph(ctx context.Context, orgID, serviceID string) (*DependencyGraph, error) {
	var out DependencyGraph
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dependency-graph", orgID, serviceID), &out)
}

func (c *Client) UpdateServiceDependencies(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*DependencyGraph, error) {
	if err := c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dependencies/sync", orgID, serviceID), body, nil); err != nil {
		return nil, err
	}
	return c.GetServiceDependencyGraph(ctx, orgID, serviceID)
}

func (c *Client) GetDependencyGraph(ctx context.Context, orgID string) (*DependencyGraph, error) {
	var out DependencyGraph
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/dependency-graph", orgID), &out)
}

func (c *Client) GetServiceImpact(ctx context.Context, orgID, serviceID string, direction *string, maxDepth *int) (*DependencyGraph, error) {
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
	var out DependencyGraph
	return &out, c.get(ctx, path, &out)
}
