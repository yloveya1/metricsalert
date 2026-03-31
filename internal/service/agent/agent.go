package agent

import (
	"context"
	"time"

	"github.com/yloveya1/metricsalert/internal/agent"
	"github.com/yloveya1/metricsalert/internal/client"
	"github.com/yloveya1/metricsalert/internal/logger"
	models "github.com/yloveya1/metricsalert/internal/model"
	"go.uber.org/zap"
)

type Agent struct {
	cl             client.IClient
	runtimeAgent   agent.IRuntimeAgent
	reportInterval time.Duration
	pollInterval   time.Duration
}

func NewAgent(cl client.IClient, ra agent.IRuntimeAgent,
	r time.Duration, p time.Duration) *Agent {
	return &Agent{
		cl:             cl,
		runtimeAgent:   ra,
		reportInterval: r,
		pollInterval:   p,
	}
}

func (a *Agent) StartAgent(ctx context.Context) error {
	ticker := time.NewTicker(a.reportInterval)
	defer ticker.Stop()

	tickerP := time.NewTicker(a.pollInterval)
	defer tickerP.Stop()

	var metrics []*models.Metrics

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tickerP.C:
			metrics = a.runtimeAgent.GetMetrics()
		case <-ticker.C:
			if err := a.sendMetric(metrics); err != nil {
				return err
			}
		}
	}
}

func (a *Agent) sendMetric(metrics []*models.Metrics) error {
	const (
		maxRetries = 3
		retryDelay = 100 * time.Millisecond
	)

	for _, m := range metrics {
		if m == nil {
			continue
		}

		for attempt := 1; attempt <= maxRetries; attempt++ {
			if err := a.cl.SendMetric(m); err != nil {
				logger.AgentLog.Warn(
					"failed to send metric",
					zap.String("id", m.ID),
					zap.String("type", m.MType),
					zap.Int("attempt", attempt),
					zap.Error(err),
				)
				time.Sleep(retryDelay)
				continue
			}

			break
		}
	}
	return nil
}
