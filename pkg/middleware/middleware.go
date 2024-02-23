package middleware

import (
	"context"
	"net/http"
)

type middlewareCtxMapKeyType string

var (
	middlewareCtxMapKey middlewareCtxMapKeyType = "middlewareCtxMapKey"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		// Authenticate user token
		next.ServeHTTP(w, r)
	})
}

func PermissionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		// Check Permissions of user token
		next.ServeHTTP(w, r)
	})
}

func LogRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		// log route
		next.ServeHTTP(w, r)
	})
}

func LogUserSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		// log route
		next.ServeHTTP(w, r)
	})
}

// Add generic map to ctx for persisting values within handlers specific to each call.
func AddCtxMapMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m := map[string]interface{}{}
		ctx = context.WithValue(ctx, middlewareCtxMapKey, m)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
