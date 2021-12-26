package log

import (
	"os"
	"testing"
	"time"
)

// Test elasticsearch log
// create index: curl -X PUT "http://localhost:9200/pango?pretty"
func TestWebhook_ESLog(t *testing.T) {
	url := os.Getenv("ES_WEBHOOK")
	if len(url) < 1 {
		t.Skip("ES_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetFormatter(NewJSONFormatter(
		`{"when": %t{2006-01-02T15:04:05.000Z07:00}, "level": %l, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`))
	log.SetWriter(NewMultiWriter(
		&WebhookWriter{
			Webhook:     url,
			ContentType: "application/json",
			Logfil:      NewLevelFilter(LevelDebug),
			Timeout:     time.Millisecond * 300,
		},
		&StreamWriter{Color: true},
	))

	log.Trace("This is a elasticsearch trace log")
	log.Debug("This is a elasticsearch debug log")
	log.Info("This is a elasticsearch info log")
	log.Warn("This is a elasticsearch warn log")
	log.Error("This is a elasticsearch error log")
	log.Fatal("This is a elasticsearch fatal log")
}
