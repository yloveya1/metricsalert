package controller

import (
	"context"

	models "github.com/yloveya1/metricsalert/internal/model"
)

type IMetricController interface {
	UpdateMetric(ctx context.Context, metrics *models.Metrics) error
	UpdateMetricList(ctx context.Context, metrics []*models.Metrics) error
	GetMetricList(ctx context.Context) ([]*models.Metrics, error)
	GetMetric(ctx context.Context, metric *models.Metrics) (models.Metrics, error)
	Ping(ctx context.Context) error
}
