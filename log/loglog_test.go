package log

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pandafw/pango/str"
)

const testRoutines = 10

type testConcurrentDetectWriter struct {
	count1 uint64
	count2 uint64
	count3 uint64

	counts [testRoutines]int64
	closed bool
	last   time.Time

	error string
}

func newTestConcurrentDetectWriter() *testConcurrentDetectWriter {
	tw := &testConcurrentDetectWriter{}
	tw.last = time.Now()
	return tw
}

func (tw *testConcurrentDetectWriter) do() {
	if tw.closed {
		atomic.AddUint64(&tw.count3, 1)
	} else {
		tw.count1++
		atomic.AddUint64(&tw.count2, 1)
	}
}

func (tw *testConcurrentDetectWriter) Write(le *Event) {
	tw.do()

	ss := str.Split(le.msg, " ")
	k, _ := strconv.Atoi(ss[0])
	c, _ := strconv.ParseInt(ss[1], 10, 64)
	c0 := tw.counts[k]
	if c0 != 0 && c0+1 != c {
		tw.error = fmt.Sprintf("[%d] %d <- %d", k, c0, c)
	}
	tw.counts[k] = c

	t := time.Now()
	if t.After(tw.last.Add(time.Second)) {
		fmt.Println(le.when, k, c, tw.count2)
		tw.last = t
	}
}

func (tw *testConcurrentDetectWriter) Flush() {
	tw.do()
}

func (tw *testConcurrentDetectWriter) Close() {
	tw.do()
	tw.closed = true
}

func testLogRoutine(log *Log, wg *sync.WaitGroup, n int, p *int64) {
	var m int64
	et := time.Now().Add(time.Second * 5)
	for time.Now().Before(et) {
		m++
		log.Infof("%d %d", n, m)
	}
	*p = m
	wg.Done()
	fmt.Println(time.Now(), "test routine ", n, " done ", m)
}

func testCheckConcurrentDetectWriter(t *testing.T, c string, tw1 *testConcurrentDetectWriter, tw2 *testConcurrentDetectWriter, cs []int64) {
	fmt.Println("tw1: ", tw1.closed, tw1.count1, tw1.count2, tw1.count3, tw1.error)
	fmt.Println("tw2: ", tw2.closed, tw2.count1, tw2.count2, tw2.count3, tw2.error)

	if tw1.error != "" {
		t.Errorf("%s(%s) error: %s", c, "tw1.error", tw1.error)
	}
	if !tw1.closed {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw1.closed", true, tw1.closed)
	}
	if tw1.count1 != tw1.count2 {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw1.count1/2", tw1.count1, tw1.count2)
	}
	if tw1.count3 != 0 {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw1.count3", 0, tw1.count3)
	}

	if tw2.error != "" {
		t.Errorf("%s(%s) error: %s", c, "tw2.error", tw2.error)
	}
	if !tw2.closed {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw2.closed", true, tw2.closed)
	}
	if tw2.count1 != tw2.count2 {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw2.count1/2", tw2.count1, tw2.count2)
	}
	if tw2.count3 != 0 {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw2.count3", 0, tw2.count3)
	}

	for i := 0; i < len(cs); i++ {
		if cs[i] != tw2.counts[i] {
			t.Errorf("%s(%s) [%d] expect: %v, actual: %v", c, "tw2.counts", i, cs[i], tw2.counts[i])
		}
	}
}

func TestAsyncToAsync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(NewAsyncWriter(tw1, 1000))

	fmt.Println(time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println(time.Now(), "REPLACE")
	tw2 := newTestConcurrentDetectWriter()
	log.SwitchWriter(NewAsyncWriter(tw2, 100))

	fmt.Println(time.Now(), "WAIT")
	wg.Wait()

	fmt.Println(time.Now(), "CLOSE")
	log.Close()

	fmt.Println(time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestAsyncToAsync", tw1, tw2, counts)
}

func TestAsyncToSync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(NewAsyncWriter(tw1, 1000))

	fmt.Println(time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println(time.Now(), "REPLACE")
	tw2 := newTestConcurrentDetectWriter()
	log.SwitchWriter(NewSyncWriter(tw2))

	fmt.Println(time.Now(), "WAIT")
	wg.Wait()

	fmt.Println(time.Now(), "CLOSE")
	log.Close()

	fmt.Println(time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestAsyncToSync", tw1, tw2, counts)
}

func TestSyncToSync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(NewSyncWriter(tw1))

	fmt.Println(time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println(time.Now(), "REPLACE")
	tw2 := newTestConcurrentDetectWriter()
	log.SwitchWriter(NewSyncWriter(tw2))

	fmt.Println(time.Now(), "WAIT")
	wg.Wait()

	fmt.Println(time.Now(), "CLOSE")
	log.Close()

	fmt.Println(time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestSyncToSync", tw1, tw2, counts)
}

func TestSyncToAsync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(NewSyncWriter(tw1))

	fmt.Println(time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println(time.Now(), "REPLACE")
	tw2 := newTestConcurrentDetectWriter()
	log.SwitchWriter(NewAsyncWriter(tw2, 1000))

	fmt.Println(time.Now(), "WAIT")
	wg.Wait()

	fmt.Println(time.Now(), "CLOSE")
	log.Close()

	fmt.Println(time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestSyncToAsync", tw1, tw2, counts)
}
