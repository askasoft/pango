package ref

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringEmptyIsNotNil(t *testing.T) {
	var s interface{}
	s = ""
	assert.True(t, s != nil)
}

func TestStringZeroIsNotNil(t *testing.T) {
	var s string
	var i interface{}
	i = s
	assert.True(t, i != nil)
	assert.True(t, i == "")
}

func TestConvertFloatToString(t *testing.T) {
	v, err := Convert(1.123, reflect.TypeOf(""))
	assert.Nil(t, err)
	assert.Equal(t, "1.123", v.(string))
}

func TestConvertStringToInt(t *testing.T) {
	v, err := Convert("0777", reflect.TypeOf(int32(0)))
	assert.Nil(t, err)
	assert.Equal(t, int32(0777), v.(int32))
}

func TestConvertNilToInt(t *testing.T) {
	v, err := Convert(nil, reflect.TypeOf(int(0)))
	assert.Nil(t, err)
	assert.Equal(t, int(0), v.(int))
}

func TestConvertZeroStrToInt32(t *testing.T) {
	v, err := Convert("", reflect.TypeOf(int32(0)))
	assert.Nil(t, err)
	assert.Equal(t, int32(0), v.(int32))
}
