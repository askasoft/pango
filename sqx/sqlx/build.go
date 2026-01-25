package sqlx

import (
	"reflect"

	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/sqx"
)

// Builder a simple sql builder
// NOTE: the arguments are strict to it's order
type Builder struct {
	mpr *ref.Mapper
	sqb sqx.Builder
}

// Clone returns a copy of the builder
func (b *Builder) Clone() *Builder {
	n := &Builder{
		mpr: b.mpr,
		sqb: *b.sqb.Clone(),
	}
	return n
}

// Rebind a SQL from the default binder (QUESTION) to the target binder.
func (b *Builder) Rebind(sql string) string {
	return b.sqb.Rebind(sql)
}

func (b *Builder) Placeholder(n int) string {
	return b.sqb.Placeholder(n)
}

func (b *Builder) Explain(sql string, args ...any) string {
	return b.sqb.Explain(sql, args...)
}

// Quotes quote string 's' in 'ss' with quote marks [2]rune, return (m[0] + s + m[1])
func (b *Builder) Quotes(ss ...string) []string {
	return b.sqb.Quotes(ss...)
}

// Quote quote string 's' with quotes [2]rune.
// Returns (quoter[0] + s + quoter[1]), if 's' does not contains any "!\"#$%&'()*+,-/:;<=>?@[\\]^`{|}~" characters.
func (b *Builder) Quote(s string) string {
	return b.sqb.Quote(s)
}

func (b *Builder) Reset() *Builder {
	b.sqb.Reset()
	return b
}

// Build returns (sql, parameters)
func (b *Builder) Build() (string, []any) {
	return b.sqb.Build()
}

// Params returns parameters
func (b *Builder) Params() []any {
	return b.sqb.Params()
}

// SQL returns sql
func (b *Builder) SQL() string {
	return b.sqb.SQL()
}

// SQLWhere returns sql after WHERE
func (b *Builder) SQLWhere() string {
	return b.sqb.SQLWhere()
}

// Count shortcut for SELECT COUNT(*)
func (b *Builder) Count(cols ...string) *Builder {
	b.sqb.Count(cols...)
	return b
}

// Count shortcut for SELECT COUNT(distinct *)
func (b *Builder) CountDistinct(cols ...string) *Builder {
	b.sqb.CountDistinct(cols...)
	return b
}

// Select add select columns
// if `cols` is not specified, default select "*"
func (b *Builder) Select(cols ...string) *Builder {
	b.sqb.Select(cols...)
	return b
}

// ForUpdate add 'FOR UPDATE' for SELECT
func (b *Builder) ForUpdate() *Builder {
	b.sqb.ForUpdate()
	return b
}

// Distinct add 'DISTINCT' for SELECT
func (b *Builder) Distinct() *Builder {
	b.sqb.Distinct()
	return b
}

func (b *Builder) Delete(tb string) *Builder {
	b.sqb.Delete(tb)
	return b
}

func (b *Builder) Insert(tb string) *Builder {
	b.sqb.Insert(tb)
	return b
}

func (b *Builder) Update(tb string) *Builder {
	b.sqb.Update(tb)
	return b
}

func (b *Builder) From(tb string, args ...any) *Builder {
	b.sqb.From(tb, args...)
	return b
}

// Columns add select columns
func (b *Builder) Columns(cols ...string) *Builder {
	b.sqb.Columns(cols...)
	return b
}

// Columnx add select column and parameters
func (b *Builder) Columnx(col string, args ...any) *Builder {
	b.sqb.Columnx(col, args...)
	return b
}

// Values add insert/update values
func (b *Builder) Values(vals ...string) *Builder {
	b.sqb.Values(vals...)
	return b
}

// Join specify Join query and conditions
//
//	sqb.Join("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "abc@example.org")
func (b *Builder) Join(query string, args ...any) *Builder {
	b.sqb.Join(query, args...)
	return b
}

func (b *Builder) Setx(col string, val string, args ...any) *Builder {
	b.sqb.Setx(col, val, args...)
	return b
}

func (b *Builder) Setc(col string, arg any) *Builder {
	b.sqb.Setc(col, arg)
	return b
}

// Name add named column and value.
// Example:
//
//	sqb.Insert("a").Name("id") // INSERT INTO a (id) VALUES (:id)
//	sqb.Update("a").Name("value").Where("id = :id") // UPDATE a SET value = :value WHERE id = :id
func (b *Builder) Name(col string) *Builder {
	b.sqb.Name(col)
	return b
}

// Names add named columns and values.
// Example:
//
//	sqb.Insert("a").Names("id", "name", "value") // INSERT INTO a (id, name, value) VALUES (:id, :name,, :value)
//	sqb.Update("a").Names("name", "value").Where("id = :id") // UPDATE a SET name = :name, value = :value WHERE id = :id
func (b *Builder) Names(cols ...string) *Builder {
	b.sqb.Names(cols...)
	return b
}

func (b *Builder) Omits(cols ...string) *Builder {
	b.sqb.Omits(cols...)
	return b
}

func (b *Builder) Where(q string, args ...any) *Builder {
	b.sqb.Where(q, args...)
	return b
}

