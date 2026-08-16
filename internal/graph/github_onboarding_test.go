package graph_test

import (
	"context"
	"errors"
	"testing"

	"github.com/uigraph/graphql/internal/graph"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

type fakeGitHubOnboardingClient struct {
	batch *uigraphapi.RepositoryOnboardingBatch
	err   error
}

func (f *fakeGitHubOnboardingClient) GetGitHubApp(context.Context, string) (*uigraphapi.GitHubAppInstallation, error) {
	return &uigraphapi.GitHubAppInstallation{}, f.err
}

func (f *fakeGitHubOnboardingClient) GetGitHubAppInstallURL(context.Context, string) (string, error) {
	return "https://github.com/apps/uigraph/installations/new", f.err
}

func (f *fakeGitHubOnboardingClient) DisconnectGitHubApp(context.Context, string) error {
	return f.err
}

func (f *fakeGitHubOnboardingClient) ListGitHubRepositories(context.Context, string) ([]uigraphapi.GitHubRepository, error) {
	return nil, f.err
}

func (f *fakeGitHubOnboardingClient) StartRepositoryOnboarding(context.Context, string, uigraphapi.StartRepositoryOnboardingInput) (*uigraphapi.RepositoryOnboardingBatch, error) {
	return f.batch, f.err
}

func (f *fakeGitHubOnboardingClient) GetRepositoryOnboarding(context.Context, string, string) (*uigraphapi.RepositoryOnboardingBatch, error) {
	return f.batch, f.err
}

func (f *fakeGitHubOnboardingClient) RecheckRepositoryOnboarding(context.Context, string, string, string) (*uigraphapi.RepositoryOnboarding, error) {
	if f.batch == nil || len(f.batch.Repositories) == 0 {
		return nil, f.err
	}
	return &f.batch.Repositories[0], f.err
}

func (f *fakeGitHubOnboardingClient) RetryRepositoryOnboarding(context.Context, string, string, string) (*uigraphapi.RepositoryOnboarding, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &f.batch.Repositories[0], nil
}

func TestRepositoryOnboardingQuery(t *testing.T) {
	client := &fakeGitHubOnboardingClient{batch: &uigraphapi.RepositoryOnboardingBatch{
		ID:       "batch-1",
		Status:   uigraphapi.RepositoryOnboardingStatusRunning,
		TeamID:   "team-1",
		TeamName: nil,
		Repositories: []uigraphapi.RepositoryOnboarding{{
			ID:     "onboarding-1",
			Status: uigraphapi.RepositoryOnboardingStatusWaitingAI,
			Repository: uigraphapi.GitHubRepository{
				ID:            "repo-node-1",
				GitHubID:      100,
				Name:          "api",
				FullName:      "acme/api",
				URL:           "https://github.com/acme/api",
				DefaultBranch: "main",
			},
			MissingAIConfiguration: []string{"OPENAI_API_KEY"},
		}},
	}}
	server := newTestServer(&graph.Resolver{GitHubAPI: client})
	defer server.Close()

	data := doGraphQL(t, server, `{
		repositoryOnboarding(orgId: "org-1", batchId: "batch-1") {
			id status teamId teamName
			repositories {
				id status setupPullRequestUrl generationRunUrl artifactsPullRequestUrl syncRunUrl error serviceId missingAIConfiguration
				repository { id githubId name fullName }
			}
		}
	}`)
	batch := data["repositoryOnboarding"].(map[string]interface{})
	if batch["status"] != "RUNNING" || batch["teamName"] != nil {
		t.Fatalf("repositoryOnboarding = %#v", batch)
	}
	onboarding := batch["repositories"].([]interface{})[0].(map[string]interface{})
	if onboarding["setupPullRequestUrl"] != nil || onboarding["serviceId"] != nil {
		t.Fatalf("nullable onboarding fields = %#v", onboarding)
	}
	if onboarding["status"] != "WAITING_AI_CONFIGURATION" {
		t.Errorf("status = %q, want WAITING_AI_CONFIGURATION", onboarding["status"])
	}
}

func TestRetryRepositoryOnboardingError(t *testing.T) {
	client := &fakeGitHubOnboardingClient{err: errors.New("retry unavailable")}
	server := newTestServer(&graph.Resolver{GitHubAPI: client})
	defer server.Close()

	out := doGraphQLRaw(t, server, `mutation {
		retryRepositoryOnboarding(orgId: "org-1", batchId: "batch-1", onboardingId: "onboarding-1") { id }
	}`)
	if _, ok := out["errors"]; !ok {
		t.Fatalf("errors missing from response: %#v", out)
	}
}

func TestGitHubAppQueryReturnsNullWhenNotConnected(t *testing.T) {
	client := &fakeGitHubOnboardingClient{err: &uigraphapi.APIError{Status: 404, Body: `{"error":"not found"}`}}
	server := newTestServer(&graph.Resolver{GitHubAPI: client})
	defer server.Close()

	data := doGraphQL(t, server, `{ githubApp(orgId: "org-1") { id } }`)
	if data["githubApp"] != nil {
		t.Fatalf("githubApp = %#v, want nil", data["githubApp"])
	}
}
