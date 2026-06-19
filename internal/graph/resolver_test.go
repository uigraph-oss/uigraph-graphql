package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

// doGraphQL sends a GraphQL query and returns the data map.
// It fails the test on HTTP errors or if GraphQL errors are present.
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

// doGraphQLRaw sends a GraphQL query and returns the raw decoded response
// (including any errors field), without failing on GraphQL errors.
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

// TestMeQuery exercises the me query through the real gqlgen executable schema
// using a fake authClient that returns a known MeResponse.
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

// TestMeQueryError verifies that an error from the authClient propagates as a
// GraphQL error in the response (not a 500).
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

// TestCreateFolderMutation exercises the createFolder mutation through the real
// schema using a fake folderClient.
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
