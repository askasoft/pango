package sqlog

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/askasoft/pango/log"
)

var (
	testDSN = ""
	testSQL = "INSERT INTO sqlogs (time, level, msg, file, line, func, trace) VALUES"
	testARG = "%t %p %m %S %L %F %T"
)

func init() {
	dsn := os.Getenv("SQLW_DSN")
	if dsn == "" {
		return
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS sqlogs")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(`
CREATE TABLE sqlogs (
	id serial NOT NULL,
	time timestamp with time zone NOT NULL,
	level char(1) NOT NULL,
	msg text NOT NULL,
	file text NOT NULL,
	line integer NOT NULL,
	func text NOT NULL,
	trace text NOT NULL
)
	`)
	if err != nil {
		fmt.Println(err)
		return
	}

	testDSN = dsn
}

// Test sql log writer
// set SQLW_DSN=host=127.0.0.1 user=pango password=pango dbname=pango port=5432 sslmode=disable
func TestSQLWriter(t *testing.T) {
	if testDSN == "" {
		t.Skip("SQLW_DSN not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	lg.SetProp("HOST", "localhost")
	lg.SetProp("VERSION", "1.0")

	sw := &SQLWriter{
		Driver:    "postgres",
		Dsn:       testDSN,
		Statement: testSQL,
	}
	sw.SetParameter(testARG)

	sw.Filter = log.NewLevelFilter(log.LevelDebug)
	lg.SetWriter(log.NewMultiWriter(
		sw,
		&log.StreamWriter{Color: true},
	))

	lg.Trace("This is a SQLWriter trace log")
	lg.Debug("This is a SQLWriter debug log")
	lg.Info("This is a SQLWriter info log")
	lg.Warn("This is a SQLWriter warn log")
	lg.Error("This is a SQLWriter error log")

	lg.Close()
}

// Test batch sql log writer
func TestWebhookBatchWriter(t *testing.T) {
	if testDSN == "" {
		t.Skip("SQLW_DSN not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	lg.SetProp("HOST", "localhost")
	lg.SetProp("VERSION", "1.0")

	sw := &SQLWriter{
		Driver:    "postgres",
		Dsn:       testDSN,
		Statement: testSQL,
		BatchSupport: log.BatchSupport{
			BatchCount: 3,
			CacheCount: 6,
			FlushLevel: log.LevelWarn,
			FlushDelta: time.Second,
		},
	}
	sw.SetParameter(testARG)

	sw.Filter = log.NewLevelFilter(log.LevelDebug)
	lg.SetWriter(log.NewMultiWriter(
		sw,
		&log.StreamWriter{Color: true},
	))

	lg.Trace("This is a SQLWriter(batch) trace log")
	lg.Debug("This is a SQLWriter(batch) debug log")
	lg.Info("This is a SQLWriter(batch) info log")
	lg.Info("This is a SQLWriter(batch) info log2, should flush by BatchCount")
	lg.Warn("This is a SQLWriter(batch) warn log, should flush by FlushLevel")

	lg.Info("This is a SQLWriter(batch) info log2")
	time.Sleep(time.Millisecond * 1200)
	lg.Info("This is a SQLWriter(batch) info log3, should flush by FlushDelta")

	lg.Error("This is a SQLWriter(batch) error log")

	lg.Close()
}
