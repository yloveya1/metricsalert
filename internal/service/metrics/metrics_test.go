package metrics

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yloveya1/metricsalert/internal/mocks"
	models "github.com/yloveya1/metricsalert/internal/model"
)

func Test_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockIStorage(ctrl)

	s := NewService(store)

	assert.NotEmpty(t, s)
}

func TestService_UpdateMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		metric  *models.Metrics
		prepare func(store *mocks.MockIStorage, metric *models.Metrics)
		wantErr bool
	}{
		{
			name: "success updating counter metric",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIStorage, metric *models.Metrics) {
				store.EXPECT().UpdateCounterMetric(metric).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success updating gauge metric",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, metric *models.Metrics) {
				store.EXPECT().UpdateGaugeMetric(metric).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed updating counter metric",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIStorage, metric *models.Metrics) {
				store.EXPECT().UpdateCounterMetric(metric).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "failed updating gauge metric",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, metric *models.Metrics) {
				store.EXPECT().UpdateGaugeMetric(metric).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "unknown metric type",
			metric: &models.Metrics{
				MType: "invalid type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockIStorage(ctrl)
			s := &Service{
				storage: store,
			}

			if tt.prepare != nil {
				tt.prepare(store, tt.metric)
			}

			err := s.UpdateMetric(tt.metric)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
