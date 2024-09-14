package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/askasoft/pango/sqx"
)

// Stmt is an sqlx wrapper around sql.Stmt with extra functionality
type Stmt struct {
	query string
	stmt  *sql.Stmt
	ext
}

// Close closes the statement.
func (s *Stmt) Close() error {
	return s.stmt.Close()
}

// Exec executes a prepared statement with the given arguments and
// returns a Result summarizing the effect of the statement.
func (s *Stmt) Exec(args ...any) (sql.Result, error) {
	qs := &qStmt{s}
	return qs.Exec(s.query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (s *Stmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	qs := &qStmt{s}
	return qs.ExecContext(ctx, s.query, args...)
}

// Query executes a prepared query statement with the given arguments
// and returns the query results as a *Rows.
func (s *Stmt) Query(args ...any) (*sql.Rows, error) {
	qs := &qStmt{s}
	return qs.Query(s.query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (s *Stmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	qs := &qStmt{s}
	return qs.QueryContext(ctx, s.query, args...)
}

// QueryRowxContext using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) QueryRowxContext(ctx context.Context, args ...any) *Row {
	qs := &qStmt{s}
	return qs.QueryRowxContext(ctx, s.query, args...)
}

// Stmt returns the underlying *sql.Stmt
func (s *Stmt) Stmt() *sql.Stmt {
	return s.stmt
}

// Unsafe returns a version of Stmt which will silently succeed to scan when
// columns in the SQL result have no fields in the destination struct.
func (s *Stmt) Unsafe() *Stmt {
	c := &Stmt{query: s.query, stmt: s.stmt, ext: s.ext}
	c.unsafe = true
	return c
}

// Select using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) Select(dest any, args ...any) error {
	return Select(&qStmt{s}, dest, s.query, args...)
}

// Get using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (s *Stmt) Get(dest any, args ...any) error {
	return Get(&qStmt{s}, dest, s.query, args...)
}

// Queryx using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) Queryx(args ...any) (*Rows, error) {
	qs := &qStmt{s}
	return qs.Queryx(s.query, args...)
}

// QueryRowx using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) QueryRowx(args ...any) *Row {
	qs := &qStmt{s}
	return qs.QueryRowx(s.query, args...)
}

// SelectContext using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) SelectContext(ctx context.Context, dest any, args ...any) error {
	return SelectContext(ctx, &qStmt{s}, dest, s.query, args...)
}

// GetContext using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (s *Stmt) GetContext(ctx context.Context, dest any, args ...any) error {
	return GetContext(ctx, &qStmt{s}, dest, s.query, args...)
}

// MustExec (panic) using this statement.  Note that the query portion of the error
// output will be blank, as Stmt does not expose its query.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) MustExec(args ...any) sql.Result {
	return sqx.MustExec(&qStmt{s}, s.query, args...)
}

// MustExecContext (panic) using this statement.  Note that the query portion of
// the error output will be blank, as Stmt does not expose its query.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) MustExecContext(ctx context.Context, args ...any) sql.Result {
	return sqx.MustExecContext(ctx, &qStmt{s}, s.query, args...)
}

// QueryxContext using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) QueryxContext(ctx context.Context, args ...any) (*Rows, error) {
	qs := &qStmt{s}
	return qs.QueryxContext(ctx, s.query, args...)
}

// qStmt is an unexposed wrapper which lets you use a Stmt as a Queryer & Execer by
// implementing those interfaces and ignoring the `query` argument.
type qStmt struct {
	s *Stmt
}

func (q *qStmt) Exec(query string, args ...any) (sql.Result, error) {
	return q.s.tracer.TraceStmtExec(q.s.stmt, query, args...)
}

func (q *qStmt) Query(query string, args ...any) (*sql.Rows, error) {
	return q.s.tracer.TraceStmtQuery(q.s.stmt, query, args...)
}

func (q *qStmt) Queryx(query string, args ...any) (*Rows, error) {
	r, err := q.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: q.s.ext}, err
}

func (q *qStmt) QueryRowx(query string, args ...any) *Row {
	rows, err := q.Query(query, args...)
	return &Row{rows: rows, err: err, ext: q.s.ext}
}
func (q *qStmt) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return q.s.tracer.TraceStmtQueryContext(ctx, q.s.stmt, query, args...)
}

func (q *qStmt) QueryxContext(ctx context.Context, query string, args ...any) (*Rows, error) {
	r, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: q.s.ext}, err
}

func (q *qStmt) QueryRowxContext(ctx context.Context, query string, args ...any) *Row {
	rows, err := q.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, ext: q.s.ext}
}

func (q *qStmt) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return q.s.tracer.TraceStmtExecContext(ctx, q.s.stmt, query, args...)
}

func getQueryStmt(stmt any) (q string, s *sql.Stmt) {
	switch v := stmt.(type) {
	case Stmt:
		q = v.query
		s = v.stmt
	case *Stmt:
		q = v.query
		s = v.stmt
	case *sql.Stmt:
		f := reflect.ValueOf(stmt).Elem().FieldByName("query")
		if f.IsZero() {
			q = "sql.Stmt"
		} else {
			pc := (*string)(unsafe.Pointer(f.UnsafeAddr()))
			q = *pc
		}
		s = v
	default:
		panic(fmt.Sprintf("non-statement type %T", stmt))
	}

	return
}
