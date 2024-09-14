package sqlx

import (
	"reflect"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/sqx"
)

// Builder a simple sql builder
// NOTE: the arguments are strict to it's order
type Builder struct {
	mpr *ref.Mapper
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
	return b.sqb.SQL()
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

// Distinct set distinct keyword only for SELECT
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
		if !asg.Contains(omits, fi.Name) {
			b.sqb.Select(fi.Name)
		}
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
		if !asg.Contains(omits, fi.Name) {
			b.sqb.Select(prefix + fi.Name)
		}
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
		if !asg.Contains(omits, fi.Name) {
			b.sqb.Names(fi.Name)
		}
	}
	return b
}
