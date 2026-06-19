// Package client provides a typed HTTP client for the uigraph-api REST API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uigraph/graphql/internal/middleware"
)

// Client calls the uigraph-api REST backend.
type Client struct {
	base string
	http *http.Client
}

// New returns a Client targeting baseURL (e.g. "http://uigraph-api:8080").
func New(baseURL string) *Client {
	return &Client{
		base: baseURL,
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

// ── internal helpers ─────────────────────────────────────────────────────────

func (c *Client) do(ctx context.Context, method, path string, body, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.base+path, bodyReader)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	middleware.ApplyAuth(ctx, req)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		return &APIError{Status: resp.StatusCode, Body: string(data)}
	}

	if out != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode: %w", err)
		}
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, out interface{}) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) post(ctx context.Context, path string, body, out interface{}) error {
	return c.do(ctx, http.MethodPost, path, body, out)
}

func (c *Client) put(ctx context.Context, path string, body, out interface{}) error {
	return c.do(ctx, http.MethodPut, path, body, out)
}

func (c *Client) del(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// APIError carries the HTTP status code returned by uigraph-api.
type APIError struct {
	Status int
	Body   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error %d: %s", e.Status, e.Body)
}

// IsNotFound returns true when the upstream returned 404.
func IsNotFound(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.Status == http.StatusNotFound
	}
	return false
}
