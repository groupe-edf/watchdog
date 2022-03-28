package postgres

import (
	"fmt"

	"github.com/groupe-edf/watchdog/internal/core/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (store *PostgresStore) AddWebhook(webhook *models.Webhook) (*models.Webhook, error) {
	var ID int64
	statement := `INSERT INTO "integrations_webhooks" (
		"integration_id",
		"group_id",
		"token",
		"url",
		"webhook_id"
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		webhook.IntegrationID,
		webhook.GroupID,
		webhook.Token,
		webhook.URL,
		webhook.WebhookID,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	webhook.ID = ID
	return webhook, nil
}

func (store *PostgresStore) DeleteIntegration(id int64) error {
	_, err := store.database.Exec(`DELETE FROM integrations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresStore) FindIntegrations(q *query.Query) (models.Paginator[models.Integration], error) {
	paginator := models.Paginator[models.Integration]{
		Items: make([]models.Integration, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"integrations"."id"`,
		`"integrations"."api_token"`,
		`"integrations"."created_at"`,
		`"integrations"."instance_name"`,
		`"integrations"."instance_url"`,
		`"integrations"."synced_at"`,
		`"integrations"."syncing_error"`,
		`"users"."id"`,
		`"users"."first_name"`,
		`"users"."last_name"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("integrations").
		Join("LEFT", "users", builder.Expression(`"integrations"."created_by" = "users"."id"`))
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
		var integration models.Integration
		var createdBy models.User
		err = rows.Scan(
			&integration.ID,
			&integration.APIToken,
			&integration.CreatedAt,
			&integration.InstanceName,
			&integration.InstanceURL,
			&integration.SyncedAt,
			&integration.SyncingError,
			&createdBy.ID,
			&createdBy.FirstName,
			&createdBy.LastName,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		integration.CreatedBy = &createdBy
		paginator.Items = append(paginator.Items, integration)
	}
	return paginator, nil
}

func (store *PostgresStore) FindIntegrationByID(id int64) (*models.Integration, error) {
	var integration models.Integration
	statement := fmt.Sprintf(`
		SELECT
			"integrations"."id",
			"integrations"."api_token",
			"integrations"."created_at",
			"integrations"."created_by",
			"integrations"."instance_name",
			"integrations"."instance_url"
		FROM "integrations"
		WHERE "integrations"."id" = '%d'`,
		id,
	)
	var createdBy models.User
	err := store.database.QueryRow(statement).Scan(
		&integration.ID,
		&integration.APIToken,
		&integration.CreatedAt,
		&createdBy.ID,
		&integration.InstanceName,
		&integration.InstanceURL,
	)
	if err != nil {
		return nil, err
	}
	integration.CreatedBy = &createdBy
	return &integration, nil
}

func (store *PostgresStore) FindWebhooks(q *query.Query) (models.Paginator[models.Webhook], error) {
	paginator := models.Paginator[models.Webhook]{
		Items: make([]models.Webhook, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"integrations_webhooks"."id"`,
		`"integrations_webhooks"."created_at"`,
		`"integrations_webhooks"."group_id"`,
		`"integrations_webhooks"."integration_id"`,
		`"integrations_webhooks"."token"`,
		`"integrations_webhooks"."url"`,
	}...).
		From("integrations_webhooks")
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
		var webhook models.Webhook
		err = rows.Scan(
			&webhook.ID,
			&webhook.CreatedAt,
			&webhook.GroupID,
			&webhook.IntegrationID,
			&webhook.Token,
			&webhook.URL,
		)
		if err != nil {
			return paginator, err
		}
		paginator.Items = append(paginator.Items, webhook)
	}
	return paginator, nil
}

func (store *PostgresStore) SaveIntegration(integration *models.Integration) (*models.Integration, error) {
	var ID int64
	statement := `INSERT INTO "integrations" (
		"api_token",
		"created_at",
		"created_by",
		"instance_name",
		"instance_url",
		"synced_at"
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		integration.APIToken,
		integration.CreatedAt,
		integration.CreatedBy.ID,
		integration.InstanceName,
		integration.InstanceURL,
		integration.SyncedAt,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	integration.ID = ID
	return integration, nil
}

func (postgres *PostgresStore) UpdateIntegration(integration *models.Integration) (*models.Integration, error) {
	statement := `UPDATE "integrations"
	SET
		"synced_at" = $2,
		"syncing_error" = $3
  WHERE id = $1`
	_, err := postgres.database.Exec(statement,
		integration.ID,
		integration.SyncedAt,
		integration.SyncingError,
	)
	if err != nil {
		return nil, err
	}
	return integration, nil
}
