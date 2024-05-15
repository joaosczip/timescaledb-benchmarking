package common

import (
	"encoding/csv"
	"errors"
	"os"

	"github.com/joaosczip/timescale/internal/dtos"
)

var ErrInvalidHeader = errors.New("invalid csv header")

type QueryParamsCsvReader[T dtos.CpuUsageQueryParams] struct{}

func NewQueryParamsCsvReader() *QueryParamsCsvReader[dtos.CpuUsageQueryParams] {
	return &QueryParamsCsvReader[dtos.CpuUsageQueryParams]{}
}

func (QueryParamsCsvReader[T]) Read(path string) ([]dtos.CpuUsageQueryParams, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()

	if err != nil {
		return nil, err
	}

	if len(header) != 3 {
		return nil, ErrInvalidHeader
	}

	valid_headers := map[int]string{
		0: "hostname",
		1: "start_time",
		2: "end_time",
	}

	for i, column := range header {
		if column != valid_headers[i] {
			return nil, ErrInvalidHeader
		}
	}

	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	var queryParams []dtos.CpuUsageQueryParams

	for _, record := range records {
		queryParams = append(queryParams, dtos.CpuUsageQueryParams{
			Host:      record[0],
			StartTime: record[1],
			EndTime:   record[2],
		})
	}

	return queryParams, nil
}
