package log

import (
	"os"
	"testing"
	"time"
)

/*
OpenSearch DevTools:

DELETE pango_logs

PUT pango_logs
{
	"mappings": {
		"properties": {
			"when": {
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

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetFormat(`json:{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "host":%x{HOST}, "version":%x{VERSON}, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)
	log.SetProp("HOST", "localhost")
	log.SetProp("VERSION", "1.0")

	ww := &HTTPWriter{
		URL:         url,
		ContentType: "application/json",
		Insecure:    true,
		Username:    "admin",
		Password:    "admin",
		Timeout:     time.Millisecond * 300,
	}

	ww.Filter = NewLevelFilter(LevelDebug)
	log.SetWriter(NewMultiWriter(
		ww,
		&StreamWriter{Color: true},
	))

	log.Trace("This is a HTTPWriter trace log")
	log.Debug("This is a HTTPWriter debug log")
	log.Info("This is a HTTPWriter info log")
	log.Warn("This is a HTTPWriter warn log")
	log.Error("This is a HTTPWriter error log")
	log.Fatal("This is a HTTPWriter fatal log")

	log.Close()
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

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetFormat(`json:{"create": {}}%n{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "host":%x{HOST}, "version":%x{VERSON}, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)
	log.SetProp("HOST", "localhost")
	log.SetProp("VERSION", "1.0")

	ww := &HTTPWriter{
		URL:         url,
		ContentType: "application/json",
		Insecure:    true,
		Username:    "admin",
		Password:    "admin",
		Timeout:     time.Millisecond * 300,
		BatchWriter: BatchWriter{
			CacheCount: 6,
			BatchCount: 3,
			FlushLevel: LevelWarn,
			FlushDelta: time.Second,
		},
	}

	ww.Filter = NewLevelFilter(LevelDebug)
	log.SetWriter(NewMultiWriter(
		ww,
		&StreamWriter{Color: true},
	))

	log.Trace("This is a HTTPWriter(batch) trace log")
	log.Debug("This is a HTTPWriter(batch) debug log")
	log.Info("This is a HTTPWriter(batch) info log")
	log.Info("This is a HTTPWriter(batch) info log2, should flush by BatchCount")
	log.Warn("This is a HTTPWriter(batch) warn log, should flush by FlushLevel")

	log.Info("This is a HTTPWriter(batch) info log2")
	time.Sleep(time.Millisecond * 1200)
	log.Info("This is a HTTPWriter(batch) info log3, should flush by FlushDelta")

	log.Error("This is a HTTPWriter(batch) error log")
	log.Fatal("This is a HTTPWriter(batch) fatal log")

	log.Close()
}
