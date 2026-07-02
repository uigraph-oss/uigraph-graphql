package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestMeToModel(t *testing.T) {
	t.Run("maps all fields with avatar", func(t *testing.T) {
		withAvatar := MeToModel(&uigraphapi.MeResponse{
			UserID: "u1", OrgID: "o1", Email: "a@b.com", Name: "Ann", Login: "ann",
			Kind: "user", Role: "server_admin", AuthProvider: "local", AvatarURL: "https://x/a.png",
		})
		if withAvatar.AvatarURL == nil || *withAvatar.AvatarURL != "https://x/a.png" {
			t.Fatalf("AvatarURL = %v, want pointer to https://x/a.png", withAvatar.AvatarURL)
		}
		if withAvatar.UserID != "u1" {
			t.Errorf("UserID = %q, want u1", withAvatar.UserID)
		}
		if withAvatar.OrgID != "o1" {
			t.Errorf("OrgID = %q, want o1", withAvatar.OrgID)
		}
		if withAvatar.Email != "a@b.com" {
			t.Errorf("Email = %q, want a@b.com", withAvatar.Email)
		}
		if withAvatar.Name != "Ann" {
			t.Errorf("Name = %q, want Ann", withAvatar.Name)
		}
		if withAvatar.Login != "ann" {
			t.Errorf("Login = %q, want ann", withAvatar.Login)
		}
		if withAvatar.Kind != "user" {
			t.Errorf("Kind = %q, want user", withAvatar.Kind)
		}
		if !withAvatar.IsServerAdmin {
			t.Errorf("IsServerAdmin = %v, want true", withAvatar.IsServerAdmin)
		}
		if withAvatar.AuthProvider != "local" {
			t.Errorf("AuthProvider = %q, want local", withAvatar.AuthProvider)
		}
	})

	t.Run("nil AvatarURL when empty string", func(t *testing.T) {
		withoutAvatar := MeToModel(&uigraphapi.MeResponse{UserID: "u2"})
		if withoutAvatar.AvatarURL != nil {
			t.Fatalf("AvatarURL = %v, want nil for empty AvatarURL", *withoutAvatar.AvatarURL)
		}
	})
}

func TestOrgSummariesToModel(t *testing.T) {
	t.Run("maps slice of two summaries", func(t *testing.T) {
		in := []uigraphapi.OrgSummary{
			{ID: "1", Name: "A", Role: "admin", Active: true},
			{ID: "2", Name: "B", Role: "member", Active: false},
		}
		out := OrgSummariesToModel(in)
		if len(out) != 2 {
			t.Fatalf("len = %d, want 2", len(out))
		}
		if out[0].ID != "1" || out[0].Name != "A" || out[0].Role != "admin" || !out[0].Active {
			t.Errorf("out[0] = %+v, unexpected", out[0])
		}
		if out[1].ID != "2" || out[1].Active != false {
			t.Errorf("out[1] = %+v, unexpected", out[1])
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		out := OrgSummariesToModel([]uigraphapi.OrgSummary{})
		if len(out) != 0 {
			t.Fatalf("len = %d, want 0", len(out))
		}
	})
}
