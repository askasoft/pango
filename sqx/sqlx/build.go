package sqlx

import (
	"github.com/askasoft/pango/sqx"
)

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

func (b *Builder) Count(cols ...string) *Builder {
	b.sqb.Count(cols...)
	return b
}

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

func (b *Builder) Columns(cols ...string) *Builder {
	b.sqb.Columns(cols...)
	return b
}

func (b *Builder) From(tb string) *Builder {
	b.sqb.From(tb)
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

func (b *Builder) Where(q string, args ...any) *Builder {
	b.sqb.Where(q, args...)
	return b
}

func (b *Builder) Set(col string, args ...any) *Builder {
	b.sqb.Set(col, args...)
	return b
}

func (b *Builder) Into(col string, args ...any) *Builder {
	b.sqb.Into(col, args...)
	return b
}

func (b *Builder) In(col string, arg any) *Builder {
	b.sqb.In(col, arg)
	return b
}

func (b *Builder) NotIn(col string, arg any) *Builder {
	b.sqb.NotIn(col, arg)
	return b
}

func (b *Builder) Values(vals ...string) *Builder {
	b.sqb.Values(vals...)
	return b
}
