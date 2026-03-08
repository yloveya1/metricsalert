package metrics

import (
	"errors"
	"fmt"
	"sort"

	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
	"github.com/yloveya1/metricsalert/internal/service/controller"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
)

type Service struct {
	storage repository.IStorage
}

func (s *Service) GetMetricList() ([]models.Metrics, error) {
	metricList, err := s.storage.GetMetricList()
	if err != nil {
		return nil, err
	}

	sort.Slice(metricList, func(i, j int) bool {
		return metricList[i].ID < metricList[j].ID
	})

	return metricList, nil
}

func (s *Service) GetMetric(metric *models.Metrics) (models.Metrics, error) {
	return s.storage.GetMetricByID(metric)
}

func NewService(storage repository.IStorage) controller.IMetricController {
	return &Service{
		storage: storage,
	}
}

func (s *Service) UpdateMetric(metric *models.Metrics) error {
	switch metric.MType {
	case models.Counter:
		return s.storage.UpdateCounterMetric(metric)
	case models.Gauge:
		return s.storage.UpdateGaugeMetric(metric)
	default:
		return fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}
