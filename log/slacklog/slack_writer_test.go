package slacklog

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

func TestSlackWriter(t *testing.T) {
	wh := os.Getenv("SLACK_WEBHOOK")
	if wh == "" {
		skipTest(t, "SLACK_WEBHOOK not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	sw := &SlackWriter{Webhook: wh}
	sw.Filter = log.NewLevelFilter(log.LevelInfo)
	lg.SetWriter(sw)

	lg.Debug("This is a slack debug log")
	lg.Info("This is a slack info log")
}
