package bolt

import (
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

// Access Tokens
func (store *BoltStore) FindAccessTokens(q *query.Query) (models.Paginator[models.AccessToken], error) {
	paginator := models.Paginator[models.AccessToken]{
		Items: make([]models.AccessToken, 0),
		Query: q,
	}
	return paginator, nil
}

func (store *BoltStore) FindAccessToken(name string) (*models.AccessToken, error) {
	return nil, nil
}

func (store *BoltStore) SaveAccessToken(token *models.AccessToken) (*models.AccessToken, error) {
	return nil, nil
}
