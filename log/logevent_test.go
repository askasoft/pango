package log

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventCaller(t *testing.T) {
	le := newEvent(&logger{}, LevelInfo, "caller")
	le.When = time.Time{}
	le.Caller(2, false)
	b, err := json.Marshal(le)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, `{"level":4,"msg":"caller","when":"0001-01-01T00:00:00Z","file":"logevent_test.go","line":`+
		strconv.Itoa(le.Line)+`,"func":"log.TestEventCaller","trace":""}`, string(b))
}

func TestEventJsonMarshal(t *testing.T) {
	le := newEvent(&logger{}, LevelInfo, "marshal")
	le.When = time.Time{}
	b, err := json.Marshal(le)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, `{"level":4,"msg":"marshal","when":"0001-01-01T00:00:00Z","file":"","line":0,"func":"","trace":""}`, string(b))
}
