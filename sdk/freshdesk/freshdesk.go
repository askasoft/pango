package freshdesk

import (
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

	RetryOnRateLimited int
}

const (
	contentTypeJSON = `application/json; charset="utf-8"`
	logTimeFormat   = "2006-01-02T15:04:05.000"
)

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

func (fd *Freshdesk) logRequest(req *http.Request) (rid uint64) {
	if fd.Logger != nil && fd.Logger.IsTraceEnabled() {
		rid = rand.Uint64() //nolint: gosec
		bs, _ := httputil.DumpRequestOut(req, true)
		fd.Logger.Tracef(">>>>>>>> %s %016x >>>>>>>>", time.Now().Format(logTimeFormat), rid)
		fd.Logger.Trace(bye.UnsafeString(bs))
	}
	return
}

func (fd *Freshdesk) logResponse(res *http.Response, rid uint64) {
	if fd.Logger != nil && fd.Logger.IsTraceEnabled() {
		bs, _ := httputil.DumpResponse(res, true)
		fd.Logger.Tracef("<<<<<<<< %s %016x <<<<<<<<", time.Now().Format(logTimeFormat), rid)
		fd.Logger.Trace(bye.UnsafeString(bs))
	}
}

func (fd *Freshdesk) call(req *http.Request) (res *http.Response, err error) {
	err = fd.SleepAndRetry(func() error {
		fd.Logger.Infof("%s %s", req.Method, req.URL)

		fd.authenticate(req)
		rid := fd.logRequest(req)

		client := http.Client{
			Transport: fd.Transport,
			Timeout:   fd.Timeout,
		}

		res, err = client.Do(req)
		if err != nil {
			return err
		}
		fd.logResponse(res, rid)

		if res.StatusCode == http.StatusTooManyRequests {
			s := res.Header.Get("Retry-After")
			n, _ := strconv.Atoi(s)
			if n <= 0 {
				n = 60 // invalid number, default to 60s
			}
			iox.DrainAndClose(res.Body)
			return &RateLimitedError{StatusCode: res.StatusCode, RetryAfter: n}
		}

		return err
	}, fd.RetryOnRateLimited)

	return
}

func (fd *Freshdesk) doCall(req *http.Request, result any) error {
	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return decodeResponse(res, result)
}

func (fd *Freshdesk) doGet(url string, result any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return fd.doCall(req, result)
}

func (fd *Freshdesk) doList(url string, lo ListOption, ap any) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	if lo != nil {
		q := lo.Values()
		req.URL.RawQuery = q.Encode()
	}

	res, err := fd.call(req)
	if err != nil {
		return false, err
	}

	if err := decodeResponse(res, ap); err != nil {
		return false, err
	}

	next := res.Header.Get("Link") != ""
	return next, nil
}

func (fd *Freshdesk) doPost(url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	return fd.doCall(req, result)
}

func (fd *Freshdesk) doPut(url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	return fd.doCall(req, result)
}

func (fd *Freshdesk) doDelete(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return fd.doCall(req, nil)
}

// SleepForRetry if err is RateLimitedError, sleep Retry-After and return true
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

func (fd *Freshdesk) SleepAndRetry(api func() error, maxRetry int) (err error) {
	for i := 0; ; i++ {
		err = api()
		if i >= maxRetry {
			break
		}
		if !fd.SleepForRetry(err) {
			break
		}
	}
	return err
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fd *Freshdesk) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("%s/helpdesk/attachments/%d", fd.Domain, aid)
}

func (fd *Freshdesk) GetJob(jid string) (*Job, error) {
	url := fmt.Sprintf("%s/api/v2/jobs/%s", fd.Domain, jid)
	job := &Job{}
	err := fd.doGet(url, job)
	return job, err
}

func (fd *Freshdesk) CreateTicket(ticket *Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets", fd.Domain)
	result := &Ticket{}
	err := fd.doPost(url, ticket, result)
	return result, err
}

// GetTicket Get a Ticket
// include: conversations, requester, company, stats
func (fd *Freshdesk) GetTicket(tid int64, include ...string) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	if len(include) > 0 {
		s := strings.Join(include, ",")
		url += "?include=" + s
	}

	ticket := &Ticket{}
	err := fd.doGet(url, ticket)
	return ticket, err
}

