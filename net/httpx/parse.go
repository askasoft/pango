package httpx

import (
	"bytes"
	"strings"
)

// ParsePairs extracts key/value pairs from a comma-separated list of
// values as described by RFC 2068 and returns a map[key]value. The
// resulting values are unquoted. If a list element doesn't contain a
// "=", the key is the element itself and the value is an empty
// string.
//
// Lifted from https://code.google.com/p/gorilla/source/browse/http/parser/parser.go
func ParsePairs(value string) map[string]string {
	m := make(map[string]string)
	for _, pair := range ParseList(strings.TrimSpace(value)) {
		switch i := strings.Index(pair, "="); {
		case i < 0:
			// No '=' in pair, treat whole string as a 'key'.
			m[pair] = ""
		case i == len(pair)-1:
			// Malformed pair ('key=' with no value), keep key with empty value.
			m[pair[:i]] = ""
		default:
			v := pair[i+1:]
			if v[0] == '"' && v[len(v)-1] == '"' {
				// Unquote it.
				v = v[1 : len(v)-1]
			}
			m[pair[:i]] = v
		}
	}
	return m
}

// ParseList parses a comma-separated list of values as described by
// RFC 2068 and returns list elements.
//
// Lifted from https://code.google.com/p/gorilla/source/browse/http/parser/parser.go
// which was ported from urllib2.parse_http_list, from the Python
// standard library.
func ParseList(value string) []string {
	var list []string
	var escape, quote bool
	b := new(bytes.Buffer)
	for _, r := range value {
		switch {
		case escape:
			b.WriteRune(r)
			escape = false
		case quote:
			if r == '\\' {
				escape = true
			} else {
				if r == '"' {
					quote = false
				}
				b.WriteRune(r)
			}
		case r == ',':
			list = append(list, strings.TrimSpace(b.String()))
			b.Reset()
		case r == '"':
			quote = true
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	// Append last part.
	if s := b.String(); s != "" {
		list = append(list, strings.TrimSpace(s))
	}
	return list
}
