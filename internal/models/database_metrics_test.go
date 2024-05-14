package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseMetrics(t *testing.T) {
	t.Run("Should set the average query time", func(t *testing.T) {
		metrics := DatabaseMetrics{TotalQueries: 10, TotalTime: 100}

		metrics.SetAverageQueryTime()

		assert.Equal(t, 10.0, metrics.AvgQueryTime)
	})

	t.Run("Should set the median query time", func(t *testing.T) {
		metrics := DatabaseMetrics{TotalTime: 100, TotalQueries: 2}

		metrics.SetMedianQueryTime()

		assert.Equal(t, 50.0, metrics.MedianQueryTime)
	})

	t.Run("Adding a new query time should update the metrics", func(t *testing.T) {
		metrics := DatabaseMetrics{
			TotalQueries: 10,
			TotalTime:    100.0,
			MinQueryTime: 50.0,
			MaxQueryTime: 60.0,
		}

		metrics.AddQueryTime(80.0)

		assert.Equal(t, 11, metrics.TotalQueries)
		assert.Equal(t, 50.0, metrics.MinQueryTime)
		assert.Equal(t, 80.0, metrics.MaxQueryTime)
		assert.Equal(t, 180.0, metrics.TotalTime)

		metrics.AddQueryTime(40.0)

		assert.Equal(t, 12, metrics.TotalQueries)
		assert.Equal(t, 40.0, metrics.MinQueryTime)
		assert.Equal(t, 80.0, metrics.MaxQueryTime)
		assert.Equal(t, 220.0, metrics.TotalTime)
	})

	t.Run("Should increment the failures", func(t *testing.T) {
		metrics := DatabaseMetrics{}

		metrics.IncrementFailures()

		assert.Equal(t, 1, metrics.Failures)

		metrics.IncrementFailures()

		assert.Equal(t, 2, metrics.Failures)
	})

	t.Run("Should only set the avg query time when the total queries is greater than 0", func(t *testing.T) {
		metrics := DatabaseMetrics{}

		metrics.SetAverageQueryTime()

		assert.Equal(t, 0.0, metrics.AvgQueryTime)
	})

	t.Run("Printing all the metrics should format it nicely", func(t *testing.T) {
		metrics := DatabaseMetrics{
			TotalQueries:    10,
			Failures:        2,
			MinQueryTime:    10.0,
			MaxQueryTime:    20.0,
			AvgQueryTime:    15.0,
			MedianQueryTime: 12.0,
			TotalTime:       150.0,
		}

		expected := map[string]string{
			"total_queries":     "10",
			"failures":          "2",
			"min_query_time":    "10.00ms",
			"max_query_time":    "20.00ms",
			"avg_query_time":    "15.00ms",
			"median_query_time": "12.00ms",
			"total_time":        "150.00ms",
		}

		assert.Equal(t, expected, metrics.ToMap())
	})
}
