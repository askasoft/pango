package log

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func testConnTCPServer(sigChan chan string, finChan chan string, revChan chan string) {
	wg := &sync.WaitGroup{}

	ln, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}

	// Listen and accept incoming connections
	tl, _ := ln.(*net.TCPListener)
	for {
		tl.SetDeadline(time.Now().Add(time.Second * 1))
		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				select {
				case <-sigChan:
					wg.Wait()
					ln.Close()
					fmt.Println("Listen Done.")
					finChan <- "done"
					return
				default:
				}
				continue
			}
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		wg.Add(1)
		go testConnEcho(conn, wg, revChan)
	}
}

func testConnEcho(conn net.Conn, wg *sync.WaitGroup, revChan chan string) {
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		s := sc.Text()
		revChan <- s
		os.Stdout.Write([]byte(s))
		os.Stdout.Write([]byte("\r\n"))
		if strings.Contains(s, "Close!") {
			fmt.Println("Server Close Connection!")
			if err := conn.Close(); err != nil {
				fmt.Println(err)
			}
			wg.Done()
			return
		}
	}
}

func TestConnWriter(t *testing.T) {
	sigChan := make(chan string, 1)
	finChan := make(chan string, 1)
	revChan := make(chan string, 100)
	go testConnTCPServer(sigChan, finChan, revChan)

	time.Sleep(time.Second)
	lg := NewLog()

	cw := &ConnWriter{Addr: "localhost:9999"}
	cw.SetFormat("%m%n")
	lg.SetWriter(cw)

	ss := []string{
		"Hello Trace",
		"Hello Debug",
		"Hello Info - Close!",
		"Hello Warn",
		"Hello Error",
		"Hello Fatal - Close!",
	}

	i := 0
	lg.Trace(ss[i])
	i++
	lg.Debug(ss[i])
	i++
	lg.Info(ss[i])
	time.Sleep(time.Millisecond * 500)

	// https://gosamples.dev/broken-pipe/
	lg.Info(strings.Repeat("!missing! ", 1000))
	time.Sleep(time.Millisecond * 500)

	i++
	lg.Warn(ss[i])
	i++
	lg.Error(ss[i])
	i++
	lg.Fatal(ss[i])
	time.Sleep(time.Millisecond * 500)

	sigChan <- "done"
	<-finChan

	rs := []string{}
	for {
		if len(revChan) > 0 {
			s := <-revChan
			rs = append(rs, s)
			continue
		}
		break
	}

	if !reflect.DeepEqual(ss, rs) {
		t.Errorf("TestConnWriter() failure\n expect: %v\n actual: %v", ss, rs)
	}
}
