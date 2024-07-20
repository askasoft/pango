package log

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSyncWriter(t *testing.T) {
	sb := &strings.Builder{}

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(NewSyncWriter(&StreamWriter{Output: sb}))

	// test concurrent write
	wg := sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			time.Sleep(time.Microsecond * 10)
			for i := 1; i < 10; i++ {
				log.Info(n, i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Close()

	// test concurrent write after close
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			time.Sleep(time.Microsecond * 10)
			for i := 1; i < 10; i++ {
				log.Info(n, i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

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
		t.Errorf("TestSyncWriter\n expect: %q, actual %q", es, as)
	}
}
