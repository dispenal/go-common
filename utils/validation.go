package common_utils

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Message string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error: message %s", ve.Message)
}

type ValidationErrors struct {
	Errors     []ValidationError
	StatusCode int
}

func (ve *ValidationErrors) Error() string {
	return fmt.Sprintf("validation errors: status code %d, message %s", ve.StatusCode, ve.Errors[0].Message)
}

func PanicValidationError(errors []ValidationError, statusCode int) {
	validationErrors := ValidationErrors{
		Errors:     errors,
		StatusCode: statusCode,
	}
	panic(validationErrors)
}

func ValidateStruct(data interface{}) {
	var validationErrors []ValidationError
	validate := validator.New()
	errorValidate := validate.Struct(data)

	if errorValidate != nil {
		for _, err := range errorValidate.(validator.ValidationErrors) {
			var validationError ValidationError
			validationError.Message = strings.Split(err.Error(), "Error:")[1]
			validationErrors = append(validationErrors, validationError)
		}
		PanicValidationError(validationErrors, 400)
	}
}

func ValidateBodyPayload(body io.ReadCloser, output interface{}) {
	err := NewDecoder(body).Decode(output)
	PanicIfAppError(err, "failed when decode body payload", 400)

	ValidateStruct(output)
}
