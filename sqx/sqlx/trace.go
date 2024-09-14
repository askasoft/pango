package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/askasoft/pango/sqx"
)

type tracer struct {
	Bind  Binder
	Trace Trace
}

func (t *tracer) ExplainSQL(sql string, args ...any) string {
	return t.Bind.Explain(sql, args...)
}

func (t *tracer) TracePing(pr sqx.Pinger) error {
	start := time.Now()
	err := pr.Ping()
	if t.Trace != nil {
		t.Trace(start, "Ping()", -1, err)
	}
	return err
}

func (t *tracer) TracePingContext(ctx context.Context, pr sqx.ContextPinger) error {
	start := time.Now()
	err := pr.PingContext(ctx)
	if t.Trace != nil {
		t.Trace(start, "PingContext()", -1, err)
	}
	return err
}

func (t *tracer) TraceQuery(qr sqx.Queryer, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := qr.Query(query, args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, err)
	}
	return rows, err
}

func (t *tracer) TraceQueryRow(rqr sqx.RowQueryer, query string, args ...any) *sql.Row {
	start := time.Now()
	row := rqr.QueryRow(query, args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, row.Err())
	}
	return row
}

func (t *tracer) TraceStmtQuery(sqr sqx.StmtQueryer, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := sqr.Query(args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, err)
	}
	return rows, err
}

func (t *tracer) TraceQueryContext(ctx context.Context, cqr sqx.ContextQueryer, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := cqr.QueryContext(ctx, query, args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, err)
	}
	return rows, err
}

func (t *tracer) TraceQueryRowContext(ctx context.Context, crqr sqx.ContextRowQueryer, query string, args ...any) *sql.Row {
	start := time.Now()
	row := crqr.QueryRowContext(ctx, query, args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, row.Err())
	}
	return row
}

func (t *tracer) TraceStmtQueryContext(ctx context.Context, csqr sqx.ContextStmtQueryer, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := csqr.QueryContext(ctx, args...)
	if t.Trace != nil {
		t.Trace(start, t.ExplainSQL(query, args...), -1, err)
	}
	return rows, err
}

func (t *tracer) TraceExec(er sqx.Execer, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	sqr, err := er.Exec(query, args...)
	if t.Trace != nil {
		cnt, _ := sqr.RowsAffected()
		t.Trace(start, t.ExplainSQL(query, args...), cnt, err)
	}
	return sqr, err
}

func (t *tracer) TraceStmtExec(ser sqx.StmtExecer, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	sqr, err := ser.Exec(args...)
	if t.Trace != nil {
		cnt, _ := sqr.RowsAffected()
		t.Trace(start, t.ExplainSQL(query, args...), cnt, err)
	}
	return sqr, err
}

func (t *tracer) TraceExecContext(ctx context.Context, cer sqx.ContextExecer, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	sqr, err := cer.ExecContext(ctx, query, args...)
	if t.Trace != nil {
		cnt, _ := sqr.RowsAffected()
		t.Trace(start, t.ExplainSQL(query, args...), cnt, err)
	}
	return sqr, err
}

func (t *tracer) TraceStmtExecContext(ctx context.Context, scer sqx.ContextStmtExecer, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	sqr, err := scer.ExecContext(ctx, args...)
	if t.Trace != nil {
		cnt, _ := sqr.RowsAffected()
		t.Trace(start, t.ExplainSQL(query, args...), cnt, err)
	}
	return sqr, err
}

func (t *tracer) TracePrepare(pr sqx.Preparer, query string) (*sql.Stmt, error) {
	start := time.Now()
	stmt, err := pr.Prepare(query)
	if t.Trace != nil {
		t.Trace(start, "Prepare: "+query, -1, err)
	}
	return stmt, err
}

func (t *tracer) TracePrepareContext(ctx context.Context, cpr sqx.ContextPreparer, query string) (*sql.Stmt, error) {
	start := time.Now()
	stmt, err := cpr.PrepareContext(ctx, query)
	if t.Trace != nil {
		t.Trace(start, "PrepareContext: "+query, -1, err)
	}
	return stmt, err
}

func (t *tracer) TraceBegin(btr sqx.Beginer) (*sql.Tx, error) {
	start := time.Now()
	tx, err := btr.Begin()
	if t.Trace != nil {
		t.Trace(start, "Begin()", -1, err)
	}
	return tx, err
}

func (t *tracer) TraceBeginTx(ctx context.Context, btr sqx.BeginTxer, opts *sql.TxOptions) (*sql.Tx, error) {
	start := time.Now()
	tx, err := btr.BeginTx(ctx, opts)
	if t.Trace != nil {
		if opts == nil {
			t.Trace(start, "BeginTx(nil)", -1, err)
		} else {
			t.Trace(start, fmt.Sprintf("BeginTx(%v, %v)", opts.Isolation, opts.ReadOnly), -1, err)
		}
	}
	return tx, err
}

func (t *tracer) TraceCommit(cr sqx.Txer) error {
	start := time.Now()
	err := cr.Commit()
	if t.Trace != nil {
		t.Trace(start, "Commit()", -1, err)
	}
	return err
}

func (t *tracer) TraceRollback(rr sqx.Txer) error {
	start := time.Now()
	err := rr.Rollback()
	if t.Trace != nil {
		t.Trace(start, "Rollback()", -1, err)
	}
	return err
}
