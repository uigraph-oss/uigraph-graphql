package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/uigraph/graphql/internal/graph"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

var errUnauthorized = errors.New("unauthorized")

type fakeAuthClient struct {
	me  *uigraphapi.MeResponse
	err error
}

func (f *fakeAuthClient) Me(_ context.Context) (*uigraphapi.MeResponse, error) {
	return f.me, f.err
}
func (f *fakeAuthClient) MyOrgs(_ context.Context) ([]uigraphapi.OrgSummary, error) {
	return nil, nil
}
func (f *fakeAuthClient) SwitchOrg(_ context.Context, _ string) error { return nil }
func (f *fakeAuthClient) PrepareUserAvatarUpload(_ context.Context) (*uigraphapi.AssetUpload, error) {
	return nil, nil
}
func (f *fakeAuthClient) SetMyAvatar(_ context.Context) error { return nil }

type fakeFolderClient struct {
	created *uigraphapi.Folder
}

func (f *fakeFolderClient) ListFolders(_ context.Context, _, _, _ string) ([]uigraphapi.Folder, error) {
	return nil, nil
}
func (f *fakeFolderClient) GetFolder(_ context.Context, _, _ string) (*uigraphapi.Folder, error) {
	return nil, nil
}
func (f *fakeFolderClient) CreateFolder(_ context.Context, orgID string, body map[string]interface{}) (*uigraphapi.Folder, error) {
	f.created = &uigraphapi.Folder{
		ID:    "folder-1",
		OrgID: orgID,
		Name:  body["name"].(string),
		Type:  body["type"].(string),
	}
	return f.created, nil
}
func (f *fakeFolderClient) UpdateFolder(_ context.Context, _, _ string, _ map[string]interface{}) (*uigraphapi.Folder, error) {
	return nil, nil
}
func (f *fakeFolderClient) DeleteFolder(_ context.Context, _, _ string) error { return nil }

func newTestServer(resolver *graph.Resolver) *httptest.Server {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	srv := handler.New(schema)
	srv.AddTransport(transport.POST{})
	return httptest.NewServer(srv)
}

func doGraphQL(t *testing.T, srv *httptest.Server, query string) map[string]interface{} {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /graphql: %v", err)
	}
	defer resp.Body.Close()

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if errs, ok := out["errors"]; ok {
		t.Fatalf("graphql errors: %v", errs)
	}
	return out["data"].(map[string]interface{})
}

func doGraphQLRaw(t *testing.T, srv *httptest.Server, query string) map[string]interface{} {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /graphql: %v", err)
	}
	defer resp.Body.Close()

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return out
}

func TestMeQuery(t *testing.T) {
	resolver := &graph.Resolver{
		Auth: &fakeAuthClient{
			me: &uigraphapi.MeResponse{
				UserID: "u1",
				Email:  "a@b.com",
				Name:   "Ann",
				Role:   "admin",
			},
		},
	}
	srv := newTestServer(resolver)
	defer srv.Close()

	data := doGraphQL(t, srv, `{ me { userId email } }`)
	me, ok := data["me"].(map[string]interface{})
	if !ok {
		t.Fatalf("me field missing or wrong type: %+v", data)
	}
	if me["userId"] != "u1" {
		t.Errorf("userId: got %q, want %q", me["userId"], "u1")
	}
	if me["email"] != "a@b.com" {
		t.Errorf("email: got %q, want %q", me["email"], "a@b.com")
	}
}

func TestMeQueryError(t *testing.T) {
	resolver := &graph.Resolver{
		Auth: &fakeAuthClient{err: errUnauthorized},
	}
	srv := newTestServer(resolver)
	defer srv.Close()

	out := doGraphQLRaw(t, srv, `{ me { userId email } }`)
	errs, ok := out["errors"]
	if !ok {
		t.Fatalf("expected GraphQL errors but got none; full response: %+v", out)
	}
	errList, ok := errs.([]interface{})
	if !ok || len(errList) == 0 {
		t.Fatalf("errors field is not a non-empty array: %+v", errs)
	}
}

