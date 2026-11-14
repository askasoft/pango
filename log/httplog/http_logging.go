package httplog

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/ran"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

const (
	logTimeFormat = "2006-01-02T15:04:05.000Z07:00"
)

type RoundTripperLogger struct {
	Logger       log.Logger
	Transport    http.RoundTripper
	RedactHeader func(http.Header)
}

func LoggingRoundTripper(logger log.Logger, transport ...http.RoundTripper) *RoundTripperLogger {
	return &RoundTripperLogger{
		Logger:       logger,
		Transport:    asg.First(transport, http.DefaultTransport),
		RedactHeader: redactHeader,
	}
}

func (lrt *RoundTripperLogger) do(req *http.Request) (*http.Response, error) {
	return lrt.Transport.RoundTrip(req)
}

func (lrt *RoundTripperLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	if logger := lrt.Logger; logger != nil {
		if logger.IsTraceEnabled() {
			return lrt.traceDo(logger, req)
		}

		if logger.IsDebugEnabled() {
			return lrt.debugDo(logger, req)
		}
	}

	return lrt.do(req)
}

func (lrt *RoundTripperLogger) debugDo(logger log.Logger, req *http.Request) (*http.Response, error) {
	st := time.Now()
	res, err := lrt.do(req)
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

func (lrt *RoundTripperLogger) traceDo(logger log.Logger, req *http.Request) (*http.Response, error) {
	rid := ran.RandInt63()
	bso, _ := lrt.dumpRequest(req)

	st := time.Now()
	logger.Tracef(">>>>>>>> %s %016x >>>>>>>>\n%s", st.Format(logTimeFormat), rid, str.UnsafeString(bso))

	res, err := lrt.do(req)
	et := time.Now()

	var s string
	if err != nil {
		s = fmt.Sprintf("%s %s - %v", req.Method, req.URL, err)
	} else {
		bsr, _ := lrt.dumpResponse(res)
		s = str.UnsafeString(bsr)
	}
	logger.Tracef("<<<<<<<< %s %016x <<<<<<<< [%s]\n%s", et.Format(logTimeFormat), rid, tmu.HumanDuration(et.Sub(st)), s)

	return res, err
}

func (lrt *RoundTripperLogger) dumpRequest(req *http.Request) ([]byte, error) {
	hh := req.Header

	defer func() {
		req.Header = hh
	}()

	if redact := lrt.RedactHeader; redact != nil {
		req.Header = hh.Clone()
		redact(req.Header)
	}
	return httputil.DumpRequestOut(req, true)
}

func (lrt *RoundTripperLogger) dumpResponse(res *http.Response) ([]byte, error) {
	hh := res.Header

	defer func() {
		res.Header = hh
	}()

	if redact := lrt.RedactHeader; redact != nil {
		res.Header = hh.Clone()
		redact(res.Header)
	}

	return httputil.DumpResponse(res, true)
}

func redactHeader(hh http.Header) {
	httpx.RedactHeader(hh, httpx.SensitiveHeader, httpx.RedactedValue)
}
