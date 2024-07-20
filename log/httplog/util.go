package httplog

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
)

const (
	logTimeFormat = "2006-01-02T15:04:05.000"
)

func TraceHttpRequest(logger log.Logger, req *http.Request) (rid uint64) {
	if logger != nil && logger.IsTraceEnabled() {
		rid = rand.Uint64() //nolint: gosec
		bs, _ := httputil.DumpRequestOut(req, true)
		logger.Tracef(">>>>>>>> %s %016x >>>>>>>>", time.Now().Format(logTimeFormat), rid)
		logger.Trace(str.UnsafeString(bs))
	}
	return
}

func TraceHttpResponse(logger log.Logger, res *http.Response, rid uint64) {
	if logger != nil && logger.IsTraceEnabled() {
		bs, _ := httputil.DumpResponse(res, true)
		logger.Tracef("<<<<<<<< %s %016x <<<<<<<<", time.Now().Format(logTimeFormat), rid)
		logger.Trace(str.UnsafeString(bs))
	}
}
