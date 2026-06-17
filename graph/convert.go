package graph

import (
	"context"
	"encoding/json"

	"github.com/uigraph/graphql/client"
	"github.com/uigraph/graphql/graph/model"
)

// resolveActor resolves a single created_by / updated_by id within an org to
// its public actor info, returning nil when id is empty or matches no actor.
func (r *Resolver) resolveActor(ctx context.Context, orgID, id string) (*model.Actor, error) {
	if id == "" {
		return nil, nil
	}
	actors, err := r.Client.ResolveActors(ctx, orgID, []string{id})
	if err != nil {
		return nil, err
	}
	a := actors[id]
	if a == nil {
		return nil, nil
	}
	m := &model.Actor{ID: a.ID, Type: a.Type, Name: a.Name, Disabled: a.Disabled}
	if a.Email != "" {
		m.Email = &a.Email
	}
	return m, nil
}

// resolveAssetURL resolves a single asset id within an org to a presigned GET
// URL, returning nil when id is empty or no url is produced.
func (r *Resolver) resolveAssetURL(ctx context.Context, orgID, assetID string) (*string, error) {
	if assetID == "" {
		return nil, nil
	}
	urls, err := r.Client.ResolveAssetURLs(ctx, orgID, []string{assetID})
	if err != nil {
		return nil, err
	}
	u, ok := urls[assetID]
	if !ok || u == "" {
		return nil, nil
	}
	return &u, nil
}

// toMap JSON-round-trips a struct into map[string]interface{}.
// This correctly handles optional fields: nil pointer fields are omitted
// from the resulting map (because of omitempty in the input struct tags).
func toMap(v interface{}) map[string]interface{} {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

// rawStr returns the JSON string of a raw message, defaulting to "{}".
func rawStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "{}"
	}
	return string(b)
}

func rawArrStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "[]"
	}
	return string(b)
}

// ── Auth ─────────────────────────────────────────────────────────────────────

func meToModel(m *client.MeResponse) *model.Me {
	return &model.Me{
		UserID: m.UserID, OrgID: m.OrgID,
		Email: m.Email, Name: m.Name, Login: m.Login,
		Kind: m.Kind, Role: m.Role, AuthProvider: m.AuthProvider,
	}
}

func orgSummaryToModel(o client.OrgSummary) *model.OrgSummary {
	return &model.OrgSummary{ID: o.ID, Name: o.Name, Slug: o.Slug, Role: o.Role, Active: o.Active}
}

// ── Org ───────────────────────────────────────────────────────────────────────

func orgToModel(o *client.Org) *model.Org {
	return &model.Org{ID: o.ID, Name: o.Name, Slug: o.Slug, Disabled: o.Disabled, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt}
}

