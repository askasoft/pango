package log

import (
	"fmt"
	"io"
	"net"
	"time"
)

// ConnWriter implements Writer.
// it writes messages in keep-live tcp connection.
type ConnWriter struct {
	LogFilter
	LogFormatter

	Net     string
	Addr    string
	Timeout time.Duration

	conn io.WriteCloser
}

// SetTimeout set timeout
func (cw *ConnWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("ConnkWriter - Invalid timeout: %w", err)
	}
	cw.Timeout = tmo
	return nil
}

// Write write logger message to connection.
func (cw *ConnWriter) Write(le *Event) (err error) {
	if cw.Reject(le) {
		return
	}

	err = cw.dial()
	if err != nil {
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

		_, err = cw.conn.Write(cw.bb.Bytes())
		if err != nil {
			err = fmt.Errorf("ConnWriter(%s:%s) - Write([%d]): %w", cw.Net, cw.Addr, len(cw.bb.Bytes()), err)
			cw.Close()
		}
	}
	return
}

// Flush implementing method. empty.
func (cw *ConnWriter) Flush() {
}

// Close close the file description, close file writer.
func (cw *ConnWriter) Close() {
	if cw.conn != nil {
		err := cw.conn.Close()
		if err != nil {
			perrorf("ConnWriter(%s:%s) - Close(): %v", cw.Net, cw.Addr, err)
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
		return fmt.Errorf("ConnWriter(%s:%s) - Dial(): %w", cw.Net, cw.Addr, err)
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err = tcpConn.SetKeepAlive(true)
		if err != nil {
			return fmt.Errorf("ConnWriter(%s:%s) - SetKeepAlive(): %w", cw.Net, cw.Addr, err)
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
