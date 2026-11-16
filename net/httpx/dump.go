package httpx

import (
	"net/http"
	"net/http/httputil"
)

func redactHeader(h http.Header) {
	RedactHeader(h, SensitiveHeader, RedactedValue)
}

// SafeDumpRequest redacts matching headers, and dumps it.
func SafeDumpRequest(req *http.Request, body bool) ([]byte, error) {
	hh := req.Header

	defer func() {
		req.Header = hh
	}()

	req.Header = hh.Clone()
	redactHeader(req.Header)

	return httputil.DumpRequest(req, body)
}

// SafeDumpRequestOut redacts matching headers, and dumps it.
func SafeDumpRequestOut(req *http.Request, body bool) ([]byte, error) {
	hh := req.Header

	defer func() {
		req.Header = hh
	}()

	req.Header = hh.Clone()
	redactHeader(req.Header)

	return httputil.DumpRequestOut(req, body)
}

// SafeDumpResponse redacts matching headers, and dumps it.
func SafeDumpResponse(res *http.Response, body bool) ([]byte, error) {
	hh := res.Header

	defer func() {
		res.Header = hh
	}()

	res.Header = hh.Clone()
	redactHeader(res.Header)

	return httputil.DumpResponse(res, body)
}