func (fd *Freshdesk) ListTickets(lto *ListTicketsOption) ([]*Ticket, bool, error) {
	url := fmt.Sprintf("%s/api/v2/tickets", fd.Domain)
	tickets := []*Ticket{}
	next, err := fd.doList(url, lto, &tickets)
	return tickets, next, err
}

func (fd *Freshdesk) IterTickets(lto *ListTicketsOption, itf func(*Ticket) bool) error {
	if lto == nil {
		lto = &ListTicketsOption{}
	}
	if lto.Page < 1 {
		lto.Page = 1
	}
	if lto.PerPage < 1 {
		lto.PerPage = 100
	}

	for {
		tickets, next, err := fd.ListTickets(lto)
		if err != nil {
			return err
		}
		for _, t := range tickets {
			if !itf(t) {
				return nil
			}
		}
		if !next {
			break
		}
		lto.Page++
	}
	return nil
}

// FilterTickets
// Use custom ticket fields that you have created in your account to filter through the tickets and get a list of tickets matching the specified ticket fields.
// Query Format: "(ticket_field:integer OR ticket_field:'string') AND ticket_field:boolean"
func (fd *Freshdesk) FilterTickets(fto *FilterTicketsOption) ([]*Ticket, bool, error) {
	url := fmt.Sprintf("%s/api/v2/search/tickets", fd.Domain)
	tickets := []*Ticket{}
	next, err := fd.doList(url, fto, &tickets)
	return tickets, next, err
}

func (fd *Freshdesk) UpdateTicket(tid int64, ticket *Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	result := &Ticket{}
	err := fd.doPut(url, ticket, result)
	return result, err
}

// BulkUpdateTickets returns job id
func (fd *Freshdesk) BulkUpdateTickets(tids []int64, properties *TicketProperties) (string, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/bulk_update", fd.Domain)
	data := map[string]any{
		"bulk_action": map[string]any{
			"ids":        tids,
			"properties": properties,
		},
	}
	result := map[string]string{}
	err := fd.doPut(url, data, &result)
	return result["job_id"], err
}

func (fd *Freshdesk) ForwardTicket(tid int64, tf *TicketForward) (*ForwardResult, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/forward", fd.Domain, tid)
	result := &ForwardResult{}
	err := fd.doPost(url, tf, result)
	return result, err
}

// MergeTickets
// Sometimes, a customer might try to get your attention regarding a particular issue by contacting you through separate channels.
// Sometimes, the same issue might be reported by different people in the team or someone might accidentally open a new ticket instead of following up on an existing one.
// To avoid conflicts, you can merge all related tickets together and keep the communication streamlined.
func (fd *Freshdesk) MergeTickets(tm *TicketsMerge) error {
	url := fmt.Sprintf("%s/api/v2/tickets/merge", fd.Domain)
	err := fd.doPut(url, tm, nil)
	return err
}

func (fd *Freshdesk) ListTicketWatchers(tid int64) ([]int64, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/watchers", fd.Domain, tid)
	result := &TicketWatchers{}
	err := fd.doGet(url, result)
	return result.WatcherIDs, err
}

func (fd *Freshdesk) AddTicketWatcher(tid, uid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/watchers", fd.Domain, tid)
	data := map[string]any{
		"user_id": uid,
	}
	return fd.doPost(url, data, nil)
}

func (fd *Freshdesk) UnwatchTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/unwatch", fd.Domain, tid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) BulkWatchTickets(tids []int64, uid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/buld_watch", fd.Domain)
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.doPut(url, data, nil)
}

func (fd *Freshdesk) BulkUnwatchTickets(tids []int64, uid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/buld_unwatch", fd.Domain)
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.doPut(url, data, nil)
}

func (fd *Freshdesk) RestoreTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/restore", fd.Domain, tid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) DeleteTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) BulkDeleteTickets(tids []int64) (string, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/bulk_delete", fd.Domain)
	data := map[string]any{
		"bulk_action": map[string]any{
			"ids": tids,
		},
	}
	result := map[string]string{}
	err := fd.doPut(url, data, &result)
	return result["job_id"], err
}

