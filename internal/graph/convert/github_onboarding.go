package convert

import (
	"fmt"
	"strconv"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func GitHubAppInstallationToModel(installation *uigraphapi.GitHubAppInstallation) *model.GitHubAppInstallation {
	return &model.GitHubAppInstallation{
		InstallationID: strconv.FormatInt(installation.InstallationID, 10),
		AccountLogin:   installation.AccountLogin,
		AccountType:    installation.AccountType,
		Suspended:      installation.Suspended,
	}
}

func GitHubRepositoryToModel(repository uigraphapi.GitHubRepository) *model.GitHubRepository {
	return &model.GitHubRepository{
		GithubID:      strconv.FormatInt(repository.GitHubID, 10),
		Owner:         repository.Owner,
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
		ID:             value.ID,
		GithubOwnerID:  strconv.FormatInt(value.GitHubOwnerID, 10),
		GithubRepo:     value.GitHubRepo,
		Status:         status,
		Steps:          steps,
		TeamID:         value.TeamID,
		TeamName:       value.TeamName,
		Branch:         value.Branch,
		RunURL:         value.RunURL,
		PullRequestURL: value.PullRequestURL,
		Error:          value.Error,
		ServiceID:      value.ServiceID,
		CreatedAt:      value.CreatedAt,
		RunStartedAt:   value.RunStartedAt,
		RunCompletedAt: value.RunCompletedAt,
	}, nil
}
