package middleware

import (
	"net/http"
	"time"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupMiddleware(route *chi.Mux, config *common_utils.BaseConfig) {
	if config.ServiceEnv == common_utils.TEST || config.ServiceEnv == common_utils.DEVELOPMENT {
		route.Use(TraceHttp)
		route.Use(RecoveryTracer)
	} else {
		route.Use(middleware.RequestID)
		route.Use(middleware.RealIP)
		route.Use(middleware.Logger)
		route.Use(middleware.Timeout(60 * time.Second))

		route.Use(TraceHttp)
		route.Use(RecoveryTracer)
	}

	route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		common_utils.GenerateJsonResponse(w, nil, http.StatusNotFound, "Not Found")
	})
}
