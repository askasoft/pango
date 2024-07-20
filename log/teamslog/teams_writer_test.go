package teamslog

import (
	"fmt"
	"os"
	"testing"

	"github.com/askasoft/pango/log"
)

func skipTest(t *testing.T, msg string) {
	fmt.Println(msg)
	t.Skip(msg)
}

func TestTeamsWriter(t *testing.T) {
	wh := os.Getenv("TEAMS_WEBHOOK")
	if wh == "" {
		skipTest(t, "TEAMS_WEBHOOK not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	sw := &TeamsWriter{Webhook: wh}
	sw.Filter = log.NewLevelFilter(log.LevelInfo)
	lg.SetWriter(sw)

	lg.Debug("This is a teams debug log")
	lg.Info("This is a teams info log")
}
