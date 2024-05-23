package sqlx

import (
	"database/sql"
	"errors"
	"fmt"

	"reflect"

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
// that uses sqx, you should be aware that the name mappings you expect
// can be overridden by your user's application.

// NameMapper is used to map column names to struct field names.  By default,
// it uses str.SnakeCase to snakecase struct field names.  It can be set
// to whatever you want, but it is encouraged to be set before sqx is used
// as name-to-field mappings are cached after first use on a type.
var NameMapper = ref.NewMapperFunc("db", str.SnakeCase)

//------------------------------------------------
// GO database/sql interface
//

// Queryer is an interface used by Get and Select
type Queryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// Execer is an interface used by MustExec
type Execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

// Sql the basic interface for sql.DB, sql.Tx
type Sql interface {
	Queryer
	Execer
	Preparer
}

//------------------------------------------------
// sqx interface
//

// ColScanner is an interface used by MapScan and SliceScan
type ColScanner interface {
	Columns() ([]string, error)
	Scan(dest ...any) error
	Err() error
}

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

// Bind is an interface for something which can bind queries (Tx, DB)
type Bind interface {
	Rebind(string) string
	BindNamed(string, any) (string, []any, error)
}

type Sqlx interface {
	Sql
	Supporter
	Bind
	Selector
	Queryerx
	NamedQueryer
	NamedExecer
	Preparerx
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

// determine if any of our extensions are unsafe
func isUnsafe(i any) bool {
	if us, ok := i.(unsafer); ok {
		return us.IsUnsafe()
	}
	return false
}

type ext struct {
	driverName string
	binder     Binder
	quoter     sqx.Quoter
	mapper     *ref.Mapper
	unsafe     bool
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

var _scannerInterface = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

// Row is a reimplementation of sql.Row in order to gain access to the underlying
// sql.Rows.Columns() data, necessary for StructScan.
type Row struct {
	rows *sql.Rows
	err  error
	ext
}

// Scan is a fixed implementation of sql.Row.Scan, which does not discard the
// underlying error from the internal rows object if it exists.
func (r *Row) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	// TODO(bradfitz): for now we need to defensively clone all
	// []byte that the driver returned (not permitting
	// *RawBytes in Rows.Scan), since we're about to close
	// the Rows in our defer, when we return from this function.
	// the contract with the driver.Next(...) interface is that it
	// can return slices into read-only temporary memory that's
	// only valid until the next Scan/Close.  But the TODO is that
	// for a lot of drivers, this copy will be unnecessary.  We
	// should provide an optional interface for drivers to
	// implement to say, "don't worry, the []bytes that I return
	// from Next will not be modified again." (for instance, if
	// they were obtained from the network anyway) But for now we
	// don't care.
	defer r.rows.Close()
	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return ErrNoRows
	}
	err := r.rows.Scan(dest...)
	if err != nil {
		return err
	}
	// Make sure the query can be processed to completion with no errors.
	if err := r.rows.Close(); err != nil {
		return err
	}
	return nil
}

// Columns returns the underlying sql.Rows.Columns(), or the deferred error usually
// returned by Row.Scan()
func (r *Row) Columns() ([]string, error) {
	if r.err != nil {
		return []string{}, r.err
	}
	return r.rows.Columns()
}

// ColumnTypes returns the underlying sql.Rows.ColumnTypes(), or the deferred error
func (r *Row) ColumnTypes() ([]*sql.ColumnType, error) {
	if r.err != nil {
		return []*sql.ColumnType{}, r.err
	}
	return r.rows.ColumnTypes()
}

// Err returns the error encountered while scanning.
func (r *Row) Err() error {
	return r.err
}

// DB is a wrapper around sql.DB which keeps track of the driverName upon Open,
// used mostly to automatically bind named queries using the right bindvars.
type DB struct {
	*sql.DB
	ext
}

// NewDB returns a new sqx DB wrapper for a pre-existing *sql.DB.  The
// driverName of the original database is required for named query support.
func NewDB(db *sql.DB, driverName string) *DB {
	return &DB{DB: db, ext: ext{driverName: driverName, binder: GetBinder(driverName), quoter: sqx.GetQuoter(driverName), mapper: NameMapper}}
}

// Open is the same as sql.Open, but returns an *sqx.DB instead.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return NewDB(db, driverName), err
}

