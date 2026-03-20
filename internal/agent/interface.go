package agent

//go:generate mockgen -source=interface.go -destination=../mocks/agent.go -package=mocks
type IRuntimeAgent interface {
	GetCounterMetrics() map[string]int64
	GetGaugeMetrics() map[string]float64
}
