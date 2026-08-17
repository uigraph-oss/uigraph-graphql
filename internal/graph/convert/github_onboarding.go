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

func repositoryImportStatusToModel(status uigraphapi.RepositoryImportStatus) (model.RepositoryImportStatus, error) {
	switch status {
	case uigraphapi.RepositoryImportStatusSelected:
		return model.RepositoryImportStatusSelected, nil
	case uigraphapi.RepositoryImportStatusCheckingAI:
		return model.RepositoryImportStatusCheckingAiConfiguration, nil
	case uigraphapi.RepositoryImportStatusWaitingAI:
		return model.RepositoryImportStatusWaitingAiConfiguration, nil
	case uigraphapi.RepositoryImportStatusRunQueued:
		return model.RepositoryImportStatusRunQueued, nil
	case uigraphapi.RepositoryImportStatusRunRunning:
		return model.RepositoryImportStatusRunRunning, nil
	case uigraphapi.RepositoryImportStatusCompleted:
		return model.RepositoryImportStatusCompleted, nil
	case uigraphapi.RepositoryImportStatusFailed:
		return model.RepositoryImportStatusFailed, nil
	default:
		return "", fmt.Errorf("unsupported repository import status %q", status)
	}
}

func RepositoryImportToModel(value uigraphapi.RepositoryImport) (*model.RepositoryImport, error) {
	status, err := repositoryImportStatusToModel(value.Status)
	if err != nil {
		return nil, err
	}
	steps := make([]*model.RepositoryImportStep, len(value.Steps))
	for i, step := range value.Steps {
		steps[i] = &model.RepositoryImportStep{
			Number:      step.Number,
			Name:        step.Name,
			Status:      step.Status,
			StartedAt:   step.StartedAt,
			CompletedAt: step.CompletedAt,
		}
		if step.Conclusion != "" {
			steps[i].Conclusion = &step.Conclusion
		}
	}
	return &model.RepositoryImport{
		ID:                     value.ID,
		Repository:             GitHubRepositoryToModel(value.Repository),
		Status:                 status,
		Steps:                  steps,
		TeamID:                 value.TeamID,
		TeamName:               value.TeamName,
		Branch:                 value.Branch,
		RunURL:                 value.RunURL,
		PullRequestURL:         value.PullRequestURL,
		MissingAIConfiguration: value.MissingAIConfiguration,
		Error:                  value.Error,
		ServiceID:              value.ServiceID,
		CreatedAt:              value.CreatedAt,
		RunStartedAt:           value.RunStartedAt,
		RunCompletedAt:         value.RunCompletedAt,
	}, nil
}

func RepositoryImportsToModel(values []uigraphapi.RepositoryImport) ([]*model.RepositoryImport, error) {
	out := make([]*model.RepositoryImport, len(values))
	for i, value := range values {
		converted, err := RepositoryImportToModel(value)
		if err != nil {
			return nil, err
		}
		out[i] = converted
	}
	return out, nil
}
