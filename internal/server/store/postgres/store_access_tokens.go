package postgres

import (
	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

// Access Tokens
func (store *PostgresStore) FindAccessTokens(q *query.Query) (models.Paginator[models.AccessToken], error) {
	paginator := models.Paginator[models.AccessToken]{
		Items: make([]models.AccessToken, 0),
		Query: q,
	}
	return paginator, nil
}

func (store *PostgresStore) FindAccessToken(token string) (*models.AccessToken, error) {
	queryBuilder := builder.Select([]string{
		`"access_tokens"."id"`,
		`"access_tokens"."expires_at"`,
		`"access_tokens"."name"`,
		`"access_tokens"."revoked"`,
		`"access_tokens"."token"`,
		`"access_tokens"."user_id"`,
	}...).
		From("access_tokens").
		Where(builder.Equal{`"access_tokens"."token"`: token})
	statement, err := queryBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	var accessToken models.AccessToken
	err = store.database.QueryRow(statement).Scan(
		&accessToken.ID,
		&accessToken.ExpiresAt,
		&accessToken.Name,
		&accessToken.Revoked,
		&accessToken.Token,
		&accessToken.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func (store *PostgresStore) SaveAccessToken(token *models.AccessToken) (*models.AccessToken, error) {
	var ID int64
	statement := `INSERT INTO "access_tokens" (
		"id",
		"created_at",
		"name",
		"token",
		"user_id"
	) VALUES (NEXTVAL('access_tokens_id_seq'), $1, $2, $3, $4)
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		token.CreatedAt,
		token.Name,
		token.Token,
		token.UserID,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	token.ID = ID
	return token, nil
}
