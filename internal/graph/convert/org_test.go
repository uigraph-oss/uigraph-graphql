package convert

import (
	"testing"
	"time"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestTeamToModel(t *testing.T) {
	t.Run("maps all fields including optional email and externalID", func(t *testing.T) {
		now := time.Now()
		withExtras := TeamToModel(&uigraphapi.Team{
			ID: "t1", OrgID: "o1", Name: "Platform", Email: "platform@x.com", ExternalID: "ext-1",
			CreatedAt: now, UpdatedAt: now,
		})
		if withExtras.ID != "t1" {
			t.Errorf("ID = %q, want t1", withExtras.ID)
		}
		if withExtras.OrgID != "o1" {
			t.Errorf("OrgID = %q, want o1", withExtras.OrgID)
		}
		if withExtras.Name != "Platform" {
			t.Errorf("Name = %q, want Platform", withExtras.Name)
		}
		if withExtras.Email == nil || *withExtras.Email != "platform@x.com" {
			t.Fatalf("Email = %v, want pointer to platform@x.com", withExtras.Email)
		}
		if withExtras.ExternalID == nil || *withExtras.ExternalID != "ext-1" {
			t.Fatalf("ExternalID = %v, want pointer to ext-1", withExtras.ExternalID)
		}
	})

	t.Run("nil Email and ExternalID when empty strings", func(t *testing.T) {
		bare := TeamToModel(&uigraphapi.Team{ID: "t2", OrgID: "o1", Name: "Bare"})
		if bare.Email != nil {
			t.Fatalf("Email = %v, want nil for empty input", bare.Email)
		}
		if bare.ExternalID != nil {
			t.Fatalf("ExternalID = %v, want nil for empty input", bare.ExternalID)
		}
	})
}

func TestInvitationToModel(t *testing.T) {
	t.Run("maps ExpiresAt pointer when set", func(t *testing.T) {
		expires := time.Now().Add(24 * time.Hour)
		withExpiry := InvitationToModel(uigraphapi.Invitation{
			ID: "i1", OrgID: "o1", Email: "u@x.com", Role: "admin",
			Code: "abc", CreatedBy: "u0", ExpiresAt: &expires,
		})
		if withExpiry.ID != "i1" {
			t.Errorf("ID = %q, want i1", withExpiry.ID)
		}
		if withExpiry.ExpiresAt == nil || !withExpiry.ExpiresAt.Equal(expires) {
			t.Fatalf("ExpiresAt = %v, want %v", withExpiry.ExpiresAt, expires)
		}
	})

	t.Run("nil ExpiresAt when input pointer is nil", func(t *testing.T) {
		noExpiry := InvitationToModel(uigraphapi.Invitation{ID: "i2"})
		if noExpiry.ExpiresAt != nil {
			t.Fatalf("ExpiresAt = %v, want nil", noExpiry.ExpiresAt)
		}
	})
}

func TestInvitationsToModel(t *testing.T) {
	t.Run("converts slice preserving ExpiresAt", func(t *testing.T) {
		expires := time.Now().Add(time.Hour)
		in := []uigraphapi.Invitation{
			{ID: "i1", ExpiresAt: &expires},
			{ID: "i2"},
		}
		out := InvitationsToModel(in)
		if len(out) != 2 {
			t.Fatalf("len = %d, want 2", len(out))
		}
		if out[0].ExpiresAt == nil {
			t.Errorf("out[0].ExpiresAt is nil, want non-nil")
		}
		if out[1].ExpiresAt != nil {
			t.Errorf("out[1].ExpiresAt = %v, want nil", out[1].ExpiresAt)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		out := InvitationsToModel([]uigraphapi.Invitation{})
		if len(out) != 0 {
			t.Fatalf("len = %d, want 0", len(out))
		}
	})
}
