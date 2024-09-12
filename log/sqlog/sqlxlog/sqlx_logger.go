package sqlxlog

import (
	"errors"
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx/sqlx"
)

type SqlxLogger struct {
	Logger         log.Logger
	Level          log.Level
	ErrorSQLLevel  log.Level
	SlowSQLLevel   log.Level
	SlowThreshold  time.Duration
	TraceErrNoRows bool
}

func NewSqlxLogger(logger log.Logger, slowSQL ...time.Duration) *SqlxLogger {
	sl := &SqlxLogger{
		Logger:        logger,
		Level:         log.LevelDebug,
		ErrorSQLLevel: log.LevelWarn,
		SlowSQLLevel:  log.LevelWarn,
		SlowThreshold: time.Second,
	}

	if len(slowSQL) > 0 {
		sl.SlowThreshold = slowSQL[0]
	}

	return sl
}

func (sl *SqlxLogger) printf(lvl log.Level, msg string, data ...any) {
	if sl.Logger.IsLevelEnabled(lvl) {
		le := log.Event{
			Logger: sl.Logger,
			Level:  lvl,
			Msg:    fmt.Sprintf(msg, data...),
			Time:   time.Now(),
		}
		le.CallerStop("/pango/sqx/sqlx/", sl.Logger.GetTraceLevel() >= lvl)

		sl.Logger.Write(le)
	}
}

// Trace print sql message
func (sl *SqlxLogger) Trace(begin time.Time, sql string, rows int64, err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil && (!errors.Is(err, sqlx.ErrNoRows) || sl.TraceErrNoRows) && sl.Logger.IsLevelEnabled(sl.ErrorSQLLevel):
		if rows < 0 {
			sl.printf(sl.ErrorSQLLevel, "%s [%v] %s", err, elapsed, sql)
		} else {
			sl.printf(sl.ErrorSQLLevel, "%s [%d: %v] %s", err, rows, elapsed, sql)
		}
	case sl.SlowThreshold != 0 && elapsed > sl.SlowThreshold && sl.Logger.IsLevelEnabled(sl.SlowSQLLevel):
		if rows < 0 {
			sl.printf(sl.SlowSQLLevel, "SLOW >= %v [%v] %s", sl.SlowThreshold, elapsed, sql)
		} else {
			sl.printf(sl.SlowSQLLevel, "SLOW >= %v [%d: %v] %s", sl.SlowThreshold, rows, elapsed, sql)
		}
	case sl.Logger.IsLevelEnabled(sl.Level):
		if rows < 0 {
			sl.printf(sl.Level, "[%v] %s", elapsed, sql)
		} else {
			sl.printf(sl.Level, "[%d: %v] %s", rows, elapsed, sql)
		}
	}
}
