package agent

import models "github.com/yloveya1/metricsalert/internal/model"

//go:generate mockgen -source=interface.go -destination=../mocks/agent.go -package=mocks
type IRuntimeAgent interface {
	GetMetrics() []*models.Metrics
}
