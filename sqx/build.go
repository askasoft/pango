package sqx

import (
	"reflect"
	"strings"

	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
)

// Builder a simple sql builder
// NOTE: the arguments are strict to it's order
type Builder struct {
	Quoter   Quoter
	command  sqlcmd
	distinct bool
	table    string
	columns  []string
	values   []string
	joins    []string
	wheres   []string
	orders   []string
	returns  []string
	params   []any
	offset   int
	limit    int
}

func (b *Builder) Reset() *Builder {
	b.command = 0
	b.distinct = false
	b.table = ""
	b.columns = b.columns[:0]
	b.values = b.values[:0]
	b.joins = b.joins[:0]
	b.wheres = b.wheres[:0]
	b.orders = b.orders[:0]
	b.returns = b.returns[:0]
	b.params = b.params[:0]
	b.offset = 0
	b.limit = 0

	return b
}

func (b *Builder) Build() (string, []any) {
	return b.SQL(), b.Params()
}

func (b *Builder) Params() []any {
	return b.params
}

func (b *Builder) SQL() string {
	switch b.command {
	case cselect:
		return b.buildSelect()
	case cdelete:
		return b.buildDelete()
	case cinsert:
		return b.buildInsert()
	case cupdate:
		return b.buildUpdate()
	default:
		return ""
	}
}

// Count shortcut for SELECT COUNT(*)
func (b *Builder) Count(cols ...string) *Builder {
	b.command = cselect
	if len(cols) == 0 {
		return b.Columns("COUNT(*)")
	}
	return b.Columns("COUNT(" + str.Join(cols, ", ") + ")")
}

// Count shortcut for SELECT COUNT(distinct *)
func (b *Builder) CountDistinct(cols ...string) *Builder {
	b.command = cselect
	if len(cols) == 0 {
		return b.Columns("COUNT(distinct *)")
	}
	return b.Columns("COUNT(distinct " + str.Join(cols, ", ") + ")")
}

// Select add select columns
// if `cols` is not specified, default select "*"
func (b *Builder) Select(cols ...string) *Builder {
	b.command = cselect
	if len(cols) == 0 {
		return b.Columns("*")
	}
	return b.Columns(cols...)
}

// Select add distinct select columns
// if `cols` is not specified, default select "*"
func (b *Builder) SelectDistinct(cols ...string) *Builder {
	b.distinct = true
	return b.Select(cols...)
}

func (b *Builder) Delete(tb string) *Builder {
	b.command = cdelete
	b.table = tb
	return b
}

func (b *Builder) Insert(tb string) *Builder {
	b.command = cinsert
	b.table = tb
	return b
}

func (b *Builder) Update(tb string) *Builder {
	b.command = cupdate
	b.table = tb
	return b
}

func (b *Builder) From(tb string) *Builder {
	b.table = tb
	return b
}

func (b *Builder) Columns(cols ...string) *Builder {
	b.columns = append(b.columns, cols...)
	return b
}

func (b *Builder) Values(vals ...string) *Builder {
	b.values = append(b.values, vals...)
	return b
}

// Join specify Join query and conditions
//
//	sqb.Join("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "abc@example.org")
func (b *Builder) Join(query string, args ...any) *Builder {
	b.joins = append(b.joins, query)
	b.params = append(b.params, args...)
	return b
}

func (b *Builder) Where(q string, args ...any) *Builder {
	b.wheres = append(b.wheres, q)
	b.params = append(b.params, args...)
	return b
}

func (b *Builder) Setx(col string, val string, args ...any) *Builder {
	b.columns = append(b.columns, col)
	b.values = append(b.values, val)
	b.params = append(b.params, args...)
	return b
}

func (b *Builder) Setc(col string, arg any) *Builder {
	return b.Setx(col, "?", arg)
}

// Names add named columns and values.
// Example:
//
//	sqb.Insert("a").Names("id", "name", "value") // INSERT INTO a (id, name) VALUES (:id, :name)
//	sqb.Update("a").Names("name", "value").Where("id = :id") // UPDATE a SET name = :name, value = :value WHERE id = :id
func (b *Builder) Names(cols ...string) *Builder {
	for _, col := range cols {
		b.columns = append(b.columns, col)
		b.values = append(b.values, ":"+col)
	}
	return b
}

func (b *Builder) In(col string, val any) *Builder {
	return b.in("IN", col, val)
}

func (b *Builder) NotIn(col string, val any) *Builder {
	return b.in("NOT IN", col, val)
}

func (b *Builder) in(op, col string, val any) *Builder {
	sql, args := in(op, col, val)
	b.wheres = append(b.wheres, sql)
	b.params = append(b.params, args...)
	return b
}

func (b *Builder) Order(order string, desc ...bool) *Builder {
	if len(desc) > 0 {
		order += str.If(desc[0], " DESC", " ASC")
	}
	b.orders = append(b.orders, order)
	return b
}

func (b *Builder) Offset(offset int) *Builder {
	b.offset = offset
	return b
}

func (b *Builder) Limit(limit int) *Builder {
	b.limit = limit
	return b
}

func (b *Builder) Returns(cols ...string) *Builder {
	if len(cols) == 0 {
		b.returns = append(b.returns, "*")
	} else {
		b.returns = append(b.returns, cols...)
	}
	return b
}

