package models

import "time"

type DatabaseMetrics struct {
	TotalQueries    int
	Failures        int
	TotalTime       time.Duration
	MinQueryTime    time.Duration
	MaxQueryTime    time.Duration
	AvgQueryTime    time.Duration
	MedianQueryTime time.Duration
}

func NewDatabaseMetrics() *DatabaseMetrics {
	return &DatabaseMetrics{}
}

func (m *DatabaseMetrics) SetMinQueryTime(newQueryTime time.Duration) {
	if newQueryTime < m.MinQueryTime || m.MinQueryTime == 0 {
		m.MinQueryTime = newQueryTime
	}
}

func (m *DatabaseMetrics) SetMaxQueryTime(newQueryTime time.Duration) {
	if newQueryTime > m.MaxQueryTime {
		m.MaxQueryTime = newQueryTime
	}
}

func (m *DatabaseMetrics) SetAverageQueryTime() {
	m.AvgQueryTime = m.TotalTime / time.Duration(m.TotalQueries)
}

func (m *DatabaseMetrics) SetMedianQueryTime() {
	m.MedianQueryTime = m.TotalTime / 2
}

func (m *DatabaseMetrics) IncrementTotalQueries() {
	m.TotalQueries++
}

func (m *DatabaseMetrics) IncrementTotalTime(newQueryTime time.Duration) {
	m.TotalTime += newQueryTime
}

func (m *DatabaseMetrics) IncrementFailures() {
	m.Failures++
}
