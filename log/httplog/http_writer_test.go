package httplog

import (
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

/*
OpenSearch DevTools:

DELETE pango_logs

PUT pango_logs
{
	"mappings": {
		"properties": {
			"time": {
				"type": "date",
				"format": "date_time"
			}
		}
	}
}

GET pango_logs

GET pango_logs/_search
{
	"query": {
		"match_all": {}
	}
}

POST pango_logs/_delete_by_query
{
	"query": {
		"match_all": {}
	}
}
*/

// Test OpenSearch log writer
// set OPENSEARCH_URL=https://localhost:9200/pango_logs
func TestHTTPWriter(t *testing.T) {
	url := os.Getenv("OPENSEARCH_URL")
	if len(url) < 1 {
		t.Skip("OPENSEARCH_URL not set")
		return
	}

	url += "/_doc"

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	lg.SetFormat(`json:{"time": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "host":%x{HOST}, "version":%x{VERSON}, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)
	lg.SetProp("HOST", "localhost")
	lg.SetProp("VERSION", "1.0")

	ww := &HTTPWriter{
		URL:         url,
		ContentType: "application/json",
		Insecure:    true,
		Username:    "admin",
		Password:    "admin",
		Timeout:     time.Millisecond * 300,
	}

	ww.Filter = log.NewLevelFilter(log.LevelDebug)
	lg.SetWriter(log.NewMultiWriter(
		ww,
		&log.StreamWriter{Color: true},
	))

	lg.Trace("This is a HTTPWriter trace log")
	lg.Debug("This is a HTTPWriter debug log")
	lg.Info("This is a HTTPWriter info log")
	lg.Warn("This is a HTTPWriter warn log")
	lg.Error("This is a HTTPWriter error log")
	lg.Fatal("This is a HTTPWriter fatal log")

	lg.Close()
}

// Test OpenSearch batch log writer
// set OPENSEARCH_URL=https://localhost:9200/pango_logs
func TestWebhookBatchWriter(t *testing.T) {
	url := os.Getenv("OPENSEARCH_URL")
	if len(url) < 1 {
		t.Skip("OPENSEARCH_URL not set")
		return
	}

	url += "/_bulk"

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	lg.SetFormat(`json:{"create": {}}%n{"time": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "host":%x{HOST}, "version":%x{VERSON}, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)
	lg.SetProp("HOST", "localhost")
	lg.SetProp("VERSION", "1.0")

	ww := &HTTPWriter{
		URL:         url,
		ContentType: "application/json",
		Insecure:    true,
		Username:    "admin",
		Password:    "admin",
		Timeout:     time.Millisecond * 300,
		BatchWriter: log.BatchWriter{
			CacheCount: 6,
			BatchCount: 3,
			FlushLevel: log.LevelWarn,
			FlushDelta: time.Second,
		},
	}

	ww.Filter = log.NewLevelFilter(log.LevelDebug)
	lg.SetWriter(log.NewMultiWriter(
		ww,
		&log.StreamWriter{Color: true},
	))

	lg.Trace("This is a HTTPWriter(batch) trace log")
	lg.Debug("This is a HTTPWriter(batch) debug log")
	lg.Info("This is a HTTPWriter(batch) info log")
	lg.Info("This is a HTTPWriter(batch) info log2, should flush by BatchCount")
	lg.Warn("This is a HTTPWriter(batch) warn log, should flush by FlushLevel")

	lg.Info("This is a HTTPWriter(batch) info log2")
	time.Sleep(time.Millisecond * 1200)
	lg.Info("This is a HTTPWriter(batch) info log3, should flush by FlushDelta")

	lg.Error("This is a HTTPWriter(batch) error log")
	lg.Fatal("This is a HTTPWriter(batch) fatal log")

	lg.Close()
}
