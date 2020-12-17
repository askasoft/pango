package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// ElasticSearchWriter implements log Writer Interface and send log message to Elastic Search.
type ElasticSearchWriter struct {
	URL     string
	Timeout time.Duration
	Logfmt  Formatter // log formatter
	Logfil  Filter    // log filter

	hc *http.Client
	bb bytes.Buffer
}

// SetFormat set the log formatter
func (ew *ElasticSearchWriter) SetFormat(format string) {
	ew.Logfmt = NewJSONFormatter(format)
}

// SetFilter set the log filter
func (ew *ElasticSearchWriter) SetFilter(filter string) {
	ew.Logfil = NewLogFilter(filter)
}

// SetTimeout set timeout
func (ew *ElasticSearchWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("ElasticSearchWriter - Invalid timeout: %v", err)
	}
	ew.Timeout = tmo
	return nil
}

// Write send log message to elasticsearch
func (ew *ElasticSearchWriter) Write(le *Event) {
	if ew.Logfil != nil && ew.Logfil.Reject(le) {
		return
	}

	if ew.Logfmt == nil {
		ew.Logfmt = le.Logger.GetFormatter()
	}

	if ew.hc == nil {
		ew.hc = &http.Client{Timeout: ew.Timeout}
	}

	// format msg
	ew.bb.Reset()
	ew.Logfmt.Write(&ew.bb, le)

	req, err := http.NewRequest("POST", ew.URL, &ew.bb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ElasticSearchWriter(%q) - NewRequest(): %v\n", ew.URL, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := ew.hc.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ElasticSearchWriter(%q) - POST(): %v\n", ew.URL, err)
		return
	}
	io.Copy(ioutil.Discard, res.Body)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		fmt.Fprintf(os.Stderr, "ElasticSearchWriter(%q) - POST(): %s\n", ew.URL, res.Status)
	}
}

// Flush implementing method. empty.
func (ew *ElasticSearchWriter) Flush() {
}

// Close implementing method. empty.
func (ew *ElasticSearchWriter) Close() {
}

func init() {
	RegisterWriter("es", func() Writer {
		return &ElasticSearchWriter{}
	})
}
