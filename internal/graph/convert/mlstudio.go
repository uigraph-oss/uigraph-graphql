package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func MLProjectToModel(p *uigraphapi.MLProject) *model.MlProject {
	out := &model.MlProject{
		ID: p.ID, Name: p.Name, Type: p.Type, Description: p.Description,
		SourceType: p.SourceType, SourceURL: p.SourceURL, TeamID: p.TeamID,
	}
	if p.Stats != nil {
		out.Stats = &model.MlProjectStats{
			ModelCount:      p.Stats.ModelCount,
			ExperimentCount: p.Stats.ExperimentCount,
			RunCount:        p.Stats.RunCount,
		}
	}
	return out
}

func MLProjectsToModel(in []uigraphapi.MLProject) []*model.MlProject {
	out := make([]*model.MlProject, len(in))
	for i := range in {
		out[i] = MLProjectToModel(&in[i])
	}
	return out
}

func MLModelToModel(m *uigraphapi.MLModel) *model.MlModel {
	return &model.MlModel{
		ID: m.ID, ProjectID: m.ProjectID, Name: m.Name, Description: m.Description,
		Domain: m.Domain, ProblemType: m.ProblemType, Tags: m.Tags,
		Owners: m.Owners, License: m.License, References: m.References,
		IntendedUse: m.IntendedUse, Limitations: m.Limitations,
		EthicalConsiderations: m.EthicalConsiderations, Caveats: m.Caveats,
		ProductionVersionID: m.ProductionVersionID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func MLModelsToModel(in []uigraphapi.MLModel) []*model.MlModel {
	out := make([]*model.MlModel, len(in))
	for i := range in {
		out[i] = MLModelToModel(&in[i])
	}
	return out
}

func MLModelVersionToModel(v *uigraphapi.MLModelVersion) *model.MlModelVersion {
	return &model.MlModelVersion{
		ID: v.ID, ModelID: v.ModelID, Version: v.Version,
		Description: v.Description, DeploymentStatus: v.DeploymentStatus, RunID: v.RunID,
		CreatedAt: v.CreatedAt,
	}
}

func MLModelVersionsToModel(in []uigraphapi.MLModelVersion) []*model.MlModelVersion {
	out := make([]*model.MlModelVersion, len(in))
	for i := range in {
		out[i] = MLModelVersionToModel(&in[i])
	}
	return out
}

func MLVersionDeploymentUpdateToModel(u *uigraphapi.MLVersionDeploymentUpdate) *model.MlVersionDeploymentUpdate {
	return &model.MlVersionDeploymentUpdate{
		ID: u.ID, VersionID: u.VersionID, FromStatus: u.FromStatus,
		ToStatus: u.ToStatus, ChangedBy: u.ChangedBy, ChangedAt: u.ChangedAt,
	}
}

func MLVersionDeploymentUpdatesToModel(in []uigraphapi.MLVersionDeploymentUpdate) []*model.MlVersionDeploymentUpdate {
	out := make([]*model.MlVersionDeploymentUpdate, len(in))
	for i := range in {
		out[i] = MLVersionDeploymentUpdateToModel(&in[i])
	}
	return out
}

func MLExperimentToModel(e *uigraphapi.MLExperiment) *model.MlExperiment {
	return &model.MlExperiment{
		ID: e.ID, ProjectID: e.ProjectID, Name: e.Name, Description: e.Description, Status: e.Status,
		StartedAt: e.StartedAt,
	}
}

func MLExperimentsToModel(in []uigraphapi.MLExperiment) []*model.MlExperiment {
	out := make([]*model.MlExperiment, len(in))
	for i := range in {
		out[i] = MLExperimentToModel(&in[i])
	}
	return out
}

func MLRunToModel(run *uigraphapi.MLRun) *model.MlRun {
	params := map[string]any{}
	if run.Parameters != nil {
		params = run.Parameters
	}
	metrics := map[string]any{}
	if run.Metrics != nil {
		metrics = run.Metrics
	}
	return &model.MlRun{
		ID: run.ID, OrgID: run.OrgID, ExperimentID: run.ExperimentID, Name: run.Name, Status: run.Status,
		StartedAt: run.StartedAt, EndedAt: run.EndedAt, Duration: run.Duration, Notes: run.Notes,
		Parameters: params, Metrics: metrics, DatasetID: run.DatasetID,
		UpdatedAt: run.UpdatedAt, SyncedAt: run.SyncedAt,
	}
}

func MLRunsToModel(in []uigraphapi.MLRun) []*model.MlRun {
	out := make([]*model.MlRun, len(in))
	for i := range in {
		out[i] = MLRunToModel(&in[i])
	}
	return out
}

func MLRunSeriesToJSON(points []uigraphapi.MLMetricPoint) map[string]any {
	series := map[string][]map[string]any{}
	for _, p := range points {
		series[p.Key] = append(series[p.Key], map[string]any{"step": p.Step, "value": p.Value})
	}
	out := map[string]any{}
	for k, v := range series {
		out[k] = v
	}
	return out
}

func MLArtifactToModel(a *uigraphapi.MLArtifact) *model.MlArtifact {
	return &model.MlArtifact{
		ID: a.ID, RunID: a.RunID, Name: a.Name, Type: a.Type,
		URI: a.URI, Size: a.Size, Format: a.Format,
		UpdatedAt: a.UpdatedAt, SyncedAt: a.SyncedAt,
	}
}

func MLArtifactsToModel(in []uigraphapi.MLArtifact) []*model.MlArtifact {
	out := make([]*model.MlArtifact, len(in))
	for i := range in {
		out[i] = MLArtifactToModel(&in[i])
	}
	return out
}

func MLDatasetToModel(ds *uigraphapi.MLDataset) *model.MlDataset {
	schema := make([]*model.MlSchemaField, len(ds.Schema))
	for i := range ds.Schema {
		schema[i] = &model.MlSchemaField{
			Name: ds.Schema[i].Name, Type: ds.Schema[i].Type, Description: ds.Schema[i].Description,
		}
	}
	tags := map[string]any{}
	for k, v := range ds.Tags {
		tags[k] = v
	}
	return &model.MlDataset{
		ID: ds.ID, ExperimentID: ds.ExperimentID, Name: ds.Name, Digest: ds.Digest,
		Source: ds.Source, SourceType: ds.SourceType, Context: ds.Context,
		RowCount: int(ds.RowCount), Schema: schema, Tags: tags,
	}
}

func MLDatasetsToModel(in []uigraphapi.MLDataset) []*model.MlDataset {
	out := make([]*model.MlDataset, len(in))
	for i := range in {
		out[i] = MLDatasetToModel(&in[i])
	}
	return out
}

func MLDeploymentToModel(d *uigraphapi.MLDeployment) *model.MlDeployment {
	return &model.MlDeployment{
		ID: d.ID, ModelID: d.ModelID, VersionID: d.VersionID, Name: d.Name, Environment: d.Environment,
		Status: d.Status, Endpoint: d.Endpoint, Region: d.Region,
		DeployedAt: d.DeployedAt, RolledBackAt: d.RolledBackAt,
	}
}

func MLDeploymentsToModel(in []uigraphapi.MLDeployment) []*model.MlDeployment {
	out := make([]*model.MlDeployment, len(in))
	for i := range in {
		out[i] = MLDeploymentToModel(&in[i])
	}
	return out
}

func MLFindingToModel(f *uigraphapi.MLFinding) *model.MlFinding {
	runIDs := f.RunIDs
	if runIDs == nil {
		runIDs = []string{}
	}
	return &model.MlFinding{
		ID: f.ID, ModelID: f.ModelID, VersionID: f.VersionID, Title: f.Title,
		Summary: f.Summary, Description: f.Description, RunIds: runIDs,
	}
}

func MLFindingsToModel(in []uigraphapi.MLFinding) []*model.MlFinding {
	out := make([]*model.MlFinding, len(in))
	for i := range in {
		out[i] = MLFindingToModel(&in[i])
	}
	return out
}
