package teams

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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
	io.Copy(ioutil.Discard, res.Body) //nolint: errcheck
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		e := &Error{Status: res.Status, StatusCode: res.StatusCode}
		ra := res.Header.Get("retry-after")
		if ra != "" {
			e.RetryAfter, _ = strconv.Atoi(ra)
		}
		return e
	}
	return nil
}

// Error teams api error
type Error struct {
	StatusCode int
	Status     string
	RetryAfter int
}

// Error return error string
func (e *Error) Error() string {
	if e.RetryAfter != 0 {
		return e.Status + " - RetryAfter: " + strconv.Itoa(e.RetryAfter)
	}
	return e.Status
}

// Message teams message
type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}