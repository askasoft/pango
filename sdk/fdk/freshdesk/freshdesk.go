package freshdesk

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/askasoft/pango/sdk/fdk"
)

type FreshDesk fdk.FDK

func (fd *FreshDesk) doGet(url string, result any) error {
	return (*fdk.FDK)(fd).DoGet(url, result)
}

func (fd *FreshDesk) doList(url string, lo ListOption, ap any) (bool, error) {
	return (*fdk.FDK)(fd).DoList(url, lo, ap)
}

func (fd *FreshDesk) doPost(url string, source, result any) error {
	return (*fdk.FDK)(fd).DoPost(url, source, result)
}

func (fd *FreshDesk) doPut(url string, source, result any) error {
	return (*fdk.FDK)(fd).DoPut(url, source, result)
}

func (fd *FreshDesk) doDelete(url string) error {
	return (*fdk.FDK)(fd).DoDelete(url)
}

func (fd *FreshDesk) Download(url string) ([]byte, error) {
	return (*fdk.FDK)(fd).DoDownload(url)
}

func (fd *FreshDesk) SaveFile(url string, filename string) error {
	return (*fdk.FDK)(fd).DoSave(url, filename)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fd *FreshDesk) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("%s/helpdesk/attachments/%d", fd.Domain, aid)
}

func (fd *FreshDesk) GetJob(jid string) (*Job, error) {
	url := fmt.Sprintf("%s/api/v2/jobs/%s", fd.Domain, jid)
	job := &Job{}
	err := fd.doGet(url, job)
	return job, err
}

func (fd *FreshDesk) CreateTicket(ticket *Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets", fd.Domain)
	result := &Ticket{}
	err := fd.doPost(url, ticket, result)
	return result, err
}

// GetTicket Get a Ticket
// include: conversations, requester, company, stats
func (fd *FreshDesk) GetTicket(tid int64, include ...string) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	if len(include) > 0 {
		s := strings.Join(include, ",")
		url += "?include=" + s
	}

	ticket := &Ticket{}
	err := fd.doGet(url, ticket)
	return ticket, err
}

func (fd *FreshDesk) ListTickets(lto *ListTicketsOption) ([]*Ticket, bool, error) {
	url := fmt.Sprintf("%s/api/v2/tickets", fd.Domain)
	tickets := []*Ticket{}
	next, err := fd.doList(url, lto, &tickets)
	return tickets, next, err
}

func (fd *FreshDesk) IterTickets(lto *ListTicketsOption, itf func(*Ticket) bool) error {
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
func (fd *FreshDesk) FilterTickets(fto *FilterTicketsOption) ([]*Ticket, bool, error) {
	url := fmt.Sprintf("%s/api/v2/search/tickets", fd.Domain)
	tickets := []*Ticket{}
	next, err := fd.doList(url, fto, &tickets)
	return tickets, next, err
}

func (fd *FreshDesk) UpdateTicket(tid int64, ticket *Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	result := &Ticket{}
	err := fd.doPut(url, ticket, result)
	return result, err
}

// BulkUpdateTickets returns job id
func (fd *FreshDesk) BulkUpdateTickets(tids []int64, properties *TicketProperties) (string, error) {
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

func (fd *FreshDesk) ForwardTicket(tid int64, tf *TicketForward) (*ForwardResult, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/forward", fd.Domain, tid)
	result := &ForwardResult{}
	err := fd.doPost(url, tf, result)
	return result, err
}

// MergeTickets
// Sometimes, a customer might try to get your attention regarding a particular issue by contacting you through separate channels.
// Sometimes, the same issue might be reported by different people in the team or someone might accidentally open a new ticket instead of following up on an existing one.
// To avoid conflicts, you can merge all related tickets together and keep the communication streamlined.
func (fd *FreshDesk) MergeTickets(tm *TicketsMerge) error {
	url := fmt.Sprintf("%s/api/v2/tickets/merge", fd.Domain)
	err := fd.doPut(url, tm, nil)
	return err
}

func (fd *FreshDesk) ListTicketWatchers(tid int64) ([]int64, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/watchers", fd.Domain, tid)
	result := &TicketWatchers{}
	err := fd.doGet(url, result)
	return result.WatcherIDs, err
}

func (fd *FreshDesk) AddTicketWatcher(tid, uid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/watchers", fd.Domain, tid)
	data := map[string]any{
		"user_id": uid,
	}
	return fd.doPost(url, data, nil)
}

func (fd *FreshDesk) UnwatchTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/unwatch", fd.Domain, tid)
	return fd.doPut(url, nil, nil)
}

func (fd *FreshDesk) BulkWatchTickets(tids []int64, uid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/buld_watch", fd.Domain)
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.doPut(url, data, nil)
}

func (fd *FreshDesk) BulkUnwatchTickets(tids []int64, uid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/buld_unwatch", fd.Domain)
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.doPut(url, data, nil)
}

func (fd *FreshDesk) RestoreTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/restore", fd.Domain, tid)
	return fd.doPut(url, nil, nil)
}

