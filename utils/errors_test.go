package common_utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errApp = errors.New("panic if app error")
)

func TestPanicIfError(t *testing.T) {
	assert.PanicsWithValue(t, AppError{
		Message:    "panic if error",
		StatusCode: 500,
	}, func() {
		err := errors.New("panic if error")
		PanicIfError(err)
	})
}

func TestPanicIfAppError(t *testing.T) {
	assert.PanicsWithValue(t, AppError{
		Message:    fmt.Sprintf("%s|testing", errApp),
		StatusCode: 422,
	}, func() {
		err := errApp
		PanicIfAppError(err, "testing", 422)
	})
}

func TestPanicAppError(t *testing.T) {
	assert.PanicsWithValue(t, AppError{
		Message:    fmt.Sprintf("|%s", errApp.Error()),
		StatusCode: 422,
	}, func() {
		PanicAppError(errApp.Error(), 422)
	})
}

func TestDeferCheck(t *testing.T) {
	DeferCheck(func() error {
		return errors.New("error")
	})
}

func TestLogIfError(t *testing.T) {
	LogIfError(errors.New("error"))
}

func TestPrintError(t *testing.T) {
	t.Run("App error", func(t *testing.T) {
		err := AppError{
			Message:    "from app error",
			StatusCode: 422,
		}

		assert.Equal(t, "app error: status code 422, message from app error", err.Error())
	})

	t.Run("Validation error", func(t *testing.T) {
		err := ValidationError{
			Message: "from validation error",
		}

		assert.Equal(t, "validation error: message from validation error", err.Error())
	})

	t.Run("Validations error", func(t *testing.T) {
		err := ValidationErrors{
			Errors: []ValidationError{
				{
					Message: "from validation error",
				},
			},
			StatusCode: 400,
		}

		assert.Equal(t, "validation errors: status code 400, message from validation error", err.Error())
	})
}

func TestLogAndPanicIfError(t *testing.T) {
	err := errors.New("error")
	assert.PanicsWithValue(t, err, func() {
		LogAndPanicIfError(err, "error occured")
	})
}
