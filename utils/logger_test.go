package common_utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	LogInfo("info")
	LogDebug("debug")
	LogError("error")
	assert.PanicsWithValue(t, "panic", func() {
		LogPanic("panic", zap.Error(errors.New("panic-error")))
	})
}