// MustOpen is the same as sql.Open, but returns an *sqx.DB instead and panics on error.
func MustOpen(driverName, dataSourceName string) *DB {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}

// MapperFunc sets a new mapper for this db using the default sqx struct tag
// and the provided mapper function.
func (db *DB) MapperFunc(mf func(string) string) {
	db.mapper = ref.NewMapperFunc("db", mf)
}

// Unsafe returns a version of DB which will silently succeed to scan when
// columns in the SQL result have no fields in the destination struct.
// sqx.Stmt and sqx.Tx which are created from this DB will inherit its
// safety behavior.
func (db *DB) Unsafe() *DB {
	ndb := &DB{DB: db.DB, ext: db.ext}
	ndb.unsafe = true
	return ndb
}

// BindNamed binds a query using the DB driver's bindvar type.
func (db *DB) BindNamed(query string, arg any) (string, []any, error) {
	return db.binder.bindNamedMapper(query, arg, db.mapper)
}

// NamedQuery using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedQuery(query string, arg any) (*Rows, error) {
	return namedQuery(db, query, arg)
}

// NamedQueryRow using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedQueryRow(query string, arg any) *Row {
	return namedQueryRow(db, query, arg)
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *DB) NamedExec(query string, arg any) (sql.Result, error) {
	return namedExec(db, query, arg)
}

// Select using this DB.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) Select(dest any, query string, args ...any) error {
	return Select(db, dest, query, args...)
}

// Get using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) Get(dest any, query string, args ...any) error {
	return Get(db, dest, query, args...)
}

// MustBegin starts a transaction, and panics on error.  Returns an *sqx.Tx instead
// of an *sql.Tx.
func (db *DB) MustBegin() *Tx {
	tx, err := db.Beginx()
	if err != nil {
		panic(err)
	}
	return tx
}

// Beginx begins a transaction and returns an *sqx.Tx instead of an *sql.Tx.
func (db *DB) Beginx() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx, ext: db.ext}, err
}

// Queryx queries the database and returns an *sqx.Rows.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) Queryx(query string, args ...any) (*Rows, error) {
	r, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: db.ext}, err
}

// QueryRowx queries the database and returns an *sqx.Row.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) QueryRowx(query string, args ...any) *Row {
	rows, err := db.DB.Query(query, args...)
	return &Row{rows: rows, err: err, ext: db.ext}
}

// MustExec (panic) runs MustExec using this database.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) MustExec(query string, args ...any) sql.Result {
	return MustExec(db, query, args...)
}

// Preparex returns an sqx.Stmt instead of a sql.Stmt
func (db *DB) Preparex(query string) (*Stmt, error) {
	s, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{Stmt: s, ext: db.ext}, err
}

// PrepareNamed returns an sqx.NamedStmt
func (db *DB) PrepareNamed(query string) (*NamedStmt, error) {
	return prepareNamed(db, query)
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit. Transaction executes an
// arbitrary number of commands in fc within a transaction. On success the changes are committed; if an error occurs
// they are rolled back.
func (db *DB) Transaction(fc func(tx *Tx) error) (err error) {
	return Transaction(db, fc)
}

// Conn is a wrapper around sql.Conn with extra functionality
type Conn struct {
	*sql.Conn
	ext
}

// Tx is an sqx wrapper around sql.Tx with extra functionality
type Tx struct {
	*sql.Tx
	ext
}

// Unsafe returns a version of Tx which will silently succeed to scan when
// columns in the SQL result have no fields in the destination struct.
func (tx *Tx) Unsafe() *Tx {
	ntx := &Tx{Tx: tx.Tx, ext: tx.ext}
	ntx.unsafe = true
	return ntx
}

// BindNamed binds a query within a transaction's bindvar type.
func (tx *Tx) BindNamed(query string, arg any) (string, []any, error) {
	return tx.binder.bindNamedMapper(query, arg, tx.mapper)
}

// NamedQuery within a transaction.
// Any named placeholder parameters are replaced with fields from arg.
func (tx *Tx) NamedQuery(query string, arg any) (*Rows, error) {
	return namedQuery(tx, query, arg)
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

// Select within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) Select(dest any, query string, args ...any) error {
	return Select(tx, dest, query, args...)
}

// Queryx within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) Queryx(query string, args ...any) (*Rows, error) {
	r, err := tx.Tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: tx.ext}, err
}

