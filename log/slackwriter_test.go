package log

import (
	"os"
	"testing"
)

// Test slack log
func TestSlackLog(t *testing.T) {
	wh := os.Getenv("SLACK_WEBHOOK")
	if len(wh) < 1 {
		skipTest(t, "SLACK_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	log.SetWriter(&SlackWriter{Webhook: wh, Username: "gotest", Logfil: NewLevelFilter(LevelInfo)})

	log.Debug("This is a slack debug log")
	log.Info("This is a slack info log")
}
