package uigraphapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Assertion struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AuthConfig struct {
	Type          string  `json:"type"`
	BearerToken   *string `json:"bearerToken,omitempty"`
	APIKeyHeader  *string `json:"apiKeyHeader,omitempty"`
	APIKeyValue   *string `json:"apiKeyValue,omitempty"`
	BasicUsername *string `json:"basicUsername,omitempty"`
	BasicPassword *string `json:"basicPassword,omitempty"`
}

type TestCaseStep struct {
	Order          int    `json:"order"`
	Action         string `json:"action"`
	ExpectedResult string `json:"expectedResult"`
}

type ManualTestCase struct {
	Preconditions   *string        `json:"preconditions,omitempty"`
	TestData        *string        `json:"testData,omitempty"`
	Steps           []TestCaseStep `json:"steps,omitempty"`
	ExpectedOutcome *string        `json:"expectedOutcome,omitempty"`
	Postconditions  *string        `json:"postconditions,omitempty"`
}

type APITestCase struct {
	HTTPMethod         string      `json:"httpMethod"`
	APISpecID          *string     `json:"apiSpecId,omitempty"`
	OperationID        *string     `json:"operationId,omitempty"`
	Auth               *AuthConfig `json:"auth,omitempty"`
	RequestHeaders     []KeyValue  `json:"requestHeaders,omitempty"`
	QueryParams        []KeyValue  `json:"queryParams,omitempty"`
	RequestBody        *string     `json:"requestBody,omitempty"`
	ExpectedStatusCode *int        `json:"expectedStatusCode,omitempty"`
	MaxResponseTimeMs  *int        `json:"maxResponseTimeMs,omitempty"`
	ResponseBody       *string     `json:"responseBody,omitempty"`
	Assertions         []Assertion `json:"assertions,omitempty"`
}

type GraphQLTestCase struct {
	OperationType string      `json:"operationType"`
	OperationName *string     `json:"operationName,omitempty"`
	Query         string      `json:"query"`
	Variables     *string     `json:"variables,omitempty"`
	ResponseBody  *string     `json:"responseBody,omitempty"`
	Assertions    []Assertion `json:"assertions,omitempty"`
	ExpectError   bool        `json:"expectError"`
}

type DatabaseTestCase struct {
	Dialect       string      `json:"dialect"`
	SchemaID      *string     `json:"schemaId,omitempty"`
	Query         string      `json:"query"`
	Assertions    []Assertion `json:"assertions,omitempty"`
	SetupQuery    *string     `json:"setupQuery,omitempty"`
	TeardownQuery *string     `json:"teardownQuery,omitempty"`
}

type GRPCTestCase struct {
	ServiceName    string      `json:"serviceName"`
	MethodName     string      `json:"methodName"`
	CallMode       string      `json:"callMode"`
	ProtoFileID    *string     `json:"protoFileId,omitempty"`
	ServerAddress  *string     `json:"serverAddress,omitempty"`
	RequestMessage *string     `json:"requestMessage,omitempty"`
	Metadata       []KeyValue  `json:"metadata,omitempty"`
	ExpectedStatus string      `json:"expectedStatus"`
	DeadlineMs     *int        `json:"deadlineMs,omitempty"`
	ResponseBody   *string     `json:"responseBody,omitempty"`
	Assertions     []Assertion `json:"assertions,omitempty"`
	UseTLS         bool        `json:"useTLS"`
	ExpectError    bool        `json:"expectError"`
}

