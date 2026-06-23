package convert

import (
	"encoding/json"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

type parsedServiceDBSchema struct {
	Tables       []*model.DbTable `json:"tables"`
	NoSQLSchema  any              `json:"noSQLSchema"`
	DbDiagramID  *string          `json:"dbDiagramId"`
	PgDumpFileID *string          `json:"pgDumpFileId"`
}

func parseServiceDBSchema(b json.RawMessage) parsedServiceDBSchema {
	p := parsedServiceDBSchema{Tables: []*model.DbTable{}}
	if len(b) == 0 {
		return p
	}
	raw := b
	var asString string
	if json.Unmarshal(b, &asString) == nil {
		raw = json.RawMessage(asString)
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return parsedServiceDBSchema{Tables: []*model.DbTable{}}
	}
	if p.Tables == nil {
		p.Tables = []*model.DbTable{}
	}
	return p
}

func ServiceToModel(s *uigraphapi.Service) *model.Service {
	return &model.Service{
		ID: s.ID, OrgID: s.OrgID, FolderID: s.FolderID, TeamID: s.TeamID,
		Name: s.Name, Slug: s.Slug, Description: s.Description,
		Status: s.Status, Tier: s.Tier, Category: s.Category, Language: s.Language,
		GitRepoURL: s.GitRepoURL, JiraProjectURL: s.JiraProjectURL,
		SlackChannelURL: s.SlackChannelURL, LastCommitSha: s.LastCommitSha,
		Labels:    s.Labels,
		Metadata:  RawStr(s.Metadata),
		CreatedBy: s.CreatedBy, UpdatedBy: s.UpdatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func ServiceStatsToModel(s uigraphapi.ServiceStats) *model.ServiceStats {
	return &model.ServiceStats{
		ServiceID:     s.ServiceID,
		EndpointCount: s.EndpointCount,
		DiagramCount:  s.DiagramCount,
		DocCount:      s.DocCount,
		DbTableCount:  s.DBTableCount,
		TestCaseCount: s.TestCaseCount,
	}
}

func APIGroupToModel(g *uigraphapi.APIGroup) *model.APIGroup {
	return &model.APIGroup{
		ID: g.ID, ServiceID: g.ServiceID, OrgID: g.OrgID,
		Name: g.Name, Version: g.Version, Label: g.Label, Protocol: g.Protocol,
		SpecKey: g.SpecKey, SpecHash: g.SpecHash,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func APIGroupVersionToModel(orgID string, v uigraphapi.APIGroupVersion) *model.APIGroupVersion {
	return &model.APIGroupVersion{
		ID: v.ID, OrgID: orgID, APIGroupID: v.APIGroupID, VersionNumber: v.VersionNumber,
		Label: v.Label, SpecKey: v.SpecKey, SpecHash: v.SpecHash,
		IsAutoVersion: v.IsAutoVersion, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func ServiceDocToModel(d *uigraphapi.ServiceDoc) *model.ServiceDoc {
	out := &model.ServiceDoc{
		ServiceID: d.ServiceID,
		DocID:     d.DocID,
		OrgID:     d.OrgID,
		CreatedBy: d.CreatedBy,
		UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
	if d.Doc != nil {
		out.Doc = DocToModel(d.Doc)
	}
	return out
}

func ServiceDiagramToModel(d *uigraphapi.ServiceDiagram) *model.ServiceDiagram {
	out := &model.ServiceDiagram{
		ServiceID: d.ServiceID,
		DiagramID: d.DiagramID,
		OrgID:     d.OrgID,
		CreatedBy: d.CreatedBy,
		UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
	if d.Diagram != nil {
		out.Diagram = DiagramToModel(d.Diagram)
	}
	return out
}

func ServiceDBToModel(d *uigraphapi.ServiceDB) *model.ServiceDb {
	p := parseServiceDBSchema(d.SchemaJSON)
	return &model.ServiceDb{
		ID: d.ID, ServiceID: d.ServiceID, OrgID: d.OrgID,
		DbName: d.DBName, DbType: d.DBType, Dialect: d.Dialect,
		SchemaJSON:   RawStr(d.SchemaJSON),
		Tables:       p.Tables,
		NoSQLSchema:  p.NoSQLSchema,
		DbDiagramID:  p.DbDiagramID,
		PgDumpFileID: p.PgDumpFileID,
		Source:       d.Source, SourceTs: d.SourceTS,
		CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func ServiceDBVersionToModel(orgID string, v uigraphapi.ServiceDBVersion) *model.ServiceDBVersion {
	p := parseServiceDBSchema(v.SchemaJSON)
	return &model.ServiceDBVersion{
		ID: v.ID, OrgID: orgID, ServiceDbID: v.ServiceDBID, VersionNumber: v.VersionNumber,
		Label: v.Label, SchemaJSON: RawStr(v.SchemaJSON),
		Tables:       p.Tables,
		NoSQLSchema:  p.NoSQLSchema,
		DbDiagramID:  p.DbDiagramID,
		PgDumpFileID: p.PgDumpFileID,
		Source:       v.Source, SourceTs: v.SourceTS,
		IsAutoVersion: v.IsAutoVersion, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func APIEndpointToModel(e *uigraphapi.APIEndpoint) *model.APIEndpoint {
	return &model.APIEndpoint{
		ID: e.ID, APIGroupID: e.APIGroupID, ServiceID: e.ServiceID, OrgID: e.OrgID,
		OperationID: e.OperationID, Method: e.Method, Path: e.Path,
		Summary: e.Summary, Description: e.Description, Tags: e.Tags,
		Parameters:  RawArrStr(e.Parameters),
		RequestBody: RawStr(e.RequestBody),
		Responses:   RawStr(e.Responses),
		Order:       e.Order,
		CreatedBy:   e.CreatedBy, UpdatedBy: e.UpdatedBy, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

func ServicesToModel(services []uigraphapi.Service) []*model.Service {
	out := make([]*model.Service, len(services))
	for i := range services {
		out[i] = ServiceToModel(&services[i])
	}
	return out
}

func ServiceStatsListToModel(stats []uigraphapi.ServiceStats) []*model.ServiceStats {
	out := make([]*model.ServiceStats, len(stats))
	for i, s := range stats {
		out[i] = ServiceStatsToModel(s)
	}
	return out
}

func APIGroupsToModel(groups []uigraphapi.APIGroup) []*model.APIGroup {
	out := make([]*model.APIGroup, len(groups))
	for i := range groups {
		out[i] = APIGroupToModel(&groups[i])
	}
	return out
}

func APIGroupVersionsToModel(orgID string, versions []uigraphapi.APIGroupVersion) []*model.APIGroupVersion {
	out := make([]*model.APIGroupVersion, len(versions))
	for i, v := range versions {
		out[i] = APIGroupVersionToModel(orgID, v)
	}
	return out
}

func ServiceDocsToModel(docs []uigraphapi.ServiceDoc) []*model.ServiceDoc {
	out := make([]*model.ServiceDoc, len(docs))
	for i := range docs {
		out[i] = ServiceDocToModel(&docs[i])
	}
	return out
}

func ServiceDiagramsToModel(diagrams []uigraphapi.ServiceDiagram) []*model.ServiceDiagram {
	out := make([]*model.ServiceDiagram, len(diagrams))
	for i := range diagrams {
		out[i] = ServiceDiagramToModel(&diagrams[i])
	}
	return out
}

func ServiceDBsToModel(dbs []uigraphapi.ServiceDB) []*model.ServiceDb {
	out := make([]*model.ServiceDb, len(dbs))
	for i := range dbs {
		out[i] = ServiceDBToModel(&dbs[i])
	}
	return out
}

func ServiceDBVersionsToModel(orgID string, versions []uigraphapi.ServiceDBVersion) []*model.ServiceDBVersion {
	out := make([]*model.ServiceDBVersion, len(versions))
	for i, v := range versions {
		out[i] = ServiceDBVersionToModel(orgID, v)
	}
	return out
}

func APIEndpointsToModel(endpoints []uigraphapi.APIEndpoint) []*model.APIEndpoint {
	out := make([]*model.APIEndpoint, len(endpoints))
	for i := range endpoints {
		out[i] = APIEndpointToModel(&endpoints[i])
	}
	return out
}
