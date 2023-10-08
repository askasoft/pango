package freshdesk

import (
	"fmt"

	"github.com/askasoft/pango/sdk/fdk"
)

type RateLimitedError = fdk.RateLimitedError
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

func (fd *Freshdesk) endpoint(format string, a ...any) string {
	return (*fdk.FDK)(fd).Endpoint(format, a...)
}

func (fd *Freshdesk) doGet(url string, result any) error {
	return (*fdk.FDK)(fd).DoGet(url, result)
}

func (fd *Freshdesk) doList(url string, lo ListOption, result any) (bool, error) {
	return (*fdk.FDK)(fd).DoList(url, lo, result)
}

func (fd *Freshdesk) doPost(url string, source, result any) error {
	return (*fdk.FDK)(fd).DoPost(url, source, result)
}

func (fd *Freshdesk) doPut(url string, source, result any) error {
	return (*fdk.FDK)(fd).DoPut(url, source, result)
}

func (fd *Freshdesk) doDelete(url string) error {
	return (*fdk.FDK)(fd).DoDelete(url)
}

func (fd *Freshdesk) Download(url string) ([]byte, error) {
	return (*fdk.FDK)(fd).DoDownload(url)
}

func (fd *Freshdesk) SaveFile(url string, path string) error {
	return (*fdk.FDK)(fd).DoSaveFile(url, path)
}

func (fd *Freshdesk) DownloadNoAuth(url string) ([]byte, error) {
	return (*fdk.FDK)(fd).DoDownloadNoAuth(url)
}

func (fd *Freshdesk) SaveFileNoAuth(url string, path string) error {
	return (*fdk.FDK)(fd).DoSaveFileNoAuth(url, path)
}

// GetSolutionArticleURL return a permlink for solution article URL
func (fd *Freshdesk) GetSolutionArticleURL(aid int64) string {
	return fmt.Sprintf("https://%s/support/solutions/articles/%d", fd.Domain, aid)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fd *Freshdesk) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("https://%s/helpdesk/attachments/%d", fd.Domain, aid)
}

func (fd *Freshdesk) DeleteAttachment(aid int64) error {
	url := fd.endpoint("/attachments/%d", aid)
	return fd.doDelete(url)
}

// GetJob get job detail
func (fd *Freshdesk) GetJob(jid string) (*Job, error) {
	url := fd.endpoint("/jobs/%s", jid)
	job := &Job{}
	err := fd.doGet(url, job)
	return job, err
}
