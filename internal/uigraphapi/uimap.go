package uigraphapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type UIMap struct {
	ID                  string     `json:"id"`
	OrgID               string     `json:"orgId"`
	FolderID            *string    `json:"folderId,omitempty"`
	TeamID              *string    `json:"teamId,omitempty"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Status              string     `json:"status"`
	CreatedBy           string     `json:"createdBy"`
	UpdatedBy           *string    `json:"updatedBy,omitempty"`
	CreatedByCommitHash *string    `json:"createdByCommitHash,omitempty"`
	UpdatedByCommitHash *string    `json:"updatedByCommitHash,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	DeletedAt           *time.Time `json:"deletedAt,omitempty"`
}

type Frame struct {
	ID                    string     `json:"id"`
	MapID                 string     `json:"mapId"`
	OrgID                 string     `json:"orgId"`
	ParentFrameID         *string    `json:"parentFrameId,omitempty"`
	Name                  string     `json:"name"`
	Description           string     `json:"description"`
	TemplateType          string     `json:"templateType"`
	ScreenshotAssetID     *string    `json:"screenshotAssetId,omitempty"`
	ScreenshotContentHash *string    `json:"screenshotContentHash,omitempty"`
	Status                string     `json:"status"`
	Order                 float64    `json:"order"`
	Source                *string    `json:"source,omitempty"`
	CreatedBy             string     `json:"createdBy"`
	UpdatedBy             *string    `json:"updatedBy,omitempty"`
	CreatedByCommitHash   *string    `json:"createdByCommitHash,omitempty"`
	UpdatedByCommitHash   *string    `json:"updatedByCommitHash,omitempty"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
	DeletedAt             *time.Time `json:"deletedAt,omitempty"`
	FocalPointCount       int        `json:"focalPointCount"`
}

type FocalPoint struct {
	ID                  string     `json:"id"`
	FrameID             string     `json:"frameId"`
	OrgID               string     `json:"orgId"`
	Name                string     `json:"name"`
	LocationX           float64    `json:"locationX"`
	LocationY           float64    `json:"locationY"`
	Visibility          string     `json:"visibility"`
	IsActive            bool       `json:"isActive"`
	CreatedBy           string     `json:"createdBy"`
	UpdatedBy           *string    `json:"updatedBy,omitempty"`
	CreatedByCommitHash *string    `json:"createdByCommitHash,omitempty"`
	UpdatedByCommitHash *string    `json:"updatedByCommitHash,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	DeletedAt           *time.Time `json:"deletedAt,omitempty"`
}

