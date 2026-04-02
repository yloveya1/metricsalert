package metrics

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yloveya1/metricsalert/internal/config"
	"github.com/yloveya1/metricsalert/internal/mocks"
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository/filestore"
)

func Test_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockIStorage(ctrl)
	filestorage := filestore.NewFileStorage("test")
	s := NewService(context.Background(), store, filestorage, config.ServerCfg{})

	assert.NotEmpty(t, s)
}

func TestService_UpdateMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		metric  *models.Metrics
		prepare func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics)
		wantErr bool
	}{
		{
			name: "success updating counter metric",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList().Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(nil)
				store.EXPECT().UpdateCounterMetric(metric).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success updating gauge metric",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList().Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(nil)
				store.EXPECT().UpdateGaugeMetric(metric).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed updating counter metric",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList().Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(nil)
				store.EXPECT().UpdateCounterMetric(metric).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "failed updating gauge metric",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList().Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(nil)
				store.EXPECT().UpdateGaugeMetric(metric).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "failed getting metric list",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList().Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "failed writing metrics",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList().Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "unknown metric type",
			metric: &models.Metrics{
				MType: "invalid type",
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList().Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockIStorage(ctrl)
			filestorage := mocks.NewMockIFile(ctrl)
			s := &Service{
				fileStorage: filestorage,
				storage:     store,
			}

			if tt.prepare != nil {
				tt.prepare(store, filestorage, tt.metric)
			}

			err := s.UpdateMetric(tt.metric)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