// QueryRowx within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryRowx(query string, args ...any) *Row {
	rows, err := tx.Tx.Query(query, args...)
	return &Row{rows: rows, err: err, ext: tx.ext}
}

// Get within a transaction.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (tx *Tx) Get(dest any, query string, args ...any) error {
	return Get(tx, dest, query, args...)
}

// MustExec runs MustExec within a transaction.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) MustExec(query string, args ...any) sql.Result {
	return MustExec(tx, query, args...)
}

// Preparex  a statement within a transaction.
func (tx *Tx) Preparex(query string) (*Stmt, error) {
	s, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{Stmt: s, ext: tx.ext}, err
}

// Stmtx returns a version of the prepared statement which runs within a transaction.  Provided
// stmt can be either *sql.Stmt or *sqx.Stmt.
func (tx *Tx) Stmtx(stmt any) *Stmt {
	var s *sql.Stmt
	switch v := stmt.(type) {
	case Stmt:
		s = v.Stmt
	case *Stmt:
		s = v.Stmt
	case *sql.Stmt:
		s = v
	default:
		panic(fmt.Sprintf("non-statement type %v passed to Stmtx", reflect.ValueOf(stmt).Type()))
	}
	return &Stmt{Stmt: tx.Stmt(s), ext: tx.ext}
}

// NamedStmt returns a version of the prepared statement which runs within a transaction.
func (tx *Tx) NamedStmt(stmt *NamedStmt) *NamedStmt {
	return &NamedStmt{
		QueryString: stmt.QueryString,
		Params:      stmt.Params,
		Stmt:        tx.Stmtx(stmt.Stmt),
	}
}

// PrepareNamed returns an sqx.NamedStmt
func (tx *Tx) PrepareNamed(query string) (*NamedStmt, error) {
	return prepareNamed(tx, query)
}

// Stmt is an sqx wrapper around sql.Stmt with extra functionality
type Stmt struct {
	*sql.Stmt
	ext
}

// Unsafe returns a version of Stmt which will silently succeed to scan when
// columns in the SQL result have no fields in the destination struct.
func (s *Stmt) Unsafe() *Stmt {
	c := &Stmt{Stmt: s.Stmt, ext: s.ext}
	c.unsafe = true
	return c
}

// Select using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) Select(dest any, args ...any) error {
	return Select(&qStmt{s}, dest, "", args...)
}

// Get using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (s *Stmt) Get(dest any, args ...any) error {
	return Get(&qStmt{s}, dest, "", args...)
}

// MustExec (panic) using this statement.  Note that the query portion of the error
// output will be blank, as Stmt does not expose its query.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) MustExec(args ...any) sql.Result {
	return MustExec(&qStmt{s}, "", args...)
}

// QueryRowx using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) QueryRowx(args ...any) *Row {
	qs := &qStmt{s}
	return qs.QueryRowx("", args...)
}

// Queryx using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) Queryx(args ...any) (*Rows, error) {
	qs := &qStmt{s}
	return qs.Queryx("", args...)
}

// qStmt is an unexposed wrapper which lets you use a Stmt as a Queryer & Execer by
// implementing those interfaces and ignoring the `query` argument.
type qStmt struct{ *Stmt }

func (q *qStmt) Query(query string, args ...any) (*sql.Rows, error) {
	return q.Stmt.Query(args...)
}

func (q *qStmt) Queryx(query string, args ...any) (*Rows, error) {
	r, err := q.Stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, ext: q.Stmt.ext}, err
}

func (q *qStmt) QueryRowx(query string, args ...any) *Row {
	rows, err := q.Stmt.Query(args...)
	return &Row{rows: rows, err: err, ext: q.Stmt.ext}
}

func (q *qStmt) Exec(query string, args ...any) (sql.Result, error) {
	return q.Stmt.Exec(args...)
}

// Rows is a wrapper around sql.Rows which caches costly reflect operations
// during a looped StructScan
type Rows struct {
	*sql.Rows
	ext

	// these fields cache memory use for a rows during iteration w/ structScan
	started bool
	fields  [][]int
	values  []any
}

// SliceScan using this Rows.
func (r *Rows) SliceScan() ([]any, error) {
	return SliceScan(r)
}

// MapScan using this Rows.
func (r *Rows) MapScan(dest map[string]any) error {
	return MapScan(r, dest)
}

