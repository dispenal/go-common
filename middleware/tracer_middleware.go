package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TraceHttp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		operation := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		otelhttp.NewHandler(next, operation).ServeHTTP(rw, r)
	})
}
