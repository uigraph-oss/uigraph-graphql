package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ── Folders ───────────────────────────────────────────────────────────────────

func (c *Client) ListFolders(ctx context.Context, orgID, folderType, parentID string) ([]Folder, error) {
	q := url.Values{}
	if folderType != "" {
		q.Set("type", folderType)
	}
	if parentID != "" {
		q.Set("parentId", parentID)
	}
	path := "/api/v1/orgs/" + orgID + "/folders"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Folders []Folder `json:"folders"`
	}
	return out.Folders, c.get(ctx, path, &out)
}

func (c *Client) GetFolder(ctx context.Context, orgID, id string) (*Folder, error) {
	var out Folder
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id), &out)
}

func (c *Client) CreateFolder(ctx context.Context, orgID string, body map[string]interface{}) (*Folder, error) {
	var out Folder
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/folders", body, &out)
}

func (c *Client) UpdateFolder(ctx context.Context, orgID, id string, body map[string]interface{}) (*Folder, error) {
	var out Folder
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id), body, &out)
}

func (c *Client) DeleteFolder(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/folders/%s", orgID, id))
}

// ── Diagrams ──────────────────────────────────────────────────────────────────

func (c *Client) ListDiagrams(ctx context.Context, orgID, folderID string) ([]Diagram, error) {
	path := "/api/v1/orgs/" + orgID + "/diagrams"
	if folderID != "" {
		path += "?folderId=" + folderID
	}
	var out struct {
		Diagrams []Diagram `json:"diagrams"`
	}
	return out.Diagrams, c.get(ctx, path, &out)
}

func (c *Client) GetDiagram(ctx context.Context, orgID, id string) (*Diagram, error) {
	var out Diagram
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s", orgID, id), &out)
}

func (c *Client) GetDiagramContent(ctx context.Context, orgID, id string) (string, error) {
	var out struct {
		DiagramID string `json:"diagramId"`
		Content   string `json:"content"`
	}
	return out.Content, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/content", orgID, id), &out)
}

func (c *Client) CreateDiagram(ctx context.Context, orgID string, body map[string]interface{}) (*Diagram, error) {
	var out Diagram
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/diagrams", body, &out)
}

func (c *Client) UpdateDiagram(ctx context.Context, orgID, id string, body map[string]interface{}) (*Diagram, error) {
	var out Diagram
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s", orgID, id), body, &out)
}

func (c *Client) DeleteDiagram(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s", orgID, id))
}

func (c *Client) UpdateDiagramThumbnail(ctx context.Context, orgID, id string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/thumbnail", orgID, id), body, &out)
}

func (c *Client) ListFlowDiagramComponents(ctx context.Context, orgID string) (*FlowComponents, error) {
	var out FlowComponents
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/flow-diagram-components", orgID), &out)
}

func (c *Client) ListDiagramImages(ctx context.Context, orgID, diagramID string) ([]DiagramImage, error) {
	var out struct {
		Images []DiagramImage `json:"images"`
	}
	return out.Images, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/images", orgID, diagramID), &out)
}

func (c *Client) CreateDiagramImage(ctx context.Context, orgID, diagramID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/images", orgID, diagramID), body, &out)
}

func (c *Client) SyncDiagram(ctx context.Context, orgID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, "/api/v1/orgs/"+orgID+"/diagrams/sync", body, &out)
}

func (c *Client) ListDiagramVersions(ctx context.Context, orgID, diagramID string) ([]DiagramVersion, error) {
	var out struct {
		Versions []DiagramVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions", orgID, diagramID), &out)
}

func (c *Client) CreateDiagramVersion(ctx context.Context, orgID, diagramID string, body map[string]interface{}) (*DiagramVersion, error) {
	var out DiagramVersion
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions", orgID, diagramID), body, &out)
}

func (c *Client) GetDiagramVersionContent(ctx context.Context, orgID, diagramID, versionID string) (string, error) {
	var out struct {
		VersionID string `json:"versionId"`
		Content   string `json:"content"`
	}
	return out.Content, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions/%s/content", orgID, diagramID, versionID), &out)
}

func (c *Client) RestoreDiagramVersion(ctx context.Context, orgID, diagramID, versionID string) (*Diagram, error) {
	var out Diagram
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/versions/%s/restore", orgID, diagramID, versionID), nil, &out)
}

// ── Maps ──────────────────────────────────────────────────────────────────────

func (c *Client) ListMaps(ctx context.Context, orgID, folderID string) ([]UIMap, error) {
	path := "/api/v1/orgs/" + orgID + "/maps"
	if folderID != "" {
		path += "?folderId=" + folderID
	}
	var out struct {
		Maps []UIMap `json:"maps"`
	}
	return out.Maps, c.get(ctx, path, &out)
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

// ── Frames ────────────────────────────────────────────────────────────────────

func (c *Client) ListFrames(ctx context.Context, orgID, mapID string) ([]Frame, error) {
	var out struct {
		Frames []Frame `json:"frames"`
	}
	return out.Frames, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames", orgID, mapID), &out)
}

func (c *Client) GetFrame(ctx context.Context, orgID, mapID, id string) (*Frame, error) {
	var out Frame
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/frames/%s", orgID, mapID, id), &out)
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

// ── Focal Points ──────────────────────────────────────────────────────────────

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

// ── Canvas ────────────────────────────────────────────────────────────────────

func (c *Client) GetCanvas(ctx context.Context, orgID, mapID string) (*Canvas, error) {
	var out Canvas
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/canvas", orgID, mapID), &out)
}

func (c *Client) UpsertCanvas(ctx context.Context, orgID, mapID string, body map[string]interface{}) (*Canvas, error) {
	var out Canvas
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/maps/%s/canvas", orgID, mapID), body, &out)
}

// rawJSON returns the JSON string of a raw message, defaulting to "{}".
func rawJSON(b json.RawMessage) string {
	if len(b) == 0 {
		return "{}"
	}
	return string(b)
}
