package smtplog

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/net/netx"
)

func skipTest(t *testing.T, msg string) {
	fmt.Println(msg)
	t.Skip(msg)
}

func TestSmtpWriter(t *testing.T) {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		skipTest(t, "SMTP_HOST not set")
		return
	}

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	user := os.Getenv("SMTP_USER")
	if user == "" {
		skipTest(t, "SMTP_USER not set")
		return
	}

	pass := os.Getenv("SMTP_PASS")
	if pass == "" {
		skipTest(t, "SMTP_PASS not set")
		return
	}

	from := os.Getenv("SMTP_FROM")
	if from == "" {
		skipTest(t, "SMTP_FROM not set")
		return
	}

	to := os.Getenv("SMTP_TO")
	if to == "" {
		skipTest(t, "SMTP_TO not set")
		return
	}

	lg := log.NewLog()
	sw := &SMTPWriter{
		Host:     host,
		Port:     port,
		Insecure: true,
		Username: user,
		Password: pass,
	}
	sw.SetFrom(from)
	sw.SetTo(to)
	sw.SetSubject("%t [%l] %m")
	lg.SetWriter(sw)

	sw.initSender()
	f := os.Stdout
	w := iox.SyncWriter(f)
	sw.sender.ConnDebug = func(conn net.Conn) net.Conn {
		return netx.DumpConn(
			conn,
			iox.WrapWriter(w, iox.ConsoleColor.Magenta+"< ", iox.ConsoleColor.Reset),
			iox.WrapWriter(w, iox.ConsoleColor.Yellow+"> ", iox.ConsoleColor.Reset),
		)
	}

	lg.Fatal("smtp log fatal test!")
}
