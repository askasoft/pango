package log

import (
	"os"
	"testing"
)

// Test elasticsearch log
// create index: curl -X PUT "http://localhost:9200/pango?pretty"
func TestWebhook_ESLog(t *testing.T) {
	//os.Setenv("ES_WEBHOOK", "http://localhost:9200/pango/logs")

	url := os.Getenv("ES_WEBHOOK")
	if len(url) < 1 {
		t.Skip("ES_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetFormatter(JSONFmtDefault)
	log.SetWriter(NewMultiWriter(
		&WebhookWriter{
			Webhook:     url,
			ContentType: "application/json",
			Logfil:      NewLevelFilter(LevelDebug),
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
