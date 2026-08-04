package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type TimelineTouch struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Kind  string `json:"kind"`
}

type TimelineEvent struct {
	ID                string          `json:"id"`
	OrgID             string          `json:"orgId"`
	ServiceID         string          `json:"serviceId"`
	Type              string          `json:"type"`
	Title             string          `json:"title"`
	Summary           string          `json:"summary"`
	EventDate         time.Time       `json:"eventDate"`
	Version           *string         `json:"version,omitempty"`
	ADRNumber         *string         `json:"adrNumber,omitempty"`
	DecisionStatus    *string         `json:"decisionStatus,omitempty"`
	SourceLabel       *string         `json:"sourceLabel,omitempty"`
	SourceURL         *string         `json:"sourceUrl,omitempty"`
	IsAgentSummarized bool            `json:"isAgentSummarized"`
	Origin            string          `json:"origin"`
	Touches           []TimelineTouch `json:"touches"`

	AttachmentAssetID  *string `json:"attachmentAssetId,omitempty"`
	AttachmentFileName *string `json:"attachmentFileName,omitempty"`
	AttachmentFileType *string `json:"attachmentFileType,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Client) ListServiceTimelineEvents(ctx context.Context, orgID, serviceID string) ([]TimelineEvent, error) {
	var out struct {
		Events []TimelineEvent `json:"events"`
	}
	return out.Events, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/timeline", orgID, serviceID), &out)
}

func (c *Client) CreateTimelineEvent(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TimelineEvent, error) {
	var out TimelineEvent
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/timeline", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTimelineEvent(ctx context.Context, orgID, serviceID, eventID string, body map[string]interface{}) (*TimelineEvent, error) {
	var out TimelineEvent
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/timeline/%s", orgID, serviceID, eventID), body, &out)
}

func (c *Client) DeleteTimelineEvent(ctx context.Context, orgID, serviceID, eventID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/timeline/%s", orgID, serviceID, eventID))
}
