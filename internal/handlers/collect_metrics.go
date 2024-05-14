package handlers

import (
	"fmt"
	"runtime"
	"time"

	"github.com/joaosczip/timescale/internal/database/repositories"
	"github.com/joaosczip/timescale/internal/dtos"
	"github.com/joaosczip/timescale/internal/models"
)

type CollectCpuUsageMetricsHandler struct {
	repository repositories.CpuUsage
}

func NewCollectCpuUsageMetricsHandler(repository repositories.CpuUsage) *CollectCpuUsageMetricsHandler {
	return &CollectCpuUsageMetricsHandler{repository}
}

func (h *CollectCpuUsageMetricsHandler) Handle(queryParams []dtos.CpuUsageQueryParams) *models.DatabaseMetrics {
	workers := runtime.NumCPU()
	queryDurationCh := make(chan float64, workers)
	errCh := make(chan error, workers)

	for _, queryParams := range queryParams {
		go h.queryStatistics(queryDurationCh, errCh, queryParams)
	}

	metrics := models.NewDatabaseMetrics()

	for i := 0; i < len(queryParams); i++ {
		select {
		case queryDuration := <-queryDurationCh:
			metrics.AddQueryTime(queryDuration)
		case err := <-errCh:
			fmt.Printf("Error querying statistics: %v\n", err)
			metrics.IncrementFailures()
		}
	}

	metrics.SetAverageQueryTime()
	metrics.SetMedianQueryTime()

	return metrics
}

func (h *CollectCpuUsageMetricsHandler) queryStatistics(queryDurationCh chan<- float64, errCh chan<- error, queryParams dtos.CpuUsageQueryParams) {
	start := time.Now()

	_, err := h.repository.QueryStatistics(queryParams.Host, queryParams.StartTime, queryParams.EndTime)

	if err != nil {
		errCh <- err
	} else {
		duration := time.Since(start)
		queryDurationCh <- duration.Seconds()
	}
}
