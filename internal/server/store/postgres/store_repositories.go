package postgres

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/core/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (postgres *PostgresStore) DeleteRepository(id uuid.UUID) error {
	_, err := postgres.database.Exec(`DELETE FROM repositories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (postgres *PostgresStore) FindRepositories(q *query.Query) (models.Paginator[models.Repository], error) {
	paginator := models.Paginator[models.Repository]{
		Items: make([]models.Repository, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"repositories"."id"`,
		`"repositories"."created_at"`,
		`"repositories"."created_by"`,
		`"repositories"."repository_url"`,
		`"repositories"."visibility"`,
		`"integrations"."id"`,
		`"integrations"."instance_name"`,
		`"integrations"."instance_url"`,
		`"last_analysis"."id"`,
		`"last_analysis"."duration"`,
		`"last_analysis"."finished_at"`,
		`"last_analysis"."last_commit_hash"`,
		`"last_analysis"."severity"`,
		`"last_analysis"."started_at"`,
		`"last_analysis"."state"`,
		`"last_analysis"."state_message"`,
		`"last_analysis"."total_issues"`,
		`COALESCE("last_analysis"."trigger", '')`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("repositories").
		Join("LEFT", `(
			SELECT
				DISTINCT ON ("repositories_analyzes"."repository_id") *
			FROM "repositories_analyzes"
			ORDER BY "repositories_analyzes"."repository_id", "repositories_analyzes"."created_at" desc
		) AS "last_analysis"`, builder.Expression(`"repositories"."id" = "last_analysis"."repository_id"`)).
		Join("LEFT", "integrations", builder.Expression(`"repositories"."integration_id" = "integrations"."id"`)).
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
		var analysis models.Analysis
		var analysisDuration sql.NullInt64
		var analysisLastCommitHash sql.NullString
		var analysisState sql.NullString
		var analysisStateMessage sql.NullString
		var analysisTotalIssues sql.NullInt32
		var repository models.Repository
		var integrationID sql.NullInt64
		var integrationInstanceName sql.NullString
		var integrationInstanceURL sql.NullString
		err = rows.Scan(
			&repository.ID,
			&repository.CreatedAt,
			&repository.CreatedBy,
			&repository.RepositoryURL,
			&repository.Visibility,
			&integrationID,
			&integrationInstanceName,
			&integrationInstanceURL,
			&analysis.ID,
			&analysisDuration,
			&analysis.FinishedAt,
			&analysisLastCommitHash,
			&analysis.Severity,
			&analysis.StartedAt,
			&analysisState,
			&analysisStateMessage,
			&analysisTotalIssues,
			&analysis.Trigger,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		if analysis.ID != nil {
			analysis.Duration = time.Duration(analysisDuration.Int64)
			analysis.LastCommitHash = analysisLastCommitHash.String
			analysis.State = models.AnalysisState(analysisState.String)
			analysis.StateMessage = analysisStateMessage.String
			analysis.TotalIssues = int(analysisTotalIssues.Int32)
			repository.LastAnalysis = &analysis
		}
		if integrationID.Valid {
			repository.Integration = &models.Integration{
				ID:           integrationID.Int64,
				InstanceName: integrationInstanceName.String,
				InstanceURL:  integrationInstanceURL.String,
			}
		}
		paginator.Items = append(paginator.Items, repository)
	}
	return paginator, nil
}

func (postgres *PostgresStore) FindRepositoryByID(id *uuid.UUID) (*models.Repository, error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"repositories"."id"`,
		Operator: query.Equal,
		Value:    id.String(),
	})
	q.Limit = 1
	paginator, err := postgres.FindRepositories(q)
	if err != nil {
		return nil, err
	}
	repository := paginator.Items[0]
	return &repository, nil
}

func (store *PostgresStore) FindRepositoryByURI(uri string) *models.Repository {
	return nil
}

func (store *PostgresStore) SaveRepository(repository *models.Repository) (*models.Repository, error) {
	var ID uuid.UUID
	var integrationID *int64
	if repository.Integration != nil {
		integrationID = &repository.Integration.ID
	}
	statement := `INSERT INTO "repositories" (
		"id",
		"created_at",
		"created_by",
		"enable_monitoring",
		"integration_id",
		"repository_url",
		"visibility"
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT ON CONSTRAINT "repositories_repository_url_key" DO
	UPDATE SET
		"enable_monitoring" = EXCLUDED."enable_monitoring",
		"repository_url" = EXCLUDED."repository_url"
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		repository.ID,
		repository.CreatedAt,
		repository.CreatedBy,
		repository.EnableMonitoring,
		integrationID,
		repository.RepositoryURL,
		repository.Visibility,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	repository.ID = &ID
	return repository, nil
}
