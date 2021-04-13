package col

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicFeatures(t *testing.T) {
	n := 100
	om := NewOrderedMap()

	// set(i, 2 * i)
	for i := 0; i < n; i++ {
		assertLenEqual(t, om, i)
		old, ok := om.Set(i, 2*i)
		assertLenEqual(t, om, i+1)

		assert.Nil(t, old)
		assert.False(t, ok)
	}

	// get what we just set
	for i := 0; i < n; i++ {
		value, ok := om.Get(i)

		assert.Equal(t, 2*i, value)
		assert.True(t, ok)
	}

	// get entries of what we just set
	for i := 0; i < n; i++ {
		entry := om.GetEntry(i)

		assert.NotNil(t, entry)
		assert.Equal(t, 2*i, entry.Value)
	}

	// forward iteration
	i := 0
	for entry := om.Front(); entry != nil; entry = entry.Next() {
		assert.Equal(t, i, entry.key)
		assert.Equal(t, 2*i, entry.Value)
		i++
	}

	// backward iteration
	i = n - 1
	for entry := om.Back(); entry != nil; entry = entry.Prev() {
		assert.Equal(t, i, entry.key)
		assert.Equal(t, 2*i, entry.Value)
		i--
	}

	// forward iteration starting from known key
	i = 42
	for entry := om.GetEntry(i); entry != nil; entry = entry.Next() {
		assert.Equal(t, i, entry.key)
		assert.Equal(t, 2*i, entry.Value)
		i++
	}

	// double values for entries with even keys
	for j := 0; j < n/2; j++ {
		i = 2 * j
		old, ok := om.Set(i, 4*i)

		assert.Equal(t, 2*i, old)
		assert.True(t, ok)
	}

	// and delete entries with odd keys
	for j := 0; j < n/2; j++ {
		i = 2*j + 1
		assertLenEqual(t, om, n-j)
		value, ok := om.Remove(i)
		assertLenEqual(t, om, n-j-1)

		assert.Equal(t, 2*i, value)
		assert.True(t, ok)

		// deleting again shouldn't change anything
		value, ok = om.Remove(i)
		assertLenEqual(t, om, n-j-1)
		assert.Nil(t, value)
		assert.False(t, ok)
	}

	// get the whole range
	for j := 0; j < n/2; j++ {
		i = 2 * j
		value, ok := om.Get(i)
		assert.Equal(t, 4*i, value)
		assert.True(t, ok)

		i = 2*j + 1
		value, ok = om.Get(i)
		assert.Nil(t, value)
		assert.False(t, ok)
	}

	// check iterations again
	i = 0
	for entry := om.Front(); entry != nil; entry = entry.Next() {
		assert.Equal(t, i, entry.key)
		assert.Equal(t, 4*i, entry.Value)
		i += 2
	}
	i = 2 * ((n - 1) / 2)
	for entry := om.Back(); entry != nil; entry = entry.Prev() {
		assert.Equal(t, i, entry.key)
		assert.Equal(t, 4*i, entry.Value)
		i -= 2
	}
}

func TestUpdatingDoesntChangePairsOrder(t *testing.T) {
	om := NewOrderedMap()
	om.Set("foo", "bar")
	om.Set(12, 28)
	om.Set(78, 100)
	om.Set("bar", "baz")

	old, ok := om.Set(78, 102)
	assert.Equal(t, 100, old)
	assert.True(t, ok)

	assertOrderedPairsEqual(t, om,
		[]interface{}{"foo", 12, 78, "bar"},
		[]interface{}{"bar", 28, 102, "baz"})
}

func TestDeletingAndReinsertingChangesPairsOrder(t *testing.T) {
	om := NewOrderedMap()
	om.Set("foo", "bar")
	om.Set(12, 28)
	om.Set(78, 100)
	om.Set("bar", "baz")

	// delete a entry
	old, ok := om.Remove(78)
	assert.Equal(t, 100, old)
	assert.True(t, ok)

	// re-insert the same entry
	old, ok = om.Set(78, 100)
	assert.Nil(t, old)
	assert.False(t, ok)

	assertOrderedPairsEqual(t, om,
		[]interface{}{"foo", 12, "bar", 78},
		[]interface{}{"bar", 28, "baz", 100})
}

