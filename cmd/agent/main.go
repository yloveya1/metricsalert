package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/yloveya1/metricsalert/internal/app"
	"github.com/yloveya1/metricsalert/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := logger.Initialize(zap.InfoLevel.String()); err != nil {
		log.Fatalf("failed to initialize logger, err: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			stop()
			return nil
		case <-egCtx.Done():
			return nil
		}
	})

	eg.Go(func() error { return app.RunAgent(egCtx) })

	if err := eg.Wait(); err != nil {
		logger.AgentLog.Error("failed to run agent", zap.Error(err))
		return
	}
}
