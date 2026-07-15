package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func (r *diagramResolver) PreviewImageURL(ctx context.Context, obj *model.Diagram) (*string, error) {
	if obj.PreviewAssetID == nil {
		return nil, nil
	}
	return r.resolveAssetURL(ctx, obj.OrgID, *obj.PreviewAssetID)
}

func (r *diagramResolver) CreatedByActor(ctx context.Context, obj *model.Diagram) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *diagramResolver) UpdatedByActor(ctx context.Context, obj *model.Diagram) (*model.Actor, error) {
	if obj.UpdatedBy == nil {
		return nil, nil
	}
	return r.resolveActor(ctx, obj.OrgID, *obj.UpdatedBy)
}

func (r *diagramImageResolver) ImageURL(ctx context.Context, obj *model.DiagramImage) (*string, error) {
	return r.resolveAssetURL(ctx, obj.OrgID, obj.AssetID)
}

func (r *diagramVersionResolver) CreatedByActor(ctx context.Context, obj *model.DiagramVersion) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *mutationResolver) CreateDiagram(ctx context.Context, orgID string, input model.CreateDiagramInput) (*model.Diagram, error) {
	d, err := r.DiagramAPI.CreateDiagram(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.DiagramToModel(d), nil
}

func (r *mutationResolver) UpdateDiagram(ctx context.Context, orgID string, id string, input model.UpdateDiagramInput) (*model.Diagram, error) {
	d, err := r.DiagramAPI.UpdateDiagram(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.DiagramToModel(d), nil
}

func (r *mutationResolver) DeleteDiagram(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.DiagramAPI.DeleteDiagram(ctx, orgID, id)
}

func (r *mutationResolver) SyncDiagram(ctx context.Context, orgID string, input model.SyncDiagramInput) (*model.SyncDiagramResult, error) {
	out, err := r.DiagramAPI.SyncDiagram(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	res := &model.SyncDiagramResult{
		DiagramID:      convert.StrFromMap(out, "diagramId"),
		VersionCreated: convert.BoolFromMap(out, "versionCreated"),
	}
	if v := convert.OptStrFromMap(out, "versionId"); v != nil {
		res.VersionID = v
	}
	return res, nil
}

func (r *mutationResolver) CreateDiagramVersion(ctx context.Context, orgID string, diagramID string, label *string) (*model.DiagramVersion, error) {
	body := map[string]interface{}{}
	if label != nil {
		body["label"] = *label
	}
	v, err := r.DiagramAPI.CreateDiagramVersion(ctx, orgID, diagramID, body)
	if err != nil {
		return nil, err
	}
	return convert.DiagramVersionToModel(orgID, *v), nil
}

func (r *mutationResolver) RestoreDiagramVersion(ctx context.Context, orgID string, diagramID string, versionID string) (*model.Diagram, error) {
	d, err := r.DiagramAPI.RestoreDiagramVersion(ctx, orgID, diagramID, versionID)
	if err != nil {
		return nil, err
	}
	return convert.DiagramToModel(d), nil
}

func (r *mutationResolver) PrepareDiagramThumbnailUpload(ctx context.Context, orgID string, diagramID string) (*model.DiagramThumbnailUpload, error) {
	out, err := r.DiagramAPI.PrepareDiagramThumbnailUpload(ctx, orgID, diagramID)
	if err != nil {
		return nil, err
	}
	return &model.DiagramThumbnailUpload{UploadURL: out.UploadURL, AssetID: out.AssetID}, nil
}

func (r *mutationResolver) ConfirmDiagramThumbnailUpload(ctx context.Context, orgID string, diagramID string, contentHash string) (bool, error) {
	return true, r.DiagramAPI.ConfirmDiagramThumbnailUpload(ctx, orgID, diagramID, contentHash)
}

func (r *mutationResolver) CreateDiagramImage(ctx context.Context, orgID string, diagramID string, input model.CreateDiagramImageInput) (*model.DiagramImage, error) {
	img, err := r.DiagramAPI.CreateDiagramImage(ctx, orgID, diagramID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.DiagramImageToModel(*img), nil
}

func (r *queryResolver) Diagrams(ctx context.Context, orgID string, folderID *string, teamID *string, serviceID *string, search *string, sortBy *string, sortDir *string, limit *int, offset *int) (*model.DiagramPage, error) {
	p := uigraphapi.ListParams{
		FolderID:  derefStr(folderID),
		TeamID:    derefStr(teamID),
		ServiceID: derefStr(serviceID),
		Search:    derefStr(search),
		SortBy:    derefStr(sortBy),
		SortDir:   derefStr(sortDir),
		Limit:     limit,
		Offset:    offset,
	}
	diagrams, total, err := r.DiagramAPI.ListDiagrams(ctx, orgID, p)
	if err != nil {
		return nil, err
	}
	return &model.DiagramPage{Items: convert.DiagramsToModel(diagrams), TotalCount: total}, nil
}

func (r *queryResolver) Diagram(ctx context.Context, orgID string, id string) (*model.Diagram, error) {
	d, err := r.DiagramAPI.GetDiagram(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.DiagramToModel(d), nil
}

func (r *queryResolver) DiagramContent(ctx context.Context, orgID string, id string) (*model.DiagramContent, error) {
	content, err := r.DiagramAPI.GetDiagramContent(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return &model.DiagramContent{DiagramID: id, Content: content}, nil
}

func (r *queryResolver) DiagramVersions(ctx context.Context, orgID string, diagramID string) ([]*model.DiagramVersion, error) {
	versions, err := r.DiagramAPI.ListDiagramVersions(ctx, orgID, diagramID)
	if err != nil {
		return nil, err
	}
	return convert.DiagramVersionsToModel(orgID, versions), nil
}

func (r *queryResolver) DiagramVersionContent(ctx context.Context, orgID string, diagramID string, versionID string) (*model.DiagramContent, error) {
	content, err := r.DiagramAPI.GetDiagramVersionContent(ctx, orgID, diagramID, versionID)
	if err != nil {
		return nil, err
	}
	return &model.DiagramContent{DiagramID: diagramID, Content: content}, nil
}

func (r *queryResolver) DiagramImages(ctx context.Context, orgID string, diagramID string) ([]*model.DiagramImage, error) {
	images, err := r.DiagramAPI.ListDiagramImages(ctx, orgID, diagramID)
	if err != nil {
		return nil, err
	}
	return convert.DiagramImagesToModel(images), nil
}

func (r *Resolver) Diagram() generated.DiagramResolver { return &diagramResolver{r} }

func (r *Resolver) DiagramImage() generated.DiagramImageResolver { return &diagramImageResolver{r} }

func (r *Resolver) DiagramVersion() generated.DiagramVersionResolver {
	return &diagramVersionResolver{r}
}

type diagramResolver struct{ *Resolver }
type diagramImageResolver struct{ *Resolver }
type diagramVersionResolver struct{ *Resolver }
