package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestFolderToModel(t *testing.T) {
	t.Run("maps all fields with parentID set", func(t *testing.T) {
		parentID := "parent-1"
		in := &uigraphapi.Folder{
			ID: "f1", OrgID: "o1", ParentID: &parentID, Type: "diagrams",
			Name: "My Folder", Order: 1.5, CreatedBy: "u1",
		}
		out := FolderToModel(in)
		if out.ID != "f1" {
			t.Errorf("ID = %q, want f1", out.ID)
		}
		if out.OrgID != "o1" {
			t.Errorf("OrgID = %q, want o1", out.OrgID)
		}
		if out.Name != "My Folder" {
			t.Errorf("Name = %q, want My Folder", out.Name)
		}
		if out.Order != 1.5 {
			t.Errorf("Order = %v, want 1.5", out.Order)
		}
		if out.Type != "diagrams" {
			t.Errorf("Type = %q, want diagrams", out.Type)
		}
		if out.CreatedBy != "u1" {
			t.Errorf("CreatedBy = %q, want u1", out.CreatedBy)
		}
		if out.ParentID == nil || *out.ParentID != parentID {
			t.Fatalf("ParentID = %v, want pointer to %q", out.ParentID, parentID)
		}
		if out.TeamID != nil {
			t.Fatalf("TeamID = %v, want nil", out.TeamID)
		}
	})

	t.Run("nil ParentID and TeamID when inputs are nil", func(t *testing.T) {
		in := &uigraphapi.Folder{ID: "f2", OrgID: "o1", Name: "Root"}
		out := FolderToModel(in)
		if out.ParentID != nil {
			t.Fatalf("ParentID = %v, want nil", out.ParentID)
		}
		if out.TeamID != nil {
			t.Fatalf("TeamID = %v, want nil", out.TeamID)
		}
	})

	t.Run("TeamID propagated when set", func(t *testing.T) {
		teamID := "team-99"
		in := &uigraphapi.Folder{ID: "f3", OrgID: "o1", Name: "Team Folder", TeamID: &teamID}
		out := FolderToModel(in)
		if out.TeamID == nil || *out.TeamID != teamID {
			t.Fatalf("TeamID = %v, want pointer to %q", out.TeamID, teamID)
		}
	})
}

func TestFoldersToModel(t *testing.T) {
	t.Run("converts slice preserving pointer fields", func(t *testing.T) {
		parentID := "p1"
		in := []uigraphapi.Folder{
			{ID: "f1", OrgID: "o1", Name: "A", ParentID: &parentID},
			{ID: "f2", OrgID: "o1", Name: "B"},
		}
		out := FoldersToModel(in)
		if len(out) != 2 {
			t.Fatalf("len = %d, want 2", len(out))
		}
		if out[0].ParentID == nil || *out[0].ParentID != "p1" {
			t.Errorf("out[0].ParentID = %v, want pointer to p1", out[0].ParentID)
		}
		if out[1].ParentID != nil {
			t.Errorf("out[1].ParentID = %v, want nil", out[1].ParentID)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		out := FoldersToModel([]uigraphapi.Folder{})
		if len(out) != 0 {
			t.Fatalf("len = %d, want 0", len(out))
		}
	})
}
