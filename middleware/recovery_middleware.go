package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dispenal/go-common/tracer"
	common_utils "github.com/dispenal/go-common/utils"
	"go.opentelemetry.io/otel/codes"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				var errorMsgs []map[string]interface{}
				var statusCode int

				appErr, isAppErr := err.(common_utils.AppError)
				validationErr, isValidationErr := err.(common_utils.ValidationErrors)

				if isAppErr {
					messages := strings.Split(appErr.Message, "|")
					common_utils.LogError(fmt.Sprintf("APP ERROR (PANIC) %s", messages[0]))

					errorMsgs = []map[string]interface{}{
						{"message": messages[1]},
					}
					statusCode = appErr.StatusCode
				} else if isValidationErr {
					common_utils.LogError(fmt.Sprintf("VALIDATION ERROR (PANIC) %v", validationErr))

					for _, err := range validationErr.Errors {
						errorMsg := map[string]interface{}{
							"message": err.Message,
						}
						errorMsgs = append(errorMsgs, errorMsg)
					}
					statusCode = validationErr.StatusCode
				} else {
					common_utils.LogError(fmt.Sprintf("UNKNOWN ERROR (PANIC) %v", validationErr))
					errorMsgs = []map[string]interface{}{
						{"message": "internal server error"},
					}
					statusCode = 500
				}

				common_utils.GenerateJsonResponse(w, nil, statusCode, errorMsgs[0]["message"].(string))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RecoveryTracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				_, span := tracer.StartAndTraceHttp(r, "panic.recovery")
				defer span.End()

				var errorMsgs []map[string]interface{}
				var appError error
				var statusCode int

				appErr, isAppErr := err.(common_utils.AppError)
				validationErr, isValidationErr := err.(common_utils.ValidationErrors)

				if isAppErr {
					messages := strings.Split(appErr.Message, "|")
					errorMsg := fmt.Sprintf("APP ERROR (PANIC) %s", messages[0])
					common_utils.LogError(errorMsg)

					errorMsgs = []map[string]interface{}{
						{"message": messages[1]},
					}
					statusCode = appErr.StatusCode
					appError = errors.New(errorMsg)
				} else if isValidationErr {
					errorMsg := fmt.Sprintf("VALIDATION ERROR (PANIC) %v", validationErr)
					common_utils.LogError(errorMsg)

					for _, err := range validationErr.Errors {
						errorMsg := map[string]interface{}{
							"message": err.Message,
						}
						errorMsgs = append(errorMsgs, errorMsg)
					}
					statusCode = validationErr.StatusCode
					appError = errors.New(errorMsg)
				} else {
					errorMsg := fmt.Sprintf("UNKNOWN ERROR (PANIC) %v", validationErr)
					common_utils.LogError(errorMsg)
					errorMsgs = []map[string]interface{}{
						{"message": "internal server error"},
					}
					statusCode = 500
					appError = errors.New(errorMsg)
				}

				span.RecordError(appError)
				span.SetStatus(codes.Error, appError.Error())

				common_utils.GenerateJsonResponse(w, nil, statusCode, errorMsgs[0]["message"].(string))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