func (fd *Freshdesk) DeleteAttachment(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/attachments/%d", fd.Domain, aid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateReply(tid int64, reply *Reply) (*Reply, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/reply", fd.Domain, tid)
	result := &Reply{}
	err := fd.doPost(url, reply, result)
	return result, err
}

func (fd *Freshdesk) CreateNote(tid int64, note *Note) (*Note, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/notes", fd.Domain, tid)
	result := &Note{}
	err := fd.doPost(url, note, result)
	return result, err
}

// UpdateConversation only public & private notes can be edited.
func (fd *Freshdesk) UpdateConversation(cid int64, conversation *Conversation) (*Conversation, error) {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)
	result := &Conversation{}
	err := fd.doPut(url, conversation, result)
	return result, err
}

// DeleteConversation delete a conversation (Incoming Reply can not be deleted)
func (fd *Freshdesk) DeleteConversation(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) ReplyToForward(tid int64, rf *ReplyForward) (*ForwardResult, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/reply_to_forward", fd.Domain, tid)
	result := &ForwardResult{}
	err := fd.doPost(url, rf, result)
	return result, err
}

func (fd *Freshdesk) GetAgent(aid int64) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fd.Domain, aid)
	agent := &Agent{}
	err := fd.doGet(url, agent)
	return agent, err
}

func (fd *Freshdesk) ListAgents(lao *ListAgentsOption) ([]*Agent, bool, error) {
	url := fmt.Sprintf("%s/api/v2/agents", fd.Domain)
	agents := []*Agent{}
	next, err := fd.doList(url, lao, &agents)
	return agents, next, err
}

func (fd *Freshdesk) IterAgents(lao *ListAgentsOption, iaf func(*Agent) bool) error {
	if lao == nil {
		lao = &ListAgentsOption{}
	}
	if lao.Page < 1 {
		lao.Page = 1
	}
	if lao.PerPage < 1 {
		lao.PerPage = 100
	}

	for {
		agents, next, err := fd.ListAgents(lao)
		if err != nil {
			return err
		}
		for _, c := range agents {
			if !iaf(c) {
				return nil
			}
		}
		if !next {
			break
		}
		lao.Page++
	}
	return nil
}

func (fd *Freshdesk) CreateAgent(agent *AgentRequest) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents", fd.Domain)
	result := &Agent{}
	err := fd.doPost(url, agent, result)
	return result, err
}

func (fd *Freshdesk) UpdateAgent(aid int64, agent *AgentRequest) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fd.Domain, aid)
	result := &Agent{}
	err := fd.doPut(url, agent, result)
	return result, err
}

func (fd *Freshdesk) DeleteAgent(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fd.Domain, aid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) GetCurrentAgent() (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/me", fd.Domain)
	agent := &Agent{}
	err := fd.doGet(url, agent)
	return agent, err
}

func (fd *Freshdesk) SearchAgents(keyword string) ([]*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/autocomplete?term=%s", fd.Domain, url.QueryEscape(keyword))
	agents := []*Agent{}
	err := fd.doGet(url, &agents)
	return agents, err
}

func (fd *Freshdesk) CreateContact(contact *Contact) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts", fd.Domain)
	result := &Contact{}
	err := fd.doPost(url, contact, result)
	return result, err
}

func (fd *Freshdesk) UpdateContact(cid int64, contact *Contact) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/%d", fd.Domain, cid)
	result := &Contact{}
	err := fd.doPut(url, contact, result)
	return result, err
}

func (fd *Freshdesk) GetContact(cid int64) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/%d", fd.Domain, cid)
	contact := &Contact{}
	err := fd.doGet(url, contact)
	return contact, err
}

func (fd *Freshdesk) DeleteContact(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d", fd.Domain, cid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) HardDeleteContact(cid int64, force ...bool) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/hard_delete", fd.Domain, cid)
	if len(force) > 0 && force[0] {
		url += "?force=true"
	}
	return fd.doDelete(url)
}

