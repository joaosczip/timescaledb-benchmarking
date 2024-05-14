package repositories

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joaosczip/timescale/internal/dtos"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCpuUsageRepository_QueryStatistics(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	repo := NewCpuUsageRepository(db)

	testCases := []struct {
		host          string
		startTime     string
		endTime       string
		expectedQuery *sqlmock.ExpectedQuery
		expected      []dtos.CpuUsageStatistics
	}{
		{"host00001", "2024-05-01 08:00:00", "2024-05-01 08:03:00", mock.
			ExpectQuery(".*").
			WillReturnRows(sqlmock.NewRows([]string{"host", "window_start", "max_usage", "min_usage"}).
				AddRows([]driver.Value{"host00001", "00:00:00", 37.1, 10.2}).
				AddRows([]driver.Value{"host00001", "00:01:00", 92.13, 26.77}).
				AddRows([]driver.Value{"host00001", "00:02:00", 20.19, 19.68})), []dtos.CpuUsageStatistics{
			{Host: "host00001", WindowStart: "00:00:00", Max: 37.1, Min: 10.2},
			{Host: "host00001", WindowStart: "00:01:00", Max: 92.13, Min: 26.77},
			{Host: "host00001", WindowStart: "00:02:00", Max: 20.19, Min: 19.68},
		}},
		{"host00002", "2024-05-01 08:00:00", "2024-05-01 08:03:00", mock.
			ExpectQuery(".*").
			WillReturnRows(sqlmock.NewRows([]string{"host", "window_start", "max_usage", "min_usage"}).
				AddRows([]driver.Value{"host00002", "00:00:00", 87.3, 80.01}).
				AddRows([]driver.Value{"host00002", "00:01:00", 40.41, 34.90}).
				AddRows([]driver.Value{"host00002", "00:02:00", 29.20, 10.05})), []dtos.CpuUsageStatistics{
			{Host: "host00002", WindowStart: "00:00:00", Max: 87.3, Min: 80.01},
			{Host: "host00002", WindowStart: "00:01:00", Max: 40.41, Min: 34.90},
			{Host: "host00002", WindowStart: "00:02:00", Max: 29.20, Min: 10.05},
		}},
	}

	for _, testCase := range testCases {
		cpuUsage, err := repo.QueryStatistics(testCase.host, testCase.startTime, testCase.endTime)
		assert.Nil(t, err)
		assert.Len(t, cpuUsage, 3)
		assert.Equal(t, cpuUsage, testCase.expected)
	}
}
