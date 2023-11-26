package common_utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

type ParamType interface {
	uuid.UUID | string
}

func ValidateUrlParamUUID(r *http.Request, paramName string) uuid.UUID {
	param := chi.URLParam(r, paramName)

	uuid, err := uuid.Parse(param)
	if err != nil {
		PanicIfError(CustomErrorWithTrace(err, fmt.Sprintf("invalid %s param", paramName), 400))
	}

	return uuid
}

func ValidateQueryParamInt(r *http.Request, queryName string) int {
	query := r.URL.Query().Get(queryName)

	queryInt, err := strconv.Atoi(query)
	if err != nil {
		PanicIfError(CustomErrorWithTrace(err, fmt.Sprintf("invalid %s query", queryName), 400))
	}

	if queryInt < 0 {
		PanicIfError(CustomError(fmt.Sprintf("invalid %s query", queryName), 400))
	}

	return queryInt
}

func EnableCORS(r *chi.Mux) {
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "wwww.*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}
