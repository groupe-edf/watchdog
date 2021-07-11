package postgres

import "github.com/groupe-edf/watchdog/internal/models"

func (store *PostgresStore) RefreshAnalytics() error {
	statement := `SELECT analytics_aggregation()`
	_, err := store.database.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) GetAnalytics() ([]models.AnalyticsData, error) {
	statement := `SELECT
		"window_value" AS x,
		"container_count" AS y
	FROM "analytics"
	ORDER BY "window_value" ASC`
	return store.queryAnalytics(statement)
}

func (store *PostgresStore) GetLeakCountBySeverity() ([]models.AnalyticsData, error) {
	statement := `SELECT
		SEVERITY,
		COUNT(*)
	FROM "repositories_leaks"
	GROUP BY "severity"`
	return store.queryAnalytics(statement)
}

func (store *PostgresStore) queryAnalytics(statement string) ([]models.AnalyticsData, error) {
	data := make([]models.AnalyticsData, 0)
	rows, err := store.database.Query(statement)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry models.AnalyticsData
		err = rows.Scan(
			&entry.X,
			&entry.Y,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, entry)
	}
	return data, nil
}
