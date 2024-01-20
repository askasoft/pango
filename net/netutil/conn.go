package netutil

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/askasoft/pango/iox"
)

// DumpConn wrap a net.conn for dump
func DumpConn(conn net.Conn, recv io.Writer, send io.Writer) net.Conn {
	return &ConnDumper{conn, recv, send}
}

// ConnDumper a connection dump utility
type ConnDumper struct {
	net.Conn
	Recv io.Writer
	Send io.Writer
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (cd *ConnDumper) Read(b []byte) (int, error) {
	n, err := cd.Conn.Read(b)
	if cd.Recv != nil {
		if n > 0 {
			cd.Recv.Write(b[:n]) //nolint: errcheck
		}
		if err != nil {
			iox.WriteString(cd.Recv, err.Error()) //nolint: errcheck
		}
	}
	return n, err
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (cd *ConnDumper) Write(b []byte) (int, error) {
	n, err := cd.Conn.Write(b)
	if cd.Send != nil {
		if n > 0 {
			cd.Send.Write(b[:n]) //nolint: errcheck
		}
		if err != nil {
			iox.WriteString(cd.Send, err.Error()) //nolint: errcheck
		}
	}
	return n, err
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (cd *ConnDumper) Close() (err error) {
	err = cd.Conn.Close()
	if cd.Recv != nil {
		if c, ok := cd.Recv.(io.Closer); ok {
			c.Close()
		}
	}
	if cd.Send != nil {
		if c, ok := cd.Send.(io.Closer); ok {
			c.Close()
		}
	}
	return
}

// -----------------------------------------------
const (
	debugTimeFormat = "2006-01-02T15:04:05.000"
	stateRecv       = -1
	stateSend       = 1
)

// ConnDebugger a connection debug utility
type ConnDebugger struct {
	net.Conn

	Writer     io.Writer
	RecvPrefix string
	RecvSuffix string
	SendPrefix string
	SendSuffix string
	Timestamp  bool

	state int
}

func (cd *ConnDebugger) timestamp(s string) string {
	if cd.Timestamp {
		s = fmt.Sprintf(s, time.Now().Format(debugTimeFormat))
	}
	return s
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (cd *ConnDebugger) Read(b []byte) (int, error) {
	n, err := cd.Conn.Read(b)

	if cd.state == stateSend {
		iox.WriteString(cd.Writer, cd.timestamp(cd.SendSuffix)) //nolint: errcheck
	}
	if cd.state != stateRecv {
		iox.WriteString(cd.Writer, cd.timestamp(cd.RecvPrefix)) //nolint: errcheck
	}

	cd.state = stateRecv
	if n > 0 {
		cd.Writer.Write(b[:n]) //nolint: errcheck
	}
	if err != nil {
		iox.WriteString(cd.Writer, err.Error()) //nolint: errcheck
		cd.Writer.Write([]byte{'\r', '\n'})     //nolint: errcheck
	}
	return n, err
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (cd *ConnDebugger) Write(b []byte) (int, error) {
	n, err := cd.Conn.Write(b)

	if cd.state == stateRecv {
		iox.WriteString(cd.Writer, cd.timestamp(cd.RecvSuffix)) //nolint: errcheck
	}
	if cd.state != stateSend {
		iox.WriteString(cd.Writer, cd.timestamp(cd.SendPrefix)) //nolint: errcheck
	}

	cd.state = stateSend
	if n > 0 {
		cd.Writer.Write(b[:n]) //nolint: errcheck
	}
	if err != nil {
		iox.WriteString(cd.Writer, err.Error()) //nolint: errcheck
		cd.Writer.Write([]byte{'\r', '\n'})     //nolint: errcheck
	}

	return n, err
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (cd *ConnDebugger) Close() error {
	switch cd.state {
	case stateRecv:
		iox.WriteString(cd.Writer, cd.timestamp(cd.RecvSuffix)) //nolint: errcheck
	case stateSend:
		iox.WriteString(cd.Writer, cd.timestamp(cd.SendSuffix)) //nolint: errcheck
	}

	if c, ok := cd.Writer.(io.Closer); ok {
		c.Close()
	}
	return cd.Conn.Close()
}
