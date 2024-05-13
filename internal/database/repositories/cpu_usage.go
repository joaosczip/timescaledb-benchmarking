package repositories

import (
	"database/sql"

	"github.com/joaosczip/timescale/internal/dtos"
)

type CpuUsageRepository struct {
	db *sql.DB
}

func NewCpuUsageRepository(db *sql.DB) *CpuUsageRepository {
	return &CpuUsageRepository{db}
}

func (r *CpuUsageRepository) QueryStatistics(host string, startTime string, endTime string) ([]dtos.CpuUsageStatistics, error) {
	query := `
        SELECT 
            strftime('%Y-%m-%d %H:%M:00', ts) AS start_time,
            strftime('%Y-%m-%d %H:%M:59', ts) AS end_time,
            MAX(usage) AS max_usage,
            MIN(usage) AS min_usage
        FROM 
            cpu_usage
        WHERE 
            host = ?
            AND ts BETWEEN ? AND ?
        GROUP BY 
            strftime('%Y-%m-%d %H:%M', ts)
        ORDER BY 
            start_time
    `

	rows, err := r.db.Query(query, host, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []dtos.CpuUsageStatistics
	for rows.Next() {
		var stat dtos.CpuUsageStatistics
		err := rows.Scan(&stat.StartTime, &stat.EndTime, &stat.Max, &stat.Min)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}
