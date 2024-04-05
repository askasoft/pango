package sqlx

import (
	"context"
	"database/sql"
	"fmt"
)

// Beginxer is an interface used by Transaction
type Beginxer interface {
	Beginx() (*Tx, error)
}

// BeginTxxer is an interface used by Transactionx
type BeginTxxer interface {
	BeginTxx(context.Context, *sql.TxOptions) (*Tx, error)
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func Transaction(db Beginxer, fc func(tx *Tx) error) (err error) {
	var tx *Tx

	tx, err = db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		// Make sure to rollback when panic
		if x := recover(); x != nil {
			err = fmt.Errorf("panic: %v", x)
			_ = tx.Rollback()
		}
	}()

	if err = fc(tx); err == nil {
		return tx.Commit()
	}

	return nil
}

// Transactionx start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func Transactionx(db BeginTxxer, ctx context.Context, opts *sql.TxOptions, fc func(tx *Tx) error) (err error) { //nolint: all
	var tx *Tx

	tx, err = db.BeginTxx(ctx, opts)
	if err != nil {
		return
	}

	defer func() {
		// Make sure to rollback when panic
		if x := recover(); x != nil {
			err = fmt.Errorf("panic: %v", x)
			_ = tx.Rollback()
		}
	}()

	if err = fc(tx); err == nil {
		return tx.Commit()
	}

	return nil
}
