package convert

import (
	"testing"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestUserToModel(t *testing.T) {
	t.Run("maps LastSeenAt pointer when set", func(t *testing.T) {
		lastSeen := time.Now()
		active := UserToModel(&uigraphapi.User{
			ID: "u1", Email: "a@b.com", Name: "Alice", Login: "alice",
			Disabled: false, Role: "admin", LastSeenAt: &lastSeen,
		})
		if active.ID != "u1" {
			t.Errorf("ID = %q, want u1", active.ID)
		}
		if active.Email != "a@b.com" {
			t.Errorf("Email = %q, want a@b.com", active.Email)
		}
		if active.Name != "Alice" {
			t.Errorf("Name = %q, want Alice", active.Name)
		}
		if active.Login != "alice" {
			t.Errorf("Login = %q, want alice", active.Login)
		}
		if active.Role != "admin" {
			t.Errorf("Role = %q, want admin", active.Role)
		}
		if active.Disabled != false {
			t.Errorf("Disabled = %v, want false", active.Disabled)
		}
		if active.LastSeenAt == nil || !active.LastSeenAt.Equal(lastSeen) {
			t.Fatalf("LastSeenAt = %v, want %v", active.LastSeenAt, lastSeen)
		}
	})

	t.Run("nil LastSeenAt when input pointer is nil", func(t *testing.T) {
		neverSeen := UserToModel(&uigraphapi.User{ID: "u2"})
		if neverSeen.LastSeenAt != nil {
			t.Fatalf("LastSeenAt = %v, want nil", neverSeen.LastSeenAt)
		}
	})
}

func TestUsersToModel(t *testing.T) {
	t.Run("converts slice preserving LastSeenAt", func(t *testing.T) {
		lastSeen := time.Now()
		in := []uigraphapi.User{
			{ID: "u1", LastSeenAt: &lastSeen},
			{ID: "u2"},
		}
		out := UsersToModel(in)
		if len(out) != 2 {
			t.Fatalf("len = %d, want 2", len(out))
		}
		if out[0].LastSeenAt == nil {
			t.Errorf("out[0].LastSeenAt is nil, want non-nil")
		}
		if out[1].LastSeenAt != nil {
			t.Errorf("out[1].LastSeenAt = %v, want nil", out[1].LastSeenAt)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		out := UsersToModel([]uigraphapi.User{})
		if len(out) != 0 {
			t.Fatalf("len = %d, want 0", len(out))
		}
	})
}
