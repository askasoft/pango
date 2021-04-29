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
		cd.Recv.Write(b[:n])
	}
	return n, err
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (cd *connDumper) Write(b []byte) (int, error) {
	cd.Send.Write(b)
	return cd.Conn.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (cd *connDumper) Close() error {
	return cd.Conn.Close()
}
