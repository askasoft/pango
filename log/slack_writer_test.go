package log

import (
	"os"
	"testing"
)

func TestSlackWriter(t *testing.T) {
	wh := os.Getenv("SLACK_WEBHOOK")
	if wh == "" {
		skipTest(t, "SLACK_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	sw := &SlackWriter{Webhook: wh}
	sw.Filter = NewLevelFilter(LevelInfo)
	log.SetWriter(sw)

	log.Debug("This is a slack debug log")
	log.Info("This is a slack info log")
}
