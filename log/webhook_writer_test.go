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
// set OPENSEARCH_WEBHOOK=https://localhost:9200/pango_logs/_doc
func TestWebhookWriter(t *testing.T) {
	url := os.Getenv("OPENSEARCH_WEBHOOK")
	if len(url) < 1 {
		t.Skip("OPENSEARCH_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetFormat(`json:{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "host":%x{HOST}, "version":%x{VERSON}, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)
	log.SetProp("HOST", "localhost")
	log.SetProp("VERSION", "1.0")

	ww := &WebhookWriter{
		Webhook:     url,
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

	log.Trace("This is a webhook trace log")
	log.Debug("This is a webhook debug log")
	log.Info("This is a webhook info log")
	log.Warn("This is a webhook warn log")
	log.Error("This is a webhook error log")
	log.Fatal("This is a webhook fatal log")

	log.Close()
}
