package log

import (
	"os"
	"testing"
)

// Test teams log
func TestTeamsLog(t *testing.T) {
	wh := os.Getenv("TEAMS_WEBHOOK")
	if wh == "" {
		skipTest(t, "TEAMS_WEBHOOK not set")
		return
	}

	log := NewLog()
	log.SetLevel(LevelTrace)
	sw := &TeamsWriter{Webhook: wh, Logfil: NewLevelFilter(LevelInfo)}
	log.SetWriter(sw)

	log.Debug("This is a teams debug log")
	log.Info("This is a teams info log")
}
