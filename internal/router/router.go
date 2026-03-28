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

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetricFromBody)
		r.Post("/{type}/{name}/{value}", h.UpdateMetricFromPath)
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricFromBody)
		r.Get("/{type}/{name}", h.GetMetricFromPath)
	})

	r.Get("/", h.GetMetricList)

	return r
}
