package postgres

import (
	"database/sql"
	"time"

	"github.com/groupe-edf/watchdog/internal/core/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (store *PostgresStore) AddPolicyCondition(condition *models.Condition) (*models.Condition, error) {
	var ID int64
	statement := `INSERT INTO "policies_conditions" (
		"pattern",
		"policy_id",
		"type"
	) VALUES ($1, $2, $3)
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		condition.Pattern,
		condition.PolicyID,
		condition.Type,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	condition.ID = ID
	return condition, nil
}

func (store *PostgresStore) DeletePolicy(policyID int64) error {
	_, err := store.database.Exec(`DELETE FROM policies WHERE id = $1`, policyID)
	return err
}

func (store *PostgresStore) DeletePolicyCondition(policyID int64, conditionID int64) error {
	_, err := store.database.Exec(`DELETE FROM policies_conditions WHERE id = $1 AND policy_id = $2`, conditionID, policyID)
	return err
}

func (store *PostgresStore) FindPolicies(q *query.Query) (models.Paginator[models.Policy], error) {
	paginator := models.Paginator[models.Policy]{
		Items: make([]models.Policy, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"policies"."id"`,
		`"policies"."created_at"`,
		`"policies"."description"`,
		`"policies"."display_name"`,
		`"policies"."enabled"`,
		`COALESCE("policies"."name", '')`,
		`"policies"."severity"`,
		`"policies"."type"`,
		`"policies_conditions"."id"`,
		`"policies_conditions"."pattern"`,
		`"policies_conditions"."type"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("policies").
		Join("LEFT", "policies_conditions", builder.Expression(`"policies"."id" = "policies_conditions"."policy_id"`)).
		OrderBy(`"policies"."name" ASC`).
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
	policiesMap := make(map[int64]*models.Policy)
	for rows.Next() {
		policy := struct {
			ID                int64
			Conditions        []models.Condition
			CreatedAt         *time.Time
			Description       string
			DisplayName       string
			Enabled           bool
			Name              string
			Type              models.PolicyType
			Severity          models.Severity
			ConditionID       sql.NullInt64
			CondititonPattern sql.NullString
			ConditionType     sql.NullString
		}{}
		err = rows.Scan(
			&policy.ID,
			&policy.CreatedAt,
			&policy.Description,
			&policy.DisplayName,
			&policy.Enabled,
			&policy.Name,
			&policy.Severity,
			&policy.Type,
			&policy.ConditionID,
			&policy.CondititonPattern,
			&policy.ConditionType,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		if _, ok := policiesMap[policy.ID]; !ok {
			policiesMap[policy.ID] = &models.Policy{
				ID:          policy.ID,
				Conditions:  make([]models.Condition, 0, 1),
				Description: policy.Description,
				DisplayName: policy.DisplayName,
				Enabled:     policy.Enabled,
				Name:        policy.Name,
				Severity:    policy.Severity,
				Type:        policy.Type,
			}
		}
		if policy.ConditionID.Valid {
			policiesMap[policy.ID].Conditions = append(policiesMap[policy.ID].Conditions, models.Condition{
				ID:      policy.ConditionID.Int64,
				Pattern: policy.CondititonPattern.String,
				Type:    models.ConditionType(policy.ConditionType.String),
			})
		}
	}
	for _, value := range policiesMap {
		paginator.Items = append(paginator.Items, *value)
	}
	return paginator, nil
}

func (store *PostgresStore) FindPolicyByID(id int64) (policy *models.Policy, err error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"policies"."id"`,
		Operator: query.Equal,
		Value:    id,
	})
	q.Limit = 1
	paginator, err := store.FindPolicies(q)
	if err != nil {
		return policy, err
	}
	policy = &paginator.Items[0]
	return policy, err
}

func (store *PostgresStore) SavePolicy(policy *models.Policy) (*models.Policy, error) {
	var ID int64
	statement := `INSERT INTO "policies" (
		"description",
		"display_name",
		"type"
	) VALUES ($1, $2, $3)
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		policy.Description,
		policy.DisplayName,
		policy.Type,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	policy.ID = ID
	return policy, nil
}

func (store *PostgresStore) TogglePolicy(policyID int64, enabled bool) error {
	statement := `UPDATE "policies" SET "enabled" = ($2) WHERE id = $1`
	_, err := store.database.Exec(statement,
		policyID,
		enabled,
	)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) UpdatePolicy(policy *models.Policy) (*models.Policy, error) {
	statement := `UPDATE "policies"
	SET
		"description" = $2,
		"display_name" = $3,
		"enabled" = $4
  WHERE id = $1`
	_, err := store.database.Exec(statement,
		policy.ID,
		policy.Description,
		policy.DisplayName,
		policy.Enabled,
	)
	if err != nil {
		return nil, err
	}
	return policy, nil
}
