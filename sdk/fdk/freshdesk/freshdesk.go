package freshdesk

import (
	"fmt"

	"github.com/askasoft/pango/sdk/fdk"
)

type Freshdesk fdk.FDK

func (fd *Freshdesk) endpoint(format string, a ...any) string {
	return (*fdk.FDK)(fd).Endpoint(format, a...)
}

func (fd *Freshdesk) doGet(url string, result any) error {
	return (*fdk.FDK)(fd).DoGet(url, result)
}

func (fd *Freshdesk) doList(url string, lo ListOption, ap any) (bool, error) {
	return (*fdk.FDK)(fd).DoList(url, lo, ap)
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

func (fd *Freshdesk) SaveFile(url string, filename string) error {
	return (*fdk.FDK)(fd).DoSave(url, filename)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fd *Freshdesk) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("https://%s/helpdesk/attachments/%d", fd.Domain, aid)
}

// GetJob get job detail
func (fd *Freshdesk) GetJob(jid string) (*Job, error) {
	url := fd.endpoint("/jobs/%s", jid)
	job := &Job{}
	err := fd.doGet(url, job)
	return job, err
}
