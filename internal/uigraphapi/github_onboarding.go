package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type GitHubAppInstallation struct {
	InstallationID int64  `json:"installationId"`
	AccountLogin   string `json:"account"`
	AccountType    string `json:"accountType"`
	Suspended      bool   `json:"suspended"`
}

type GitHubRepository struct {
	GitHubID      int64  `json:"githubId"`
	Owner         string `json:"owner"`
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
)

type RepositoryImportStep struct {
	Number      int        `json:"number"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Conclusion  string     `json:"conclusion"`
	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

type RepositoryImport struct {
	ID                     string                 `json:"id"`
	GitHubOwnerID          int64                  `json:"githubOwnerId"`
	GitHubRepo             string                 `json:"githubRepo"`
	Status                 RepositoryImportStatus `json:"status"`
	Steps                  []RepositoryImportStep `json:"steps"`
	TeamID                 string                 `json:"teamId"`
	TeamName               *string                `json:"team"`
	Branch                 string                 `json:"branch"`
	RunURL                 *string                `json:"runUrl"`
	PullRequestURL         *string                `json:"prUrl"`
	MissingAIConfiguration []string               `json:"missingAIConfiguration"`
	Error                  *string                `json:"error"`
	ServiceID              *string                `json:"serviceId"`
	CreatedAt              time.Time              `json:"createdAt"`
	RunStartedAt           *time.Time             `json:"runStartedAt"`
	RunCompletedAt         *time.Time             `json:"runCompletedAt"`
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

type RepositoryAIConfiguration struct {
	Missing []string `json:"missing"`
	Ready   bool     `json:"ready"`
}

func (c *Client) GetRepositoryAIConfiguration(ctx context.Context, orgID, owner, repo string) (*RepositoryAIConfiguration, error) {
	var out RepositoryAIConfiguration
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/github-app/repositories/%s/%s/ai-configuration", orgID, owner, repo), &out)
}

func (c *Client) StartRepositoryImport(ctx context.Context, orgID, teamID, owner, repo string) (*RepositoryImport, error) {
	var out RepositoryImport
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports", orgID), map[string]string{"teamId": teamID, "owner": owner, "repo": repo}, &out)
}

func (c *Client) GetRepositoryImport(ctx context.Context, orgID, importID string) (*RepositoryImport, error) {
	var out RepositoryImport
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/repository-imports/%s", orgID, importID), &out)
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
