package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/str"
)

var (
	ErrConnDone = sql.ErrConnDone
	ErrNoRows   = sql.ErrNoRows
	ErrTxDone   = sql.ErrTxDone
)

type Result = sql.Result

// Although the NameMapper is convenient, in practice it should not
// be relied on except for application code.  If you are writing a library
// that uses sqlx, you should be aware that the name mappings you expect
// can be overridden by your user's application.

// NameMapper is used to map column names to struct field names.  By default,
// it uses str.SnakeCase to snakecase struct field names.  It can be set
// to whatever you want, but it is encouraged to be set before sqlx is used
// as name-to-field mappings are cached after first use on a type.
var NameMapper = ref.NewMapperFunc("db", str.SnakeCase)

//------------------------------------------------
// sqlx interface
//

type Trace func(start time.Time, sql string, rows int64, err error)

type Supporter interface {
	SupportLastInsertID() bool
}

type Selector interface {
	Get(dest any, query string, args ...any) error
	Select(dest any, query string, args ...any) error
}

type Queryerx interface {
	Queryx(query string, args ...any) (*Rows, error)
	QueryRowx(query string, args ...any) *Row
}

type NamedQueryer interface {
	NamedQuery(query string, arg any) (*Rows, error)
	NamedQueryRow(query string, arg any) *Row
}

// NamedExecer is an interface used by MustExec
type NamedExecer interface {
	NamedExec(query string, arg any) (sql.Result, error)
}

type Preparerx interface {
	Preparex(query string) (*Stmt, error)
}

// ContextPreparerx is an interface used by PreparexContext.
type ContextPreparerx interface {
	PreparexContext(ctx context.Context, query string) (*Stmt, error)
}

// Beginxer is an interface used by Transaction
type Beginxer interface {
	Beginx() (*Tx, error)
}

// BeginTxxer is an interface used by Transactionx
type BeginTxxer interface {
	BeginTxx(context.Context, *sql.TxOptions) (*Tx, error)
}

// ContextQueryerx is an interface used by GetContext and SelectContext
type ContextQueryerx interface {
	QueryxContext(ctx context.Context, query string, args ...any) (*Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *Row
}

// ExtContext is a union interface which can bind, query, and exec, with Context
// used by NamedQueryContext and NamedExecContext.
type ExtContext interface {
	binder
	mapper
	sqx.ContextQueryer
	sqx.ContextExecer
	ContextQueryerx
}

// Bind is an interface for something which can bind queries (Tx, DB)
type Bind interface {
	Rebind(string) string
	BindNamed(string, any) (string, []any, error)
}

type Sqlx interface {
	sqx.Sql
	Supporter
	Bind
	Selector
	Queryerx
	NamedQueryer
	NamedExecer
	Preparerx
}

type Transactioner interface {
	Transaction(func(tx *Tx) error) error
}

type Transactionerx interface {
	Transactionx(ctx context.Context, opts *sql.TxOptions, fc func(tx *Tx) error) error
}

type Sqltx interface {
	Sqlx
	Transactioner
}

//------------------------------------------------
// internal interface
//

type unsafer interface {
	IsUnsafe() bool
}

type mapper interface {
	Mapper() *ref.Mapper
}

// Binder is an interface for something which can bind queries (Tx, DB)
type binder interface {
	Binder() Binder
}

// isqlx is a union interface which can bind, query, and exec, used by NamedQuery and NamedExec.
type isqlx interface {
	binder
	unsafer
	mapper
	Sqlx
}

type ext struct {
	driverName string
	unsafe     bool
	binder     Binder
	quoter     sqx.Quoter
	mapper     *ref.Mapper
	tracer     tracer
}

// DriverName returns the driverName passed to the Open function for this DB.
func (ext *ext) DriverName() string {
	return ext.driverName
}

// Binder returns the binder by driverName passed to the Open function for this DB.
func (ext *ext) Binder() Binder {
	return ext.binder
}

// Rebind transforms a query from QUESTION to the DB driver's bindvar type.
func (ext *ext) Rebind(query string) string {
	return ext.binder.Rebind(query)
}

// Quoter returns the quoter by driverName passed to the Open function for this DB.
func (ext *ext) Quoter() sqx.Quoter {
	return ext.quoter
}

// IsUnsafe returns the unsafe
func (ext *ext) IsUnsafe() bool {
	return ext.unsafe
}

// Quote quote string 's' with quote marks [2]rune, return (m[0] + s + m[1])
func (ext *ext) Quote(s string) string {
	return ext.quoter.Quote(s)
}

// Mapper returns the mapper
func (ext *ext) Mapper() *ref.Mapper {
	return ext.mapper
}

// SupportRetuning check sql driver support "RETUNING"
func (ext *ext) SupportLastInsertID() bool {
	return ext.binder != BindDollar
}

// Builder returns a new sql builder
func (ext *ext) Builder() *Builder {
	return &Builder{bid: ext.binder}
}

// NewDB returns a new sqlx DB wrapper for a pre-existing *sql.DB.  The
// driverName of the original database is required for named query support.
func NewDB(db *sql.DB, driverName string, trace ...Trace) *DB {
	ext := ext{
		driverName: driverName,
		binder:     GetBinder(driverName),
		quoter:     sqx.GetQuoter(driverName),
		mapper:     NameMapper,
	}
	if len(trace) > 0 {
		ext.tracer = tracer{Bind: ext.binder, Trace: trace[0]}
	}

	return &DB{db: db, ext: ext}
}

// Open is the same as sql.Open, but returns an *sqlx.DB instead.
func Open(driverName, dataSourceName string, trace ...Trace) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return NewDB(db, driverName, trace...), err
}

// MustOpen is the same as sql.Open, but returns an *sqlx.DB instead and panics on error.
func MustOpen(driverName, dataSourceName string) *DB {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}

// Connect to a database and verify with a ping.
func Connect(driverName, dataSourceName string, trace ...Trace) (*DB, error) {
	db, err := Open(driverName, dataSourceName, trace...)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// ConnectContext to a database and verify with a ping.
func ConnectContext(ctx context.Context, driverName, dataSourceName string) (*DB, error) {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return db, err
	}
	err = db.PingContext(ctx)
	return db, err
}

// MustConnect connects to a database and panics on error.
func MustConnect(driverName, dataSourceName string) *DB {
	db, err := Connect(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}

// Select executes a query using the provided Queryer, and StructScans each row
// into dest, which must be a slice.  If the slice elements are scannable, then
// the result set must have only one column.  Otherwise, StructScan is used.
// The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func Select(q Queryerx, dest any, query string, args ...any) error {
	rows, err := q.Queryx(query, args...)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// SelectContext executes a query using the provided Queryer, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func SelectContext(ctx context.Context, q ContextQueryerx, dest any, query string, args ...any) error {
	rows, err := q.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// Get does a QueryRow using the provided Queryer, and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func Get(q Queryerx, dest any, query string, args ...any) error {
	r := q.QueryRowx(query, args...)
	return r.scanAny(dest, false)
}

// GetContext does a QueryRow using the provided Queryer, and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func GetContext(ctx context.Context, q ContextQueryerx, dest any, query string, args ...any) error {
	r := q.QueryRowxContext(ctx, query, args...)
	return r.scanAny(dest, false)
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

	return
}

// Transactionx start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func Transactionx(ctx context.Context, db BeginTxxer, opts *sql.TxOptions, fc func(tx *Tx) error) (err error) {
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

	return
}