type Canvas struct {
	MapID          string          `json:"mapId"`
	OrgID          string          `json:"orgId"`
	Zoom           float64         `json:"zoom"`
	NavigationX    float64         `json:"navigationX"`
	NavigationY    float64         `json:"navigationY"`
	FramePositions json.RawMessage `json:"framePositions"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type FrameGroup struct {
	ID          string     `json:"id"`
	FrameID     string     `json:"frameId"`
	OrgID       string     `json:"orgId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	LocationX   float64    `json:"locationX"`
	LocationY   float64    `json:"locationY"`
	Width       float64    `json:"width"`
	Height      float64    `json:"height"`
	Order       float64    `json:"order"`
	IsActive    bool       `json:"isActive"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type FrameLink struct {
	ID            string     `json:"id"`
	FrameID       string     `json:"frameId"`
	OrgID         string     `json:"orgId"`
	Kind          string     `json:"kind"`
	TargetFrameID *string    `json:"targetFrameId,omitempty"`
	TargetMapID   *string    `json:"targetMapId,omitempty"`
	Label         string     `json:"label"`
	LocationX     float64    `json:"locationX"`
	LocationY     float64    `json:"locationY"`
	IsActive      bool       `json:"isActive"`
	CreatedBy     string     `json:"createdBy"`
	UpdatedBy     *string    `json:"updatedBy,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

type FocalPointMeta struct {
	ID                         string          `json:"id"`
	FocalPointID               string          `json:"focalPointId"`
	OrgID                      string          `json:"orgId"`
	FrameID                    string          `json:"frameId"`
	ComponentID                string          `json:"componentId"`
	ComponentLinkDiagramID     *string         `json:"componentLinkDiagramId,omitempty"`
	ComponentLinkAPIEndpointID *string         `json:"componentLinkApiEndpointId,omitempty"`
	ComponentLinkTestPackID    *string         `json:"componentLinkTestPackId,omitempty"`
	ComponentLinkServiceDocID  *string         `json:"componentLinkServiceDocId,omitempty"`
	ComponentModalFields       json.RawMessage `json:"componentModalFields"`
	CreatedBy                  string          `json:"createdBy"`
	UpdatedBy                  *string         `json:"updatedBy,omitempty"`
	CreatedByCommitHash        *string         `json:"createdByCommitHash,omitempty"`
	UpdatedByCommitHash        *string         `json:"updatedByCommitHash,omitempty"`
	CreatedAt                  time.Time       `json:"createdAt"`
	UpdatedAt                  time.Time       `json:"updatedAt"`
	DeletedAt                  *time.Time      `json:"deletedAt,omitempty"`
}

type ComponentLinkUsage struct {
	MetaID            string  `json:"metaId"`
	OrgID             string  `json:"orgId"`
	ComponentID       string  `json:"componentId"`
	MapID             string  `json:"mapId"`
	MapName           string  `json:"mapName"`
	FrameID           string  `json:"frameId"`
	FrameName         string  `json:"frameName"`
	ScreenshotAssetID *string `json:"screenshotAssetId,omitempty"`
	FocalPointID      string  `json:"focalPointId"`
	FocalPointName    string  `json:"focalPointName"`
	LocationX         float64 `json:"locationX"`
	LocationY         float64 `json:"locationY"`
}

func (c *Client) ListMaps(ctx context.Context, orgID string, p ListParams) ([]UIMap, int, error) {
	path := "/api/v1/orgs/" + orgID + "/maps" + listQuery(p)
	var out struct {
		Maps  []UIMap `json:"maps"`
		Total int     `json:"total"`
	}
	return out.Maps, out.Total, c.get(ctx, path, &out)
}

func (c *Client) GetMap(ctx context.Context, orgID, id string) (*UIMap, error) {
	var out UIMap
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s", orgID, id), &out)
}

func (c *Client) CreateMap(ctx context.Context, orgID string, body map[string]interface{}) (*UIMap, error) {
	var out UIMap
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/maps", body, &out)
}

func (c *Client) UpdateMap(ctx context.Context, orgID, id string, body map[string]interface{}) (*UIMap, error) {
	var out UIMap
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s", orgID, id), body, &out)
}

func (c *Client) DeleteMap(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s", orgID, id))
}

func (c *Client) ListFrames(ctx context.Context, orgID, mapID string, p ListParams) ([]Frame, int, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames", orgID, mapID) + listQuery(p)
	var out struct {
		Frames []Frame `json:"frames"`
		Total  int     `json:"total"`
	}
	return out.Frames, out.Total, c.get(ctx, path, &out)
}

func (c *Client) GetFrame(ctx context.Context, orgID, mapID, id string) (*Frame, error) {
	var out Frame
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s", orgID, mapID, id), &out)
}

func (c *Client) GetFrameByID(ctx context.Context, orgID, id string) (*Frame, error) {
	var out Frame
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/frames/%s", orgID, id), &out)
}

func (c *Client) CreateFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*Frame, error) {
	var out Frame
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames", orgID, mapID), body, &out)
}

func (c *Client) UpdateFrame(ctx context.Context, orgID, mapID, id string, body map[string]interface{}) (*Frame, error) {
	var out Frame
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s", orgID, mapID, id), body, &out)
}

func (c *Client) DeleteFrame(ctx context.Context, orgID, mapID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s", orgID, mapID, id))
}

func (c *Client) SyncFrame(ctx context.Context, orgID, mapID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/sync", orgID, mapID), body, &out)
}

func (c *Client) ListFocalPoints(ctx context.Context, orgID, mapID, frameID string) ([]FocalPoint, error) {
	var out struct {
		FocalPoints []FocalPoint `json:"focalPoints"`
	}
	return out.FocalPoints, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points", orgID, mapID, frameID), &out)
}

