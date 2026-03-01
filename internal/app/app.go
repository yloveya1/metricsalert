package app

import (
	"context"
	"net/http"

	"github.com/yloveya1/metricsalert/internal/agent/runtimemetrics"
	"github.com/yloveya1/metricsalert/internal/client/httpclient"
	"github.com/yloveya1/metricsalert/internal/handler"
	"github.com/yloveya1/metricsalert/internal/repository/memory"
	"github.com/yloveya1/metricsalert/internal/router"
	"github.com/yloveya1/metricsalert/internal/service/agent"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
)

func RunServer() error {
	storage := memory.NewMemStorage()
	service := metrics.NewService(storage)

	h := handler.New(service)
	r := router.New(h)

	return http.ListenAndServe(":8080", r)
}

func RunAgent(ctx context.Context) error {
	cl := httpclient.NewClient(httpclient.Config{Host: "http://localhost:8080"})
	rc := runtimemetrics.NewRuntimeCollector()
	go rc.CollectMetrics(ctx)

	ag := agent.NewAgent(cl, rc)
	return ag.StartAgent(ctx)
}
