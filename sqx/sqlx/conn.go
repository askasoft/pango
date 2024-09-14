package sqlx

import (
	"context"
	"database/sql"
)

// Conn is a wrapper around sql.Conn with extra functionality
type Conn struct {
	conn *sql.Conn
	ext
}

// BeginTx starts a transaction.
//
// The provided context is used until the transaction is committed or rolled back.
// If the context is canceled, the sql package will roll back
// the transaction. Tx.Commit will return an error if the context provided to
// BeginTx is canceled.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
func (c *Conn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.tracer.TraceBeginTx(ctx, c.conn, opts)
}

// Close returns the connection to the connection pool.
// All operations after a Close will return with ErrConnDone.
// Close is safe to call concurrently with other operations and will
// block until all other operations finish. It may be useful to first
// cancel any used context and then call close directly after.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (c *Conn) Exec(query string, args ...any) (sql.Result, error) {
	return c.ExecContext(context.Background(), query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (c *Conn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.tracer.TraceExecContext(ctx, c.conn, query, args...)
}

// Ping verifies the connection to the database is still alive.
func (c *Conn) Ping() error {
	return c.PingContext(context.Background())
}

// PingContext verifies the connection to the database is still alive.
func (c *Conn) PingContext(ctx context.Context) error {
	return c.tracer.TracePingContext(ctx, c.conn)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The provided context is used for the preparation of the statement, not for the
// execution of the statement.
func (c *Conn) Prepare(query string) (*sql.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The provided context is used for the preparation of the statement, not for the
// execution of the statement.
func (c *Conn) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.tracer.TracePrepareContext(ctx, c.conn, query)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (c *Conn) Query(query string, args ...any) (*sql.Rows, error) {
	return c.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (c *Conn) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.tracer.TraceQueryContext(ctx, c.conn, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (c *Conn) QueryRow(query string, args ...any) *sql.Row {
	return c.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return c.tracer.TraceQueryRowContext(ctx, c.conn, query, args...)
}

// Raw executes f exposing the underlying driver connection for the
// duration of f. The driverConn must not be used outside of f.
//
// Once f returns and err is not driver.ErrBadConn, the Conn will continue to be usable
// until Conn.Close is called.
func (c *Conn) Raw(f func(driverConn any) error) error {
	return c.conn.Raw(f)
}

// Conn returns the underlying *sql.Conn
func (c *Conn) Conn() *sql.Conn {
	return c.conn
}

// BindNamed binds a query using the DB driver's bindvar type.
func (c *Conn) BindNamed(query string, arg any) (string, []any, error) {
	return bindNamedMapper(c.binder, query, arg, c.mapper)
}

// Beginx begins a transaction and returns an *sqlx.Tx instead of an *sql.Tx.
func (c *Conn) Beginx() (*Tx, error) {
	tx, err := c.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx, ext: c.ext}, err
}

// BeginTxx begins a transaction and returns an *sqlx.Tx instead of an
// *sql.Tx.
//
// The provided context is used until the transaction is committed or rolled
// back. If the context is canceled, the sql package will roll back the
// transaction. Tx.Commit will return an error if the context provided to
// BeginxContext is canceled.
func (c *Conn) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := c.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx, ext: c.ext}, err
}

// Preparex returns an sqlx.Stmt instead of a sql.Stmt.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (c *Conn) Preparex(query string) (*Stmt, error) {
	s, err := c.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{query: query, stmt: s, ext: c.ext}, err
}

// PreparexContext returns an sqlx.Stmt instead of a sql.Stmt.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (c *Conn) PreparexContext(ctx context.Context, query string) (*Stmt, error) {
	s, err := c.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{query: query, stmt: s, ext: c.ext}, err
}

// Queryx queries the database and returns an *sqlx.Rows.
// Any placeholder parameters are replaced with supplied args.
func (c *Conn) Queryx(query string, args ...any) (*Rows, error) {
	r, err := c.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: c.ext}, err
}

// QueryxContext queries the database and returns an *sqlx.Rows.
// Any placeholder parameters are replaced with supplied args.
func (c *Conn) QueryxContext(ctx context.Context, query string, args ...any) (*Rows, error) {
	r, err := c.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: c.ext}, err
}

// QueryRowx queries the database and returns an *sqlx.Row.
// Any placeholder parameters are replaced with supplied args.
func (c *Conn) QueryRowx(query string, args ...any) *Row {
	rows, err := c.Query(query, args...)
	return &Row{rows: rows, err: err, ext: c.ext}
}

