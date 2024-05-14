package dtos

type CpuUsageStatistics struct {
	StartTime string
	EndTime   string
	Max       float64
	Min       float64
}

type CpuUsageQueryParams struct {
	Host      string
	StartTime string
	EndTime   string
}
