package xin

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/askasoft/pango/ran"
	"github.com/stretchr/testify/assert"
)

// testRunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (i.e. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func testRunUnix(engine *Engine, file string) (err error) {
	engine.Logger.Infof("Listening and serving HTTP on unix:/%s", file)

	var listener net.Listener
	listener, err = net.Listen("unix", file)
	if err != nil {
		engine.Logger.Errorf("Listen on unix:/%s failed: %v", file, err)
		return
	}
	defer listener.Close()
	defer os.Remove(file)

	server := testNewHttpServer(engine)
	err = server.Serve(listener)
	if err != nil {
		engine.Logger.Errorf("Serve on unix:/%s failed: %v", file, err)
	}
	return
}

func TestUnixSocket(t *testing.T) {
	router := New()

	unixTestSocket := filepath.Join(os.TempDir(), "xin_unix_test_"+ran.RandNumbers(8))

	defer os.Remove(unixTestSocket)

	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, testRunUnix(router, unixTestSocket))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	c, err := net.Dial("unix", unixTestSocket)
	assert.NoError(t, err)

	fmt.Fprint(c, "GET /example HTTP/1.0\r\n\r\n")
	scanner := bufio.NewScanner(c)
	var response string
	for scanner.Scan() {
		response += scanner.Text()
	}
	assert.Contains(t, response, "HTTP/1.0 200", "should get a 200")
	assert.Contains(t, response, "it worked", "resp body should match")
}

func TestBadUnixSocket(t *testing.T) {
	router := New()
	assert.Error(t, testRunUnix(router, "#/tmp/unix_unit_test"))
}

// testRunFd attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified file descriptor.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func testRunFd(engine *Engine, fd int) (err error) {
	engine.Logger.Infof("Listening and serving HTTP on fd@%d", fd)

	var listener net.Listener

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd@%d", fd))
	listener, err = net.FileListener(f)
	if err != nil {
		engine.Logger.Errorf("Listen on fd@%d failed: %v", fd, err)
		return
	}
	defer listener.Close()

	err = testRunListener(engine, listener)
	if err != nil {
		engine.Logger.Errorf("Listen on fd@%d failed: %v", fd, err)
	}
	return
}

func TestFileDescriptor(t *testing.T) {
	router := New()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	assert.NoError(t, err)
	listener, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	socketFile, err := listener.File()
	assert.NoError(t, err)

	go func() {
		router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })
		assert.NoError(t, testRunFd(router, int(socketFile.Fd())))
	}()
	// have to wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(5 * time.Millisecond)

	c, err := net.Dial("tcp", listener.Addr().String())
	assert.NoError(t, err)

	fmt.Fprintf(c, "GET /example HTTP/1.0\r\n\r\n")
	scanner := bufio.NewScanner(c)
	var response string
	for scanner.Scan() {
		response += scanner.Text()
	}
	assert.Contains(t, response, "HTTP/1.0 200", "should get a 200")
	assert.Contains(t, response, "it worked", "resp body should match")
}

func TestBadFileDescriptor(t *testing.T) {
	router := New()
	assert.Error(t, testRunFd(router, 0))
}
