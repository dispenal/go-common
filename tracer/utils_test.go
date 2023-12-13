package tracer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
)

func TestInjectTextMapCarrier(t *testing.T) {
	type args struct {
		spanCtx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    propagation.MapCarrier
		wantErr bool
	}{
		{
			name: "Success injecting text map carrier",
			args: args{
				spanCtx: nil,
			},
			want:    propagation.MapCarrier{},
			wantErr: false,
		},
		{
			name: "Failed injecting text map carrier",
			args: args{
				spanCtx: nil,
			},
			want:    propagation.MapCarrier{},
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
