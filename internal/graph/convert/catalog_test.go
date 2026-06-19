package convert

import (
	"encoding/json"
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestServiceToModel(t *testing.T) {
	t.Run("populated Metadata passes through as string", func(t *testing.T) {
		out := ServiceToModel(&uigraphapi.Service{
			ID: "s1", OrgID: "o1", Name: "Checkout",
			Metadata: json.RawMessage(`{"team":"core"}`),
		})
		if out.ID != "s1" {
			t.Errorf("ID = %q, want s1", out.ID)
		}
		if out.OrgID != "o1" {
			t.Errorf("OrgID = %q, want o1", out.OrgID)
		}
		if out.Name != "Checkout" {
			t.Errorf("Name = %q, want Checkout", out.Name)
		}
		if out.Metadata != `{"team":"core"}` {
			t.Errorf("Metadata = %q, want passthrough", out.Metadata)
		}
	})

	t.Run("empty Metadata defaults to empty object string", func(t *testing.T) {
		out := ServiceToModel(&uigraphapi.Service{ID: "s2"})
		if out.Metadata != "{}" {
			t.Errorf("Metadata = %q, want %q for empty RawMessage", out.Metadata, "{}")
		}
	})

	t.Run("all string fields mapped correctly", func(t *testing.T) {
		folderID := "folder-1"
		teamID := "team-1"
		gitURL := "https://github.com/org/repo"
		out := ServiceToModel(&uigraphapi.Service{
			ID: "s3", OrgID: "o1", FolderID: &folderID, TeamID: &teamID,
			Name: "API", Slug: "api", Description: "core api",
			Status: "active", Tier: "tier-1", Category: "backend",
			Language: "go", GitRepoURL: &gitURL,
		})
		if out.Slug != "api" {
			t.Errorf("Slug = %q, want api", out.Slug)
		}
		if out.Description != "core api" {
			t.Errorf("Description = %q, want core api", out.Description)
		}
		if out.Status != "active" {
			t.Errorf("Status = %q, want active", out.Status)
		}
		if out.FolderID == nil || *out.FolderID != folderID {
			t.Errorf("FolderID = %v, want pointer to %q", out.FolderID, folderID)
		}
		if out.TeamID == nil || *out.TeamID != teamID {
			t.Errorf("TeamID = %v, want pointer to %q", out.TeamID, teamID)
		}
	})
}

func TestServiceDiagramToModel(t *testing.T) {
	t.Run("nested Diagram converted when non-nil", func(t *testing.T) {
		out := ServiceDiagramToModel(&uigraphapi.ServiceDiagram{
			ServiceID: "svc1", DiagramID: "d1",
			Diagram: &uigraphapi.Diagram{ID: "d1", Name: "Checkout"},
		})
		if out.ServiceID != "svc1" {
			t.Errorf("ServiceID = %q, want svc1", out.ServiceID)
		}
		if out.DiagramID != "d1" {
			t.Errorf("DiagramID = %q, want d1", out.DiagramID)
		}
		if out.Diagram == nil {
			t.Fatal("Diagram = nil, want a converted Diagram")
		}
		if out.Diagram.Name != "Checkout" {
			t.Errorf("Diagram.Name = %q, want Checkout", out.Diagram.Name)
		}
	})

	t.Run("nested Diagram is nil when source Diagram pointer is nil", func(t *testing.T) {
		out := ServiceDiagramToModel(&uigraphapi.ServiceDiagram{ServiceID: "svc1", DiagramID: "d2"})
		if out.Diagram != nil {
			t.Errorf("Diagram = %v, want nil when source Diagram pointer is nil", out.Diagram)
		}
	})
}

func TestAPIGroupToModel(t *testing.T) {
	t.Run("maps all primitive fields", func(t *testing.T) {
		label := "stable"
		specKey := "spec/key"
		specHash := "abc123"
		out := APIGroupToModel(&uigraphapi.APIGroup{
			ID: "ag1", ServiceID: "svc1", OrgID: "o1",
			Name: "User API", Version: "v1", Label: &label,
			Protocol: "REST", SpecKey: &specKey, SpecHash: &specHash,
		})
		if out.ID != "ag1" {
			t.Errorf("ID = %q, want ag1", out.ID)
		}
		if out.ServiceID != "svc1" {
			t.Errorf("ServiceID = %q, want svc1", out.ServiceID)
		}
		if out.Name != "User API" {
			t.Errorf("Name = %q, want User API", out.Name)
		}
		if out.Version != "v1" {
			t.Errorf("Version = %q, want v1", out.Version)
		}
		if out.Protocol != "REST" {
			t.Errorf("Protocol = %q, want REST", out.Protocol)
		}
		if out.SpecKey == nil || *out.SpecKey != "spec/key" {
			t.Errorf("SpecKey = %v, want pointer to spec/key", out.SpecKey)
		}
		if out.SpecHash == nil || *out.SpecHash != "abc123" {
			t.Errorf("SpecHash = %v, want pointer to abc123", out.SpecHash)
		}
		if out.Label == nil || *out.Label != "stable" {
			t.Errorf("Label = %v, want pointer to stable", out.Label)
		}
	})

	t.Run("nil optional fields when source has none", func(t *testing.T) {
		out := APIGroupToModel(&uigraphapi.APIGroup{ID: "ag2", ServiceID: "svc1", OrgID: "o1"})
		if out.Label != nil {
			t.Errorf("Label = %v, want nil", out.Label)
		}
		if out.SpecKey != nil {
			t.Errorf("SpecKey = %v, want nil", out.SpecKey)
		}
		if out.SpecHash != nil {
			t.Errorf("SpecHash = %v, want nil", out.SpecHash)
		}
	})
}

func TestServicesToModel(t *testing.T) {
	t.Run("converts slice", func(t *testing.T) {
		in := []uigraphapi.Service{
			{ID: "s1", OrgID: "o1", Name: "A"},
			{ID: "s2", OrgID: "o1", Name: "B"},
		}
		out := ServicesToModel(in)
		if len(out) != 2 {
			t.Fatalf("len = %d, want 2", len(out))
		}
		if out[0].ID != "s1" {
			t.Errorf("out[0].ID = %q, want s1", out[0].ID)
		}
		if out[1].ID != "s2" {
			t.Errorf("out[1].ID = %q, want s2", out[1].ID)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		out := ServicesToModel([]uigraphapi.Service{})
		if len(out) != 0 {
			t.Fatalf("len = %d, want 0", len(out))
		}
	})
}
