package repository

import (
	"context"

	models "github.com/yloveya1/metricsalert/internal/model"
)

//go:generate mockgen -source=interface.go -destination=../mocks/storage.go -package=mocks
type IStorage interface {
	UpdateCounterMetric(ctx context.Context, metrics *models.Metrics) error
	UpdateGaugeMetric(ctx context.Context, metrics *models.Metrics) error
	GetMetricList(ctx context.Context) ([]*models.Metrics, error)
	GetMetricByID(ctx context.Context, metrics *models.Metrics) (models.Metrics, error)
	Ping(ctx context.Context) error
}

type IFile interface {
	WriteMetrics(metrics []*models.Metrics) error
	UploadMetrics() ([]*models.Metrics, error)
}
