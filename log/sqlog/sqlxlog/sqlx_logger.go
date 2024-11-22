package sqlxlog

import (
	"errors"
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
)

type SqlxLogger struct {
	Logger          log.Logger
	DefaultSQLLevel log.Level
	ErrorSQLLevel   log.Level
	SelectSQLLevel  log.Level
	SlowSQLLevel    log.Level
	SlowThreshold   time.Duration
	TraceErrNoRows  bool
	GetErrLogLevel  func(error) log.Level
	GetSQLLogLevel  func(string) log.Level
}

func NewSqlxLogger(logger log.Logger, slowSQL time.Duration) *SqlxLogger {
	sl := &SqlxLogger{
		Logger:          logger,
		DefaultSQLLevel: log.LevelInfo,
		ErrorSQLLevel:   log.LevelError,
		SelectSQLLevel:  log.LevelDebug,
		SlowSQLLevel:    log.LevelWarn,
		SlowThreshold:   slowSQL,
	}

	return sl
}

func (sl *SqlxLogger) printf(lvl log.Level, msg string, data ...any) {
	if sl.Logger.IsLevelEnabled(lvl) {
		le := &log.Event{
			Name:  sl.Logger.GetName(),
			Props: sl.Logger.GetProps(),
			Level: lvl,
			Msg:   fmt.Sprintf(msg, data...),
			Time:  time.Now(),
		}
		le.CallerStop("/pango/sqx/sqlx/", sl.Logger.GetTraceLevel() >= lvl)

		sl.Logger.Write(le)
	}
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
				sl.printf(lvl, "%s [%v] %s", err, elapsed, sql)
			} else {
				sl.printf(lvl, "%s [%d: %v] %s", err, rows, elapsed, sql)
			}
		}
	case sl.SlowThreshold != 0 && elapsed > sl.SlowThreshold && sl.Logger.IsLevelEnabled(sl.SlowSQLLevel):
		if rows < 0 {
			sl.printf(sl.SlowSQLLevel, "SLOW >= %v [%v] %s", sl.SlowThreshold, elapsed, sql)
		} else {
			sl.printf(sl.SlowSQLLevel, "SLOW >= %v [%d: %v] %s", sl.SlowThreshold, rows, elapsed, sql)
		}
	default:
		lvl := sl.DefaultSQLLevel
		if f := sl.GetSQLLogLevel; f != nil {
			lvl = f(sql)
		} else if str.StartsWithFold(str.StripLeft(sql), "SELECT") {
			lvl = sl.SelectSQLLevel
		}
		if sl.Logger.IsLevelEnabled(lvl) {
			if rows < 0 {
				sl.printf(lvl, "[%v] %s", elapsed, sql)
			} else {
				sl.printf(lvl, "[%d: %v] %s", rows, elapsed, sql)
			}
		}
	}
}
