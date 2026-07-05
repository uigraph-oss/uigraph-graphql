package convert

import (
	"strings"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func savedQueryScopeToModel(scope string) model.SavedQueryScope {
	if strings.EqualFold(scope, "team") {
		return model.SavedQueryScopeTeam
	}
	return model.SavedQueryScopePersonal
}

func SavedQueryFolderToModel(f *uigraphapi.SavedQueryFolder) *model.SavedQueryFolder {
	return &model.SavedQueryFolder{
		ID: f.ID, OrgID: f.OrgID, ServiceDbID: f.ServiceDBID, Scope: savedQueryScopeToModel(f.Scope),
		OwnerUserID: f.OwnerUserID, TeamID: f.TeamID, Name: f.Name,
		CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func SavedQueryFoldersToModel(folders []uigraphapi.SavedQueryFolder) []*model.SavedQueryFolder {
	out := make([]*model.SavedQueryFolder, len(folders))
	for i := range folders {
		out[i] = SavedQueryFolderToModel(&folders[i])
	}
	return out
}

func SavedQueryToModel(q *uigraphapi.SavedQuery) *model.SavedQuery {
	return &model.SavedQuery{
		ID: q.ID, OrgID: q.OrgID, ServiceDbID: q.ServiceDBID, FolderID: q.FolderID,
		Scope: savedQueryScopeToModel(q.Scope), OwnerUserID: q.OwnerUserID, TeamID: q.TeamID,
		Title: q.Title, Description: q.Description, QueryText: q.QueryText, Tags: q.Tags,
		Source:    q.Source,
		CreatedBy: q.CreatedBy, UpdatedBy: q.UpdatedBy,
		CreatedByCommitHash: q.CreatedByCommitHash, UpdatedByCommitHash: q.UpdatedByCommitHash,
		CreatedAt: q.CreatedAt, UpdatedAt: q.UpdatedAt,
	}
}

func SavedQueriesToModel(queries []uigraphapi.SavedQuery) []*model.SavedQuery {
	out := make([]*model.SavedQuery, len(queries))
	for i := range queries {
		out[i] = SavedQueryToModel(&queries[i])
	}
	return out
}
