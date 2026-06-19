package convert

import (
	"testing"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TestAuthConfigToModel(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		if got := AuthConfigToModel(nil); got != nil {
			t.Fatalf("AuthConfigToModel(nil) = %v, want nil", got)
		}
	})

	t.Run("bearer token type with pointer field", func(t *testing.T) {
		token := "secret"
		got := AuthConfigToModel(&uigraphapi.AuthConfig{Type: "bearer", BearerToken: &token})
		if got == nil {
			t.Fatal("got = nil, want non-nil")
		}
		if got.Type != "bearer" {
			t.Errorf("Type = %q, want bearer", got.Type)
		}
		if got.BearerToken == nil || *got.BearerToken != token {
			t.Errorf("BearerToken = %v, want pointer to %q", got.BearerToken, token)
		}
	})

	t.Run("nil optional pointer fields when source has none", func(t *testing.T) {
		got := AuthConfigToModel(&uigraphapi.AuthConfig{Type: "none"})
		if got.BearerToken != nil {
			t.Errorf("BearerToken = %v, want nil", got.BearerToken)
		}
		if got.APIKeyHeader != nil {
			t.Errorf("APIKeyHeader = %v, want nil", got.APIKeyHeader)
		}
		if got.APIKeyValue != nil {
			t.Errorf("APIKeyValue = %v, want nil", got.APIKeyValue)
		}
		if got.BasicUsername != nil {
			t.Errorf("BasicUsername = %v, want nil", got.BasicUsername)
		}
		if got.BasicPassword != nil {
			t.Errorf("BasicPassword = %v, want nil", got.BasicPassword)
		}
	})
}

func TestManualTestCaseToModel(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		if got := ManualTestCaseToModel(nil); got != nil {
			t.Fatalf("ManualTestCaseToModel(nil) = %v, want nil", got)
		}
	})

	t.Run("steps are converted correctly", func(t *testing.T) {
		preconditions := "user is logged in"
		expectedOutcome := "success"
		in := &uigraphapi.ManualTestCase{
			Preconditions:   &preconditions,
			ExpectedOutcome: &expectedOutcome,
			Steps: []uigraphapi.TestCaseStep{
				{Order: 1, Action: "open page", ExpectedResult: "page loads"},
				{Order: 2, Action: "click button", ExpectedResult: "modal opens"},
			},
		}
		got := ManualTestCaseToModel(in)
		if got == nil {
			t.Fatal("got = nil, want non-nil")
		}
		if len(got.Steps) != 2 {
			t.Fatalf("len(Steps) = %d, want 2", len(got.Steps))
		}
		if got.Steps[0].Action != "open page" {
			t.Errorf("Steps[0].Action = %q, want open page", got.Steps[0].Action)
		}
		if got.Steps[0].ExpectedResult != "page loads" {
			t.Errorf("Steps[0].ExpectedResult = %q, want page loads", got.Steps[0].ExpectedResult)
		}
		if got.Steps[1].Order != 2 {
			t.Errorf("Steps[1].Order = %d, want 2", got.Steps[1].Order)
		}
	})

	t.Run("empty steps slice returns empty slice", func(t *testing.T) {
		got := ManualTestCaseToModel(&uigraphapi.ManualTestCase{})
		if len(got.Steps) != 0 {
			t.Errorf("len(Steps) = %d, want 0", len(got.Steps))
		}
	})
}

func TestTestRunResultToModel(t *testing.T) {
	t.Run("ResponseTimeMs converted from int64 to int when set", func(t *testing.T) {
		var ms int64 = 1500
		got := TestRunResultToModel(&uigraphapi.TestRunResult{TestRunResultID: "r1", ResponseTimeMs: &ms})
		if got.TestRunResultID != "r1" {
			t.Errorf("TestRunResultID = %q, want r1", got.TestRunResultID)
		}
		if got.ResponseTimeMs == nil {
			t.Fatal("ResponseTimeMs = nil, want pointer to 1500")
		}
		if *got.ResponseTimeMs != 1500 {
			t.Errorf("*ResponseTimeMs = %d, want 1500 (converted from *int64 to *int)", *got.ResponseTimeMs)
		}
	})

	t.Run("nil ResponseTimeMs remains nil", func(t *testing.T) {
		got := TestRunResultToModel(&uigraphapi.TestRunResult{TestRunResultID: "r2"})
		if got.ResponseTimeMs != nil {
			t.Errorf("ResponseTimeMs = %v, want nil", got.ResponseTimeMs)
		}
	})
}

