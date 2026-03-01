package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yloveya1/metricsalert/internal/mocks"
	models "github.com/yloveya1/metricsalert/internal/model"
)

func Test_NewAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCl := mocks.NewMockIClient(ctrl)
	mockCollector := mocks.NewMockIRuntimeAgent(ctrl)

	ag := NewAgent(mockCl, mockCollector)

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
		gauge   map[string]float64
		counter map[string]int64
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
				gauge:   map[string]float64{"test": 1.1},
				counter: map[string]int64{"test": 1},
			},
			wantErr: false,
			prepare: func(args args) {
				runtimeAgent.EXPECT().GetGaugeMetrics().Return(args.gauge)
				for id, value := range args.gauge {
					mockCl.EXPECT().SendMetric(&models.Metrics{
						ID:    id,
						MType: models.Gauge,
						Value: &value,
					}).Return(nil)
				}

				runtimeAgent.EXPECT().GetCounterMetrics().Return(args.counter)
				for id, value := range args.counter {
					mockCl.EXPECT().SendMetric(&models.Metrics{
						ID:    id,
						MType: models.Counter,
						Delta: &value,
					}).Return(nil)
				}

				args.cancel()
			},
		},
		{
			name: "gauge error",
			args: args{
				gauge:   map[string]float64{"test": 1.1},
				counter: map[string]int64{"test": 1},
			},
			wantErr: true,
			prepare: func(args args) {
				runtimeAgent.EXPECT().GetGaugeMetrics().Return(args.gauge)
				for id, value := range args.gauge {
					mockCl.EXPECT().SendMetric(&models.Metrics{
						ID:    id,
						MType: models.Gauge,
						Value: &value,
					}).Return(errors.New("gauge error"))
				}
			},
		},
		{
			name: "counter error",
			args: args{
				gauge:   map[string]float64{"test": 1.1},
				counter: map[string]int64{"test": 1},
			},
			wantErr: true,
			prepare: func(args args) {
				runtimeAgent.EXPECT().GetGaugeMetrics().Return(args.gauge)
				for id, value := range args.gauge {
					mockCl.EXPECT().SendMetric(&models.Metrics{
						ID:    id,
						MType: models.Gauge,
						Value: &value,
					}).Return(nil)
				}

				runtimeAgent.EXPECT().GetCounterMetrics().Return(args.counter)
				for id, value := range args.counter {
					mockCl.EXPECT().SendMetric(&models.Metrics{
						ID:    id,
						MType: models.Counter,
						Delta: &value,
					}).Return(errors.New("counter error"))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := Agent{cl: mockCl, runtimeAgent: runtimeAgent}
			tt.args.ctx, tt.args.cancel = context.WithCancel(context.Background())
			if tt.prepare != nil {
				tt.prepare(tt.args)
			}

			err := ag.StartAgent(tt.args.ctx)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
