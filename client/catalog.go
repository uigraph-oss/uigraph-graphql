package client

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// ── Services ──────────────────────────────────────────────────────────────────

func (c *Client) ListServices(ctx context.Context, orgID, folderID, teamID string) ([]Service, error) {
	path := "/api/v1/orgs/" + orgID + "/services"
	q := url.Values{}
	if folderID != "" {
		q.Set("folderId", folderID)
	}
	if teamID != "" {
		q.Set("teamId", teamID)
	}
	if enc := q.Encode(); enc != "" {
		path += "?" + enc
	}
	var out struct {
		Services []Service `json:"services"`
	}
	return out.Services, c.get(ctx, path, &out)
}

func (c *Client) GetService(ctx context.Context, orgID, id string) (*Service, error) {
	var out Service
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id), &out)
}

func (c *Client) CreateService(ctx context.Context, orgID string, body map[string]interface{}) (*Service, error) {
	var out Service
	return &out, c.post(ctx, "/api/v1/orgs/"+orgID+"/services", body, &out)
}

func (c *Client) UpdateService(ctx context.Context, orgID, id string, body map[string]interface{}) (*Service, error) {
	var out Service
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id), body, &out)
}

func (c *Client) DeleteService(ctx context.Context, orgID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s", orgID, id))
}

func (c *Client) ListServiceStats(ctx context.Context, orgID string, serviceID *string) ([]ServiceStats, error) {
	path := "/api/v1/orgs/" + orgID + "/services/stats"
	if serviceID != nil && *serviceID != "" {
		q := url.Values{}
		q.Set("serviceId", *serviceID)
		path += "?" + q.Encode()
	}
	var out struct {
		Stats []ServiceStats `json:"stats"`
	}
	return out.Stats, c.get(ctx, path, &out)
}

// ── API Groups ────────────────────────────────────────────────────────────────

func (c *Client) ListAPIGroups(ctx context.Context, orgID, serviceID string) ([]APIGroup, error) {
	var out struct {
		APIGroups []APIGroup `json:"apiGroups"`
	}
	return out.APIGroups, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups", orgID, serviceID), &out)
}

func (c *Client) GetAPIGroup(ctx context.Context, orgID, serviceID, id string) (*APIGroup, error) {
	var out APIGroup
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*APIGroup, error) {
	var out APIGroup
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups", orgID, serviceID), body, &out)
}

func (c *Client) UpdateAPIGroup(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*APIGroup, error) {
	var out APIGroup
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteAPIGroup(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s", orgID, serviceID, id))
}

func (c *Client) SyncAPIGroup(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	return out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/sync", orgID, serviceID), body, &out)
}

func (c *Client) ListAPIGroupVersions(ctx context.Context, orgID, serviceID, apiGroupID string) ([]APIGroupVersion, error) {
	var out struct {
		Versions []APIGroupVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/versions", orgID, serviceID, apiGroupID), &out)
}

// ── Service Docs ──────────────────────────────────────────────────────────────

func (c *Client) ListServiceDocs(ctx context.Context, orgID, serviceID string) ([]ServiceDoc, error) {
	var out struct {
		Docs []ServiceDoc `json:"docs"`
	}
	return out.Docs, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs", orgID, serviceID), &out)
}

func (c *Client) GetServiceDoc(ctx context.Context, orgID, serviceID, id string) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateServiceDoc(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs", orgID, serviceID), body, &out)
}

func (c *Client) UpdateServiceDoc(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*ServiceDoc, error) {
	var out ServiceDoc
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteServiceDoc(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/docs/%s", orgID, serviceID, id))
}

// ── Service Diagrams ──────────────────────────────────────────────────────────

func (c *Client) ListServiceDiagrams(ctx context.Context, orgID, serviceID string) ([]ServiceDiagram, error) {
	var out struct {
		Diagrams []ServiceDiagram `json:"diagrams"`
	}
	return out.Diagrams, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams", orgID, serviceID), &out)
}

func (c *Client) CreateServiceDiagram(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDiagram, error) {
	var out ServiceDiagram
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams", orgID, serviceID), body, &out)
}

func (c *Client) DeleteServiceDiagram(ctx context.Context, orgID, serviceID, diagramID string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/diagrams/%s", orgID, serviceID, diagramID))
}

// ── Service DBs ───────────────────────────────────────────────────────────────

func (c *Client) ListServiceDBs(ctx context.Context, orgID, serviceID string) ([]ServiceDB, error) {
	var out struct {
		DBs []ServiceDB `json:"dbs"`
	}
	return out.DBs, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs", orgID, serviceID), &out)
}

func (c *Client) GetServiceDB(ctx context.Context, orgID, serviceID, id string) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id), &out)
}

func (c *Client) CreateServiceDB(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs", orgID, serviceID), body, &out)
}

func (c *Client) UpdateServiceDB(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteServiceDB(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s", orgID, serviceID, id))
}

func (c *Client) ListServiceDBVersions(ctx context.Context, orgID, serviceID, serviceDBID string) ([]ServiceDBVersion, error) {
	var out struct {
		Versions []ServiceDBVersion `json:"versions"`
	}
	return out.Versions, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions", orgID, serviceID, serviceDBID), &out)
}

func (c *Client) CreateServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID string, body map[string]interface{}) (*ServiceDBVersion, error) {
	var out ServiceDBVersion
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions", orgID, serviceID, serviceDBID), body, &out)
}

