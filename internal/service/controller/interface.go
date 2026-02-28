package controller

import models "github.com/yloveya1/metricsalert/internal/model"

type IMetricController interface {
	UpdateCounterMetric(metric *models.Metrics) error
	UpdateGaugeMetric(metric *models.Metrics) error
}
