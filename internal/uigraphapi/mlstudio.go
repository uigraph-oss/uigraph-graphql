package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type MLModel struct {
	ID                    string     `json:"id"`
	OrgID                 string     `json:"orgId"`
	Name                  string     `json:"name"`
	Description           string     `json:"description"`
	Domain                string     `json:"domain"`
	ProblemType           string     `json:"problemType"`
	Tags                  []string   `json:"tags"`
	Owners                string     `json:"owners"`
	License               string     `json:"license"`
	References            []string   `json:"references"`
	IntendedUse           string     `json:"intendedUse"`
	Limitations           string     `json:"limitations"`
	EthicalConsiderations string     `json:"ethicalConsiderations"`
	Caveats               string     `json:"caveats"`
	ProductionVersionID   *string    `json:"productionVersionId,omitempty"`
	CreatedAt             *time.Time `json:"createdAt,omitempty"`
	UpdatedAt             *time.Time `json:"updatedAt,omitempty"`
}

type MLModelVersion struct {
	ID               string     `json:"id"`
	OrgID            string     `json:"orgId"`
	ModelID          string     `json:"modelId"`
	Version          string     `json:"version"`
	Description      string     `json:"description"`
	DeploymentStatus string     `json:"deploymentStatus"`
	RunID            *string    `json:"runId,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
}

type MLVersionDeploymentUpdate struct {
	ID         string     `json:"id"`
	OrgID      string     `json:"orgId"`
	VersionID  string     `json:"versionId"`
	FromStatus *string    `json:"fromStatus,omitempty"`
	ToStatus   string     `json:"toStatus"`
	ChangedBy  string     `json:"changedBy"`
	ChangedAt  *time.Time `json:"changedAt,omitempty"`
}

type MLExperiment struct {
	ID          string     `json:"id"`
	OrgID       string     `json:"orgId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
}

type MLRun struct {
	ID           string         `json:"id"`
	OrgID        string         `json:"orgId"`
	ExperimentID string         `json:"experimentId"`
	Name         string         `json:"name"`
	Status       string         `json:"status"`
	StartedAt    *time.Time     `json:"startedAt,omitempty"`
	EndedAt      *time.Time     `json:"endedAt,omitempty"`
	Duration     string         `json:"duration"`
	Notes        string         `json:"notes"`
	Parameters   map[string]any `json:"parameters"`
	Metrics      map[string]any `json:"metrics"`
	DatasetID    *string        `json:"datasetId,omitempty"`
}

type MLMetricPoint struct {
	Key   string     `json:"key"`
	Step  int64      `json:"step"`
	Value float64    `json:"value"`
	TS    *time.Time `json:"ts,omitempty"`
}

