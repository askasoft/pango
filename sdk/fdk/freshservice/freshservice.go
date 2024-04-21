package freshservice

import (
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

func (fs *Freshservice) endpoint(format string, a ...any) string {
	return (*fdk.FDK)(fs).Endpoint(format, a...)
}

func (fs *Freshservice) doGet(url string, result any) error {
	return (*fdk.FDK)(fs).DoGet(url, result)
}

func (fs *Freshservice) doList(url string, lo ListOption, result any) (bool, error) {
	return (*fdk.FDK)(fs).DoList(url, lo, result)
}

func (fs *Freshservice) doPost(url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPost(url, source, result)
}

func (fs *Freshservice) doPut(url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPut(url, source, result)
}

func (fs *Freshservice) doDelete(url string) error {
	return (*fdk.FDK)(fs).DoDelete(url)
}

func (fs *Freshservice) Download(url string) ([]byte, error) {
	return (*fdk.FDK)(fs).DoDownload(url)
}

func (fs *Freshservice) SaveFile(url string, path string) error {
	return (*fdk.FDK)(fs).DoSaveFile(url, path)
}

func (fs *Freshservice) DownloadNoAuth(url string) ([]byte, error) {
	return (*fdk.FDK)(fs).DoDownloadNoAuth(url)
}

func (fs *Freshservice) SaveFileNoAuth(url string, path string) error {
	return (*fdk.FDK)(fs).DoSaveFileNoAuth(url, path)
}

// GetSolutionArticleURL return a permlink for solution article URL
func (fs *Freshservice) GetSolutionArticleURL(aid int64) string {
	return fmt.Sprintf("https://%s/support/solutions/articles/%d", fs.Domain, aid)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fs *Freshservice) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("https://%s/helpdesk/attachments/%d", fs.Domain, aid)
}

func (fs *Freshservice) DownloadAttachment(aid int64) ([]byte, error) {
	url := fs.endpoint("/attachments/%d", aid)
	return fs.Download(url)
}

func (fs *Freshservice) SaveAttachment(aid int64, path string) error {
	url := fs.endpoint("/attachments/%d", aid)
	return fs.SaveFile(url, path)
}
