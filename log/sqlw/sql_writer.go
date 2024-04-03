package sqlw

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/askasoft/pango/log"
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
// Statement: "INSERT INTO sqlogs (time, level, msg, file, line, func, trace) VALUES ($1, $2, $3, $4, $5, $6, $7)"
// Parameter: "%t %p %m %S %L %F %T"
type SQLWriter struct {
	log.LogFilter
	log.BatchWriter

	Driver    string
	Dsn       string
	Statement string

	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration

	db *sql.DB
	ps []argfunc
}

type argfunc func(le *log.Event) any

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
	ps := make([]argfunc, 0, 10)

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
		var p argfunc
		switch format[i] {
		case 't':
			p = func(le *log.Event) any {
				return le.When
			}
		case 'c':
			p = func(le *log.Event) any {
				return le.Logger.GetName()
			}
		case 'p':
			p = func(le *log.Event) any {
				return le.Level.Prefix()
			}
		case 'l':
			p = func(le *log.Event) any {
				return le.Level.String()
			}
		case 'x':
			o := getOption(format, &i)
			if o != "" {
				p = fcProp(o)
			}
		case 'S':
			p = func(le *log.Event) any {
				return le.File
			}
		case 'L':
			p = func(le *log.Event) any {
				return le.Line
			}
		case 'F':
			p = func(le *log.Event) any {
				return le.Func
			}
		case 'T':
			p = func(le *log.Event) any {
				return le.Trace
			}
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
			ps = append(ps, p)
		}
	}

	sw.ps = ps
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

// Write cache log message, flush if needed
func (sw *SQLWriter) Write(le *log.Event) error {
	if sw.Reject(le) {
		return nil
	}

	if sw.BatchCount > 1 {
		sw.InitBuffer()
		sw.EventBuffer.Push(le)

		if sw.ShouldFlush(le) {
			if err := sw.flush(); err != nil {
				return err
			}
		}
		return nil
	}

	sw.initDB()
	return sw.write(le)
}

func (sw *SQLWriter) flush() error {
	sw.initDB()

	for le, ok := sw.EventBuffer.Peek(); ok; le, ok = sw.EventBuffer.Peek() {
		if err := sw.write(le); err != nil {
			return err
		}
		_, _ = sw.EventBuffer.Poll()
	}

	return nil
}

// Flush flush cached events
func (sw *SQLWriter) Flush() {
	if sw.EventBuffer == nil || sw.EventBuffer.IsEmpty() {
		return
	}

	if err := sw.flush(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

// Close flush and close the writer
func (sw *SQLWriter) Close() {
	sw.Flush()
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
	}
}

func (sw *SQLWriter) write(le *log.Event) error {
	args := make([]any, len(sw.ps))
	for i, f := range sw.ps {
		args[i] = f(le)
	}

	_, err := sw.db.Exec(sw.Statement, args...)
	return err
}

func fcString(s string) argfunc {
	return func(le *log.Event) any {
		return s
	}
}

func fcProp(key string) argfunc {
	return func(le *log.Event) any {
		return le.Logger.GetProp(key)
	}
}

func init() {
	log.RegisterWriter("sql", func() log.Writer {
		return &SQLWriter{}
	})
}
