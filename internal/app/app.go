package app

import (
	"context"
	"errors"
	"fmt"
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
	"golang.org/x/sync/errgroup"
)

func RunServer(ctx context.Context) error {
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

	eg, egCtx := errgroup.WithContext(ctx)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	eg.Go(func() error {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start server: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()
		shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(shCtx)
	})

	return eg.Wait()
}

func RunAgent(ctx context.Context) error {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return err
	}

	cl := httpclient.NewClient(httpclient.Config{Host: "http://" + cfg.Address})
	rc := runtimemetrics.NewRuntimeCollector()

	ag := agent.NewAgent(cl, rc, time.Duration(cfg.ReportInterval)*time.Second, time.Duration(cfg.PollInterval)*time.Second)
	return ag.StartAgent(ctx)
}
