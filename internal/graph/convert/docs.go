package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func DocToModel(d *uigraphapi.Doc) *model.Doc {
	return &model.Doc{
		ID: d.ID, OrgID: d.OrgID, FolderID: d.FolderID, TeamID: d.TeamID,
		FileAssetID: d.FileAssetID, FileName: d.FileName, FileType: d.FileType,
		Description: d.Description, ContentHash: d.ContentHash,
		CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func DocsToModel(docs []uigraphapi.Doc) []*model.Doc {
	out := make([]*model.Doc, len(docs))
	for i := range docs {
		out[i] = DocToModel(&docs[i])
	}
	return out
}
