package repositories

import "github.com/joaosczip/timescale/internal/dtos"

type CpuUsage interface {
	QueryStatistics(host string, startTime string, endTime string) ([]dtos.CpuUsageStatistics, error)
}
