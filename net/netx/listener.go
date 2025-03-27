package netx

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DumpListener a listener dump utility
type DumpListener struct {
	net.Listener
	Path       string // dump path
	RecvPrefix string
	RecvSuffix string
	SendPrefix string
	SendSuffix string
	Timestamp  bool

	disabled bool // disable the dumper
}

// NewDumpListener wrap a net.conn for dump
func NewDumpListener(listener net.Listener, path string) *DumpListener {
	return &DumpListener{
		Listener:   listener,
		Path:       path,
		RecvPrefix: ">>>>>>>> %s >>>>>>>>\r\n",
		RecvSuffix: "\r\n%.s\r\n",
		SendPrefix: "<<<<<<<< %s <<<<<<<<\r\n",
		SendSuffix: "\r\n%.s\r\n",
		Timestamp:  true,
	}
}

// Disable disable the dumper or not
func (dl *DumpListener) Disable(disabled bool) {
	dl.disabled = disabled
}

// Accept waits for and returns the next connection to the listener.
func (dl *DumpListener) Accept() (conn net.Conn, err error) {
	conn, err = dl.Listener.Accept()
	if err != nil || dl.disabled {
		return
	}

	conn = dl.dump(conn)
	return
}

func (dl *DumpListener) dump(conn net.Conn) net.Conn {
	err := os.MkdirAll(dl.Path, os.FileMode(0770))
	if err != nil {
		// ignore the dump error
		return conn
	}

	fn := fmt.Sprintf("%s_%s.log", time.Now().Format("20060102150405.999999999"), strings.ReplaceAll(conn.RemoteAddr().String(), ":", "_"))
	fn = filepath.Join(dl.Path, fn)
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0660))
	if err != nil {
		// ignore the dump error
		return conn
	}

	dcon := &ConnDebugger{
		Conn:       conn,
		Writer:     file,
		RecvPrefix: dl.RecvPrefix,
		RecvSuffix: dl.RecvSuffix,
		SendPrefix: dl.SendPrefix,
		SendSuffix: dl.SendSuffix,
		Timestamp:  true,
	}
	return dcon
}
