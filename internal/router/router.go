package router

import (
	"net/http"

	"github.com/yloveya1/metricsalert/internal/handler"
)

func New(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", h.UpdateMetric)
	return mux
}
