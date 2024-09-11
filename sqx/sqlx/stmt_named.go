package sqlx

import (
	"context"
	"database/sql"
)

// NamedStmt is a prepared statement that executes named queries.  Prepare it
// how you would execute a NamedQuery, but pass in a struct or map when executing.
type NamedStmt struct {
	stmt   *Stmt
	query  string
	params []string
}

// IsUnsafe return unsafe
func (ns *NamedStmt) IsUnsafe() bool {
	return ns.stmt.IsUnsafe()
}

// Unsafe creates an unsafe version of the NamedStmt
func (ns *NamedStmt) Unsafe() *NamedStmt {
	r := &NamedStmt{params: ns.params, stmt: ns.stmt, query: ns.query}
	r.stmt.unsafe = true
	return r
}

// Close closes the named statement.
func (ns *NamedStmt) Close() error {
	return ns.stmt.Close()
}

// Exec executes a named statement using the struct passed.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) Exec(arg any) (sql.Result, error) {
	args, err := bindAnyArgs(ns.params, arg, ns.stmt.mapper)
	if err != nil {
		return *new(sql.Result), err
	}
	return ns.stmt.Exec(args...)
}

// Query executes a named statement using the struct argument, returning rows.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) Query(arg any) (*sql.Rows, error) {
	args, err := bindAnyArgs(ns.params, arg, ns.stmt.mapper)
	if err != nil {
		return nil, err
	}
	return ns.stmt.Query(args...)
}

// QueryRow executes a named statement against the database.  Because sqlx cannot
// create a *sql.Row with an error condition pre-set for binding errors, sqlx
// returns a *sqlx.Row instead.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) QueryRow(arg any) *Row {
	args, err := bindAnyArgs(ns.params, arg, ns.stmt.mapper)
	if err != nil {
		return &Row{err: err}
	}
	return ns.stmt.QueryRowx(args...)
}

// MustExec execs a NamedStmt, panicing on error
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) MustExec(arg any) sql.Result {
	res, err := ns.Exec(arg)
	if err != nil {
		panic(err)
	}
	return res
}

// Queryx using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) Queryx(arg any) (*Rows, error) {
	r, err := ns.Query(arg)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: ns.stmt.ext}, err
}

// QueryRowx this NamedStmt.  Because of limitations with QueryRow, this is
// an alias for QueryRow.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) QueryRowx(arg any) *Row {
	return ns.QueryRow(arg)
}

// Select using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) Select(dest any, arg any) error {
	rows, err := ns.Queryx(arg)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// Get using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) Get(dest any, arg any) error {
	r := ns.QueryRowx(arg)
	return r.scanAny(dest, false)
}

// ExecContext executes a named statement using the struct passed.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) ExecContext(ctx context.Context, arg any) (sql.Result, error) {
	args, err := bindAnyArgs(ns.params, arg, ns.stmt.mapper)
	if err != nil {
		return *new(sql.Result), err
	}
	return ns.stmt.ExecContext(ctx, args...)
}

// QueryContext executes a named statement using the struct argument, returning rows.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) QueryContext(ctx context.Context, arg any) (*sql.Rows, error) {
	args, err := bindAnyArgs(ns.params, arg, ns.stmt.mapper)
	if err != nil {
		return nil, err
	}
	return ns.stmt.QueryContext(ctx, args...)
}

// QueryRowContext executes a named statement against the database.  Because sqlx cannot
// create a *sql.Row with an error condition pre-set for binding errors, sqlx
// returns a *sqlx.Row instead.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) QueryRowContext(ctx context.Context, arg any) *Row {
	args, err := bindAnyArgs(ns.params, arg, ns.stmt.mapper)
	if err != nil {
		return &Row{err: err}
	}
	return ns.stmt.QueryRowxContext(ctx, args...)
}

// MustExecContext execs a NamedStmt, panicing on error
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) MustExecContext(ctx context.Context, arg any) sql.Result {
	res, err := ns.ExecContext(ctx, arg)
	if err != nil {
		panic(err)
	}
	return res
}

// QueryxContext using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) QueryxContext(ctx context.Context, arg any) (*Rows, error) {
	r, err := ns.QueryContext(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: ns.stmt.ext}, err
}

// QueryRowxContext this NamedStmt.  Because of limitations with QueryRow, this is
// an alias for QueryRow.
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) QueryRowxContext(ctx context.Context, arg any) *Row {
	return ns.QueryRowContext(ctx, arg)
}

// SelectContext using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) SelectContext(ctx context.Context, dest any, arg any) error {
	rows, err := ns.QueryxContext(ctx, arg)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// GetContext using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (ns *NamedStmt) GetContext(ctx context.Context, dest any, arg any) error {
	r := ns.QueryRowxContext(ctx, arg)
	return r.scanAny(dest, false)
}

// A union interface of Preparerx and binder, required to be able to
// prepare named statements (as the bindtype must be determined).
type iPreparerx interface {
	Preparerx
	binder
}

func prepareNamed(p iPreparerx, query string) (*NamedStmt, error) {
	q, args, err := p.Binder().compileNamedQuery(query)
	if err != nil {
		return nil, err
	}

	stmt, err := p.Preparex(q)
	if err != nil {
		return nil, err
	}

	return &NamedStmt{
		stmt:   stmt,
		query:  q,
		params: args,
	}, nil
}

// A union interface of ContextPreparerx and binder, required to be able to
// prepare named statements with context (as the bindtype must be determined).
type iContextNamedPreparerx interface {
	ContextPreparerx
	binder
}

func prepareNamedContext(ctx context.Context, p iContextNamedPreparerx, query string) (*NamedStmt, error) {
	q, args, err := p.Binder().compileNamedQuery(query)
	if err != nil {
		return nil, err
	}
	stmt, err := p.PreparexContext(ctx, q)
	if err != nil {
		return nil, err
	}
	return &NamedStmt{
		stmt:   stmt,
		query:  q,
		params: args,
	}, nil
}
