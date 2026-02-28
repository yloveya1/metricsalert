package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/yloveya1/metricsalert/internal/model"
)

const (
	TypePath  = "type"
	NamePath  = "name"
	ValuePath = "value"
)

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metrics, err := getMetricsInfoFromRq(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	case models.Gauge:
		err = h.metric.UpdateGaugeMetric(metrics)
	case models.Counter:
		err = h.metric.UpdateCounterMetric(metrics)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getMetricsInfoFromRq(r *http.Request) (*models.Metrics, error) {
	metrics := &models.Metrics{}

	metrics.ID = r.PathValue(NamePath)
	if len(metrics.ID) == 0 {
		return nil, errors.New("no metric name provided")
	}

	metrics.MType = r.PathValue(TypePath)
	metricValue := r.PathValue(ValuePath)

	switch metrics.MType {
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing metric gauge value: %w", err)
		}
		metrics.Value = &value
		return metrics, nil

	case models.Counter:
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing metric counter value: %w", err)
		}
		metrics.Delta = &delta
		return metrics, nil
	default:
		return nil, errors.New(fmt.Sprintf("invalid metric type: %s", metrics.MType))
	}
}
