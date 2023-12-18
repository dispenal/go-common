package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func TraceHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		operation := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		option := []otelhttp.Option{
			otelhttp.WithPropagators(otel.GetTextMapPropagator()),
		}

		otelhttp.NewHandler(next, operation, option...).ServeHTTP(rw, r)
	})
}
