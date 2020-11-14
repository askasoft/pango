package log

import (
	"os"
	"testing"
)

// Test slack log
func TestSlackLog(t *testing.T) {
	wh := os.Getenv("SLACK_WEBHOOK")
	if len(wh) < 1 {
		return
	}

	log1 := NewLog()
	log1.SetLevel(LevelTrace)
	log1.SetWriter(&SlackWriter{Level: LevelTrace, Webhook: os.Getenv("SLACK_WEBHOOK")})

	log1.Info("This is a slack info log")
}
