package test

import "database/sql"

func MigrateTables(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE cpu_usage (ts TIMESTAMPZ, host TEXT, usage DOUBLE PRECISION)")
	return err
}

func PopulateSampleData(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO cpu_usage (ts, host, usage) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	data := [][]interface{}{
		{"2024-05-01 08:00:00", "host00001", 37.1},
		{"2024-05-01 08:00:10", "host00002", 80.01},
		{"2024-05-01 08:00:20", "host00002", 87.3},
		{"2024-05-01 08:00:30", "host00001", 10.2},
		{"2024-05-01 08:01:00", "host00001", 44.91},
		{"2024-05-01 08:01:10", "host00002", 35.87},
		{"2024-05-01 08:01:20", "host00002", 34.90},
		{"2024-05-01 08:01:30", "host00001", 92.13},
		{"2024-05-01 08:01:40", "host00001", 63.30},
		{"2024-05-01 08:01:50", "host00001", 26.77},
		{"2024-05-01 08:01:50", "host00002", 40.41},
		{"2024-05-01 08:02:00", "host00001", 20.19},
		{"2024-05-01 08:02:10", "host00002", 29.20},
		{"2024-05-01 08:02:20", "host00002", 20.16},
		{"2024-05-01 08:02:30", "host00001", 19.68},
		{"2024-05-01 08:02:30", "host00002", 10.05},
	}

	for _, d := range data {
		_, err = stmt.Exec(d...)
		if err != nil {
			return err
		}
	}

	return nil
}
