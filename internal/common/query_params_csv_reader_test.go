package common

import (
	"testing"

	"github.com/joaosczip/timescale/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestQueryParamsCsvReader(t *testing.T) {
	t.Run("Should return an error if the input file does not exist", func(t *testing.T) {
		reader := NewQueryParamsCsvReader()
		_, err := reader.Read("non_existent_file.csv")
		assert.NotNil(t, err)
	})

	t.Run("Should return an error if the file header has less than 3 columns", func(t *testing.T) {
		reader := NewQueryParamsCsvReader()
		_, err := reader.Read("../../test/invalid_header_size.csv")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrInvalidHeader)
	})

	t.Run("Should return an error if the file header has invalid columns", func(t *testing.T) {
		reader := NewQueryParamsCsvReader()
		_, err := reader.Read("../../test/invalid_header_columns.csv")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrInvalidHeader)
	})

	t.Run("Should return the query params when the csv is valid", func(t *testing.T) {
		reader := NewQueryParamsCsvReader()

		expected := []dtos.CpuUsageQueryParams{
			{
				Host:      "host_000008",
				StartTime: "2017-01-01 08:59:22",
				EndTime:   "2017-01-01 09:59:22",
			},
			{
				Host:      "host_000001",
				StartTime: "2017-01-02 13:02:02",
				EndTime:   "2017-01-02 14:02:02",
			},
		}

		queryParams, err := reader.Read("../../test/valid_header.csv")

		assert.Nil(t, err)
		assert.NotNil(t, queryParams)

		assert.Equal(t, expected, *queryParams)
	})
}
