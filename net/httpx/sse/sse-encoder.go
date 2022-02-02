package sse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Server-Sent Events
// W3C Working Draft 29 October 2009
// http://www.w3.org/TR/2009/WD-eventsource-20091029/

// ContentType "text/event-stream"
const ContentType = "text/event-stream"

var contentType = []string{ContentType}
var noCache = []string{"no-cache"}

var fieldReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r")

var dataReplacer = strings.NewReplacer(
	"\n", "\ndata:",
	"\r", "\\r")

// Event text event struct
type Event struct {
	Event string
	ID    string
	Retry uint
	Data  interface{}
}

// Encode encode text event
func Encode(writer io.Writer, event Event) error {
	w := checkWriter(writer)
	writeID(w, event.ID)
	writeEvent(w, event.Event)
	writeRetry(w, event.Retry)
	return writeData(w, event.Data)
}

func writeID(w stringWriter, id string) {
	if len(id) > 0 {
		w.WriteString("id:")             //nolint: errcheck
		fieldReplacer.WriteString(w, id) //nolint: errcheck
		w.WriteString("\n")              //nolint: errcheck
	}
}

func writeEvent(w stringWriter, event string) {
	if len(event) > 0 {
		w.WriteString("event:")             //nolint: errcheck
		fieldReplacer.WriteString(w, event) //nolint: errcheck
		w.WriteString("\n")                 //nolint: errcheck
	}
}

func writeRetry(w stringWriter, retry uint) {
	if retry > 0 {
		w.WriteString("retry:")                              //nolint: errcheck
		w.WriteString(strconv.FormatUint(uint64(retry), 10)) //nolint: errcheck
		w.WriteString("\n")                                  //nolint: errcheck
	}
}

func writeData(w stringWriter, data interface{}) error {
	w.WriteString("data:") //nolint: errcheck
	switch kindOfData(data) {
	case reflect.Struct, reflect.Slice, reflect.Map:
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return err
		}
		w.WriteString("\n") //nolint: errcheck
	default:
		dataReplacer.WriteString(w, fmt.Sprint(data)) //nolint: errcheck
		w.WriteString("\n\n")                         //nolint: errcheck
	}
	return nil
}

// Render write event to http response
func (r Event) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return Encode(w, r)
}

// WriteContentType write content type header to http response
func (r Event) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	header["Content-Type"] = contentType

	if _, exist := header["Cache-Control"]; !exist {
		header["Cache-Control"] = noCache
	}
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
