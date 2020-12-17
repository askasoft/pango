package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogConfigJson(t *testing.T) {
	log := Default()
	assert.Nil(t, Config(log, "testdata/log.json"))
}
