package xin

import (
	"bufio"
	"net"
	"net/http"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

// ResponseWriter ...
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Status returns the HTTP response status code of the current request.
	Status() int

	// Size returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// WriteString writes the string into the response body.
	WriteString(string) (int, error)

	// Written returns true if the response body was already written.
	Written() bool

	// WriteHeaderNow forces to write the http header (status code + headers).
	WriteHeaderNow()

	// Pusher get the http.Pusher for server push
	Pusher() http.Pusher
}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
	logger log.Logger
}

func (w *responseWriter) reset(writer http.ResponseWriter, logger log.Logger) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
	w.logger = logger
}

func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			w.logger.Warnf("Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = iox.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotifier interface.
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flusher interface.
func (w *responseWriter) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *responseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

// NewHeaderWriter create a ResponseWriter to append on WriteHeader(statusCode int).
// if statusCode != 200, header will not append.
// a existing header will not be overwriten.
func NewHeaderWriter(w ResponseWriter, key, value string, overwrite ...bool) ResponseWriter {
	return &headerWriter{w, key, value, asg.First(overwrite)}
}

// headerWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headerWriter struct {
	ResponseWriter
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

// NewHeadersWriter create a ResponseWriter to append on WriteHeader(statusCode int).
// if statusCode != 200, header will not append.
// a existing header will not be overwriten.
func NewHeadersWriter(w ResponseWriter, headers map[string]string, overwrite ...bool) ResponseWriter {
	return &headersWriter{w, headers, asg.First(overwrite)}
}

// headersWriter write header when statusCode == 200 on WriteHeader(statusCode int)
// a existing header will not be overwriten.
type headersWriter struct {
	ResponseWriter
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
