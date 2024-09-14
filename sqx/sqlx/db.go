package sqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/sqx"
)

// DB is a wrapper around sql.DB which keeps track of the driverName upon Open,
// used mostly to automatically bind named queries using the right bindvars.
type DB struct {
	db *sql.DB
	ext
}

// Begin starts a transaction. The default isolation level is dependent on
// the driver.
func (db *DB) Begin() (*sql.Tx, error) {
	return db.tracer.TraceBegin(db.db)
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
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.tracer.TraceBeginTx(ctx, db.db, opts)
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (db *DB) Close() error {
	return db.db.Close()
}

// Conn returns a single connection by either opening a new connection
// or returning an existing connection from the connection pool. Conn will
// block until either a connection is returned or ctx is canceled.
// Queries run on the same Conn will be run in the same database session.
//
// Every Conn must be returned to the database pool after use by
// calling Conn.Close.
func (db *DB) Conn(ctx context.Context) (*sql.Conn, error) {
	return db.db.Conn(ctx)
}

// Driver returns the database's underlying driver.
func (db *DB) Driver() driver.Driver {
	return db.db.Driver()
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (db *DB) Exec(query string, args ...any) (sql.Result, error) {
	return db.tracer.TraceExec(db.db, query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.tracer.TraceExecContext(ctx, db.db, query, args...)
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (db *DB) Ping() error {
	return db.tracer.TracePing(db.db)
}

// PingContext verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (db *DB) PingContext(ctx context.Context) error {
	return db.tracer.TracePingContext(ctx, db.db)
}

func (db *DB) Prepare(query string) (*sql.Stmt, error) {
	return db.tracer.TracePrepare(db.db, query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
func (db *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.db.PrepareContext(ctx, query)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) Query(query string, args ...any) (*sql.Rows, error) {
	return db.tracer.TraceQuery(db.db, query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.tracer.TraceQueryContext(ctx, db.db, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (db *DB) QueryRow(query string, args ...any) *sql.Row {
	return db.tracer.TraceQueryRow(db.db, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.tracer.TraceQueryRowContext(ctx, db.db, query, args...)
}

// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are not closed due to a connection's idle time.
func (db *DB) SetConnMaxIdleTime(d time.Duration) {
	db.db.SetConnMaxIdleTime(d)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are not closed due to a connection's age.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	db.db.SetConnMaxLifetime(d)
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
//
// If n <= 0, no idle connections are retained.
//
// The default max idle connections is currently 2. This may change in
// a future release.
func (db *DB) SetMaxIdleConns(n int) {
	db.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (db *DB) SetMaxOpenConns(n int) {
	db.db.SetMaxOpenConns(n)
}

// Stats returns database statistics.
func (db *DB) Stats() sql.DBStats {
	return db.db.Stats()
}

// DB returns the underlying *sql.DB
func (db *DB) DB() *sql.DB {
	return db.db
}

// MapperFunc sets a new mapper for this db using the default sqlx struct tag
// and the provided mapper function.
func (db *DB) MapperFunc(mf func(string) string) {
	db.mapper = ref.NewMapperFunc("db", mf)
}

// Unsafe returns a version of DB which will silently succeed to scan when
// columns in the SQL result have no fields in the destination struct.
// sqlx.Stmt and sqlx.Tx which are created from this DB will inherit its
// safety behavior.
func (db *DB) Unsafe() *DB {
	ndb := &DB{db: db.db, ext: db.ext}
	ndb.unsafe = true
	return ndb
}

// BindNamed binds a query using the DB driver's bindvar type.
func (db *DB) BindNamed(query string, arg any) (string, []any, error) {
	return bindNamedMapper(db.binder, query, arg, db.mapper)
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedExec(query string, arg any) (sql.Result, error) {
	return namedExec(db, query, arg)
}

// NamedExecContext using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	return namedExecContext(ctx, db, query, arg)
}

// NamedQuery using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedQuery(query string, arg any) (*Rows, error) {
	return namedQuery(db, query, arg)
}

// NamedQueryContext using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedQueryContext(ctx context.Context, query string, arg any) (*Rows, error) {
	return namedQueryContext(ctx, db, query, arg)
}

// NamedQueryRow using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedQueryRow(query string, arg any) *Row {
	return namedQueryRow(db, query, arg)
}

// NamedQueryRowContext using the DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedQueryRowContext(ctx context.Context, query string, arg any) *Row {
	return namedQueryRowContext(ctx, db, query, arg)
}

// Queryx queries the database and returns an *sqlx.Rows.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) Queryx(query string, args ...any) (*Rows, error) {
	r, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: db.ext}, err
}

// QueryxContext queries the database and returns an *sqlx.Rows.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) QueryxContext(ctx context.Context, query string, args ...any) (*Rows, error) {
	r, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: db.ext}, err
}

// QueryRowx queries the database and returns an *sqlx.Row.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) QueryRowx(query string, args ...any) *Row {
	rows, err := db.Query(query, args...)
	return &Row{rows: rows, err: err, ext: db.ext}
}

// QueryRowxContext queries the database and returns an *sqlx.Row.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) QueryRowxContext(ctx context.Context, query string, args ...any) *Row {
	rows, err := db.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, ext: db.ext}
}

// Preparex returns an sqlx.Stmt instead of a sql.Stmt
func (db *DB) Preparex(query string) (*Stmt, error) {
	s, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{query: query, stmt: s, ext: db.ext}, err
}

