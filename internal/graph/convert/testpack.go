package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func KeyValueToModel(v uigraphapi.KeyValue) *model.KeyValue {
	return &model.KeyValue{Key: v.Key, Value: v.Value}
}

func AssertionToModel(a uigraphapi.Assertion) *model.Assertion {
	return &model.Assertion{Field: a.Field, Type: a.Type, Value: a.Value}
}

func AuthConfigToModel(a *uigraphapi.AuthConfig) *model.AuthConfig {
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

func TestCaseStepToModel(s uigraphapi.TestCaseStep) *model.TestCaseStep {
	return &model.TestCaseStep{Order: s.Order, Action: s.Action, ExpectedResult: s.ExpectedResult}
}

func ManualTestCaseToModel(m *uigraphapi.ManualTestCase) *model.ManualTestCase {
	if m == nil {
		return nil
	}
	steps := make([]*model.TestCaseStep, len(m.Steps))
	for i, s := range m.Steps {
		steps[i] = TestCaseStepToModel(s)
	}
	return &model.ManualTestCase{
		Preconditions:   m.Preconditions,
		TestData:        m.TestData,
		Steps:           steps,
		ExpectedOutcome: m.ExpectedOutcome,
		Postconditions:  m.Postconditions,
	}
}

func APITestCaseToModel(a *uigraphapi.APITestCase) *model.APITestCase {
	if a == nil {
		return nil
	}
	headers := make([]*model.KeyValue, len(a.RequestHeaders))
	for i, v := range a.RequestHeaders {
		headers[i] = KeyValueToModel(v)
	}
	params := make([]*model.KeyValue, len(a.QueryParams))
	for i, v := range a.QueryParams {
		params[i] = KeyValueToModel(v)
	}
	assertions := make([]*model.Assertion, len(a.Assertions))
	for i, v := range a.Assertions {
		assertions[i] = AssertionToModel(v)
	}
	return &model.APITestCase{
		HTTPMethod:         a.HTTPMethod,
		APISpecID:          a.APISpecID,
		OperationID:        a.OperationID,
		Auth:               AuthConfigToModel(a.Auth),
		RequestHeaders:     headers,
		QueryParams:        params,
		RequestBody:        a.RequestBody,
		ExpectedStatusCode: a.ExpectedStatusCode,
		MaxResponseTimeMs:  a.MaxResponseTimeMs,
		ResponseBody:       a.ResponseBody,
		Assertions:         assertions,
	}
}

func GraphQLTestCaseToModel(g *uigraphapi.GraphQLTestCase) *model.GraphQLTestCase {
	if g == nil {
		return nil
	}
	assertions := make([]*model.Assertion, len(g.Assertions))
	for i, v := range g.Assertions {
		assertions[i] = AssertionToModel(v)
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

func DatabaseTestCaseToModel(d *uigraphapi.DatabaseTestCase) *model.DatabaseTestCase {
	if d == nil {
		return nil
	}
	assertions := make([]*model.Assertion, len(d.Assertions))
	for i, v := range d.Assertions {
		assertions[i] = AssertionToModel(v)
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

func GRPCTestCaseToModel(g *uigraphapi.GRPCTestCase) *model.GRPCTestCase {
	if g == nil {
		return nil
	}
	metadata := make([]*model.KeyValue, len(g.Metadata))
	for i, v := range g.Metadata {
		metadata[i] = KeyValueToModel(v)
	}
	assertions := make([]*model.Assertion, len(g.Assertions))
	for i, v := range g.Assertions {
		assertions[i] = AssertionToModel(v)
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

func TestPackToModel(p *uigraphapi.TestPack) *model.TestPack {
	return &model.TestPack{
		TestPackID: p.TestPackID, ServiceID: p.ServiceID, OrgID: p.OrgID,
		Name: p.Name, Type: p.Type,
		CreatedBy: p.CreatedBy, UpdatedBy: p.UpdatedBy, DeletedBy: p.DeletedBy,
		CreatedByCommitHash: p.CreatedByCommitHash, UpdatedByCommitHash: p.UpdatedByCommitHash,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, DeletedAt: p.DeletedAt,
	}
}

func TestCaseToModel(tc *uigraphapi.TestCase) *model.TestCase {
	return &model.TestCase{
		TestCaseID: tc.TestCaseID, TestPackID: tc.TestPackID, ServiceID: tc.ServiceID, OrgID: tc.OrgID,
		Title: tc.Title, Order: tc.Order, Type: tc.Type, Description: tc.Description, Priority: tc.Priority,
		Labels: tc.Labels, LinkedTicket: tc.LinkedTicket, EstimatedDurationMins: tc.EstimatedDurationMins,
		TestOwner: tc.TestOwner, LinkedMapNodeID: tc.LinkedMapNodeID, IsCritical: tc.IsCritical, EvidenceRequired: tc.EvidenceRequired,
		Manual: ManualTestCaseToModel(tc.Manual), API: APITestCaseToModel(tc.API),
		Graphql: GraphQLTestCaseToModel(tc.GraphQL), Database: DatabaseTestCaseToModel(tc.Database), Grpc: GRPCTestCaseToModel(tc.GRPC),
		Status: tc.Status, Version: tc.Version, BaselineRunResultID: tc.BaselineRunResultID, Dependencies: tc.Dependencies,
		CreatedBy: tc.CreatedBy, UpdatedBy: tc.UpdatedBy, DeletedBy: tc.DeletedBy,
		CreatedByCommitHash: tc.CreatedByCommitHash, UpdatedByCommitHash: tc.UpdatedByCommitHash,
		CreatedAt: tc.CreatedAt, UpdatedAt: tc.UpdatedAt, DeletedAt: tc.DeletedAt,
	}
}

func TestRunToModel(tr *uigraphapi.TestRun) *model.TestRun {
	return &model.TestRun{
		TestRunID: tr.TestRunID, TestPackID: tr.TestPackID, ServiceID: tr.ServiceID, OrgID: tr.OrgID,
		Environment: tr.Environment, ReleaseLabel: tr.ReleaseLabel, StartedAt: tr.StartedAt, CompletedAt: tr.CompletedAt,
		Status: tr.Status, StartedBy: tr.StartedBy, ExecutedBy: tr.ExecutedBy, ExecutedAt: tr.ExecutedAt, OverallStatus: tr.OverallStatus,
	}
}

func TestRunSummaryToModel(s uigraphapi.TestRunSummary) *model.TestRunSummary {
	return &model.TestRunSummary{
		TestRunID: s.TestRunID, TestPackID: s.TestPackID, ServiceID: s.ServiceID,
		Environment: s.Environment, ReleaseLabel: s.ReleaseLabel, StartedAt: s.StartedAt, CompletedAt: s.CompletedAt,
		Status: s.Status, StartedBy: s.StartedBy, ExecutedBy: s.ExecutedBy, ExecutedAt: s.ExecutedAt, OverallStatus: s.OverallStatus,
		PassedCount: s.PassedCount, FailedCount: s.FailedCount, SkippedCount: s.SkippedCount, BlockedCount: s.BlockedCount,
	}
}

func TestRunResultToModel(rr *uigraphapi.TestRunResult) *model.TestRunResult {
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

func TestPacksToModel(packs []uigraphapi.TestPack) []*model.TestPack {
	out := make([]*model.TestPack, len(packs))
	for i := range packs {
		out[i] = TestPackToModel(&packs[i])
	}
	return out
}

func TestCasesToModel(cases []uigraphapi.TestCase) []*model.TestCase {
	out := make([]*model.TestCase, len(cases))
	for i := range cases {
		out[i] = TestCaseToModel(&cases[i])
	}
	return out
}

func TestRunsToModel(runs []uigraphapi.TestRun) []*model.TestRun {
	out := make([]*model.TestRun, len(runs))
	for i := range runs {
		out[i] = TestRunToModel(&runs[i])
	}
	return out
}

func TestRunSummariesToModel(summaries []uigraphapi.TestRunSummary) []*model.TestRunSummary {
	out := make([]*model.TestRunSummary, len(summaries))
	for i, s := range summaries {
		out[i] = TestRunSummaryToModel(s)
	}
	return out
}

func TestRunResultsToModel(results []uigraphapi.TestRunResult) []*model.TestRunResult {
	out := make([]*model.TestRunResult, len(results))
	for i := range results {
		out[i] = TestRunResultToModel(&results[i])
	}
	return out
}
