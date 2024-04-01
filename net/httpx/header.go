package httpx

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/askasoft/pango/str"
)

// NewHeaderWriter create a http.ResponseWriter to append on WriteHeader(statusCode int).
// if statusCode != 200, header will not append.
// a existing header will not be overwriten.
func NewHeaderWriter(w http.ResponseWriter, key, value string) http.ResponseWriter {
	return &headerWriter{w, key, value}
}

// headerWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headerWriter struct {
	http.ResponseWriter
	key   string
	value string
}

// WriteHeader append header when statusCode == 200
func (hw *headerWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		if hw.Header().Get(hw.key) == "" {
			hw.Header().Add(hw.key, hw.value)
		}
	}
	hw.ResponseWriter.WriteHeader(statusCode)
}

// NewHeadersWriter create a http.ResponseWriter to append on WriteHeader(statusCode int).
// if statusCode != 200, header will not append.
// a existing header will not be overwriten.
func NewHeadersWriter(w http.ResponseWriter, h map[string]string) http.ResponseWriter {
	return &headersWriter{w, h}
}

// headersWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headersWriter struct {
	http.ResponseWriter
	header map[string]string
}

// WriteHeader append header when statusCode == 200
func (hw *headersWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		for k, v := range hw.header {
			if hw.Header().Get(k) == "" {
				hw.Header().Add(k, v)
			}
		}
	}
	hw.ResponseWriter.WriteHeader(statusCode)
}

// SetAttachmentHeader set header Content-Disposition: attachment; filename=...
func SetAttachmentHeader(header http.Header, filename string) {
	var v string

	if str.IsASCII(filename) {
		v = `attachment; filename=` + strconv.Quote(filename)
	} else {
		v = `attachment; filename*=UTF-8''` + url.QueryEscape(filename)
	}

	header.Set("Content-Disposition", v)
}
