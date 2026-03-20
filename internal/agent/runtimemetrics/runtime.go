package runtimemetrics

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	Alloc         = "Alloc"
	BuckHashSys   = "BuckHashSys"
	Frees         = "Frees"
	GCCPUFraction = "GCCPUFraction"
	GCSys         = "GCSys"
	HeapAlloc     = "HeapAlloc"
	HeapIdle      = "HeapIdle"
	HeapInuse     = "HeapInuse"
	HeapObjects   = "HeapObjects"
	HeapReleased  = "HeapReleased"
	HeapSys       = "HeapSys"
	LastGC        = "LastGC"
	Lookups       = "Lookups"
	MCacheInuse   = "MCacheInuse"
	MCacheSys     = "MCacheSys"
	MSpanInuse    = "MSpanInuse"
	MSpanSys      = "MSpanSys"
	Mallocs       = "Mallocs"
	NextGC        = "NextGC"
	NumForcedGC   = "NumForcedGC"
	NumGC         = "NumGC"
	OtherSys      = "OtherSys"
	PauseTotalNs  = "PauseTotalNs"
	StackInuse    = "StackInuse"
	StackSys      = "StackSys"
	Sys           = "Sys"
	TotalAlloc    = "TotalAlloc"
	RandomValue   = "RandomValue"

	PollCount = "PollCount"
)

type RuntimeCollector struct {
	counter      map[string]int64
	gauge        map[string]float64
	mu           sync.RWMutex
	pollInterval time.Duration
}

func NewRuntimeCollector(pollInterval time.Duration) *RuntimeCollector {
	return &RuntimeCollector{
		counter:      make(map[string]int64),
		gauge:        make(map[string]float64),
		pollInterval: pollInterval,
	}
}

func (rc *RuntimeCollector) GetCounterMetrics() map[string]int64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	res := make(map[string]int64, len(rc.counter))
	for k, v := range rc.counter {
		res[k] = v
	}
	return res
}

func (rc *RuntimeCollector) GetGaugeMetrics() map[string]float64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	res := make(map[string]float64, len(rc.gauge))
	for k, v := range rc.gauge {
		res[k] = v
	}
	return res
}

func (rc *RuntimeCollector) CollectMetrics(ctx context.Context) {
	ticker := time.NewTicker(rc.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			rc.mu.Lock()
			rc.setGaugeMetrics(memStats)
			rc.setCounterMetrics()
			rc.mu.Unlock()
		}
	}
}

func (rc *RuntimeCollector) setGaugeMetrics(stats runtime.MemStats) {
	rc.gauge[Alloc] = float64(stats.Alloc)
	rc.gauge[BuckHashSys] = float64(stats.BuckHashSys)
	rc.gauge[Frees] = float64(stats.Frees)
	rc.gauge[GCCPUFraction] = stats.GCCPUFraction
	rc.gauge[GCSys] = float64(stats.GCSys)
	rc.gauge[HeapAlloc] = float64(stats.HeapAlloc)
	rc.gauge[HeapIdle] = float64(stats.HeapIdle)
	rc.gauge[HeapInuse] = float64(stats.HeapInuse)
	rc.gauge[HeapObjects] = float64(stats.HeapObjects)
	rc.gauge[HeapReleased] = float64(stats.HeapReleased)
	rc.gauge[HeapSys] = float64(stats.HeapSys)
	rc.gauge[LastGC] = float64(stats.LastGC)
	rc.gauge[Lookups] = float64(stats.Lookups)
	rc.gauge[MCacheInuse] = float64(stats.MCacheInuse)
	rc.gauge[MCacheSys] = float64(stats.MCacheSys)
	rc.gauge[Mallocs] = float64(stats.Mallocs)
	rc.gauge[NextGC] = float64(stats.NextGC)
	rc.gauge[NumForcedGC] = float64(stats.NumForcedGC)
	rc.gauge[NumGC] = float64(stats.NumGC)
	rc.gauge[OtherSys] = float64(stats.OtherSys)
	rc.gauge[PauseTotalNs] = float64(stats.PauseTotalNs)
	rc.gauge[StackInuse] = float64(stats.StackInuse)
	rc.gauge[StackSys] = float64(stats.StackSys)
	rc.gauge[Sys] = float64(stats.Sys)
	rc.gauge[TotalAlloc] = float64(stats.TotalAlloc)
	rc.gauge[MSpanSys] = float64(stats.MSpanSys)
	rc.gauge[MSpanInuse] = float64(stats.MSpanInuse)
	rc.gauge[RandomValue] = rand.Float64()
}

func (rc *RuntimeCollector) setCounterMetrics() {
	rc.counter[PollCount] = rc.counter[PollCount] + 1
}