func (b *Builder) buildSelect() string {
	sb := &strings.Builder{}

	sb.WriteString("SELECT ")
	if b.distinct {
		sb.WriteString("DISTINCT ")
	}
	for i, col := range b.columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(b.Quoter.Quote(col))
	}
	sb.WriteString(" FROM ")
	sb.WriteString(b.Quoter.Quote(b.table))

	for _, j := range b.joins {
		sb.WriteByte(' ')
		sb.WriteString(j)
	}

	b.appendWhere(sb)

	for i, o := range b.orders {
		sb.WriteString(str.If(i == 0, " ORDER BY ", ", "))
		sb.WriteString(o)
	}

	if b.limit > 0 {
		sb.WriteString(" LIMIT ")
		sb.WriteString(num.Itoa(b.limit))
	}
	if b.offset > 0 {
		sb.WriteString(" OFFSET ")
		sb.WriteString(num.Itoa(b.offset))
	}
	return sb.String()
}

func (b *Builder) buildUpdate() string {
	sb := &strings.Builder{}

	sb.WriteString("UPDATE ")
	sb.WriteString(b.Quoter.Quote(b.table))
	sb.WriteString(" SET ")

	for i, col := range b.columns {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(b.Quoter.Quote(col))
		sb.WriteString(" = ")

		if i < len(b.values) {
			sb.WriteString(b.values[i])
		} else {
			sb.WriteByte('?')
		}
	}

	b.appendWhere(sb)
	b.appendRetuning(sb)

	return sb.String()
}

func (b *Builder) buildInsert() string {
	sb := &strings.Builder{}

	sb.WriteString("INSERT INTO ")
	sb.WriteString(b.Quoter.Quote(b.table))
	if len(b.columns) > 0 {
		sb.WriteString(" (")
		for i, col := range b.columns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(b.Quoter.Quote(col))
		}
		sb.WriteString(")")
	}

	sb.WriteString(" VALUES (")
	for i, val := range b.values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(val)
	}
	sb.WriteString(")")

	b.appendRetuning(sb)

	return sb.String()
}

func (b *Builder) buildDelete() string {
	sb := &strings.Builder{}

	sb.WriteString("DELETE FROM ")
	sb.WriteString(b.Quoter.Quote(b.table))

	b.appendWhere(sb)
	b.appendRetuning(sb)

	return sb.String()
}

func (b *Builder) appendWhere(sb *strings.Builder) {
	for i, w := range b.wheres {
		sb.WriteString(str.If(i == 0, " WHERE ", " AND "))
		sb.WriteString(w)
	}
}

func (b *Builder) appendRetuning(sb *strings.Builder) {
	if len(b.returns) > 0 {
		sb.WriteString(" RETUNING ")
		for i, col := range b.returns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(b.Quoter.Quote(col))
		}
	}
}

type sqlcmd int

const (
	cselect sqlcmd = 'S'
	cinsert sqlcmd = 'I'
	cdelete sqlcmd = 'D'
	cupdate sqlcmd = 'U'
)

func (c sqlcmd) String() string {
	switch c {
	case cselect:
		return "SELECT"
	case cinsert:
		return "INSERT"
	case cdelete:
		return "DELETE"
	case cupdate:
		return "UPDATE"
	default:
		return "UNKNOWN"
	}
}

func In(col string, val any) (sql string, args []any) {
	return in("IN", col, val)
}

func NotIn(col string, val any) (sql string, args []any) {
	return in("NOT IN", col, val)
}

func in(op, col string, val any) (sql string, args []any) {
	if v, ok := asSliceForIn(val); ok {
		z := v.Len()

		qs := str.Repeat("?,", z)
		sql = col + " " + op + " (" + qs[:len(qs)-1] + ")"
		args = appendReflectSlice(args, v, z)
		return
	}

	sql = col + " " + op + " (?)"
	args = append(args, val)
	return
}

func asSliceForIn(i any) (v reflect.Value, ok bool) {
	if i == nil {
		return reflect.Value{}, false
	}

	v = reflect.ValueOf(i)
	t := ref.DerefType(v.Type())

	// Only expand slices
	if t.Kind() != reflect.Slice {
		return reflect.Value{}, false
	}

	// []byte is a driver.Value type so it should not be expanded
	if t == reflect.TypeOf([]byte{}) {
		return reflect.Value{}, false

	}

	return v, true
}

func appendReflectSlice(args []any, v reflect.Value, vlen int) []any {
	switch val := v.Interface().(type) {
	case []any:
		args = append(args, val...)
	case []int:
		for i := range val {
			args = append(args, val[i])
		}
	case []int32:
		for i := range val {
			args = append(args, val[i])
		}
	case []int64:
		for i := range val {
			args = append(args, val[i])
		}
	case []string:
		for i := range val {
			args = append(args, val[i])
		}
	default:
		for si := 0; si < vlen; si++ {
			args = append(args, v.Index(si).Interface())
		}
	}

	return args
}

func Question(n int) string {
	sb := &strings.Builder{}
	sb.Grow(n * 2)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('?')
	}
	return sb.String()
}

func Questions(n int) []string {
	qs := make([]string, n)
	for i := 0; i < n; i++ {
		qs[i] = "?"
	}
	return qs
}
