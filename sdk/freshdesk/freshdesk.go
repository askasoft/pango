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

func (fd *Freshdesk) call(req *http.Request) (*http.Response, error) {
	fd.authenticate(req)
	rid := fd.logRequest(req)

	client := http.Client{
		Transport: fd.Transport,
		Timeout:   fd.Timeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fd.logResponse(res, rid)

	if res.StatusCode == http.StatusTooManyRequests {
		s := res.Header.Get("Retry-After")
		n, _ := strconv.Atoi(s)
		if n <= 0 {
			n = 60 // invalid number, default to 60s
		}
		iox.DrainAndClose(res.Body)
		return res, &RateLimitedError{StatusCode: res.StatusCode, RetryAfter: n}
	}

	return res, err
}

func (fd *Freshdesk) doGet(url string, obj any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return decodeResponse(res, http.StatusOK, obj)
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

	if err := decodeResponse(res, http.StatusOK, ap); err != nil {
		return false, err
	}

	next := res.Header.Get("Link") != ""
	return next, nil
}

func (fd *Freshdesk) doCreate(url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", ct)

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return decodeResponse(res, http.StatusCreated, result)
}

func (fd *Freshdesk) doUpdate(url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", ct)

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return decodeResponse(res, http.StatusOK, result)
}

func (fd *Freshdesk) doDelete(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	res, err := fd.call(req)
	if err != nil {
		return err
	}

	return decodeResponse(res, http.StatusNoContent, nil)
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
	url := fmt.Sprintf("%s/api/v2/tickets", fd.Domain)
	result := &Ticket{}
	err := fd.doCreate(url, ticket, result)
	return result, err
}

func (fd *Freshdesk) UpdateTicket(tid int64, ticket *Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	result := &Ticket{}
	err := fd.doUpdate(url, ticket, result)
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

	return decodeResponse(res, http.StatusNoContent, nil)
}

func (fd *Freshdesk) DeleteTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) DeleteAttachment(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/attachments/%d", fd.Domain, aid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateReply(tid int64, reply *Reply) (*Reply, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/reply", fd.Domain, tid)
	result := &Reply{}
	err := fd.doCreate(url, reply, result)
	return result, err
}

func (fd *Freshdesk) CreateNote(tid int64, note *Note) (*Note, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/notes", fd.Domain, tid)
	result := &Note{}
	err := fd.doCreate(url, note, result)
	return result, err
}

func (fd *Freshdesk) UpdateConversation(cid int64, conversation *Conversation) (*Conversation, error) {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)
	result := &Conversation{}
	err := fd.doUpdate(url, conversation, result)
	return result, err
}

func (fd *Freshdesk) DeleteConversation(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)
	return fd.doDelete(url)
}

type ListTicketsOption struct {
	Filter           string // The various filters available are new_and_my_open, watching, spam, deleted.
	RequestID        int64
	Email            string
	UniqueExternalID string
	CompanyID        int64
	UpdatedSince     Time
	Include          string // stats, requester, description
	OrderBy          string // created_at, due_by, updated_at, status
	OrderType        string // asc, desc (default)
	Page             int
	PerPage          int
}

func (lto *ListTicketsOption) Values() Values {
	q := Values{}
	q.SetString("filter", lto.Filter)
	q.SetInt64("request_id", lto.RequestID)
	q.SetString("email", lto.Email)
	q.SetString("unique_external_id", lto.UniqueExternalID)
	q.SetInt64("company_id", lto.CompanyID)
	q.SetTime("updated_since", lto.UpdatedSince)
	q.SetString("include", lto.Include)
	q.SetString("order_by", lto.OrderBy)
	q.SetString("order_type", lto.OrderType)
	q.SetInt("page", lto.Page)
	q.SetInt("per_page", lto.PerPage)
	return q
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
		if fd.SleepForRetry(err) {
			continue
		}
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

func (fd *Freshdesk) CreateContact(contact *Contact) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts", fd.Domain)
	result := &Contact{}
	err := fd.doCreate(url, contact, result)
	return result, err
}

func (fd *Freshdesk) UpdateContact(cid int64, contact *Contact) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/%d", fd.Domain, cid)
	result := &Contact{}
	err := fd.doUpdate(url, contact, result)
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

type ListContactsOption struct {
	Email        string
	Mobile       string
	Phone        string
	CompanyID    int64
	UpdatedSince Time
	State        string // [blocked/deleted/unverified/verified]
	Page         int
	PerPage      int
}

func (lco *ListContactsOption) Values() Values {
	q := Values{}
	q.SetString("email", lco.Email)
	q.SetString("mobile", lco.Mobile)
	q.SetString("phone", lco.Phone)
	q.SetInt64("company_id", lco.CompanyID)
	q.SetString("state", lco.State)
	q.SetTime("updated_since", lco.UpdatedSince)
	q.SetInt("page", lco.Page)
	q.SetInt("per_page", lco.PerPage)
	return q
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
		if fd.SleepForRetry(err) {
			continue
		}
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

func (fd *Freshdesk) CreateCategory(category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fd.Domain)
	result := &Category{}
	err := fd.doCreate(url, category, result)
	return result, err
}

func (fd *Freshdesk) CreateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	result := &Category{}
	err := fd.doCreate(url, category, result)
	return result, err
}

func (fd *Freshdesk) UpdateCategory(cid int64, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fd.Domain, cid)
	result := &Category{}
	err := fd.doUpdate(url, category, result)
	return result, err
}

func (fd *Freshdesk) UpdateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	result := &Category{}
	err := fd.doUpdate(url, category, result)
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
	err := fd.doCreate(url, folder, result)
	return result, err
}

func (fd *Freshdesk) CreateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Folder{}
	err := fd.doCreate(url, folder, result)
	return result, err
}

func (fd *Freshdesk) UpdateFolder(fid int64, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	result := &Folder{}
	err := fd.doUpdate(url, folder, result)
	return result, err
}

func (fd *Freshdesk) UpdateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Folder{}
	err := fd.doUpdate(url, folder, result)
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
	err := fd.doCreate(url, article, result)
	return result, err
}

func (fd *Freshdesk) CreateArticleTranslated(fid int64, lang string, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Article{}
	err := fd.doCreate(url, article, result)
	return result, err
}

func (fd *Freshdesk) UpdateArticle(aid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fd.Domain, aid)
	result := &Article{}
	err := fd.doUpdate(url, article, result)
	return result, err
}

func (fd *Freshdesk) UpdateArticleTranslated(aid int64, lang string, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/%s", fd.Domain, aid, lang)
	result := &Article{}
	err := fd.doUpdate(url, article, result)
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

func (fd *Freshdesk) GetArticleAttachmentURL(aid int64) string {
	return fmt.Sprintf("%s/helpdesk/attachments/%d", fd.Domain, aid)
}