func (fd *Freshdesk) ListContacts(lco *ListContactsOption) ([]*Contact, bool, error) {
	url := fmt.Sprintf("%s/api/v2/contacts", fd.Domain)
	contacts := []*Contact{}
	next, err := fd.doList(url, lco, &contacts)
	return contacts, next, err
}

func (fd *Freshdesk) IterContacts(lco *ListContactsOption, itf func(*Contact) bool) error {
	if lco == nil {
		lco = &ListContactsOption{}
	}
	if lco.Page < 1 {
		lco.Page = 1
	}
	if lco.PerPage < 1 {
		lco.PerPage = 100
	}

	for {
		contacts, next, err := fd.ListContacts(lco)
		if err != nil {
			return err
		}
		for _, c := range contacts {
			if !itf(c) {
				return nil
			}
		}
		if !next {
			break
		}
		lco.Page++
	}
	return nil
}

func (fd *Freshdesk) SearchContacts(keyword string) ([]*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/autocomplete?term=%s", fd.Domain, url.QueryEscape(keyword))
	contacts := []*Contact{}
	err := fd.doGet(url, &contacts)
	return contacts, err
}

func (fd *Freshdesk) RestoreContact(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/restore", fd.Domain, cid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) InviteContact(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/send_invite", fd.Domain, cid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) MergeContacts(cm *ContactsMerge) error {
	url := fmt.Sprintf("%s/api/v2/contacts/merge", fd.Domain)
	return fd.doPost(url, nil, nil)
}

// ExportContacts return a job id, call GetExportedContactsURL() to get the job detail
func (fd *Freshdesk) ExportContacts(defaultFields, customFields []string) (string, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/export", fd.Domain)
	data := map[string]any{
		"fields": &ContactsExport{defaultFields, customFields},
	}
	result := map[string]string{}
	err := fd.doPost(url, data, &result)
	return result["id"], err
}

// GetExportedContactsURL get the exported contacts url
func (fd *Freshdesk) GetExportedContactsURL(jid string) (*Job, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/export/%s", fd.Domain, jid)
	job := &Job{}
	err := fd.doGet(url, job)
	return job, err
}

func (fd *Freshdesk) MakeAgent(cid int64, agent *Agent) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/make_agent", fd.Domain, cid)
	result := &Contact{}
	err := fd.doPut(url, agent, result)
	return result, err
}

func (fd *Freshdesk) CreateCategory(category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fd.Domain)
	result := &Category{}
	err := fd.doPost(url, category, result)
	return result, err
}

func (fd *Freshdesk) CreateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	result := &Category{}
	err := fd.doPost(url, category, result)
	return result, err
}

func (fd *Freshdesk) UpdateCategory(cid int64, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fd.Domain, cid)
	result := &Category{}
	err := fd.doPut(url, category, result)
	return result, err
}

func (fd *Freshdesk) UpdateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	result := &Category{}
	err := fd.doPut(url, category, result)
	return result, err
}

func (fd *Freshdesk) GetCategory(cid int64) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fd.Domain, cid)
	cat := &Category{}
	err := fd.doGet(url, cat)
	return cat, err
}

func (fd *Freshdesk) GetCategoryTranslated(cid int64, lang string) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	cat := &Category{}
	err := fd.doGet(url, cat)
	return cat, err
}

func (fd *Freshdesk) ListCategories() ([]*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fd.Domain)
	categories := []*Category{}
	err := fd.doGet(url, &categories)
	return categories, err
}

func (fd *Freshdesk) ListCategoriesTranslated(lang string) ([]*Category, error) {
	url := fd.Domain + "/api/v2/solutions/categories/" + lang
	categories := []*Category{}
	err := fd.doGet(url, &categories)
	return categories, err
}

func (fd *Freshdesk) DeleteCategory(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fd.Domain, cid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateFolder(cid int64, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/folders", fd.Domain, cid)
	result := &Folder{}
	err := fd.doPost(url, folder, result)
	return result, err
}

func (fd *Freshdesk) CreateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Folder{}
	err := fd.doPost(url, folder, result)
	return result, err
}

func (fd *Freshdesk) UpdateFolder(fid int64, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	result := &Folder{}
	err := fd.doPut(url, folder, result)
	return result, err
}

