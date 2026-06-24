package convert

import (
	"testing"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestSavingsSummaryToModel(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		if got := SavingsSummaryToModel(nil); got != nil {
			t.Fatalf("SavingsSummaryToModel(nil) = %v, want nil", got)
		}
	})

	t.Run("empty ModelID becomes nil (blended)", func(t *testing.T) {
		got := SavingsSummaryToModel(&uigraphapi.SavingsSummary{OrgID: "o1", Period: "7d", ModelID: "", TotalCalls: 3})
		if got.ModelID != nil {
			t.Errorf("ModelID = %v, want nil for blended summary", got.ModelID)
		}
		if got.TotalCalls != 3 {
			t.Errorf("TotalCalls = %d, want 3", got.TotalCalls)
		}
	})

	t.Run("non-empty ModelID is preserved as pointer", func(t *testing.T) {
		got := SavingsSummaryToModel(&uigraphapi.SavingsSummary{ModelID: "claude-sonnet-4-6"})
		if got.ModelID == nil || *got.ModelID != "claude-sonnet-4-6" {
			t.Errorf("ModelID = %v, want pointer to claude-sonnet-4-6", got.ModelID)
		}
	})
}

func TestDailySavingsListToModel(t *testing.T) {
	t.Run("maps each row in order", func(t *testing.T) {
		now := time.Now()
		got := DailySavingsListToModel([]uigraphapi.DailySavings{
			{Date: now, TotalCalls: 1, TotalTokensSaved: 10},
			{Date: now.AddDate(0, 0, 1), TotalCalls: 2, TotalTokensSaved: 20},
		})
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
		if got[0].TotalCalls != 1 || got[1].TotalCalls != 2 {
			t.Errorf("TotalCalls = [%d, %d], want [1, 2]", got[0].TotalCalls, got[1].TotalCalls)
		}
	})

	t.Run("empty input returns empty slice, not nil", func(t *testing.T) {
		got := DailySavingsListToModel(nil)
		if got == nil {
			t.Fatal("got nil, want empty slice")
		}
		if len(got) != 0 {
			t.Errorf("len = %d, want 0", len(got))
		}
	})
}

func TestUserSavingsToModel(t *testing.T) {
	t.Run("resolves display name from actor map for a user", func(t *testing.T) {
		uid := "u1"
		actors := map[string]*uigraphapi.Actor{"u1": {ID: "u1", Name: "Ada Lovelace"}}
		got := UserSavingsToModel(uigraphapi.UserSavings{UserID: &uid, TotalCalls: 5}, actors)
		if got.DisplayName != "Ada Lovelace" {
			t.Errorf("DisplayName = %q, want Ada Lovelace", got.DisplayName)
		}
	})

	t.Run("falls back to Service Account when actor not found", func(t *testing.T) {
		said := "sa1"
		got := UserSavingsToModel(uigraphapi.UserSavings{ServiceAccountID: &said}, map[string]*uigraphapi.Actor{})
		if got.DisplayName != "Service Account" {
			t.Errorf("DisplayName = %q, want Service Account", got.DisplayName)
		}
	})

	t.Run("falls back to Unknown User when neither id nor actor resolves", func(t *testing.T) {
		got := UserSavingsToModel(uigraphapi.UserSavings{}, map[string]*uigraphapi.Actor{})
		if got.DisplayName != "Unknown User" {
			t.Errorf("DisplayName = %q, want Unknown User", got.DisplayName)
		}
	})
}
