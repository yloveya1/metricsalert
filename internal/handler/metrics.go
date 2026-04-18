package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yloveya1/metricsalert/internal/logger"
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
	"go.uber.org/zap"
)

const (
	TypePath  = "type"
	NamePath  = "name"
	ValuePath = "value"
)

func (h *Handler) UpdateMetricFromPath(w http.ResponseWriter, r *http.Request) {
	metric, err := getMetricInfoFromRq(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.metricCtrl.UpdateMetric(r.Context(), metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateMetricFromBody(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		http.Error(w, "unsupported content type", http.StatusBadRequest)
		return
	}

	var metric models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.metricCtrl.UpdateMetric(r.Context(), &metric); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(metric)
}

func (h *Handler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		http.Error(w, "unsupported content type", http.StatusBadRequest)
		return
	}

	var bodyList []models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&bodyList); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var metricList []*models.Metrics
	for _, m := range bodyList {
		metricList = append(metricList, &m)
	}

	if err := h.metricCtrl.UpdateMetricList(r.Context(), metricList); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(metricList)
}

func (h *Handler) GetMetricFromBody(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		http.Error(w, "unsupported content type", http.StatusBadRequest)
		return
	}

	var metric models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.metricCtrl.GetMetric(r.Context(), &metric)
	if err != nil {
		if errors.Is(err, metrics.ErrMetricNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		logger.ServerLog.Error("encode metric error", zap.Error(err))
	}
}

func (h *Handler) GetMetricFromPath(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, TypePath)
	if mType != models.Gauge && mType != models.Counter {
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	name := chi.URLParam(r, NamePath)

	resp, err := h.metricCtrl.GetMetric(r.Context(), &models.Metrics{
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
		_, err = w.Write([]byte(fmt.Sprintf("%g", *resp.Value)))
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
	resp, err := h.metricCtrl.GetMetricList(r.Context())
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

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.metricCtrl.Ping(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to ping metric, err: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
