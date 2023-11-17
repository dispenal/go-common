package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBaseConfig(t *testing.T) {
	assert.NotPanics(t, func() {
		config, err := LoadBaseConfig("../", "test")

		assert.NotNil(t, config)
		assert.Nil(t, err)
	})
}

func TestCheckAndSetConfig(t *testing.T) {
	t.Run("Load local config", func(t *testing.T) {
		assert.NotPanics(t, func() {
			config := CheckAndSetConfig("../", "local")
			assert.NotNil(t, config)
			assert.Equal(t, "local", config.ServiceEnv)
		})
	})

	t.Run("Load test config", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "test")
		assert.NotPanics(t, func() {
			config := CheckAndSetConfig("../", "test")
			assert.NotNil(t, config)
			assert.Equal(t, "test", config.ServiceEnv)
		})
	})
}
