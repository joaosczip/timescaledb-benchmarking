package handlers

import (
	"runtime"
	"time"

	"github.com/joaosczip/timescale/internal/database/repositories"
	"github.com/joaosczip/timescale/internal/dtos"
)

type Clock interface {
	Since(time.Time) time.Duration
}

type CollectCpuUsageMetricsHandler struct {
	repository repositories.CpuUsage
	clock      Clock
}

func NewCollectCpuUsageMetricsHandler(repository repositories.CpuUsage, clock Clock) *CollectCpuUsageMetricsHandler {
	return &CollectCpuUsageMetricsHandler{repository, clock}
}

func (h *CollectCpuUsageMetricsHandler) Handle(queryParams []dtos.CpuUsageQueryParams) (dtos.DatabaseQueriesMetrics, error) {
	workers := runtime.NumCPU()
	metricsCh := make(chan dtos.DatabaseQueriesMetrics, workers)
	errCh := make(chan error, workers)

	for _, queryParams := range queryParams {
		go h.queryStatistics(metricsCh, errCh, queryParams)
	}

	var metrics dtos.DatabaseQueriesMetrics
	for i := 0; i < workers; i++ {
		select {
		case m := <-metricsCh:
			metrics.TotalQueries += m.TotalQueries
			metrics.TotalTime += m.TotalTime

			if m.MinQueryTime < metrics.MinQueryTime || metrics.MinQueryTime == 0 {
				metrics.MinQueryTime = m.MinQueryTime
			}

			if m.MaxQueryTime > metrics.MaxQueryTime {
				metrics.MaxQueryTime = m.MaxQueryTime
			}

			metrics.AvgQueryTime = metrics.TotalTime / time.Duration(metrics.TotalQueries)
			metrics.MedianQueryTime = metrics.TotalTime / 2

		case err := <-errCh:
			return dtos.DatabaseQueriesMetrics{}, err
		}
	}

	return metrics, nil
}

func (h *CollectCpuUsageMetricsHandler) queryStatistics(metricsCh chan<- dtos.DatabaseQueriesMetrics, errCh chan<- error, queryParams dtos.CpuUsageQueryParams) {
	start := time.Now()
	_, err := h.repository.QueryStatistics(queryParams.Host, queryParams.StartTime, queryParams.EndTime)
	if err != nil {
		errCh <- err
		return
	}
	queryDuration := h.clock.Since(start)

	metrics := dtos.DatabaseQueriesMetrics{
		TotalQueries: 1,
		TotalTime:    queryDuration,
	}

	metricsCh <- metrics
}
