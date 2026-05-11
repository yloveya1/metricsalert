package client

import models "github.com/yloveya1/metricsalert/internal/model"

//go:generate mockgen -source=interface.go -destination=../mocks/client.go -package=mocks
type IClient interface {
	SendMetric(metric *models.Metrics) error
	SendMetricList(metric []*models.Metrics) error
}
