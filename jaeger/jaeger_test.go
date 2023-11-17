package jaeger

import (
	"testing"

	common_utils "github.com/dispenal/go-common/utils"
)

func TestJaegerConnection(t *testing.T) {
	t.Run("Success connecting to jaeger", func(t *testing.T) {
		_, _, err := NewJaegerTracer(&common_utils.BaseConfig{
			ServiceName: "test",
			JaegerHost:  "127.0.0.1",
			JaegerPort:  "6831",
		})

		if err != nil {
			t.Error("Failed connecting to jaeger")
		}
	})

	t.Run("Failed connecting to jaeger", func(t *testing.T) {
		_, _, err := NewJaegerTracer(&common_utils.BaseConfig{
			ServiceName: "test",
			JaegerHost:  "127.0.0.2",
			JaegerPort:  "6832",
		})

		if err == nil {
			t.Log("Success connecting to jaeger")
		}
	})

}
