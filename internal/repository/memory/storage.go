package memory

import (
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
)

type MemStorage struct {
	counter map[string]models.Metrics
	gauge   map[string]models.Metrics
}

func NewMemStorage() repository.IStorage {
	return &MemStorage{
		counter: make(map[string]models.Metrics),
		gauge:   make(map[string]models.Metrics),
	}
}

func (ms *MemStorage) GetMetricList() ([]models.Metrics, error) {
	metricList := make([]models.Metrics, 0, len(ms.counter)+len(ms.gauge))
	for _, v := range ms.counter {
		metricList = append(metricList, v)
	}

	for _, v := range ms.gauge {
		metricList = append(metricList, v)
	}

	return metricList, nil
}

func (ms *MemStorage) GetMetricByID(metric *models.Metrics) (models.Metrics, error) {
	switch metric.MType {
	case models.Counter:
		if m, ok := ms.counter[metric.ID]; ok {
			return m, nil
		}
	case models.Gauge:
		if m, ok := ms.gauge[metric.ID]; ok {
			return m, nil
		}
	}
	return models.Metrics{}, metrics.ErrMetricNotFound
}

func (ms *MemStorage) UpdateGaugeMetric(metric *models.Metrics) error {
	ms.gauge[metric.ID] = *metric
	return nil
}

func (ms *MemStorage) UpdateCounterMetric(metric *models.Metrics) error {
	if _, ok := ms.counter[metric.ID]; !ok {
		ms.counter[metric.ID] = *metric
		return nil
	}

	*ms.counter[metric.ID].Delta += *metric.Delta
	return nil
}
