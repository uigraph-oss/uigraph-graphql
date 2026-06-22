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
