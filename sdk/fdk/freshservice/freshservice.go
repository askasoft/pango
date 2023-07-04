package freshservice

import (
	"fmt"

	"github.com/askasoft/pango/sdk/fdk"
)

type RateLimitedError = fdk.RateLimitedError

type Date = fdk.Date

type Time = fdk.Time

type Attachment = fdk.Attachment

type Attachments = fdk.Attachments

type ListOption = fdk.ListOption

type PageOption = fdk.PageOption

type File = fdk.File

type Files = fdk.Files

type WithFiles = fdk.WithFiles

type Values = fdk.Values

type OrderType string
type ApprovalStatus int

const (
	OrderAsc  OrderType = "asc"
	OrderDesc OrderType = "desc"

	ApprovalStatusApproved ApprovalStatus = 1
	ApprovalStatusRejected ApprovalStatus = 2
)

func NewAttachment(file string, data ...[]byte) *Attachment {
	return fdk.NewAttachment(file, data...)
}

func toString(o any) string {
	return fdk.ToString(o)
}

type Freshservice fdk.FDK

func (fs *Freshservice) endpoint(format string, a ...any) string {
	return (*fdk.FDK)(fs).Endpoint(format, a...)
}

func (fs *Freshservice) doGet(url string, result any) error {
	return (*fdk.FDK)(fs).DoGet(url, result)
}

func (fs *Freshservice) doList(url string, lo ListOption, ap any) (bool, error) {
	return (*fdk.FDK)(fs).DoList(url, lo, ap)
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

func (fs *Freshservice) SaveFile(url string, filename string) error {
	return (*fdk.FDK)(fs).DoSaveFile(url, filename)
}

func (fs *Freshservice) DownloadNoAuth(url string) ([]byte, error) {
	return (*fdk.FDK)(fs).DoDownloadNoAuth(url)
}

func (fs *Freshservice) SaveFileNoAuth(url string, filename string) error {
	return (*fdk.FDK)(fs).DoSaveFileNoAuth(url, filename)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fs *Freshservice) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("https://%s/helpdesk/attachments/%d", fs.Domain, aid)
}

func (fs *Freshservice) DownloadAttachment(aid int64) ([]byte, error) {
	url := fs.endpoint("/attachments/%d", aid)
	return fs.Download(url)
}

func (fs *Freshservice) SaveAttachment(aid int64, filename string) error {
	url := fs.endpoint("/attachments/%d", aid)
	return fs.SaveFile(url, filename)
}
