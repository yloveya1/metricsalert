package handler

import "github.com/yloveya1/metricsalert/internal/service/controller"

type Handler struct {
	metric controller.IMetricController
}

func New(metric controller.IMetricController) *Handler {
	return &Handler{
		metric: metric,
	}
}