// QueryRowxContext queries the database and returns an *sqlx.Row.
// Any placeholder parameters are replaced with supplied args.
func (c *Conn) QueryRowxContext(ctx context.Context, query string, args ...any) *Row {
	rows, err := c.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, ext: c.ext}
}

// NamedExec using this Conn.
// Any named placeholder parameters are replaced with fields from arg.
func (c *Conn) NamedExec(query string, arg any) (sql.Result, error) {
	return namedExec(c, query, arg)
}

// NamedExecContext using this Conn.
// Any named placeholder parameters are replaced with fields from arg.
func (c *Conn) NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	return namedExecContext(ctx, c, query, arg)
}

// NamedQuery using this Conn.
// Any named placeholder parameters are replaced with fields from arg.
func (c *Conn) NamedQuery(query string, arg any) (*Rows, error) {
	return namedQuery(c, query, arg)
}

// NamedQueryContext using this Conn.
// Any named placeholder parameters are replaced with fields from arg.
func (c *Conn) NamedQueryContext(ctx context.Context, query string, arg any) (*Rows, error) {
	return namedQueryContext(ctx, c, query, arg)
}

// NamedQueryRow using this Conn.
// Any named placeholder parameters are replaced with fields from arg.
func (c *Conn) NamedQueryRow(query string, arg any) *Row {
	return namedQueryRow(c, query, arg)
}

// NamedQueryRowContext using the Conn.
// Any named placeholder parameters are replaced with fields from arg.
func (c *Conn) NamedQueryRowContext(ctx context.Context, query string, arg any) *Row {
	return namedQueryRowContext(ctx, c, query, arg)
}

// Select executes a query, and StructScans each row
// into dest, which must be a slice.  If the slice elements are scannable, then
// the result set must have only one column.  Otherwise, StructScan is used.
// The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func (c *Conn) Select(dest any, query string, args ...any) error {
	return Select(c, dest, query, args...)
}

// NamedSelect executes a query, and StructScans each row
// into dest, which must be a slice.  If the slice elements are scannable, then
// the result set must have only one column.  Otherwise, StructScan is used.
// The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func (c *Conn) NamedSelect(dest any, query string, arg any) error {
	return NamedSelect(c, dest, query, arg)
}

// SelectContext executes a query, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied arg.
func (c *Conn) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return SelectContext(ctx, c, dest, query, args...)
}

// NamedSelectContext executes a query, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied arg.
func (c *Conn) NamedSelectContext(ctx context.Context, dest any, query string, arg any) error {
	return NamedSelectContext(ctx, c, dest, query, arg)
}

// Get does a QueryRowx() and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (c *Conn) Get(dest any, query string, args ...any) error {
	return Get(c, dest, query, args...)
}

// NamedGet does a NamedQueryRow() and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (c *Conn) NamedGet(dest any, query string, arg any) error {
	return NamedGet(c, dest, query, arg)
}

// GetContext does a QueryRowxContext() and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (c *Conn) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return GetContext(ctx, c, dest, query, args...)
}

// NamedGetContext does a NamedQueryRowContext() and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (c *Conn) NamedGetContext(ctx context.Context, dest any, query string, arg any) error {
	return NamedGetContext(ctx, c, dest, query, arg)
}

// Create does a QueryRowx() and scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() and returns Result.LastInsertId().
func (c *Conn) Create(query string, args ...any) (int64, error) {
	return Create(c, query, args...)
}

// NamedCreate does a NamedQueryRowx() and scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() return Result.LastInsertId().
func (c *Conn) NamedCreate(query string, arg any) (int64, error) {
	return NamedCreate(c, query, arg)
}

// CreateContext does a QueryRowxContext() scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() return Result.LastInsertId().
func (c *Conn) CreateContext(ctx context.Context, query string, args ...any) (int64, error) {
	return CreateContext(ctx, c, query, args...)
}

// NamedCreateContext does a NamedQueryRow() and scans the resulting row
// returns the last inserted ID.
// If the db supports LastInsertId(), does a NamedExecContext() and returns Result.LastInsertId().
func (c *Conn) NamedCreateContext(ctx context.Context, query string, arg any) (int64, error) {
	return NamedCreateContext(ctx, c, query, arg)
}

// Transactionx start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func (c *Conn) Transactionx(ctx context.Context, opts *sql.TxOptions, fc func(tx *Tx) error) (err error) {
	return Transactionx(ctx, c, opts, fc)
}
