package postgres

import (
	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (postgres *PostgresStore) DeleteAnalysis(id uuid.UUID) error {
	_, err := postgres.database.Exec(`DELETE FROM repositories_analyzes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) FindAnalyzes(q *query.Query) (models.Paginator[models.Analysis], error) {
	paginator := models.Paginator[models.Analysis]{
		Items: make([]models.Analysis, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"repositories"."id"`,
		`"repositories"."repository_url"`,
		`"repositories_analyzes"."id"`,
		`"repositories_analyzes"."created_at"`,
		`"repositories_analyzes"."duration"`,
		`"repositories_analyzes"."finished_at"`,
		`"repositories_analyzes"."last_commit_hash"`,
		`"repositories_analyzes"."severity"`,
		`"repositories_analyzes"."started_at"`,
		`"repositories_analyzes"."state"`,
		`"repositories_analyzes"."state_message"`,
		`"repositories_analyzes"."total_issues"`,
		`"repositories_analyzes"."trigger"`,
		`"users"."id"`,
		`"users"."created_at"`,
		`"users"."first_name"`,
		`"users"."last_name"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("repositories_analyzes").
		Join("LEFT", "users", builder.Expression(`"repositories_analyzes"."created_by" = "users"."id"`)).
		Join("LEFT", "repositories", builder.Expression(`"repositories_analyzes"."repository_id" = "repositories"."id"`)).
		Limit(q.Limit, q.Offset).
		WithRouteQuery(q)
	statement, err := queryBuilder.ToBoundSQL()
	if err != nil {
		return paginator, err
	}
	rows, err := store.database.Query(statement)
	if err != nil {
		return paginator, err
	}
	defer rows.Close()
	for rows.Next() {
		var analysis models.Analysis
		var repository models.Repository
		var user models.User
		err = rows.Scan(
			&repository.ID,
			&repository.RepositoryURL,
			&analysis.ID,
			&analysis.CreatedAt,
			&analysis.Duration,
			&analysis.FinishedAt,
			&analysis.LastCommitHash,
			&analysis.Severity,
			&analysis.StartedAt,
			&analysis.State,
			&analysis.StateMessage,
			&analysis.TotalIssues,
			&analysis.Trigger,
			&user.ID,
			&user.CreatedAt,
			&user.FirstName,
			&user.LastName,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		analysis.CreatedBy = &user
		analysis.Repository = &repository
		paginator.Items = append(paginator.Items, analysis)
	}
	return paginator, nil
}

func (postgres *PostgresStore) FindAnalysisByID(id *uuid.UUID) (*models.Analysis, error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"repositories_analyzes"."id"`,
		Operator: query.Equal,
		Value:    id.String(),
	})
	q.Limit = 1
	paginator, err := postgres.FindAnalyzes(q)
	if err != nil {
		return nil, err
	}
	analysis := paginator.Items[0]
	return &analysis, nil
}

func (store *PostgresStore) SaveAnalysis(analysis *models.Analysis) (*models.Analysis, error) {
	var ID uuid.UUID
	statement := `INSERT INTO "repositories_analyzes" (
		"id",
		"created_at",
		"created_by",
		"duration",
		"finished_at",
		"last_commit_hash",
		"repository_id",
		"severity",
		"started_at",
		"state",
		"state_message",
		"total_issues",
		"trigger"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT ON CONSTRAINT "repositories_analyzes_primary_key" DO
	UPDATE SET
		"duration" = EXCLUDED."duration",
		"finished_at" = EXCLUDED."finished_at",
		"last_commit_hash" = EXCLUDED."last_commit_hash",
		"started_at" = EXCLUDED."started_at",
		"state" = EXCLUDED."state",
		"state_message" = EXCLUDED."state_message",
		"total_issues" = EXCLUDED."total_issues"
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		analysis.ID,
		analysis.CreatedAt,
		analysis.CreatedBy.ID,
		analysis.Duration,
		analysis.FinishedAt,
		analysis.LastCommitHash,
		analysis.Repository.ID,
		analysis.Severity,
		analysis.StartedAt,
		analysis.State,
		analysis.StateMessage,
		analysis.TotalIssues,
		analysis.Trigger,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

func (store *PostgresStore) UpdateAnalysis(analysis *models.Analysis) (*models.Analysis, error) {
	return nil, nil
}
