package netutil

import (
	"io"
	"net"
)

// connDumper a connection dump utility
type connDumper struct {
	net.Conn
	Recv io.Writer
	Send io.Writer
}

// DumpConn wrap a net.conn for dump
func DumpConn(conn net.Conn, recv io.Writer, send io.Writer) net.Conn {
	return &connDumper{conn, recv, send}
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (cd *connDumper) Read(b []byte) (int, error) {
	n, err := cd.Conn.Read(b)
	if n > 0 {
		cd.Recv.Write(b[:n]) //nolint: errcheck
	}
	if err != nil {
		io.WriteString(cd.Recv, err.Error()) //nolint: errcheck
	}
	return n, err
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (cd *connDumper) Write(b []byte) (int, error) {
	n, err := cd.Conn.Write(b)
	if n > 0 {
		cd.Send.Write(b[:n]) //nolint: errcheck
	}
	if err != nil {
		io.WriteString(cd.Send, err.Error()) //nolint: errcheck
	}
	return n, err
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (cd *connDumper) Close() error {
	if c, ok := cd.Recv.(io.Closer); ok {
		c.Close()
	}
	if c, ok := cd.Send.(io.Closer); ok {
		c.Close()
	}
	return cd.Conn.Close()
}

//-----------------------------------------------
const stateRecv = -1
const stateSend = 1

// connDumper1 a connection dump utility
type connDumper1 struct {
	net.Conn

	Writer     io.Writer
	RecvPrefix string
	RecvSuffix string
	SendPrefix string
	SendSuffix string

	state int
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (cd *connDumper1) Read(b []byte) (int, error) {
	n, err := cd.Conn.Read(b)

	if cd.state == stateSend {
		io.WriteString(cd.Writer, cd.SendSuffix) //nolint: errcheck
	}
	if cd.state != stateRecv {
		io.WriteString(cd.Writer, cd.RecvPrefix) //nolint: errcheck
	}

	cd.state = stateRecv
	if n > 0 {
		cd.Writer.Write(b[:n]) //nolint: errcheck
	}
	if err != nil {
		io.WriteString(cd.Writer, err.Error()) //nolint: errcheck
		cd.Writer.Write([]byte{'\r', '\n'})    //nolint: errcheck
	}
	return n, err
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (cd *connDumper1) Write(b []byte) (int, error) {
	n, err := cd.Conn.Write(b)

	if cd.state == stateRecv {
		io.WriteString(cd.Writer, cd.RecvSuffix) //nolint: errcheck
	}
	if cd.state != stateSend {
		io.WriteString(cd.Writer, cd.SendPrefix) //nolint: errcheck
	}

	cd.state = stateSend
	if n > 0 {
		cd.Writer.Write(b[:n]) //nolint: errcheck
	}
	if err != nil {
		io.WriteString(cd.Writer, err.Error()) //nolint: errcheck
		cd.Writer.Write([]byte{'\r', '\n'})    //nolint: errcheck
	}

	return n, err
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (cd *connDumper1) Close() error {
	switch cd.state {
	case stateRecv:
		io.WriteString(cd.Writer, cd.RecvSuffix) //nolint: errcheck
	case stateSend:
		io.WriteString(cd.Writer, cd.SendSuffix) //nolint: errcheck
	}

	if c, ok := cd.Writer.(io.Closer); ok {
		c.Close()
	}
	return cd.Conn.Close()
}
