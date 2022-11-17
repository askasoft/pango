package freshdesk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"

	"github.com/pandafw/pango/net/httpx"
	"github.com/pandafw/pango/num"
)

type WithAttachments interface {
	GetAttachments() []*Attachment
	Values() Values
}

type Values url.Values

func (vs Values) Map() map[string][]string {
	return (map[string][]string)(vs)
}

func (vs Values) SetBool(name string, value bool) {
	s := "false"
	if value {
		s = "true"
	}
	(url.Values)(vs).Set(name, s)
}

func (vs Values) SetString(name string, value string) {
	if value != "" {
		(url.Values)(vs).Set(name, value)
	}
}

func (vs Values) SetStrings(name string, value []string) {
	if len(value) > 0 {
		for _, s := range value {
			(url.Values)(vs).Add(name, s)
		}
	}
}

func (vs Values) SetTime(name string, value Time) {
	if !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetTimePtr(name string, value *Time) {
	if value != nil && !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetInt(name string, value int) {
	if value != 0 {
		(url.Values)(vs).Set(name, num.Itoa(value))
	}
}

func (vs Values) SetInt64(name string, value int64) {
	if value != 0 {
		(url.Values)(vs).Set(name, num.Ltoa(value))
	}
}

func addMultipartValues(mw *httpx.MultipartWriter, vs Values) error {
	return mw.WriteFields(url.Values(vs))
}

func addMultipartAttachments(mw *httpx.MultipartWriter, as []*Attachment) (err error) {
	for _, a := range as {
		if len(a.data) > 0 {
			err = mw.WriteFileData("attachments[]", a.path, a.data)
		} else {
			err = mw.WriteFile("attachments[]", a.path)
		}
		if err != nil {
			return
		}
	}
	return
}

func buildAttachmentsRequest(wa WithAttachments) (io.Reader, string, error) {
	if len(wa.GetAttachments()) == 0 {
		body, err := json.Marshal(wa)
		if err != nil {
			return nil, "", err
		}

		buf := bytes.NewReader(body)
		return buf, contentTypeJSON, nil
	}

	buf := &bytes.Buffer{}
	mw := httpx.NewMultipartWriter(buf)

	contentType := mw.FormDataContentType()

	if err := addMultipartValues(mw, wa.Values()); err != nil {
		return nil, "", err
	}
	if err := addMultipartAttachments(mw, wa.GetAttachments()); err != nil {
		return nil, "", err
	}
	if err := mw.Close(); err != nil {
		return nil, "", err
	}

	return buf, contentType, nil
}