func (fd *FreshDesk) DeleteTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fd.Domain, tid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) BulkDeleteTickets(tids []int64) (string, error) {
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

func (fd *FreshDesk) DeleteAttachment(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/attachments/%d", fd.Domain, aid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) ListTicketConversations(tid int64, lco *ListConversationsOption) ([]*Conversation, bool, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/conversations", fd.Domain, tid)
	conversations := []*Conversation{}
	next, err := fd.doList(url, lco, &conversations)
	return conversations, next, err
}

func (fd *FreshDesk) IterTicketConversations(tid int64, lco *ListConversationsOption, icf func(*Conversation) bool) error {
	if lco == nil {
		lco = &ListConversationsOption{}
	}
	if lco.Page < 1 {
		lco.Page = 1
	}
	if lco.PerPage < 1 {
		lco.PerPage = 100
	}

	for {
		conversations, next, err := fd.ListTicketConversations(tid, lco)
		if err != nil {
			return err
		}
		for _, c := range conversations {
			if !icf(c) {
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

func (fd *FreshDesk) CreateReply(tid int64, reply *Reply) (*Reply, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/reply", fd.Domain, tid)
	result := &Reply{}
	err := fd.doPost(url, reply, result)
	return result, err
}

func (fd *FreshDesk) CreateNote(tid int64, note *Note) (*Note, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/notes", fd.Domain, tid)
	result := &Note{}
	err := fd.doPost(url, note, result)
	return result, err
}

// UpdateConversation only public & private notes can be edited.
func (fd *FreshDesk) UpdateConversation(cid int64, conversation *Conversation) (*Conversation, error) {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)
	result := &Conversation{}
	err := fd.doPut(url, conversation, result)
	return result, err
}

// DeleteConversation delete a conversation (Incoming Reply can not be deleted)
func (fd *FreshDesk) DeleteConversation(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fd.Domain, cid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) ReplyToForward(tid int64, rf *ReplyForward) (*ForwardResult, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/reply_to_forward", fd.Domain, tid)
	result := &ForwardResult{}
	err := fd.doPost(url, rf, result)
	return result, err
}

func (fd *FreshDesk) GetAgent(aid int64) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fd.Domain, aid)
	agent := &Agent{}
	err := fd.doGet(url, agent)
	return agent, err
}

func (fd *FreshDesk) ListAgents(lao *ListAgentsOption) ([]*Agent, bool, error) {
	url := fmt.Sprintf("%s/api/v2/agents", fd.Domain)
	agents := []*Agent{}
	next, err := fd.doList(url, lao, &agents)
	return agents, next, err
}

func (fd *FreshDesk) IterAgents(lao *ListAgentsOption, iaf func(*Agent) bool) error {
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

func (fd *FreshDesk) CreateAgent(agent *AgentRequest) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents", fd.Domain)
	result := &Agent{}
	err := fd.doPost(url, agent, result)
	return result, err
}

func (fd *FreshDesk) UpdateAgent(aid int64, agent *AgentRequest) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fd.Domain, aid)
	result := &Agent{}
	err := fd.doPut(url, agent, result)
	return result, err
}

func (fd *FreshDesk) DeleteAgent(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fd.Domain, aid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) GetCurrentAgent() (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/me", fd.Domain)
	agent := &Agent{}
	err := fd.doGet(url, agent)
	return agent, err
}

func (fd *FreshDesk) SearchAgents(keyword string) ([]*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/autocomplete?term=%s", fd.Domain, url.QueryEscape(keyword))
	agents := []*Agent{}
	err := fd.doGet(url, &agents)
	return agents, err
}

func (fd *FreshDesk) CreateContact(contact *Contact) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts", fd.Domain)
	result := &Contact{}
	err := fd.doPost(url, contact, result)
	return result, err
}

func (fd *FreshDesk) UpdateContact(cid int64, contact *Contact) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/%d", fd.Domain, cid)
	result := &Contact{}
	err := fd.doPut(url, contact, result)
	return result, err
}

func (fd *FreshDesk) GetContact(cid int64) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/%d", fd.Domain, cid)
	contact := &Contact{}
	err := fd.doGet(url, contact)
	return contact, err
}

func (fd *FreshDesk) DeleteContact(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d", fd.Domain, cid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) HardDeleteContact(cid int64, force ...bool) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/hard_delete", fd.Domain, cid)
	if len(force) > 0 && force[0] {
		url += "?force=true"
	}
	return fd.doDelete(url)
}

