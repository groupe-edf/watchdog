package postgres

import (
	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (postgres *PostgresStore) GetWhitelist(q *query.Query) (models.Paginator[models.Whitelist], error) {
	paginator := models.Paginator[models.Whitelist]{
		Items: make([]models.Whitelist, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"allowed_items"."id"`,
		`"allowed_items"."paths"`,
	}...).
		From("allowed_items")
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
		var whitelist models.Whitelist
		err = rows.Scan(
			&whitelist.ID,
			&whitelist.Paths,
		)
		if err != nil {
			return paginator, err
		}
		paginator.Items = append(paginator.Items, whitelist)
	}
	return paginator, nil
}
