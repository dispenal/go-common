package middleware

import (
	"net/http"

	"github.com/dispenal/go-common/tracer"
	"go.opentelemetry.io/otel/trace"
)

func TraceHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())

		span.SetAttributes(tracer.BuildAttribute(r.Header)...)

		defer span.End()

		next.ServeHTTP(w, r)
	})
}
