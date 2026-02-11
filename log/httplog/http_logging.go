package httplog

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/ran"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

const (
	logTimeFormat = "2006-01-02T15:04:05.000Z07:00"
)

type iDo interface {
	Do(req *http.Request) (*http.Response, error)
}

func traceDo(logger log.Logger, do iDo, req *http.Request) (*http.Response, error) {
	rid := ran.RandInt63()
	bso, _ := httputil.DumpRequestOut(req, true)

	st := time.Now()
	logger.Tracef(">>>>>>>> %s %016x >>>>>>>>\n%s", st.Format(logTimeFormat), rid, str.UnsafeString(bso))

	res, err := do.Do(req)
	et := time.Now()

	if err != nil {
		logger.Debugf("[%s] %s %s - %v", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, err)
		return res, err
	}

	bsr, _ := httputil.DumpResponse(res, true)
	logger.Tracef("<<<<<<<< %s %016x <<<<<<<<\n%s", et.Format(logTimeFormat), rid, str.UnsafeString(bsr))

	if res.ContentLength >= 0 {
		logger.Debugf("[%s] %s %s - %s (%s)", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, res.Status, num.HumanSize(res.ContentLength))
	} else {
		logger.Debugf("[%s] %s %s - %s", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, res.Status)
	}
	return res, err
}

func debugDo(logger log.Logger, do iDo, req *http.Request) (*http.Response, error) {
	st := time.Now()
	res, err := do.Do(req)
	et := time.Now()

	if err != nil {
		logger.Debugf("[%s] %s %s - %v", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, err)
		return res, err
	}

	if res.ContentLength >= 0 {
		logger.Debugf("[%s] %s %s - %s (%s)", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, res.Status, num.HumanSize(res.ContentLength))
	} else {
		logger.Debugf("[%s] %s %s - %s", tmu.HumanDuration(et.Sub(st)), req.Method, req.URL, res.Status)
	}
	return res, err
}

func TraceClientDo(logger log.Logger, hc *http.Client, req *http.Request) (*http.Response, error) {
	if logger == nil {
		return hc.Do(req)
	}

	if logger.IsTraceEnabled() {
		return traceDo(logger, hc, req)
	}

	if logger.IsDebugEnabled() {
		return debugDo(logger, hc, req)
	}

	return hc.Do(req)
}

type LoggingRoundTripper struct {
	Logger    log.Logger
	Transport http.RoundTripper
}

func NewLoggingRoundTripper(logger log.Logger, transport ...http.RoundTripper) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		Logger:    logger,
		Transport: asg.First(transport),
	}
}

func (lrt *LoggingRoundTripper) Do(req *http.Request) (*http.Response, error) {
	return lrt.transport().RoundTrip(req)
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if lrt.Logger != nil {
		if lrt.Logger.IsTraceEnabled() {
			return traceDo(lrt.Logger, lrt, req)
		}

		if lrt.Logger.IsDebugEnabled() {
			return debugDo(lrt.Logger, lrt, req)
		}
	}

	return lrt.Do(req)
}

func (lrt *LoggingRoundTripper) transport() http.RoundTripper {
	if lrt.Transport != nil {
		return lrt.Transport
	}
	return http.DefaultTransport
}
