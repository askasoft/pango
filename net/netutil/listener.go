package netutil

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// ListenerDumper a listener dump utility
type ListenerDumper struct {
	net.Listener
	Path       string // dump path
	RecvPrefix string
	RecvSuffix string
	SendPrefix string
	SendSuffix string

	disabled bool   // disable the dumper
	sequence uint32 // accepted conn sequence
}

// DumpListener wrap a net.conn for dump
func DumpListener(listener net.Listener, path string) *ListenerDumper {
	return &ListenerDumper{
		Listener:   listener,
		Path:       path,
		RecvPrefix: "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\r\n",
		RecvSuffix: "\r\n\r\n",
		SendPrefix: ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\r\n",
		SendSuffix: "\r\n\r\n",
	}
}

// Disable disable the dumper or not
func (ld *ListenerDumper) Disable(disabled bool) {
	ld.disabled = disabled
}

// Accept waits for and returns the next connection to the listener.
func (ld *ListenerDumper) Accept() (net.Conn, error) {
	conn, err := ld.Listener.Accept()
	if err != nil || ld.disabled {
		return conn, err
	}

	ld.sequence++
	os.MkdirAll(ld.Path, os.FileMode(0770))

	fn := filepath.Join(ld.Path, fmt.Sprintf("%08d.log", ld.sequence))
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0660))
	if err != nil {
		// ignore the dump error
		return conn, nil
	}

	dcon := &connDumper1{
		Conn:       conn,
		Writer:     file,
		RecvPrefix: ld.RecvPrefix,
		RecvSuffix: ld.RecvSuffix,
		SendPrefix: ld.SendPrefix,
		SendSuffix: ld.SendSuffix,
	}
	return dcon, nil
}
