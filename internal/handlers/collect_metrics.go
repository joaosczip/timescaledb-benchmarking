package handlers

import (
	"fmt"
	"runtime"
	"time"

	"github.com/joaosczip/timescale/internal/database/repositories"
	"github.com/joaosczip/timescale/internal/dtos"
	"github.com/joaosczip/timescale/internal/models"
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

func (h *CollectCpuUsageMetricsHandler) Handle(queryParams []dtos.CpuUsageQueryParams) *models.DatabaseMetrics {
	workers := runtime.NumCPU()
	queryDurationCh := make(chan time.Duration, workers)
	errCh := make(chan error, workers)

	for _, queryParams := range queryParams {
		go h.queryStatistics(queryDurationCh, errCh, queryParams)
	}

	metrics := models.NewDatabaseMetrics()

	for i := 0; i < workers; i++ {
		select {
		case queryDuration := <-queryDurationCh:
			metrics.IncrementTotalQueries()
			metrics.IncrementTotalTime(queryDuration)
			metrics.SetMaxQueryTime(queryDuration)
			metrics.SetMinQueryTime(queryDuration)
		case err := <-errCh:
			fmt.Printf("Error querying statistics: %v\n", err)
			metrics.IncrementFailures()
		}
	}

	metrics.SetAverageQueryTime()
	metrics.SetMedianQueryTime()

	return metrics
}

func (h *CollectCpuUsageMetricsHandler) queryStatistics(queryDurationCh chan<- time.Duration, errCh chan<- error, queryParams dtos.CpuUsageQueryParams) {
	start := time.Now()

	_, err := h.repository.QueryStatistics(queryParams.Host, queryParams.StartTime, queryParams.EndTime)

	if err != nil {
		errCh <- err
	} else {
		queryDurationCh <- h.clock.Since(start)
	}
}
