package gormlog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger                   log.Logger
	Level                    log.Level
	ErrorSQLLevel            log.Level
	SlowSQLLevel             log.Level
	SlowThreshold            time.Duration
	TraceRecordNotFoundError bool
	ParameterizedQueries     bool
}

func NewGormLogger(logger log.Logger, slowSQL ...time.Duration) *GormLogger {
	gl := &GormLogger{
		Logger:        logger,
		Level:         log.LevelDebug,
		ErrorSQLLevel: log.LevelWarn,
		SlowSQLLevel:  log.LevelWarn,
		SlowThreshold: time.Second,
	}

	if len(slowSQL) > 0 {
		gl.SlowThreshold = slowSQL[0]
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
		le := log.Event{
			Logger: gl.Logger,
			Level:  lvl,
			Msg:    fmt.Sprintf(msg, data...),
			Time:   time.Now(),
		}
		le.CallerStop("gorm.io", gl.Logger.GetTraceLevel() >= lvl)

		gl.Logger.Write(le)
	}
}

// Trace print sql message
func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || gl.TraceRecordNotFoundError) && gl.Logger.IsLevelEnabled(gl.ErrorSQLLevel):
		sql, rows := fc()
		gl.printf(gl.ErrorSQLLevel, "%s [%d: %v] %s", err, rows, elapsed, sql)
	case gl.SlowThreshold != 0 && elapsed > gl.SlowThreshold && gl.Logger.IsLevelEnabled(gl.SlowSQLLevel):
		sql, rows := fc()
		gl.printf(gl.SlowSQLLevel, "SLOW >= %v [%d: %v] %s", gl.SlowThreshold, rows, elapsed, sql)
	case gl.Logger.IsLevelEnabled(gl.Level):
		sql, rows := fc()
		gl.printf(gl.Level, "[%d: %v] %s", rows, elapsed, sql)
	}
}

// Trace print sql message
func (gl *GormLogger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if gl.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
