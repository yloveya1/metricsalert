package repository

import (
	"context"

	models "github.com/yloveya1/metricsalert/internal/model"
)

//go:generate mockgen -source=interface.go -destination=../mocks/storage.go -package=mocks
type IStorage interface {
	UpdateCounterMetric(metric *models.Metrics) error
	UpdateGaugeMetric(metric *models.Metrics) error
	GetMetricList() ([]*models.Metrics, error)
	GetMetricByID(metric *models.Metrics) (models.Metrics, error)
	Ping(ctx context.Context) error
}

type IFile interface {
	WriteMetrics(metrics []*models.Metrics) error
	UploadMetrics() ([]*models.Metrics, error)
}