type MLArtifact struct {
	ID     string `json:"id"`
	OrgID  string `json:"orgId"`
	RunID  string `json:"runId"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	URI    string `json:"uri"`
	Size   string `json:"size"`
	Format string `json:"format"`
}

type MLSchemaField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type MLDataset struct {
	ID           string            `json:"id"`
	OrgID        string            `json:"orgId"`
	ExperimentID string            `json:"experimentId"`
	Name         string            `json:"name"`
	Digest       string            `json:"digest"`
	Source       string            `json:"source"`
	SourceType   string            `json:"sourceType"`
	Context      string            `json:"context"`
	RowCount     int64             `json:"rowCount"`
	Schema       []MLSchemaField   `json:"schema"`
	Tags         map[string]string `json:"tags"`
}

type MLDeployment struct {
	ID           string     `json:"id"`
	OrgID        string     `json:"orgId"`
	ModelID      string     `json:"modelId"`
	VersionID    string     `json:"versionId"`
	Name         string     `json:"name"`
	Environment  string     `json:"environment"`
	Status       string     `json:"status"`
	Endpoint     string     `json:"endpoint"`
	Region       string     `json:"region"`
	DeployedAt   *time.Time `json:"deployedAt,omitempty"`
	RolledBackAt *time.Time `json:"rolledBackAt,omitempty"`
}

type MLFinding struct {
	ID          string   `json:"id"`
	OrgID       string   `json:"orgId"`
	ModelID     string   `json:"modelId"`
	VersionID   *string  `json:"versionId,omitempty"`
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	RunIDs      []string `json:"runIds"`
}

func mlBase(orgID string) string {
	return "/api/v1/orgs/" + orgID + "/ml"
}

func (c *Client) ListMLModels(ctx context.Context, orgID string) ([]MLModel, error) {
	var out struct {
		Models []MLModel `json:"models"`
	}
	return out.Models, c.get(ctx, mlBase(orgID)+"/models", &out)
}

func (c *Client) GetMLModel(ctx context.Context, orgID, id string) (*MLModel, error) {
	var out MLModel
	return &out, c.get(ctx, mlBase(orgID)+"/models/"+id, &out)
}

func (c *Client) UpdateMLModel(ctx context.Context, orgID, id string, body map[string]interface{}) (*MLModel, error) {
	var out MLModel
	return &out, c.patch(ctx, mlBase(orgID)+"/models/"+id, body, &out)
}

func (c *Client) ListMLModelVersions(ctx context.Context, orgID, modelID string) ([]MLModelVersion, error) {
	q := url.Values{}
	if modelID != "" {
		q.Set("modelId", modelID)
	}
	path := mlBase(orgID) + "/versions"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Versions []MLModelVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, path, &out)
}

func (c *Client) GetMLModelVersion(ctx context.Context, orgID, id string) (*MLModelVersion, error) {
	var out MLModelVersion
	return &out, c.get(ctx, mlBase(orgID)+"/versions/"+id, &out)
}

func (c *Client) ListMLExperiments(ctx context.Context, orgID string) ([]MLExperiment, error) {
	var out struct {
		Experiments []MLExperiment `json:"experiments"`
	}
	return out.Experiments, c.get(ctx, mlBase(orgID)+"/experiments", &out)
}

func (c *Client) GetMLExperiment(ctx context.Context, orgID, id string) (*MLExperiment, error) {
	var out MLExperiment
	return &out, c.get(ctx, mlBase(orgID)+"/experiments/"+id, &out)
}

func (c *Client) ListMLRuns(ctx context.Context, orgID, experimentID string) ([]MLRun, error) {
	q := url.Values{}
	if experimentID != "" {
		q.Set("experimentId", experimentID)
	}
	path := mlBase(orgID) + "/runs"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Runs []MLRun `json:"runs"`
	}
	return out.Runs, c.get(ctx, path, &out)
}

func (c *Client) GetMLRun(ctx context.Context, orgID, id string) (*MLRun, error) {
	var out MLRun
	return &out, c.get(ctx, mlBase(orgID)+"/runs/"+id, &out)
}

func (c *Client) ListMLRunSeries(ctx context.Context, orgID, runID string) ([]MLMetricPoint, error) {
	var out struct {
		Points []MLMetricPoint `json:"points"`
	}
	return out.Points, c.get(ctx, fmt.Sprintf("%s/runs/%s/series", mlBase(orgID), runID), &out)
}

func (c *Client) ListMLArtifacts(ctx context.Context, orgID, runID string) ([]MLArtifact, error) {
	q := url.Values{}
	if runID != "" {
		q.Set("runId", runID)
	}
	path := mlBase(orgID) + "/artifacts"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Artifacts []MLArtifact `json:"artifacts"`
	}
	return out.Artifacts, c.get(ctx, path, &out)
}

func (c *Client) ListMLDatasets(ctx context.Context, orgID, experimentID string) ([]MLDataset, error) {
	q := url.Values{}
	if experimentID != "" {
		q.Set("experimentId", experimentID)
	}
	path := mlBase(orgID) + "/datasets"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Datasets []MLDataset `json:"datasets"`
	}
	return out.Datasets, c.get(ctx, path, &out)
}

func (c *Client) GetMLDataset(ctx context.Context, orgID, id string) (*MLDataset, error) {
	var out MLDataset
	return &out, c.get(ctx, mlBase(orgID)+"/datasets/"+id, &out)
}

func (c *Client) ListMLDeployments(ctx context.Context, orgID, modelID, versionID string) ([]MLDeployment, error) {
	q := url.Values{}
	if modelID != "" {
		q.Set("modelId", modelID)
	}
	if versionID != "" {
		q.Set("versionId", versionID)
	}
	path := mlBase(orgID) + "/deployments"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Deployments []MLDeployment `json:"deployments"`
	}
	return out.Deployments, c.get(ctx, path, &out)
}

func (c *Client) GetMLDeployment(ctx context.Context, orgID, id string) (*MLDeployment, error) {
	var out MLDeployment
	return &out, c.get(ctx, mlBase(orgID)+"/deployments/"+id, &out)
}

func (c *Client) CreateMLDeployment(ctx context.Context, orgID string, body map[string]interface{}) (*MLDeployment, error) {
	var out MLDeployment
	return &out, c.post(ctx, mlBase(orgID)+"/deployments", body, &out)
}

func (c *Client) UpdateMLDeployment(ctx context.Context, orgID, id string, body map[string]interface{}) (*MLDeployment, error) {
	var out MLDeployment
	return &out, c.put(ctx, mlBase(orgID)+"/deployments/"+id, body, &out)
}

func (c *Client) DeleteMLDeployment(ctx context.Context, orgID, id string) error {
	return c.del(ctx, mlBase(orgID)+"/deployments/"+id)
}

func (c *Client) ListVersionDeploymentUpdates(ctx context.Context, orgID, versionID string) ([]MLVersionDeploymentUpdate, error) {
	var out struct {
		Updates []MLVersionDeploymentUpdate `json:"updates"`
	}
	return out.Updates, c.get(ctx, fmt.Sprintf("%s/versions/%s/deployment-updates", mlBase(orgID), versionID), &out)
}

func (c *Client) CreateVersionDeploymentUpdate(ctx context.Context, orgID, versionID string, body map[string]interface{}) (*MLVersionDeploymentUpdate, error) {
	var out MLVersionDeploymentUpdate
	return &out, c.post(ctx, fmt.Sprintf("%s/versions/%s/deployment-updates", mlBase(orgID), versionID), body, &out)
}

func (c *Client) ListMLFindings(ctx context.Context, orgID, modelID string) ([]MLFinding, error) {
	q := url.Values{}
	if modelID != "" {
		q.Set("modelId", modelID)
	}
	path := mlBase(orgID) + "/findings"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		Findings []MLFinding `json:"findings"`
	}
	return out.Findings, c.get(ctx, path, &out)
}

func (c *Client) GetMLFinding(ctx context.Context, orgID, id string) (*MLFinding, error) {
	var out MLFinding
	return &out, c.get(ctx, mlBase(orgID)+"/findings/"+id, &out)
}

func (c *Client) CreateMLFinding(ctx context.Context, orgID string, body map[string]interface{}) (*MLFinding, error) {
	var out MLFinding
	return &out, c.post(ctx, mlBase(orgID)+"/findings", body, &out)
}

func (c *Client) UpdateMLFinding(ctx context.Context, orgID, id string, body map[string]interface{}) (*MLFinding, error) {
	var out MLFinding
	return &out, c.put(ctx, mlBase(orgID)+"/findings/"+id, body, &out)
}

func (c *Client) DeleteMLFinding(ctx context.Context, orgID, id string) error {
	return c.del(ctx, mlBase(orgID)+"/findings/"+id)
}
