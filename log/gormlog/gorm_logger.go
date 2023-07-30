package gormlog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger                   log.Logger
	SlowThreshold            time.Duration
	TraceRecordNotFoundError bool
	ParameterizedQueries     bool
}

// LogMode log mode
func (gl *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return gl
}

// Info print info
func (gl *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	gl.printf(log.LevelInfo, msg, data...)
}

// Warn print warn messages
func (gl *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	gl.printf(log.LevelWarn, msg, data...)
}

// Error print error messages
func (gl *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	gl.printf(log.LevelError, msg, data...)
}

func (gl *GormLogger) printf(lvl log.Level, msg string, data ...any) {
	if gl.Logger.IsLevelEnabled(lvl) {
		le := &log.Event{
			Logger: gl.Logger,
			Level:  lvl,
			Msg:    fmt.Sprintf(msg, data...),
			When:   time.Now(),
		}
		le.CallerStop("gorm.io", gl.Logger.GetTraceLevel() >= lvl)

		gl.Logger.Write(le)
	}
}

func getSqlAndRows(fc func() (string, int64)) (sql, rows string) {
	var n int64
	sql, n = fc()
	rows = str.If(n == -1, "-", num.Ltoa(n))
	return
}

// Trace print sql message
func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !gl.Logger.IsErrorEnabled() {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || gl.TraceRecordNotFoundError):
		sql, rows := getSqlAndRows(fc)
		gl.printf(log.LevelError, "%s [%.3fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case gl.SlowThreshold != 0 && elapsed > gl.SlowThreshold && gl.Logger.IsWarnEnabled():
		sql, rows := getSqlAndRows(fc)
		slowLog := fmt.Sprintf("SLOW SQL >= %v", gl.SlowThreshold)
		gl.printf(log.LevelWarn, "%s [%.3fms] [rows:%v] %s", slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case gl.Logger.IsInfoEnabled():
		sql, rows := getSqlAndRows(fc)
		gl.printf(log.LevelInfo, "[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}

// Trace print sql message
func (gl *GormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if gl.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
