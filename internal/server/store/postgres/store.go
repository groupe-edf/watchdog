package postgres

import (
	"database/sql"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/database"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/pkg/query"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrations embed.FS
var _ models.Store = &PostgresStore{}

type PostgresStore struct {
	database *database.SQLStore
}

func (postgres *PostgresStore) Count(container string, q *query.Query) (count int, err error) {
	queryBuilder := builder.Select("COUNT(*)").From(container)
	statement, err := queryBuilder.ToBoundSQL()
	if err != nil {
		return count, err
	}
	postgres.database.QueryRow(statement).Scan(&count)
	return count, err
}

func (postgres *PostgresStore) GetHealth() error {
	_, err := postgres.database.Exec("SELECT 1")
	return err
}

// NewMigrater return new migrater
func NewMigrater(database *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	directory, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, err
	}
	migrater, err := migrate.NewWithInstance("iofs", directory, "postgres", driver)
	if err != nil {
		return nil, err
	}
	return migrater, nil
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	migrater, err := NewMigrater(db)
	if err != nil {
		return nil, err
	}
	if err := migrater.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	return &PostgresStore{
		database: &database.SQLStore{
			DB: db,
		},
	}, nil
}
