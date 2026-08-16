package uigraphapi

import (
	"context"
	"fmt"
)

type GitHubAppInstallation struct {
	ID           string `json:"id"`
	AccountLogin string `json:"account"`
	AccountType  string `json:"accountType"`
	Status       string `json:"status"`
}

type GitHubRepository struct {
	ID            string `json:"id"`
	GitHubID      int64  `json:"githubId"`
	Name          string `json:"name"`
	FullName      string `json:"fullName"`
	URL           string `json:"url"`
	DefaultBranch string `json:"defaultBranch"`
	Private       bool   `json:"private"`
	Archived      bool   `json:"archived"`
}

type RepositoryOnboardingStatus string

const (
	RepositoryOnboardingStatusRunning    RepositoryOnboardingStatus = "running"
	RepositoryOnboardingStatusSelected   RepositoryOnboardingStatus = "selected"
	RepositoryOnboardingStatusCheckingAI RepositoryOnboardingStatus = "checking_ai_configuration"
	RepositoryOnboardingStatusWaitingAI  RepositoryOnboardingStatus = "waiting_ai_configuration"
	RepositoryOnboardingStatusRunQueued  RepositoryOnboardingStatus = "run_queued"
	RepositoryOnboardingStatusRunRunning RepositoryOnboardingStatus = "run_running"
	RepositoryOnboardingStatusCompleted  RepositoryOnboardingStatus = "completed"
	RepositoryOnboardingStatusFailed     RepositoryOnboardingStatus = "failed"
	RepositoryOnboardingStatusCancelled  RepositoryOnboardingStatus = "cancelled"
)

type RepositoryOnboardingBatch struct {
	ID           string                     `json:"id"`
	Status       RepositoryOnboardingStatus `json:"status"`
	TeamID       string                     `json:"teamId"`
	TeamName     *string                    `json:"team,omitempty"`
	Repositories []RepositoryOnboarding     `json:"repositories"`
}

type RepositoryOnboarding struct {
	ID                     string                     `json:"id"`
	Repository             GitHubRepository           `json:"repository"`
	Status                 RepositoryOnboardingStatus `json:"status"`
	Branch                 string                     `json:"branch"`
	RunURL                 *string                    `json:"runUrl,omitempty"`
	PullRequestURL         *string                    `json:"prUrl,omitempty"`
	MissingAIConfiguration []string                   `json:"missingAIConfiguration,omitempty"`
	Error                  *string                    `json:"error,omitempty"`
	ServiceID              *string                    `json:"serviceId,omitempty"`
}

type StartRepositoryOnboardingInput struct {
	TeamID        string   `json:"teamId"`
	RepositoryIDs []string `json:"repositoryIds"`
}

func (c *Client) GetGitHubApp(ctx context.Context, orgID string) (*GitHubAppInstallation, error) {
	var out struct {
		Installation *GitHubAppInstallation `json:"installation"`
	}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/github-app", orgID), &out); err != nil {
		return nil, err
	}
	return out.Installation, nil
}

func (c *Client) GetGitHubAppInstallURL(ctx context.Context, orgID string) (string, error) {
	var out struct {
		URL string `json:"url"`
	}
	return out.URL, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/github-app/install", orgID), nil, &out)
}

func (c *Client) DisconnectGitHubApp(ctx context.Context, orgID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/github-app", orgID))
}

func (c *Client) ListGitHubRepositories(ctx context.Context, orgID string) ([]GitHubRepository, error) {
	var out struct {
		Repositories []GitHubRepository `json:"repositories"`
	}
	return out.Repositories, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/github-app/repositories", orgID), &out)
}

func (c *Client) StartRepositoryOnboarding(ctx context.Context, orgID string, input StartRepositoryOnboardingInput) (*RepositoryOnboardingBatch, error) {
	var out RepositoryOnboardingBatch
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-onboarding", orgID), input, &out)
}

func (c *Client) GetRepositoryOnboarding(ctx context.Context, orgID, batchID string) (*RepositoryOnboardingBatch, error) {
	var out RepositoryOnboardingBatch
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-onboarding/%s", orgID, batchID), &out)
}

func (c *Client) GetLatestRepositoryOnboarding(ctx context.Context, orgID string) (*RepositoryOnboardingBatch, error) {
	var out struct {
		Batch *RepositoryOnboardingBatch `json:"batch"`
	}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-onboarding", orgID), &out); err != nil {
		return nil, err
	}
	return out.Batch, nil
}

func (c *Client) RecheckRepositoryOnboarding(ctx context.Context, orgID, batchID, onboardingID string) (*RepositoryOnboarding, error) {
	var out RepositoryOnboarding
	path := fmt.Sprintf("/api/v1/orgs/%s/repository-onboarding/%s/repositories/%s/recheck", orgID, batchID, onboardingID)
	return &out, c.post(ctx, path, nil, &out)
}

func (c *Client) RetryRepositoryOnboarding(ctx context.Context, orgID, batchID, onboardingID string) (*RepositoryOnboarding, error) {
	var out RepositoryOnboarding
	path := fmt.Sprintf("/api/v1/orgs/%s/repository-onboarding/%s/repositories/%s/retry", orgID, batchID, onboardingID)
	return &out, c.post(ctx, path, nil, &out)
}
