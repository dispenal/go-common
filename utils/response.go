package common_utils

import (
	"net/http"
)

// swagger:model Response
type Response[T any] struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       T      `json:"data,omitempty"`
}

// FOR TESTING PURPOSE
type ResponseMap struct {
	StatusCode int                    `json:"status_code"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// FOR TESTING PURPOSE
type ResponseSlice struct {
	StatusCode int                      `json:"status_code"`
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Data       []map[string]interface{} `json:"data,omitempty"`
}

func GenerateJsonResponse(w http.ResponseWriter, data interface{}, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response[interface{}]{
		Message:    message,
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Data:       data,
	}

	responseEncode, err := Marshal(response)
	PanicIfAppError(err, "failed when marshar response", 500)
	w.Write(responseEncode)
}
