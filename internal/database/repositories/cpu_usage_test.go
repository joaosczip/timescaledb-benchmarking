package repositories

import (
	"database/sql"
	"testing"

	"github.com/joaosczip/timescale/internal/dtos"
	"github.com/joaosczip/timescale/test"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCpuUsageRepository_QueryStatistics(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.Nil(t, err)
	defer db.Close()

	err = test.MigrateTables(db)
	assert.Nil(t, err)

	err = test.PopulateSampleData(db)
	assert.Nil(t, err)

	repo := NewCpuUsageRepository(db)

	testCases := []struct {
		host      string
		startTime string
		endTime   string
		expected  []dtos.CpuUsageStatistics
	}{
		{"host00001", "2024-05-01 08:00:00", "2024-05-01 08:03:00", []dtos.CpuUsageStatistics{
			{StartTime: "2024-05-01 08:00:00", EndTime: "2024-05-01 08:00:59", Max: 37.1, Min: 10.2},
			{StartTime: "2024-05-01 08:01:00", EndTime: "2024-05-01 08:01:59", Max: 92.13, Min: 26.77},
			{StartTime: "2024-05-01 08:02:00", EndTime: "2024-05-01 08:02:59", Max: 20.19, Min: 19.68},
		}},
		{"host00002", "2024-05-01 08:00:00", "2024-05-01 08:03:00", []dtos.CpuUsageStatistics{
			{StartTime: "2024-05-01 08:00:00", EndTime: "2024-05-01 08:00:59", Max: 87.3, Min: 80.01},
			{StartTime: "2024-05-01 08:01:00", EndTime: "2024-05-01 08:01:59", Max: 40.41, Min: 34.90},
			{StartTime: "2024-05-01 08:02:00", EndTime: "2024-05-01 08:02:59", Max: 29.20, Min: 10.05},
		}},
	}

	for _, testCase := range testCases {
		cpuUsage, err := repo.QueryStatistics(testCase.host, testCase.startTime, testCase.endTime)
		assert.Nil(t, err)
		assert.Len(t, cpuUsage, 3)
		assert.Equal(t, cpuUsage, testCase.expected)
	}
}
