package commands

import (
	"github.com/joaosczip/timescale/internal/common"
	"github.com/joaosczip/timescale/internal/dtos"
	"github.com/joaosczip/timescale/internal/handlers"
	"github.com/joaosczip/timescale/internal/models"
)

type DatabaseMetricsCommand struct {
	handler      handlers.Handler[dtos.CpuUsageQueryParams, models.DatabaseMetrics]
	csvReader    common.CsvReader[dtos.CpuUsageQueryParams]
	stdoutWriter common.StdoutWriter
}

func NewDatabaseMetricsCommand(
	handler handlers.Handler[dtos.CpuUsageQueryParams, models.DatabaseMetrics],
	csvReader common.CsvReader[dtos.CpuUsageQueryParams],
	stdoutWriter common.StdoutWriter,
) *DatabaseMetricsCommand {
	return &DatabaseMetricsCommand{
		handler:      handler,
		csvReader:    csvReader,
		stdoutWriter: stdoutWriter,
	}
}

func (c *DatabaseMetricsCommand) Run(queryParamsFilePath string) error {
	queryParams, err := c.csvReader.Read(queryParamsFilePath)

	if err != nil {
		return err
	}

	metrics := c.handler.Handle(queryParams)

	if err := c.stdoutWriter.Write(metrics.ToMap()); err != nil {
		return err
	}

	return nil
}
