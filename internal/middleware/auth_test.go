package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth_PropagatesHeadersToOutgoingRequest(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	inbound.Header.Set("Authorization", "Bearer abc123")
	inbound.Header.Set("Cookie", "session=xyz")

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	outbound, _ := http.NewRequest(http.MethodGet, "http://upstream/api/v1/auth/me", nil)
	ApplyAuth(capturedCtx, outbound)

	if got := outbound.Header.Get("Authorization"); got != "Bearer abc123" {
		t.Fatalf("Authorization = %q, want %q", got, "Bearer abc123")
	}
	if got := outbound.Header.Get("Cookie"); got != "session=xyz" {
		t.Fatalf("Cookie = %q, want %q", got, "session=xyz")
	}
}

func TestAuth_NoHeadersMeansNothingPropagated(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	outbound, _ := http.NewRequest(http.MethodGet, "http://upstream/api/v1/auth/me", nil)
	ApplyAuth(capturedCtx, outbound)

	if got := outbound.Header.Get("Authorization"); got != "" {
		t.Fatalf("Authorization = %q, want empty", got)
	}
}

func TestBearerToken(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	inbound.Header.Set("Authorization", "Bearer abc123")

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	if got := BearerToken(capturedCtx); got != "abc123" {
		t.Fatalf("BearerToken() = %q, want %q", got, "abc123")
	}
}

func TestBearerToken_NonBearerAuthReturnsEmpty(t *testing.T) {
	inbound := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	inbound.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

	var capturedCtx context.Context
	handler := Auth(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
	}))
	handler.ServeHTTP(httptest.NewRecorder(), inbound)

	if got := BearerToken(capturedCtx); got != "" {
		t.Fatalf("BearerToken() = %q, want empty for non-Bearer Authorization header", got)
	}
}
