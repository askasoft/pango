package sqx

import (
	"strconv"
	"strings"
	"sync"
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

// FIXME: this should be able to be tolerant of escaped ?'s in queries without
// losing much speed, and should be to avoid confusion.

// Rebind a query from the default bindtype (QUESTION) to the target bindtype.
func (binder Binder) Rebind(query string) string {
	switch binder {
	case BindQuestion, BindUnknown:
		return query
	}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		switch binder {
		case BindDollar:
			rqb = append(rqb, '$')
		case BindColon:
			rqb = append(rqb, ':', 'a', 'r', 'g')
		case BindAt:
			rqb = append(rqb, '@', 'p')
		}

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}