type TestPack struct {
	TestPackID          string     `json:"testPackId"`
	ServiceID           string     `json:"serviceId"`
	OrgID               string     `json:"orgId"`
	Name                string     `json:"name"`
	Type                string     `json:"type"`
	CreatedBy           string     `json:"createdBy"`
	UpdatedBy           *string    `json:"updatedBy,omitempty"`
	CreatedByCommitHash *string    `json:"createdByCommitHash,omitempty"`
	UpdatedByCommitHash *string    `json:"updatedByCommitHash,omitempty"`
	DeletedBy           *string    `json:"deletedBy,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	DeletedAt           *time.Time `json:"deletedAt,omitempty"`
}

type TestCase struct {
	TestCaseID            string            `json:"testCaseId"`
	TestPackID            string            `json:"testPackId"`
	ServiceID             string            `json:"serviceId"`
	OrgID                 string            `json:"orgId"`
	Title                 string            `json:"title"`
	Order                 float64           `json:"order"`
	Type                  string            `json:"type"`
	Description           *string           `json:"description,omitempty"`
	Priority              *string           `json:"priority,omitempty"`
	Labels                []string          `json:"labels,omitempty"`
	LinkedTicket          *string           `json:"linkedTicket,omitempty"`
	EstimatedDurationMins *int              `json:"estimatedDurationMins,omitempty"`
	TestOwner             *string           `json:"testOwner,omitempty"`
	LinkedMapNodeID       *string           `json:"linkedMapNodeId,omitempty"`
	IsCritical            bool              `json:"isCritical"`
	EvidenceRequired      bool              `json:"evidenceRequired"`
	Manual                *ManualTestCase   `json:"manual,omitempty"`
	API                   *APITestCase      `json:"api,omitempty"`
	GraphQL               *GraphQLTestCase  `json:"graphql,omitempty"`
	Database              *DatabaseTestCase `json:"database,omitempty"`
	GRPC                  *GRPCTestCase     `json:"grpc,omitempty"`
	Status                string            `json:"status"`
	Version               int               `json:"version"`
	BaselineRunResultID   *string           `json:"baselineRunResultId,omitempty"`
	Dependencies          []string          `json:"dependencies,omitempty"`
	CreatedBy             string            `json:"createdBy"`
	UpdatedBy             *string           `json:"updatedBy,omitempty"`
	CreatedByCommitHash   *string           `json:"createdByCommitHash,omitempty"`
	UpdatedByCommitHash   *string           `json:"updatedByCommitHash,omitempty"`
	DeletedBy             *string           `json:"deletedBy,omitempty"`
	CreatedAt             time.Time         `json:"createdAt"`
	UpdatedAt             time.Time         `json:"updatedAt"`
	DeletedAt             *time.Time        `json:"deletedAt,omitempty"`
}

type TestRun struct {
	TestRunID     string     `json:"testRunId"`
	TestPackID    string     `json:"testPackId"`
	ServiceID     string     `json:"serviceId"`
	OrgID         string     `json:"orgId"`
	Environment   string     `json:"environment"`
	ReleaseLabel  *string    `json:"releaseLabel,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	Status        string     `json:"status"`
	StartedBy     *string    `json:"startedBy,omitempty"`
	ExecutedBy    string     `json:"executedBy"`
	ExecutedAt    time.Time  `json:"executedAt"`
	OverallStatus string     `json:"overallStatus"`
}

type TestRunSummary struct {
	TestRunID     string     `json:"testRunId"`
	TestPackID    string     `json:"testPackId"`
	ServiceID     string     `json:"serviceId"`
	Environment   string     `json:"environment"`
	ReleaseLabel  *string    `json:"releaseLabel,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	Status        string     `json:"status"`
	StartedBy     *string    `json:"startedBy,omitempty"`
	ExecutedBy    string     `json:"executedBy"`
	ExecutedAt    time.Time  `json:"executedAt"`
	OverallStatus string     `json:"overallStatus"`
	PassedCount   int        `json:"passedCount"`
	FailedCount   int        `json:"failedCount"`
	SkippedCount  int        `json:"skippedCount"`
	BlockedCount  int        `json:"blockedCount"`
}

type TestRunResult struct {
	TestRunResultID string    `json:"testRunResultId"`
	TestRunID       string    `json:"testRunId"`
	TestCaseID      string    `json:"testCaseId"`
	ServiceID       string    `json:"serviceId"`
	OrgID           string    `json:"orgId"`
	Status          string    `json:"status"`
	BlockedReason   *string   `json:"blockedReason,omitempty"`
	ResponseStatus  *int      `json:"responseStatus,omitempty"`
	ResponseBody    *string   `json:"responseBody,omitempty"`
	ResponseTimeMs  *int64    `json:"responseTimeMs,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	ScreenshotURLs  []string  `json:"screenshotUrls,omitempty"`
	ExecutedAt      time.Time `json:"executedAt"`
	ExecutedBy      string    `json:"executedBy"`
}

func (c *Client) ListTestPacks(ctx context.Context, orgID, serviceID string) ([]TestPack, error) {
	var out struct {
		TestPacks []TestPack `json:"testPacks"`
	}
	return out.TestPacks, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/services/%s/test-packs", orgID, serviceID), &out)
}

func (c *Client) GetTestPackByID(ctx context.Context, orgID, id string) (*TestPack, error) {
	var out TestPack
	return &out, c.get(ctx, fmt.Sprintf("/api/v1/orgs/%s/test-packs/%s", orgID, id), &out)
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
