package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/askasoft/pango/iox"
)

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
	Type        string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Text        string `json:"text,omitempty"`
	Attachments []any  `json:"attachments,omitempty"`
}

// Post post teams message
func Post(url string, timeout time.Duration, msg any) error {
	bs, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Println(string(bs))
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

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		e := &Error{Status: res.Status, StatusCode: res.StatusCode}
		ra := res.Header.Get("Retry-After")
		if ra != "" {
			e.RetryAfter, _ = strconv.Atoi(ra)
		}
		return e
	}
	return nil
}
