package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func UIMapToModel(m *uigraphapi.UIMap) *model.UIMap {
	return &model.UIMap{
		ID: m.ID, OrgID: m.OrgID, FolderID: m.FolderID, TeamID: m.TeamID,
		Name: m.Name, Description: m.Description, Status: m.Status,
		CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func FrameToModel(f *uigraphapi.Frame) *model.Frame {
	return &model.Frame{
		ID: f.ID, MapID: f.MapID, OrgID: f.OrgID, ParentFrameID: f.ParentFrameID,
		Name: f.Name, Description: f.Description, TemplateType: f.TemplateType,
		ScreenshotAssetID: f.ScreenshotAssetID, ScreenshotContentHash: f.ScreenshotContentHash,
		Status: f.Status, Order: f.Order, Source: f.Source,
		CreatedBy: f.CreatedBy, UpdatedBy: f.UpdatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func FocalPointToModel(fp *uigraphapi.FocalPoint) *model.FocalPoint {
	return &model.FocalPoint{
		ID: fp.ID, FrameID: fp.FrameID, OrgID: fp.OrgID,
		Name: fp.Name, LocationX: fp.LocationX, LocationY: fp.LocationY,
		Visibility: fp.Visibility, IsActive: fp.IsActive,
		CreatedBy: fp.CreatedBy, UpdatedBy: fp.UpdatedBy, CreatedAt: fp.CreatedAt, UpdatedAt: fp.UpdatedAt,
	}
}

func CanvasToModel(c *uigraphapi.Canvas) *model.Canvas {
	return &model.Canvas{
		MapID: c.MapID, OrgID: c.OrgID,
		Zoom: c.Zoom, NavigationX: c.NavigationX, NavigationY: c.NavigationY,
		FramePositions: RawStr(c.FramePositions),
		UpdatedAt:      c.UpdatedAt,
	}
}

func FrameGroupToModel(g *uigraphapi.FrameGroup) *model.FrameGroup {
	return &model.FrameGroup{
		ID: g.ID, FrameID: g.FrameID, OrgID: g.OrgID,
		Name: g.Name, Description: g.Description,
		LocationX: g.LocationX, LocationY: g.LocationY,
		Width: g.Width, Height: g.Height, Order: g.Order, IsActive: g.IsActive,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy,
		CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func FrameGroupsToModel(gs []uigraphapi.FrameGroup) []*model.FrameGroup {
	out := make([]*model.FrameGroup, len(gs))
	for i := range gs {
		out[i] = FrameGroupToModel(&gs[i])
	}
	return out
}

func FrameLinkToModel(l *uigraphapi.FrameLink) *model.FrameLink {
	return &model.FrameLink{
		ID: l.ID, FrameID: l.FrameID, OrgID: l.OrgID, Kind: l.Kind,
		TargetFrameID: l.TargetFrameID, TargetMapID: l.TargetMapID,
		Label: l.Label, LocationX: l.LocationX, LocationY: l.LocationY, IsActive: l.IsActive,
		CreatedBy: l.CreatedBy, UpdatedBy: l.UpdatedBy,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func FrameLinksToModel(ls []uigraphapi.FrameLink) []*model.FrameLink {
	out := make([]*model.FrameLink, len(ls))
	for i := range ls {
		out[i] = FrameLinkToModel(&ls[i])
	}
	return out
}

func FocalPointMetaToModel(m *uigraphapi.FocalPointMeta) *model.FocalPointMeta {
	return &model.FocalPointMeta{
		ID: m.ID, FocalPointID: m.FocalPointID, OrgID: m.OrgID, FrameID: m.FrameID,
		ComponentID: m.ComponentID, ComponentLinkID: m.ComponentLinkID,
		ComponentImages:      RawArrStr(m.ComponentImages),
		ComponentFlowDiagram: m.ComponentFlowDiagram,
		ComponentModalFields: RawArrStr(m.ComponentModalFields),
		CreatedBy:            m.CreatedBy,
		UpdatedBy:            m.UpdatedBy,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func FocalPointMetasToModel(ms []uigraphapi.FocalPointMeta) []*model.FocalPointMeta {
	out := make([]*model.FocalPointMeta, len(ms))
	for i := range ms {
		out[i] = FocalPointMetaToModel(&ms[i])
	}
	return out
}

func FocalPointMetaBody(body map[string]interface{}) map[string]interface{} {
	for _, key := range []string{"componentImages", "componentModalFields"} {
		if s, ok := body[key].(string); ok {
			var raw interface{}
			if err := UnmarshalJSONString(s, &raw); err == nil {
				body[key] = raw
			}
		}
	}
	return body
}

func UIMapsToModel(maps []uigraphapi.UIMap) []*model.UIMap {
	out := make([]*model.UIMap, len(maps))
	for i := range maps {
		out[i] = UIMapToModel(&maps[i])
	}
	return out
}

func FramesToModel(frames []uigraphapi.Frame) []*model.Frame {
	out := make([]*model.Frame, len(frames))
	for i := range frames {
		out[i] = FrameToModel(&frames[i])
	}
	return out
}

func FocalPointsToModel(fps []uigraphapi.FocalPoint) []*model.FocalPoint {
	out := make([]*model.FocalPoint, len(fps))
	for i := range fps {
		out[i] = FocalPointToModel(&fps[i])
	}
	return out
}
