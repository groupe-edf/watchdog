package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/security"
	_ "github.com/lib/pq"
)

const (
	// ConnectionTimeout connection timeout
	ConnectionTimeout = 30 * time.Second
)

type Repository struct {
	ID            *uuid.UUID `json:"id"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	LastAnalysis  *Analysis  `json:"last_analysis,omitempty"`
	RepositoryURL string     `json:"repository_url"`
}

type Analysis struct {
	ID           *uuid.UUID    `json:"id"`
	CreatedAt    *time.Time    `json:"created_at,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	FinishedAt   *time.Time    `json:"finished_at,omitempty"`
	RepositoryID *uuid.UUID    `json:"repository_id,omitempty"`
	Severity     issue.Score   `json:"severity"`
	StartedAt    *time.Time    `json:"started_at,omitempty"`
	State        string        `json:"state,omitempty"`
	TotalIssues  int           `json:"total_issues"`
}

type Store struct {
	database *sql.DB
	logger   logging.Interface
}

func (store *Store) FindIssues() ([]issue.Issue, error) {
	var issues []issue.Issue
	statement := `
		SELECT
			"repositories_issues"."id",
			"repositories_issues"."author",
			"repositories_issues"."commit_hash",
			"repositories_issues"."severity"
		FROM "repositories_issues"`
	rows, err := store.database.Query(statement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row issue.Issue
		err = rows.Scan(
			&row.ID,
			&row.Author,
			&row.Commit,
			&row.Severity,
		)
		if err != nil {
			return nil, err
		}
		issues = append(issues, row)
	}
	return issues, nil
}

func (store *Store) FindRepositories() ([]Repository, error) {
	var repositories []Repository
	statement := `
		SELECT
			"repositories"."id",
			"repositories"."created_by",
			"repositories"."repository_url",
			"last_analysis"."id",
			"last_analysis"."duration",
			"last_analysis"."finished_at",
			"last_analysis"."severity",
			"last_analysis"."started_at",
			"last_analysis"."state",
			"last_analysis"."total_issues"
		FROM "repositories"
		JOIN (
			SELECT
				DISTINCT ON ("repositories_analyzes"."repository_id") *
			FROM "repositories_analyzes"
			ORDER BY "repositories_analyzes"."repository_id", "repositories_analyzes"."created_at" desc
		) AS "last_analysis" ON (
			"repositories"."id" = "last_analysis"."repository_id"
		)`
	rows, err := store.database.Query(statement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var analysis Analysis
		var repository Repository
		err = rows.Scan(
			&repository.ID,
			&repository.CreatedBy,
			&repository.RepositoryURL,
			&analysis.ID,
			&analysis.Duration,
			&analysis.FinishedAt,
			&analysis.Severity,
			&analysis.StartedAt,
			&analysis.State,
			&analysis.TotalIssues,
		)
		if err != nil {
			return nil, err
		}
		repository.LastAnalysis = &analysis
		repositories = append(repositories, repository)
	}
	return repositories, nil
}

func (store *Store) FindRepositoryByID(id *uuid.UUID) (*Repository, error) {
	var repository Repository
	statement := fmt.Sprintf(`
		SELECT
			"repositories"."id",
			"repositories"."created_by",
			"repositories"."repository_url"
		FROM "repositories"
		WHERE "repositories"."id" = '%s'`,
		id.String(),
	)
	err := store.database.QueryRow(statement).Scan(
		&repository.ID,
		&repository.CreatedBy,
		&repository.RepositoryURL,
	)
	if err != nil {
		return nil, err
	}
	return &repository, nil
}

func (store *Store) FindRepositoryByURI(uri string) Repository {
	var repository Repository
	return repository
}

func (store *Store) GetRules(uri string) []security.Rule {
	var rules []security.Rule
	return rules
}

func (store *Store) SaveAnalysis(analysis *Analysis) (*Analysis, error) {
	var ID uuid.UUID
	statement := `INSERT INTO "repositories_analyzes" (
		"id",
		"created_at",
		"duration",
		"finished_at",
		"repository_id",
		"severity",
		"started_at",
		"state",
		"total_issues"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT ON CONSTRAINT "repositories_analyzes_primary_key" DO
	UPDATE SET
		"duration" = EXCLUDED."duration",
		"finished_at" = EXCLUDED."finished_at",
		"state" = EXCLUDED."state",
		"total_issues" = EXCLUDED."total_issues"
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		analysis.ID,
		analysis.CreatedAt,
		analysis.Duration,
		analysis.FinishedAt,
		analysis.RepositoryID,
		analysis.Severity,
		analysis.StartedAt,
		analysis.State,
		analysis.TotalIssues,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

func (store *Store) SaveIssue(issueID *uuid.UUID, analysisID *uuid.UUID, data issue.Issue) error {
	for _, leak := range data.Leaks {
		statement := `
		INSERT INTO "repositories_issues" (
			"id",
			"author",
			"commit_hash",
			"email",
			"file",
			"line",
			"line_number",
			"analysis_id",
			"severity"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		err := store.database.QueryRow(statement,
			issueID,
			data.Author,
			data.Commit,
			data.Email,
			leak.File,
			leak.Line,
			leak.LineNumber,
			analysisID,
			leak.Severity,
		)
		if err != nil {
			return err.Err()
		}
	}
	return nil
}

func (store *Store) SaveRepository(repository *Repository) (*Repository, error) {
	var ID uuid.UUID
	statement := `INSERT INTO "repositories" (
		"id",
		"repository_url"
	) VALUES ($1, $2)
	ON CONFLICT ON CONSTRAINT "repositories_repository_url_key" DO
	UPDATE SET
		"repository_url" = EXCLUDED."repository_url"
	RETURNING "id"`
	err := store.database.QueryRow(statement,
		repository.ID,
		repository.RepositoryURL,
	).Scan(&ID)
	if err != nil {
		return nil, err
	}
	repository.ID = &ID
	return repository, nil
}

func NewStore(logger logging.Interface, settings *config.Database) (*Store, error) {
	logger.Info("setting up database connection")
	database, err := createLoopConnection(settings)
	if err != nil {
		return nil, err
	}
	logger.Info("applying all up migrations")
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	migrater, err := migrate.NewWithDatabaseInstance("file://web/database/migrations", "postgres", driver)
	if err := migrater.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	return &Store{
		database: database,
		logger:   logger,
	}, nil
}

func createLoopConnection(settings *config.Database) (*sql.DB, error) {
	var err error
	dataSourceName := buildConnectionString(settings)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(ConnectionTimeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, err
		case <-ticker.C:
			db, _ := sql.Open("postgres", dataSourceName)
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
	}
}

func buildConnectionString(settings *config.Database) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		settings.Host,
		settings.Port,
		settings.Username,
		settings.Password,
		settings.Name)
}
