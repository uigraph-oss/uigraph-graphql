package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestRepositoryOnboardingStatusToModel(t *testing.T) {
	tests := []struct {
		apiStatus uigraphapi.RepositoryOnboardingStatus
		want      model.RepositoryOnboardingStatus
	}{
		{uigraphapi.RepositoryOnboardingStatusRunning, model.RepositoryOnboardingStatusRunning},
		{uigraphapi.RepositoryOnboardingStatusSelected, model.RepositoryOnboardingStatusSelected},
		{uigraphapi.RepositoryOnboardingStatusCheckingAI, model.RepositoryOnboardingStatusCheckingAiConfiguration},
		{uigraphapi.RepositoryOnboardingStatusWaitingAI, model.RepositoryOnboardingStatusWaitingAiConfiguration},
		{uigraphapi.RepositoryOnboardingStatusRunQueued, model.RepositoryOnboardingStatusRunQueued},
		{uigraphapi.RepositoryOnboardingStatusRunRunning, model.RepositoryOnboardingStatusRunRunning},
		{uigraphapi.RepositoryOnboardingStatusCompleted, model.RepositoryOnboardingStatusCompleted},
		{uigraphapi.RepositoryOnboardingStatusFailed, model.RepositoryOnboardingStatusFailed},
		{uigraphapi.RepositoryOnboardingStatusCancelled, model.RepositoryOnboardingStatusCancelled},
	}
	for _, tt := range tests {
		got, err := RepositoryOnboardingStatusToModel(tt.apiStatus)
		if err != nil {
			t.Fatalf("RepositoryOnboardingStatusToModel(%q) error = %v", tt.apiStatus, err)
		}
		if got != tt.want {
			t.Errorf("RepositoryOnboardingStatusToModel(%q) = %q, want %q", tt.apiStatus, got, tt.want)
		}
	}
}

func TestRepositoryOnboardingStatusToModelRejectsUnknownStatus(t *testing.T) {
	if _, err := RepositoryOnboardingStatusToModel("unknown"); err == nil {
		t.Fatal("RepositoryOnboardingStatusToModel(unknown) error = nil")
	}
}

func TestRepositoryOnboardingToModelPreservesNullableFields(t *testing.T) {
	got, err := RepositoryOnboardingToModel(uigraphapi.RepositoryOnboarding{
		ID:     "onboarding-1",
		Status: uigraphapi.RepositoryOnboardingStatusSelected,
		Repository: uigraphapi.GitHubRepository{
			ID:       "repo-node-1",
			GitHubID: 100,
		},
	})
	if err != nil {
		t.Fatalf("RepositoryOnboardingToModel() error = %v", err)
	}
	if got.RunURL != nil || got.PullRequestURL != nil || got.Error != nil || got.ServiceID != nil {
		t.Fatalf("nullable fields were not nil: %#v", got)
	}
	if got.MissingAIConfiguration == nil || len(got.MissingAIConfiguration) != 0 {
		t.Errorf("MissingAIConfiguration = %#v, want non-nil empty slice", got.MissingAIConfiguration)
	}
	if got.Repository.GithubID != "100" {
		t.Errorf("GithubID = %q, want 100", got.Repository.GithubID)
	}
}
