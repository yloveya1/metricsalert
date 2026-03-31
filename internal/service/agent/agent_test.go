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
	defer ctrl.Finish()

	mockCl := mocks.NewMockIClient(ctrl)
	runtimeAgent := mocks.NewMockIRuntimeAgent(ctrl)

	tests := []struct {
		name    string
		prepare func(
			cl *mocks.MockIClient,
			ra *mocks.MockIRuntimeAgent,
			cancel context.CancelFunc,
		)
		wantErr bool
	}{
		{
			name: "success send metric",
			prepare: func(cl *mocks.MockIClient, ra *mocks.MockIRuntimeAgent, cancel context.CancelFunc) {
				metrics := []*models.Metrics{
					{
						ID:    "test",
						MType: models.Gauge,
					},
				}

				ra.EXPECT().
					GetMetrics().
					Return(metrics).
					AnyTimes()

				called := false
				cl.EXPECT().
					SendMetric(metrics[0]).
					DoAndReturn(func(*models.Metrics) error {
						if !called {
							called = true
							cancel()
						}
						return nil
					}).
					AnyTimes()
			},
			wantErr: false,
		},
		{
			name: "client send error does not stop agent",
			prepare: func(cl *mocks.MockIClient, ra *mocks.MockIRuntimeAgent, cancel context.CancelFunc) {
				metrics := []*models.Metrics{
					{
						ID:    "test",
						MType: models.Gauge,
					},
				}

				ra.EXPECT().
					GetMetrics().
					Return(metrics).
					AnyTimes()

				called := false
				cl.EXPECT().
					SendMetric(metrics[0]).
					DoAndReturn(func(*models.Metrics) error {
						if !called {
							called = true
							cancel()
						}
						return errors.New("send error")
					}).
					AnyTimes()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.prepare != nil {
				tt.prepare(mockCl, runtimeAgent, cancel)
			}

			ag := Agent{
				cl:             mockCl,
				runtimeAgent:   runtimeAgent,
				reportInterval: 20 * time.Millisecond,
				pollInterval:   10 * time.Millisecond,
			}

			errCh := make(chan error, 1)
			go func() {
				errCh <- ag.StartAgent(ctx)
			}()

			select {
			case err := <-errCh:
				assert.Equal(t, tt.wantErr, err != nil)
			case <-time.After(500 * time.Millisecond):
				t.Fatal("StartAgent did not stop")
			}
		})
	}
}
