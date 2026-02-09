package lineworkslog

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

func TestLineWorksWriter(t *testing.T) {
	wh := os.Getenv("LINEWORKS_WEBHOOK")
	if wh == "" {
		skipTest(t, "LINEWORKS_WEBHOOK not set")
		return
	}

	lg := log.NewLog()
	lg.SetLevel(log.LevelTrace)
	sw := &LineWorksWriter{Webhook: wh}
	sw.Filter = log.NewLevelFilter(log.LevelInfo)
	lg.SetWriter(sw)

	lg.Info("This is a info log")
}
