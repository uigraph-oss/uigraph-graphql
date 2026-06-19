package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func DiagramToModel(d *uigraphapi.Diagram) *model.Diagram {
	return &model.Diagram{
		ID: d.ID, OrgID: d.OrgID, FolderID: d.FolderID, TeamID: d.TeamID,
		Name: d.Name, ContentKey: d.ContentKey, ContentHash: d.ContentHash,
		PreviewAssetID: d.PreviewAssetID, PreviewContentHash: d.PreviewContentHash,
		Source: d.Source, CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func DiagramVersionToModel(orgID string, v uigraphapi.DiagramVersion) *model.DiagramVersion {
	return &model.DiagramVersion{
		ID: v.ID, OrgID: orgID, DiagramID: v.DiagramID, VersionNumber: v.VersionNumber,
		Label: v.Label, ContentKey: v.ContentKey, ContentHash: v.ContentHash,
		IsAutoVersion: v.IsAutoVersion, Source: v.Source, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func DiagramImageToModel(img uigraphapi.DiagramImage) *model.DiagramImage {
	return &model.DiagramImage{
		DiagramImageID: img.DiagramImageID, DiagramID: img.DiagramID,
		OrgID: img.OrgID, AssetID: img.AssetID, FileName: img.FileName,
		Order: img.Order, CreatedBy: img.CreatedBy, CreatedAt: img.CreatedAt,
	}
}

func DiagramsToModel(diagrams []uigraphapi.Diagram) []*model.Diagram {
	out := make([]*model.Diagram, len(diagrams))
	for i := range diagrams {
		out[i] = DiagramToModel(&diagrams[i])
	}
	return out
}

func DiagramVersionsToModel(orgID string, versions []uigraphapi.DiagramVersion) []*model.DiagramVersion {
	out := make([]*model.DiagramVersion, len(versions))
	for i, v := range versions {
		out[i] = DiagramVersionToModel(orgID, v)
	}
	return out
}

func DiagramImagesToModel(images []uigraphapi.DiagramImage) []*model.DiagramImage {
	out := make([]*model.DiagramImage, len(images))
	for i, img := range images {
		out[i] = DiagramImageToModel(img)
	}
	return out
}
