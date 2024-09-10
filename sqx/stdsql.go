package sqx

import (
	"context"
	"database/sql"
	"fmt"
)

//------------------------------------------------
// GO database/sql interface
//

// Pinger is an interface for Ping()
type Pinger interface {
	Ping() error
}

// Queryer is an interface used by Get and Select
type Queryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type RowQueryer interface {
	QueryRow(query string, args ...any) *sql.Row
}

// Execer is an interface used by MustExec
type Execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

// StmtQueryer is an interface used by Get and Select
type StmtQueryer interface {
	Query(args ...any) (*sql.Rows, error)
}

// StmtExecer is an interface used by MustExec
type StmtExecer interface {
	Exec(args ...any) (sql.Result, error)
}

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

type Commiter interface {
	Commit() error
}

type Rollbacker interface {
	Rollback() error
}

// Sql the basic interface for sql.DB, sql.Tx
type Sql interface {
	Queryer
	Execer
	Preparer
}

// Beginer is an interface for Begin()
type Beginer interface {
	Begin() (*sql.Tx, error)
}

// BeginTxer is an interface used by Transaction
type BeginTxer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// MustExec execs the query using e and panics if there was an error.
// Any placeholder parameters are replaced with supplied args.
func MustExec(e Execer, query string, args ...any) sql.Result {
	res, err := e.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func Transaction(db Beginer, fc func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx

	tx, err = db.Begin()
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

	return
}

// Transactionx start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func Transactionx(ctx context.Context, db BeginTxer, opts *sql.TxOptions, fc func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx

	tx, err = db.BeginTx(ctx, opts)
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

	return
}
