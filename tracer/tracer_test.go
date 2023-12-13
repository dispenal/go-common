package tracer

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	common_utils "github.com/dispenal/go-common/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
)

func TestJaegerConnection(t *testing.T) {
	t.Run("Span context", func(t *testing.T) {
		shutdown := NewTracer(&common_utils.BaseConfig{
			ServiceName: "test",
			JaegerHost:  "127.0.0.1",
			JaegerPort:  "4317",
		})
		defer shutdown()

		tracer := otel.GetTracerProvider().Tracer("demo-client-tracer")
		meter := otel.Meter("demo-client-meter")

		method, _ := baggage.NewMember("method", "repl")
		client, _ := baggage.NewMember("client", "cli")
		bag, _ := baggage.New(method, client)

		requestLatency, _ := meter.Float64Histogram(
			"demo_client/request_latency",
			metric.WithDescription("The latency of the requests"),
		)

		lineLengths, _ := meter.Int64Histogram(
			"demo_client/line_lengths",
			metric.WithDescription("The lengths of the various lines in"),
		)

		lineCounts, _ := meter.Int64Counter(
			"demo_client/line_counts",
			metric.WithDescription("The counts of the lines in"),
		)

		requestCount, _ := meter.Int64Counter(
			"demo_client/request_counts",
			metric.WithDescription("The number of requests processed"),
		)

		commonLabels := []attribute.KeyValue{
			attribute.String("method", "repl"),
			attribute.String("client", "cli"),
		}

		defaultCtx := baggage.ContextWithBaggage(context.Background(), bag)
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			startTime := time.Now()
			ctx, span := tracer.Start(defaultCtx, "ExecuteRequest")
			time.Sleep(time.Duration(2) * time.Second)
			span.End()
			latencyMs := float64(time.Since(startTime)) / 1e6
			nr := int(rng.Int31n(7))
			for i := 0; i < nr; i++ {
				randLineLength := rng.Int63n(999)
				lineCounts.Add(ctx, 1, metric.WithAttributes(commonLabels...))
				lineLengths.Record(ctx, randLineLength, metric.WithAttributes(commonLabels...))
				fmt.Printf("#%d: LineLength: %dBy\n", i, randLineLength)
			}

			requestLatency.Record(ctx, latencyMs, metric.WithAttributes(commonLabels...))
			requestCount.Add(ctx, 1, metric.WithAttributes(commonLabels...))

			fmt.Printf("Latency: %.3fms\n", latencyMs)
			time.Sleep(time.Duration(1) * time.Second)
		}

	})

}
