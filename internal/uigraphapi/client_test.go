package uigraphapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientGet_DecodesResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/orgs/org-1" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Org{ID: "org-1", Name: "Acme"})
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.GetOrg(context.Background(), "org-1")
	if err != nil {
		t.Fatalf("GetOrg() error = %v", err)
	}
	if got.ID != "org-1" || got.Name != "Acme" {
		t.Fatalf("GetOrg() = %+v, want ID=org-1 Name=Acme", got)
	}
}

func TestClientGet_404ReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.GetOrg(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected an error for a 404 response, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.Status != http.StatusNotFound {
		t.Fatalf("APIError.Status = %d, want 404", apiErr.Status)
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(err) = false, want true for err = %v", err)
	}
}

func TestClientGet_500IsNotIsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.GetOrg(context.Background(), "x")
	if err == nil {
		t.Fatal("expected an error for a 500 response, got nil")
	}
	if IsNotFound(err) {
		t.Fatal("IsNotFound(err) = true, want false for a 500 response")
	}
}

func TestClientPost_SendsBody(t *testing.T) {
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/orgs" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Org{ID: "org-2", Name: gotBody["name"].(string)})
	}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.CreateOrg(context.Background(), map[string]interface{}{"name": "Globex"})
	if err != nil {
		t.Fatalf("CreateOrg() error = %v", err)
	}
	if gotBody["name"] != "Globex" {
		t.Fatalf("server received body %v, want name=Globex", gotBody)
	}
	if got.Name != "Globex" {
		t.Fatalf("CreateOrg() = %+v, want Name=Globex", got)
	}
}

func TestAPIError_Error(t *testing.T) {
	e := &APIError{Status: 422, Body: `{"error":"unprocessable"}`}
	msg := e.Error()
	if !strings.Contains(msg, "422") {
		t.Errorf("APIError.Error() = %q, want it to contain the status code 422", msg)
	}
	if !strings.Contains(msg, "unprocessable") {
		t.Errorf("APIError.Error() = %q, want it to contain the body content", msg)
	}
}

func TestClientPing_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(srv.URL)
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping() error = %v, want nil", err)
	}
}

func TestClientPing_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(srv.URL)
	if err := c.Ping(context.Background()); err == nil {
		t.Fatal("Ping() error = nil, want non-nil for 500 response")
	}
}
