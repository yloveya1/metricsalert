package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/yloveya1/metricsalert/internal/agent"
	"github.com/yloveya1/metricsalert/internal/client"
	models "github.com/yloveya1/metricsalert/internal/model"
)

type Agent struct {
	cl             client.IClient
	runtimeAgent   agent.IRuntimeAgent
	reportInterval time.Duration
}

func NewAgent(cl client.IClient, ra agent.IRuntimeAgent, r time.Duration) *Agent {
	return &Agent{
		cl:             cl,
		runtimeAgent:   ra,
		reportInterval: r,
	}
}

func (a *Agent) StartAgent(ctx context.Context) error {
	ticker := time.NewTicker(a.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := a.sendMetric(); err != nil {
				return err
			}
		}
	}
}

func (a *Agent) sendMetric() error {
	for id, value := range a.runtimeAgent.GetGaugeMetrics() {
		metric := &models.Metrics{
			ID:    id,
			MType: models.Gauge,
			Value: &value,
		}
		err := a.cl.SendMetric(metric)
		if err != nil {
			return fmt.Errorf("failed to send gauge metric: %w", err)
		}
	}

	for id, value := range a.runtimeAgent.GetCounterMetrics() {
		metric := &models.Metrics{
			ID:    id,
			MType: models.Counter,
			Delta: &value,
		}
		err := a.cl.SendMetric(metric)
		if err != nil {
			return fmt.Errorf("failed to send counter metric: %w", err)
		}
	}

	return nil
}
