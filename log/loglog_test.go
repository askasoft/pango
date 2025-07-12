package log

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/askasoft/pango/str"
)

const testRoutines = 10

type testConcurrentDetectWriter struct {
	countStack  uint64
	countAtomic uint64
	countClosed uint64

	counts [testRoutines]int64
	closed bool
	last   time.Time
}

func newTestConcurrentDetectWriter() *testConcurrentDetectWriter {
	tw := &testConcurrentDetectWriter{}
	tw.last = time.Now()
	return tw
}

func (tw *testConcurrentDetectWriter) do() {
	if tw.closed {
		atomic.AddUint64(&tw.countClosed, 1)
	} else {
		tw.countStack++
		atomic.AddUint64(&tw.countAtomic, 1)
	}
}

func (tw *testConcurrentDetectWriter) Write(le *Event) {
	tw.do()

	ss := str.Split(le.Message, " ")
	k, _ := strconv.Atoi(ss[0])
	c, _ := strconv.ParseInt(ss[1], 10, 64)
	tw.counts[k] = c

	t := time.Now()
	if t.After(tw.last.Add(time.Second)) {
		fmt.Println(le.Time, k, c, tw.countAtomic)
		tw.last = t
	}
}

func (tw *testConcurrentDetectWriter) Flush() {
	tw.do()
}

func (tw *testConcurrentDetectWriter) Close() {
	tw.closed = true
	tw.do()
}

func testLogRoutine(log *Log, wg *sync.WaitGroup, n int, p *int64) {
	var m int64

	et := time.Now().Add(time.Second * 5)
	for time.Now().Before(et) {
		m++
		log.Infof("%d %d", n, m)
		time.Sleep(time.Millisecond * 10)
	}
	*p = m

	fmt.Println(time.Now(), "test routine ", n, " done ", m)
	wg.Done()
}

func testCheckConcurrentDetectWriter(t *testing.T, c string, tw1 *testConcurrentDetectWriter, tw2 *testConcurrentDetectWriter, cs []int64) {
	fmt.Printf("tw1: %v, CS: %d, CA: %d, CC:%d\n", tw1.closed, tw1.countStack, tw1.countAtomic, tw1.countClosed)
	fmt.Println("tw1: ", tw1.counts)
	fmt.Printf("tw2: %v, CS: %d, CA: %d, CC:%d\n", tw2.closed, tw2.countStack, tw2.countAtomic, tw2.countClosed)
	fmt.Println("tw2: ", tw2.counts)

	if !tw1.closed {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw1.closed", true, tw1.closed)
	}
	if tw1.countStack != tw1.countAtomic {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw1.count(stack != atomic)", tw1.countStack, tw1.countAtomic)
	}
	if tw1.countClosed != 1 {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw1.countClosed", 1, tw1.countClosed)
	}

	if !tw2.closed {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw2.closed", true, tw2.closed)
	}
	if tw2.countStack != tw2.countAtomic {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw2.count(stack != atomic)", tw2.countStack, tw2.countAtomic)
	}
	if tw2.countClosed != 1 {
		t.Errorf("%s(%s) expect: %v, actual: %v", c, "tw2.countClosed", 1, tw2.countClosed)
	}

	for i := 0; i < len(cs); i++ {
		if cs[i] != tw2.counts[i] {
			t.Errorf("%s(%s) [%d] expect: %v, actual: %v", c, "tw2.counts", i, cs[i], tw2.counts[i])
		}
	}
}

func TestAsyncToAsync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()
	aw1 := NewAsyncWriter(tw1, 1000)

	log := NewLog()
	log.SetWriter(aw1)

	fmt.Println("A2A", time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println("A2A", time.Now(), "SWITCH")
	tw2 := newTestConcurrentDetectWriter()
	aw2 := NewAsyncWriter(tw2, 100)
	log.SwitchWriter(aw2)

	fmt.Println("A2A", time.Now(), "WAIT")
	wg.Wait()

	fmt.Println("A2A", time.Now(), "CLOSE")
	log.Close()

	aw1.Wait()
	aw2.Wait()

	fmt.Println("A2A", time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestAsyncToAsync", tw1, tw2, counts)
}

func TestAsyncToSync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()
	aw1 := NewAsyncWriter(tw1, 1000)

	log := NewLog()
	log.SetWriter(aw1)

	fmt.Println("A2S", time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println("A2S", time.Now(), "SWITCH")
	tw2 := newTestConcurrentDetectWriter()
	log.SwitchWriter(NewSyncWriter(tw2))

	fmt.Println("A2S", time.Now(), "WAIT")
	wg.Wait()

	fmt.Println("A2S", time.Now(), "CLOSE")
	log.Close()

	aw1.Wait()

	fmt.Println("A2S", time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestAsyncToSync", tw1, tw2, counts)
}

func TestSyncToSync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()

	log := NewLog()
	log.SetWriter(NewSyncWriter(tw1))

	fmt.Println("S2S", time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println("S2S", time.Now(), "SWITCH")
	tw2 := newTestConcurrentDetectWriter()
	log.SwitchWriter(NewSyncWriter(tw2))

	fmt.Println("S2S", time.Now(), "WAIT")
	wg.Wait()

	fmt.Println("S2S", time.Now(), "CLOSE")
	log.Close()

	fmt.Println("S2S", time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestSyncToSync", tw1, tw2, counts)
}

func TestSyncToAsync(t *testing.T) {
	tw1 := newTestConcurrentDetectWriter()

	log := NewLog()
	log.SetWriter(NewSyncWriter(tw1))

	fmt.Println("S2A", time.Now(), "START")
	wg := &sync.WaitGroup{}
	counts := make([]int64, testRoutines)
	for i := 0; i < len(counts); i++ {
		wg.Add(1)
		go testLogRoutine(log, wg, i, &counts[i])
	}

	time.Sleep(time.Second * 2)

	fmt.Println("S2A", time.Now(), "SWITCH")
	tw2 := newTestConcurrentDetectWriter()
	aw2 := NewAsyncWriter(tw2, 1000)
	log.SwitchWriter(aw2)

	fmt.Println("S2A", time.Now(), "WAIT")
	wg.Wait()

	fmt.Println("S2A", time.Now(), "CLOSE")
	log.Close()

	aw2.Wait()

	fmt.Println("S2A", time.Now(), "END")

	testCheckConcurrentDetectWriter(t, "TestSyncToAsync", tw1, tw2, counts)
}
