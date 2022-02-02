package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test syncwriter
func TestSyncWriteFile(t *testing.T) {
	path := "TestSyncWrite/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(NewSyncWriter(&FileWriter{Path: path}))

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
	bs, _ := ioutil.ReadFile(path)
	as := strings.Split(strings.TrimSuffix(string(bs), EOL), EOL)
	sort.Strings(as)

	// expected data
	es := []string{}
	for n := 1; n < 10; n++ {
		for i := 1; i < 10; i++ {
			es = append(es, fmt.Sprint("[I] ", n, i))
		}
	}

	if !reflect.DeepEqual(as, es) {
		t.Errorf("TestSyncWriteFile\n expect: %q, actual %q", es, as)
	}
}
