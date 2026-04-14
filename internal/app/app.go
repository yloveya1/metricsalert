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
	"github.com/yloveya1/metricsalert/internal/repository"
	"github.com/yloveya1/metricsalert/internal/repository/filestore"
	"github.com/yloveya1/metricsalert/internal/repository/memory"
	"github.com/yloveya1/metricsalert/internal/repository/pg"
	"github.com/yloveya1/metricsalert/internal/router"
	"github.com/yloveya1/metricsalert/internal/service/agent"
	"github.com/yloveya1/metricsalert/internal/service/metrics"
	"golang.org/x/sync/errgroup"
)

func RunServer(ctx context.Context) error {
	cfg, err := config.GetServerConfig()
	if err != nil {
		return err
	}

	var storage repository.IStorage
	if *cfg.DBConn != "" {
		storage, err = pg.NewDatabase(ctx, *cfg.DBConn)
		if err != nil {
			return fmt.Errorf("failed to connect database, err: %w", err)
		}
	} else {
		storage = memory.NewMemStorage()
	}

	fileStorage := filestore.NewFileStorage(*cfg.FileStoragePath)

	service := metrics.NewService(ctx, storage, fileStorage, cfg)

	h := handler.New(service)
	r := router.New(h)

	eg, egCtx := errgroup.WithContext(ctx)

	server := &http.Server{
		Addr:    *cfg.Address,
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

	cl := httpclient.NewClient(httpclient.Config{Host: "http://" + *cfg.Address})
	rc := runtimemetrics.NewRuntimeCollector()

	ag := agent.NewAgent(cl, rc, time.Duration(*cfg.ReportInterval)*time.Second, time.Duration(*cfg.PollInterval)*time.Second)
	return ag.StartAgent(ctx)
}
