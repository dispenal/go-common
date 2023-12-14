package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TraceHttp(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http-request", otelhttp.WithPublicEndpoint())
}
