package lineworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/num"
)

// Error lineworks webhook error
type Error struct {
	Status     string
	StatusCode int
	RetryAfter time.Duration
}

func (e *Error) GetRetryAfter() time.Duration {
	return e.RetryAfter
}

// Error return error string
func (e *Error) Error() string {
	if e.RetryAfter != 0 {
		return fmt.Sprintf("%d %s (Retry-After: %s)", e.StatusCode, e.Status, e.RetryAfter)
	}
	return fmt.Sprintf("%d %s", e.StatusCode, e.Status)
}

type Body struct {
	Text string `json:"text,omitempty"`
}

type Button struct {
	Label string `json:"label,omitempty"`
	URL   string `json:"url,omitempty"`
}

// Message lineworks message
type Message struct {
	Title  string `json:"title,omitempty"`
	Body   Body   `json:"body,omitzero"`
	Button Button `json:"button,omitzero"`
}

// Post post lineworks message
func Post(url string, timeout time.Duration, msg any) error {
	bs, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: timeout}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		e := &Error{Status: res.Status, StatusCode: res.StatusCode}
		if res.StatusCode == http.StatusTooManyRequests {
			ra := res.Header.Get("RateLimit-Reset")
			if ra != "" {
				e.RetryAfter = time.Duration(num.Atoi(ra)) * time.Second
			}
		}
		return e
	}
	return nil
}