// PreparexContext returns an sqlx.Stmt instead of a sql.Stmt.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (db *DB) PreparexContext(ctx context.Context, query string) (*Stmt, error) {
	s, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{query: query, stmt: s, ext: db.ext}, err
}

// PrepareNamed returns an sqlx.NamedStmt
func (db *DB) PrepareNamed(query string) (*NamedStmt, error) {
	return prepareNamed(db, query)
}

// PrepareNamedContext returns an sqlx.NamedStmt
func (db *DB) PrepareNamedContext(ctx context.Context, query string) (*NamedStmt, error) {
	return prepareNamedContext(ctx, db, query)
}

// Select executes a query, and StructScans each row
// into dest, which must be a slice.  If the slice elements are scannable, then
// the result set must have only one column.  Otherwise, StructScan is used.
// The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) Select(dest any, query string, args ...any) error {
	return Select(db, dest, query, args...)
}

// NamedSelect executes a query, and StructScans each row
// into dest, which must be a slice.  If the slice elements are scannable, then
// the result set must have only one column.  Otherwise, StructScan is used.
// The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) NamedSelect(dest any, query string, arg any) error {
	return NamedSelect(db, dest, query, arg)
}

// SelectContext executes a query, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied arg.
func (db *DB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return SelectContext(ctx, db, dest, query, args...)
}

// NamedSelectContext executes a query, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied arg.
func (db *DB) NamedSelectContext(ctx context.Context, dest any, query string, arg any) error {
	return NamedSelectContext(ctx, db, dest, query, arg)
}

// Get does a QueryRowx() and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) Get(dest any, query string, args ...any) error {
	return Get(db, dest, query, args...)
}

// NamedGet does a NamedQueryRow() and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) NamedGet(dest any, query string, arg any) error {
	return NamedGet(db, dest, query, arg)
}

// GetContext does a QueryRowxContext() and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return GetContext(ctx, db, dest, query, args...)
}

// NamedGetContext does a NamedQueryRowContext() and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) NamedGetContext(ctx context.Context, dest any, query string, arg any) error {
	return NamedGetContext(ctx, db, dest, query, arg)
}

// Create does a QueryRowx() and scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() and returns Result.LastInsertId().
func (db *DB) Create(query string, args ...any) (int64, error) {
	return Create(db, query, args...)
}

// NamedCreate does a NamedQueryRowx() and scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() return Result.LastInsertId().
func (db *DB) NamedCreate(query string, arg any) (int64, error) {
	return NamedCreate(db, query, arg)
}

// CreateContext does a QueryRowxContext() scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() return Result.LastInsertId().
func (db *DB) CreateContext(ctx context.Context, query string, args ...any) (int64, error) {
	return CreateContext(ctx, db, query, args...)
}

// NamedCreateContext does a NamedQueryRow() and scans the resulting row
// returns the last inserted ID.
// If the db supports LastInsertId(), does a NamedExecContext() and returns Result.LastInsertId().
func (db *DB) NamedCreateContext(ctx context.Context, query string, arg any) (int64, error) {
	return NamedCreateContext(ctx, db, query, arg)
}

// Beginx begins a transaction and returns an *sqlx.Tx instead of an *sql.Tx.
func (db *DB) Beginx() (*Tx, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx, ext: db.ext}, err
}

// BeginTxx begins a transaction and returns an *sqlx.Tx instead of an
// *sql.Tx.
//
// The provided context is used until the transaction is committed or rolled
// back. If the context is canceled, the sql package will roll back the
// transaction. Tx.Commit will return an error if the context provided to
// BeginxContext is canceled.
func (db *DB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx, ext: db.ext}, err
}

// Connx returns an *sqlx.Conn instead of an *sql.Conn.
func (db *DB) Connx(ctx context.Context) (*Conn, error) {
	conn, err := db.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	return &Conn{conn: conn, ext: db.ext}, nil
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func (db *DB) Transaction(fc func(tx *Tx) error) (err error) {
	return Transaction(db, fc)
}

// Transactionx start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func (db *DB) Transactionx(ctx context.Context, opts *sql.TxOptions, fc func(tx *Tx) error) (err error) {
	return Transactionx(ctx, db, opts, fc)
}

// MustBeginx starts a transaction, and panics on error.  Returns an *sqlx.Tx instead
// of an *sql.Tx.
func (db *DB) MustBeginx() *Tx {
	tx, err := db.Beginx()
	if err != nil {
		panic(err)
	}
	return tx
}

// MustBeginTx starts a transaction, and panics on error.  Returns an *sqlx.Tx instead
// of an *sql.Tx.
//
// The provided context is used until the transaction is committed or rolled
// back. If the context is canceled, the sql package will roll back the
// transaction. Tx.Commit will return an error if the context provided to
// MustBeginContext is canceled.
func (db *DB) MustBeginTx(ctx context.Context, opts *sql.TxOptions) *Tx {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		panic(err)
	}
	return tx
}

// MustExec (panic) runs MustExec using this database.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) MustExec(query string, args ...any) sql.Result {
	return sqx.MustExec(db, query, args...)
}

// MustExecContext (panic) runs MustExec using this database.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) MustExecContext(ctx context.Context, query string, args ...any) sql.Result {
	return sqx.MustExecContext(ctx, db, query, args...)
}
