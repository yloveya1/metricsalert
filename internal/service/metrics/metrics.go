package metrics

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/yloveya1/metricsalert/internal/config"
	"github.com/yloveya1/metricsalert/internal/logger"
	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
	"github.com/yloveya1/metricsalert/internal/service/controller"
	"go.uber.org/zap"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
)

type Service struct {
	storage     repository.IStorage
	fileStorage repository.IFile
	cfg         config.ServerCfg
}

func (s *Service) GetMetricList(ctx context.Context) ([]*models.Metrics, error) {
	metricList, err := s.storage.GetMetricList(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(metricList, func(i, j int) bool {
		return metricList[i].ID < metricList[j].ID
	})

	return metricList, nil
}

func (s *Service) GetMetric(ctx context.Context, metric *models.Metrics) (models.Metrics, error) {
	return s.storage.GetMetricByID(ctx, metric)
}

func NewService(ctx context.Context, storage repository.IStorage, fileStorage repository.IFile,
	cfg config.ServerCfg) controller.IMetricController {
	srv := &Service{
		storage:     storage,
		fileStorage: fileStorage,
		cfg:         cfg,
	}

	if *cfg.StoreInterval > 0 {
		go srv.runPeriodSafe(ctx)
	}

	if !*cfg.Restore {
		return srv
	}

	metricList, err := srv.fileStorage.UploadMetrics()
	if err != nil {
		logger.ServerLog.Error("failed to upload metrics", zap.Error(err))
		return srv
	}

	for _, metric := range metricList {
		if err = srv.UpdateMetric(ctx, metric); err != nil {
			logger.ServerLog.Error("failed to update metric", zap.Error(err))
		}
	}

	return srv
}

func (s *Service) runPeriodSafe(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(*s.cfg.StoreInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := s.saveMetrics(ctx); err != nil {
				logger.ServerLog.Error("failed to save metrics", zap.Error(err))
			} else {
				logger.ServerLog.Info("saved metrics")
			}
		}
	}
}

func (s *Service) saveMetrics(ctx context.Context) error {
	metrics, err := s.storage.GetMetricList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get metrics, err: %w", err)
	}

	err = s.fileStorage.WriteMetrics(metrics)
	if err != nil {
		return fmt.Errorf("failed to write metrics, err: %w", err)
	}

	return nil
}

func (s *Service) UpdateMetric(ctx context.Context, metric *models.Metrics) error {
	switch metric.MType {
	case models.Counter:
		if err := s.storage.UpdateCounterMetric(ctx, metric); err != nil {
			return err
		}
	case models.Gauge:
		if err := s.storage.UpdateGaugeMetric(ctx, metric); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown metric type: %s", metric.MType)
	}

	if *s.cfg.StoreInterval == 0 {
		if err := s.saveMetrics(ctx); err != nil {
			return fmt.Errorf("failed to save to file: %w", err)
		}
	}

	return nil
}

func (s *Service) Ping(ctx context.Context) error {
	err := s.storage.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping storage, err: %w", err)
	}
	return nil
}
