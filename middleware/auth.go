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
)

// Auth extracts the Authorization header and session cookie from the incoming
// request and stores them in the context so client calls can forward them.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if v := r.Header.Get("Authorization"); v != "" {
			ctx = context.WithValue(ctx, authHeaderKey, v)
		}
		if v := r.Header.Get("Cookie"); v != "" {
			ctx = context.WithValue(ctx, cookieKey, v)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ApplyAuth copies stored auth headers onto an outgoing request.
func ApplyAuth(ctx context.Context, req *http.Request) {
	if v, ok := ctx.Value(authHeaderKey).(string); ok && v != "" {
		req.Header.Set("Authorization", v)
	}
	if v, ok := ctx.Value(cookieKey).(string); ok && v != "" {
		req.Header.Set("Cookie", v)
	}
}

// BearerToken extracts the raw token from "Bearer <token>" in the context.
// Returns empty string if absent.
func BearerToken(ctx context.Context) string {
	v, _ := ctx.Value(authHeaderKey).(string)
	if strings.HasPrefix(v, "Bearer ") {
		return strings.TrimPrefix(v, "Bearer ")
	}
	return ""
}
