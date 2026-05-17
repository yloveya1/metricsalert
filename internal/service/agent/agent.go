package agent

import (
	"context"
	"sync"
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
	rateLimit      int
}

func NewAgent(cl client.IClient, ra agent.IRuntimeAgent,
	r time.Duration, p time.Duration, rateLimit int) *Agent {
	return &Agent{
		cl:             cl,
		runtimeAgent:   ra,
		reportInterval: r,
		pollInterval:   p,
		rateLimit:      rateLimit,
	}
}

func (a *Agent) StartAgent(ctx context.Context) error {
	var wg sync.WaitGroup

	jobs := make(chan []*models.Metrics, a.rateLimit)

	var mu sync.Mutex
	var currentMetrics []*models.Metrics

	for i := 0; i < a.rateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for metricsBatch := range jobs {
				if err := a.cl.SendMetricList(metricsBatch); err != nil {
					logger.AgentLog.Warn("failed to send metric list", zap.Error(err))
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		tickerP := time.NewTicker(a.pollInterval)
		defer tickerP.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerP.C:
				mu.Lock()
				currentMetrics = a.runtimeAgent.GetMetrics()
				mu.Unlock()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(jobs)

		tickerR := time.NewTicker(a.reportInterval)
		defer tickerR.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerR.C:
				mu.Lock()
				m := currentMetrics
				mu.Unlock()

				if len(m) == 0 {
					continue
				}

				select {
				case jobs <- m:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	<-ctx.Done()

	wg.Wait()

	return nil
}
