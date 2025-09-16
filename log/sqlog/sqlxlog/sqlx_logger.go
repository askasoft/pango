package sqlxlog

import (
	"errors"
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

type SqlxLogger struct {
	Logger          log.Logger
	DefaultSQLLevel log.Level
	ErrorSQLLevel   log.Level
	WriteSQLLevel   log.Level
	SlowSQLLevel    log.Level
	SlowThreshold   time.Duration
	TraceErrNoRows  bool
	GetErrLogLevel  func(error) log.Level
	GetSQLLogLevel  func(string) log.Level
}

func NewSqlxLogger(logger log.Logger, slowSQL time.Duration) *SqlxLogger {
	sl := &SqlxLogger{
		Logger:          logger,
		DefaultSQLLevel: log.LevelDebug,
		ErrorSQLLevel:   log.LevelError,
		WriteSQLLevel:   log.LevelInfo,
		SlowSQLLevel:    log.LevelWarn,
		SlowThreshold:   slowSQL,
	}

	return sl
}

func (sl *SqlxLogger) printf(lvl log.Level, msg string, data ...any) {
	if sl.Logger.IsLevelEnabled(lvl) {
		le := log.NewEvent(sl.Logger, lvl, fmt.Sprintf(msg, data...))
		le.CallerStop("/sqx/sqlx/", sl.Logger.GetTraceLevel() >= lvl)
		sl.Logger.Write(le)
	}
}

func (sl *SqlxLogger) getSQLLogLevel(sql string) log.Level {
	sql = str.StripLeft(sql)
	if str.StartsWithFold(sql, "SELECT") || str.StartsWithFold(sql, "Prepare") ||
		str.StartsWithFold(sql, "Begin") || str.StartsWithFold(sql, "Commit") {
		return sl.DefaultSQLLevel
	}
	return sl.WriteSQLLevel
}

// Trace print sql message
func (sl *SqlxLogger) Trace(begin time.Time, sql string, rows int64, err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil && (sl.TraceErrNoRows || !errors.Is(err, sqlx.ErrNoRows)):
		lvl := sl.ErrorSQLLevel
		if f := sl.GetErrLogLevel; f != nil {
			lvl = f(err)
		}
		if sl.Logger.IsLevelEnabled(lvl) {
			if rows < 0 {
				sl.printf(lvl, "%s [%s] %s", err, tmu.HumanDuration(elapsed), sql)
			} else {
				sl.printf(lvl, "%s [%d: %s] %s", err, rows, tmu.HumanDuration(elapsed), sql)
			}
		}
	case sl.SlowThreshold != 0 && elapsed > sl.SlowThreshold && sl.Logger.IsLevelEnabled(sl.SlowSQLLevel):
		if rows < 0 {
			sl.printf(sl.SlowSQLLevel, "SLOW >= %s [%s] %s", tmu.HumanDuration(sl.SlowThreshold), tmu.HumanDuration(elapsed), sql)
		} else {
			sl.printf(sl.SlowSQLLevel, "SLOW >= %s [%d: %s] %s", tmu.HumanDuration(sl.SlowThreshold), rows, tmu.HumanDuration(elapsed), sql)
		}
	default:
		f := sl.GetSQLLogLevel
		if f == nil {
			f = sl.getSQLLogLevel
		}

		lvl := f(sql)
		if sl.Logger.IsLevelEnabled(lvl) {
			if rows < 0 {
				sl.printf(lvl, "[%s] %s", tmu.HumanDuration(elapsed), sql)
			} else {
				sl.printf(lvl, "[%d: %s] %s", rows, tmu.HumanDuration(elapsed), sql)
			}
		}
	}
}
