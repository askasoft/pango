package email

import (
	"crypto/tls"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pandafw/pango/iox"
)

func testSendEmail(t *testing.T, m *Email) {
	var err error

	s := &Sender{Timeout: time.Second * 5, Helo: "localhost"}
	// f, err = os.OpenFile("sender.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer f.Close()
	f := os.Stdout
	w := iox.SyncWriter(f)
	s.ConnDebug = func(conn net.Conn) net.Conn {
		return &iox.ConnDump{
			Conn: conn,
			Recv: &iox.WrapWriter{Writer: w, Prefix: iox.ConsoleColor.Magenta + "< ", Suffix: iox.ConsoleColor.Reset},
			Send: &iox.WrapWriter{Writer: w, Prefix: iox.ConsoleColor.Yellow + "> ", Suffix: iox.ConsoleColor.Reset},
		}
	}
	// s.DataDebug = func(w io.Writer) io.Writer {
	// 	return io.MultiWriter(os.Stdout, w)
	// }

	s.Host = os.Getenv("SMTP_HOST")
	if len(s.Host) < 1 {
		t.Skip("SMTP_HOST not set")
		return
	}

	s.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	s.Username = os.Getenv("SMTP_USER")
	if len(s.Username) < 1 {
		t.Skip("SMTP_USER not set")
		return
	}

	s.Password = os.Getenv("SMTP_PASS")
	if len(s.Password) < 1 {
		t.Skip("SMTP_PASS not set")
		return
	}

	sf := os.Getenv("SMTP_FROM")
	if len(sf) < 1 {
		t.Skip("SMTP_FROM not set")
		return
	}

	st := os.Getenv("SMTP_TO")
	if len(st) < 1 {
		t.Skip("SMTP_TO not set")
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
	m.Subject = "test subject " + time.Now().String() + strings.Repeat(" あいうえお", 10)

	s.TLSConfig = &tls.Config{ServerName: s.Host, InsecureSkipVerify: true}
	err = s.DialAndSend(m)
	if err != nil {
		t.Error("DialAndSend", err)
		return
	}
}

func TestSendTextEmailOnly(t *testing.T) {
	email := &Email{}
	email.Message = ".\nthis is a test email " + time.Now().String() + " from example.com. あいうえお"
	testSendEmail(t, email)
}

func TestSendTextEmailAttach(t *testing.T) {
	email := &Email{}

	email.Message = ".\nthis is a test email " + time.Now().String() + " from example.com. あいうえお"
	email.AttachString("string.txt", strings.Repeat("abcdefgあいうえお\r\n", 10))
	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}

	testSendEmail(t, email)
}

func TestSendHtmlEmailOnly(t *testing.T) {
	email := &Email{}
	email.SetHTMLMsg("<pre><font color=red>.\nthis is a test email " + time.Now().String() + " from example.com. あいうえお</font></pre>")

	testSendEmail(t, email)
}

func TestSendHtmlEmailAttach(t *testing.T) {
	email := &Email{}

	email.SetHTMLMsg("<pre><IMG src=\"cid:panda.png\"> <font color=red>.\nthis is a test email " + time.Now().String() + " from example.com. あいうえお</font></pre>")
	email.AttachString("string.txt", strings.Repeat("abcdefgあいうえお\r\n", 10))

	err := email.EmbedFile("panda.png", "testdata/panda.png")
	if err != nil {
		t.Error(err)
		return
	}
	testSendEmail(t, email)
}
