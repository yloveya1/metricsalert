package app

import (
	"context"
	"net/http"
	"time"

	"github.com/yloveya1/metricsalert/internal/agent/runtimemetrics"
	"github.com/yloveya1/metricsalert/internal/client/httpclient"
	"github.com/yloveya1/metricsalert/internal/config"
	"github.com/yloveya1/metricsalert/internal/handler"
	"github.com/yloveya1/metricsalert/internal/logger"
	"github.com/yloveya1/metricsalert/internal/repository/memory"
	"github.com/yloveya1/metricsalert/internal/router"
	"github.com/yloveya1/metricsalert/internal/service/agent"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
	"go.uber.org/zap"
)

func RunServer() error {
	if err := logger.Initialize(zap.InfoLevel.String()); err != nil {
		return err
	}

	cfg, err := config.GetServerConfig()
	if err != nil {
		return err
	}

	storage := memory.NewMemStorage()
	service := metrics.NewService(storage)

	h := handler.New(service)
	r := router.New(h)

	return http.ListenAndServe(cfg.Address, r)
}

func RunAgent(ctx context.Context) error {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return err
	}

	cl := httpclient.NewClient(httpclient.Config{Host: "http://" + cfg.Address})
	rc := runtimemetrics.NewRuntimeCollector(time.Duration(cfg.PollInterval) * time.Second)
	go rc.CollectMetrics(ctx)

	ag := agent.NewAgent(cl, rc, time.Duration(cfg.ReportInterval)*time.Second)
	return ag.StartAgent(ctx)
}
