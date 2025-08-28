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

// DumpListener a payload dump listener
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

// NewDumpListener wrap a net.conn for dump payload
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

// LimitListener create a listener that accepts at most n simultaneous connections.
func LimitListener(l net.Listener, n int) net.Listener {
	return NewLimitedListener(l, n)
}

// LimitedListener a Listener that accepts at most n simultaneous
// connections from the provided Listener.
type LimitedListener struct {
	net.Listener
	Semaphore chan struct{}
	closeOnce sync.Once     // ensures the close chan is only closed once
	done      chan struct{} // no values sent; closed when Close is called
}

// NewLimitedListener create a Listener that accepts at most cap(sema) simultaneous
// connections from the provided Listener. n must greater or equal 0. n = 0 means unlimit.
func NewLimitedListener(l net.Listener, n int) *LimitedListener {
	if n < 0 {
		panic("netx: LimitedListener(n) must greater or equal 0")
	}
	return &LimitedListener{
		Listener:  l,
		Semaphore: make(chan struct{}, n),
		done:      make(chan struct{}),
	}
}

// acquire acquires the limiting semaphore. Returns true if successfully
// acquired, false if the listener is closed and the semaphore is not
// acquired.
func (ll *LimitedListener) acquire(semaphore chan<- struct{}) bool {
	select {
	case <-ll.done:
		return false
	case semaphore <- struct{}{}:
		return true
	}
}

func (ll *LimitedListener) Accept() (net.Conn, error) {
	semaphore := ll.Semaphore
	if cap(semaphore) == 0 {
		return ll.Listener.Accept()
	}

	if !ll.acquire(semaphore) {
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

	conn, err := ll.Listener.Accept()
	if err != nil {
		<-semaphore
		return nil, err
	}
	return &limitedListenerConn{Conn: conn, semaphore: semaphore}, nil
}

func (ll *LimitedListener) Close() (err error) {
	err = ll.Listener.Close()
	ll.closeOnce.Do(func() {
		close(ll.done)
	})
	return
}

type limitedListenerConn struct {
	net.Conn
	closeOnce sync.Once
	semaphore <-chan struct{}
}

func (lc *limitedListenerConn) Close() (err error) {
	err = lc.Conn.Close()
	lc.closeOnce.Do(func() {
		<-lc.semaphore
	})
	return
}
