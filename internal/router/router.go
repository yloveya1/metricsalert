package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yloveya1/metricsalert/internal/handler"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(h.WithLogging())
	r.Use(middleware.Recoverer)

	r.Get("/", h.GetMetricList)

	r.Post("/update/", h.UpdateMetricFromBody)
	r.Post("/update/{type}/{name}/{value}", h.UpdateMetricFromPath)

	r.Post("/value/", h.GetMetricFromBody)
	r.Get("/value/{type}/{name}", h.GetMetricFromPath)

	return r
}
