package freshdesk

import (
	"context"
	"fmt"

	"github.com/askasoft/pango/sdk/fdk"
)

type Date = fdk.Date
type Time = fdk.Time
type TimeSpent = fdk.TimeSpent
type Attachment = fdk.Attachment
type Attachments = fdk.Attachments
type ListOption = fdk.ListOption
type PageOption = fdk.PageOption
type File = fdk.File
type Files = fdk.Files
type WithFiles = fdk.WithFiles
type Values = fdk.Values

type OrderType string

const (
	OrderAsc  OrderType = "asc"
	OrderDesc OrderType = "desc"
)

func ParseDate(s string) (*Date, error) {
	return fdk.ParseDate(s)
}

func ParseTime(s string) (*Time, error) {
	return fdk.ParseTime(s)
}

func ParseTimeSpent(s string) (TimeSpent, error) {
	return fdk.ParseTimeSpent(s)
}

func NewAttachment(file string, data ...[]byte) *Attachment {
	return fdk.NewAttachment(file, data...)
}

func toString(o any) string {
	return fdk.ToJSONIndent(o)
}

func ToJSON(o any) string {
	return fdk.ToJSON(o)
}

func ToJSONIndent(o any) string {
	return fdk.ToJSONIndent(o)
}

type Freshdesk fdk.FDK

func (fd *Freshdesk) Endpoint(format string, a ...any) string {
	return (*fdk.FDK)(fd).Endpoint(format, a...)
}

func (fd *Freshdesk) DoGet(ctx context.Context, url string, result any) error {
	return (*fdk.FDK)(fd).DoGet(ctx, url, result)
}

func (fd *Freshdesk) DoList(ctx context.Context, url string, lo ListOption, result any) (bool, error) {
	return (*fdk.FDK)(fd).DoList(ctx, url, lo, result)
}

func (fd *Freshdesk) DoPost(ctx context.Context, url string, source, result any) error {
	return (*fdk.FDK)(fd).DoPost(ctx, url, source, result)
}

func (fd *Freshdesk) DoPut(ctx context.Context, url string, source, result any) error {
	return (*fdk.FDK)(fd).DoPut(ctx, url, source, result)
}

func (fd *Freshdesk) DoDelete(ctx context.Context, url string) error {
	return (*fdk.FDK)(fd).DoDelete(ctx, url)
}

func (fd *Freshdesk) Download(ctx context.Context, url string) ([]byte, error) {
	return (*fdk.FDK)(fd).DoDownload(ctx, url)
}

func (fd *Freshdesk) SaveFile(ctx context.Context, url string, path string) error {
	return (*fdk.FDK)(fd).DoSaveFile(ctx, url, path)
}

func (fd *Freshdesk) DownloadNoAuth(ctx context.Context, url string) ([]byte, error) {
	return (*fdk.FDK)(fd).DoDownloadNoAuth(ctx, url)
}

func (fd *Freshdesk) SaveFileNoAuth(ctx context.Context, url string, path string) error {
	return (*fdk.FDK)(fd).DoSaveFileNoAuth(ctx, url, path)
}

func (fd *Freshdesk) DeleteAttachment(ctx context.Context, aid int64) error {
	url := fd.Endpoint("/attachments/%d", aid)
	return fd.DoDelete(ctx, url)
}

// GetJob get job detail
func (fd *Freshdesk) GetJob(ctx context.Context, jid string) (*Job, error) {
	url := fd.Endpoint("/jobs/%s", jid)
	job := &Job{}
	err := fd.DoGet(ctx, url, job)
	return job, err
}

// GetAgentTicketURL return a permlink for agent ticket URL
func (fd *Freshdesk) GetAgentTicketURL(tid int64) string {
	return GetAgentTicketURL(fd.Domain, tid)
}

// GetSolutionArticleURL return a permlink for solution article URL
func (fd *Freshdesk) GetSolutionArticleURL(aid int64) string {
	return GetSolutionArticleURL(fd.Domain, aid)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fd *Freshdesk) GetHelpdeskAttachmentURL(aid int64) string {
	return GetHelpdeskAttachmentURL(fd.Domain, aid)
}

// GetAgentTicketURL return a permlink for agent ticket URL
func GetAgentTicketURL(domain string, tid int64) string {
	return fmt.Sprintf("https://%s/a/tickets/%d", domain, tid)
}

// GetSolutionArticleURL return a permlink for solution article URL
func GetSolutionArticleURL(domain string, aid int64) string {
	return fmt.Sprintf("https://%s/support/solutions/articles/%d", domain, aid)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func GetHelpdeskAttachmentURL(domain string, aid int64) string {
	return fmt.Sprintf("https://%s/helpdesk/attachments/%d", domain, aid)
}
