package freshservice

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

type Freshservice fdk.FDK

func (fs *Freshservice) Endpoint(format string, a ...any) string {
	return (*fdk.FDK)(fs).Endpoint(format, a...)
}

func (fs *Freshservice) DoGet(ctx context.Context, url string, result any) error {
	return (*fdk.FDK)(fs).DoGet(ctx, url, result)
}

func (fs *Freshservice) DoList(ctx context.Context, url string, lo ListOption, result any) (bool, error) {
	return (*fdk.FDK)(fs).DoList(ctx, url, lo, result)
}

func (fs *Freshservice) DoPost(ctx context.Context, url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPost(ctx, url, source, result)
}

func (fs *Freshservice) DoPut(ctx context.Context, url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPut(ctx, url, source, result)
}

func (fs *Freshservice) DoDelete(ctx context.Context, url string) error {
	return (*fdk.FDK)(fs).DoDelete(ctx, url)
}

func (fs *Freshservice) Download(ctx context.Context, url string) ([]byte, error) {
	return (*fdk.FDK)(fs).DoDownload(ctx, url)
}

func (fs *Freshservice) SaveFile(ctx context.Context, url string, path string) error {
	return (*fdk.FDK)(fs).DoSaveFile(ctx, url, path)
}

func (fs *Freshservice) DownloadNoAuth(ctx context.Context, url string) ([]byte, error) {
	return (*fdk.FDK)(fs).DoDownloadNoAuth(ctx, url)
}

func (fs *Freshservice) SaveFileNoAuth(ctx context.Context, url string, path string) error {
	return (*fdk.FDK)(fs).DoSaveFileNoAuth(ctx, url, path)
}

func (fs *Freshservice) DownloadAttachment(ctx context.Context, aid int64) ([]byte, error) {
	url := fs.Endpoint("/attachments/%d", aid)
	return fs.Download(ctx, url)
}

func (fs *Freshservice) SaveAttachment(ctx context.Context, aid int64, path string) error {
	url := fs.Endpoint("/attachments/%d", aid)
	return fs.SaveFile(ctx, url, path)
}

// GetAgentTicketURL return a permlink for agent ticket URL
func (fs *Freshservice) GetAgentTicketURL(tid int64) string {
	return GetAgentTicketURL(fs.Domain, tid)
}

// GetSolutionArticleURL return a permlink for solution article URL
func (fs *Freshservice) GetSolutionArticleURL(aid int64) string {
	return GetSolutionArticleURL(fs.Domain, aid)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fs *Freshservice) GetHelpdeskAttachmentURL(aid int64) string {
	return GetHelpdeskAttachmentURL(fs.Domain, aid)
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
