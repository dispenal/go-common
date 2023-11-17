package jaeger

import (
	"fmt"
	"io"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
)

func NewJaegerTracer(baseConfig *common_utils.BaseConfig) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: baseConfig.ServiceName,

		// "const" sampler is a binary sampling strategy: 0=never sample, 1=always sample.
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		// Log the emitted spans to stdout.
		Reporter: &config.ReporterConfig{
			LogSpans:           baseConfig.JaegerLogSpans,
			LocalAgentHostPort: fmt.Sprintf("%s:%s", baseConfig.JaegerHost, baseConfig.JaegerPort),
		},
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	return cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		config.Injector(opentracing.TextMap, zipkinPropagator),
		config.Injector(opentracing.Binary, zipkinPropagator),
		config.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		config.Extractor(opentracing.TextMap, zipkinPropagator),
		config.Extractor(opentracing.Binary, zipkinPropagator),
		config.ZipkinSharedRPCSpan(false),
	)
}