func (c *Client) RestoreServiceDBVersion(ctx context.Context, orgID, serviceID, serviceDBID, versionID string) (*ServiceDB, error) {
	var out ServiceDB
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/dbs/%s/versions/%s/restore", orgID, serviceID, serviceDBID, versionID), nil, &out)
}

// ── API Endpoints ─────────────────────────────────────────────────────────────

func (c *Client) ListAPIEndpoints(ctx context.Context, orgID, serviceID, apiGroupID string) ([]APIEndpoint, error) {
	var out struct {
		Endpoints []APIEndpoint `json:"endpoints"`
	}
	return out.Endpoints, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints", orgID, serviceID, apiGroupID), &out)
}

func (c *Client) GetAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id), &out)
}

func (c *Client) CreateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID string, body map[string]interface{}) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints", orgID, serviceID, apiGroupID), body, &out)
}

func (c *Client) UpdateAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string, body map[string]interface{}) (*APIEndpoint, error) {
	var out APIEndpoint
	return &out, c.put(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id), body, &out)
}

func (c *Client) DeleteAPIEndpoint(ctx context.Context, orgID, serviceID, apiGroupID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/api-groups/%s/endpoints/%s", orgID, serviceID, apiGroupID, id))
}

// ── Test Packs ────────────────────────────────────────────────────────────────

func (c *Client) ListTestPacks(ctx context.Context, orgID, serviceID string) ([]TestPack, error) {
	var out struct {
		TestPacks []TestPack `json:"testPacks"`
	}
	return out.TestPacks, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-packs", orgID, serviceID), &out)
}

func (c *Client) CreateTestPack(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestPack, error) {
	var out TestPack
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-pack", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestPack(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestPack, error) {
	var out TestPack
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-pack/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteTestPack(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-pack/%s", orgID, serviceID, id))
}

// ── Test Cases ────────────────────────────────────────────────────────────────

func (c *Client) ListTestCases(ctx context.Context, orgID, serviceID string, testPackID *string) ([]TestCase, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-cases", orgID, serviceID)
	if testPackID != nil && *testPackID != "" {
		q := url.Values{}
		q.Set("testPackId", *testPackID)
		path += "?" + q.Encode()
	}
	var out struct {
		TestCases []TestCase `json:"testCases"`
	}
	return out.TestCases, c.get(ctx, path, &out)
}

func (c *Client) CreateTestCase(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestCase, error) {
	var out TestCase
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-case", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestCase(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestCase, error) {
	var out TestCase
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-case/%s", orgID, serviceID, id), body, &out)
}

func (c *Client) DeleteTestCase(ctx context.Context, orgID, serviceID, id string) error {
	return c.del(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-case/%s", orgID, serviceID, id))
}

// ── Test Runs ────────────────────────────────────────────────────────────────

func (c *Client) GetTestRun(ctx context.Context, orgID, serviceID, id string) (*TestRun, error) {
	var out TestRun
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run/%s", orgID, serviceID, id), &out)
}

func (c *Client) ListTestRuns(ctx context.Context, orgID, serviceID string, testPackID *string) ([]TestRun, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-runs", orgID, serviceID)
	if testPackID != nil && *testPackID != "" {
		q := url.Values{}
		q.Set("testPackId", *testPackID)
		path += "?" + q.Encode()
	}
	var out struct {
		TestRuns []TestRun `json:"testRuns"`
	}
	return out.TestRuns, c.get(ctx, path, &out)
}

func (c *Client) ListTestRunsSummary(
	ctx context.Context,
	orgID, serviceID string,
	testPackID *string,
	environment *string,
	status *string,
	executedBy *string,
	fromDate *time.Time,
	toDate *time.Time,
) ([]TestRunSummary, error) {
	q := url.Values{}
	if testPackID != nil && *testPackID != "" {
		q.Set("testPackId", *testPackID)
	}
	if environment != nil && *environment != "" {
		q.Set("environment", *environment)
	}
	if status != nil && *status != "" {
		q.Set("status", *status)
	}
	if executedBy != nil && *executedBy != "" {
		q.Set("executedBy", *executedBy)
	}
	if fromDate != nil {
		q.Set("fromDate", fromDate.UTC().Format(time.RFC3339))
	}
	if toDate != nil {
		q.Set("toDate", toDate.UTC().Format(time.RFC3339))
	}
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-runs-summary", orgID, serviceID)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out struct {
		TestRunsSummary []TestRunSummary `json:"testRunsSummary"`
	}
	return out.TestRunsSummary, c.get(ctx, path, &out)
}

func (c *Client) CreateTestRun(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestRun, error) {
	var out TestRun
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestRun(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestRun, error) {
	var out TestRun
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run/%s", orgID, serviceID, id), body, &out)
}

// ── Test Run Results ─────────────────────────────────────────────────────────

func (c *Client) ListTestRunResults(ctx context.Context, orgID, serviceID, testRunID string) ([]TestRunResult, error) {
	q := url.Values{}
	q.Set("testRunId", testRunID)
	path := fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run-results?%s", orgID, serviceID, q.Encode())
	var out struct {
		TestRunResults []TestRunResult `json:"testRunResults"`
	}
	return out.TestRunResults, c.get(ctx, path, &out)
}

func (c *Client) CreateTestRunResult(ctx context.Context, orgID, serviceID string, body map[string]interface{}) (*TestRunResult, error) {
	var out TestRunResult
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run-result", orgID, serviceID), body, &out)
}

func (c *Client) UpdateTestRunResult(ctx context.Context, orgID, serviceID, id string, body map[string]interface{}) (*TestRunResult, error) {
	var out TestRunResult
	return &out, c.post(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-run-result/%s", orgID, serviceID, id), body, &out)
}
