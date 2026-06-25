package uigraphapi

import (
	"context"
	"fmt"
	"time"
)

type Diagram struct {
	ID                 string     `json:"id"`
	OrgID              string     `json:"orgId"`
	FolderID           *string    `json:"folderId,omitempty"`
	TeamID             *string    `json:"teamId,omitempty"`
	Name               string     `json:"name"`
	ContentKey         string     `json:"contentKey"`
	ContentHash        string     `json:"contentHash"`
	PreviewAssetID     *string    `json:"previewAssetId,omitempty"`
	PreviewContentHash *string    `json:"previewContentHash,omitempty"`
	Source             *string    `json:"source,omitempty"`
	CreatedBy          string     `json:"createdBy"`
	UpdatedBy          *string    `json:"updatedBy,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt,omitempty"`
}

type DiagramImage struct {
	DiagramImageID string    `json:"diagramImageId"`
	DiagramID      string    `json:"diagramId"`
	OrgID          string    `json:"orgId"`
	AssetID        string    `json:"assetId"`
	FileName       *string   `json:"fileName,omitempty"`
	Order          int       `json:"order"`
	CreatedBy      string    `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
}

type DiagramVersion struct {
	ID            string    `json:"id"`
	DiagramID     string    `json:"diagramId"`
	VersionNumber int       `json:"versionNumber"`
	Label         *string   `json:"label,omitempty"`
	ContentKey    string    `json:"contentKey"`
	ContentHash   string    `json:"contentHash"`
	IsAutoVersion bool      `json:"isAutoVersion"`
	Source        *string   `json:"source,omitempty"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

type ListParams struct {
	FolderID string
	TeamID   string
	Search   string
	SortBy   string
	SortDir  string
	Limit    *int
	Offset   *int
}

func (c *Client) ListDiagrams(ctx context.Context, orgID string, p ListParams) ([]Diagram, int, error) {
	path := "/api/v1/orgs/" + orgID + "/diagrams" + listQuery(p)
	var out struct {
		Diagrams []Diagram `json:"diagrams"`
		Total    int       `json:"total"`
	}
	return out.Diagrams, out.Total, c.get(ctx, path, &out)
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

func (c *Client) ListDiagramImages(ctx context.Context, orgID, diagramID string) ([]DiagramImage, error) {
	var out struct {
		Images []DiagramImage `json:"images"`
	}
	return out.Images, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/images", orgID, diagramID), &out)
}

func (c *Client) CreateDiagramImage(ctx context.Context, orgID, diagramID string, body map[string]interface{}) (*DiagramImage, error) {
	var out DiagramImage
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/images", orgID, diagramID), body, &out)
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

// DiagramThumbnailUpload is the response from PrepareDiagramThumbnailUpload.
type DiagramThumbnailUpload struct {
	UploadURL string `json:"uploadUrl"`
	AssetID   string `json:"assetId"`
}

// PrepareDiagramThumbnailUpload calls POST /thumbnail/prepare and returns a
// presigned PUT URL plus the deterministic asset ID for the diagram thumbnail.
func (c *Client) PrepareDiagramThumbnailUpload(ctx context.Context, orgID, diagramID string) (*DiagramThumbnailUpload, error) {
	var out DiagramThumbnailUpload
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/thumbnail/prepare", orgID, diagramID), nil, &out)
}

// ConfirmDiagramThumbnailUpload calls POST /thumbnail/confirm after the client
// has uploaded the file directly to storage via the presigned PUT URL.
func (c *Client) ConfirmDiagramThumbnailUpload(ctx context.Context, orgID, diagramID, contentHash string) error {
	return c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/diagrams/%s/thumbnail/confirm", orgID, diagramID),
		map[string]any{"contentHash": contentHash}, nil)
}
