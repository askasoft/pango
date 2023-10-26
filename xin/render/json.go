package render

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/askasoft/pango/str"
)

// JSON contains the given interface object.
type JSON struct {
	Data any
}

// IndentedJSON contains the given interface object.
type IndentedJSON struct {
	Data any
}

// SecureJSON contains the given interface object and its prefix.
type SecureJSON struct {
	Prefix string
	Data   any
}

// JsonpJSON contains the given interface object its callback.
type JsonpJSON struct {
	Callback string
	Data     any
}

// AsciiJSON contains the given interface object.
type AsciiJSON struct {
	Data any
}

// PureJSON contains the given interface object.
type PureJSON struct {
	Data any
}

var (
	jsonContentType      = "application/json; charset=utf-8"
	jsonpContentType     = "application/javascript; charset=utf-8"
	jsonASCIIContentType = "application/json"
)

// Render (JSON) writes data with custom ContentType.
func (r JSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return json.NewEncoder(w).Encode(r.Data)
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Render (IndentedJSON) marshals the given interface object and writes it with custom ContentType.
func (r IndentedJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	je := json.NewEncoder(w)
	je.SetIndent("", "  ")
	return je.Encode(r.Data)
}

// WriteContentType (IndentedJSON) writes JSON ContentType.
func (r IndentedJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Render (SecureJSON) marshals the given interface object and writes it with custom ContentType.
func (r SecureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if _, err := w.Write(str.UnsafeBytes(r.Prefix)); err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(r.Data)
}

// WriteContentType (SecureJSON) writes JSON ContentType.
func (r SecureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// Render (JsonpJSON) marshals the given interface object and writes it and its callback with custom ContentType.
func (r JsonpJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	callback := template.JSEscapeString(r.Callback)
	if _, err := w.Write(str.UnsafeBytes(callback)); err != nil {
		return err
	}

	if _, err := w.Write([]byte{'('}); err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(r.Data); err != nil {
		return err
	}

	if _, err := w.Write([]byte{')', ';'}); err != nil {
		return err
	}

	return nil
}

// WriteContentType (JsonpJSON) writes Javascript ContentType.
func (r JsonpJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpContentType)
}

// Render (AsciiJSON) marshals the given interface object and writes it with custom ContentType.
func (r AsciiJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bs, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	for len(bs) > 0 {
		r, n := utf8.DecodeRune(bs)
		if r >= unicode.MaxASCII {
			_, err = w.Write(str.UnsafeBytes(fmt.Sprintf("\\u%04x", int64(r))))
		} else {
			_, err = w.Write([]byte{byte(r)})
		}
		if err != nil {
			return err
		}

		bs = bs[n:]
	}

	return nil
}

// WriteContentType (AsciiJSON) writes JSON ContentType.
func (r AsciiJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonASCIIContentType)
}

// Render (PureJSON) writes custom ContentType and encodes the given interface object.
func (r PureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	je := json.NewEncoder(w)
	je.SetEscapeHTML(false)
	return je.Encode(r.Data)
}

// WriteContentType (PureJSON) writes custom ContentType.
func (r PureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
