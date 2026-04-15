package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
)

type MemStorage struct {
	counter map[string]models.Metrics
	gauge   map[string]models.Metrics
	mu      sync.RWMutex
}

func NewMemStorage() repository.IStorage {
	return &MemStorage{
		counter: make(map[string]models.Metrics),
		gauge:   make(map[string]models.Metrics),
	}
}

func (ms *MemStorage) GetMetricList(ctx context.Context) ([]*models.Metrics, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metricList := make([]*models.Metrics, 0, len(ms.counter)+len(ms.gauge))
	for _, v := range ms.counter {
		value := v
		metricList = append(metricList, &value)
	}

	for _, v := range ms.gauge {
		value := v
		metricList = append(metricList, &value)
	}

	return metricList, nil
}

func (ms *MemStorage) GetMetricByID(ctx context.Context, metric *models.Metrics) (models.Metrics, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

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

func (ms *MemStorage) UpdateGaugeMetric(ctx context.Context, metric *models.Metrics) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.gauge[metric.ID] = *metric
	return nil
}

func (ms *MemStorage) UpdateCounterMetric(ctx context.Context, metric *models.Metrics) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if metric.Delta == nil {
		return errors.New("metrics counter delta is nil")
	}

	val, ok := ms.counter[metric.ID]
	if !ok {
		ms.counter[metric.ID] = *metric
		return nil
	}

	*val.Delta += *metric.Delta
	ms.counter[metric.ID] = val

	return nil
}
func (ms *MemStorage) UpdateMetricList(ctx context.Context, metrics []*models.Metrics) error {
	for _, m := range metrics {
		switch m.MType {
		case models.Counter:
			err := ms.UpdateCounterMetric(ctx, m)
			if err != nil {
				return fmt.Errorf("failed to update counter metric, err: %w", err)
			}
		case models.Gauge:
			err := ms.UpdateGaugeMetric(ctx, m)
			if err != nil {
				return fmt.Errorf("failed to update gauge metric, err: %w", err)
			}
		default:
			return fmt.Errorf("unknown metric type: %s", m.MType)
		}
	}

	return nil
}

func (ms *MemStorage) Ping(ctx context.Context) error {
	return nil
}