// StructScan is like sql.Rows.Scan, but scans a single Row into a single Struct.
// Use this and iterate over Rows manually when the memory load of Select() might be
// prohibitive.  *Rows.StructScan caches the reflect work of matching up column
// positions to fields to avoid that overhead per scan, which means it is not safe
// to run StructScan on the same Rows instance with different struct types.
func (r *Rows) StructScan(dest any) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	v = v.Elem()

	if !r.started {
		columns, err := r.Columns()
		if err != nil {
			return err
		}
		m := r.mapper

		r.fields = m.TraversalsByName(v.Type(), columns)
		// if we are not unsafe and are missing fields, return an error
		if f, err := missingFields(r.fields); err != nil && !r.unsafe {
			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
		}
		r.values = make([]any, len(columns))
		r.started = true
	}

	err := fieldsByTraversal(v, r.fields, r.values, true)
	if err != nil {
		return err
	}
	// scan into the struct field pointers and append to our results
	err = r.Scan(r.values...)
	if err != nil {
		return err
	}
	return r.Err()
}

// Connect to a database and verify with a ping.
func Connect(driverName, dataSourceName string) (*DB, error) {
	db, err := Open(driverName, dataSourceName)
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

// Get does a QueryRow using the provided Queryer, and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func Get(q Queryerx, dest any, query string, args ...any) error {
	r := q.QueryRowx(query, args...)
	return r.scanAny(dest, false)
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

// SliceScan using this Rows.
func (r *Row) SliceScan() ([]any, error) {
	return SliceScan(r)
}

// MapScan using this Rows.
func (r *Row) MapScan(dest map[string]any) error {
	return MapScan(r, dest)
}

// isScannable takes the reflect.Type and the actual dest value and returns
// whether or not it's Scannable.  Something is scannable if:
//   - it is not a struct
//   - it implements sql.Scanner
//   - it has no exported fields
func isScannable(t reflect.Type) bool {
	if reflect.PtrTo(t).Implements(_scannerInterface) {
		return true
	}
	if t.Kind() != reflect.Struct {
		return true
	}

	// it's not important that we use the right mapper for this particular object,
	// we're only concerned on how many exported fields this struct has
	return len(NameMapper.TypeMap(t).Index) == 0
}

func (r *Row) scanAny(dest any, structOnly bool) error {
	if r.err != nil {
		return r.err
	}
	if r.rows == nil {
		r.err = ErrNoRows
		return r.err
	}
	defer r.rows.Close()

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	base := ref.Deref(v.Type())
	scannable := isScannable(base)

	if structOnly && scannable {
		return structOnlyError(base)
	}

	columns, err := r.Columns()
	if err != nil {
		return err
	}

	if scannable && len(columns) > 1 {
		return fmt.Errorf("scannable dest type %s with >1 columns (%d) in result", base.Kind(), len(columns))
	}

	if scannable {
		return r.Scan(dest)
	}

	m := r.mapper

	fields := m.TraversalsByName(v.Type(), columns)
	// if we are not unsafe and are missing fields, return an error
	if f, err := missingFields(fields); err != nil && !r.unsafe {
		return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
	}
	values := make([]any, len(columns))

	err = fieldsByTraversal(v, fields, values, true)
	if err != nil {
		return err
	}
	// scan into the struct field pointers and append to our results
	return r.Scan(values...)
}

// StructScan a single Row into dest.
func (r *Row) StructScan(dest any) error {
	return r.scanAny(dest, true)
}

// SliceScan a row, returning a []any with values similar to MapScan.
// This function is primarily intended for use where the number of columns
// is not known.  Because you can pass an []any directly to Scan,
// it's recommended that you do that as it will not have to allocate new
// slices per row.
func SliceScan(r ColScanner) ([]any, error) {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
	if err != nil {
		return []any{}, err
	}

	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(any)
	}

	err = r.Scan(values...)

	if err != nil {
		return values, err
	}

	for i := range columns {
		values[i] = *(values[i].(*any))
	}

	return values, r.Err()
}

// MapScan scans a single Row into the dest map[string]any.
// Use this to get results for SQL that might not be under your control
// (for instance, if you're building an interface for an SQL server that
// executes SQL from input).  Please do not use this as a primary interface!
// This will modify the map sent to it in place, so reuse the same map with
// care.  Columns which occur more than once in the result will overwrite
// each other!
func MapScan(r ColScanner, dest map[string]any) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(any)
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*any))
	}

	return r.Err()
}

