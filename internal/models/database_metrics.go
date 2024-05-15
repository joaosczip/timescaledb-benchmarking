package models

import (
	"fmt"
	"sort"
)

type DatabaseMetrics struct {
	TotalQueries    int
	Failures        int
	TotalTime       float64
	MinQueryTime    float64
	MaxQueryTime    float64
	AvgQueryTime    float64
	MedianQueryTime float64
	queries         []float64
}

func NewDatabaseMetrics() *DatabaseMetrics {
	return &DatabaseMetrics{}
}

func (m *DatabaseMetrics) SetAverageQueryTime() {
	if m.TotalQueries > 0 {
		m.AvgQueryTime = m.TotalTime / float64(m.TotalQueries)
	}
}

func (m *DatabaseMetrics) SetMedianQueryTime() {
	if len(m.queries) == 0 {
		return
	}
	
	sort.Float64s(m.queries)
	midIndex := len(m.queries) / 2

	if len(m.queries)%2 == 0 {
		m.MedianQueryTime = (m.queries[midIndex-1] + m.queries[midIndex]) / 2
	} else {
		m.MedianQueryTime = m.queries[midIndex]
	}
}

func (m *DatabaseMetrics) AddQueryTime(newQueryTime float64) {
	m.TotalQueries++

	if newQueryTime > m.MaxQueryTime {
		m.MaxQueryTime = newQueryTime
	} else if newQueryTime < m.MinQueryTime || m.MinQueryTime == 0 {
		m.MinQueryTime = newQueryTime
	}

	m.TotalTime += newQueryTime
	m.queries = append(m.queries, newQueryTime)
}

func (m *DatabaseMetrics) IncrementFailures() {
	m.Failures++
}

func (m DatabaseMetrics) ToMap() map[string]string {
	return map[string]string{
		"total_queries":     fmt.Sprintf("%d", m.TotalQueries),
		"failures":          fmt.Sprintf("%d", m.Failures),
		"min_query_time":    fmt.Sprintf("%.2fms", m.MinQueryTime),
		"max_query_time":    fmt.Sprintf("%.2fms", m.MaxQueryTime),
		"avg_query_time":    fmt.Sprintf("%.2fms", m.AvgQueryTime),
		"median_query_time": fmt.Sprintf("%.2fms", m.MedianQueryTime),
		"total_time":        fmt.Sprintf("%.2fms", m.TotalTime),
	}
}
