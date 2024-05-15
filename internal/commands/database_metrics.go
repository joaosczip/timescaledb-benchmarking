package commands

import (
	"github.com/joaosczip/timescale/internal/common"
	"github.com/joaosczip/timescale/internal/dtos"
	"github.com/joaosczip/timescale/internal/handlers"
	"github.com/joaosczip/timescale/internal/models"
)

type DatabaseMetricsCommand struct {
	handler       handlers.Handler[dtos.CpuUsageQueryParams, models.DatabaseMetrics]
	csvReader     common.CsvReader[dtos.CpuUsageQueryParams]
	stdoutPrinter common.StdoutPrinter
}

func NewDatabaseMetricsCommand(
	handler handlers.Handler[dtos.CpuUsageQueryParams, models.DatabaseMetrics],
	csvReader common.CsvReader[dtos.CpuUsageQueryParams],
	stdoutPrinter common.StdoutPrinter,
) *DatabaseMetricsCommand {
	return &DatabaseMetricsCommand{
		handler:       handler,
		csvReader:     csvReader,
		stdoutPrinter: stdoutPrinter,
	}
}

func (c *DatabaseMetricsCommand) Run(queryParamsFilePath string) error {
	queryParams, err := c.csvReader.Read(queryParamsFilePath)

	if err != nil {
		return err
	}

	metrics := c.handler.Handle(queryParams)

	if err := c.stdoutPrinter.Print(metrics.ToMap()); err != nil {
		return err
	}

	return nil
}
