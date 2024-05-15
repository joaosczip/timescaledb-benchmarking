package commands

import (
	"testing"

	"github.com/joaosczip/timescale/internal/common"
	"github.com/joaosczip/timescale/internal/dtos"
	"github.com/joaosczip/timescale/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCollectCpuUsageMetricsHandler struct {
	mock.Mock
}

func (m *MockCollectCpuUsageMetricsHandler) Handle(queryParams []dtos.CpuUsageQueryParams) *models.DatabaseMetrics {
	args := m.Called()
	return args.Get(0).(*models.DatabaseMetrics)
}

type MockCsvReader struct {
	mock.Mock
}

func (m *MockCsvReader) Read(path string) ([]dtos.CpuUsageQueryParams, error) {
	args := m.Called()
	return args.Get(0).([]dtos.CpuUsageQueryParams), args.Error(1)
}

type MockStdoutPrinter struct {
	data map[string]string
	mock.Mock
}

func (m *MockStdoutPrinter) Print(data map[string]string) error {
	m.data = data
	args := m.Called()
	return args.Error(0)
}

func TestDatabaseMetricsCommand_Run(t *testing.T) {
	t.Run("Should return an error when there's an error reading the csv", func(t *testing.T) {
		handler := new(MockCollectCpuUsageMetricsHandler)
		csvReader := new(MockCsvReader)
		stdoutPrinter := new(MockStdoutPrinter)

		csvReader.On("Read").Return([]dtos.CpuUsageQueryParams{}, common.ErrInvalidHeader)

		command := NewDatabaseMetricsCommand(handler, csvReader, stdoutPrinter)

		err := command.Run("incorrect.csv")

		assert.NotNil(t, err)
		assert.Equal(t, common.ErrInvalidHeader, err)
	})

	t.Run("Should use the csv data to call the handler", func(t *testing.T) {
		handler := new(MockCollectCpuUsageMetricsHandler)
		csvReader := new(MockCsvReader)
		stdoutPrinter := new(MockStdoutPrinter)

		csvReader.On("Read").Return([]dtos.CpuUsageQueryParams{
			{
				Host:      "host",
				StartTime: "start",
				EndTime:   "end",
			},
		}, nil)
		handler.On("Handle").Return(&models.DatabaseMetrics{})
		stdoutPrinter.On("Print", mock.Anything).Return(nil)

		command := NewDatabaseMetricsCommand(handler, csvReader, stdoutPrinter)

		err := command.Run("correct.csv")

		assert.Nil(t, err)
		handler.AssertExpectations(t)
	})

	t.Run("Should use the data returned by the handler to print to stdout", func(t *testing.T) {
		handler := new(MockCollectCpuUsageMetricsHandler)
		csvReader := new(MockCsvReader)
		stdoutPrinter := new(MockStdoutPrinter)

		csvReader.On("Read").Return([]dtos.CpuUsageQueryParams{
			{
				Host:      "host",
				StartTime: "start",
				EndTime:   "end",
			},
		}, nil)

		databaseMetrics := &models.DatabaseMetrics{
			TotalQueries:    2,
			Failures:        0,
			TotalTime:       2.0,
			MinQueryTime:    1.0,
			MaxQueryTime:    2.0,
			AvgQueryTime:    1.5,
			MedianQueryTime: 1.5,
		}

		handler.On("Handle").Return(databaseMetrics)
		stdoutPrinter.On("Print", mock.Anything).Return(nil)

		command := NewDatabaseMetricsCommand(handler, csvReader, stdoutPrinter)

		err := command.Run("correct.csv")

		assert.Nil(t, err)
		csvReader.AssertExpectations(t)
		handler.AssertExpectations(t)

		assert.Equal(t, map[string]string{
			"total_queries":     "2",
			"failures":          "0",
			"total_time":        "2.00ms",
			"min_query_time":    "1.00ms",
			"max_query_time":    "2.00ms",
			"avg_query_time":    "1.50ms",
			"median_query_time": "1.50ms",
		}, stdoutPrinter.data)
	})

	t.Run("Should return an error when there's an error printing to stdout", func(t *testing.T) {
		handler := new(MockCollectCpuUsageMetricsHandler)
		csvReader := new(MockCsvReader)
		stdoutPrinter := new(MockStdoutPrinter)

		csvReader.On("Read").Return([]dtos.CpuUsageQueryParams{
			{
				Host:      "host",
				StartTime: "start",
				EndTime:   "end",
			},
		}, nil)

		handler.On("Handle").Return(&models.DatabaseMetrics{})
		stdoutPrinter.On("Print", mock.Anything).Return(assert.AnError)

		command := NewDatabaseMetricsCommand(handler, csvReader, stdoutPrinter)

		err := command.Run("correct.csv")

		assert.NotNil(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}