func TestCreateFolderMutation(t *testing.T) {
	folders := &fakeFolderClient{}
	resolver := &graph.Resolver{FolderAPI: folders}
	srv := newTestServer(resolver)
	defer srv.Close()

	data := doGraphQL(t, srv, `mutation {
		createFolder(orgId: "org-1", input: { name: "Diagrams", type: "diagram" }) {
			id name type orgId
		}
	}`)
	created, ok := data["createFolder"].(map[string]interface{})
	if !ok {
		t.Fatalf("createFolder field missing or wrong type: %+v", data)
	}
	if created["name"] != "Diagrams" {
		t.Errorf("name: got %q, want %q", created["name"], "Diagrams")
	}
	if created["type"] != "diagram" {
		t.Errorf("type: got %q, want %q", created["type"], "diagram")
	}
	if created["orgId"] != "org-1" {
		t.Errorf("orgId: got %q, want %q", created["orgId"], "org-1")
	}
	if folders.created == nil {
		t.Fatal("fakeFolderClient.created is nil — CreateFolder was not called")
	}
	if folders.created.OrgID != "org-1" {
		t.Errorf("fakeFolderClient.created.OrgID: got %q, want %q", folders.created.OrgID, "org-1")
	}
}

type fakeDependencyClient struct {
	direction   *string
	criticality *string
	maxDepth    *int
}

func (f *fakeDependencyClient) ListDependencies(_ context.Context, _, _ string, direction, criticality *string) ([]uigraphapi.Dependency, error) {
	f.direction = direction
	f.criticality = criticality
	return []uigraphapi.Dependency{{
		ID:              "dependency-1",
		Name:            "Payments",
		ConsumerService: uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"},
	}}, nil
}

