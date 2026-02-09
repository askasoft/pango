package slack

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSlackError(t *testing.T) {
	e := &Error{Status: "Too Many Requests", StatusCode: 429, RetryAfter: 30 * time.Second}
	ev := "429 Too Many Requests (Retry-After: 30s)"
	av := fmt.Sprintf("%v", e)
	if ev != av {
		t.Errorf("Got %v, want %v", av, ev)
	}
}

func postSlack(t *testing.T, sm *Message) {
	url := os.Getenv("SLACK_WEBHOOK")
	if len(url) < 1 {
		t.Skip("SLACK_WEBHOOK not set")
		return
	}
	err := Post(url, time.Second*5, sm)
	if err != nil {
		t.Error(err)
	}
}

// Test post slack message
func TestSlackPostText(t *testing.T) {
	sm := &Message{Text: "TestSlackPost"}
	postSlack(t, sm)
}

// Test post slack message with icon
func TestSlackPostWithIcon(t *testing.T) {
	sm := &Message{IconEmoji: ":bug:", Text: "**TestSlackPostWithIcon**"}
	postSlack(t, sm)
}

// Test post slack message with attach
func TestSlackPostWithAttach(t *testing.T) {
	sm := &Message{IconEmoji: ":fire:", Text: "**TestSlackPostWithAttach**"}
	sm.AddAttachment(&Attachment{Text: "**attachment text**"})
	postSlack(t, sm)
}
