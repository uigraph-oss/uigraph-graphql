package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type ChatSession struct {
	ID           string    `json:"id"`
	OrgID        string    `json:"orgId"`
	OwnerUserID  string    `json:"ownerUserId"`
	Title        string    `json:"title"`
	IsPinned     bool      `json:"isPinned"`
	MessageCount int       `json:"messageCount"`
	CreatedBy    string    `json:"createdBy"`
	UpdatedBy    *string   `json:"updatedBy,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ChatMessage struct {
	ID            string    `json:"id"`
	OrgID         string    `json:"orgId"`
	ChatSessionID string    `json:"chatSessionId"`
	Role          string    `json:"role"`
	Content       string    `json:"content"`
	Parts         any       `json:"parts,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

func chatSessionsBasePath(orgID string) string {
	return fmt.Sprintf("/api/v1/orgs/%s/chat-sessions", orgID)
}

func (c *Client) ListChatSessions(ctx context.Context, orgID string) ([]ChatSession, error) {
	var out struct {
		Sessions []ChatSession `json:"sessions"`
	}
	return out.Sessions, c.get(ctx, chatSessionsBasePath(orgID), &out)
}

func (c *Client) CreateChatSession(ctx context.Context, orgID string, body map[string]interface{}) (*ChatSession, error) {
	var out ChatSession
	return &out, c.post(ctx, chatSessionsBasePath(orgID), body, &out)
}

func (c *Client) GetChatSession(ctx context.Context, orgID, id string) (*ChatSession, []ChatMessage, error) {
	var out struct {
		Session  ChatSession   `json:"session"`
		Messages []ChatMessage `json:"messages"`
	}
	err := c.get(ctx, fmt.Sprintf("%s/%s", chatSessionsBasePath(orgID), id), &out)
	if err != nil {
		return nil, nil, err
	}
	return &out.Session, out.Messages, nil
}

func (c *Client) UpdateChatSession(ctx context.Context, orgID, id string, body map[string]interface{}) (*ChatSession, error) {
	var out ChatSession
	return &out, c.put(ctx, fmt.Sprintf("%s/%s", chatSessionsBasePath(orgID), id), body, &out)
}

func (c *Client) DeleteChatSession(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("%s/%s", chatSessionsBasePath(orgID), id))
}

func (c *Client) ListChatMessages(ctx context.Context, orgID, sessionID string) ([]ChatMessage, error) {
	var out struct {
		Messages []ChatMessage `json:"messages"`
	}
	return out.Messages, c.get(ctx, fmt.Sprintf("%s/%s/messages", chatSessionsBasePath(orgID), sessionID), &out)
}

func (c *Client) CreateChatMessage(ctx context.Context, orgID, sessionID string, body map[string]interface{}) (*ChatMessage, error) {
	var out ChatMessage
	return &out, c.post(ctx, fmt.Sprintf("%s/%s/messages", chatSessionsBasePath(orgID), sessionID), body, &out)
}
