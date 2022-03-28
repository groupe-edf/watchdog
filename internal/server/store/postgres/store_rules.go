package postgres

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (postgres *PostgresStore) FindRules(q *query.Query) (models.Paginator[models.Rule], error) {
	paginator := models.Paginator[models.Rule]{
		Items: make([]models.Rule, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"rules"."id"`,
		`"rules"."created_at"`,
		`"rules"."created_by"`,
		`"rules"."display_name"`,
		`"rules"."description"`,
		`"rules"."enabled"`,
		`"rules"."file"`,
		`COALESCE("rules"."name", '')`,
		`"rules"."pattern"`,
		`"rules"."severity"`,
		`"rules"."tags"`,
	}...).
		From("rules").
		WithRouteQuery(q)
	statement, err := queryBuilder.ToBoundSQL()
	if err != nil {
		return paginator, err
	}
	rows, err := postgres.database.Query(statement)
	if err != nil {
		return paginator, err
	}
	defer rows.Close()
	for rows.Next() {
		var rule models.Rule
		err = rows.Scan(
			&rule.ID,
			&rule.CreatedAt,
			&rule.CreatedBy,
			&rule.DisplayName,
			&rule.Description,
			&rule.Enabled,
			&rule.File,
			&rule.Name,
			&rule.Pattern,
			&rule.Severity,
			&rule.Tags,
		)
		if err != nil {
			return paginator, err
		}
		paginator.Items = append(paginator.Items, rule)
	}
	return paginator, nil
}

func (store *PostgresStore) SaveRule(rule *models.Rule) (*models.Rule, error) {
	var ID int64
	statement := `INSERT INTO "rules" (
		"created_at",
		"created_by",
		"display_name",
		"description",
		"enabled",
		"pattern",
		"severity"
	) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING "id"`
	err := store.database.QueryRow(statement,
		rule.CreatedAt,
		rule.CreatedBy,
		rule.DisplayName,
		rule.Description,
		rule.Enabled,
		rule.Pattern,
		rule.Severity,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	rule.ID = ID
	return rule, nil
}

func (store *PostgresStore) ToggleRule(ruleID int64, enabled bool) error {
	statement := `
	UPDATE "rules"
	SET "enabled" = ($2)
  WHERE id = $1`
	_, err := store.database.Exec(statement,
		ruleID,
		enabled,
	)
	if err != nil {
		return err
	}
	return nil
}
