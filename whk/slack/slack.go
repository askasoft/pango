package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/iox"
)

// Error slack api error
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

var slackEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
	`>`, "&gt;",
)

// EscapeString escapes special characters.
// `&` => "&amp;"
// `<` => "&lt;"
// `>` => "&gt;"
func EscapeString(s string) string {
	return slackEscaper.Replace(s)
}

// Post post slack message
func Post(url string, timeout time.Duration, sm *Message) error {
	bs, err := json.Marshal(sm)
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