func (fd *FreshDesk) ListContacts(lco *ListContactsOption) ([]*Contact, bool, error) {
	url := fmt.Sprintf("%s/api/v2/contacts", fd.Domain)
	contacts := []*Contact{}
	next, err := fd.doList(url, lco, &contacts)
	return contacts, next, err
}

func (fd *FreshDesk) IterContacts(lco *ListContactsOption, itf func(*Contact) bool) error {
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

func (fd *FreshDesk) SearchContacts(keyword string) ([]*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/autocomplete?term=%s", fd.Domain, url.QueryEscape(keyword))
	contacts := []*Contact{}
	err := fd.doGet(url, &contacts)
	return contacts, err
}

func (fd *FreshDesk) RestoreContact(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/restore", fd.Domain, cid)
	return fd.doPut(url, nil, nil)
}

func (fd *FreshDesk) InviteContact(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/send_invite", fd.Domain, cid)
	return fd.doPut(url, nil, nil)
}

func (fd *FreshDesk) MergeContacts(cm *ContactsMerge) error {
	url := fmt.Sprintf("%s/api/v2/contacts/merge", fd.Domain)
	return fd.doPost(url, nil, nil)
}

// ExportContacts return a job id, call GetExportedContactsURL() to get the job detail
func (fd *FreshDesk) ExportContacts(defaultFields, customFields []string) (string, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/export", fd.Domain)
	data := map[string]any{
		"fields": &ContactsExport{defaultFields, customFields},
	}
	result := map[string]string{}
	err := fd.doPost(url, data, &result)
	return result["id"], err
}

// GetExportedContactsURL get the exported contacts url
func (fd *FreshDesk) GetExportedContactsURL(jid string) (*Job, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/export/%s", fd.Domain, jid)
	job := &Job{}
	err := fd.doGet(url, job)
	return job, err
}

func (fd *FreshDesk) MakeAgent(cid int64, agent *Agent) (*Contact, error) {
	url := fmt.Sprintf("%s/api/v2/contacts/%d/make_agent", fd.Domain, cid)
	result := &Contact{}
	err := fd.doPut(url, agent, result)
	return result, err
}

func (fd *FreshDesk) CreateCategory(category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fd.Domain)
	result := &Category{}
	err := fd.doPost(url, category, result)
	return result, err
}

func (fd *FreshDesk) CreateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	result := &Category{}
	err := fd.doPost(url, category, result)
	return result, err
}

func (fd *FreshDesk) UpdateCategory(cid int64, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fd.Domain, cid)
	result := &Category{}
	err := fd.doPut(url, category, result)
	return result, err
}

func (fd *FreshDesk) UpdateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	result := &Category{}
	err := fd.doPut(url, category, result)
	return result, err
}

func (fd *FreshDesk) GetCategory(cid int64) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fd.Domain, cid)
	cat := &Category{}
	err := fd.doGet(url, cat)
	return cat, err
}

func (fd *FreshDesk) GetCategoryTranslated(cid int64, lang string) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/%s", fd.Domain, cid, lang)
	cat := &Category{}
	err := fd.doGet(url, cat)
	return cat, err
}

func (fd *FreshDesk) ListCategories() ([]*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fd.Domain)
	categories := []*Category{}
	err := fd.doGet(url, &categories)
	return categories, err
}

func (fd *FreshDesk) ListCategoriesTranslated(lang string) ([]*Category, error) {
	url := fd.Domain + "/api/v2/solutions/categories/" + lang
	categories := []*Category{}
	err := fd.doGet(url, &categories)
	return categories, err
}

func (fd *FreshDesk) DeleteCategory(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fd.Domain, cid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) CreateFolder(cid int64, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/folders", fd.Domain, cid)
	result := &Folder{}
	err := fd.doPost(url, folder, result)
	return result, err
}

func (fd *FreshDesk) CreateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Folder{}
	err := fd.doPost(url, folder, result)
	return result, err
}

func (fd *FreshDesk) UpdateFolder(fid int64, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	result := &Folder{}
	err := fd.doPut(url, folder, result)
	return result, err
}

func (fd *FreshDesk) UpdateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	result := &Folder{}
	err := fd.doPut(url, folder, result)
	return result, err
}

func (fd *FreshDesk) GetFolder(fid int64) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	folder := &Folder{}
	err := fd.doGet(url, folder)
	return folder, err
}

func (fd *FreshDesk) GetFolderTranslated(fid int64, lang string) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/%s", fd.Domain, fid, lang)
	folder := &Folder{}
	err := fd.doGet(url, folder)
	return folder, err
}

func (fd *FreshDesk) ListCategoryFolders(cid int64) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/folders", fd.Domain, cid)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *FreshDesk) ListCategoryFoldersTranslated(cid int64, lang string) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d/folders/%s", fd.Domain, cid, lang)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *FreshDesk) ListSubFolders(fid int64) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/subfolders", fd.Domain, fid)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *FreshDesk) ListSubFoldersTranslated(fid int64, lang string) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/subfolders/%s", fd.Domain, fid, lang)
	folders := []*Folder{}
	err := fd.doGet(url, &folders)
	return folders, err
}

