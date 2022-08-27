package log

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

// ConnWriter implements Writer.
// it writes messages in keep-live tcp connection.
type ConnWriter struct {
	Net     string
	Addr    string
	Timeout time.Duration
	Logfmt  Formatter // log formatter
	Logfil  Filter    // log filter

	conn io.WriteCloser
	bb   bytes.Buffer
}

// SetFormat set the log formatter
func (cw *ConnWriter) SetFormat(format string) {
	cw.Logfmt = NewLogFormatter(format)
}

// SetFilter set the log filter
func (cw *ConnWriter) SetFilter(filter string) {
	cw.Logfil = NewLogFilter(filter)
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
	if cw.Logfil != nil && cw.Logfil.Reject(le) {
		return
	}

	lf := cw.Logfmt
	if lf == nil {
		lf = le.Logger().GetFormatter()
		if lf == nil {
			lf = TextFmtDefault
		}
	}

	err = cw.dial()
	if err != nil {
		return
	}

	// format msg
	cw.bb.Reset()
	lf.Write(&cw.bb, le)

	// write log
	_, err = cw.conn.Write(cw.bb.Bytes())
	if err != nil {
		// This is probably due to a timeout, so reconnect and try again.
		cw.Close()
		err = cw.dial()
		if err != nil {
			return
		}

		_, err = cw.conn.Write(cw.bb.Bytes())
		if err != nil {
			err = fmt.Errorf("ConnWriter(%q) - Write(%s): %w", cw.Addr, cw.bb.Bytes(), err)
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
			perrorf("ConnWriter(%q) - Close(): %v", cw.Addr, err)
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
		return fmt.Errorf("ConnWriter(%q) - Dial(%q): %w", cw.Addr, cw.Net, err)
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		err = tcpConn.SetKeepAlive(true)
		if err != nil {
			return fmt.Errorf("ConnWriter(%q) - SetKeepAlive(%q): %w", cw.Addr, cw.Net, err)
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
