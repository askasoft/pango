package sqlx

import (
	"context"
	"database/sql"

	"github.com/askasoft/pango/sqx"
)

// Tx is an sqlx wrapper around sql.Tx with extra functionality
type Tx struct {
	tx *sql.Tx
	ext
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	return tx.tracer.TraceCommit(tx.tx)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (tx *Tx) Exec(query string, args ...any) (sql.Result, error) {
	return tx.tracer.TraceExec(tx.tx, query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.tracer.TraceExecContext(ctx, tx.tx, query, args...)
}

// Prepare creates a prepared statement for use within a transaction.
//
// The returned statement operates within the transaction and will be closed
// when the transaction has been committed or rolled back.
//
// To use an existing prepared statement on this transaction, see Tx.Stmt.
func (tx *Tx) Prepare(query string) (*sql.Stmt, error) {
	return tx.tracer.TracePrepare(tx, query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
func (tx *Tx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return tx.tx.PrepareContext(ctx, query)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (tx *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	return tx.tracer.TraceQuery(tx.tx, query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tracer.TraceQueryContext(ctx, tx.tx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (tx *Tx) QueryRow(query string, args ...any) *sql.Row {
	return tx.tracer.TraceQueryRow(tx.tx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return tx.tracer.TraceQueryRowContext(ctx, tx.tx, query, args)
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return tx.tracer.TraceRollback(tx.tx)
}

// Stmt returns a transaction-specific prepared statement from
// an existing statement.
func (tx *Tx) Stmt(stmt *sql.Stmt) *sql.Stmt {
	return tx.tx.Stmt(stmt)
}

// StmtContext returns a transaction-specific prepared statement from
// an existing statement.
func (tx *Tx) StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt {
	return tx.tx.StmtContext(ctx, stmt)
}

// Tx returns the underlying *sql.Tx
func (tx *Tx) Tx() *sql.Tx {
	return tx.tx
}

// Unsafe returns a version of Tx which will silently succeed to scan when
// columns in the SQL result have no fields in the destination struct.
func (tx *Tx) Unsafe() *Tx {
	ntx := &Tx{tx: tx.tx, ext: tx.ext}
	ntx.unsafe = true
	return ntx
}

// BindNamed binds a query within a transaction's bindvar type.
func (tx *Tx) BindNamed(query string, arg any) (string, []any, error) {
	return bindNamedMapper(tx.binder, query, arg, tx.mapper)
}

// NamedQuery within a transaction.
// Any named placeholder parameters are replaced with fields from arg.
func (tx *Tx) NamedQuery(query string, arg any) (*Rows, error) {
	return namedQuery(tx, query, arg)
}

// NamedQueryContext using this Tx.
// Any named placeholder parameters are replaced with fields from arg.
func (tx *Tx) NamedQueryContext(ctx context.Context, query string, arg any) (*Rows, error) {
	return namedQueryContext(ctx, tx, query, arg)
}

// NamedQueryRowContext within a transaction.
// Any named placeholder parameters are replaced with fields from arg.
func (tx *Tx) NamedQueryRowContext(ctx context.Context, query string, arg any) *Row {
	return namedQueryRowContext(ctx, tx, query, arg)
}

// NamedQueryRow within a transaction.
// Any named placeholder parameters are replaced with fields from arg.
func (tx *Tx) NamedQueryRow(query string, arg any) *Row {
	return namedQueryRow(tx, query, arg)
}

// NamedExec a named query within a transaction.
// Any named placeholder parameters are replaced with fields from arg.
func (tx *Tx) NamedExec(query string, arg any) (sql.Result, error) {
	return namedExec(tx, query, arg)
}

// NamedExecContext using this Tx.
// Any named placeholder parameters are replaced with fields from arg.
func (tx *Tx) NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	return namedExecContext(ctx, tx, query, arg)
}

// Queryx within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) Queryx(query string, args ...any) (*Rows, error) {
	r, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: tx.ext}, err
}

// QueryxContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryxContext(ctx context.Context, query string, args ...any) (*Rows, error) {
	r, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: tx.ext}, err
}

// QueryRowx within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryRowx(query string, args ...any) *Row {
	rows, err := tx.Query(query, args...)
	return &Row{rows: rows, err: err, ext: tx.ext}
}

// QueryRowxContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryRowxContext(ctx context.Context, query string, args ...any) *Row {
	rows, err := tx.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, ext: tx.ext}
}

// Preparex  a statement within a transaction.
func (tx *Tx) Preparex(query string) (*Stmt, error) {
	s, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{query: query, stmt: s, ext: tx.ext}, err
}

// PreparexContext returns an sqlx.Stmt instead of a sql.Stmt.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (tx *Tx) PreparexContext(ctx context.Context, query string) (*Stmt, error) {
	s, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{query: query, stmt: s, ext: tx.ext}, err
}

