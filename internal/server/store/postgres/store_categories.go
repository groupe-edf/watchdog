package postgres

import (
	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (postgres *PostgresStore) FindCategories(q *query.Query) (models.Paginator[models.Category], error) {
	paginator := models.Paginator[models.Category]{
		Items: make([]models.Category, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"categories"."id"`,
		`"categories"."extension"`,
		`"categories"."level"`,
		`"categories"."lft"`,
		`"categories"."rgt"`,
		`"categories"."title"`,
		`"categories"."value"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("categories")
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
		var category models.Category
		err = rows.Scan(
			&category.ID,
			&category.Extension,
			&category.Level,
			&category.Left,
			&category.Right,
			&category.Title,
			&category.Value,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		paginator.Items = append(paginator.Items, category)
	}
	return paginator, nil
}
