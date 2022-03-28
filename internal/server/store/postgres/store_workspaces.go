package postgres

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (postgres *PostgresStore) FindWorkspaces(q *query.Query) (models.Paginator[models.Workspace], error) {
	paginator := models.Paginator[models.Workspace]{
		Items: make([]models.Workspace, 0),
		Query: q,
	}
	return paginator, nil
}
