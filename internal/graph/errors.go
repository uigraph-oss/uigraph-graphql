package graph

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

// ErrorPresenter logs every resolver error server-side, then returns a
// sanitized message to the GraphQL client so upstream REST error bodies
// (which may contain internal details) are never forwarded verbatim.
// gqlgen validation/parsing errors are passed through unchanged.
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	// Log the real error server-side before any sanitization.
	slog.ErrorContext(ctx, "resolver error", "err", err)

	// Upstream REST API errors: sanitize based on HTTP status.
	var apiErr *uigraphapi.APIError
	if errors.As(err, &apiErr) {
		return &gqlerror.Error{Message: sanitize(apiErr.Status)}
	}

	// gqlgen validation/parsing errors: pass through unchanged.
	var gqlErr *gqlerror.Error
	if errors.As(err, &gqlErr) {
		return gqlErr
	}

	// All other unexpected errors: generic message.
	return &gqlerror.Error{Message: "internal server error"}
}

func sanitize(status int) string {
	switch status {
	case http.StatusNotFound:
		return "not found"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusBadRequest, http.StatusConflict, http.StatusUnprocessableEntity:
		return "invalid request"
	default:
		return "internal server error"
	}
}
