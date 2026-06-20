package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type FlowDiagramComponentField struct {
	FlowDiagramComponentFieldID string   `json:"flowDiagramComponentFieldId"`
	Label                       string   `json:"label"`
	Type                        string   `json:"type"`
	Required                    bool     `json:"required"`
	Readonly                    *bool    `json:"readonly,omitempty"`
	Options                     []string `json:"options,omitempty"`
	Order                       int      `json:"order"`
}

type FlowDiagramComponent struct {
	ComponentID                string                      `json:"componentId"`
	Type                       string                      `json:"type"`
	Name                       string                      `json:"name"`
	Description                string                      `json:"description"`
	Category                   string                      `json:"category"`
	Tags                       []string                    `json:"tags"`
	Slug                       string                      `json:"slug"`
	PreviewImageJpg            string                      `json:"previewImageJpg"`
	IsActive                   bool                        `json:"isActive"`
	Order                      int                         `json:"order"`
	OrganizationID             *string                     `json:"organizationId,omitempty"`
	FlowDiagramComponentFields []FlowDiagramComponentField `json:"flowDiagramComponentFields"`
}

type FlowComponents struct {
	Components       []FlowDiagramComponent `json:"components"`
	CustomComponents []FlowDiagramComponent `json:"customComponents"`
}

type ComponentField struct {
	ComponentFieldID string   `json:"componentFieldId"`
	Label            string   `json:"label"`
	Type             string   `json:"type"`
	Required         bool     `json:"required"`
	Readonly         *bool    `json:"readonly,omitempty"`
	Options          []string `json:"options,omitempty"`
	Order            int      `json:"order"`
}

type Component struct {
	ComponentID     string           `json:"componentId"`
	Type            string           `json:"type"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Category        string           `json:"category"`
	Tags            []string         `json:"tags"`
	Slug            string           `json:"slug"`
	PreviewImageJpg string           `json:"previewImageJpg"`
	IsActive        bool             `json:"isActive"`
	Order           int              `json:"order"`
	ComponentFields []ComponentField `json:"componentFields"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type Components struct {
	Components       []Component `json:"components"`
	CustomComponents []Component `json:"customComponents"`
}

func (c *Client) ListFlowDiagramComponents(ctx context.Context, orgID string) (*FlowComponents, error) {
	var out FlowComponents
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/flow-diagram-components", orgID), &out)
}

func (c *Client) ListComponents(ctx context.Context, orgID string) (*Components, error) {
	var out Components
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/components", orgID), &out)
}

func (c *Client) CreateCustomComponent(ctx context.Context, orgID string, body map[string]interface{}) (*Component, error) {
	var out Component
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/components", orgID), body, &out)
}

func (c *Client) UpdateCustomComponent(ctx context.Context, orgID, id string, body map[string]interface{}) (*Component, error) {
	var out Component
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/components/%s", orgID, id), body, &out)
}

func (c *Client) DeleteCustomComponent(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/components/%s", orgID, id))
}
