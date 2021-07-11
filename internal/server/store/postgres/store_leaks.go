package postgres

import (
	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
	"github.com/lib/pq"
)

func (postgres *PostgresStore) FindLeaks(q *query.Query) (models.Paginator[models.Leak], error) {
	paginator := models.Paginator[models.Leak]{
		Items: make([]models.Leak, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"repositories_leaks"."id"`,
		`"repositories_leaks"."author"`,
		`"repositories_leaks"."commit_hash"`,
		`"repositories_leaks"."created_at"`,
		`"repositories_leaks"."file"`,
		`"repositories_leaks"."secret_hash"`,
		`"repositories_leaks"."line"`,
		`"repositories_leaks"."line_number"`,
		`"repositories_leaks"."offender"`,
		`"repositories_leaks"."severity"`,
		`"unique_leaks"."occurence"`,
		`"repositories"."id"`,
		`"repositories"."repository_url"`,
		`"rules"."display_name"`,
		`"rules"."tags"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("repositories_leaks").
		Join("INNER", `(
			SELECT
				"secret_hash",
				MAX("created_at") AS "created_at",
				COUNT("secret_hash") AS "occurence"
			FROM "repositories_leaks"
			GROUP BY "secret_hash"
		) AS "unique_leaks"`, builder.Expression(`
			"repositories_leaks"."secret_hash" = "unique_leaks"."secret_hash" AND
			"repositories_leaks"."created_at" = "unique_leaks"."created_at"`)).
		Join("LEFT", "repositories", builder.Expression(`"repositories_leaks"."repository_id" = "repositories"."id"`)).
		Join("LEFT", "rules", builder.Expression(`"repositories_leaks"."rule_id" = "rules"."id"`)).
		Limit(q.Limit, q.Offset).
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
		var leak models.Leak
		var repository models.Repository
		var rule models.Rule
		err = rows.Scan(
			&leak.ID,
			&leak.Author,
			&leak.CommitHash,
			&leak.CreatedAt,
			&leak.File,
			&leak.SecretHash,
			&leak.Line,
			&leak.LineNumber,
			&leak.Offender,
			&leak.Severity,
			&leak.Occurence,
			&repository.ID,
			&repository.RepositoryURL,
			&rule.DisplayName,
			&rule.Tags,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		leak.Repository = repository
		leak.Rule = rule
		paginator.Items = append(paginator.Items, leak)
	}
	return paginator, nil
}

func (postgres *PostgresStore) FindLeakByID(id int64) (leak *models.Leak, err error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"repositories_leaks"."id"`,
		Operator: query.Equal,
		Value:    id,
	})
	q.Limit = 1
	paginator, err := postgres.FindLeaks(q)
	if err != nil {
		return leak, err
	}
	leak = &paginator.Items[0]
	return leak, err
}

func (store *PostgresStore) SaveLeaks(repositoryID *uuid.UUID, analysisID *uuid.UUID, leaks []models.Leak) error {
	database, err := store.database.EnableTx()
	if err != nil {
		return err
	}
	statement, err := database.Prepare(pq.CopyIn("repositories_leaks",
		"analysis_id",
		"author",
		"commit_hash",
		"created_at",
		"file",
		"secret_hash",
		"line",
		"line_number",
		"offender",
		"repository_id",
		"rule_id",
		"severity",
	))
	if err != nil {
		return err
	}
	for _, leak := range leaks {
		_, err = statement.Exec(
			analysisID,
			leak.Author,
			leak.CommitHash,
			leak.CreatedAt,
			leak.File,
			leak.SecretHash,
			leak.Line,
			leak.LineNumber,
			leak.Offender,
			repositoryID,
			leak.Rule.ID,
			leak.Severity,
		)
		if err != nil {
			return err
		}
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	err = statement.Close()
	if err != nil {
		return err
	}
	err = database.Commit()
	if err != nil {
		return err
	}
	return nil
}
