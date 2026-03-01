package repository

import models "github.com/yloveya1/metricsalert/internal/model"

//go:generate mockgen -source=interface.go -destination=../mocks/storage.go -package=mocks
type IStorage interface {
	UpdateCounterMetric(metric *models.Metrics) error
	UpdateGaugeMetric(metric *models.Metrics) error
}
