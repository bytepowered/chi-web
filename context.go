package chiweb

import (
	"context"
	"net/http"
)

func RegisterContextAttrs(attrs []ContextAttr) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		wrap := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for _, attr := range attrs {
				ctx = context.WithValue(ctx, attr.Key, attr.Value)
			}
			if ctx != nil && ctx != r.Context() {
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(wrap)
	}
}

type ContextAttr struct {
	Key   any
	Value any
}
