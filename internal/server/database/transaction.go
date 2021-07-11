package database

import (
	"database/sql"
)

// Transaction is an interface that models the standard transaction in `database/sql`.
type Transaction interface {
	// Commit a transaction
	Commit() error
	// Rollback a transaction
	Rollback() error
	// TxEnd commits a transaction if no errors, otherwise rollback. TxFunc is the operations wrapped in a transaction
	TxEnd(TxFunc func() error) error
}

// TxFunc is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFunc func() error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFunc`
func WithTransaction(db *sql.DB, transaction TxFunc) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
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
	err = transaction()
	return err
}
