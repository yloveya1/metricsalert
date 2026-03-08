package controller

import models "github.com/yloveya1/metricsalert/internal/model"

type IMetricController interface {
	UpdateMetric(metric *models.Metrics) error
	GetMetricList() ([]models.Metrics, error)
	GetMetric(metric *models.Metrics) (models.Metrics, error)
}
