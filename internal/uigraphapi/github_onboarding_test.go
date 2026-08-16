package uigraphapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGitHubOnboardingClientRequests(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		response string
		wantBody map[string]interface{}
		request  func(*Client) error
	}{
		{
			name:     "get installation",
			method:   http.MethodGet,
			path:     "/api/v1/orgs/org-1/github-app",
			response: `{"id":"installation-1","account":"acme","accountType":"Organization","status":"active"}`,
			request: func(c *Client) error {
				_, err := c.GetGitHubApp(context.Background(), "org-1")
				return err
			},
		},
		{
			name:     "create install url",
			method:   http.MethodPost,
			path:     "/api/v1/orgs/org-1/github-app/install",
			response: `{"url":"https://github.com/apps/uigraph/installations/new"}`,
			request: func(c *Client) error {
				_, err := c.GetGitHubAppInstallURL(context.Background(), "org-1")
				return err
			},
		},
		{
			name:   "disconnect installation",
			method: http.MethodDelete,
			path:   "/api/v1/orgs/org-1/github-app",
			request: func(c *Client) error {
				return c.DisconnectGitHubApp(context.Background(), "org-1")
			},
		},
		{
			name:     "list repositories",
			method:   http.MethodGet,
			path:     "/api/v1/orgs/org-1/github-app/repositories",
			response: `{"repositories":[]}`,
			request: func(c *Client) error {
				_, err := c.ListGitHubRepositories(context.Background(), "org-1")
				return err
			},
		},
		{
			name:     "start onboarding",
			method:   http.MethodPost,
			path:     "/api/v1/orgs/org-1/repository-onboarding",
			response: `{"id":"batch-1","status":"running","teamId":"team-1","repositories":[]}`,
			wantBody: map[string]interface{}{
				"teamId":        "team-1",
				"repositoryIds": []interface{}{"repo-1", "repo-2"},
			},
			request: func(c *Client) error {
				_, err := c.StartRepositoryOnboarding(context.Background(), "org-1", StartRepositoryOnboardingInput{
					TeamID:        "team-1",
					RepositoryIDs: []string{"repo-1", "repo-2"},
				})
				return err
			},
		},
		{
			name:     "get onboarding",
			method:   http.MethodGet,
			path:     "/api/v1/orgs/org-1/repository-onboarding/batch-1",
			response: `{"id":"batch-1","status":"running","teamId":"team-1","repositories":[]}`,
			request: func(c *Client) error {
				_, err := c.GetRepositoryOnboarding(context.Background(), "org-1", "batch-1")
				return err
			},
		},
		{
			name:     "recheck repository",
			method:   http.MethodPost,
			path:     "/api/v1/orgs/org-1/repository-onboarding/batch-1/repositories/onboarding-1/recheck",
			response: onboardingResponseJSON,
			request: func(c *Client) error {
				_, err := c.RecheckRepositoryOnboarding(context.Background(), "org-1", "batch-1", "onboarding-1")
				return err
			},
		},
		{
			name:     "retry repository",
			method:   http.MethodPost,
			path:     "/api/v1/orgs/org-1/repository-onboarding/batch-1/repositories/onboarding-1/retry",
			response: onboardingResponseJSON,
			request: func(c *Client) error {
				_, err := c.RetryRepositoryOnboarding(context.Background(), "org-1", "batch-1", "onboarding-1")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method || r.URL.Path != tt.path {
					t.Errorf("request = %s %s, want %s %s", r.Method, r.URL.Path, tt.method, tt.path)
				}
				if tt.wantBody != nil {
					var body map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode request body: %v", err)
					}
					if !reflect.DeepEqual(body, tt.wantBody) {
						t.Errorf("body = %#v, want %#v", body, tt.wantBody)
					}
				}
				if tt.response == "" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			if err := tt.request(New(server.URL)); err != nil {
				t.Fatalf("request error = %v", err)
			}
		})
	}
}

func TestGetRepositoryOnboardingDecodesNullableFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"batch-1",
			"status":"running",
			"teamId":"team-1",
			"team":null,
			"repositories":[{
				"id":"onboarding-1",
				"repository":{"id":"repo-node-1","githubId":100,"name":"api","fullName":"acme/api","url":"https://github.com/acme/api","defaultBranch":"main","private":true,"archived":false},
				"status":"waiting_ai_configuration",
				"setupPrUrl":null,
				"missingAIConfiguration":["OPENAI_API_KEY"],
				"error":null
			}]
		}`))
	}))
	defer server.Close()

	batch, err := New(server.URL).GetRepositoryOnboarding(context.Background(), "org-1", "batch-1")
	if err != nil {
		t.Fatalf("GetRepositoryOnboarding() error = %v", err)
	}
	if batch.TeamName != nil {
		t.Errorf("TeamName = %q, want nil", *batch.TeamName)
	}
	onboarding := batch.Repositories[0]
	if onboarding.SetupPullRequestURL != nil || onboarding.Error != nil || onboarding.ServiceID != nil {
		t.Fatalf("nullable fields = %#v, want nil", onboarding)
	}
	if !reflect.DeepEqual(onboarding.MissingAIConfiguration, []string{"OPENAI_API_KEY"}) {
		t.Errorf("MissingAIConfiguration = %#v", onboarding.MissingAIConfiguration)
	}
	if onboarding.Status != RepositoryOnboardingStatusWaitingAI {
		t.Errorf("Status = %q, want %q", onboarding.Status, RepositoryOnboardingStatusWaitingAI)
	}
}

const onboardingResponseJSON = `{"id":"onboarding-1","repository":{"id":"repo-node-1","githubId":100,"name":"api","fullName":"acme/api","url":"https://github.com/acme/api","defaultBranch":"main","private":true,"archived":false},"status":"generation_running","missingAIConfiguration":[]}`