type rowsi interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...any) error
}

// structOnlyError returns an error appropriate for type when a non-scannable
// struct is expected but something else is given
func structOnlyError(t reflect.Type) error {
	isStruct := t.Kind() == reflect.Struct
	isScanner := reflect.PtrTo(t).Implements(_scannerInterface)
	if !isStruct {
		return fmt.Errorf("expected %s but got %s", reflect.Struct, t.Kind())
	}
	if isScanner {
		return fmt.Errorf("structscan expects a struct dest but the provided struct type %s implements scanner", t.Name())
	}
	return fmt.Errorf("expected a struct, but struct %s has no exported fields", t.Name())
}

// scanAll scans all rows into a destination, which must be a slice of any
// type.  It resets the slice length to zero before appending each element to
// the slice.  If the destination slice type is a Struct, then StructScan will
// be used on each row.  If the destination is some other kind of base type,
// then each row must only have one column which can scan into that type.  This
// allows you to do something like:
//
//	rows, _ := db.Query("select id from people;")
//	var ids []int
//	scanAll(rows, &ids, false)
//
// and ids will be a list of the id results.  I realize that this is a desirable
// interface to expose to users, but for now it will only be exposed via changes
// to `Get` and `Select`.  The reason that this has been implemented like this is
// this is the only way to not duplicate reflect work in the new API while
// maintaining backwards compatibility.
func scanAll(rows rowsi, dest any, structOnly bool) error {
	var v, vp reflect.Value

	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	direct := reflect.Indirect(value)

	slice, err := baseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}
	direct.SetLen(0)

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := ref.Deref(slice.Elem())
	scannable := isScannable(base)

	if structOnly && scannable {
		return structOnlyError(base)
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// if it's a base type make sure it only has 1 column;  if not return an error
	if scannable && len(columns) > 1 {
		return fmt.Errorf("non-struct dest type %s with >1 columns (%d)", base.Kind(), len(columns))
	}

	if !scannable {
		var values []any
		var m *ref.Mapper

		if rs, ok := rows.(*Rows); ok {
			m = rs.mapper
		} else {
			m = NameMapper
		}

		fields := m.TraversalsByName(base, columns)
		// if we are not unsafe and are missing fields, return an error
		if f, err := missingFields(fields); err != nil && !isUnsafe(rows) {
			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
		}
		values = make([]any, len(columns))

		for rows.Next() {
			// create a new struct type (which returns PtrTo) and indirect it
			vp = reflect.New(base)
			v = reflect.Indirect(vp)

			err = fieldsByTraversal(v, fields, values, true)
			if err != nil {
				return err
			}

			// scan into the struct field pointers and append to our results
			err = rows.Scan(values...)
			if err != nil {
				return err
			}

			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, v))
			}
		}
	} else {
		for rows.Next() {
			vp = reflect.New(base)
			err = rows.Scan(vp.Interface())
			if err != nil {
				return err
			}
			// append
			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
			}
		}
	}

	return rows.Err()
}

// FIXME: StructScan was the very first bit of API in sqx, and now unfortunately
// it doesn't really feel like it's named properly.  There is an incongruency
// between this and the way that StructScan (which might better be ScanStruct
// anyway) works on a rows object.

// StructScan all rows from an sql.Rows or an sqx.Rows into the dest slice.
// StructScan will scan in the entire rows result, so if you do not want to
// allocate structs for the entire result, use Queryx and see sqx.Rows.StructScan.
// If rows is sqx.Rows, it will use its mapper, otherwise it will use the default.
func StructScan(rows rowsi, dest any) error {
	return scanAll(rows, dest, true)

}

// reflect helpers

func baseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = ref.Deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

// fieldsByName fills a values interface with fields from the passed value based
// on the traversals in int.  If ptrs is true, return addresses instead of values.
// We write this instead of using FieldsByName to save allocations and map lookups
// when iterating over many rows.  Empty traversals will get an interface pointer.
// Because of the necessity of requesting ptrs or values, it's considered a bit too
// specialized for inclusion in ref itself.
func fieldsByTraversal(v reflect.Value, traversals [][]int, values []any, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(any)
			continue
		}
		f := ref.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}

func missingFields(transversals [][]int) (field int, err error) {
	for i, t := range transversals {
		if len(t) == 0 {
			return i, errors.New("missing field")
		}
	}
	return 0, nil
}