func TestEmptyMapOperations(t *testing.T) {
	om := NewOrderedMap()

	old, ok := om.Get("foo")
	assert.Nil(t, old)
	assert.False(t, ok)

	old, ok = om.Remove("bar")
	assert.Nil(t, old)
	assert.False(t, ok)

	assertLenEqual(t, om, 0)

	assert.Nil(t, om.Front())
	assert.Nil(t, om.Back())
}

type dummyTestStruct struct {
	value string
}

func TestPackUnpackStructs(t *testing.T) {
	om := NewOrderedMap()
	om.Set("foo", dummyTestStruct{"foo!"})
	om.Set("bar", dummyTestStruct{"bar!"})

	value, ok := om.Get("foo")
	assert.True(t, ok)
	if assert.NotNil(t, value) {
		assert.Equal(t, "foo!", value.(dummyTestStruct).value)
	}

	value, ok = om.Set("bar", dummyTestStruct{"baz!"})
	assert.True(t, ok)
	if assert.NotNil(t, value) {
		assert.Equal(t, "bar!", value.(dummyTestStruct).value)
	}

	value, ok = om.Get("bar")
	assert.True(t, ok)
	if assert.NotNil(t, value) {
		assert.Equal(t, "baz!", value.(dummyTestStruct).value)
	}
}

func TestShuffle(t *testing.T) {
	ranLen := 100

	for _, n := range []int{0, 10, 20, 100, 1000, 10000} {
		t.Run(fmt.Sprintf("shuffle test with %d items", n), func(t *testing.T) {
			om := NewOrderedMap()

			keys := make([]interface{}, n)
			values := make([]interface{}, n)

			for i := 0; i < n; i++ {
				// we prefix with the number to ensure that we don't get any duplicates
				keys[i] = fmt.Sprintf("%d_%s", i, randomHexString(t, ranLen))
				values[i] = randomHexString(t, ranLen)

				value, ok := om.Set(keys[i], values[i])
				assert.Nil(t, value)
				assert.False(t, ok)
			}

			assertOrderedPairsEqual(t, om, keys, values)
		})
	}
}

/* Test helpers */

func assertOrderedPairsEqual(t *testing.T, om *OrderedMap, expectedKeys, expectedValues []interface{}) {
	assertOrderedPairsEqualFromNewest(t, om, expectedKeys, expectedValues)
	assertOrderedPairsEqualFromOldest(t, om, expectedKeys, expectedValues)
}

func assertOrderedPairsEqualFromNewest(t *testing.T, om *OrderedMap, expectedKeys, expectedValues []interface{}) {
	if assert.Equal(t, len(expectedKeys), len(expectedValues)) && assert.Equal(t, len(expectedKeys), om.Len()) {
		i := om.Len() - 1
		for entry := om.Back(); entry != nil; entry = entry.Prev() {
			assert.Equal(t, expectedKeys[i], entry.key)
			assert.Equal(t, expectedValues[i], entry.Value)
			i--
		}
	}
}

func assertOrderedPairsEqualFromOldest(t *testing.T, om *OrderedMap, expectedKeys, expectedValues []interface{}) {
	if assert.Equal(t, len(expectedKeys), len(expectedValues)) && assert.Equal(t, len(expectedKeys), om.Len()) {
		i := om.Len() - 1
		for entry := om.Back(); entry != nil; entry = entry.Prev() {
			assert.Equal(t, expectedKeys[i], entry.key)
			assert.Equal(t, expectedValues[i], entry.Value)
			i--
		}
	}
}

func assertLenEqual(t *testing.T, om *OrderedMap, expectedLen int) {
	assert.Equal(t, expectedLen, om.Len())

	// also check the list length, for good measure
	assert.Equal(t, expectedLen, om.list.Len())
}

func randomHexString(t *testing.T, length int) string {
	b := length / 2
	randBytes := make([]byte, b)

	if n, err := rand.Read(randBytes); err != nil || n != b {
		if err == nil {
			err = fmt.Errorf("only got %v random bytes, expected %v", n, b)
		}
		t.Fatal(err)
	}

	return hex.EncodeToString(randBytes)
}
