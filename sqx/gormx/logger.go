package gormx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type GormSQLPrinter struct {
	sb strings.Builder
}

func (gsp *GormSQLPrinter) SQL() string {
	return gsp.sb.String()
}

func (gsp *GormSQLPrinter) Printf(msg string, data ...any) {
	s := fmt.Sprintf(msg, data...) + ";\n"

	fmt.Print(s)
	gsp.sb.WriteString(s)
}

// LogMode log mode
func (gsp *GormSQLPrinter) LogMode(level logger.LogLevel) logger.Interface {
	return gsp
}

// Info print info
func (gsp *GormSQLPrinter) Info(ctx context.Context, msg string, data ...any) {
	gsp.Printf(msg, data...)
}

// Warn print warn messages
func (gsp *GormSQLPrinter) Warn(ctx context.Context, msg string, data ...any) {
	gsp.Printf(msg, data...)
}

// Error print error messages
func (gsp *GormSQLPrinter) Error(ctx context.Context, msg string, data ...any) {
	gsp.Printf(msg, data...)
}

// Trace print sql message
func (gsp *GormSQLPrinter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	gsp.Printf("%s", sql)
}

// Trace print sql message
func (gsp *GormSQLPrinter) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	return sql, params
}
