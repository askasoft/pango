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

// https://techcommunity.microsoft.com/discussions/teamsdeveloper/simple-workflow-to-replace-teams-incoming-webhooks/4225270
func TestTeamsWriter(t *testing.T) {
	url := os.Getenv("TEAMS_WEBHOOK")
	if url == "" {
		skipTest(t, "TEAMS_WEBHOOK not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	sw := &TeamsWriter{Webhook: url}
	sw.Filter = log.NewLevelFilter(log.LevelInfo)
	lg.SetWriter(sw)

	lg.Debug("This is a teams **debug** log")
	lg.Info("This is a teams **info** log. detail: This is detail message.")
}

func TestTeamsWriterAdaptive(t *testing.T) {
	url := os.Getenv("TEAMS_WEBHOOK_ADAPTIVE")
	if url == "" {
		skipTest(t, "TEAMS_WEBHOOK not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	sw := &TeamsWriter{Webhook: url, Style: "adaptive"}
	sw.Filter = log.NewLevelFilter(log.LevelInfo)
	lg.SetWriter(sw)

	lg.Debug("This is a teams **debug** log")
	lg.Info("This is a teams **info** log. \ndetail: This is detail message.")
}
