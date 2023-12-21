package common_utils

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicValidationError(t *testing.T) {
	defer func() {
		err := recover()

		validationErrors, isValidationErrors := err.(ValidationErrors)
		assert.Equal(t, true, isValidationErrors)
		assert.Equal(t, validationErrors.StatusCode, 400)
		assert.Equal(t, "validation errors", validationErrors.Errors[0].Message)
	}()

	validationError := []ValidationError{{Message: "validation errors"}}
	PanicValidationError(validationError, 400)
}

func TestValidateStruct(t *testing.T) {
	defer func() {
		err := recover()

		validationErrors, isValidationErrors := err.(ValidationErrors)
		assert.Equal(t, true, isValidationErrors)
		assert.Equal(t, 400, validationErrors.StatusCode)
		assert.Equal(t, "Field validation for 'Email' failed on the 'email' tag", validationErrors.Errors[0].Message)

	}()
	type testInput struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"min=3,max=10"`
		Password string `json:"password" validate:"min=8,max=15"`
	}

	input := testInput{
		Email:    "test@test",
		Username: "Testing",
		Password: "xxxxxxxx",
	}

	ValidateStruct(&input)
}

func TestValidateBodyPayload(t *testing.T) {
	type testInput struct {
		Success bool `json:"success"`
	}
	input := testInput{
		Success: true,
	}
	body, err := Marshal(input)
	assert.NoError(t, err)
	reader := bytes.NewReader(body)
	var output testInput

	ValidateBodyPayload(io.NopCloser(reader), &output)
	assert.Equal(t, true, output.Success)
}
