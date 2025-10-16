package sqlx

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/askasoft/pango/asg"
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

type Binder = sqx.Binder

const (
	BindUnknown  = sqx.BindUnknown
	BindQuestion = sqx.BindQuestion
	BindDollar   = sqx.BindDollar
	BindColon    = sqx.BindColon
	BindAt       = sqx.BindAt
)

type Quoter = sqx.Quoter

var (
	QuoteDefault   = sqx.QuoteDefault
	QuoteBackticks = sqx.QuoteBackticks
	QuoteBrackets  = sqx.QuoteBrackets
)

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

type Queryerx interface {
	Queryx(query string, args ...any) (*Rows, error)
	QueryRowx(query string, args ...any) *Row
}
type ContextQueryerx interface {
	QueryxContext(ctx context.Context, query string, args ...any) (*Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *Row
}

type NamedQueryer interface {
	NamedQuery(query string, arg any) (*Rows, error)
	NamedQueryRow(query string, arg any) *Row
}
type ContextNamedQueryer interface {
	NamedQueryContext(ctx context.Context, query string, arg any) (*Rows, error)
	NamedQueryRowContext(ctx context.Context, query string, arg any) *Row
}

type NamedExecer interface {
	NamedExec(query string, arg any) (sql.Result, error)
}
type ContextNamedExecer interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

type Getter interface {
	Get(dest any, query string, args ...any) error
}
type ContextGetter interface {
	GetContext(ctx context.Context, dest any, query string, args ...any) error
}

type Selector interface {
	Select(dest any, query string, args ...any) error
}
type ContextSelector interface {
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type Creator interface {
	Create(query string, args ...any) (int64, error)
}
type ContextCreator interface {
	CreateContext(ctx context.Context, query string, args ...any) (int64, error)
}

type NamedGetter interface {
	NamedGet(dest any, query string, arg any) error
}
type ContextNamedGetter interface {
	NamedGetContext(ctx context.Context, dest any, query string, arg any) error
}

type NamedSelector interface {
	NamedSelect(dest any, query string, arg any) error
}
type ContextNamedSelector interface {
	NamedSelectContext(ctx context.Context, dest any, query string, arg any) error
}

type NamedCreator interface {
	NamedCreate(query string, arg any) (int64, error)
}
type ContextNamedCreator interface {
	NamedCreateContext(ctx context.Context, query string, arg any) (int64, error)
}

type Preparerx interface {
	Preparex(query string) (*Stmt, error)
}
type ContextPreparerx interface {
	PreparexContext(ctx context.Context, query string) (*Stmt, error)
}

type Beginxer interface {
	Beginx() (*Tx, error)
}
type BeginTxxer interface {
	BeginTxx(context.Context, *sql.TxOptions) (*Tx, error)
}

type BindNamed interface {
	BindNamed(string, any) (string, []any, error)
}

type Build interface {
	Builder() *Builder
}

type Sqlx interface {
	sqx.Quote
	sqx.Rebind
	sqx.Sql
	sqx.Sqlc

	Supporter
	BindNamed
	Build

	Queryerx
	Preparerx

	Getter
	Selector
	Creator

	NamedQueryer
	NamedExecer

	NamedGetter
	NamedSelector
	NamedCreator

	ContextQueryerx
	ContextNamedQueryer
	ContextNamedExecer
	ContextPreparerx

	ContextGetter
	ContextSelector
	ContextCreator

	ContextNamedGetter
	ContextNamedSelector
	ContextNamedCreator
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
	Supporter
	Queryerx
	sqx.Execer
}

// icsqlx is a union interface which can bind, query, and exec, used by NamedQueryContext and NamedExecContext.
type icsqlx interface {
	binder
	unsafer
	mapper
	Supporter
	ContextQueryerx
	sqx.ContextExecer
}

type ext struct {
	driverName string
	unsafe     bool
	binder     Binder
	quoter     Quoter
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

// Explain generate SQL string with given parameters.
func (ext *ext) Explain(sql string, args ...any) string {
	return ext.binder.Explain(sql, args...)
}

// Placeholder generate a place holder mark with No. n.
func (ext *ext) Placeholder(n int) string {
	return ext.binder.Placeholder(n)
}

// Quoter returns the quoter by driverName passed to the Open function for this DB.
func (ext *ext) Quoter() Quoter {
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

// SupportLastInsertID check sql driver support "LastInsertID()"
func (ext *ext) SupportLastInsertID() bool {
	return ext.binder != BindDollar
}

// Builder returns a new sql builder
func (ext *ext) Builder() *Builder {
	return &Builder{mpr: ext.mapper, sqb: sqx.Builder{Binder: ext.binder, Quoter: ext.quoter}}
}

// StructFields returns struct mapped fields
func (ext *ext) StructFields(a any, omits ...string) (fields []string) {
	sm := ext.mapper.TypeMap(reflect.TypeOf(a))
	for _, fi := range sm.Index {
		if isIgnoredField(fi, omits...) {
			continue
		}
		fields = append(fields, fi.Name)
	}
	return
}

func isIgnoredField(fi *ref.FieldInfo, omits ...string) bool {
	if fi.Embedded || asg.Contains(omits, fi.Name) {
		return true
	}

	for fi = fi.Parent; fi != nil && fi.Path != ""; fi = fi.Parent {
		if !fi.Embedded && !asg.Contains(omits, fi.Name) {
			return true
		}
	}

	return false
}

// NewDB returns a new sqlx DB wrapper for a pre-existing *sql.DB.  The
// driverName of the original database is required for named query support.
func NewDB(db *sql.DB, driverName string, trace ...Trace) *DB {
	ext := ext{
		driverName: driverName,
		binder:     sqx.GetBinder(driverName),
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

// Select executes a query, and StructScans each row
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

// NamedSelect executes a query, and StructScans each row
// into dest, which must be a slice.  If the slice elements are scannable, then
// the result set must have only one column.  Otherwise, StructScan is used.
// The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func NamedSelect(q NamedQueryer, dest any, query string, arg any) error {
	rows, err := q.NamedQuery(query, arg)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// SelectContext executes a query, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied arg.
func SelectContext(ctx context.Context, q ContextQueryerx, dest any, query string, args ...any) error {
	rows, err := q.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// NamedSelectContext executes a query, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied arg.
func NamedSelectContext(ctx context.Context, q ContextNamedQueryer, dest any, query string, arg any) error {
	rows, err := q.NamedQueryContext(ctx, query, arg)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// Get does a QueryRowx() and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func Get(q Queryerx, dest any, query string, args ...any) error {
	r := q.QueryRowx(query, args...)
	return r.scanAny(dest, false)
}

// NamedGet does a NamedQueryRow() and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func NamedGet(q NamedQueryer, dest any, query string, arg any) error {
	r := q.NamedQueryRow(query, arg)
	return r.scanAny(dest, false)
}

// GetContext does a QueryRowxContext() and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func GetContext(ctx context.Context, q ContextQueryerx, dest any, query string, args ...any) error {
	r := q.QueryRowxContext(ctx, query, args...)
	return r.scanAny(dest, false)
}

// NamedGetContext does a NamedQueryRowContext() and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func NamedGetContext(ctx context.Context, q ContextNamedQueryer, dest any, query string, arg any) error {
	r := q.NamedQueryRowContext(ctx, query, arg)
	return r.scanAny(dest, false)
}

// Create does a QueryRowx() and scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() and returns Result.LastInsertId().
func Create(x Sqlx, query string, args ...any) (int64, error) {
	if x.SupportLastInsertID() {
		r, err := x.Exec(query, args...)
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	}

	r := x.QueryRowx(query, args...)
	if r.Err() != nil {
		return 0, r.Err()
	}

	var id int64
	err := r.Scan(&id)
	return id, err
}

// NamedCreate does a NamedQueryRowx() and scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() return Result.LastInsertId().
func NamedCreate(x Sqlx, query string, arg any) (int64, error) {
	if x.SupportLastInsertID() {
		r, err := x.NamedExec(query, arg)
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	}

	r := x.NamedQueryRow(query, arg)
	if r.Err() != nil {
		return 0, r.Err()
	}

	var id int64
	err := r.Scan(&id)
	return id, err
}

// CreateContext does a QueryRowxContext() scans the resulting row returns the last inserted ID.
// If the db supports LastInsertId(), do a Exec() return Result.LastInsertId().
func CreateContext(ctx context.Context, x Sqlx, query string, args ...any) (int64, error) {
	if x.SupportLastInsertID() {
		r, err := x.ExecContext(ctx, query, args...)
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	}

	r := x.QueryRowxContext(ctx, query, args...)
	if r.Err() != nil {
		return 0, r.Err()
	}

	var id int64
	err := r.Scan(&id)
	return id, err
}

// NamedCreateContext does a NamedQueryRow() and scans the resulting row
// returns the last inserted ID.
// If the db supports LastInsertId(), does a NamedExecContext() and returns Result.LastInsertId().
func NamedCreateContext(ctx context.Context, x Sqlx, query string, arg any) (int64, error) {
	if x.SupportLastInsertID() {
		r, err := x.NamedExecContext(ctx, query, arg)
		if err != nil {
			return 0, err
		}
		return r.LastInsertId()
	}

	r := x.NamedQueryRowContext(ctx, query, arg)
	if r.Err() != nil {
		return 0, r.Err()
	}

	var id int64
	err := r.Scan(&id)
	return id, err
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func Transaction(db Beginxer, fc func(tx *Tx) error) (err error) {
	var tx *Tx
	var done bool

	tx, err = db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		// Make sure to rollback when panic
		if err != nil || !done {
			_ = tx.Rollback()
		}
	}()

	if err = fc(tx); err == nil {
		err = tx.Commit()
	}

	done = true
	return
}

// Transactionx start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func Transactionx(ctx context.Context, db BeginTxxer, opts *sql.TxOptions, fc func(tx *Tx) error) (err error) {
	var tx *Tx
	var done bool

	tx, err = db.BeginTxx(ctx, opts)
	if err != nil {
		return
	}

	defer func() {
		// Make sure to rollback when panic
		if err != nil || !done {
			_ = tx.Rollback()
		}
	}()

	if err = fc(tx); err == nil {
		err = tx.Commit()
	}

	done = true
	return
}

// Named takes a query using named parameters and an argument and
// returns a new query with a list of args that can be executed by
// a database.  The return value uses the `?` bindvar.
func Named(query string, arg any) (string, []any, error) {
	return bindNamedMapper(BindQuestion, query, arg, NameMapper)
}

// namedExec uses BindStruct to get a query executable by the driver and
// then runs Exec on the result.  Returns an error from the binding
// or the query execution itself.
func namedExec(x isqlx, query string, arg any) (sql.Result, error) {
	q, args, err := bindNamedMapper(x.Binder(), query, arg, x.Mapper())
	if err != nil {
		return nil, err
	}
	return x.Exec(q, args...)
}

// namedExecContext uses BindStruct to get a query executable by the driver and
// then runs Exec on the result.  Returns an error from the binding
// or the query execution itself.
func namedExecContext(ctx context.Context, x icsqlx, query string, arg any) (sql.Result, error) {
	q, args, err := bindNamedMapper(x.Binder(), query, arg, x.Mapper())
	if err != nil {
		return nil, err
	}
	return x.ExecContext(ctx, q, args...)
}

// namedQuery binds a named query and then runs Query on the result using the
// provided Ext (sqlx.Tx, sqlx.Db).  It works with both structs and with
// map[string]any types.
func namedQuery(x isqlx, query string, arg any) (*Rows, error) {
	q, args, err := bindNamedMapper(x.Binder(), query, arg, x.Mapper())
	if err != nil {
		return nil, err
	}
	return x.Queryx(q, args...)
}

// namedQueryContext binds a named query and then runs Query on the result using the
// provided Ext (sqlx.Tx, sqlx.Db).  It works with both structs and with
// map[string]any types.
func namedQueryContext(ctx context.Context, x icsqlx, query string, arg any) (*Rows, error) {
	q, args, err := bindNamedMapper(x.Binder(), query, arg, x.Mapper())
	if err != nil {
		return nil, err
	}
	return x.QueryxContext(ctx, q, args...)
}

// namedQueryRow binds a named query and then runs Query on the result using the
// provided Ext (sqlx.Tx, sqlx.Db).  It works with both structs and with
// map[string]any types.
func namedQueryRow(x isqlx, query string, arg any) *Row {
	q, args, err := bindNamedMapper(x.Binder(), query, arg, x.Mapper())
	if err != nil {
		return &Row{err: err}
	}
	return x.QueryRowx(q, args...)
}

// namedQueryRowContext binds a named query and then runs Query on the result using the
// provided Ext (sqlx.Tx, sqlx.Db).  It works with both structs and with
// map[string]any types.
func namedQueryRowContext(ctx context.Context, x icsqlx, query string, arg any) *Row {
	q, args, err := bindNamedMapper(x.Binder(), query, arg, x.Mapper())
	if err != nil {
		return &Row{err: err}
	}
	return x.QueryRowxContext(ctx, q, args...)
}
