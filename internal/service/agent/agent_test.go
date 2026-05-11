package agent

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	tests := []struct {
		name    string
		prepare func(
			t *testing.T,
			cl *mocks.MockIClient,
			ra *mocks.MockIRuntimeAgent,
			cancel context.CancelFunc,
		)
	}{
		{
			name: "success send metric",
			prepare: func(
				t *testing.T,
				cl *mocks.MockIClient,
				ra *mocks.MockIRuntimeAgent,
				cancel context.CancelFunc,
			) {
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

				var once sync.Once
				cl.EXPECT().
					SendMetricList(gomock.Any()).
					DoAndReturn(func(m []*models.Metrics) error {
						require.NotNil(t, m)
						require.Equal(t, "test", m[0].ID)
						require.Equal(t, models.Gauge, m[0].MType)

						once.Do(func() {
							cancel()
						})

						return nil
					}).
					AnyTimes()
			},
		},
		{
			name: "client error does not stop agent",
			prepare: func(
				t *testing.T,
				cl *mocks.MockIClient,
				ra *mocks.MockIRuntimeAgent,
				cancel context.CancelFunc,
			) {
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

				var once sync.Once
				cl.EXPECT().
					SendMetricList(gomock.Any()).
					DoAndReturn(func(m []*models.Metrics) error {
						require.NotNil(t, m)
						require.Equal(t, "test", m[0].ID)
						require.Equal(t, models.Gauge, m[0].MType)

						once.Do(func() {
							cancel()
						})

						return errors.New("send error")
					}).
					AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCl := mocks.NewMockIClient(ctrl)
			runtimeAgent := mocks.NewMockIRuntimeAgent(ctrl)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.prepare(t, mockCl, runtimeAgent, cancel)

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
				require.NoError(t, err)
			case <-time.After(1 * time.Second):
				t.Fatal("StartAgent did not stop")
			}
		})
	}
}
