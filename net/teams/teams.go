package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pandafw/pango/iox"
)

// Post post teams message
func Post(url string, timeout time.Duration, tm *Message) error {
	bs, err := json.Marshal(tm)
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
		ra := res.Header.Get("Retry-After")
		if ra != "" {
			e.RetryAfter, _ = strconv.Atoi(ra)
		}
		return e
	}
	return nil
}

// Error teams api error
type Error struct {
	Status     string
	StatusCode int
	RetryAfter int
}

// Error return error string
func (e *Error) Error() string {
	if e.RetryAfter != 0 {
		return fmt.Sprintf("%d %s (Retry-After: %d)", e.StatusCode, e.Status, e.RetryAfter)
	}
	return fmt.Sprintf("%d %s", e.StatusCode, e.Status)
}

// Message teams message
type Message struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Text  string `json:"text"`
}
