package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging_SetsRequestIDInContext(t *testing.T) {
	var gotID string
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = RequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if gotID == "" {
		t.Fatal("RequestID(ctx) is empty, want a generated request id")
	}
}

func TestLogging_RecordsResponseStatus(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusTeapot {
		t.Fatalf("recorder status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}

func TestLogging_EchoesRequestIDInResponseHeader(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if got := rec.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("X-Request-ID response header is empty, want a generated request id")
	}
}

func TestLogging_HonoursIncomingRequestID(t *testing.T) {
	const incomingID = "my-trace-id-42"

	var gotID string
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = RequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-ID", incomingID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if gotID != incomingID {
		t.Fatalf("RequestID(ctx) = %q, want %q", gotID, incomingID)
	}
	if got := rec.Header().Get("X-Request-ID"); got != incomingID {
		t.Fatalf("X-Request-ID response header = %q, want %q", got, incomingID)
	}
}

func TestLogging_GeneratesFreshUUIDWhenNoHeader(t *testing.T) {
	var gotID string
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = RequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	// No X-Request-ID header set.
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if gotID == "" {
		t.Fatal("RequestID(ctx) is empty, want a generated UUID")
	}
	if got := rec.Header().Get("X-Request-ID"); got != gotID {
		t.Fatalf("X-Request-ID response header = %q, want %q", got, gotID)
	}
}
