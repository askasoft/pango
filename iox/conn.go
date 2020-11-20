package iox

import (
	"io"
	"net"
	"time"
)

// ConnWrapFunc a connect wrapper function
type ConnWrapFunc func(conn net.Conn) net.Conn

// ConnDump a connection dump utility
type ConnDump struct {
	Conn net.Conn
	Send io.Writer
	Recv io.Writer
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (dc *ConnDump) Read(b []byte) (int, error) {
	n, err := dc.Conn.Read(b)
	if n > 0 {
		dc.Recv.Write(b[:n])
	}
	return n, err
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (dc *ConnDump) Write(b []byte) (int, error) {
	dc.Send.Write(b)
	return dc.Conn.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (dc *ConnDump) Close() error {
	return dc.Conn.Close()
}

// LocalAddr returns the local network address.
func (dc *ConnDump) LocalAddr() net.Addr {
	return dc.Conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (dc *ConnDump) RemoteAddr() net.Addr {
	return dc.Conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (dc *ConnDump) SetDeadline(t time.Time) error {
	return dc.Conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (dc *ConnDump) SetReadDeadline(t time.Time) error {
	return dc.Conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (dc *ConnDump) SetWriteDeadline(t time.Time) error {
	return dc.Conn.SetWriteDeadline(t)
}
