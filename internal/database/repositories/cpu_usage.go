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
		host,
		FLOOR(EXTRACT(MINUTE FROM ts)) * INTERVAL '1 minute' AS window_start,
		MAX(usage) AS max_usage,
		MIN(usage) AS min_usage
	FROM 
		cpu_usage
	WHERE 
		host = $1
	AND ts BETWEEN $2 AND $3
	GROUP BY 
		host,
		FLOOR(EXTRACT(MINUTE FROM ts)) * INTERVAL '1 minute'
	ORDER BY 
		host,
		window_start;
    `

	rows, err := r.db.Query(query, host, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []dtos.CpuUsageStatistics
	for rows.Next() {
		var stat dtos.CpuUsageStatistics
		err := rows.Scan(&stat.Host, &stat.WindowStart, &stat.Max, &stat.Min)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}
