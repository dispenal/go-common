package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dispenal/go-common/tracer"
	common_utils "github.com/dispenal/go-common/utils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
				_, span := tracer.StartSpan(r.Context())
				defer span.End(trace.WithStackTrace(true))

				var errorMsgs []map[string]interface{}
				var appError error
				var statusCode int

				appErr, isAppErr := err.(common_utils.AppError)
				validationErr, isValidationErr := err.(common_utils.ValidationErrors)

				if isAppErr {
					messages := strings.Split(appErr.Message, "|")
					errorMsg := fmt.Sprintf("APP ERROR (PANIC) %s", messages[0])
					common_utils.LogError(errorMsg)
					appError = errors.New(errorMsg)

					errorMsgs = []map[string]interface{}{
						{"message": messages[1]},
					}
					statusCode = appErr.StatusCode
				} else if isValidationErr {
					errorMsg := fmt.Sprintf("VALIDATION ERROR (PANIC) %v", validationErr)
					common_utils.LogError(errorMsg)
					appError = errors.New(errorMsg)

					for _, err := range validationErr.Errors {
						errorMsg := map[string]interface{}{
							"message": err.Message,
						}
						errorMsgs = append(errorMsgs, errorMsg)
					}
					statusCode = validationErr.StatusCode
				} else {
					errorMsg := fmt.Sprintf("UNKNOWN ERROR (PANIC) %v", validationErr)
					common_utils.LogError(errorMsg)
					appError = errors.New(errorMsg)

					errorMsgs = []map[string]interface{}{
						{"message": "internal server error"},
					}
					statusCode = 500
				}

				span.RecordError(appError)
				span.SetStatus(codes.Error, appError.Error())

				headers := make(map[string]string)
				for k, v := range r.Header {
					headers[k] = v[0]
				}
				attributes := tracer.BuildAttribute(headers)

				span.SetAttributes(attributes...)

				common_utils.GenerateJsonResponse(w, nil, statusCode, errorMsgs[0]["message"].(string))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
