package httpx

import (
	"net/http"
	"net/url"
	"strings"
)

// AppendPrefix returns a handler that serves HTTP requests by appending the
// given prefix from the request URL's Path (and RawPath if set) and invoking
// the handler hh.
func AppendPrefix(prefix string, hh http.Handler) http.Handler {
	if prefix == "" {
		return hh
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := prefix + r.URL.Path
		rp := prefix + r.URL.RawPath
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}

// URLReplace returns a handler that serves HTTP requests by replacing the
// request URL's Path (and RawPath if set) (use strings.Replace(path, src, des) and invoking
// the handler hh.
func URLReplace(src, des string, hh http.Handler) http.Handler {
	if src == "" || src == des {
		return hh
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.Replace(r.URL.Path, src, des, 1)
		rp := strings.Replace(r.URL.RawPath, src, des, 1)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}
