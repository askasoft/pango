package httpx

import (
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/str"
)

// NewHeaderWriter create a http.ResponseWriter to append on WriteHeader(statusCode int).
// if statusCode != 200, header will not append.
// a existing header will not be overwriten.
func NewHeaderWriter(w http.ResponseWriter, key, value string, overwrite ...bool) http.ResponseWriter {
	return &headerWriter{w, key, value, asg.First(overwrite)}
}

// headerWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headerWriter struct {
	http.ResponseWriter
	key       string
	value     string
	overwrite bool
}

// WriteHeader append header when statusCode == 200
func (hw *headerWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		if hw.overwrite || hw.Header().Get(hw.key) == "" {
			hw.Header().Set(hw.key, hw.value)
		}
	}
	hw.ResponseWriter.WriteHeader(statusCode)
}

// NewHeadersWriter create a http.ResponseWriter to append on WriteHeader(statusCode int).
// if statusCode != 200, header will not append.
// a existing header will not be overwriten.
func NewHeadersWriter(w http.ResponseWriter, headers map[string]string, overwrite ...bool) http.ResponseWriter {
	return &headersWriter{w, headers, asg.First(overwrite)}
}

// headersWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headersWriter struct {
	http.ResponseWriter
	headers   map[string]string
	overwrite bool
}

// WriteHeader append header when statusCode == 200
func (hw *headersWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		for k, v := range hw.headers {
			if hw.overwrite || hw.Header().Get(k) == "" {
				hw.Header().Set(k, v)
			}
		}
	}
	hw.ResponseWriter.WriteHeader(statusCode)
}

// SetAttachmentHeader set header Content-Disposition: attachment; filename=... and Content-Type: mime.TypeByExtension(filename)
func SetAttachmentHeader(hh http.Header, filename string) {
	if _, ok := hh["Content-Type"]; !ok {
		ct := mime.TypeByExtension(filepath.Ext(filename))
		if ct != "" {
			hh.Set("Content-Type", ct)
		}
	}

	var cd string
	if str.IsASCII(filename) {
		cd = `attachment; filename=` + strconv.Quote(filename)
	} else {
		cd = `attachment; filename*=UTF-8''` + url.QueryEscape(filename)
	}
	hh.Set("Content-Disposition", cd)
}
