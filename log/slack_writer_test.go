package log

import (
	"os"
	"testing"
)

// Test slack log
func TestSlackLog(t *testing.T) {
	wh := os.Getenv("SLACK_WEBHOOK")
	if wh == "" {
		skipTest(t, "SLACK_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	sw := &SlackWriter{Webhook: wh, Logfil: NewLevelFilter(LevelInfo)}
	log.SetWriter(sw)

	log.Debug("This is a slack debug log")
	log.Info("This is a slack info log")
}
