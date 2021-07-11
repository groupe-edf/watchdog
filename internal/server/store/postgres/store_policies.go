package postgres

import (
	"database/sql"
	"time"

	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (postgres *PostgresStore) FindPolicies(q *query.Query) (models.Paginator[models.Policy], error) {
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
		`"policies"."name"`,
		`"policies"."type"`,
		`"policies_conditions"."pattern"`,
		`"policies_conditions"."type"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("policies").
		Join("LEFT", "policies_conditions", builder.Expression(`"policies"."id" = "policies_conditions"."policy_id"`)).
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
			&policy.Type,
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
				Type:        policy.Type,
			}
		}
		policiesMap[policy.ID].Conditions = append(policiesMap[policy.ID].Conditions, models.Condition{
			Pattern: policy.CondititonPattern.String,
			Type:    models.ConditionType(policy.ConditionType.String),
		})
	}
	for _, value := range policiesMap {
		paginator.Items = append(paginator.Items, *value)
	}
	return paginator, nil
}

func (postgres *PostgresStore) FindPolicyByID(id int64) (policy *models.Policy, err error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"policies"."id"`,
		Operator: query.Equal,
		Value:    id,
	})
	q.Limit = 1
	paginator, err := postgres.FindPolicies(q)
	if err != nil {
		return policy, err
	}
	policy = &paginator.Items[0]
	return policy, err
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
