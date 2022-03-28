package postgres

import (
	"database/sql"

	"github.com/groupe-edf/watchdog/internal/core/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (postgres *PostgresStore) GetSettings(q *query.Query) ([]*models.Setting, error) {
	var settings = make([]*models.Setting, 0)
	queryBuilder := builder.Select([]string{
		`"settings"."id"`,
		`"settings"."container_id"`,
		`"settings"."container_type"`,
		`"settings"."setting_key"`,
		`"settings"."setting_type"`,
		`"settings"."setting_value"`,
	}...).
		From("settings").
		WithRouteQuery(q)
	statement, err := queryBuilder.ToBoundSQL()
	if err != nil {
		return settings, err
	}
	rows, err := postgres.database.Query(statement)
	if err != nil {
		return settings, err
	}
	defer rows.Close()
	for rows.Next() {
		var containerID sql.NullInt64
		var setting models.Setting
		err = rows.Scan(
			&setting.ID,
			&containerID,
			&setting.ContainerType,
			&setting.SettingKey,
			&setting.SettingType,
			&setting.SettingValue,
		)
		setting.ContainerID = containerID.Int64
		settings = append(settings, &setting)
		if err != nil {
			return settings, err
		}
	}
	return settings, nil
}
