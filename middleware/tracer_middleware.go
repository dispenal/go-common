package middleware

import (
	"net/http"

	"github.com/dispenal/go-common/tracer"
)

func TraceHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanCtx, span := tracer.StartAndTraceHttp(r, "middleware.TraceHttp")

		defer span.End()

		ctx := r.WithContext(spanCtx)
		next.ServeHTTP(w, ctx)
	})
}
