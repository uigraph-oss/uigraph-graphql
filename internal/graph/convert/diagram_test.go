package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestDiagramVersionToModel(t *testing.T) {
	t.Run("maps orgID parameter and version fields", func(t *testing.T) {
		out := DiagramVersionToModel("org-1", uigraphapi.DiagramVersion{
			ID: "v1", DiagramID: "d1", VersionNumber: 2,
		})
		if out.OrgID != "org-1" {
			t.Errorf("OrgID = %q, want org-1 (the orgID parameter, since DiagramVersion carries no OrgID field)", out.OrgID)
		}
		if out.ID != "v1" {
			t.Errorf("ID = %q, want v1", out.ID)
		}
		if out.DiagramID != "d1" {
			t.Errorf("DiagramID = %q, want d1", out.DiagramID)
		}
		if out.VersionNumber != 2 {
			t.Errorf("VersionNumber = %d, want 2", out.VersionNumber)
		}
	})

	t.Run("nil optional pointer fields remain nil", func(t *testing.T) {
		out := DiagramVersionToModel("o1", uigraphapi.DiagramVersion{ID: "v2"})
		if out.Label != nil {
			t.Errorf("Label = %v, want nil", out.Label)
		}
		if out.Source != nil {
			t.Errorf("Source = %v, want nil", out.Source)
		}
	})

	t.Run("optional pointer fields propagated when set", func(t *testing.T) {
		label := "v1.0"
		src := "ci"
		out := DiagramVersionToModel("o1", uigraphapi.DiagramVersion{
			ID: "v3", Label: &label, Source: &src, IsAutoVersion: true,
		})
		if out.Label == nil || *out.Label != label {
			t.Errorf("Label = %v, want pointer to %q", out.Label, label)
		}
		if out.Source == nil || *out.Source != src {
			t.Errorf("Source = %v, want pointer to %q", out.Source, src)
		}
		if !out.IsAutoVersion {
			t.Errorf("IsAutoVersion = false, want true")
		}
	})
}

func TestDiagramToModel(t *testing.T) {
	t.Run("propagates optional pointer fields when set", func(t *testing.T) {
		previewID := "asset-1"
		folderID := "folder-1"
		teamID := "team-1"
		source := "figma"
		out := DiagramToModel(&uigraphapi.Diagram{
			ID: "d1", OrgID: "o1", Name: "Checkout",
			FolderID: &folderID, TeamID: &teamID,
			PreviewAssetID: &previewID, Source: &source,
		})
		if out.ID != "d1" {
			t.Errorf("ID = %q, want d1", out.ID)
		}
		if out.OrgID != "o1" {
			t.Errorf("OrgID = %q, want o1", out.OrgID)
		}
		if out.Name != "Checkout" {
			t.Errorf("Name = %q, want Checkout", out.Name)
		}
		if out.PreviewAssetID == nil || *out.PreviewAssetID != previewID {
			t.Errorf("PreviewAssetID = %v, want pointer to %q", out.PreviewAssetID, previewID)
		}
		if out.FolderID == nil || *out.FolderID != folderID {
			t.Errorf("FolderID = %v, want pointer to %q", out.FolderID, folderID)
		}
		if out.TeamID == nil || *out.TeamID != teamID {
			t.Errorf("TeamID = %v, want pointer to %q", out.TeamID, teamID)
		}
		if out.Source == nil || *out.Source != source {
			t.Errorf("Source = %v, want pointer to %q", out.Source, source)
		}
	})

	t.Run("nil optional pointer fields when source struct has no values", func(t *testing.T) {
		out := DiagramToModel(&uigraphapi.Diagram{ID: "d2", OrgID: "o1", Name: "Empty"})
		if out.FolderID != nil {
			t.Errorf("FolderID = %v, want nil", out.FolderID)
		}
		if out.TeamID != nil {
			t.Errorf("TeamID = %v, want nil", out.TeamID)
		}
		if out.PreviewAssetID != nil {
			t.Errorf("PreviewAssetID = %v, want nil", out.PreviewAssetID)
		}
		if out.Source != nil {
			t.Errorf("Source = %v, want nil", out.Source)
		}
		if out.UpdatedBy != nil {
			t.Errorf("UpdatedBy = %v, want nil", out.UpdatedBy)
		}
	})
}

func TestDiagramsToModel(t *testing.T) {
	t.Run("converts slice preserving pointer fields", func(t *testing.T) {
		previewID := "prev-1"
		in := []uigraphapi.Diagram{
			{ID: "d1", OrgID: "o1", Name: "A", PreviewAssetID: &previewID},
			{ID: "d2", OrgID: "o1", Name: "B"},
		}
		out := DiagramsToModel(in)
		if len(out) != 2 {
			t.Fatalf("len = %d, want 2", len(out))
		}
		if out[0].PreviewAssetID == nil || *out[0].PreviewAssetID != "prev-1" {
			t.Errorf("out[0].PreviewAssetID = %v, want pointer to prev-1", out[0].PreviewAssetID)
		}
		if out[1].PreviewAssetID != nil {
			t.Errorf("out[1].PreviewAssetID = %v, want nil", out[1].PreviewAssetID)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		out := DiagramsToModel([]uigraphapi.Diagram{})
		if len(out) != 0 {
			t.Fatalf("len = %d, want 0", len(out))
		}
	})
}
