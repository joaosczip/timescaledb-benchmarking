package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseMetrics(t *testing.T) {
	t.Run("Should set the new min query time if the new value is lower than the current one", func(t *testing.T) {
		metrics := DatabaseMetrics{MinQueryTime: 10}

		metrics.SetMinQueryTime(5)

		assert.Equal(t, time.Duration(5), metrics.MinQueryTime)
	})

	t.Run("Should not set the new min query time if the new value is greater than the current one", func(t *testing.T) {
		metrics := DatabaseMetrics{MinQueryTime: 10}

		metrics.SetMinQueryTime(15)

		assert.Equal(t, time.Duration(10), metrics.MinQueryTime)
	})

	t.Run("Should set the new max query time if the new value is greater than the current one", func(t *testing.T) {
		metrics := DatabaseMetrics{MaxQueryTime: 10}

		metrics.SetMaxQueryTime(15)

		assert.Equal(t, time.Duration(15), metrics.MaxQueryTime)
	})

	t.Run("Should not set the new max query time if the new value is lower than the current one", func(t *testing.T) {
		metrics := DatabaseMetrics{MaxQueryTime: 10}

		metrics.SetMaxQueryTime(5)

		assert.Equal(t, time.Duration(10), metrics.MaxQueryTime)
	})

	t.Run("Should set the average query time", func(t *testing.T) {
		metrics := DatabaseMetrics{TotalQueries: 10, TotalTime: 100}

		metrics.SetAverageQueryTime()

		assert.Equal(t, time.Duration(10), metrics.AvgQueryTime)
	})

	t.Run("Should set the median query time", func(t *testing.T) {
		metrics := DatabaseMetrics{TotalTime: 100}

		metrics.SetMedianQueryTime()

		assert.Equal(t, time.Duration(50), metrics.MedianQueryTime)
	})

	t.Run("Should increment the total queries", func(t *testing.T) {
		metrics := DatabaseMetrics{}

		metrics.IncrementTotalQueries()

		assert.Equal(t, 1, metrics.TotalQueries)

		metrics.IncrementTotalQueries()

		assert.Equal(t, 2, metrics.TotalQueries)
	})

	t.Run("Should increment the total time", func(t *testing.T) {
		metrics := DatabaseMetrics{}

		metrics.IncrementTotalTime(time.Duration(10))

		assert.Equal(t, time.Duration(10), metrics.TotalTime)

		metrics.IncrementTotalTime(time.Duration(20))

		assert.Equal(t, time.Duration(30), metrics.TotalTime)
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

		assert.Equal(t, time.Duration(0), metrics.AvgQueryTime)
	})
}
