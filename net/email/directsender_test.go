package email

import (
	"os"
	"testing"
	"time"
)

func testDirectSendEmail(t *testing.T, m *Email) {
	st := os.Getenv("SMTP_TO")
	if len(st) < 1 {
		t.Skip("SMTP_TO not set")
		return
	}

	m.SetFrom("testテスター <test@example.com>")
	m.AddTo(st)
	m.Subject = "test subject あいうえお " + time.Now().String()

	s := DirectSender{Timeout: time.Second * 5}
	err := s.DirectSend(m)
	if err != nil {
		t.Error(err)
	}
}

func TestDirectSendTextEmailOnly(t *testing.T) {
	email := &Email{}
	email.Message = ".\nthis is a test email " + time.Now().String() + " from example.com. あいうえお"
	testDirectSendEmail(t, email)
}

func TestDirectSendTextEmailAttach(t *testing.T) {
	email := &Email{}

	email.Message = ".\nthis is a test email " + time.Now().String() + " from example.com. あいうえお"
	email.AttachString("string.txt", "abcdefg")
	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}

	testDirectSendEmail(t, email)
}

func TestDirectSendHtmlEmailOnly(t *testing.T) {
	email := &Email{}
	email.SetHTMLMsg("<pre><font color=red>.\nthis is a test email " + time.Now().String() + " from example.com. あいうえお</font></pre>")

	testDirectSendEmail(t, email)
}

func TestDirectSendHtmlEmailAttach(t *testing.T) {
	email := &Email{}

	email.SetHTMLMsg("<pre><IMG src=\"cid:panda.png\"> <font color=red>.\nthis is a test email " + time.Now().String() + " from example.com. あいうえお</font></pre>")
	email.AttachString("test.txt", "abcdefg")
	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}

	testDirectSendEmail(t, email)
}
