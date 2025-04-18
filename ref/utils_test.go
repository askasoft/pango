package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func somefunction() {
	// this empty function is used by TestFunctionName()
}

func TestNameOfFunc(t *testing.T) {
	assert.Regexp(t, `^(.*/vendor/)?github.com/askasoft/pango/ref.somefunction$`, NameOfFunc(somefunction))
}