func memberToModel(m client.Member) *model.Member {
	return &model.Member{UserID: m.UserID, OrgID: m.OrgID, Role: m.Role, Source: m.Source, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func teamToModel(t *client.Team) *model.Team {
	m := &model.Team{ID: t.ID, OrgID: t.OrgID, Name: t.Name, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
	if t.Email != "" {
		m.Email = &t.Email
	}
	if t.ExternalID != "" {
		m.ExternalID = &t.ExternalID
	}
	return m
}

func teamMemberToModel(m client.TeamMember) *model.TeamMember {
	return &model.TeamMember{TeamID: m.TeamID, UserID: m.UserID, Permission: m.Permission, CreatedAt: m.CreatedAt}
}

func invitationToModel(i client.Invitation) *model.Invitation {
	return &model.Invitation{
		ID: i.ID, OrgID: i.OrgID, Email: i.Email, Role: i.Role,
		Code: i.Code, CreatedBy: i.CreatedBy, CreatedAt: i.CreatedAt, ExpiresAt: i.ExpiresAt,
	}
}

func serviceAccountToModel(sa client.ServiceAccount) *model.ServiceAccount {
	return &model.ServiceAccount{
		ID: sa.ID, OrgID: sa.OrgID, Name: sa.Name, Description: sa.Description,
		Role: sa.Role, Disabled: sa.Disabled, CreatedAt: sa.CreatedAt, UpdatedAt: sa.UpdatedAt,
	}
}

func saTokenToModel(t client.ServiceAccountToken) *model.ServiceAccountToken {
	return &model.ServiceAccountToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name, Prefix: t.Prefix,
		ExpiresAt: t.ExpiresAt, LastUsedAt: t.LastUsedAt, Revoked: t.Revoked, CreatedAt: t.CreatedAt,
	}
}

func createdTokenToModel(t *client.CreatedToken) *model.CreatedToken {
	return &model.CreatedToken{
		ID: t.ID, ServiceAccountID: t.ServiceAccountID, Name: t.Name,
		Prefix: t.Prefix, Token: t.Token, CreatedAt: t.CreatedAt,
	}
}

// ── Content ───────────────────────────────────────────────────────────────────

func folderToModel(f *client.Folder) *model.Folder {
	return &model.Folder{
		ID: f.ID, OrgID: f.OrgID, ParentID: f.ParentID, TeamID: f.TeamID, Type: f.Type,
		Name: f.Name, Order: f.Order, CreatedBy: f.CreatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func diagramToModel(d *client.Diagram) *model.Diagram {
	return &model.Diagram{
		ID: d.ID, OrgID: d.OrgID, FolderID: d.FolderID, TeamID: d.TeamID,
		Name: d.Name, ContentKey: d.ContentKey, ContentHash: d.ContentHash,
		PreviewAssetID: d.PreviewAssetID, PreviewContentHash: d.PreviewContentHash,
		Source: d.Source, CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func diagramVersionToModel(v client.DiagramVersion) *model.DiagramVersion {
	return &model.DiagramVersion{
		ID: v.ID, DiagramID: v.DiagramID, VersionNumber: v.VersionNumber,
		Label: v.Label, ContentKey: v.ContentKey, ContentHash: v.ContentHash,
		IsAutoVersion: v.IsAutoVersion, Source: v.Source, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func flowComponentFieldToModel(f client.FlowDiagramComponentField) *model.FlowDiagramComponentField {
	return &model.FlowDiagramComponentField{
		FlowDiagramComponentFieldID: f.FlowDiagramComponentFieldID,
		Label:                       f.Label,
		Type:                        f.Type,
		Required:                    f.Required,
		Readonly:                    f.Readonly,
		Options:                     f.Options,
		Order:                       f.Order,
	}
}

func flowComponentToModel(c client.FlowDiagramComponent) *model.FlowDiagramComponent {
	fields := make([]*model.FlowDiagramComponentField, len(c.FlowDiagramComponentFields))
	for i, f := range c.FlowDiagramComponentFields {
		fields[i] = flowComponentFieldToModel(f)
	}
	return &model.FlowDiagramComponent{
		ComponentID: c.ComponentID, Type: c.Type, Name: c.Name,
		Description: c.Description, Category: c.Category, Tags: c.Tags,
		Slug: c.Slug, PreviewImageJpg: c.PreviewImageJpg, IsActive: c.IsActive,
		Order: c.Order, OrganizationID: c.OrganizationID,
		FlowDiagramComponentFields: fields,
	}
}

func flowComponentsToModel(components []client.FlowDiagramComponent) []*model.FlowDiagramComponent {
	out := make([]*model.FlowDiagramComponent, len(components))
	for i, c := range components {
		out[i] = flowComponentToModel(c)
	}
	return out
}

func componentFieldToModel(f client.ComponentField) *model.ComponentField {
	return &model.ComponentField{
		ComponentFieldID: f.ComponentFieldID,
		Label:            f.Label,
		Type:             f.Type,
		Required:         f.Required,
		Readonly:         f.Readonly,
		Options:          f.Options,
		Order:            f.Order,
	}
}

func componentToModel(c client.Component) *model.Component {
	fields := make([]*model.ComponentField, len(c.ComponentFields))
	for i, f := range c.ComponentFields {
		fields[i] = componentFieldToModel(f)
	}
	return &model.Component{
		ComponentID: c.ComponentID, Type: c.Type, Name: c.Name,
		Description: c.Description, Category: c.Category, Tags: c.Tags,
		Slug: c.Slug, PreviewImageJpg: c.PreviewImageJpg, IsActive: c.IsActive,
		Order: c.Order, ComponentFields: fields,
	}
}

func componentsToModel(components []client.Component) []*model.Component {
	out := make([]*model.Component, len(components))
	for i, c := range components {
		out[i] = componentToModel(c)
	}
	return out
}

func diagramImageToModel(img client.DiagramImage) *model.DiagramImage {
	return &model.DiagramImage{
		DiagramImageID: img.DiagramImageID, DiagramID: img.DiagramID,
		OrgID: img.OrgID, AssetID: img.AssetID, FileName: img.FileName,
		Order: img.Order, CreatedBy: img.CreatedBy, CreatedAt: img.CreatedAt,
	}
}

func diagramImagesToModel(images []client.DiagramImage) []*model.DiagramImage {
	out := make([]*model.DiagramImage, len(images))
	for i, img := range images {
		out[i] = diagramImageToModel(img)
	}
	return out
}

func uimapToModel(m *client.UIMap) *model.UIMap {
	return &model.UIMap{
		ID: m.ID, OrgID: m.OrgID, FolderID: m.FolderID, TeamID: m.TeamID,
		Name: m.Name, Description: m.Description, Status: m.Status,
		CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func frameToModel(f *client.Frame) *model.Frame {
	return &model.Frame{
		ID: f.ID, MapID: f.MapID, OrgID: f.OrgID, ParentFrameID: f.ParentFrameID,
		Name: f.Name, Description: f.Description, TemplateType: f.TemplateType,
		ScreenshotAssetID: f.ScreenshotAssetID, ScreenshotContentHash: f.ScreenshotContentHash,
		Status: f.Status, Order: f.Order, Source: f.Source,
		CreatedBy: f.CreatedBy, UpdatedBy: f.UpdatedBy, CreatedAt: f.CreatedAt, UpdatedAt: f.UpdatedAt,
	}
}

func focalPointToModel(fp *client.FocalPoint) *model.FocalPoint {
	return &model.FocalPoint{
		ID: fp.ID, FrameID: fp.FrameID, OrgID: fp.OrgID,
		Name: fp.Name, LocationX: fp.LocationX, LocationY: fp.LocationY,
		Visibility: fp.Visibility, IsActive: fp.IsActive,
		CreatedBy: fp.CreatedBy, UpdatedBy: fp.UpdatedBy, CreatedAt: fp.CreatedAt, UpdatedAt: fp.UpdatedAt,
	}
}

func canvasToModel(c *client.Canvas) *model.Canvas {
	return &model.Canvas{
		MapID: c.MapID, OrgID: c.OrgID,
		Zoom: c.Zoom, NavigationX: c.NavigationX, NavigationY: c.NavigationY,
		FramePositions: rawStr(c.FramePositions),
		UpdatedAt:      c.UpdatedAt,
	}
}

func frameGroupToModel(g *client.FrameGroup) *model.FrameGroup {
	return &model.FrameGroup{
		ID: g.ID, FrameID: g.FrameID, OrgID: g.OrgID,
		Name: g.Name, Description: g.Description,
		LocationX: g.LocationX, LocationY: g.LocationY,
		Width: g.Width, Height: g.Height, Order: g.Order, IsActive: g.IsActive,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy,
		CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func frameGroupsToModel(gs []client.FrameGroup) []*model.FrameGroup {
	out := make([]*model.FrameGroup, len(gs))
	for i := range gs {
		out[i] = frameGroupToModel(&gs[i])
	}
	return out
}

func frameLinkToModel(l *client.FrameLink) *model.FrameLink {
	return &model.FrameLink{
		ID: l.ID, FrameID: l.FrameID, OrgID: l.OrgID, Kind: l.Kind,
		TargetFrameID: l.TargetFrameID, TargetMapID: l.TargetMapID,
		Label: l.Label, LocationX: l.LocationX, LocationY: l.LocationY, IsActive: l.IsActive,
		CreatedBy: l.CreatedBy, UpdatedBy: l.UpdatedBy,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func frameLinksToModel(ls []client.FrameLink) []*model.FrameLink {
	out := make([]*model.FrameLink, len(ls))
	for i := range ls {
		out[i] = frameLinkToModel(&ls[i])
	}
	return out
}

func focalPointMetaToModel(m *client.FocalPointMeta) *model.FocalPointMeta {
	return &model.FocalPointMeta{
		ID: m.ID, FocalPointID: m.FocalPointID, OrgID: m.OrgID, FrameID: m.FrameID,
		ComponentID: m.ComponentID, ComponentLinkID: m.ComponentLinkID,
		ComponentImages:      rawArrStr(m.ComponentImages),
		ComponentFlowDiagram: m.ComponentFlowDiagram,
		ComponentModalFields: rawArrStr(m.ComponentModalFields),
		CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func focalPointMetasToModel(ms []client.FocalPointMeta) []*model.FocalPointMeta {
	out := make([]*model.FocalPointMeta, len(ms))
	for i := range ms {
		out[i] = focalPointMetaToModel(&ms[i])
	}
	return out
}

func focalPointMetaBody(body map[string]interface{}) map[string]interface{} {
	for _, key := range []string{"componentImages", "componentModalFields"} {
		if s, ok := body[key].(string); ok {
			var raw interface{}
			if err := unmarshalJSONString(s, &raw); err == nil {
				body[key] = raw
			}
		}
	}
	return body
}

func serviceToModel(s *client.Service) *model.Service {
	return &model.Service{
		ID: s.ID, OrgID: s.OrgID, FolderID: s.FolderID, TeamID: s.TeamID,
		Name: s.Name, Slug: s.Slug, Description: s.Description,
		Status: s.Status, Tier: s.Tier, Category: s.Category, Language: s.Language,
		GitRepoURL: s.GitRepoURL, JiraProjectURL: s.JiraProjectURL,
		SlackChannelURL: s.SlackChannelURL, LastCommitSha: s.LastCommitSha,
		Labels:    s.Labels,
		Metadata:  rawStr(s.Metadata),
		CreatedBy: s.CreatedBy, UpdatedBy: s.UpdatedBy, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func serviceStatsToModel(s client.ServiceStats) *model.ServiceStats {
	return &model.ServiceStats{
		ServiceID:     s.ServiceID,
		EndpointCount: s.EndpointCount,
		DiagramCount:  s.DiagramCount,
		DocCount:      s.DocCount,
		DbTableCount:  s.DBTableCount,
		TestCaseCount: s.TestCaseCount,
	}
}

func apiGroupToModel(g *client.APIGroup) *model.APIGroup {
	return &model.APIGroup{
		ID: g.ID, ServiceID: g.ServiceID, OrgID: g.OrgID,
		Name: g.Name, Version: g.Version, Label: g.Label, Protocol: g.Protocol,
		SpecKey: g.SpecKey, SpecHash: g.SpecHash,
		CreatedBy: g.CreatedBy, UpdatedBy: g.UpdatedBy, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt,
	}
}

func apiGroupVersionToModel(v client.APIGroupVersion) *model.APIGroupVersion {
	return &model.APIGroupVersion{
		ID: v.ID, APIGroupID: v.APIGroupID, VersionNumber: v.VersionNumber,
		Label: v.Label, SpecKey: v.SpecKey, SpecHash: v.SpecHash,
		IsAutoVersion: v.IsAutoVersion, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func serviceDocToModel(d *client.ServiceDoc) *model.ServiceDoc {
	return &model.ServiceDoc{
		ID: d.ID, ServiceID: d.ServiceID, OrgID: d.OrgID,
		FileKey: d.FileKey, FileName: d.FileName, FileType: d.FileType,
		Description: d.Description, ContentHash: d.ContentHash,
		CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func serviceDiagramToModel(d *client.ServiceDiagram) *model.ServiceDiagram {
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
		out.Diagram = diagramToModel(d.Diagram)
	}
	return out
}

func serviceDBToModel(d *client.ServiceDB) *model.ServiceDb {
	return &model.ServiceDb{
		ID: d.ID, ServiceID: d.ServiceID, OrgID: d.OrgID,
		DbName: d.DBName, DbType: d.DBType, Dialect: d.Dialect,
		SchemaJSON: rawStr(d.SchemaJSON),
		Source:     d.Source, SourceTs: d.SourceTS,
		CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func serviceDBVersionToModel(v client.ServiceDBVersion) *model.ServiceDBVersion {
	return &model.ServiceDBVersion{
		ID: v.ID, ServiceDbID: v.ServiceDBID, VersionNumber: v.VersionNumber,
		Label: v.Label, SchemaJSON: rawStr(v.SchemaJSON),
		Source: v.Source, SourceTs: v.SourceTS,
		IsAutoVersion: v.IsAutoVersion, CreatedBy: v.CreatedBy, CreatedAt: v.CreatedAt,
	}
}

func apiEndpointToModel(e *client.APIEndpoint) *model.APIEndpoint {
	return &model.APIEndpoint{
		ID: e.ID, APIGroupID: e.APIGroupID, ServiceID: e.ServiceID, OrgID: e.OrgID,
		OperationID: e.OperationID, Method: e.Method, Path: e.Path,
		Summary: e.Summary, Description: e.Description, Tags: e.Tags,
		Parameters:  rawArrStr(e.Parameters),
		RequestBody: rawStr(e.RequestBody),
		Responses:   rawStr(e.Responses),
		Order:       e.Order,
		CreatedBy:   e.CreatedBy, UpdatedBy: e.UpdatedBy, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

func keyValueToModel(v client.KeyValue) *model.KeyValue {
	return &model.KeyValue{Key: v.Key, Value: v.Value}
}

func assertionToModel(a client.Assertion) *model.Assertion {
	return &model.Assertion{Field: a.Field, Type: a.Type, Value: a.Value}
}

func authConfigToModel(a *client.AuthConfig) *model.AuthConfig {
	if a == nil {
		return nil
	}
	return &model.AuthConfig{
		Type:          a.Type,
		BearerToken:   a.BearerToken,
		APIKeyHeader:  a.APIKeyHeader,
		APIKeyValue:   a.APIKeyValue,
		BasicUsername: a.BasicUsername,
		BasicPassword: a.BasicPassword,
	}
}

func testCaseStepToModel(s client.TestCaseStep) *model.TestCaseStep {
	return &model.TestCaseStep{Order: s.Order, Action: s.Action, ExpectedResult: s.ExpectedResult}
}

func manualTestCaseToModel(m *client.ManualTestCase) *model.ManualTestCase {
	if m == nil {
		return nil
	}
	steps := make([]*model.TestCaseStep, len(m.Steps))
	for i, s := range m.Steps {
		steps[i] = testCaseStepToModel(s)
	}
	return &model.ManualTestCase{
		Preconditions:   m.Preconditions,
		TestData:        m.TestData,
		Steps:           steps,
		ExpectedOutcome: m.ExpectedOutcome,
		Postconditions:  m.Postconditions,
	}
}

func apiTestCaseToModel(a *client.APITestCase) *model.APITestCase {
	if a == nil {
		return nil
	}
	headers := make([]*model.KeyValue, len(a.RequestHeaders))
	for i, v := range a.RequestHeaders {
		headers[i] = keyValueToModel(v)
	}
	params := make([]*model.KeyValue, len(a.QueryParams))
	for i, v := range a.QueryParams {
		params[i] = keyValueToModel(v)
	}
	assertions := make([]*model.Assertion, len(a.Assertions))
	for i, v := range a.Assertions {
		assertions[i] = assertionToModel(v)
	}
	return &model.APITestCase{
		HTTPMethod:         a.HTTPMethod,
		APISpecID:          a.APISpecID,
		OperationID:        a.OperationID,
		Auth:               authConfigToModel(a.Auth),
		RequestHeaders:     headers,
		QueryParams:        params,
		RequestBody:        a.RequestBody,
		ExpectedStatusCode: a.ExpectedStatusCode,
		MaxResponseTimeMs:  a.MaxResponseTimeMs,
		ResponseBody:       a.ResponseBody,
		Assertions:         assertions,
	}
}

func graphQLTestCaseToModel(g *client.GraphQLTestCase) *model.GraphQLTestCase {
	if g == nil {
		return nil
	}
	assertions := make([]*model.Assertion, len(g.Assertions))
	for i, v := range g.Assertions {
		assertions[i] = assertionToModel(v)
	}
	return &model.GraphQLTestCase{
		OperationType: g.OperationType,
		OperationName: g.OperationName,
		Query:         g.Query,
		Variables:     g.Variables,
		ResponseBody:  g.ResponseBody,
		Assertions:    assertions,
		ExpectError:   g.ExpectError,
	}
}

func databaseTestCaseToModel(d *client.DatabaseTestCase) *model.DatabaseTestCase {
	if d == nil {
		return nil
	}
	assertions := make([]*model.Assertion, len(d.Assertions))
	for i, v := range d.Assertions {
		assertions[i] = assertionToModel(v)
	}
	return &model.DatabaseTestCase{
		Dialect:       d.Dialect,
		SchemaID:      d.SchemaID,
		Query:         d.Query,
		Assertions:    assertions,
		SetupQuery:    d.SetupQuery,
		TeardownQuery: d.TeardownQuery,
	}
}

func grpcTestCaseToModel(g *client.GRPCTestCase) *model.GRPCTestCase {
	if g == nil {
		return nil
	}
	metadata := make([]*model.KeyValue, len(g.Metadata))
	for i, v := range g.Metadata {
		metadata[i] = keyValueToModel(v)
	}
	assertions := make([]*model.Assertion, len(g.Assertions))
	for i, v := range g.Assertions {
		assertions[i] = assertionToModel(v)
	}
	return &model.GRPCTestCase{
		ServiceName:    g.ServiceName,
		MethodName:     g.MethodName,
		CallMode:       g.CallMode,
		ProtoFileID:    g.ProtoFileID,
		ServerAddress:  g.ServerAddress,
		RequestMessage: g.RequestMessage,
		Metadata:       metadata,
		ExpectedStatus: g.ExpectedStatus,
		DeadlineMs:     g.DeadlineMs,
		ResponseBody:   g.ResponseBody,
		Assertions:     assertions,
		UseTLS:         g.UseTLS,
		ExpectError:    g.ExpectError,
	}
}

func testPackToModel(p *client.TestPack) *model.TestPack {
	return &model.TestPack{
		TestPackID: p.TestPackID, ServiceID: p.ServiceID, OrgID: p.OrgID,
		Name: p.Name, Type: p.Type,
		CreatedBy: p.CreatedBy, UpdatedBy: p.UpdatedBy, DeletedBy: p.DeletedBy,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, DeletedAt: p.DeletedAt,
	}
}

func testCaseToModel(tc *client.TestCase) *model.TestCase {
	return &model.TestCase{
		TestCaseID: tc.TestCaseID, TestPackID: tc.TestPackID, ServiceID: tc.ServiceID, OrgID: tc.OrgID,
		Title: tc.Title, Order: tc.Order, Type: tc.Type, Description: tc.Description, Priority: tc.Priority,
		Labels: tc.Labels, LinkedTicket: tc.LinkedTicket, EstimatedDurationMins: tc.EstimatedDurationMins,
		TestOwner: tc.TestOwner, LinkedMapNodeID: tc.LinkedMapNodeID, IsCritical: tc.IsCritical, EvidenceRequired: tc.EvidenceRequired,
		Manual: manualTestCaseToModel(tc.Manual), API: apiTestCaseToModel(tc.API),
		Graphql: graphQLTestCaseToModel(tc.GraphQL), Database: databaseTestCaseToModel(tc.Database), Grpc: grpcTestCaseToModel(tc.GRPC),
		Status: tc.Status, Version: tc.Version, BaselineRunResultID: tc.BaselineRunResultID, Dependencies: tc.Dependencies,
		CreatedBy: tc.CreatedBy, UpdatedBy: tc.UpdatedBy, DeletedBy: tc.DeletedBy, CreatedAt: tc.CreatedAt, UpdatedAt: tc.UpdatedAt, DeletedAt: tc.DeletedAt,
	}
}

func testRunToModel(tr *client.TestRun) *model.TestRun {
	return &model.TestRun{
		TestRunID: tr.TestRunID, TestPackID: tr.TestPackID, ServiceID: tr.ServiceID, OrgID: tr.OrgID,
		Environment: tr.Environment, ReleaseLabel: tr.ReleaseLabel, StartedAt: tr.StartedAt, CompletedAt: tr.CompletedAt,
		Status: tr.Status, StartedBy: tr.StartedBy, ExecutedBy: tr.ExecutedBy, ExecutedAt: tr.ExecutedAt, OverallStatus: tr.OverallStatus,
	}
}

func testRunSummaryToModel(s client.TestRunSummary) *model.TestRunSummary {
	return &model.TestRunSummary{
		TestRunID: s.TestRunID, TestPackID: s.TestPackID, ServiceID: s.ServiceID,
		Environment: s.Environment, ReleaseLabel: s.ReleaseLabel, StartedAt: s.StartedAt, CompletedAt: s.CompletedAt,
		Status: s.Status, StartedBy: s.StartedBy, ExecutedBy: s.ExecutedBy, ExecutedAt: s.ExecutedAt, OverallStatus: s.OverallStatus,
		PassedCount: s.PassedCount, FailedCount: s.FailedCount, SkippedCount: s.SkippedCount, BlockedCount: s.BlockedCount,
	}
}

func testRunResultToModel(rr *client.TestRunResult) *model.TestRunResult {
	var responseTimeMs *int
	if rr.ResponseTimeMs != nil {
		v := int(*rr.ResponseTimeMs)
		responseTimeMs = &v
	}
	return &model.TestRunResult{
		TestRunResultID: rr.TestRunResultID, TestRunID: rr.TestRunID, TestCaseID: rr.TestCaseID,
		ServiceID: rr.ServiceID, OrgID: rr.OrgID, Status: rr.Status, BlockedReason: rr.BlockedReason,
		ResponseStatus: rr.ResponseStatus, ResponseBody: rr.ResponseBody, ResponseTimeMs: responseTimeMs,
		Notes: rr.Notes, ScreenshotUrls: rr.ScreenshotURLs, ExecutedAt: rr.ExecutedAt, ExecutedBy: rr.ExecutedBy,
	}
}

// ── List helpers ──────────────────────────────────────────────────────────────

func orgsToModel(orgs []client.Org) []*model.Org {
	out := make([]*model.Org, len(orgs))
	for i := range orgs {
		out[i] = orgToModel(&orgs[i])
	}
	return out
}

func orgSummariesToModel(orgs []client.OrgSummary) []*model.OrgSummary {
	out := make([]*model.OrgSummary, len(orgs))
	for i, o := range orgs {
		out[i] = orgSummaryToModel(o)
	}
	return out
}

func membersToModel(members []client.Member) []*model.Member {
	out := make([]*model.Member, len(members))
	for i, m := range members {
		out[i] = memberToModel(m)
	}
	return out
}

func teamsToModel(teams []client.Team) []*model.Team {
	out := make([]*model.Team, len(teams))
	for i := range teams {
		out[i] = teamToModel(&teams[i])
	}
	return out
}

func teamMembersToModel(members []client.TeamMember) []*model.TeamMember {
	out := make([]*model.TeamMember, len(members))
	for i, m := range members {
		out[i] = teamMemberToModel(m)
	}
	return out
}

func invitationsToModel(invs []client.Invitation) []*model.Invitation {
	out := make([]*model.Invitation, len(invs))
	for i, inv := range invs {
		out[i] = invitationToModel(inv)
	}
	return out
}

func serviceAccountsToModel(sas []client.ServiceAccount) []*model.ServiceAccount {
	out := make([]*model.ServiceAccount, len(sas))
	for i, sa := range sas {
		out[i] = serviceAccountToModel(sa)
	}
	return out
}

func saTokensToModel(tokens []client.ServiceAccountToken) []*model.ServiceAccountToken {
	out := make([]*model.ServiceAccountToken, len(tokens))
	for i, t := range tokens {
		out[i] = saTokenToModel(t)
	}
	return out
}

func foldersToModel(folders []client.Folder) []*model.Folder {
	out := make([]*model.Folder, len(folders))
	for i := range folders {
		out[i] = folderToModel(&folders[i])
	}
	return out
}

func diagramsToModel(diagrams []client.Diagram) []*model.Diagram {
	out := make([]*model.Diagram, len(diagrams))
	for i := range diagrams {
		out[i] = diagramToModel(&diagrams[i])
	}
	return out
}

func diagramVersionsToModel(versions []client.DiagramVersion) []*model.DiagramVersion {
	out := make([]*model.DiagramVersion, len(versions))
	for i, v := range versions {
		out[i] = diagramVersionToModel(v)
	}
	return out
}

func uimapsToModel(maps []client.UIMap) []*model.UIMap {
	out := make([]*model.UIMap, len(maps))
	for i := range maps {
		out[i] = uimapToModel(&maps[i])
	}
	return out
}

func framesToModel(frames []client.Frame) []*model.Frame {
	out := make([]*model.Frame, len(frames))
	for i := range frames {
		out[i] = frameToModel(&frames[i])
	}
	return out
}

func focalPointsToModel(fps []client.FocalPoint) []*model.FocalPoint {
	out := make([]*model.FocalPoint, len(fps))
	for i := range fps {
		out[i] = focalPointToModel(&fps[i])
	}
	return out
}

func servicesToModel(services []client.Service) []*model.Service {
	out := make([]*model.Service, len(services))
	for i := range services {
		out[i] = serviceToModel(&services[i])
	}
	return out
}

func serviceStatsListToModel(stats []client.ServiceStats) []*model.ServiceStats {
	out := make([]*model.ServiceStats, len(stats))
	for i, s := range stats {
		out[i] = serviceStatsToModel(s)
	}
	return out
}

func apiGroupsToModel(groups []client.APIGroup) []*model.APIGroup {
	out := make([]*model.APIGroup, len(groups))
	for i := range groups {
		out[i] = apiGroupToModel(&groups[i])
	}
	return out
}

func apiGroupVersionsToModel(versions []client.APIGroupVersion) []*model.APIGroupVersion {
	out := make([]*model.APIGroupVersion, len(versions))
	for i, v := range versions {
		out[i] = apiGroupVersionToModel(v)
	}
	return out
}

func serviceDocsToModel(docs []client.ServiceDoc) []*model.ServiceDoc {
	out := make([]*model.ServiceDoc, len(docs))
	for i := range docs {
		out[i] = serviceDocToModel(&docs[i])
	}
	return out
}

func serviceDiagramsToModel(diagrams []client.ServiceDiagram) []*model.ServiceDiagram {
	out := make([]*model.ServiceDiagram, len(diagrams))
	for i := range diagrams {
		out[i] = serviceDiagramToModel(&diagrams[i])
	}
	return out
}

func serviceDBsToModel(dbs []client.ServiceDB) []*model.ServiceDb {
	out := make([]*model.ServiceDb, len(dbs))
	for i := range dbs {
		out[i] = serviceDBToModel(&dbs[i])
	}
	return out
}

func serviceDBVersionsToModel(versions []client.ServiceDBVersion) []*model.ServiceDBVersion {
	out := make([]*model.ServiceDBVersion, len(versions))
	for i, v := range versions {
		out[i] = serviceDBVersionToModel(v)
	}
	return out
}

func apiEndpointsToModel(endpoints []client.APIEndpoint) []*model.APIEndpoint {
	out := make([]*model.APIEndpoint, len(endpoints))
	for i := range endpoints {
		out[i] = apiEndpointToModel(&endpoints[i])
	}
	return out
}

func testPacksToModel(packs []client.TestPack) []*model.TestPack {
	out := make([]*model.TestPack, len(packs))
	for i := range packs {
		out[i] = testPackToModel(&packs[i])
	}
	return out
}

func testCasesToModel(cases []client.TestCase) []*model.TestCase {
	out := make([]*model.TestCase, len(cases))
	for i := range cases {
		out[i] = testCaseToModel(&cases[i])
	}
	return out
}

func testRunsToModel(runs []client.TestRun) []*model.TestRun {
	out := make([]*model.TestRun, len(runs))
	for i := range runs {
		out[i] = testRunToModel(&runs[i])
	}
	return out
}

func testRunSummariesToModel(summaries []client.TestRunSummary) []*model.TestRunSummary {
	out := make([]*model.TestRunSummary, len(summaries))
	for i, s := range summaries {
		out[i] = testRunSummaryToModel(s)
	}
	return out
}

func testRunResultsToModel(results []client.TestRunResult) []*model.TestRunResult {
	out := make([]*model.TestRunResult, len(results))
	for i := range results {
		out[i] = testRunResultToModel(&results[i])
	}
	return out
}

// ── Users ─────────────────────────────────────────────────────────────────────

func userToModel(u *client.User) *model.User {
	return &model.User{
		ID: u.ID, Email: u.Email, Name: u.Name, Login: u.Login,
		Disabled: u.Disabled, Role: u.Role, LastSeenAt: u.LastSeenAt,
		CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}

func usersToModel(users []client.User) []*model.User {
	out := make([]*model.User, len(users))
	for i := range users {
		out[i] = userToModel(&users[i])
	}
	return out
}

// ── SSO ───────────────────────────────────────────────────────────────────────

func oauthProviderToModel(p client.OAuthProvider) *model.OAuthProvider {
	return &model.OAuthProvider{
		ID: p.ID, ProviderName: p.ProviderName, Type: p.Type, DisplayName: p.DisplayName,
		ClientID: p.ClientID, ClientSecret: p.ClientSecret,
		AuthURL: p.AuthURL, TokenURL: p.TokenURL, UserinfoURL: p.UserinfoURL, APIURL: p.APIURL,
		Scopes: p.Scopes, AllowedDomains: p.AllowedDomains, AllowSignUp: p.AllowSignUp,
		EmailClaim: p.EmailClaim, NameClaim: p.NameClaim, SubClaim: p.SubClaim,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func oauthProvidersToModel(providers []client.OAuthProvider) []*model.OAuthProvider {
	out := make([]*model.OAuthProvider, len(providers))
	for i := range providers {
		out[i] = oauthProviderToModel(providers[i])
	}
	return out
}

func roleMappingToModel(m client.RoleMapping) *model.RoleMapping {
	return &model.RoleMapping{
		ID: m.ID, OrganizationID: m.OrganizationID,
		ClaimKey: m.ClaimKey, ClaimValue: m.ClaimValue, Role: m.Role, Scope: m.Scope,
		ResourceType: m.ResourceType, ResourceID: m.ResourceID,
	}
}

func roleMappingsToModel(mappings []client.RoleMapping) []*model.RoleMapping {
	out := make([]*model.RoleMapping, len(mappings))
	for i := range mappings {
		out[i] = roleMappingToModel(mappings[i])
	}
	return out
}

func ldapToModel(l *client.LDAPConfig) *model.LDAPConfig {
	return &model.LDAPConfig{
		ID: l.ID, Host: l.Host, Port: l.Port,
		UseSsl: l.UseSSL, StartTLS: l.StartTLS, SkipTLSVerify: l.SkipTLSVerify,
		BindDn: l.BindDN, BindPassword: l.BindPassword,
		SearchBaseDn: l.SearchBaseDN, SearchFilter: l.SearchFilter,
		EmailAttribute: l.EmailAttribute, NameAttribute: l.NameAttribute,
		UsernameAttribute: l.UsernameAttribute, MemberOfAttribute: l.MemberOfAttribute,
		AllowSignUp: l.AllowSignUp, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func samlToModel(s *client.SAMLConfig) *model.SAMLConfig {
	return &model.SAMLConfig{
		ID: s.ID, IdpMetadataURL: s.IDPMetadataURL, IdpMetadataXML: s.IDPMetadataXML,
		IdpEntityID: s.IDPEntityID, IdpSsoURL: s.IDPSsoURL, IdpCert: s.IDPCert,
		SpEntityID: s.SPEntityID, SpCert: s.SPCert, SpKey: s.SPKey,
		SignRequests: s.SignRequests, NameIDFormat: s.NameIDFormat,
		EmailAttribute: s.EmailAttribute, NameAttribute: s.NameAttribute, LoginAttribute: s.LoginAttribute,
		GroupsAttribute: s.GroupsAttribute, AllowSignUp: s.AllowSignUp,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}
