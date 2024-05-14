package dtos

type CpuUsageStatistics struct {
	Host        string
	WindowStart string
	Max         float64
	Min         float64
}

type CpuUsageQueryParams struct {
	Host      string
	StartTime string
	EndTime   string
}
