package tracer

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"cloud.google.com/go/pubsub"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func StartAndTrace(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer("")
	spanCtx, span := tracer.Start(ctx, spanName)

	return spanCtx, span
}

func StartAndTraceWithData(ctx context.Context, spanName string, data ...any) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer("")
	spanCtx, span := tracer.Start(ctx, spanName)

	bag := BuildBaggage(data)
	defaultCtx := baggage.ContextWithBaggage(spanCtx, bag)

	return defaultCtx, span
}

func StartAndTraceHttp(r *http.Request, spanName string) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer("")
	spanCtx, span := tracer.Start(r.Context(), spanName)

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(spanCtx, propagation.HeaderCarrier(r.Header))

	return spanCtx, span
}

func StartAndTracePubsub(ctx context.Context, spanName string, data *pubsub.Message) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer("")
	spanCtx, span := tracer.Start(ctx, spanName)

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(spanCtx, propagation.MapCarrier(data.Attributes))

	return spanCtx, span
}

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

func BuildBaggage(args ...any) baggage.Baggage {
	members := make([]baggage.Member, 0)

	for _, arg := range args {
		v := reflect.ValueOf(arg)
		for i := 0; i < v.NumField(); i++ {
			isEmpty := v.Field(i).Interface() == ""
			if isEmpty {
				continue
			}

			if v.Field(i).Kind() == reflect.String {
				member, _ := baggage.NewMember(v.Type().Field(i).Name, v.Field(i).Interface().(string))
				members = append(members, member)
			}

		}
	}

	bag, _ := baggage.New(members...)

	return bag
}
