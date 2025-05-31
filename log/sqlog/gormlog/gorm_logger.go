package gormlog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger                   log.Logger
	DefaultSQLLevel          log.Level
	ErrorSQLLevel            log.Level
	WriteSQLLevel            log.Level
	SlowSQLLevel             log.Level
	SlowThreshold            time.Duration
	TraceRecordNotFoundError bool
	ParameterizedQueries     bool
	GetErrLogLevel           func(error) log.Level
	GetSQLLogLevel           func(string) log.Level
}

func NewGormLogger(logger log.Logger, slowSQL time.Duration) *GormLogger {
	gl := &GormLogger{
		Logger:          logger,
		DefaultSQLLevel: log.LevelDebug,
		ErrorSQLLevel:   log.LevelError,
		WriteSQLLevel:   log.LevelInfo,
		SlowSQLLevel:    log.LevelWarn,
		SlowThreshold:   slowSQL,
	}

	return gl
}

// LogMode log mode
func (gl *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return gl
}

// Info print info
func (gl *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	gl.printf(log.LevelInfo, msg, data...)
}

// Warn print warn messages
func (gl *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	gl.printf(log.LevelWarn, msg, data...)
}

// Error print error messages
func (gl *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	gl.printf(log.LevelError, msg, data...)
}

func (gl *GormLogger) printf(lvl log.Level, msg string, data ...any) {
	if gl.Logger.IsLevelEnabled(lvl) {
		le := &log.Event{
			Name:  gl.Logger.GetName(),
			Props: gl.Logger.GetProps(),
			Level: lvl,
			Msg:   fmt.Sprintf(msg, data...),
			Time:  time.Now(),
		}
		le.CallerStop("gorm.io", gl.Logger.GetTraceLevel() >= lvl)

		gl.Logger.Write(le)
	}
}

func (gl *GormLogger) getSQLLogLevel(sql string) log.Level {
	sql = str.StripLeft(sql)
	if str.StartsWithFold(sql, "SELECT") {
		return gl.DefaultSQLLevel
	}
	return gl.WriteSQLLevel
}

// Trace print sql message
func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || gl.TraceRecordNotFoundError):
		lvl := gl.ErrorSQLLevel
		if f := gl.GetErrLogLevel; f != nil {
			lvl = f(err)
		}
		if gl.Logger.IsLevelEnabled(lvl) {
			sql, rows := fc()
			if rows < 0 {
				gl.printf(lvl, "%s [%v] %s", err, elapsed, sql)
			} else {
				gl.printf(lvl, "%s [%d: %v] %s", err, rows, elapsed, sql)
			}
		}
	case gl.SlowThreshold != 0 && elapsed > gl.SlowThreshold && gl.Logger.IsLevelEnabled(gl.SlowSQLLevel):
		sql, rows := fc()
		if rows < 0 {
			gl.printf(gl.SlowSQLLevel, "SLOW >= %v [%v] %s", gl.SlowThreshold, elapsed, sql)
		} else {
			gl.printf(gl.SlowSQLLevel, "SLOW >= %v [%d: %v] %s", gl.SlowThreshold, rows, elapsed, sql)
		}
	default:
		sql, rows := fc()

		f := gl.GetSQLLogLevel
		if f == nil {
			f = gl.getSQLLogLevel
		}

		lvl := f(sql)
		if gl.Logger.IsLevelEnabled(lvl) {
			if rows < 0 {
				gl.printf(lvl, "[%v] %s", elapsed, sql)
			} else {
				gl.printf(lvl, "[%d: %v] %s", rows, elapsed, sql)
			}

		}
	}
}

// Trace print sql message
func (gl *GormLogger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if gl.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
