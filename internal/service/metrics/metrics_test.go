package metrics

import (
	"context"
	"errors"
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
	s := NewService(context.Background(), store, filestorage, config.ServerCfg{
		Address:         ptr("address"),
		StoreInterval:   ptr(0),
		FileStoragePath: ptr("filapath"),
		Restore:         ptr(false),
	})

	assert.NotEmpty(t, s)
}

func ptr[T any](v T) *T {
	return &v
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
				store.EXPECT().GetMetricList(gomock.Any()).Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(nil)
				store.EXPECT().UpdateCounterMetric(gomock.Any(), metric).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success updating gauge metric",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList(gomock.Any()).Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(nil)
				store.EXPECT().UpdateGaugeMetric(gomock.Any(), metric).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed updating counter metric",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().UpdateCounterMetric(gomock.Any(), metric).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "failed updating gauge metric",
			metric: &models.Metrics{
				MType: models.Gauge,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().UpdateGaugeMetric(gomock.Any(), metric).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "failed getting metric list",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList(gomock.Any()).Return([]*models.Metrics{metric}, errors.New("error"))
				store.EXPECT().UpdateCounterMetric(gomock.Any(), metric).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "failed write metrics",
			metric: &models.Metrics{
				MType: models.Counter,
			},
			prepare: func(store *mocks.MockIStorage, filestorage *mocks.MockIFile, metric *models.Metrics) {
				store.EXPECT().GetMetricList(gomock.Any()).Return([]*models.Metrics{metric}, nil)
				filestorage.EXPECT().WriteMetrics([]*models.Metrics{metric}).Return(errors.New("error"))
				store.EXPECT().UpdateCounterMetric(gomock.Any(), metric).Return(nil)
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
			filestorage := mocks.NewMockIFile(ctrl)
			s := &Service{
				fileStorage: filestorage,
				storage:     store,
				cfg: config.ServerCfg{
					Address:         ptr("address"),
					StoreInterval:   ptr(0),
					FileStoragePath: ptr("filapath"),
					Restore:         ptr(false),
				},
			}

			if tt.prepare != nil {
				tt.prepare(store, filestorage, tt.metric)
			}

			err := s.UpdateMetric(context.Background(), tt.metric)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
