package log

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/askasoft/pango/tmu"
)

// ConnWriter implements Writer.
// it writes messages in keep-live tcp connection.
type ConnWriter struct {
	FilterSupport
	FormatSupport

	Net     string
	Addr    string
	Timeout time.Duration

	conn io.WriteCloser
}

// SetTimeout set timeout
func (cw *ConnWriter) SetTimeout(timeout string) error {
	tmo, err := tmu.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("connlog: invalid timeout: %w", err)
	}
	cw.Timeout = tmo
	return nil
}

// Write write logger message to connection.
func (cw *ConnWriter) Write(le *Event) {
	if cw.Reject(le) {
		return
	}

	if err := cw.write(le); err != nil {
		Perror(err)
	}
}

func (cw *ConnWriter) write(le *Event) (err error) {
	if err = cw.dial(); err != nil {
		return
	}

	// format msg
	bs := cw.Format(le)

	// write log
	_, err = cw.conn.Write(bs)
	if err != nil {
		// This is probably due to a timeout, so reconnect and try again.
		cw.Close()
		err = cw.dial()
		if err != nil {
			return
		}

		_, err = cw.conn.Write(cw.Buffer.Bytes())
		if err != nil {
			err = fmt.Errorf("connlog: (%s:%s) Write([%d]): %w", cw.Net, cw.Addr, cw.Buffer.Len(), err)
			cw.Close()
		}
	}
	return
}

// Flush do nothing.
func (cw *ConnWriter) Flush() {
}

// Close close the connection.
func (cw *ConnWriter) Close() {
	if cw.conn != nil {
		err := cw.conn.Close()
		if err != nil {
			Perrorf("connlog: (%s:%s) Close(): %v", cw.Net, cw.Addr, err)
		}
		cw.conn = nil
	}
}

func (cw *ConnWriter) dial() error {
	if cw.conn != nil {
		return nil
	}

	if cw.Net == "" {
		cw.Net = "tcp"
	}

	if cw.Timeout.Milliseconds() == 0 {
		cw.Timeout = time.Second * 2
	}

	conn, err := net.DialTimeout(cw.Net, cw.Addr, cw.Timeout)
	if err != nil {
		return fmt.Errorf("connlog: Dial(%s:%s): %w", cw.Net, cw.Addr, err)
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err = tcpConn.SetKeepAlive(true)
		if err != nil {
			return fmt.Errorf("connlog: (%s:%s) SetKeepAlive(): %w", cw.Net, cw.Addr, err)
		}
	}

	cw.conn = conn
	return nil
}

func newConnWriter() Writer {
	return &ConnWriter{Net: "tcp", Timeout: time.Second * 2}
}

func init() {
	RegisterWriter("conn", newConnWriter)
	RegisterWriter("tcp", newConnWriter)
}
