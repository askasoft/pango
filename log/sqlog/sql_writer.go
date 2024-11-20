package sqlog

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/str"
)

// SQLWriter implements log Writer Interface and batch send log messages to database.
// Prepare:
//
//	CREATE TABLE logs (
//		id serial NOT NULL,
//		time timestamp with time zone NOT NULL,
//		level char(1) NOT NULL,
//		msg text NOT NULL,
//		file text NOT NULL,
//		line integer NOT NULL,
//		func text NOT NULL,
//		trace text NOT NULL);
//
// Driver: postgres
// Dsn: host=127.0.0.1 user=pango password=pango dbname=pango port=5432 sslmode=disable
// Statement: "INSERT INTO sqlogs (time, level, msg, file, line, func, trace) VALUES"
// Parameter: "%t %p %m %S %L %F %T"
type SQLWriter struct {
	log.BatchSupport
	log.FilterSupport

	Driver    string
	Dsn       string
	Statement string

	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration

	db     *sql.DB      // database
	binder sqx.Binder   // placeholder binder
	affs   []argFmtFunc // argument functions
	stmb   str.Builder  // statement buffer
	args   []any        // argument buffer
}

type argFmtFunc func(le *log.Event) any

// SetConnMaxIdleTime set ConnMaxIdleTime
func (sw *SQLWriter) SetConnMaxIdleTime(duration string) error {
	td, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("SQLWriter: invalid ConnMaxIdleTime %q: %w", duration, err)
	}

	sw.ConnMaxIdleTime = td
	return nil
}

// SetConnMaxLifeTime set ConnMaxLifeTime
func (sw *SQLWriter) SetConnMaxLifeTime(duration string) error {
	td, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("SQLWriter: invalid ConnMaxLifeTime %q: %w", duration, err)
	}

	sw.ConnMaxLifeTime = td
	return nil
}

// SetParameter set sql statement parameter
// %t: event time
// %c: logger name
// %p: log level prefix
// %l: log level string
// %x{key}: logger property
// %S: caller source file name
// %L: caller source line number
// %F: caller function name
// %T: caller stack trace
// %m: message
// %M{msg}: custom message
func (sw *SQLWriter) SetParameter(format string) {
	affs := make([]argFmtFunc, 0, 10)

	for i := 0; i < len(format); i++ {
		c := format[i]
		if c != '%' {
			continue
		}

		i++
		if i >= len(format) {
			break
		}

		// symbol
		var p argFmtFunc
		switch format[i] {
		case 't':
			p = ffTime
		case 'c':
			p = ffName
		case 'p':
			p = ffPrefix
		case 'l':
			p = ffLevel
		case 'x':
			o := getOption(format, &i)
			if o != "" {
				p = fcProp(o)
			}
		case 'S':
			p = ffFile
		case 'L':
			p = ffLine
		case 'F':
			p = ffFunc
		case 'T':
			p = ffTrace
		case 'm':
			p = func(le *log.Event) any {
				return le.Msg
			}
		case 'M':
			o := getOption(format, &i)
			if o != "" {
				p = fcString(o)
			}
		}

		if p != nil {
			affs = append(affs, p)
		}
	}

	sw.affs = affs
}

func getOption(format string, i *int) string {
	o := format[*i+1:]
	if len(o) > 0 && o[0] == '{' {
		e := strings.IndexByte(o, '}')
		if e > 0 {
			*i += e + 1
			return o[1:e]
		}
	}
	return ""
}

func fcString(s string) argFmtFunc {
	return func(le *log.Event) any {
		return s
	}
}

func fcProp(key string) argFmtFunc {
	return func(le *log.Event) any {
		return le.Props[key]
	}
}

func ffName(le *log.Event) any {
	return le.Name
}

func ffTime(le *log.Event) any {
	return le.Time
}

func ffPrefix(le *log.Event) any {
	return le.Level.Prefix()
}

func ffLevel(le *log.Event) any {
	return le.Level.String()
}

func ffFile(le *log.Event) any {
	return le.File
}

func ffLine(le *log.Event) any {
	return le.Line
}

func ffFunc(le *log.Event) any {
	return le.Func
}

func ffTrace(le *log.Event) any {
	return le.Trace
}

// Write cache log message, flush if needed
func (sw *SQLWriter) Write(le *log.Event) {
	if sw.Reject(le) {
		return
	}

	sw.BatchWrite(le, sw.flush)
}

// Flush flush cached events
func (sw *SQLWriter) Flush() {
	sw.BatchFlush(sw.flush)
}

// Close flush and close the writer
func (sw *SQLWriter) Close() {
	sw.Flush()
}

func (sw *SQLWriter) flush(eb *log.EventBuffer) error {
	sw.initDB()

	sw.stmb.Reset()
	sw.args = sw.args[:0]

	n := 0
	sw.stmb.WriteString(sw.Statement)
	for i, it := 0, eb.Iterator(); it.Next(); i++ {
		le := it.Value()

		sw.stmb.WriteString(str.If(i == 0, " (", ",("))
		for j, f := range sw.affs {
			sw.args = append(sw.args, f(le))
			if j > 0 {
				sw.stmb.WriteRune(',')
			}
			n++
			sw.stmb.WriteString(sw.binder.Placeholder(n))
		}
		sw.stmb.WriteRune(')')
	}

	sql := sw.stmb.String()
	_, err := sw.db.Exec(sql, sw.args...)
	return err
}

func (sw *SQLWriter) initDB() {
	if sw.db == nil {
		db, err := sql.Open(sw.Driver, sw.Dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SQLWrite(%q, %q): %v", sw.Driver, sw.Dsn, err)
			return
		}

		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		db.SetConnMaxIdleTime(sw.ConnMaxIdleTime)
		db.SetConnMaxLifetime(sw.ConnMaxLifeTime)

		sw.db = db
		sw.binder = sqx.GetBinder(sw.Driver)
	}
}

func init() {
	log.RegisterWriter("sql", func() log.Writer {
		return &SQLWriter{}
	})
}
