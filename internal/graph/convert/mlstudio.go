package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func MLModelToModel(m *uigraphapi.MLModel) *model.MlModel {
	return &model.MlModel{
		ID: m.ID, Name: m.Name, Description: m.Description,
		Domain: m.Domain, ProblemType: m.ProblemType, Tags: m.Tags,
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
		Description: v.Description, Status: v.Status, Stage: v.Stage, RunID: v.RunID,
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

func MLExperimentToModel(e *uigraphapi.MLExperiment) *model.MlExperiment {
	return &model.MlExperiment{
		ID: e.ID, Name: e.Name, Description: e.Description, Status: e.Status,
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
	return &model.MlDataset{
		ID: ds.ID, Name: ds.Name, Source: ds.Source, Type: ds.Type,
		RowCount: int(ds.RowCount), Schema: schema,
	}
}

func MLDatasetsToModel(in []uigraphapi.MLDataset) []*model.MlDataset {
	out := make([]*model.MlDataset, len(in))
	for i := range in {
		out[i] = MLDatasetToModel(&in[i])
	}
	return out
}

func MLEvaluationDatasetToModel(ds *uigraphapi.MLEvaluationDataset) *model.MlEvaluationDataset {
	schema := make([]*model.MlSchemaField, len(ds.Schema))
	for i := range ds.Schema {
		schema[i] = &model.MlSchemaField{
			Name: ds.Schema[i].Name, Type: ds.Schema[i].Type, Description: ds.Schema[i].Description,
		}
	}
	return &model.MlEvaluationDataset{
		ID: ds.ID, Name: ds.Name, Digest: ds.Digest, Source: ds.Source,
		SourceType: ds.SourceType, RowCount: int(ds.RowCount), Schema: schema,
	}
}

func MLEvaluationDatasetsToModel(in []uigraphapi.MLEvaluationDataset) []*model.MlEvaluationDataset {
	out := make([]*model.MlEvaluationDataset, len(in))
	for i := range in {
		out[i] = MLEvaluationDatasetToModel(&in[i])
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
