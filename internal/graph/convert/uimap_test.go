package convert

import (
	"encoding/json"
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestParseServiceDBSchema(t *testing.T) {
	t.Run("parses tables from schema json object", func(t *testing.T) {
		p := parseServiceDBSchema(json.RawMessage(`{"tables":[{"name":"users"}],"dbDiagramId":"d1"}`))
		if len(p.Tables) != 1 || p.Tables[0].Name == nil || *p.Tables[0].Name != "users" {
			t.Errorf("Tables = %#v, want one table named users", p.Tables)
		}
		if p.DbDiagramID == nil || *p.DbDiagramID != "d1" {
			t.Errorf("DbDiagramID = %v, want d1", p.DbDiagramID)
		}
	})

	t.Run("handles double-encoded json string", func(t *testing.T) {
		p := parseServiceDBSchema(json.RawMessage(`"{\"tables\":[{\"name\":\"orders\"}]}"`))
		if len(p.Tables) != 1 || p.Tables[0].Name == nil || *p.Tables[0].Name != "orders" {
			t.Errorf("Tables = %#v, want one table named orders", p.Tables)
		}
	})

	t.Run("empty raw yields empty tables slice", func(t *testing.T) {
		p := parseServiceDBSchema(nil)
		if p.Tables == nil || len(p.Tables) != 0 {
			t.Errorf("Tables = %#v, want empty slice", p.Tables)
		}
	})
}

func TestCanvasToModel(t *testing.T) {
	t.Run("empty FramePositions defaults to empty object string", func(t *testing.T) {
		out := CanvasToModel(&uigraphapi.Canvas{MapID: "m1", OrgID: "o1"})
		if out.MapID != "m1" {
			t.Errorf("MapID = %q, want m1", out.MapID)
		}
		if out.FramePositions != "{}" {
			t.Errorf("FramePositions = %q, want %q for empty RawMessage", out.FramePositions, "{}")
		}
	})

	t.Run("populated FramePositions passes through as string", func(t *testing.T) {
		out := CanvasToModel(&uigraphapi.Canvas{
			MapID: "m1", OrgID: "o1",
			FramePositions: json.RawMessage(`{"frame-1":{"x":10,"y":20}}`),
		})
		if out.FramePositions != `{"frame-1":{"x":10,"y":20}}` {
			t.Errorf("FramePositions = %q, want passthrough of input JSON", out.FramePositions)
		}
	})

	t.Run("numeric fields passed through", func(t *testing.T) {
		out := CanvasToModel(&uigraphapi.Canvas{MapID: "m2", OrgID: "o1", Zoom: 1.5, NavigationX: 100, NavigationY: 200})
		if out.Zoom != 1.5 {
			t.Errorf("Zoom = %v, want 1.5", out.Zoom)
		}
		if out.NavigationX != 100 {
			t.Errorf("NavigationX = %v, want 100", out.NavigationX)
		}
		if out.NavigationY != 200 {
			t.Errorf("NavigationY = %v, want 200", out.NavigationY)
		}
	})
}

func TestFocalPointMetaToModel(t *testing.T) {
	t.Run("populated ComponentImages parsed into string slice", func(t *testing.T) {
		out := FocalPointMetaToModel(&uigraphapi.FocalPointMeta{
			ID: "fpm1", ComponentImages: json.RawMessage(`["a.png"]`),
		})
		if out.ID != "fpm1" {
			t.Errorf("ID = %q, want fpm1", out.ID)
		}
		if len(out.ComponentImages) != 1 || out.ComponentImages[0] != "a.png" {
			t.Errorf("ComponentImages = %#v, want [a.png]", out.ComponentImages)
		}
	})

	t.Run("empty ComponentImages defaults to empty slice", func(t *testing.T) {
		out := FocalPointMetaToModel(&uigraphapi.FocalPointMeta{ID: "fpm2"})
		if out.ComponentImages == nil || len(out.ComponentImages) != 0 {
			t.Errorf("ComponentImages = %#v, want empty slice", out.ComponentImages)
		}
	})

	t.Run("populated ComponentModalFields parsed into structured fields", func(t *testing.T) {
		out := FocalPointMetaToModel(&uigraphapi.FocalPointMeta{
			ID:                   "fpm3",
			ComponentModalFields: json.RawMessage(`[{"componentFieldId":"foo","label":"Foo","order":2}]`),
		})
		if len(out.ComponentModalFields) != 1 {
			t.Fatalf("ComponentModalFields = %#v, want one element", out.ComponentModalFields)
		}
		f := out.ComponentModalFields[0]
		if f.ComponentFieldID == nil || *f.ComponentFieldID != "foo" || f.Order == nil || *f.Order != 2 {
			t.Errorf("ComponentModalFields[0] = %#v, want {foo, order 2}", f)
		}
	})

	t.Run("empty ComponentModalFields defaults to empty slice", func(t *testing.T) {
		out := FocalPointMetaToModel(&uigraphapi.FocalPointMeta{ID: "fpm4"})
		if out.ComponentModalFields == nil || len(out.ComponentModalFields) != 0 {
			t.Errorf("ComponentModalFields = %#v, want empty slice", out.ComponentModalFields)
		}
	})

	t.Run("optional pointer fields nil when source has none", func(t *testing.T) {
		out := FocalPointMetaToModel(&uigraphapi.FocalPointMeta{ID: "fpm5", OrgID: "o1"})
		if out.ComponentLinkID != nil {
			t.Errorf("ComponentLinkID = %v, want nil", out.ComponentLinkID)
		}
		if out.ComponentFlowDiagram != nil {
			t.Errorf("ComponentFlowDiagram = %v, want nil", out.ComponentFlowDiagram)
		}
	})
}

func TestFocalPointMetaBody(t *testing.T) {
	t.Run("decodes JSON string fields into slices", func(t *testing.T) {
		body := map[string]interface{}{
			"componentImages":      `["a.png","b.png"]`,
			"componentModalFields": `{"foo":"bar"}`,
			"componentId":          "c1",
		}
		out := FocalPointMetaBody(body)

		images, ok := out["componentImages"].([]interface{})
		if !ok || len(images) != 2 {
			t.Errorf("componentImages = %#v, want a 2-element slice decoded from JSON", out["componentImages"])
		}
		if out["componentId"] != "c1" {
			t.Errorf("componentId = %v, want unchanged passthrough", out["componentId"])
		}
	})

	t.Run("non-string values are left unchanged", func(t *testing.T) {
		body := map[string]interface{}{
			"componentImages": 42,
			"otherKey":        "value",
		}
		out := FocalPointMetaBody(body)
		if out["componentImages"] != 42 {
			t.Errorf("componentImages = %v, want 42 (unchanged non-string)", out["componentImages"])
		}
		if out["otherKey"] != "value" {
			t.Errorf("otherKey = %v, want value", out["otherKey"])
		}
	})
}

func TestFrameToModel(t *testing.T) {
	t.Run("maps all fields with optional pointers set", func(t *testing.T) {
		parentID := "parent-1"
		assetID := "asset-1"
		source := "sync"
		updatedBy := "user-2"
		out := FrameToModel(&uigraphapi.Frame{
			ID: "f1", MapID: "m1", OrgID: "o1",
			ParentFrameID: &parentID, ScreenshotAssetID: &assetID,
			Source: &source, UpdatedBy: &updatedBy,
			Name: "Home", Status: "active",
		})
		if out.ID != "f1" {
			t.Errorf("ID = %q, want f1", out.ID)
		}
		if out.ParentFrameID == nil || *out.ParentFrameID != parentID {
			t.Errorf("ParentFrameID = %v, want pointer to %q", out.ParentFrameID, parentID)
		}
		if out.ScreenshotAssetID == nil || *out.ScreenshotAssetID != assetID {
			t.Errorf("ScreenshotAssetID = %v, want pointer to %q", out.ScreenshotAssetID, assetID)
		}
		if out.Source == nil || *out.Source != source {
			t.Errorf("Source = %v, want pointer to %q", out.Source, source)
		}
		if out.UpdatedBy == nil || *out.UpdatedBy != updatedBy {
			t.Errorf("UpdatedBy = %v, want pointer to %q", out.UpdatedBy, updatedBy)
		}
	})

	t.Run("nil optional pointer fields when source struct has none", func(t *testing.T) {
		out := FrameToModel(&uigraphapi.Frame{ID: "f2", MapID: "m1", OrgID: "o1", Name: "Empty"})
		if out.ParentFrameID != nil {
			t.Errorf("ParentFrameID = %v, want nil", out.ParentFrameID)
		}
		if out.ScreenshotAssetID != nil {
			t.Errorf("ScreenshotAssetID = %v, want nil", out.ScreenshotAssetID)
		}
		if out.Source != nil {
			t.Errorf("Source = %v, want nil", out.Source)
		}
		if out.UpdatedBy != nil {
			t.Errorf("UpdatedBy = %v, want nil", out.UpdatedBy)
		}
	})
}
