package handlers

import (
	"sync"
	"testing"

	"github.com/joaosczip/timescale/internal/dtos"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CpuUsageRepositoryMock struct {
	mock.Mock
	wg sync.WaitGroup
}

func (m *CpuUsageRepositoryMock) QueryStatistics(host string, startTime string, endTime string) ([]dtos.CpuUsageStatistics, error) {
	defer m.wg.Done()
	args := m.Called(host, startTime, endTime)
	return args.Get(0).([]dtos.CpuUsageStatistics), args.Error(1)
}

func TestCollectCpuUsageMetrics(t *testing.T) {
	t.Run("Should not call the repository when given an empty input", func(t *testing.T) {
		repository := new(CpuUsageRepositoryMock)

		handler := NewCollectCpuUsageMetricsHandler(repository)
		metrics := handler.Handle([]dtos.CpuUsageQueryParams{})

		assert.NotNil(t, metrics)
		repository.AssertNotCalled(t, "QueryStatistics")
	})

	t.Run("Should increment the failure metric when the repository returns an error", func(t *testing.T) {
		repository := new(CpuUsageRepositoryMock)
		repository.On("QueryStatistics", "host", "start", "end").Return([]dtos.CpuUsageStatistics{}, assert.AnError)

		handler := NewCollectCpuUsageMetricsHandler(repository)

		repository.wg.Add(1)

		handlerInput := []dtos.CpuUsageQueryParams{
			{
				Host:      "host",
				StartTime: "start",
				EndTime:   "end",
			},
		}
		metrics := handler.Handle(handlerInput)

		repository.wg.Wait()

		assert.NotNil(t, metrics)
		assert.Equal(t, 1, metrics.Failures)
		repository.AssertExpectations(t)
	})

	t.Run("Should call the repository for each one of the provided query params", func(t *testing.T) {
		handlerInput := []dtos.CpuUsageQueryParams{
			{
				Host:      "host00001",
				StartTime: "2024-05-01 08:00:00",
				EndTime:   "2024-05-01 08:03:00",
			},
			{
				Host:      "host00002",
				StartTime: "2024-05-01 08:00:00",
				EndTime:   "2024-05-01 08:03:00",
			},
			{
				Host:      "host00003",
				StartTime: "2024-05-01 08:00:00",
				EndTime:   "2024-05-01 08:03:00",
			},
			{
				Host:      "host00003",
				StartTime: "2024-05-01 09:00:00",
				EndTime:   "2024-05-01 09:03:00",
			},
		}

		repository := new(CpuUsageRepositoryMock)
		repository.
			On("QueryStatistics", handlerInput[0].Host, handlerInput[0].StartTime, handlerInput[0].EndTime).
			Return([]dtos.CpuUsageStatistics{}, nil).
			On("QueryStatistics", handlerInput[1].Host, handlerInput[1].StartTime, handlerInput[1].EndTime).
			Return([]dtos.CpuUsageStatistics{}, nil).
			On("QueryStatistics", handlerInput[2].Host, handlerInput[2].StartTime, handlerInput[2].EndTime).
			Return([]dtos.CpuUsageStatistics{}, nil).
			On("QueryStatistics", handlerInput[3].Host, handlerInput[3].StartTime, handlerInput[3].EndTime).
			Return([]dtos.CpuUsageStatistics{}, nil)

		handler := NewCollectCpuUsageMetricsHandler(repository)

		repository.wg.Add(4)
		metrics := handler.Handle(handlerInput)
		repository.wg.Wait()

		assert.NotNil(t, metrics)
		repository.AssertExpectations(t)
	})
}
