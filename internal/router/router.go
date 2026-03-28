package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yloveya1/metricsalert/internal/handler"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(h.WithLogging())

	r.Post("/update", h.UpdateMetricFromBody)
	r.Post("/update/", h.UpdateMetricFromBody)
	r.Post("/", h.UpdateMetricFromBody)

	r.Post("/value", h.GetMetricFromBody)
	r.Post("/value/", h.GetMetricFromBody)
	r.Get("/value/{type}/{name}", h.GetMetricFromPath)

	r.Get("/", h.GetMetricList)

	return r
}
