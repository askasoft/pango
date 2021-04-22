package ref

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToInt(t *testing.T) {
	var i int32
	v, err := Convert("0777", reflect.TypeOf(i))
	assert.Nil(t, err)
	assert.Equal(t, int32(0777), int32(v.Int()))
}
