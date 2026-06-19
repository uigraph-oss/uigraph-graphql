package graph

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

// ErrorPresenter logs every resolver error server-side, then returns a
// sanitized message to the GraphQL client so upstream REST error bodies
// (which may contain internal details) are never forwarded verbatim.
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	gqlErr := graphql.DefaultErrorPresenter(ctx, err)

	var apiErr *uigraphapi.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Status {
		case http.StatusNotFound:
			gqlErr.Message = "not found"
			return gqlErr
		case http.StatusUnauthorized:
			gqlErr.Message = "unauthorized"
			return gqlErr
		case http.StatusForbidden:
			gqlErr.Message = "forbidden"
			return gqlErr
		case http.StatusBadRequest, http.StatusConflict, http.StatusUnprocessableEntity:
			gqlErr.Message = "invalid request"
			return gqlErr
		}
	}

	slog.ErrorContext(ctx, "graphql resolver error", "err", err, "path", graphql.GetPath(ctx).String())
	gqlErr.Message = "internal error"
	return gqlErr
}
