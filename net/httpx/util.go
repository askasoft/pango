package httpx

import (
	"net/http"
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
