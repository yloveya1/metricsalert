package app

import (
	"net/http"

	"github.com/yloveya1/metricsalert/internal/handler"
	"github.com/yloveya1/metricsalert/internal/repository/memory"
	"github.com/yloveya1/metricsalert/internal/router"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
)

func Run() error {
	storage := memory.NewMemStorage()
	service := metrics.NewService(storage)

	h := handler.New(service)
	r := router.New(h)

	return http.ListenAndServe(":8080", r)
}
