package common_utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	ID uuid.UUID
}

func TestConvertInterfaceE(t *testing.T) {
	copy := testData{}

	data := testData{
		ID: uuid.New(),
	}

	t.Run("Success converting", func(t *testing.T) {
		err := ConvertInterfaceE(data, &copy)

		assert.NoError(t, err)
		assert.Equal(t, data.ID, copy.ID)
	})

	t.Run("Failed when converting", func(t *testing.T) {
		err := ConvertInterfaceE(data, &[]testData{})

		assert.Error(t, err)
	})

	t.Run("Failed when unmarshal", func(t *testing.T) {
		err := ConvertInterfaceE(func() {}, &testData{})

		assert.Error(t, err)
	})
}

func TestConvertInterfaceP(t *testing.T) {
	copy := testData{}

	data := testData{
		ID: uuid.New(),
	}

	t.Run("Success converting", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ConvertInterfaceP(data, &copy)
		})
		assert.Equal(t, data.ID, copy.ID)
	})

	t.Run("Failed when converting", func(t *testing.T) {
		assert.PanicsWithValue(t, AppError{
			Message:    "json: cannot unmarshal object into Go value of type []common_utils.testData|failed when unmarshal interface",
			StatusCode: 422,
		}, func() {
			ConvertInterfaceP(data, &[]testData{})
		})
	})

	t.Run("Failed when unmarshal", func(t *testing.T) {
		assert.PanicsWithValue(t, AppError{
			Message:    "json: unsupported type: func()|failed when marshal interface",
			StatusCode: 422,
		}, func() {
			ConvertInterfaceP(func() {}, &testData{})
		})
	})
}
