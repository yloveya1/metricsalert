package memory

import (
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() repository.IStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGaugeMetric(metric *models.Metrics) error {
	ms.gauge[metric.ID] = *metric.Value
	return nil
}

func (ms *MemStorage) UpdateCounterMetric(metric *models.Metrics) error {
	ms.counter[metric.ID] += *metric.Delta
	return nil
}