func TestTestPackToModel(t *testing.T) {
	t.Run("maps all fields", func(t *testing.T) {
		out := TestPackToModel(&uigraphapi.TestPack{
			TestPackID: "tp1", ServiceID: "svc1", OrgID: "o1",
			Name: "Smoke Tests", Type: "manual",
			CreatedBy: "user-1",
		})
		if out.TestPackID != "tp1" {
			t.Errorf("TestPackID = %q, want tp1", out.TestPackID)
		}
		if out.ServiceID != "svc1" {
			t.Errorf("ServiceID = %q, want svc1", out.ServiceID)
		}
		if out.OrgID != "o1" {
			t.Errorf("OrgID = %q, want o1", out.OrgID)
		}
		if out.Name != "Smoke Tests" {
			t.Errorf("Name = %q, want Smoke Tests", out.Name)
		}
		if out.Type != "manual" {
			t.Errorf("Type = %q, want manual", out.Type)
		}
		if out.CreatedBy != "user-1" {
			t.Errorf("CreatedBy = %q, want user-1", out.CreatedBy)
		}
	})
}

func TestTestCaseToModel(t *testing.T) {
	t.Run("maps all fields including optional nil converters", func(t *testing.T) {
		priority := "high"
		out := TestCaseToModel(&uigraphapi.TestCase{
			TestCaseID: "tc1", TestPackID: "tp1", ServiceID: "svc1", OrgID: "o1",
			Title: "Login test", Order: 1, Type: "manual",
			Priority: &priority, IsCritical: true, EvidenceRequired: false,
		})
		if out.TestCaseID != "tc1" {
			t.Errorf("TestCaseID = %q, want tc1", out.TestCaseID)
		}
		if out.Title != "Login test" {
			t.Errorf("Title = %q, want Login test", out.Title)
		}
		if out.Order != 1 {
			t.Errorf("Order = %v, want 1", out.Order)
		}
		if out.Priority == nil || *out.Priority != "high" {
			t.Errorf("Priority = %v, want pointer to high", out.Priority)
		}
		if !out.IsCritical {
			t.Errorf("IsCritical = false, want true")
		}
		// nil sub-converters should return nil
		if out.Manual != nil {
			t.Errorf("Manual = %v, want nil when source Manual is nil", out.Manual)
		}
		if out.API != nil {
			t.Errorf("API = %v, want nil when source API is nil", out.API)
		}
	})
}

func TestTestRunToModel(t *testing.T) {
	t.Run("maps all fields", func(t *testing.T) {
		releaseLabel := "v1.2"
		startedBy := "user-1"
		out := TestRunToModel(&uigraphapi.TestRun{
			TestRunID: "tr1", TestPackID: "tp1", ServiceID: "svc1", OrgID: "o1",
			Environment: "staging", ReleaseLabel: &releaseLabel, Status: "completed",
			StartedBy: &startedBy, OverallStatus: "passed",
		})
		if out.TestRunID != "tr1" {
			t.Errorf("TestRunID = %q, want tr1", out.TestRunID)
		}
		if out.TestPackID != "tp1" {
			t.Errorf("TestPackID = %q, want tp1", out.TestPackID)
		}
		if out.Environment != "staging" {
			t.Errorf("Environment = %q, want staging", out.Environment)
		}
		if out.ReleaseLabel == nil || *out.ReleaseLabel != "v1.2" {
			t.Errorf("ReleaseLabel = %v, want pointer to v1.2", out.ReleaseLabel)
		}
		if out.Status != "completed" {
			t.Errorf("Status = %q, want completed", out.Status)
		}
		if out.OverallStatus != "passed" {
			t.Errorf("OverallStatus = %q, want passed", out.OverallStatus)
		}
		if out.StartedBy == nil || *out.StartedBy != "user-1" {
			t.Errorf("StartedBy = %v, want pointer to user-1", out.StartedBy)
		}
	})
}
