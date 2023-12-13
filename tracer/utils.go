package tracer

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
