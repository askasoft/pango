package sqlx

import (
	"github.com/askasoft/pango/sqx"
)

// Builder a simple sql builder
// NOTE: the arguments are strict to it's order
type Builder struct {
	bid Binder
	sqb sqx.Builder
}

func (b *Builder) Reset() *Builder {
	b.sqb.Reset()
	return b
}

func (b *Builder) Build() (string, []any) {
	return b.SQL(), b.Params()
}

func (b *Builder) SQL() string {
	return b.bid.Rebind(b.sqb.SQL())
}

func (b *Builder) Params() []any {
	return b.sqb.Params()
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

// Select add distinct select columns
// if `cols` is not specified, default select "*"
func (b *Builder) SelectDistinct(cols ...string) *Builder {
	b.sqb.SelectDistinct(cols...)
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

func (b *Builder) From(tb string) *Builder {
	b.sqb.From(tb)
	return b
}

func (b *Builder) Columns(cols ...string) *Builder {
	b.sqb.Columns(cols...)
	return b
}

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

func (b *Builder) Where(q string, args ...any) *Builder {
	b.sqb.Where(q, args...)
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

// Names add named columns and values.
// Example:
//
//	sqb.Insert("a").Names("id", "name", "value") // INSERT INTO a (id, name) VALUES (:id, :name)
//	sqb.Update("a").Names("name", "value").Where("id = :id") // UPDATE a SET name = :name, value = :value WHERE id = :id
func (b *Builder) Names(cols ...string) *Builder {
	b.sqb.Names(cols...)
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

func (b *Builder) Order(order string, desc ...bool) *Builder {
	b.sqb.Order(order, desc...)
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
