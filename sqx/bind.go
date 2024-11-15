package sqx

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/mag"
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

// Rebind a query from the default binder (QUESTION) to the target binder.
func (binder Binder) Rebind(query string) string {
	switch binder {
	case BindQuestion, BindUnknown:
		return query
	}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	n := 0
	for i := str.IndexByte(query, '?'); i >= 0; i = str.IndexByte(query, '?') {
		rqb = append(rqb, query[:i]...)

		n++
		rqb = binder.Append(rqb, n)

		query = query[i+1:]
	}
	rqb = append(rqb, query...)

	return str.UnsafeString(rqb)
}

func (binder Binder) Append(q []byte, n int) []byte {
	switch binder {
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

func (binder Binder) Placeholder(n int) string {
	switch binder {
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

func (binder Binder) matcher() *regexp.Regexp {
	switch binder {
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

// Explain generate SQL string with given parameters, the generated SQL is expected to be used in logger, execute it might introduce a SQL injection vulnerability
func (binder Binder) Explain(sql string, args ...any) string {
	if len(args) == 0 {
		return sql
	}

	vars := make([]string, len(args))
	for i, v := range args {
		vars[i] = convert(v)
	}

	rep := binder.matcher()

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
	tmFmtWithMS = "2006-01-02 15:04:05.999"
	tmFmtZero   = "0000-00-00 00:00:00"
	nullStr     = "NULL"
)

// A list of Go types that should be converted to SQL primitives
var convertibleTypes = []reflect.Type{reflect.TypeOf(time.Time{}), reflect.TypeOf(false), reflect.TypeOf([]byte{})}

func convert(v any) string {
	switch v := v.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case []byte:
		if s := str.UnsafeString(v); str.IsUTFPrintable(s) {
			return "'" + strings.ReplaceAll(s, "'", "''") + "'"
		}
		return "'" + "<binary>" + "'"
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		s, _ := cas.ToString(v)
		return s
	case time.Time:
		if v.IsZero() {
			return "'" + tmFmtZero + "'"
		}
		return "'" + v.Format(tmFmtWithMS) + "'"
	case *time.Time:
		if v != nil {
			if v.IsZero() {
				return "'" + tmFmtZero + "'"
			}
			return "'" + v.Format(tmFmtWithMS) + "'"
		}
		return nullStr
	case driver.Valuer:
		rv := reflect.ValueOf(v)
		if v != nil && rv.IsValid() && ((rv.Kind() == reflect.Ptr && !rv.IsNil()) || rv.Kind() != reflect.Ptr) {
			r, _ := v.Value()
			return convert(r)
		}
		return nullStr
	case fmt.Stringer:
		rv := reflect.ValueOf(v)
		if v != nil && rv.IsValid() && ((rv.Kind() == reflect.Ptr && !rv.IsNil()) || rv.Kind() != reflect.Ptr) {
			return "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
		}
		return nullStr
	default:
		rv := reflect.ValueOf(v)
		if v == nil || !rv.IsValid() || rv.Kind() == reflect.Ptr && rv.IsNil() {
			return nullStr
		}
		if valuer, ok := v.(driver.Valuer); ok {
			v, _ = valuer.Value()
			return convert(v)
		}
		if rv.Kind() == reflect.Ptr && !rv.IsZero() {
			return convert(reflect.Indirect(rv).Interface())
		}
		for _, t := range convertibleTypes {
			if rv.Type().ConvertibleTo(t) {
				return convert(rv.Convert(t).Interface())
			}
		}
		return "'" + strings.ReplaceAll(fmt.Sprint(v), "'", "''") + "'"
	}
}
