package sqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/str"
	"gorm.io/gorm/utils"
)

type tracer struct {
	Bind  Binder
	Trace Trace
}

func (t *tracer) ExplainSQL(sql string, args ...any) string {
	if len(args) == 0 {
		return sql
	}

	return ExplainSQL(sql, t.Bind.Placeholder(), args...)
}

func (t *tracer) TracePing(pr Pinger) error {
	start := time.Now()
	err := pr.Ping()
	if t.Trace != nil {
		t.Trace(start, "Ping()", -1, err)
	}
	return err
}

func (t *tracer) TracePingContext(ctx context.Context, pr ContextPinger) error {
	start := time.Now()
	err := pr.PingContext(ctx)
	if t.Trace != nil {
		t.Trace(start, "PingContext()", -1, err)
	}
	return err
}

func (t *tracer) TraceQuery(qr Queryer, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := qr.Query(query, args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, err)
	}
	return rows, err
}

func (t *tracer) TraceQueryContext(ctx context.Context, qr ContextQueryer, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := qr.QueryContext(ctx, query, args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, err)
	}
	return rows, err
}

func (t *tracer) TraceExec(er Execer, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	sqr, err := er.Exec(query, args...)
	if t.Trace != nil {
		cnt, _ := sqr.RowsAffected()
		t.Trace(start, t.ExplainSQL(query, args...), cnt, err)
	}
	return sqr, err
}

func (t *tracer) TraceExecContext(ctx context.Context, er ContextExecer, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	sqr, err := er.ExecContext(ctx, query, args...)
	if t.Trace != nil {
		cnt, _ := sqr.RowsAffected()
		t.Trace(start, t.ExplainSQL(query, args...), cnt, err)
	}
	return sqr, err
}

const (
	tmFmtWithMS = "2006-01-02 15:04:05.999"
	tmFmtZero   = "0000-00-00 00:00:00"
	nullStr     = "NULL"
)

// A list of Go types that should be converted to SQL primitives
var convertibleTypes = []reflect.Type{reflect.TypeOf(time.Time{}), reflect.TypeOf(false), reflect.TypeOf([]byte{})}

func isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func convert(v any) string {
	switch v := v.(type) {
	case bool:
		return strconv.FormatBool(v)
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
		reflectValue := reflect.ValueOf(v)
		if v != nil && reflectValue.IsValid() && ((reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil()) || reflectValue.Kind() != reflect.Ptr) {
			r, _ := v.Value()
			return convert(r)
		}
		return nullStr
	case fmt.Stringer:
		reflectValue := reflect.ValueOf(v)
		switch reflectValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fmt.Sprintf("%d", reflectValue.Interface())
		case reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%.6f", reflectValue.Interface())
		case reflect.Bool:
			return fmt.Sprintf("%t", reflectValue.Interface())
		case reflect.String:
			return "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
		default:
			if v != nil && reflectValue.IsValid() && ((reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil()) || reflectValue.Kind() != reflect.Ptr) {
				return "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
			}
			return nullStr
		}
	case []byte:
		if s := str.UnsafeString(v); str.IsUTFPrintable(s) {
			return "'" + strings.ReplaceAll(s, "'", "''") + "'"
		}
		return "'" + "<binary>" + "'"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return utils.ToString(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
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
		if isNumeric(rv.Kind()) {
			if rv.CanInt() || rv.CanUint() {
				return fmt.Sprintf("%d", rv.Interface())
			}
			return fmt.Sprintf("%.6f", rv.Interface())
		}
		for _, t := range convertibleTypes {
			if rv.Type().ConvertibleTo(t) {
				return convert(rv.Convert(t).Interface())
			}
		}
		return "'" + strings.ReplaceAll(fmt.Sprint(v), "'", "''") + "'"
	}
}

// ExplainSQL generate SQL string with given parameters, the generated SQL is expected to be used in logger, execute it might introduce a SQL injection vulnerability
func ExplainSQL(sql string, numericPlaceholder *regexp.Regexp, args ...any) string {
	vars := make([]string, len(args))
	for i, v := range args {
		vars[i] = convert(v)
	}

	if numericPlaceholder == nil {
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

	sql = numericPlaceholder.ReplaceAllStringFunc(sql, func(p string) string {
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
