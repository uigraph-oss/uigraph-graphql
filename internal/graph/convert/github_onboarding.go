package convert

import (
	"fmt"
	"strconv"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func GitHubAppInstallationToModel(installation *uigraphapi.GitHubAppInstallation) *model.GitHubAppInstallation {
	return &model.GitHubAppInstallation{
		ID:           installation.ID,
		AccountLogin: installation.AccountLogin,
		AccountType:  installation.AccountType,
		Status:       installation.Status,
	}
}

func GitHubRepositoryToModel(repository uigraphapi.GitHubRepository) *model.GitHubRepository {
	return &model.GitHubRepository{
		ID:            repository.ID,
		GithubID:      strconv.FormatInt(repository.GitHubID, 10),
		Name:          repository.Name,
		FullName:      repository.FullName,
		URL:           repository.URL,
		DefaultBranch: repository.DefaultBranch,
		Private:       repository.Private,
		Archived:      repository.Archived,
	}
}

func GitHubRepositoriesToModel(repositories []uigraphapi.GitHubRepository) []*model.GitHubRepository {
	out := make([]*model.GitHubRepository, len(repositories))
	for i, repository := range repositories {
		out[i] = GitHubRepositoryToModel(repository)
	}
	return out
}

func RepositoryOnboardingStatusToModel(status uigraphapi.RepositoryOnboardingStatus) (model.RepositoryOnboardingStatus, error) {
	switch status {
	case uigraphapi.RepositoryOnboardingStatusRunning:
		return model.RepositoryOnboardingStatusRunning, nil
	case uigraphapi.RepositoryOnboardingStatusSelected:
		return model.RepositoryOnboardingStatusSelected, nil
	case uigraphapi.RepositoryOnboardingStatusSetupPRCreating:
		return model.RepositoryOnboardingStatusSetupPrCreating, nil
	case uigraphapi.RepositoryOnboardingStatusSetupPROpen:
		return model.RepositoryOnboardingStatusSetupPrOpen, nil
	case uigraphapi.RepositoryOnboardingStatusWaitingSetupMerge:
		return model.RepositoryOnboardingStatusWaitingSetupMerge, nil
	case uigraphapi.RepositoryOnboardingStatusCheckingAI:
		return model.RepositoryOnboardingStatusCheckingAiConfiguration, nil
	case uigraphapi.RepositoryOnboardingStatusWaitingAI:
		return model.RepositoryOnboardingStatusWaitingAiConfiguration, nil
	case uigraphapi.RepositoryOnboardingStatusGenerationQueued:
		return model.RepositoryOnboardingStatusGenerationQueued, nil
	case uigraphapi.RepositoryOnboardingStatusGenerationRunning:
		return model.RepositoryOnboardingStatusGenerationRunning, nil
	case uigraphapi.RepositoryOnboardingStatusArtifactsPROpen:
		return model.RepositoryOnboardingStatusArtifactsPrOpen, nil
	case uigraphapi.RepositoryOnboardingStatusWaitingArtifactsMerge:
		return model.RepositoryOnboardingStatusWaitingArtifactsMerge, nil
	case uigraphapi.RepositoryOnboardingStatusSyncQueued:
		return model.RepositoryOnboardingStatusSyncQueued, nil
	case uigraphapi.RepositoryOnboardingStatusSyncRunning:
		return model.RepositoryOnboardingStatusSyncRunning, nil
	case uigraphapi.RepositoryOnboardingStatusCompleted:
		return model.RepositoryOnboardingStatusCompleted, nil
	case uigraphapi.RepositoryOnboardingStatusFailed:
		return model.RepositoryOnboardingStatusFailed, nil
	case uigraphapi.RepositoryOnboardingStatusCancelled:
		return model.RepositoryOnboardingStatusCancelled, nil
	default:
		return "", fmt.Errorf("unsupported repository onboarding status %q", status)
	}
}

func RepositoryOnboardingToModel(onboarding uigraphapi.RepositoryOnboarding) (*model.RepositoryOnboarding, error) {
	status, err := RepositoryOnboardingStatusToModel(onboarding.Status)
	if err != nil {
		return nil, err
	}
	missingAIConfiguration := onboarding.MissingAIConfiguration
	if missingAIConfiguration == nil {
		missingAIConfiguration = []string{}
	}
	return &model.RepositoryOnboarding{
		ID:                      onboarding.ID,
		Repository:              GitHubRepositoryToModel(onboarding.Repository),
		Status:                  status,
		SetupPullRequestURL:     onboarding.SetupPullRequestURL,
		GenerationRunURL:        onboarding.GenerationRunURL,
		ArtifactsPullRequestURL: onboarding.ArtifactsPullRequestURL,
		SyncRunURL:              onboarding.SyncRunURL,
		MissingAIConfiguration:  missingAIConfiguration,
		Error:                   onboarding.Error,
		ServiceID:               onboarding.ServiceID,
	}, nil
}

func RepositoryOnboardingBatchToModel(batch *uigraphapi.RepositoryOnboardingBatch) (*model.RepositoryOnboardingBatch, error) {
	status, err := RepositoryOnboardingStatusToModel(batch.Status)
	if err != nil {
		return nil, err
	}
	repositories := make([]*model.RepositoryOnboarding, len(batch.Repositories))
	for i, repository := range batch.Repositories {
		repositories[i], err = RepositoryOnboardingToModel(repository)
		if err != nil {
			return nil, err
		}
	}
	return &model.RepositoryOnboardingBatch{
		ID:           batch.ID,
		Status:       status,
		TeamID:       batch.TeamID,
		TeamName:     batch.TeamName,
		Repositories: repositories,
	}, nil
}
