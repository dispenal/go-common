package tracer

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

type IPgxTracer interface {
	TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context

	TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData)
}

type PgxTracer struct {
}

func (t *PgxTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	spanCtx, span := StartAndTraceWithData(ctx, "pgx.Query.Start", map[string]any{
		"sql":  data.SQL,
		"args": data.Args,
	})
	defer span.End()

	return spanCtx
}

func (t *PgxTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	tracer := trace.SpanFromContext(ctx)
	defer tracer.End()

	tracer.SetName("pgx.Query.End")

	tracer.SetAttributes(BuildAttribute(data)...)
}
