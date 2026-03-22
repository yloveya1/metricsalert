package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/yloveya1/metricsalert/internal/handler"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(h.WithLogging())
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	r.Get("/value/{type}/{name}", h.GetMetric)

	r.Get("/", h.GetMetricList)

	return r
}
