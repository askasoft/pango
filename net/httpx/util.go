package httpx

import (
	"net/http"
	"net/url"

	"github.com/askasoft/pango/str"
)

// NewHeaderAppender create a http.ResponseWriter to append on WriteHeader(statusCode int).
// if statusCode != 200, header will not append.
// a existing header will not be overwriten.
func NewHeaderAppender(w http.ResponseWriter, h map[string]string) http.ResponseWriter {
	return &headerWriter{w, h}
}

// headerWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headerWriter struct {
	http.ResponseWriter
	header map[string]string
}

// WriteHeader append header when statusCode == 200
func (hw *headerWriter) WriteHeader(statusCode int) {
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
		v = `attachment; filename="` + filename + `"`
	} else {
		v = `attachment; filename*=UTF-8''` + url.QueryEscape(filename)
	}

	header.Set("Content-Disposition", v)
}
