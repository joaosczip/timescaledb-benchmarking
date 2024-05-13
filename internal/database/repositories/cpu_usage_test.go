package repositories

import (
	"database/sql"
	"testing"

	"github.com/joaosczip/timescale/internal/dtos"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func populateSampleData(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO cpu_usage (ts, host, usage) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	data := [][]interface{}{
		{"2024-05-01 08:00:00", "host00001", 37.1},
		{"2024-05-01 08:00:10", "host00002", 80.01},
		{"2024-05-01 08:00:20", "host00002", 87.3},
		{"2024-05-01 08:00:30", "host00001", 10.2},
		{"2024-05-01 08:01:00", "host00001", 44.91},
		{"2024-05-01 08:01:10", "host00002", 35.87},
		{"2024-05-01 08:01:20", "host00002", 34.90},
		{"2024-05-01 08:01:30", "host00001", 92.13},
		{"2024-05-01 08:01:40", "host00001", 63.30},
		{"2024-05-01 08:01:50", "host00001", 26.77},
		{"2024-05-01 08:01:50", "host00002", 40.41},
		{"2024-05-01 08:02:00", "host00001", 20.19},
		{"2024-05-01 08:02:10", "host00002", 29.20},
		{"2024-05-01 08:02:20", "host00002", 20.16},
		{"2024-05-01 08:02:30", "host00001", 19.68},
		{"2024-05-01 08:02:30", "host00002", 10.05},
	}

	for _, d := range data {
		_, err = stmt.Exec(d...)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestCpuUsageRepository_QueryStatistics(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.Nil(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE cpu_usage (ts TIMESTAMPZ, host TEXT, usage DOUBLE PRECISION)")
	assert.Nil(t, err)

	err = populateSampleData(db)
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
