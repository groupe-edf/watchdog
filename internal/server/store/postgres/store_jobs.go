package postgres

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/groupe-edf/watchdog/internal/models"
	builder "github.com/groupe-edf/watchdog/internal/server/database/query"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (postgres *PostgresStore) DeleteJob(job *models.Job) error {
	job.Lock.Lock()
	defer job.Lock.Unlock()
	if job.Deleted {
		return nil
	}
	statement := `
		DELETE FROM jobs
		WHERE
			id = $1 AND
			priority = $2 AND
			queue = $3 AND
			started_at = $4
	`
	_, err := postgres.database.Exec(statement,
		job.ID,
		job.Priority,
		job.Queue,
		job.StartedAt,
	)
	if err != nil {
		return err
	}
	job.Deleted = true
	return nil
}

func (postgres *PostgresStore) DoneJob(job *models.Job) {
	job.Lock.Lock()
	defer job.Lock.Unlock()
	statement := fmt.Sprintf(`SELECT pg_advisory_unlock(%v)`, job.ID)
	var ok bool
	_ = postgres.database.QueryRow(statement).Scan(ok)
}

func (postgres *PostgresStore) Enqueue(job *models.Job) error {
	statement := `INSERT INTO "jobs" (
		"args",
		"priority",
		"queue",
		"started_at",
		"type"
	) VALUES ($1, $2, $3, $4, $5) RETURNING "id"`
	_, err := postgres.database.Exec(statement,
		job.Args,
		job.Priority,
		job.Queue,
		job.StartedAt,
		job.Type,
	)
	if err != nil {
		return err
	}
	return nil
}

func (postgres *PostgresStore) FindJobs(q *query.Query) (models.Paginator[*models.Job], error) {
	paginator := models.Paginator[*models.Job]{
		Items: make([]*models.Job, 0),
		Query: q,
	}
	queryBuilder := builder.Select([]string{
		`"jobs"."id"`,
		`"jobs"."args"`,
		`"jobs"."error_count"`,
		`"jobs"."last_error"`,
		`"jobs"."priority"`,
		`"jobs"."queue"`,
		`"jobs"."started_at"`,
		`"jobs"."type"`,
		`COUNT(*) OVER() AS total_items`,
	}...).
		From("jobs").
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
		var job models.Job
		var lastError sql.NullString
		err = rows.Scan(
			&job.ID,
			&job.Args,
			&job.ErrorCount,
			&lastError,
			&job.Priority,
			&job.Queue,
			&job.StartedAt,
			&job.Type,
			&paginator.TotalItems,
		)
		if err != nil {
			return paginator, err
		}
		job.LastError = lastError.String
		paginator.Items = append(paginator.Items, &job)
	}
	return paginator, nil
}

func (postgres *PostgresStore) LockJob(queueName string) (*models.Job, error) {
	job := &models.Job{}
	statement := fmt.Sprintf(`
	WITH RECURSIVE locked_jobs AS (
		SELECT (job).*, pg_try_advisory_lock((job).id) AS locked
		FROM (
			SELECT job
			FROM jobs AS job
			WHERE queue = '%[1]s'::text
			AND started_at <= now()
			ORDER BY priority, started_at, id
			LIMIT 1
		) AS jobs_view
		UNION ALL (
			SELECT (job).*, pg_try_advisory_lock((job).id) AS locked
			FROM (
				SELECT (
					SELECT job
					FROM jobs AS job
					WHERE queue = '%[1]s'::text
					AND started_at <= now()
					AND (priority, started_at, id) > (jobs.priority, jobs.started_at, jobs.id)
					ORDER BY priority, started_at, id
					LIMIT 1
				) AS job
				FROM jobs
				WHERE jobs.id IS NOT NULL
				LIMIT 1
			) AS jobs_view
		)
	)
	SELECT id, args::jsonb, error_count, priority, queue, started_at, type
	FROM locked_jobs
	WHERE locked
	LIMIT 1
	`, queueName)
	err := postgres.database.QueryRow(statement).Scan(
		&job.ID,
		&job.Args,
		&job.ErrorCount,
		&job.Priority,
		&job.Queue,
		&job.StartedAt,
		&job.Type,
	)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (postgres *PostgresStore) SaveJobError(job *models.Job, message string) error {
	errorCount := job.ErrorCount + 1
	delay := int(math.Pow(float64(errorCount), float64(4))) + 3
	statement := `
		UPDATE jobs SET
			error_count = $1::integer,
			started_at = now() + $2::bigint * '1 second'::interval,
			last_error = $3::text
		WHERE
			queue = $4::text AND
			priority = $5::smallint AND
			started_at = $6::timestamptz AND
			id = $7::bigint
	`
	_, err := postgres.database.Exec(statement,
		errorCount,
		delay,
		message,
		job.Queue,
		job.Priority,
		job.StartedAt,
		job.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
