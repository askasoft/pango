package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/net/netutil"
	"github.com/askasoft/pango/str"
)

func skipTest(t *testing.T, msg string) {
	fmt.Println(msg)
	t.Skip(msg)
}

func testSendMailTime() string {
	return time.Now().Format("2006-01-02T15:04:05.000000")
}

func testSendEmail(t *testing.T, m *Email) {
	var err error

	ss := &SMTPSender{}
	ss.Timeout = time.Second * 5
	ss.Helo = "localhost"

	// f, ef := os.OpenFile("D:\\sender.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))
	// if ef != nil {
	// 	fmt.Println(ef)
	// 	return
	// }
	// defer f.Close()
	f := os.Stdout
	w := iox.SyncWriter(f)
	ss.ConnDebug = func(conn net.Conn) net.Conn {
		return netutil.DumpConn(
			conn,
			iox.WrapWriter(w, iox.ConsoleColor.Magenta+"< ", iox.ConsoleColor.Reset),
			iox.WrapWriter(w, iox.ConsoleColor.Yellow+"> ", iox.ConsoleColor.Reset),
		)
	}
	// ss.DataDebug = func(w io.Writer) io.Writer {
	// 	return io.MultiWriter(os.Stdout, w)
	// }

	ss.Host = os.Getenv("SMTP_HOST")
	if ss.Host == "" {
		skipTest(t, "SMTP_HOST not set")
		return
	}

	ss.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	ss.Username = os.Getenv("SMTP_USER")
	if ss.Username == "" {
		skipTest(t, "SMTP_USER not set")
		return
	}

	ss.Password = os.Getenv("SMTP_PASS")
	if ss.Password == "" {
		skipTest(t, "SMTP_PASS not set")
		return
	}

	sf := os.Getenv("SMTP_FROM")
	if sf == "" {
		skipTest(t, "SMTP_FROM not set")
		return
	}

	sts := str.RemoveEmptys(str.TrimSpaces(str.Split(os.Getenv("SMTP_TO"), ";")))
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

	fmt.Printf("SMTP send %s -> %s\n", m.from, m.GetTos()[0])
	ss.TLSConfig = &tls.Config{ServerName: ss.Host, InsecureSkipVerify: true}
	err = ss.DialAndSend(m)
	if err != nil {
		t.Error("DialAndSend", err)
		return
	}
}

func TestSendTextEmailEmpty(t *testing.T) {
	email := &Email{}
	email.Subject = "TestSendTextEmailEmpty " + testSendMailTime()
	testSendEmail(t, email)
}

func TestSendTextEmailAscii(t *testing.T) {
	email := &Email{}
	email.Subject = "TestSendTextEmailAscii " + testSendMailTime()
	email.Message = strings.Repeat("this is a test email "+testSendMailTime()+".\n", 10)
	testSendEmail(t, email)
}

func TestSendTextEmailAsciiLong(t *testing.T) {
	email := &Email{}
	email.Subject = "TestSendTextEmailAsciiLong " + testSendMailTime()
	email.Message = strings.Repeat("this is a test email "+testSendMailTime()+".\n", 2)
	email.Message += "\r\n+++++++++\r\n"
	email.Message += strings.Repeat("this is a test email "+testSendMailTime()+".  ", 20) // 50 * 20
	email.Message += "\n----------"
	testSendEmail(t, email)
}

func TestSendTextEmailAttach(t *testing.T) {
	email := &Email{}

	email.Subject = "TestSendTextEmailAttach " + testSendMailTime() + strings.Repeat(" 一二三四五", 10)
	email.Message = ".\nthis is a test email " + testSendMailTime() + " from example.com. 一二三四五"
	email.AttachString("string.txt", strings.Repeat("abcdefg一二三四五\r\n", 10))
	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}

	testSendEmail(t, email)
}

func TestSendHtmlEmailOnly(t *testing.T) {
	email := &Email{}

	email.Subject = "TestSendHtmlEmailOnly " + testSendMailTime() + strings.Repeat(" 一二三四五", 10)
	email.SetHTMLMsg("<pre><font color=red>.\nthis is a test email " + testSendMailTime() + " from example.com. 一二三四五</font></pre>")

	testSendEmail(t, email)
}

func TestSendHtmlEmailAttach(t *testing.T) {
	email := &Email{}

	email.Subject = "TestSendHtmlEmailAttach " + testSendMailTime() + strings.Repeat(" 一二三四五", 10)
	email.SetHTMLMsg("<pre>Image: <img src=\"cid:panda.png\">\r\n<font color=red>This is a test email " + testSendMailTime() + " from example.com. 一二三四五</font></pre>")
	email.AttachString("string.txt", strings.Repeat("abcdefg一二三四五\r\n", 10))

	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}
	testSendEmail(t, email)
}
