package uigraphapi

import (
	"context"
	"fmt"
	"time"
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

type RepositoryImportStatus string

const (
	RepositoryImportStatusSelected   RepositoryImportStatus = "selected"
	RepositoryImportStatusCheckingAI RepositoryImportStatus = "checking_ai_configuration"
	RepositoryImportStatusWaitingAI  RepositoryImportStatus = "waiting_ai_configuration"
	RepositoryImportStatusRunQueued  RepositoryImportStatus = "run_queued"
	RepositoryImportStatusRunRunning RepositoryImportStatus = "run_running"
	RepositoryImportStatusCompleted  RepositoryImportStatus = "completed"
	RepositoryImportStatusFailed     RepositoryImportStatus = "failed"
	RepositoryImportStatusCancelled  RepositoryImportStatus = "cancelled"
)

type RepositoryImportStep struct {
	Number      int        `json:"number"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Conclusion  string     `json:"conclusion,omitempty"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

type RepositoryImport struct {
	ID                     string                 `json:"id"`
	Repository             GitHubRepository       `json:"repository"`
	Status                 RepositoryImportStatus `json:"status"`
	Steps                  []RepositoryImportStep `json:"steps"`
	TeamID                 string                 `json:"teamId"`
	TeamName               *string                `json:"team,omitempty"`
	Branch                 string                 `json:"branch"`
	RunURL                 *string                `json:"runUrl,omitempty"`
	PullRequestURL         *string                `json:"prUrl,omitempty"`
	MissingAIConfiguration []string               `json:"missingAIConfiguration,omitempty"`
	Error                  *string                `json:"error,omitempty"`
	ServiceID              *string                `json:"serviceId,omitempty"`
	CreatedAt              time.Time              `json:"createdAt"`
	RunStartedAt           *time.Time             `json:"runStartedAt,omitempty"`
	RunCompletedAt         *time.Time             `json:"runCompletedAt,omitempty"`
}

type StartRepositoryImportInput struct {
	TeamID       string `json:"teamId"`
	RepositoryID string `json:"repositoryId"`
}

func (c *Client) GetGitHubApp(ctx context.Context, orgID string) (enabled bool, installation *GitHubAppInstallation, err error) {
	var out struct {
		Enabled      bool                   `json:"enabled"`
		Installation *GitHubAppInstallation `json:"installation"`
	}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/github-app", orgID), &out); err != nil {
		return false, nil, err
	}
	return out.Enabled, out.Installation, nil
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

func (c *Client) StartRepositoryImport(ctx context.Context, orgID string, input StartRepositoryImportInput) (*RepositoryImport, error) {
	var out RepositoryImport
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports", orgID), input, &out)
}

func (c *Client) GetRepositoryImport(ctx context.Context, orgID, importID string) (*RepositoryImport, error) {
	var out RepositoryImport
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports/%s", orgID, importID), &out)
}

func (c *Client) ListRepositoryImports(ctx context.Context, orgID string) ([]RepositoryImport, error) {
	var out struct {
		Imports []RepositoryImport `json:"imports"`
	}
	return out.Imports, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports", orgID), &out)
}

func (c *Client) GetLatestRepositoryImport(ctx context.Context, orgID string) (*RepositoryImport, error) {
	var out struct {
		Import *RepositoryImport `json:"import"`
	}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports/latest", orgID), &out); err != nil {
		return nil, err
	}
	return out.Import, nil
}

func (c *Client) RecheckRepositoryImport(ctx context.Context, orgID, importID string) (*RepositoryImport, error) {
	var out RepositoryImport
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports/%s/recheck", orgID, importID), nil, &out)
}

func (c *Client) RetryRepositoryImport(ctx context.Context, orgID, importID string) (*RepositoryImport, error) {
	var out RepositoryImport
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports/%s/retry", orgID, importID), nil, &out)
}
