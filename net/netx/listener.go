package netx

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

// LimitListener a Listener that accepts at most n simultaneous
// connections from the provided Listener.
type LimitListener struct {
	net.Listener
	sema      chan struct{}
	done      chan struct{} // no values sent; closed when Close is called
	closeOnce sync.Once     // ensures the close chan is only closed once
}

// NewLimitListener create a Listener that accepts at most cap(sema) simultaneous
// connections from the provided Listener.
func NewLimitListener(l net.Listener, sema chan struct{}) *LimitListener {
	return &LimitListener{
		Listener: l,
		sema:     sema,
		done:     make(chan struct{}),
	}
}

// acquire acquires the limiting semaphore. Returns true if successfully
// acquired, false if the listener is closed and the semaphore is not
// acquired.
func (ll *LimitListener) acquire() bool {
	select {
	case <-ll.done:
		return false
	case ll.sema <- struct{}{}:
		return true
	}
}

func (ll *LimitListener) Accept() (net.Conn, error) {
	if !ll.acquire() {
		// If the semaphore isn't acquired because the listener was closed, expect
		// that this call to accept won't block, but immediately return an error.
		// If it instead returns a spurious connection (due to a bug in the
		// Listener, such as https://golang.org/issue/50216), we immediately close
		// it and try again. Some buggy Listener implementations (like the one in
		// the aforementioned issue) seem to assume that Accept will be called to
		// completion, and may otherwise fail to clean up the client end of pending
		// connections.
		for {
			c, err := ll.Listener.Accept()
			if err != nil {
				return nil, err
			}
			c.Close()
		}
	}

	c, err := ll.Listener.Accept()
	if err != nil {
		ll.release()
		return nil, err
	}
	return &limitListenerConn{Conn: c, release: ll.release}, nil
}

func (ll *LimitListener) release() {
	<-ll.sema
}

func (ll *LimitListener) Close() error {
	err := ll.Listener.Close()
	ll.closeOnce.Do(func() {
		close(ll.done)
	})
	return err
}

type limitListenerConn struct {
	net.Conn
	releaseOnce sync.Once
	release     func()
}

func (lc *limitListenerConn) Close() error {
	err := lc.Conn.Close()
	lc.releaseOnce.Do(lc.release)
	return err
}
