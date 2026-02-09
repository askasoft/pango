package lineworks

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLineworksError(t *testing.T) {
	e := &Error{Status: "Too Many Requests", StatusCode: 429, RetryAfter: 30 * time.Second}
	ev := "429 Too Many Requests (Retry-After: 30s)"
	av := fmt.Sprintf("%v", e)
	if ev != av {
		t.Errorf("Got %v, want %v", av, ev)
	}
}

func postLineworks(t *testing.T, lm *Message) {
	url := os.Getenv("LINEWORKS_WEBHOOK")
	if len(url) < 1 {
		t.Skip("LINEWORKS_WEBHOOK not set")
		return
	}

	err := Post(url, time.Second*5, lm)
	if err != nil {
		t.Error(err)
	}
}

// Test post lineworks message
func TestLineworksPostText(t *testing.T) {
	lm := &Message{
		Title: time.Now().String(),
		Body:  Body{Text: "Text " + time.Now().String()},
	}
	postLineworks(t, lm)
}
