package uigraphapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type AgentSession struct {
	ID               string             `json:"id"`
	OrgID            string             `json:"orgId"`
	Type             string             `json:"type"`
	Status           string             `json:"status"`
	UserID           *string            `json:"userId,omitempty"`
	ServiceAccountID *string            `json:"serviceAccountId,omitempty"`
	Title            *string            `json:"title,omitempty"`
	ModelName        *string            `json:"modelName,omitempty"`
	Metadata         json.RawMessage    `json:"metadata,omitempty"`
	Report           *string            `json:"report,omitempty"`
	Error            *string            `json:"error,omitempty"`
	StartedAt        time.Time          `json:"startedAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
	CompletedAt      *time.Time         `json:"completedAt,omitempty"`
	Totals           AgentSessionTotals `json:"totals"`
}

type AgentSessionTotals struct {
	StepCount          int      `json:"stepCount"`
	InputTokens        int      `json:"inputTokens"`
	OutputTokens       int      `json:"outputTokens"`
	ReasoningTokens    int      `json:"reasoningTokens"`
	CachedInputTokens  int      `json:"cachedInputTokens"`
	CachedOutputTokens int      `json:"cachedOutputTokens"`
	CostUSD            *float64 `json:"costUsd,omitempty"`
	UnpricedSteps      int      `json:"unpricedSteps"`
	StepDurationMs     int      `json:"stepDurationMs"`
}

type AgentSessionStep struct {
	ID                 string          `json:"id"`
	SessionID          string          `json:"sessionId"`
	Seq                int             `json:"seq"`
	Kind               string          `json:"kind"`
	Name               *string         `json:"name,omitempty"`
	ModelName          *string         `json:"modelName,omitempty"`
	Input              json.RawMessage `json:"input,omitempty"`
	Output             json.RawMessage `json:"output,omitempty"`
	Text               *string         `json:"text,omitempty"`
	FinishReason       *string         `json:"finishReason,omitempty"`
	Error              *string         `json:"error,omitempty"`
	InputTokens        *int            `json:"inputTokens,omitempty"`
	OutputTokens       *int            `json:"outputTokens,omitempty"`
	ReasoningTokens    *int            `json:"reasoningTokens,omitempty"`
	CachedInputTokens  *int            `json:"cachedInputTokens,omitempty"`
	CachedOutputTokens *int            `json:"cachedOutputTokens,omitempty"`
	CostUSD            *float64        `json:"costUsd,omitempty"`
	StartedAt          time.Time       `json:"startedAt"`
	CompletedAt        time.Time       `json:"completedAt"`
}

type AgentSessionPage struct {
	Sessions []AgentSession `json:"sessions"`
	Total    int            `json:"total"`
	Period   string         `json:"period"`
	Limit    int            `json:"limit"`
	Offset   int            `json:"offset"`
}

type AgentSessionDetail struct {
	Session AgentSession       `json:"session"`
	Steps   []AgentSessionStep `json:"steps"`
}

type AgentSessionSummary struct {
	Period            string                    `json:"period"`
	TotalSessions     int                       `json:"totalSessions"`
	CompletedSessions int                       `json:"completedSessions"`
	FailedSessions    int                       `json:"failedSessions"`
	RunningSessions   int                       `json:"runningSessions"`
	TotalDurationMs   int                       `json:"totalDurationMs"`
	Totals            AgentSessionTotals        `json:"totals"`
	ByType            []AgentSessionTypeSummary `json:"byType"`
}

type AgentSessionTypeSummary struct {
	Type              string             `json:"type"`
	TotalSessions     int                `json:"totalSessions"`
	CompletedSessions int                `json:"completedSessions"`
	FailedSessions    int                `json:"failedSessions"`
	RunningSessions   int                `json:"runningSessions"`
	TotalDurationMs   int                `json:"totalDurationMs"`
	Totals            AgentSessionTotals `json:"totals"`
}

func (c *Client) GetAgentSessions(ctx context.Context, orgID string, sessionType, status, period *string, limit, offset *int) (*AgentSessionPage, error) {
	q := url.Values{}
	if sessionType != nil && *sessionType != "" {
		q.Set("type", *sessionType)
	}
	if status != nil && *status != "" {
		q.Set("status", *status)
	}
	if period != nil && *period != "" {
		q.Set("period", *period)
	}
	if limit != nil {
		q.Set("limit", strconv.Itoa(*limit))
	}
	if offset != nil {
		q.Set("offset", strconv.Itoa(*offset))
	}
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/agent-sessions", orgID), q)
	var out AgentSessionPage
	return &out, c.get(ctx, path, &out)
}

func (c *Client) GetAgentSession(ctx context.Context, orgID, id string) (*AgentSessionDetail, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/agent-sessions/%s", orgID, id)
	var out AgentSessionDetail
	return &out, c.get(ctx, path, &out)
}

func (c *Client) GetAgentSessionSummary(ctx context.Context, orgID string, period, sessionType *string) (*AgentSessionSummary, error) {
	q := url.Values{}
	if period != nil && *period != "" {
		q.Set("period", *period)
	}
	if sessionType != nil && *sessionType != "" {
		q.Set("type", *sessionType)
	}
	path := withQuery(fmt.Sprintf("/api/v1/orgs/%s/agent-sessions/summary", orgID), q)
	var out struct {
		Period  string               `json:"period"`
		Summary *AgentSessionSummary `json:"summary"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	if out.Summary == nil {
		return &AgentSessionSummary{Period: out.Period}, nil
	}
	out.Summary.Period = out.Period
	return out.Summary, nil
}
