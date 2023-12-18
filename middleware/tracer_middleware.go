package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func TraceHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		operation := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		ctx := r.Context()

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

		r = r.WithContext(ctx)
		otelhttp.NewHandler(next, operation).ServeHTTP(rw, r)
	})
}