func (fd *FreshDesk) DeleteFolder(fid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fd.Domain, fid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) CreateArticle(fid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/articles", fd.Domain, fid)
	result := &Article{}
	err := fd.doPost(url, article, result)
	return result, err
}

func (fd *FreshDesk) CreateArticleTranslated(aid int64, lang string, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/%s", fd.Domain, aid, lang)
	result := &Article{}
	err := fd.doPost(url, article, result)
	return result, err
}

func (fd *FreshDesk) UpdateArticle(aid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fd.Domain, aid)
	result := &Article{}
	err := fd.doPut(url, article, result)
	return result, err
}

func (fd *FreshDesk) UpdateArticleTranslated(aid int64, lang string, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/%s", fd.Domain, aid, lang)
	result := &Article{}
	err := fd.doPut(url, article, result)
	return result, err
}

func (fd *FreshDesk) GetArticle(aid int64) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fd.Domain, aid)
	article := &Article{}
	err := fd.doGet(url, article)
	return article, err
}

func (fd *FreshDesk) GetArticleTranslated(aid int64, lang string) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/%s", fd.Domain, aid, lang)
	article := &Article{}
	err := fd.doGet(url, article)
	return article, err
}

func (fd *FreshDesk) ListFolderArticles(fid int64) ([]*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/articles", fd.Domain, fid)
	articles := []*Article{}
	err := fd.doGet(url, &articles)
	return articles, err
}

func (fd *FreshDesk) ListFolderArticlesTranslated(fid int64, lang string) ([]*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d/farticles/%s", fd.Domain, fid, lang)
	articles := []*Article{}
	err := fd.doGet(url, &articles)
	return articles, err
}

func (fd *FreshDesk) DeleteArticle(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fd.Domain, aid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) SearchArticles(keyword string) ([]*ArticleEx, error) {
	url := fmt.Sprintf("%s/api/v2/search/solutions?term=%s", fd.Domain, url.QueryEscape(keyword))
	articles := []*ArticleEx{}
	err := fd.doGet(url, &articles)
	return articles, err
}

func (fd *FreshDesk) GetRole(rid int64) (*Role, error) {
	url := fmt.Sprintf("%s/api/v2/roles/%d", fd.Domain, rid)
	role := &Role{}
	err := fd.doGet(url, role)
	return role, err
}

func (fd *FreshDesk) ListRoles() ([]*Role, error) {
	url := fmt.Sprintf("%s/api/v2/roles", fd.Domain)
	roles := []*Role{}
	_, err := fd.doList(url, nil, &roles)
	return roles, err
}

func (fd *FreshDesk) GetGroup(gid int64) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	group := &Group{}
	err := fd.doGet(url, group)
	return group, err
}

func (fd *FreshDesk) CreateGroup(group *Group) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fd.Domain)
	result := &Group{}
	err := fd.doPost(url, group, result)
	return result, err
}

func (fd *FreshDesk) ListGroups() ([]*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fd.Domain)
	groups := []*Group{}
	_, err := fd.doList(url, nil, &groups)
	return groups, err
}

func (fd *FreshDesk) UpdateGroup(gid int64, group *Group) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	result := &Group{}
	err := fd.doPut(url, group, result)
	return result, err
}

func (fd *FreshDesk) DeleteGroup(gid int64) error {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) ListAutomationRules(automationTypeID int) ([]*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules", fd.Domain, automationTypeID)
	rules := []*AutomationRule{}
	_, err := fd.doList(url, nil, &rules)
	return rules, err
}

func (fd *FreshDesk) GetAutomationRule(automationTypeID int, rid int64) (*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules/%d", fd.Domain, automationTypeID, rid)
	rule := &AutomationRule{}
	err := fd.doGet(url, rule)
	return rule, err
}

func (fd *FreshDesk) DeleteAutomationRule(automationTypeID int, rid int64) error {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules/%d", fd.Domain, automationTypeID, rid)
	return fd.doDelete(url)
}

func (fd *FreshDesk) CreateAutomationRule(automationTypeID int, rule *AutomationRule) (*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules", fd.Domain, automationTypeID)
	result := &AutomationRule{}
	err := fd.doPost(url, rule, result)
	return result, err
}

func (fd *FreshDesk) UpdateAutomationRule(automationTypeID int, rid int64, rule *AutomationRule) (*AutomationRule, error) {
	url := fmt.Sprintf("%s/api/v2/automations/%d/rules/%d", fd.Domain, automationTypeID, rid)
	result := &AutomationRule{}
	err := fd.doPut(url, rule, result)
	return result, err
}
