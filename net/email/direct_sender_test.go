package email

import (
	"crypto/tls"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/net/netx"
	"github.com/askasoft/pango/str"
)

func testDirectSendEmail(t *testing.T, m *Email) {
	var err error

	s := &DirectSender{}
	s.Timeout = time.Second * 5
	s.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	f := os.Stdout
	w := iox.SyncWriter(f)
	s.ConnDebug = func(conn net.Conn) net.Conn {
		return netx.DumpConn(
			conn,
			iox.WrapWriter(w, iox.ConsoleColor.Magenta+"< ", iox.ConsoleColor.Reset),
			iox.WrapWriter(w, iox.ConsoleColor.Yellow+"> ", iox.ConsoleColor.Reset),
		)
	}

	sd := os.Getenv("SMTP_DIRECT")
	if sd != "true" {
		skipTest(t, "SMTP_DIRECT not set")
		return
	}

	sf := os.Getenv("SMTP_FROM")
	if sf == "" {
		skipTest(t, "SMTP_FROM not set")
		return
	}

	sts := str.RemoveEmpties(str.TrimSpaces(str.Split(os.Getenv("SMTP_TO"), ";")))
	if len(sts) < 1 {
		skipTest(t, "SMTP_TO not set")
		return
	}

	err = m.SetFrom(sf)
	if err != nil {
		t.Error(sf, err)
		return
	}

	err = m.AddTo(sts...)
	if err != nil {
		t.Error(sts, err)
		return
	}

	m.Subject = "direct send subject " + testSendMailTime() + strings.Repeat(" 一二三四五", 10)
	err = s.DirectSend(m)
	if err != nil {
		t.Error(err)
	}
}

func TestDirectSendTextEmailOnly(t *testing.T) {
	email := &Email{}
	email.Message = ".\nthis is a test email " + testSendMailTime() + " from example.com. 一二三四五"
	testDirectSendEmail(t, email)
}

func TestDirectSendTextEmailAttach(t *testing.T) {
	email := &Email{}

	email.Message = ".\nthis is a test email " + testSendMailTime() + " from example.com. 一二三四五"
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
	email.SetHTMLMsg("<pre><font color=red>.\nthis is a test email " + testSendMailTime() + " from example.com. 一二三四五</font></pre>")

	testDirectSendEmail(t, email)
}

func TestDirectSendHtmlEmailAttach(t *testing.T) {
	email := &Email{}

	email.SetHTMLMsg("<pre><IMG src=\"cid:panda.png\"> <font color=red>.\nthis is a test email " + testSendMailTime() + " from example.com. 一二三四五</font></pre>")
	email.AttachString("test.txt", "abcdefg")
	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}

	testDirectSendEmail(t, email)
}
