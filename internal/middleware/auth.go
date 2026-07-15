package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	authHeaderKey contextKey = "auth_header"
	cookieKey     contextKey = "cookie_header"
	apiKeyKey     contextKey = "api_key_header"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if v := r.Header.Get("Authorization"); v != "" {
			ctx = context.WithValue(ctx, authHeaderKey, v)
		}
		if v := r.Header.Get("Cookie"); v != "" {
			ctx = context.WithValue(ctx, cookieKey, v)
		}
		if v := r.Header.Get("X-API-Key"); v != "" {
			ctx = context.WithValue(ctx, apiKeyKey, v)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ApplyAuth(ctx context.Context, req *http.Request) {
	if v, ok := ctx.Value(authHeaderKey).(string); ok && v != "" {
		req.Header.Set("Authorization", v)
	}
	if v, ok := ctx.Value(cookieKey).(string); ok && v != "" {
		req.Header.Set("Cookie", v)
	}
	if v, ok := ctx.Value(apiKeyKey).(string); ok && v != "" {
		req.Header.Set("X-API-Key", v)
	}
}

func BearerToken(ctx context.Context) string {
	v, _ := ctx.Value(authHeaderKey).(string)
	if strings.HasPrefix(v, "Bearer ") {
		return strings.TrimPrefix(v, "Bearer ")
	}
	return ""
}