func (c *Client) GetFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) (*FocalPoint, error) {
	var out FocalPoint
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s", orgID, mapID, frameID, id), &out)
}

func (c *Client) CreateFocalPoint(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*FocalPoint, error) {
	var out FocalPoint
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points", orgID, mapID, frameID), body, &out)
}

func (c *Client) UpdateFocalPoint(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*FocalPoint, error) {
	var out FocalPoint
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s", orgID, mapID, frameID, id), body, &out)
}

func (c *Client) DeleteFocalPoint(ctx context.Context, orgID, mapID, frameID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s", orgID, mapID, frameID, id))
}

func (c *Client) GetCanvas(ctx context.Context, orgID, mapID string) (*Canvas, error) {
	var out Canvas
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/canvas", orgID, mapID), &out)
}

func (c *Client) UpsertCanvas(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*Canvas, error) {
	var out Canvas
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/canvas", orgID, mapID), body, &out)
}

func (c *Client) ListFrameGroups(ctx context.Context, orgID, mapID, frameID string) ([]FrameGroup, error) {
	var out struct {
		Groups []FrameGroup `json:"groups"`
	}
	return out.Groups, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups", orgID, mapID, frameID), &out)
}

func (c *Client) CreateFrameGroup(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*FrameGroup, error) {
	var out FrameGroup
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups", orgID, mapID, frameID), body, &out)
}

func (c *Client) UpdateFrameGroup(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*FrameGroup, error) {
	var out FrameGroup
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups/%s", orgID, mapID, frameID, id), body, &out)
}

func (c *Client) DeleteFrameGroup(ctx context.Context, orgID, mapID, frameID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/groups/%s", orgID, mapID, frameID, id))
}

func (c *Client) ListFrameLinks(ctx context.Context, orgID, mapID, frameID string) ([]FrameLink, error) {
	var out struct {
		Links []FrameLink `json:"links"`
	}
	return out.Links, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links", orgID, mapID, frameID), &out)
}

func (c *Client) CreateFrameLink(ctx context.Context, orgID, mapID, frameID string, body map[string]interface{}) (*FrameLink, error) {
	var out FrameLink
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links", orgID, mapID, frameID), body, &out)
}

func (c *Client) UpdateFrameLink(ctx context.Context, orgID, mapID, frameID, id string, body map[string]interface{}) (*FrameLink, error) {
	var out FrameLink
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links/%s", orgID, mapID, frameID, id), body, &out)
}

func (c *Client) DeleteFrameLink(ctx context.Context, orgID, mapID, frameID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/links/%s", orgID, mapID, frameID, id))
}

func (c *Client) ListFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string) ([]FocalPointMeta, error) {
	var out struct {
		Meta []FocalPointMeta `json:"meta"`
	}
	return out.Meta, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta", orgID, mapID, frameID, fpID), &out)
}

func (c *Client) ListFocalPointMetaByLink(ctx context.Context, orgID, linkID string) ([]FocalPointMeta, error) {
	var out struct {
		Meta []FocalPointMeta `json:"meta"`
	}
	return out.Meta, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/focal-point-meta?linkId=%s", orgID, url.QueryEscape(linkID)), &out)
}

func (c *Client) ListComponentLinkUsages(ctx context.Context, orgID, linkID string) ([]ComponentLinkUsage, error) {
	var out struct {
		Usages []ComponentLinkUsage `json:"usages"`
	}
	return out.Usages, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/component-link-usages?linkId=%s", orgID, url.QueryEscape(linkID)), &out)
}

func (c *Client) CreateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID string, body map[string]interface{}) (*FocalPointMeta, error) {
	var out FocalPointMeta
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta", orgID, mapID, frameID, fpID), body, &out)
}

func (c *Client) UpdateFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string, body map[string]interface{}) (*FocalPointMeta, error) {
	var out FocalPointMeta
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta/%s", orgID, mapID, frameID, fpID, id), body, &out)
}

func (c *Client) DeleteFocalPointMeta(ctx context.Context, orgID, mapID, frameID, fpID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s/focal-points/%s/meta/%s", orgID, mapID, frameID, fpID, id))
}
