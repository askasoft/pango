package log

import (
	"os"
	"testing"
)

// Test elasticsearch log
func TestElasticSearchLog(t *testing.T) {
	os.Setenv("ELASTICSEARCH_URL", "http://localhost:9200/pango/log")

	url := os.Getenv("ELASTICSEARCH_URL")
	if len(url) < 1 {
		t.Skip("ELASTICSEARCH_URL not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetFormatter(JSONFmtDefault)
	log.SetWriter(NewMultiWriter(
		&ElasticSearchWriter{
			URL:    url,
			Logfil: NewLevelFilter(LevelDebug),
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
