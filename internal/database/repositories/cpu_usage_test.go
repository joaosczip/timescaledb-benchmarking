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
		{"2024-05-01 08:00:30", "host00001", 10.2},
		{"2024-05-01 08:01:00", "host00001", 44.91},
		{"2024-05-01 08:01:30", "host00001", 92.13},
		{"2024-05-01 08:01:40", "host00001", 63.30},
		{"2024-05-01 08:01:50", "host00001", 26.77},
		{"2024-05-01 08:02:00", "host00001", 20.19},
		{"2024-05-01 08:02:30", "host00001", 19.68},
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

	expected := []dtos.CpuUsageStatistics{
		{StartTime: "2024-05-01 08:00:00", EndTime: "2024-05-01 08:00:59", Max: 37.1, Min: 10.2},
		{StartTime: "2024-05-01 08:01:00", EndTime: "2024-05-01 08:01:59", Max: 92.13, Min: 26.77},
		{StartTime: "2024-05-01 08:02:00", EndTime: "2024-05-01 08:02:59", Max: 20.19, Min: 19.68},
	}

	repo := NewCpuUsageRepository(db)
	cpuUsage, err := repo.QueryStatistics("host00001", "2024-05-01 08:00:00", "2024-05-01 08:03:00")
	assert.Nil(t, err)
	assert.Len(t, cpuUsage, 3)
	assert.Equal(t, cpuUsage, expected)
}
