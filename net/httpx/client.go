package httpx

import (
	"net/http"
)

// NoRedirect just return http.ErrUseLastResponse. set http.Client.CheckRedirect to disable auto redirect.
// Example: http.Client{ChecRedirect: httpx.NoRedirect }
func NoRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
