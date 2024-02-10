// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/net/netutil"
	"github.com/askasoft/pango/str"
)

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

	log := NewLog()
	sw := &SMTPWriter{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
		From:     sf,
		Tos:      sts,
	}
	sw.SetSubject("%t [%l] %m")
	log.SetWriter(sw)

	sw.initSender()
	f := os.Stdout
	w := iox.SyncWriter(f)
	sw.sender.ConnDebug = func(conn net.Conn) net.Conn {
		return netutil.DumpConn(
			conn,
			iox.WrapWriter(w, iox.ConsoleColor.Magenta+"< ", iox.ConsoleColor.Reset),
			iox.WrapWriter(w, iox.ConsoleColor.Yellow+"> ", iox.ConsoleColor.Reset),
		)
	}

	log.Fatal("smtp log fatal test!")
}
