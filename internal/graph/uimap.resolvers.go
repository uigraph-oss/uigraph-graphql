package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func (r *componentLinkUsageResolver) ScreenshotImageURL(ctx context.Context, obj *model.ComponentLinkUsage) (*string, error) {
	if obj.ScreenshotAssetID == nil {
		return nil, nil
	}
	return r.resolveAssetURL(ctx, obj.OrgID, *obj.ScreenshotAssetID)
}

func (r *frameResolver) ScreenshotImageURL(ctx context.Context, obj *model.Frame) (*string, error) {
	if obj.ScreenshotAssetID == nil {
		return nil, nil
	}
	return r.resolveAssetURL(ctx, obj.OrgID, *obj.ScreenshotAssetID)
}

func (r *frameResolver) CreatedByActor(ctx context.Context, obj *model.Frame) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *frameResolver) UpdatedByActor(ctx context.Context, obj *model.Frame) (*model.Actor, error) {
	if obj.UpdatedBy == nil {
		return nil, nil
	}
	return r.resolveActor(ctx, obj.OrgID, *obj.UpdatedBy)
}

func (r *mutationResolver) CreateMap(ctx context.Context, orgID string, input model.CreateMapInput) (*model.UIMap, error) {
	m, err := r.UIMapAPI.CreateMap(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.UIMapToModel(m), nil
}

func (r *mutationResolver) UpdateMap(ctx context.Context, orgID string, id string, input model.UpdateMapInput) (*model.UIMap, error) {
	m, err := r.UIMapAPI.UpdateMap(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.UIMapToModel(m), nil
}

func (r *mutationResolver) DeleteMap(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.UIMapAPI.DeleteMap(ctx, orgID, id)
}

func (r *mutationResolver) CreateFrame(ctx context.Context, orgID string, mapID string, input model.CreateFrameInput) (*model.Frame, error) {
	f, err := r.UIMapAPI.CreateFrame(ctx, orgID, mapID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FrameToModel(f), nil
}

func (r *mutationResolver) UpdateFrame(ctx context.Context, orgID string, mapID string, id string, input model.UpdateFrameInput) (*model.Frame, error) {
	f, err := r.UIMapAPI.UpdateFrame(ctx, orgID, mapID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FrameToModel(f), nil
}

func (r *mutationResolver) DeleteFrame(ctx context.Context, orgID string, mapID string, id string) (bool, error) {
	return true, r.UIMapAPI.DeleteFrame(ctx, orgID, mapID, id)
}

func (r *mutationResolver) SyncFrame(ctx context.Context, orgID string, mapID string, input model.SyncFrameInput) (*model.SyncFrameResult, error) {
	out, err := r.UIMapAPI.SyncFrame(ctx, orgID, mapID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return &model.SyncFrameResult{
		FrameID:        convert.StrFromMap(out, "frameId"),
		VersionCreated: convert.BoolFromMap(out, "versionCreated"),
	}, nil
}

func (r *mutationResolver) CreateFocalPoint(ctx context.Context, orgID string, mapID string, frameID string, input model.CreateFocalPointInput) (*model.FocalPoint, error) {
	fp, err := r.UIMapAPI.CreateFocalPoint(ctx, orgID, mapID, frameID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FocalPointToModel(fp), nil
}

func (r *mutationResolver) UpdateFocalPoint(ctx context.Context, orgID string, mapID string, frameID string, id string, input model.UpdateFocalPointInput) (*model.FocalPoint, error) {
	fp, err := r.UIMapAPI.UpdateFocalPoint(ctx, orgID, mapID, frameID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FocalPointToModel(fp), nil
}

func (r *mutationResolver) DeleteFocalPoint(ctx context.Context, orgID string, mapID string, frameID string, id string) (bool, error) {
	return true, r.UIMapAPI.DeleteFocalPoint(ctx, orgID, mapID, frameID, id)
}

func (r *mutationResolver) UpsertCanvas(ctx context.Context, orgID string, mapID string, input model.UpsertCanvasInput) (*model.Canvas, error) {
	body := convert.ToMap(input)
	if fp, ok := body["framePositions"].(string); ok {
		var raw interface{}
		if err := convert.UnmarshalJSONString(fp, &raw); err == nil {
			body["framePositions"] = raw
		}
	}
	c, err := r.UIMapAPI.UpsertCanvas(ctx, orgID, mapID, body)
	if err != nil {
		return nil, err
	}
	return convert.CanvasToModel(c), nil
}

func (r *mutationResolver) CreateFrameGroup(ctx context.Context, orgID string, mapID string, frameID string, input model.CreateFrameGroupInput) (*model.FrameGroup, error) {
	g, err := r.UIMapAPI.CreateFrameGroup(ctx, orgID, mapID, frameID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FrameGroupToModel(g), nil
}

func (r *mutationResolver) UpdateFrameGroup(ctx context.Context, orgID string, mapID string, frameID string, id string, input model.UpdateFrameGroupInput) (*model.FrameGroup, error) {
	g, err := r.UIMapAPI.UpdateFrameGroup(ctx, orgID, mapID, frameID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FrameGroupToModel(g), nil
}

func (r *mutationResolver) DeleteFrameGroup(ctx context.Context, orgID string, mapID string, frameID string, id string) (bool, error) {
	return true, r.UIMapAPI.DeleteFrameGroup(ctx, orgID, mapID, frameID, id)
}

func (r *mutationResolver) CreateFrameLink(ctx context.Context, orgID string, mapID string, frameID string, input model.CreateFrameLinkInput) (*model.FrameLink, error) {
	l, err := r.UIMapAPI.CreateFrameLink(ctx, orgID, mapID, frameID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FrameLinkToModel(l), nil
}

func (r *mutationResolver) UpdateFrameLink(ctx context.Context, orgID string, mapID string, frameID string, id string, input model.UpdateFrameLinkInput) (*model.FrameLink, error) {
	l, err := r.UIMapAPI.UpdateFrameLink(ctx, orgID, mapID, frameID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FrameLinkToModel(l), nil
}

func (r *mutationResolver) DeleteFrameLink(ctx context.Context, orgID string, mapID string, frameID string, id string) (bool, error) {
	return true, r.UIMapAPI.DeleteFrameLink(ctx, orgID, mapID, frameID, id)
}

func (r *mutationResolver) CreateFocalPointMeta(ctx context.Context, orgID string, mapID string, frameID string, focalPointID string, input model.CreateFocalPointMetaInput) (*model.FocalPointMeta, error) {
	m, err := r.UIMapAPI.CreateFocalPointMeta(ctx, orgID, mapID, frameID, focalPointID, convert.FocalPointMetaBody(convert.ToMap(input)))
	if err != nil {
		return nil, err
	}
	return convert.FocalPointMetaToModel(m), nil
}

func (r *mutationResolver) UpdateFocalPointMeta(ctx context.Context, orgID string, mapID string, frameID string, focalPointID string, id string, input model.UpdateFocalPointMetaInput) (*model.FocalPointMeta, error) {
	m, err := r.UIMapAPI.UpdateFocalPointMeta(ctx, orgID, mapID, frameID, focalPointID, id, convert.FocalPointMetaBody(convert.ToMap(input)))
	if err != nil {
		return nil, err
	}
	return convert.FocalPointMetaToModel(m), nil
}

func (r *mutationResolver) DeleteFocalPointMeta(ctx context.Context, orgID string, mapID string, frameID string, focalPointID string, id string) (bool, error) {
	return true, r.UIMapAPI.DeleteFocalPointMeta(ctx, orgID, mapID, frameID, focalPointID, id)
}

func (r *queryResolver) Maps(ctx context.Context, orgID string, folderID *string, teamID *string, search *string, sortBy *string, sortDir *string, limit *int, offset *int) (*model.UIMapPage, error) {
	p := uigraphapi.ListParams{
		FolderID: derefStr(folderID),
		TeamID:   derefStr(teamID),
		Search:   derefStr(search),
		SortBy:   derefStr(sortBy),
		SortDir:  derefStr(sortDir),
		Limit:    limit,
		Offset:   offset,
	}
	maps, total, err := r.UIMapAPI.ListMaps(ctx, orgID, p)
	if err != nil {
		return nil, err
	}
	return &model.UIMapPage{Items: convert.UIMapsToModel(maps), TotalCount: total}, nil
}

func (r *queryResolver) Map(ctx context.Context, orgID string, id string) (*model.UIMap, error) {
	m, err := r.UIMapAPI.GetMap(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.UIMapToModel(m), nil
}

func (r *queryResolver) Frames(ctx context.Context, orgID string, mapID string, search *string, sortBy *string, sortDir *string, limit *int, offset *int) (*model.FramePage, error) {
	p := uigraphapi.ListParams{
		Search:  derefStr(search),
		SortBy:  derefStr(sortBy),
		SortDir: derefStr(sortDir),
		Limit:   limit,
		Offset:  offset,
	}
	frames, total, err := r.UIMapAPI.ListFrames(ctx, orgID, mapID, p)
	if err != nil {
		return nil, err
	}
	return &model.FramePage{Items: convert.FramesToModel(frames), TotalCount: total}, nil
}

func (r *queryResolver) Frame(ctx context.Context, orgID string, mapID string, id string) (*model.Frame, error) {
	f, err := r.UIMapAPI.GetFrame(ctx, orgID, mapID, id)
	if err != nil {
		return nil, err
	}
	return convert.FrameToModel(f), nil
}

func (r *queryResolver) FrameByID(ctx context.Context, orgID string, id string) (*model.Frame, error) {
	f, err := r.UIMapAPI.GetFrameByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.FrameToModel(f), nil
}

func (r *queryResolver) FocalPoints(ctx context.Context, orgID string, mapID string, frameID string) ([]*model.FocalPoint, error) {
	fps, err := r.UIMapAPI.ListFocalPoints(ctx, orgID, mapID, frameID)
	if err != nil {
		return nil, err
	}
	return convert.FocalPointsToModel(fps), nil
}

func (r *queryResolver) Canvas(ctx context.Context, orgID string, mapID string) (*model.Canvas, error) {
	c, err := r.UIMapAPI.GetCanvas(ctx, orgID, mapID)
	if err != nil {
		return nil, err
	}
	return convert.CanvasToModel(c), nil
}

func (r *queryResolver) FrameGroups(ctx context.Context, orgID string, mapID string, frameID string) ([]*model.FrameGroup, error) {
	groups, err := r.UIMapAPI.ListFrameGroups(ctx, orgID, mapID, frameID)
	if err != nil {
		return nil, err
	}
	return convert.FrameGroupsToModel(groups), nil
}

func (r *queryResolver) FrameLinks(ctx context.Context, orgID string, mapID string, frameID string) ([]*model.FrameLink, error) {
	links, err := r.UIMapAPI.ListFrameLinks(ctx, orgID, mapID, frameID)
	if err != nil {
		return nil, err
	}
	return convert.FrameLinksToModel(links), nil
}

func (r *queryResolver) FocalPointMeta(ctx context.Context, orgID string, mapID string, frameID string, focalPointID string) ([]*model.FocalPointMeta, error) {
	metas, err := r.UIMapAPI.ListFocalPointMeta(ctx, orgID, mapID, frameID, focalPointID)
	if err != nil {
		return nil, err
	}
	return convert.FocalPointMetasToModel(metas), nil
}

func (r *queryResolver) FocalPointMetaByLink(ctx context.Context, orgID string, linkID string) ([]*model.FocalPointMeta, error) {
	metas, err := r.UIMapAPI.ListFocalPointMetaByLink(ctx, orgID, linkID)
	if err != nil {
		return nil, err
	}
	return convert.FocalPointMetasToModel(metas), nil
}

func (r *queryResolver) ComponentLinkUsages(ctx context.Context, orgID string, linkID string) ([]*model.ComponentLinkUsage, error) {
	usages, err := r.UIMapAPI.ListComponentLinkUsages(ctx, orgID, linkID)
	if err != nil {
		return nil, err
	}
	return convert.ComponentLinkUsagesToModel(usages), nil
}

func (r *uIMapResolver) PreviewImgUrls(ctx context.Context, obj *model.UIMap) ([]string, error) {
	frames, _, err := r.UIMapAPI.ListFrames(ctx, obj.OrgID, obj.ID, uigraphapi.ListParams{})
	if err != nil {
		return nil, err
	}

	assetIDs := make([]string, 0, len(frames))
	for _, f := range frames {
		if f.ScreenshotAssetID != nil && *f.ScreenshotAssetID != "" {
			assetIDs = append(assetIDs, *f.ScreenshotAssetID)
		}
	}
	if len(assetIDs) == 0 {
		return []string{}, nil
	}

	urlsByID, err := r.Actor.ResolveAssetURLs(ctx, obj.OrgID, assetIDs)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(assetIDs))
	for _, id := range assetIDs {
		if u, ok := urlsByID[id]; ok && u != "" {
			urls = append(urls, u)
		}
	}
	return urls, nil
}

func (r *Resolver) ComponentLinkUsage() generated.ComponentLinkUsageResolver {
	return &componentLinkUsageResolver{r}
}

func (r *Resolver) Frame() generated.FrameResolver { return &frameResolver{r} }

func (r *Resolver) UIMap() generated.UIMapResolver { return &uIMapResolver{r} }

type componentLinkUsageResolver struct{ *Resolver }
type frameResolver struct{ *Resolver }
type uIMapResolver struct{ *Resolver }
