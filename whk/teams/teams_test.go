package teams

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTeamsError(t *testing.T) {
	e := &Error{Status: "Too Many Requests", StatusCode: 429, RetryAfter: 30 * time.Second}
	ev := "429 Too Many Requests (Retry-After: 30s)"
	av := fmt.Sprintf("%v", e)
	if ev != av {
		t.Errorf("Got %v, want %v", av, ev)
	}
}

func postTeams(t *testing.T, sm *Message) {
	url := os.Getenv("TEAMS_WEBHOOK")
	if len(url) < 1 {
		t.Skip("TEAMS_WEBHOOK not set")
		return
	}
	err := Post(url, time.Second*5, sm)
	if err != nil {
		t.Error(err)
	}
}

// Test post teams message
func TestTeamsPostText(t *testing.T) {
	sm := &Message{Text: "TestTeamsPost"}
	postTeams(t, sm)
}
