package freshdesk

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pandafw/pango/bye"
	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/log"
)

type Freshdesk struct {
	Domain   string
	Apikey   string
	Username string
	Password string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger
}

const contentTypeJSON = `application/json; charset="utf-8"`

func (fd *Freshdesk) authenticate(req *http.Request) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeJSON)
	}

	if fd.Apikey != "" {
		req.SetBasicAuth(fd.Apikey, "X")
	} else {
		req.SetBasicAuth(fd.Username, fd.Password)
	}
}

func (fd *Freshdesk) logRequest(req *http.Request) (lid int64) {
	if fd.Logger != nil && fd.Logger.IsTraceEnabled() {
		lid = rand.Int63() //nolint: gosec
		bs, _ := httputil.DumpRequestOut(req, true)
		fd.Logger.Tracef(">>>>>>>> %s %08x >>>>>>>>", time.Now().Format("2006-01-02T15:04:05.000"), lid)
		fd.Logger.Trace(bye.UnsafeString(bs))
	}
	return
}

func (fd *Freshdesk) logResponse(res *http.Response, lid int64) {
	if fd.Logger != nil && fd.Logger.IsTraceEnabled() {
		bs, _ := httputil.DumpResponse(res, true)
		fd.Logger.Tracef("<<<<<<<< %s %08x <<<<<<<<", time.Now().Format("2006-01-02T15:04:05.000"), lid)
		fd.Logger.Trace(bye.UnsafeString(bs))
	}
}

func (fd *Freshdesk) call(req *http.Request) (*http.Response, error) {
	fd.authenticate(req)
	lid := fd.logRequest(req)

	client := http.Client{
		Transport: fd.Transport,
		Timeout:   fd.Timeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fd.logResponse(res, lid)

	if res.StatusCode == http.StatusTooManyRequests {
		s := res.Header.Get("Retry-After")
		n, _ := strconv.Atoi(s)
		if n <= 0 {
			n = 60 // invalid number, default to 60s
		}
		iox.DrainAndClose(res.Body)
		return res, &RateLimitedError{n}
	}

	return res, err
}

func (fd *Freshdesk) decode(res *http.Response, status int, obj any) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == status {
		if obj != nil {
			return decoder.Decode(obj)
		}
		return nil
	}

	er := &ErrorResult{}
	if err := decoder.Decode(er); err != nil {
		return err
	}
	return er
}

func (fd *Freshdesk) SleepForRetry(err error) bool {
	if err != nil {
		if rle, ok := err.(*RateLimitedError); ok { //nolint: errorlint
			if fd.Logger != nil {
				fd.Logger.Warnf("Sleep %d seconds for API Rate Limited", rle.RetryAfter)
			}
			time.Sleep(time.Duration(rle.RetryAfter) * time.Second)
			return true
		}
	}
	return false
}

func (fd *Freshdesk) CreateTicket(ticket *Ticket) (*Ticket, error) {
	buf, ct, err := ticket.BuildRequest()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fd.Domain+"/api/v2/tickets", buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ct)

	res, err := fd.call(req)
	if err != nil {
		return nil, err
	}

	nt := &Ticket{}
	if err := fd.decode(res, http.StatusCreated, nt); err != nil {
		return nil, err
	}
	return nt, nil
}

func (fd *Freshdesk) UpdateTicket(tid int64, ticket *Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)

	buf, ct, err := ticket.BuildRequest()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ct)

	res, err := fd.call(req)
	if err != nil {
		return nil, err
	}

	nt := &Ticket{}
	if err := fd.decode(res, http.StatusOK, nt); err != nil {
		return nil, err
	}
	return nt, nil
}

// GetTicket Get a Ticket
// include: conversations, requester, company, stats
func (fd *Freshdesk) GetTicket(tid int64, include ...string) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if len(include) > 0 {
		q := req.URL.Query()
		s := strings.Join(include, ",")
		q.Add("include", s)
		req.URL.RawQuery = q.Encode()
	}

	res, err := fd.call(req)
	if err != nil {
		return nil, err
	}

	ticket := &Ticket{}
	if err := fd.decode(res, http.StatusOK, ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (fd *Freshdesk) RestoreTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/restore", fd.Domain, tid)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return fd.decode(res, http.StatusNoContent, nil)
}

func (fd *Freshdesk) DeleteTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return fd.decode(res, http.StatusNoContent, nil)
}

func (fd *Freshdesk) DeleteAttachment(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/attachments/%d", fd.Domain, aid)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return fd.decode(res, http.StatusNoContent, nil)
}

func (fd *Freshdesk) CreateReply(tid int64, r *Reply) (*Reply, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/reply", fd.Domain, tid)

	buf, ct, err := r.BuildRequest()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ct)

	res, err := fd.call(req)
	if err != nil {
		return nil, err
	}

	nr := &Reply{}
	if err := fd.decode(res, http.StatusCreated, nr); err != nil {
		return nil, err
	}
	return nr, nil
}

func (fd *Freshdesk) CreateNote(tid int64, r *Note) (*Note, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/notes", fd.Domain, tid)

	buf, ct, err := r.BuildRequest()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ct)

	res, err := fd.call(req)
	if err != nil {
		return nil, err
	}

	nn := &Note{}
	if err := fd.decode(res, http.StatusCreated, nn); err != nil {
		return nil, err
	}
	return nn, nil
}

func (fd *Freshdesk) UpdateConversation(cid int64, c *Conversation) (*Conversation, error) {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)

	buf, ct, err := c.BuildRequest()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ct)

	res, err := fd.call(req)
	if err != nil {
		return nil, err
	}

	nc := &Conversation{}
	if err := fd.decode(res, http.StatusOK, nc); err != nil {
		return nil, err
	}
	return nc, nil
}

func (fd *Freshdesk) DeleteConversation(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return fd.decode(res, http.StatusNoContent, nil)
}

type ListTicketsOption struct {
	Filter           string // The various filters available are new_and_my_open, watching, spam, deleted.
	RequestID        string
	Email            string
	UniqueExternalID string
	CompanyID        string
	UpdatedSince     *Time
	Include          string // stats, requester, description
	OrderBy          string // created_at, due_by, updated_at, status
	OrderType        string // asc, desc (default)
}

func (lto *ListTicketsOption) BuildQuery() url.Values {
	q := url.Values{}
	if lto.Filter != "" {
		q.Add("filter", lto.Filter)
	}
	if lto.RequestID != "" {
		q.Add("request_id", lto.RequestID)
	}
	if lto.Email != "" {
		q.Add("email", lto.Email)
	}
	if lto.UniqueExternalID != "" {
		q.Add("unique_external_id", lto.UniqueExternalID)
	}
	if lto.CompanyID != "" {
		q.Add("company_id", lto.CompanyID)
	}
	if lto.UpdatedSince != nil {
		q.Add("updated_since", lto.UpdatedSince.String())
	}
	if lto.Include != "" {
		q.Add("include", lto.Include)
	}
	if lto.OrderBy != "" {
		q.Add("order_by", lto.OrderBy)
	}
	if lto.OrderType != "" {
		q.Add("order_type", lto.OrderType)
	}
	return q
}

func (fd *Freshdesk) ListTickets(lto *ListTicketsOption) ([]*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets", fd.Domain)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := lto.BuildQuery()
	req.URL.RawQuery = q.Encode()

	res, err := fd.call(req)
	if err != nil {
		return nil, err
	}

	tickets := []*Ticket{}
	if err := fd.decode(res, http.StatusOK, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}
