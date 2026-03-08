package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
)

const (
	TypePath  = "type"
	NamePath  = "name"
	ValuePath = "value"
)

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := getMetricInfoFromRq(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.metric.UpdateMetric(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, TypePath)
	if mType != models.Gauge && mType != models.Counter {
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	name := chi.URLParam(r, NamePath)

	resp, err := h.metric.GetMetric(&models.Metrics{
		ID:    name,
		MType: mType,
	})

	if err != nil {
		if errors.Is(err, metrics.ErrMetricNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.MType == models.Gauge {
		_, err = w.Write([]byte(fmt.Sprintf("%f", *resp.Value)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if resp.MType == models.Counter {
		_, err = w.Write([]byte(fmt.Sprintf("%d", *resp.Delta)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) GetMetricList(w http.ResponseWriter, r *http.Request) {
	resp, err := h.metric.GetMetricList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = h.metricsTmpl.Execute(w, resp)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
}

func getMetricInfoFromRq(r *http.Request) (*models.Metrics, error) {
	metric := &models.Metrics{}

	metric.ID = chi.URLParam(r, NamePath)
	if len(metric.ID) == 0 {
		return nil, errors.New("no metric name provided")
	}

	metric.MType = chi.URLParam(r, TypePath)
	metricValue := chi.URLParam(r, ValuePath)

	switch metric.MType {
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing metric gauge value: %w", err)
		}
		metric.Value = &value
		return metric, nil

	case models.Counter:
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing metric counter value: %w", err)
		}
		metric.Delta = &delta
		return metric, nil
	default:
		return nil, fmt.Errorf("invalid metric type: %s", metric.MType)
	}
}

var templFunc = template.FuncMap{
	"val": func(p *int64) int64 {
		if p != nil {
			return *p
		}
		return 0
	},
	"fval": func(p *float64) float64 {
		if p != nil {
			return *p
		}
		return 0.0
	},
}