func (f *fakeDependencyClient) GetServiceDependencyGraph(_ context.Context, _, _ string) ([]uigraphapi.Dependency, error) {
	return []uigraphapi.Dependency{{ID: "dependency-1", Name: "Payments", ConsumerService: uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"}}}, nil
}

func (f *fakeDependencyClient) GetDependencyGraph(_ context.Context, _ string) ([]uigraphapi.Dependency, error) {
	return []uigraphapi.Dependency{{ID: "dependency-1", Name: "Payments", ConsumerService: uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"}}}, nil
}

func (f *fakeDependencyClient) UpdateServiceDependencies(_ context.Context, _, _ string, _ map[string]interface{}) ([]uigraphapi.Dependency, error) {
	return []uigraphapi.Dependency{{ID: "dependency-1", Name: "Payments", ConsumerService: uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"}}}, nil
}

func (f *fakeDependencyClient) GetServiceImpact(_ context.Context, _, _ string, direction *string, maxDepth *int) ([]uigraphapi.Dependency, error) {
	f.direction = direction
	f.maxDepth = maxDepth
	return []uigraphapi.Dependency{{ID: "dependency-1", Name: "Payments", ConsumerService: uigraphapi.DependencyService{ID: "service-1", Name: "Checkout"}}}, nil
}

func TestDependencyQueries(t *testing.T) {
	dependencies := &fakeDependencyClient{}
	srv := newTestServer(&graph.Resolver{Dependency: dependencies})
	defer srv.Close()

	data := doGraphQL(t, srv, `{ dependencies(orgId: "org-1", serviceId: "service-1", direction: "outbound", criticality: "high") { id name consumerService { id name } } }`)
	items, ok := data["dependencies"].([]interface{})
	if !ok || len(items) != 1 {
		t.Fatalf("dependencies = %#v, want one item", data["dependencies"])
	}
	if dependencies.direction == nil || *dependencies.direction != "outbound" {
		t.Errorf("direction = %v, want outbound", dependencies.direction)
	}
	if dependencies.criticality == nil || *dependencies.criticality != "high" {
		t.Errorf("criticality = %v, want high", dependencies.criticality)
	}

	data = doGraphQL(t, srv, `{ serviceImpact(orgId: "org-1", serviceId: "service-1", direction: "inbound", maxDepth: 2) { id name consumerService { id name } } }`)
	impact, ok := data["serviceImpact"].([]interface{})
	if !ok || len(impact) != 1 {
		t.Fatalf("serviceImpact = %#v, want one dependency", data["serviceImpact"])
	}
	if dependencies.maxDepth == nil || *dependencies.maxDepth != 2 {
		t.Errorf("maxDepth = %v, want 2", dependencies.maxDepth)
	}
}

type fakeDiagramClient struct {
	prepareThumbnailFn func(ctx context.Context, orgID, diagramID string) (*uigraphapi.DiagramThumbnailUpload, error)
	confirmThumbnailFn func(ctx context.Context, orgID, diagramID, hash string) error
}

func (f *fakeDiagramClient) ListDiagrams(_ context.Context, _ string, _ uigraphapi.ListParams) ([]uigraphapi.Diagram, int, error) {
	return nil, 0, nil
}
func (f *fakeDiagramClient) GetDiagram(_ context.Context, _, _ string) (*uigraphapi.Diagram, error) {
	return nil, nil
}
func (f *fakeDiagramClient) GetDiagramContent(_ context.Context, _, _ string) (string, error) {
	return "", nil
}
func (f *fakeDiagramClient) CreateDiagram(_ context.Context, _ string, _ map[string]interface{}) (*uigraphapi.Diagram, error) {
	return nil, nil
}
func (f *fakeDiagramClient) UpdateDiagram(_ context.Context, _, _ string, _ map[string]interface{}) (*uigraphapi.Diagram, error) {
	return nil, nil
}
func (f *fakeDiagramClient) DeleteDiagram(_ context.Context, _, _ string) error { return nil }
func (f *fakeDiagramClient) ListDiagramImages(_ context.Context, _, _ string) ([]uigraphapi.DiagramImage, error) {
	return nil, nil
}
func (f *fakeDiagramClient) CreateDiagramImage(_ context.Context, _, _ string, _ map[string]interface{}) (*uigraphapi.DiagramImage, error) {
	return nil, nil
}
func (f *fakeDiagramClient) SyncDiagram(_ context.Context, _ string, _ map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
func (f *fakeDiagramClient) ListDiagramVersions(_ context.Context, _, _ string) ([]uigraphapi.DiagramVersion, error) {
	return nil, nil
}
func (f *fakeDiagramClient) CreateDiagramVersion(_ context.Context, _, _ string, _ map[string]interface{}) (*uigraphapi.DiagramVersion, error) {
	return nil, nil
}
func (f *fakeDiagramClient) GetDiagramVersionContent(_ context.Context, _, _, _ string) (string, error) {
	return "", nil
}
func (f *fakeDiagramClient) RestoreDiagramVersion(_ context.Context, _, _, _ string) (*uigraphapi.Diagram, error) {
	return nil, nil
}
func (f *fakeDiagramClient) PrepareDiagramThumbnailUpload(ctx context.Context, orgID, diagramID string) (*uigraphapi.DiagramThumbnailUpload, error) {
	if f.prepareThumbnailFn != nil {
		return f.prepareThumbnailFn(ctx, orgID, diagramID)
	}
	return &uigraphapi.DiagramThumbnailUpload{UploadURL: "https://storage.example.com/put", AssetID: "diagram_" + diagramID}, nil
}
func (f *fakeDiagramClient) ConfirmDiagramThumbnailUpload(ctx context.Context, orgID, diagramID, hash string) error {
	if f.confirmThumbnailFn != nil {
		return f.confirmThumbnailFn(ctx, orgID, diagramID, hash)
	}
	return nil
}

func TestPrepareDiagramThumbnailUpload_returnsUploadURL(t *testing.T) {
	dc := &fakeDiagramClient{}
	r := &graph.Resolver{DiagramAPI: dc}
	srv := newTestServer(r)
	defer srv.Close()

	body := `{"query":"mutation { prepareDiagramThumbnailUpload(orgId:\"org-1\", diagramId:\"d1\") { uploadUrl assetId } }"}`
	resp, err := http.Post(srv.URL+"/query", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			PrepareDiagramThumbnailUpload struct {
				UploadURL string `json:"uploadUrl"`
				AssetID   string `json:"assetId"`
			} `json:"prepareDiagramThumbnailUpload"`
		} `json:"data"`
		Errors []struct{ Message string } `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result.Errors) > 0 {
		t.Fatalf("graphql errors: %v", result.Errors)
	}
	if result.Data.PrepareDiagramThumbnailUpload.UploadURL == "" {
		t.Fatal("expected uploadUrl")
	}
}

func TestConfirmDiagramThumbnailUpload_returnsTrue(t *testing.T) {
	dc := &fakeDiagramClient{}
	r := &graph.Resolver{DiagramAPI: dc}
	srv := newTestServer(r)
	defer srv.Close()

	body := `{"query":"mutation { confirmDiagramThumbnailUpload(orgId:\"org-1\", diagramId:\"d1\", contentHash:\"abc\") }"}`
	resp, err := http.Post(srv.URL+"/query", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ConfirmDiagramThumbnailUpload bool `json:"confirmDiagramThumbnailUpload"`
		} `json:"data"`
		Errors []struct{ Message string } `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result.Errors) > 0 {
		t.Fatalf("graphql errors: %v", result.Errors)
	}
	if !result.Data.ConfirmDiagramThumbnailUpload {
		t.Fatal("expected true")
	}
}
