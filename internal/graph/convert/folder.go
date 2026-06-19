package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func FolderToModel(f *uigraphapi.Folder) *model.Folder {
	return &model.Folder{
		ID: f.ID, OrgID: f.OrgID, ParentID: f.ParentID, TeamID: f.TeamID, Type: f.Type,
		Name: f.Name, Order: f.Order, CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func FoldersToModel(folders []uigraphapi.Folder) []*model.Folder {
	out := make([]*model.Folder, len(folders))
	for i := range folders {
		out[i] = FolderToModel(&folders[i])
	}
	return out
}