// Stmtx returns a version of the prepared statement which runs within a transaction.  Provided
// stmt can be either *sql.Stmt or *sqlx.Stmt.
func (tx *Tx) Stmtx(stmt any) *Stmt {
	q, s := getQueryStmt(stmt)
	return &Stmt{query: q, stmt: tx.Stmt(s), ext: tx.ext}
}

// StmtxContext returns a version of the prepared statement which runs within a
// transaction. Provided stmt can be either *sql.Stmt or *sqlx.Stmt.
func (tx *Tx) StmtxContext(ctx context.Context, stmt any) *Stmt {
	q, s := getQueryStmt(stmt)
	return &Stmt{query: q, stmt: tx.StmtContext(ctx, s), ext: tx.ext}
}

// NamedStmt returns a version of the prepared statement which runs within a transaction.
func (tx *Tx) NamedStmt(stmt *NamedStmt) *NamedStmt {
	return &NamedStmt{
		stmt:   tx.Stmtx(stmt.stmt),
		query:  stmt.query,
		params: stmt.params,
	}
}

// NamedStmtContext returns a version of the prepared statement which runs
// within a transaction.
func (tx *Tx) NamedStmtContext(ctx context.Context, stmt *NamedStmt) *NamedStmt {
	return &NamedStmt{
		stmt:   tx.StmtxContext(ctx, stmt.stmt),
		query:  stmt.query,
		params: stmt.params,
	}
}

// PrepareNamed returns an sqlx.NamedStmt
func (tx *Tx) PrepareNamed(query string) (*NamedStmt, error) {
	return prepareNamed(tx, query)
}

// PrepareNamedContext returns an sqlx.NamedStmt
func (tx *Tx) PrepareNamedContext(ctx context.Context, query string) (*NamedStmt, error) {
	return prepareNamedContext(ctx, tx, query)
}

// Select within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) Select(dest any, query string, args ...any) error {
	return Select(tx, dest, query, args...)
}

// SelectContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return SelectContext(ctx, tx, dest, query, args...)
}

// Get within a transaction.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (tx *Tx) Get(dest any, query string, args ...any) error {
	return Get(tx, dest, query, args...)
}

// GetContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (tx *Tx) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return GetContext(ctx, tx, dest, query, args...)
}

// Create does a QueryRow using the provided Queryer, and scans the resulting row
// returns the last inserted ID.
// If the db supports LastInsertId(), return Result.LastInsertId().
func (tx *Tx) Create(query string, args ...any) (int64, error) {
	return Create(tx, query, args...)
}

// Create does a QueryRow using the provided Queryer, and scans the resulting row
// returns the last inserted ID.
// If the db supports LastInsertId(), return Result.LastInsertId().
func (tx *Tx) CreateContext(ctx context.Context, query string, args ...any) (int64, error) {
	return CreateContext(ctx, tx, query, args...)
}

// MustExec runs MustExec within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) MustExec(query string, args ...any) sql.Result {
	return sqx.MustExec(tx, query, args...)
}

// MustExecContext runs MustExecContext within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) MustExecContext(ctx context.Context, query string, args ...any) sql.Result {
	return sqx.MustExecContext(ctx, tx, query, args...)
}