func (fd *Freshdesk) UpdateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Folder{}
	err := fd.doPut(url, folder, result)
	return result, err
}

func (fd *Freshdesk) GetFolder(fid int64) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	folder := &Folder{}
	err := fd.doGet(url, folder)
	return folder, err
}

func (fd *Freshdesk) GetFolderTranslated(fid int64, lang string) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	folder := &Folder{}
	err := fd.doGet(url, folder)
	return folder, err
}

func (fd *Freshdesk) ListCategoryFolders(cid int64) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/folders", fd.Domain, cid)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *Freshdesk) ListCategoryFoldersTranslated(cid int64, lang string) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/folders/%s", fd.Domain, cid, lang)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *Freshdesk) ListSubFolders(fid int64) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/subfolders", fd.Domain, fid)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *Freshdesk) ListSubFoldersTranslated(fid int64, lang string) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/subfolders/%s", fd.Domain, fid, lang)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *Freshdesk) DeleteFolder(fid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateArticle(fid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	result := &Article{}
	err := fd.doPost(url, article, result)
	return result, err
}

func (fd *Freshdesk) CreateArticleTranslated(fid int64, lang string, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Article{}
	err := fd.doPost(url, article, result)
	return result, err
}

func (fd *Freshdesk) UpdateArticle(aid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fd.Domain, aid)
	result := &Article{}
	err := fd.doPut(url, article, result)
	return result, err
}

func (fd *Freshdesk) UpdateArticleTranslated(aid int64, lang string, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/%s", fd.Domain, aid, lang)
	result := &Article{}
	err := fd.doPut(url, article, result)
	return result, err
}

func (fd *Freshdesk) GetArticle(aid int64) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fd.Domain, aid)
	article := &Article{}
	err := fd.doGet(url, article)
	return article, err
}

func (fd *Freshdesk) GetArticleTranslated(aid int64, lang string) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/%s", fd.Domain, aid, lang)
	article := &Article{}
	err := fd.doGet(url, article)
	return article, err
}

func (fd *Freshdesk) ListFolderArticles(fid int64) ([]*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/articles", fd.Domain, fid)
	articles := []*Article{}
	err := fd.doGet(url, &articles)
	return articles, err
}

func (fd *Freshdesk) ListFolderArticlesTranslated(fid int64, lang string) ([]*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/farticles/%s", fd.Domain, fid, lang)
	articles := []*Article{}
	err := fd.doGet(url, &articles)
	return articles, err
}

func (fd *Freshdesk) DeleteArticle(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fd.Domain, aid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) SearchArticles(keyword string) ([]*ArticleEx, error) {
	url := fmt.Sprintf("%s/api/v2/search/solutions?term=%s", fd.Domain, url.QueryEscape(keyword))
	articles := []*ArticleEx{}
	err := fd.doGet(url, &articles)
	return articles, err
}

func (fd *Freshdesk) GetRole(rid int64) (*Role, error) {
	url := fmt.Sprintf("%s/api/v2/roles/%d", fd.Domain, rid)
	role := &Role{}
	err := fd.doGet(url, role)
	return role, err
}

func (fd *Freshdesk) ListRoles() ([]*Role, error) {
	url := fmt.Sprintf("%s/api/v2/roles", fd.Domain)
	roles := []*Role{}
	_, err := fd.doList(url, nil, &roles)
	return roles, err
}

func (fd *Freshdesk) GetGroup(gid int64) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	group := &Group{}
	err := fd.doGet(url, group)
	return group, err
}

func (fd *Freshdesk) CreateGroup(group *Group) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fd.Domain)
	result := &Group{}
	err := fd.doPost(url, group, result)
	return result, err
}

func (fd *Freshdesk) ListGroups() ([]*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fd.Domain)
	groups := []*Group{}
	_, err := fd.doList(url, nil, &groups)
	return groups, err
}

func (fd *Freshdesk) UpdateGroup(gid int64, group *Group) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	result := &Group{}
	err := fd.doPut(url, group, result)
	return result, err
}

func (fd *Freshdesk) DeleteGroup(gid int64) error {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	return fd.doDelete(url)
}
