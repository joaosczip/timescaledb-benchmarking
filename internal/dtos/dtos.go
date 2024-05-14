package dtos

import "time"

type CpuUsageStatistics struct {
	StartTime string
	EndTime   string
	Max       float64
	Min       float64
}

type CpuUsageQueryParams struct {
	Host      string
	StartTime string
	EndTime   string
}

type DatabaseQueriesMetrics struct {
	TotalQueries    int
	TotalTime       time.Duration
	MinQueryTime    time.Duration
	MaxQueryTime    time.Duration
	AvgQueryTime    time.Duration
	MedianQueryTime time.Duration
}
