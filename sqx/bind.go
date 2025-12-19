package sqx

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
)

type Rebind interface {
	Rebind(string) string
}

type Binder int

const (
	BindUnknown Binder = iota
	BindQuestion
	BindDollar
	BindColon
	BindAt
)

var (
	reDollar = regexp.MustCompile(`\$(\d+)`)
	reColon  = regexp.MustCompile(`:arg(\d+)`)
	reAt     = regexp.MustCompile(`@p(\d+)`)
)

var binders = map[string]Binder{
	"postgres":         BindDollar,
	"pgx":              BindDollar,
	"pq-timeouts":      BindDollar,
	"cloudsqlpostgres": BindDollar,
	"ql":               BindDollar,
	"nrpostgres":       BindDollar,
	"cockroach":        BindDollar,
	"mysql":            BindQuestion,
	"sqlite3":          BindQuestion,
	"nrmysql":          BindQuestion,
	"nrsqlite3":        BindQuestion,
	"oci8":             BindColon,
	"ora":              BindColon,
	"goracle":          BindColon,
	"godror":           BindColon,
	"sqlserver":        BindAt,
	"azuresql":         BindAt,
}

// GetBinder returns the binder for a given database given a drivername.
func GetBinder(driverName string) Binder {
	binder, ok := binders[driverName]
	if !ok {
		return BindUnknown
	}
	return binder
}

// BindDriver sets the Binder for driverName to binder.
func BindDriver(driverName string, binder Binder) {
	nbs := make(map[string]Binder)
	mag.Copy(nbs, binders)

	nbs[driverName] = binder
	binders = nbs
}

// Rebind a SQL from the default binder (QUESTION) to the target binder.
func (b Binder) Rebind(sql string) string {
	switch b {
	case BindQuestion, BindUnknown:
		return sql
	}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(sql)+10)

	n := 0
	for i := str.IndexByte(sql, '?'); i >= 0; i = str.IndexByte(sql, '?') {
		rqb = append(rqb, sql[:i]...)

		n++
		rqb = b.append(rqb, n)

		sql = sql[i+1:]
	}
	rqb = append(rqb, sql...)

	return str.UnsafeString(rqb)
}

func (b Binder) append(q []byte, n int) []byte {
	switch b {
	case BindDollar:
		q = append(q, '$')
	case BindColon:
		q = append(q, ':', 'a', 'r', 'g')
	case BindAt:
		q = append(q, '@', 'p')
	default:
		n = 0
	}

	if n > 0 {
		q = strconv.AppendInt(q, int64(n), 10)
	} else {
		q = append(q, '?')
	}
	return q
}

// Placeholder generate a place holder mark with No. n.
func (b Binder) Placeholder(n int) string {
	switch b {
	case BindDollar:
		return fmt.Sprintf("$%d", n)
	case BindColon:
		return fmt.Sprintf(":arg%d", n)
	case BindAt:
		return fmt.Sprintf("@p%d", n)
	default:
		return "?"
	}
}

func (b Binder) matcher() *regexp.Regexp {
	switch b {
	case BindDollar:
		return reDollar
	case BindColon:
		return reColon
	case BindAt:
		return reAt
	default:
		return nil
	}
}

// Explain generate SQL string with given parameters.
// The generated SQL is expected to be used in logger, execute it might introduce a SQL injection vulnerability.
// If a string or binary argument's length exceeds `maxArgLength`, it will be cut and append with '...'.
func (b Binder) Explain(sql string, maxArgLength int, args ...any) string {
	if len(args) == 0 {
		return sql
	}

	vars := make([]string, len(args))
	for i, v := range args {
		vars[i] = b.convert(v, maxArgLength)
	}

	rep := b.matcher()

	if rep == nil {
		var idx int
		var sb strings.Builder

		for _, v := range str.UnsafeBytes(sql) {
			if v == '?' {
				if len(vars) > idx {
					sb.WriteString(vars[idx])
					idx++
					continue
				}
			}
			sb.WriteByte(v)
		}

		return sb.String()
	}

	sql = rep.ReplaceAllStringFunc(sql, func(p string) string {
		i := str.IndexAny(p, "123456789")
		if i >= 0 {
			n, _ := strconv.Atoi(p[i:])

			// position var start from 1 ($1, $2)
			if n > 0 && n <= len(vars) {
				return vars[n-1]
			}
		}
		return p
	})

	return sql
}

const (
	tmFormat = "2006-01-02 15:04:05.999"
	zeroTime = "0000-00-00 00:00:00"
	nullStr  = "NULL"
)

// A list of Go types that should be converted to SQL primitives
var convertibleTypes = []reflect.Type{ref.TypeTime, ref.TypeBool, ref.TypeBytes}

func (b Binder) toLiteralString(s string, limit int) string {
	if limit > 0 {
		s = str.Ellipsis(s, limit)
	}
	return "'" + str.ReplaceAll(s, "'", "''") + "'"
}

func (b Binder) toLiteralBinary(bs []byte, limit int) string {
	switch {
	case len(bs) == 0:
		return "''"
	case limit > 0 && len(bs) > limit:
		return "\\x" + hex.EncodeToString(bs[:limit]) + "...'"
	default:
		return "\\x" + hex.EncodeToString(bs) + "'"
	}
}

func (b Binder) convert(v any, limit int) string {
	if ref.IsNil(v) {
		return nullStr
	}

	switch v := v.(type) {
	case string:
		return b.toLiteralString(v, limit)
	case []byte:
		if s := str.UnsafeString(v); str.IsUTFPrintable(s) {
			return b.toLiteralString(s, limit)
		}
		return b.toLiteralBinary(v, limit)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		s, _ := cas.ToString(v)
		return s
	case time.Time:
		if v.IsZero() {
			return "'" + zeroTime + "'"
		}
		return "'" + v.Local().Format(tmFormat) + "'"
	case *time.Time:
		if v.IsZero() {
			return "'" + zeroTime + "'"
		}
		return "'" + v.Local().Format(tmFormat) + "'"
	case driver.Valuer:
		r, _ := v.Value()
		return b.convert(r, limit)
	case fmt.Stringer:
		return b.toLiteralString(v.String(), limit)
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Pointer && !rv.IsZero() {
			return b.convert(reflect.Indirect(rv).Interface(), limit)
		}
		for _, t := range convertibleTypes {
			if rv.Type().ConvertibleTo(t) {
				return b.convert(rv.Convert(t).Interface(), limit)
			}
		}
		return b.toLiteralString(fmt.Sprint(v), limit)
	}
}
