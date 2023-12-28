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
// set OPENSEARCH_BULK_WEBHOOK=https://localhost:9200/pango_logs/_bulk
func TestWebhookBatchWriter(t *testing.T) {
	url := os.Getenv("OPENSEARCH_BULK_WEBHOOK")
	if len(url) < 1 {
		t.Skip("OPENSEARCH_BULK_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetFormat(`json:{"create": {}}%n{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)

	ww := &WebhookBatchWriter{
		Webhook:     url,
		ContentType: "application/json",
		Insecure:    true,
		Username:    "admin",
		Password:    "admin",
		Timeout:     time.Millisecond * 300,
		CacheCount:  6,
		BatchCount:  3,
		FlushLevel:  LevelWarn,
		FlushDelta:  time.Second,
	}

	ww.Filter = NewLevelFilter(LevelDebug)
	log.SetWriter(NewMultiWriter(
		ww,
		&StreamWriter{Color: true},
	))

	log.Trace("This is a webhook trace log")
	log.Debug("This is a webhook debug log")
	log.Info("This is a webhook info log")
	log.Info("This is a webhook info log2, should flush by BatchCount")
	log.Warn("This is a webhook warn log, should flush by FlushLevel")

	log.Info("This is a webhook info log2")
	time.Sleep(time.Millisecond * 1200)
	log.Info("This is a webhook info log3, should flush by FlushDelta")

	log.Error("This is a webhook error log")
	log.Fatal("This is a webhook fatal log")

	log.Close()
}
