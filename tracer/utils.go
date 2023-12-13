package tracer

import (
	"context"
	"encoding/json"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func InjectTextMapCarrier(spanCtx context.Context) (propagation.MapCarrier, error) {
	m := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(spanCtx, propagation.MapCarrier{})

	return m, nil
}
func ExtractTextMapCarrier(spanCtx context.Context) propagation.MapCarrier {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return make(propagation.MapCarrier)
	}
	return textMapCarrier
}

func ExtractTextMapCarrierBytes(spanCtx context.Context) []byte {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return []byte("")
	}

	dataBytes, err := json.Marshal(&textMapCarrier)
	if err != nil {
		return []byte("")
	}
	return dataBytes
}

func TraceErr(span trace.Span, err error) {
	span.RecordError(err)
	span.SetAttributes(attribute.Bool("error", true))
	span.SetAttributes(attribute.String("error_code", err.Error()))
}

func TraceWithErr(span trace.Span, err error) error {
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
		span.SetAttributes(attribute.String("error_code", err.Error()))
	}

	return err
}

func MetricLatency(ctx context.Context, span trace.Span, meter metric.Meter, attributes []attribute.KeyValue) {
	startTime := time.Now()
	latencyMs := float64(time.Since(startTime)) / 1e6

	requestLatency, _ := meter.Float64Histogram(
		"request_latency",
		metric.WithDescription("The latency of requests processed"),
	)

	requestLatency.Record(ctx, latencyMs, metric.WithAttributes(attributes...))

	span.SetAttributes(attribute.Float64("latency", latencyMs))
}

func MetricCount(ctx context.Context, meter metric.Meter, attributes []attribute.KeyValue) {
	requestCount, _ := meter.Int64Counter(
		"request_counts",
		metric.WithDescription("The number of requests processed"),
	)

	requestCount.Add(ctx, 1, metric.WithAttributes(attributes...))

}

func MetricLineCount(ctx context.Context, meter metric.Meter, attributes []attribute.KeyValue) {
	lineCounts, _ := meter.Int64Counter(
		"line_counts",
		metric.WithDescription("The counts of the lines in"),
	)

	lineCounts.Add(ctx, 1, metric.WithAttributes(attributes...))

}