// LP append '(' to WHERE
func (b *Builder) LP() *Builder {
	b.sqb.LP()
	return b
}

// RP append ')' to WHERE
func (b *Builder) RP() *Builder {
	b.sqb.RP()
	return b
}

// OR switch to OR mode (default is AND mode)
func (b *Builder) OR() *Builder {
	b.sqb.OR()
	return b
}

// AND switch to AND mode (default is AND mode)
func (b *Builder) AND() *Builder {
	b.sqb.AND()
	return b
}

func (b *Builder) IsNull(col string) *Builder {
	b.sqb.IsNull(col)
	return b
}

func (b *Builder) NotNull(col string) *Builder {
	b.sqb.NotNull(col)
	return b
}

func (b *Builder) Eq(col string, val any) *Builder {
	b.sqb.Eq(col, val)
	return b
}

func (b *Builder) Neq(col string, val any) *Builder {
	b.sqb.Neq(col, val)
	return b
}

func (b *Builder) Gt(col string, val any) *Builder {
	b.sqb.Gt(col, val)
	return b
}

func (b *Builder) Gte(col string, val any) *Builder {
	b.sqb.Gte(col, val)
	return b
}

func (b *Builder) Lt(col string, val any) *Builder {
	b.sqb.Lt(col, val)
	return b
}

func (b *Builder) Lte(col string, val any) *Builder {
	b.sqb.Lte(col, val)
	return b
}

func (b *Builder) Like(col string, val any) *Builder {
	b.sqb.Like(col, val)
	return b
}

func (b *Builder) ILike(col string, val any) *Builder {
	b.sqb.ILike(col, val)
	return b
}

func (b *Builder) NotLike(col string, val any) *Builder {
	b.sqb.NotLike(col, val)
	return b
}

func (b *Builder) NotILike(col string, val any) *Builder {
	b.sqb.NotILike(col, val)
	return b
}

func (b *Builder) Btw(col string, vmin, vmax any) *Builder {
	b.sqb.Btw(col, vmin, vmax)
	return b
}

func (b *Builder) Nbtw(col string, vmin, vmax any) *Builder {
	b.sqb.Nbtw(col, vmin, vmax)
	return b
}

func (b *Builder) In(col string, val any) *Builder {
	b.sqb.In(col, val)
	return b
}

func (b *Builder) NotIn(col string, val any) *Builder {
	b.sqb.NotIn(col, val)
	return b
}

// Returns add RETURNING cols...
// if `cols` is not specified, RETURNING *
func (b *Builder) Returns(cols ...string) *Builder {
	b.sqb.Returns(cols...)
	return b
}

func (b *Builder) Group(cols ...string) *Builder {
	b.sqb.Group(cols...)
	return b
}

func (b *Builder) Order(order string, desc ...bool) *Builder {
	b.sqb.Order(order, desc...)
	return b
}

func (b *Builder) Orders(order string, defaults ...string) *Builder {
	b.sqb.Orders(order, defaults...)
	return b
}

func (b *Builder) Offset(offset int) *Builder {
	b.sqb.Offset(offset)
	return b
}

func (b *Builder) Limit(limit int) *Builder {
	b.sqb.Limit(limit)
	return b
}

// StructSelect add columns for Struct.
// Example:
//
//	type User struct {
//		ID    int64
//		Name  string
//		Value string
//	}
//	u := &User{}
//	sqb.StructSelect(u) // SELECT id, name, value FROM ...
func (b *Builder) StructSelect(a any, omits ...string) *Builder {
	sm := b.mpr.TypeMap(reflect.TypeOf(a))
	for _, fi := range sm.Index {
		if isIgnoredField(fi, omits...) {
			continue
		}
		b.sqb.Select(fi.Name)
	}
	return b
}

// StructPrefixSelect add columns for Struct.
// Example:
//
//	type User struct {
//		ID    int64
//		Name  string
//		Value string
//	}
//	u := &User{}
//	sqb.StructPrefixSelect(u, "u.") // SELECT u.id, u.name, value FROM ...
func (b *Builder) StructPrefixSelect(a any, prefix string, omits ...string) *Builder {
	sm := b.mpr.TypeMap(reflect.TypeOf(a))
	for _, fi := range sm.Index {
		if isIgnoredField(fi, omits...) {
			continue
		}
		b.sqb.Select(prefix + fi.Name)
	}
	return b
}

// StructNames add named columns and values for Struct.
// Example:
//
//	type User struct {
//		ID    int64
//		Name  string
//		Value string
//	}
//	u := &User{}
//	sqb.Insert("users").StructNames(u) // INSERT INTO users (id, name, value) VALUES (:id, :name, :value)
//	sqb.Update("users").StructNames(u, "id").Where("id = :id") // UPDATE users SET name = :name, value = :value WHERE id = :id
func (b *Builder) StructNames(a any, omits ...string) *Builder {
	sm := b.mpr.TypeMap(reflect.TypeOf(a))
	for _, fi := range sm.Index {
		if isIgnoredField(fi, omits...) {
			continue
		}
		b.sqb.Name(fi.Name)
	}
	return b
}
