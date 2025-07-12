package log

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type testRetryWriter struct {
	RetrySupport

	error int
	count int
	msgs  []string
}

func (trw *testRetryWriter) Write(le *Event) {
	trw.RetryWrite(le, trw.write)
}

func (trw *testRetryWriter) write(le *Event) error {
	trw.count++
	if trw.count%trw.error == 0 {
		return fmt.Errorf("testRetryWriter: %d", trw.count)
	}
	trw.msgs = append(trw.msgs, le.Message)
	return nil
}

func (trw *testRetryWriter) Flush() {
	trw.RetryFlush(trw.write)
}

func (trw *testRetryWriter) Close() {
	trw.Flush()
}

func TestRetrySupport(t *testing.T) {
	log := NewLog()

	trw := &testRetryWriter{error: 3}
	trw.Retries = 3
	log.SetWriter(trw)

	var msgs []string
	for i := 0; i < 10; i++ {
		msgs = append(msgs, strconv.Itoa(i))
		log.Error(i)
	}
	log.Close()

	if !reflect.DeepEqual(msgs, trw.msgs) {
		t.Errorf("want %v\n but %v", msgs, trw.msgs)
	}
}
