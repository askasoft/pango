package httplog

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/ran"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

const (
	logTimeFormat = "2006-01-02T15:04:05.000Z07:00"
)

func TraceClientDo(logger log.Logger, hc *http.Client, req *http.Request) (*http.Response, error) {
	if logger == nil {
		return hc.Do(req)
	}

	if logger.IsTraceEnabled() {
		rid := ran.RandInt63()
		bso, _ := httputil.DumpRequestOut(req, true)

		st := time.Now()
		logger.Tracef(">>>>>>>> %s %016x >>>>>>>>\n%s", st.Format(logTimeFormat), rid, str.UnsafeString(bso))

		res, err := hc.Do(req)
		et := time.Now()

		if err != nil {
			logger.Debugf("[%s] %s %s - %v", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, err)
			return res, err
		}

		bsr, _ := httputil.DumpResponse(res, true)
		logger.Tracef("<<<<<<<< %s %016x <<<<<<<<\n%s", et.Format(logTimeFormat), rid, str.UnsafeString(bsr))

		logger.Debugf("[%s] %s %s - %s (%s)", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, res.Status, num.HumanSize(res.ContentLength))
		return res, err
	}

	if logger.IsDebugEnabled() {
		st := time.Now()
		res, err := hc.Do(req)
		et := time.Now()

		if err != nil {
			logger.Debugf("[%s] %s %s - %v", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, err)
			return res, err
		}

		logger.Debugf("[%s] %s %s - %s (%s)", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, res.Status, num.HumanSize(res.ContentLength))
		return res, err
	}

	return hc.Do(req)
}
