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

	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/net/netutil"
)

func skipTest(t *testing.T, msg string) {
	fmt.Println(msg)
	t.Skip(msg)
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

	st := os.Getenv("SMTP_TO")
	if st == "" {
		skipTest(t, "SMTP_TO not set")
		return
	}

	err = m.SetFrom(sf)
	if err != nil {
		t.Error(sf, err)
		return
	}
	err = m.AddTo(st)
	if err != nil {
		t.Error(st, err)
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

func TestSendTextEmailOnly(t *testing.T) {
	email := &Email{}
	email.Subject = "TestSendTextEmailOnly " + time.Now().String()
	email.Message = "this is a test email " + time.Now().String()
	testSendEmail(t, email)
}

func TestSendTextEmailAttach(t *testing.T) {
	email := &Email{}

	email.Subject = "TestSendTextEmailAttach " + time.Now().String() + strings.Repeat(" 一二三四五", 10)
	email.Message = ".\nthis is a test email " + time.Now().String() + " from example.com. 一二三四五"
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

	email.Subject = "TestSendHtmlEmailOnly " + time.Now().String() + strings.Repeat(" 一二三四五", 10)
	email.SetHTMLMsg("<pre><font color=red>.\nthis is a test email " + time.Now().String() + " from example.com. 一二三四五</font></pre>")

	testSendEmail(t, email)
}

func TestSendHtmlEmailAttach(t *testing.T) {
	email := &Email{}

	email.Subject = "TestSendHtmlEmailAttach " + time.Now().String() + strings.Repeat(" 一二三四五", 10)
	email.SetHTMLMsg("<pre>Image: <img src=\"cid:panda.png\">\r\n<font color=red>This is a test email " + time.Now().String() + " from example.com. 一二三四五</font></pre>")
	email.AttachString("string.txt", strings.Repeat("abcdefg一二三四五\r\n", 10))

	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}
	testSendEmail(t, email)
}
