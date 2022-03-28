package postgres

import (
	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/core/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (postgres *PostgresStore) DeleteUsers(q *query.Query) error {
	return nil
}

func (postgres *PostgresStore) FindUsers(q *query.Query) (models.Paginator[models.User], error) {
	paginator := models.Paginator[models.User]{
		Items: make([]models.User, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"users"."id"`,
		`"users"."created_at"`,
		`"users"."email"`,
		`"users"."first_name"`,
		`"users"."last_login"`,
		`"users"."last_name"`,
		`"users"."locked"`,
		`"users"."password"`,
		`"users"."provider"`,
		`"users"."state"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("users").
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
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Email,
			&user.FirstName,
			&user.LastLogin,
			&user.LastName,
			&user.Locked,
			&user.Password,
			&user.Provider,
			&user.State,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		paginator.Items = append(paginator.Items, user)
	}
	return paginator, nil
}

func (postgres *PostgresStore) FindUserByEmail(email string) (*models.User, error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"users"."email"`,
		Operator: query.Equal,
		Value:    email,
	})
	q.Limit = 1
	paginator, err := postgres.FindUsers(q)
	if err != nil {
		return nil, err
	}
	if len(paginator.Items) > 0 {
		user := paginator.Items[0]
		return &user, nil
	}
	return nil, models.ErrUserNotFound
}

func (postgres *PostgresStore) FindUserById(id *uuid.UUID) (*models.User, error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"users"."id"`,
		Operator: query.Equal,
		Value:    id.String(),
	})
	q.Limit = 1
	paginator, err := postgres.FindUsers(q)
	if err != nil {
		return nil, err
	}
	if len(paginator.Items) > 0 {
		user := paginator.Items[0]
		return &user, nil
	}
	return nil, models.ErrUserNotFound
}

func (postgres *PostgresStore) SaveUser(user *models.User) (*models.User, error) {
	var ID uuid.UUID
	statement := `INSERT INTO "users" (
		"id",
		"created_at",
		"email",
		"first_name",
		"last_login",
		"last_name",
		"password",
		"provider",
		"username"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING "id"`
	err := postgres.database.QueryRow(statement,
		user.ID,
		user.CreatedAt,
		user.Email,
		user.FirstName,
		user.LastLogin,
		user.LastName,
		user.Password,
		user.Provider,
		user.Username,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	user.ID = &ID
	return user, nil
}

func (postgres *PostgresStore) SaveOrUpdateUser(user *models.User) (*models.User, error) {
	var ID uuid.UUID
	statement := `INSERT INTO "users" (
		"id",
		"created_at",
		"email",
		"first_name",
		"last_login",
		"last_name",
		"password",
		"provider",
		"username"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT ON CONSTRAINT "users_email_key" DO
	UPDATE SET
		"last_login" = EXCLUDED."last_login",
		"username" = EXCLUDED."username"
	RETURNING "id"`
	err := postgres.database.QueryRow(statement,
		user.ID,
		user.CreatedAt,
		user.Email,
		user.FirstName,
		user.LastLogin,
		user.LastName,
		user.Password,
		user.Provider,
		user.Username,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	user.ID = &ID
	return user, nil
}

func (postgres *PostgresStore) UpdatePassword(user *models.User) (*models.User, error) {
	statement := `UPDATE "users"
	SET
		"password" = $2
  WHERE id = $1`
	_, err := postgres.database.Exec(statement,
		user.ID,
		user.Password,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
