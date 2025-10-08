package httplog

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/askasoft/pango/gog"
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
	if logger == nil || !logger.IsWarnEnabled() {
		return hc.Do(req)
	}

	st := time.Now()

	if logger.IsTraceEnabled() {
		rid := ran.RandInt63()
		bso, _ := httputil.DumpRequestOut(req, true)
		logger.Tracef(">>>>>>>> %s %016x >>>>>>>>\n%s", st.Format(logTimeFormat), rid, str.UnsafeString(bso))

		res, err := hc.Do(req)
		if err != nil {
			logger.Warnf("%s %s - %v", req.Method, req.URL, err)
			return res, err
		}

		et := time.Now()

		bsr, _ := httputil.DumpResponse(res, true)
		logger.Tracef("<<<<<<<< %s %016x <<<<<<<<\n%s", et.Format(logTimeFormat), rid, str.UnsafeString(bsr))

		lvl := gog.If(res.StatusCode >= 400, log.LevelWarn, log.LevelDebug)
		logger.Logf(lvl, "%s %s - %s (%s) [%s]", req.Method, req.URL, res.Status, tmu.HumanDuration(et.Sub(st)), num.HumanSize(res.ContentLength))
		return res, err
	}

	res, err := hc.Do(req)
	if err != nil {
		logger.Warnf("%s %s - %v", req.Method, req.URL, err)
		return res, err
	}

	et := time.Now()

	if res.StatusCode >= 400 {
		if res.StatusCode == 404 {
			logger.Warnf("%s %s - %s (%s) [%s]", req.Method, req.URL, res.Status, tmu.HumanDuration(et.Sub(st)), num.HumanSize(res.ContentLength))
		} else {
			bs, _ := httputil.DumpResponse(res, true)
			logger.Warnf("%s %s - %s (%s) [%s]\n%s", req.Method, req.URL, res.Status, tmu.HumanDuration(et.Sub(st)), num.HumanSize(res.ContentLength), str.UnsafeString(bs))
		}
		return res, err
	}

	logger.Debugf("%s %s - %s (%s) [%s]", req.Method, req.URL, res.Status, tmu.HumanDuration(et.Sub(st)), num.HumanSize(res.ContentLength))
	return res, err
}
