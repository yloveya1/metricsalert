package agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/yloveya1/metricsalert/internal/agent"
	"github.com/yloveya1/metricsalert/internal/client"
	models "github.com/yloveya1/metricsalert/internal/model"
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
				if errors.Is(err, io.EOF) {
					continue
				}
				return err
			}
		}
	}
}

func (a *Agent) sendMetric(metrics []*models.Metrics) error {
	for _, value := range metrics {
		if value == nil {
			continue
		}
		err := a.cl.SendMetric(value)
		if err != nil {
			return fmt.Errorf("failed to send metric %v: %w", *value, err)
		}
	}

	return nil
}
