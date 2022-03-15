package postgres

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (postgres *PostgresStore) FindIssues(q *query.Query) (models.Paginator[models.Issue], error) {
	paginator := models.Paginator[models.Issue]{
		Items: make([]models.Issue, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"repositories_issues"."id"`,
		`"repositories_issues"."author"`,
		`"repositories_issues"."commit_hash"`,
		`"repositories_issues"."condition_type"`,
		`"repositories_issues"."email"`,
		`"repositories_issues"."offender_object"`,
		`"repositories_issues"."offender_operand"`,
		`"repositories_issues"."offender_operator"`,
		`"repositories_issues"."offender_value"`,
		`"repositories_issues"."policy_type"`,
		`"repositories_issues"."severity"`,
		`"repositories"."id"`,
		`"repositories"."repository_url"`,
		`"policies"."enabled"`,
		`"policies"."display_name"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("repositories_issues").
		Join("LEFT", "repositories", builder.Expression(`"repositories_issues"."repository_id" = "repositories"."id"`)).
		Join("LEFT", "policies", builder.Expression(`"repositories_issues"."policy_id" = "policies"."id"`)).
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
		var commit models.Commit
		var issue models.Issue
		var policy models.Policy
		var repository models.Repository
		offender := struct {
			Object   sql.NullString
			Operand  sql.NullString
			Operator sql.NullString
			Value    sql.NullString
		}{}
		err = rows.Scan(
			&issue.ID,
			&commit.Author,
			&commit.Hash,
			&issue.ConditionType,
			&commit.Author.Email,
			&offender.Object,
			&offender.Operand,
			&offender.Operator,
			&offender.Value,
			&issue.PolicyType,
			&issue.Severity,
			&repository.ID,
			&repository.RepositoryURL,
			&policy.Enabled,
			&policy.DisplayName,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		issue.Commit = commit
		issue.Offender = &models.Offender{
			Object:   offender.Object.String,
			Operand:  offender.Operand.String,
			Operator: offender.Operator.String,
			Value:    offender.Value.String,
		}
		issue.Policy = policy
		issue.Repository = &repository
		paginator.Items = append(paginator.Items, issue)
	}
	return paginator, nil
}

func (store *PostgresStore) SaveIssue(repositoryID *uuid.UUID, analysisID *uuid.UUID, issue models.Issue) error {
	createdAt := time.Now()
	ID := uuid.New()
	statement := `
	INSERT INTO "repositories_issues" (
		"id",
		"analysis_id",
		"author",
		"commit_hash",
		"condition_type",
		"created_at",
		"email",
		"offender_object",
		"offender_operand",
		"offender_operator",
		"offender_value",
		"policy_id",
		"policy_type",
		"repository_id",
		"severity"
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := store.database.Exec(statement,
		ID,
		analysisID,
		issue.Commit.Author.Name,
		issue.Commit.Hash,
		issue.ConditionType,
		createdAt,
		issue.Commit.Author.Email,
		issue.Offender.Object,
		issue.Offender.Operand,
		issue.Offender.Operator,
		issue.Offender.Value,
		issue.Policy.ID,
		issue.Policy.Type,
		repositoryID,
		issue.Severity,
	)
	if err != nil {
		return err
	}
	return nil
}
