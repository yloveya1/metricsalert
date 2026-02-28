package metrics

import (
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
	"github.com/yloveya1/metricsalert/internal/service/controller"
)

type Service struct {
	storage repository.IStorage
}

func NewService(storage repository.IStorage) controller.IMetricController {
	return &Service{
		storage: storage,
	}
}
func (s *Service) UpdateCounterMetric(metric *models.Metrics) error {
	return s.storage.UpdateCounterMetric(metric)
}

func (s *Service) UpdateGaugeMetric(metric *models.Metrics) error {
	return s.storage.UpdateGaugeMetric(metric)
}
