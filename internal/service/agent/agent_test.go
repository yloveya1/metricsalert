package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yloveya1/metricsalert/internal/mocks"
	models "github.com/yloveya1/metricsalert/internal/model"
)

func Test_NewAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockIClient(ctrl)
	mockCollector := mocks.NewMockIRuntimeAgent(ctrl)

	ag := NewAgent(mockCl, mockCollector, 1, 1)

	assert.Equal(t, mockCl, ag.cl)
	assert.Equal(t, mockCollector, ag.runtimeAgent)
}

func Test_StartAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockIClient(ctrl)
	runtimeAgent := mocks.NewMockIRuntimeAgent(ctrl)

	type args struct {
		ctx     context.Context
		cancel  context.CancelFunc
		metrics []*models.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args args)
	}{
		{
			name: "success send metric",
			args: args{
				metrics: []*models.Metrics{{
					ID:    "test",
					MType: "gauge",
					Value: nil,
				}},
			},
			wantErr: false,
			prepare: func(args args) {
				runtimeAgent.EXPECT().GetMetrics().Return(args.metrics)
				for _, value := range args.metrics {
					mockCl.EXPECT().SendMetric(value).Return(nil)
				}

				args.cancel()
			},
		},
		{
			name: "error",
			args: args{
				metrics: []*models.Metrics{{
					ID:    "test",
					MType: "gauge",
					Value: nil,
				}},
			},
			wantErr: true,
			prepare: func(args args) {
				runtimeAgent.EXPECT().GetMetrics().Return(args.metrics).AnyTimes()
				for _, value := range args.metrics {
					mockCl.EXPECT().SendMetric(value).Return(errors.New("gauge error"))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := Agent{cl: mockCl, runtimeAgent: runtimeAgent, reportInterval: 500 * time.Millisecond, pollInterval: 100 * time.Millisecond}
			tt.args.ctx, tt.args.cancel = context.WithCancel(context.Background())
			if tt.prepare != nil {
				tt.prepare(tt.args)
			}

			err := ag.StartAgent(tt.args.ctx)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
