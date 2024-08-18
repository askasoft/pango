package httpx

import (
	"net/http"
)

// IsStatusClientError check status is client side error (400-499)
func IsStatusClientError(status int) bool {
	return status >= 400 && status <= 499
}

// IsStatusServerError check status is server side error (500-599)
func IsStatusServerError(status int) bool {
	return status >= 500 && status <= 599
}

// NoRedirect just return http.ErrUseLastResponse. set http.Client.CheckRedirect to disable auto redirect.
// Example: http.Client{ChecRedirect: httpx.NoRedirect }
func NoRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
