package repository

import models "github.com/yloveya1/metricsalert/internal/model"

type IStorage interface {
	UpdateCounterMetric(metric *models.Metrics) error
	UpdateGaugeMetric(metric *models.Metrics) error
}
