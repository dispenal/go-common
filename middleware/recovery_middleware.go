package middleware

import (
	"fmt"
	"net/http"
	"strings"

	common_utils "github.com/dispenal/go-common/utils"
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
