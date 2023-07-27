package log

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type testFailoverWriter struct {
	error int
	count int
	msgs  []string
}

// Write do nothing.
func (tfw *testFailoverWriter) Write(le *Event) error {
	tfw.count++
	if tfw.count%tfw.error == 0 {
		return fmt.Errorf("testFailoverWriter: %d", tfw.count)
	}
	tfw.msgs = append(tfw.msgs, le.Msg)
	return nil
}

// Flush do nothing.
func (tfw *testFailoverWriter) Flush() {
}

// Close do nothing.
func (tfw *testFailoverWriter) Close() {
}

func TestFailover(t *testing.T) {
	log := NewLog()

	tfw := &testFailoverWriter{error: 3}
	log.SetWriter(NewFailoverWriter(tfw, 3))

	var msgs []string
	for i := 0; i < 10; i++ {
		msgs = append(msgs, strconv.Itoa(i))
		log.Error(i)
	}
	log.Close()

	if !reflect.DeepEqual(msgs, tfw.msgs) {
		t.Errorf("want %v\n but %v", msgs, tfw.msgs)
	}
}
