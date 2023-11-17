package jaeger

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestInjectTextMapCarrier(t *testing.T) {
	type args struct {
		spanCtx opentracing.SpanContext
	}

	tests := []struct {
		name    string
		args    args
		want    opentracing.TextMapCarrier
		wantErr bool
	}{
		{
			name: "Success injecting text map carrier",
			args: args{
				spanCtx: nil,
			},
			want:    opentracing.TextMapCarrier{},
			wantErr: false,
		},
		{
			name: "Failed injecting text map carrier",
			args: args{
				spanCtx: nil,
			},
			want:    opentracing.TextMapCarrier{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := InjectTextMapCarrier(tt.args.spanCtx)
			assert.Equal(t, got, tt.want)
		})
	}

}
