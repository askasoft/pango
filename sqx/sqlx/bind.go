package sqlx

import (
	"regexp"
	"strconv"
	"sync"

	"github.com/askasoft/pango/str"
)

type Binder int

// Binder types supported by Rebind, BindMap and BindStruct.
const (
	BindUnknown Binder = iota
	BindQuestion
	BindDollar
	BindColon
	BindAt
)

var (
	PlaceholderDollar = regexp.MustCompile(`\$(\d+)`)
	PlaceholderColon  = regexp.MustCompile(`:arg(\d+)`)
	PlaceholderAt     = regexp.MustCompile(`@p(\d+)`)
)

var binds sync.Map

func init() {
	defaultBinds := map[Binder][]string{
		BindDollar:   {"postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach"},
		BindQuestion: {"mysql", "sqlite3", "nrmysql", "nrsqlite3"},
		BindColon:    {"oci8", "ora", "goracle", "godror"},
		BindAt:       {"sqlserver"},
	}

	for bind, drivers := range defaultBinds {
		for _, driver := range drivers {
			BindDriver(driver, bind)
		}
	}
}

// GetBinder returns the binder for a given database given a drivername.
func GetBinder(driverName string) Binder {
	itype, ok := binds.Load(driverName)
	if !ok {
		return BindUnknown
	}
	return itype.(Binder)
}

// BindDriver sets the Binder for driverName to binder.
func BindDriver(driverName string, binder Binder) {
	binds.Store(driverName, binder)
}

// Placeholder returns the placeholder regexp
func (binder Binder) Placeholder() *regexp.Regexp {
	switch binder {
	case BindDollar:
		return PlaceholderDollar
	case BindColon:
		return PlaceholderColon
	case BindAt:
		return PlaceholderAt
	default:
		return nil
	}
}

// Rebind a query from the default binder (QUESTION) to the target binder.
func (binder Binder) Rebind(query string) string {
	switch binder {
	case BindQuestion, BindUnknown:
		return query
	}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	n := int64(0)

	i := str.IndexByte(query, '?')
	for ; i != -1; i = str.IndexByte(query, '?') {
		rqb = append(rqb, query[:i]...)

		switch binder {
		case BindDollar:
			rqb = append(rqb, '$')
		case BindColon:
			rqb = append(rqb, ':', 'a', 'r', 'g')
		case BindAt:
			rqb = append(rqb, '@', 'p')
		}

		n++
		rqb = strconv.AppendInt(rqb, n, 10)

		query = query[i+1:]
	}
	rqb = append(rqb, query...)

	return str.UnsafeString(rqb)
}
