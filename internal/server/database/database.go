package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	QueryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "watchdog",
			Subsystem: "postgres",
			Name:      "query_total",
			Help:      "Total number of database queries",
		},
		[]string{"query"},
	)
	QueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "watchdog",
			Subsystem: "postgres",
			Name:      "query_duration_seconds",
			Help:      "The duration of all queries",
		},
		[]string{"query"},
	)
)

// Store is an interface that models the standard transaction in `database/sql`.
type Store interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Transaction
}

// SQLStore concrete implementation of Store by using *sql.DB
type SQLStore struct {
	DB *sql.DB
}

// BeginTx start database transaction
func (store *SQLStore) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return store.DB.BeginTx(ctx, opts)
}

// Exec wrapper to sql.DB.Exec
func (store *SQLStore) Exec(query string, args ...interface{}) (sql.Result, error) {
	return store.DB.Exec(query, args...)
}

// HandleErrors handle database errors
func (store *SQLStore) HandleErrors() {
}

// Prepare wrapper to sql.DB.Prepare
func (store *SQLStore) Prepare(query string) (*sql.Stmt, error) {
	return store.DB.Prepare(query)
}

// Query wrapper to sql.DB.Query
func (store *SQLStore) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := store.DB.Query(query, args...)
	QueryCounter.WithLabelValues("query").Add(1)
	QueryDuration.WithLabelValues("query").Observe(time.Since(start).Seconds())
	return rows, err
}

// QueryRow wrapper to sql.DB.QueryRow
func (store *SQLStore) QueryRow(query string, args ...interface{}) *sql.Row {
	return store.DB.QueryRow(query, args...)
}

// EnableTx enable transaction
func (store *SQLStore) EnableTx() (Store, error) {
	ctx := context.Background()
	tx, err := store.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &SQLStoreTx{
		DB: tx,
	}, nil
}

// Commit does not commit, do nothing here
func (store *SQLStore) Commit() error {
	return nil
}

// Rollback doesn't rollback, do nothing here
func (store *SQLStore) Rollback() error {
	return nil
}

// TxEnd doesnt rollback, do nothing here
func (store *SQLStore) TxEnd(txFunc func() error) error {
	return nil
}

// SQLStoreTx concrete implementation of sqlGdbc by using *sql.Tx
type SQLStoreTx struct {
	DB *sql.Tx
}

// Exec wrapper to sql.DB.Exec
func (store *SQLStoreTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return store.DB.Exec(query, args...)
}

// Prepare wrapper to sql.DB.Prepare
func (store *SQLStoreTx) Prepare(query string) (*sql.Stmt, error) {
	return store.DB.Prepare(query)
}

// Query wrapper to sql.DB.Query
func (store *SQLStoreTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return store.DB.Query(query, args...)
}

// QueryRow wrapper to sql.DB.QueryRow
func (store *SQLStoreTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return store.DB.QueryRow(query, args...)
}

// Commit doesnt commit, do nothing here
func (store *SQLStoreTx) Commit() error {
	return store.DB.Commit()
}

// Rollback doesn't rollback, do nothing here
func (store *SQLStoreTx) Rollback() error {
	return store.DB.Rollback()
}

// TxEnd doesnt rollback, do nothing here
func (store *SQLStoreTx) TxEnd(function func() error) error {
	var err error
	tx := store.DB
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = function()
	return err
}
