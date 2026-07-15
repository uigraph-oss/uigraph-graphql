package graph

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/uigraph/graphql/internal/uigraphapi"
)

func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	slog.ErrorContext(ctx, "resolver error", "err", err)

	var apiErr *uigraphapi.APIError
	if errors.As(err, &apiErr) {
		if apiErr.Status == http.StatusBadRequest ||
			apiErr.Status == http.StatusConflict ||
			apiErr.Status == http.StatusUnprocessableEntity {
			var parsed struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}
			if json.Unmarshal([]byte(apiErr.Body), &parsed) == nil && parsed.Message != "" {
				return &gqlerror.Error{
					Message:    parsed.Message,
					Extensions: map[string]interface{}{"code": parsed.Code},
				}
			}
		}
		return &gqlerror.Error{Message: sanitize(apiErr.Status)}
	}

	var gqlErr *gqlerror.Error
	if errors.As(err, &gqlErr) {
		return gqlErr
	}

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
