//go:build go1.18
// +build go1.18

package sqx

import (
	"context"
	"database/sql"
)

// A union interface of contextPreparer and binder, required to be able to
// prepare named statements with context (as the bindtype must be determined).
type namedPreparerContext interface {
	PreparerxContext
	binder
}

func prepareNamedContext(ctx context.Context, p namedPreparerContext, query string) (*NamedStmt, error) {
	q, args, err := p.Binder().compileNamedQuery(query)
	if err != nil {
		return nil, err
	}
	stmt, err := p.PreparexContext(ctx, q)
	if err != nil {
		return nil, err
	}
	return &NamedStmt{
		QueryString: q,
		Params:      args,
		Stmt:        stmt,
	}, nil
}

// ExecContext executes a named statement using the struct passed.
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) ExecContext(ctx context.Context, arg any) (sql.Result, error) {
	args, err := bindAnyArgs(n.Params, arg, n.Stmt.mapper)
	if err != nil {
		return *new(sql.Result), err
	}
	return n.Stmt.ExecContext(ctx, args...)
}

// QueryContext executes a named statement using the struct argument, returning rows.
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) QueryContext(ctx context.Context, arg any) (*sql.Rows, error) {
	args, err := bindAnyArgs(n.Params, arg, n.Stmt.mapper)
	if err != nil {
		return nil, err
	}
	return n.Stmt.QueryContext(ctx, args...)
}

// QueryRowContext executes a named statement against the database.  Because sqx cannot
// create a *sql.Row with an error condition pre-set for binding errors, sqx
// returns a *sqx.Row instead.
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) QueryRowContext(ctx context.Context, arg any) *Row {
	args, err := bindAnyArgs(n.Params, arg, n.Stmt.mapper)
	if err != nil {
		return &Row{err: err}
	}
	return n.Stmt.QueryRowxContext(ctx, args...)
}

// MustExecContext execs a NamedStmt, panicing on error
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) MustExecContext(ctx context.Context, arg any) sql.Result {
	res, err := n.ExecContext(ctx, arg)
	if err != nil {
		panic(err)
	}
	return res
}

// QueryxContext using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) QueryxContext(ctx context.Context, arg any) (*Rows, error) {
	r, err := n.QueryContext(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: n.Stmt.ext}, err
}

// QueryRowxContext this NamedStmt.  Because of limitations with QueryRow, this is
// an alias for QueryRow.
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) QueryRowxContext(ctx context.Context, arg any) *Row {
	return n.QueryRowContext(ctx, arg)
}

// SelectContext using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) SelectContext(ctx context.Context, dest any, arg any) error {
	rows, err := n.QueryxContext(ctx, arg)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// GetContext using this NamedStmt
// Any named placeholder parameters are replaced with fields from arg.
func (n *NamedStmt) GetContext(ctx context.Context, dest any, arg any) error {
	r := n.QueryRowxContext(ctx, arg)
	return r.scanAny(dest, false)
}

// NamedQueryContext binds a named query and then runs Query on the result using the
// provided Ext (sqx.Tx, sqx.Db).  It works with both structs and with
// map[string]any types.
func NamedQueryContext(ctx context.Context, e ExtContext, query string, arg any) (*Rows, error) {
	q, args, err := e.Binder().bindNamedMapper(query, arg, e.Mapper())
	if err != nil {
		return nil, err
	}
	return e.QueryxContext(ctx, q, args...)
}

// NamedExecContext uses BindStruct to get a query executable by the driver and
// then runs Exec on the result.  Returns an error from the binding
// or the query execution itself.
func NamedExecContext(ctx context.Context, e ExtContext, query string, arg any) (sql.Result, error) {
	q, args, err := e.Binder().bindNamedMapper(query, arg, e.Mapper())
	if err != nil {
		return nil, err
	}
	return e.ExecContext(ctx, q, args...)
}
