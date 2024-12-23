package log

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAsyncWriteConsole(t *testing.T) {
	log := NewLog()
	sw := &StreamWriter{Color: true}
	sw.Formatter = NewTextFormatter("%t{2006-01-02T15:04:05.000} [%c] %l - %m%n")
	aw := NewAsyncWriter(sw, 100)
	log.SetWriter(aw)

	wg := sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		l := log.GetLogger(strconv.Itoa(i))
		wg.Add(1)
		go func() {
			testConsoleCalls(l, 10)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Close()
}

func TestAsyncWriter(t *testing.T) {
	sb := &strings.Builder{}

	sw := &StreamWriter{Output: sb}
	sw.SetFormat("[%p] %m%n")

	lg := NewLog()
	lg.SetWriter(NewAsyncWriter(sw, 10))

	wg := sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			time.Sleep(time.Microsecond * 10)
			for i := 1; i < 10; i++ {
				lg.Info(n, i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	lg.Close()

	// read actual log
	as := strings.Split(strings.TrimSuffix(sb.String(), EOL), EOL)
	sort.Strings(as)

	// expected data
	es := []string{}
	for n := 1; n < 10; n++ {
		for i := 1; i < 10; i++ {
			es = append(es, fmt.Sprint("[I] ", n, i))
		}
	}

	if !reflect.DeepEqual(as, es) {
		t.Errorf("TestAsyncWriter\n expect: %q, actual %q", es, as)
	}
}
