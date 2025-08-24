package imc

import (
	"testing"
	"time"
)

type TestStruct struct {
	Num      int
	Children []*TestStruct
}

func testNewCache() *Cache[string, any] {
	return New[string, any](0, 0)
}

func TestCache(t *testing.T) {
	tc := testNewCache()

	a, found := tc.Get("a")
	if found || a != nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, found := tc.Get("b")
	if found || b != nil {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, found := tc.Get("c")
	if found || c != nil {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	tc.Set("a", 1)
	tc.Set("b", "b")
	tc.Set("c", 3.5)

	x, found := tc.Get("a")
	if !found {
		t.Error("a was not found while getting a2")
	}
	if x == nil {
		t.Error("x for a is nil")
	} else if a2 := x.(int); a2+2 != 3 {
		t.Error("a2 (which should be 1) plus 2 does not equal 3; value:", a2)
	}

	x, found = tc.Get("b")
	if !found {
		t.Error("b was not found while getting b2")
	}
	if x == nil {
		t.Error("x for b is nil")
	} else if b2 := x.(string); b2+"B" != "bB" {
		t.Error("b2 (which should be b) plus B does not equal bB; value:", b2)
	}

	x, found = tc.Get("c")
	if !found {
		t.Error("c was not found while getting c2")
	}
	if x == nil {
		t.Error("x for c is nil")
	} else if c2 := x.(float64); c2+1.2 != 4.7 {
		t.Error("c2 (which should be 3.5) plus 1.2 does not equal 4.7; value:", c2)
	}
}

func TestCacheTimes(t *testing.T) {
	var found bool

	tc := New[string, any](time.Second, 100*time.Millisecond)
	tc.Set("a", 1)
	tc.SetWithTTL("b", 2, -1)
	tc.SetWithTTL("c", 3, 2*time.Second)
	tc.SetWithTTL("d", 4, 3*time.Second)

	time.Sleep(1500 * time.Millisecond)
	_, found = tc.Get("a")
	if found {
		t.Error("Found a when it should have been automatically deleted")
	}

	time.Sleep(1000 * time.Millisecond)
	_, found = tc.Get("c")
	if found {
		t.Error("Found c when it should have been automatically deleted")
	}

	_, found = tc.Get("b")
	if !found {
		t.Error("Did not find b even though it was set to never expire")
	}

	_, found = tc.Get("d")
	if !found {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	time.Sleep(1000 * time.Millisecond)
	_, found = tc.Get("d")
	if found {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}

func TestNewFrom(t *testing.T) {
	m := map[string]Item[int]{
		"a": {
			Val: 1,
		},
		"b": {
			Val: 2,
		},
	}
	tc := NewFrom(0, 0, m)
	a, found := tc.Get("a")
	if !found {
		t.Fatal("Did not find a")
	}
	if a != 1 {
		t.Fatal("a is not 1")
	}
	b, found := tc.Get("b")
	if !found {
		t.Fatal("Did not find b")
	}
	if b != 2 {
		t.Fatal("b is not 2")
	}
}

func TestStorePointerToStruct(t *testing.T) {
	tc := testNewCache()
	tc.Set("foo", &TestStruct{Num: 1})
	x, found := tc.Get("foo")
	if !found {
		t.Fatal("*TestStruct was not found for foo")
	}
	foo := x.(*TestStruct)
	foo.Num++

	y, found := tc.Get("foo")
	if !found {
		t.Fatal("*TestStruct was not found for foo (second time)")
	}
	bar := y.(*TestStruct)
	if bar.Num != 2 {
		t.Fatal("TestStruct.Num is not 2")
	}
}

func TestAdd(t *testing.T) {
	tc := testNewCache()

	ok := tc.Add("foo", "bar")
	if !ok {
		t.Error("Couldn't add foo even though it shouldn't exist")
	}

	ok = tc.Add("foo", "baz")
	if ok {
		t.Error("Successfully added another foo when it should have returned an error")
	}
}

func TestReplace(t *testing.T) {
	tc := testNewCache()

	ok := tc.Replace("foo", "bar")
	if ok {
		t.Error("Replaced foo when it shouldn't exist")
	}

	tc.Set("foo", "bar")
	ok = tc.Replace("foo", "bar")
	if !ok {
		t.Error("Couldn't replace existing key foo")
	}
}

func TestRemove(t *testing.T) {
	tc := testNewCache()

	tc.Set("foo", "bar")
	tc.Remove("foo")
	x, found := tc.Get("foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}
	if x != nil {
		t.Error("x is not nil:", x)
	}
}

func TestLen(t *testing.T) {
	tc := testNewCache()
	tc.Set("foo", "1")
	tc.Set("bar", "2")
	tc.Set("baz", "3")
	if n := tc.Len(); n != 3 {
		t.Errorf("Item count is not 3: %d", n)
	}
}

func TestClear(t *testing.T) {
	tc := testNewCache()
	tc.Set("foo", "bar")
	tc.Set("baz", "yes")
	tc.Clear()
	x, found := tc.Get("foo")
	if found {
		t.Error("foo was found, but it should have been deleted")
	}
	if x != nil {
		t.Error("x is not nil:", x)
	}
	x, found = tc.Get("baz")
	if found {
		t.Error("baz was found, but it should have been deleted")
	}
	if x != nil {
		t.Error("x is not nil:", x)
	}
}

func TestIncrementWithInt(t *testing.T) {
	tc := testNewCache()
	tc.Set("tint", 1)
	x := tc.Increment("tint", 2)
	if x != int(3) {
		t.Error("tint is not 3:", x)
	}
}

func TestIncrementWithInt8(t *testing.T) {
	tc := testNewCache()
	tc.Set("tint8", int8(1))
	x := tc.Increment("tint8", int8(2))
	if x != int8(3) {
		t.Error("tint8 is not 3:", x)
	}
}

func TestIncrementWithInt16(t *testing.T) {
	tc := testNewCache()
	tc.Set("tint16", int16(1))
	x := tc.Increment("tint16", int16(2))
	if x != int16(3) {
		t.Error("tint16 is not 3:", x)
	}
}

func TestIncrementWithInt32(t *testing.T) {
	tc := testNewCache()
	tc.Set("tint32", int32(1))
	x := tc.Increment("tint32", int32(2))
	if x != int32(3) {
		t.Error("tint32 is not 3:", x)
	}
}

func TestIncrementWithInt64(t *testing.T) {
	tc := testNewCache()
	tc.Set("tint64", int64(1))
	x := tc.Increment("tint64", int64(2))
	if x != int64(3) {
		t.Error("tint64 is not 3:", x)
	}
}

func TestIncrementWithUint(t *testing.T) {
	tc := testNewCache()
	tc.Set("tuint", uint(1))
	x := tc.Increment("tuint", uint(2))
	if x != uint(3) {
		t.Error("tuint is not 3:", x)
	}
}

func TestIncrementWithUint8(t *testing.T) {
	tc := testNewCache()
	tc.Set("tuint8", uint8(1))
	x := tc.Increment("tuint8", uint8(2))
	if x != uint8(3) {
		t.Error("tuint8 is not 3:", x)
	}
}

func TestIncrementWithUint16(t *testing.T) {
	tc := testNewCache()
	tc.Set("tuint16", uint16(1))
	x := tc.Increment("tuint16", uint16(2))
	if x != uint16(3) {
		t.Error("tuint16 is not 3:", x)
	}
}

func TestIncrementWithUint32(t *testing.T) {
	tc := testNewCache()
	tc.Set("tuint32", uint32(1))
	x := tc.Increment("tuint32", uint32(2))
	if x != uint32(3) {
		t.Error("tuint32 is not 3:", x)
	}
}

func TestIncrementWithUint64(t *testing.T) {
	tc := testNewCache()
	tc.Set("tuint64", uint64(1))
	x := tc.Increment("tuint64", uint64(2))
	if x != uint64(3) {
		t.Error("tuint64 is not 3:", x)
	}
}

func TestIncrementWithFloat32(t *testing.T) {
	tc := testNewCache()
	tc.Set("float32", float32(1.5))
	x := tc.Increment("float32", 2)
	if x != float32(3.5) {
		t.Error("float32 is not 3.5:", x)
	}
}

func TestIncrementWithFloat64(t *testing.T) {
	tc := testNewCache()
	tc.Set("float64", float64(1.5))
	x := tc.Increment("float64", 2)
	if x != float64(3.5) {
		t.Error("float64 is not 3.5:", x)
	}
}

func TestGetWithTTL(t *testing.T) {
	tc := New[string, any](0, 0)

	a, expiration, found := tc.GetWithTTL("a")
	if found || a != nil || !expiration.IsZero() {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, expiration, found := tc.GetWithTTL("b")
	if found || b != nil || !expiration.IsZero() {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, expiration, found := tc.GetWithTTL("c")
	if found || c != nil || !expiration.IsZero() {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	tc.Set("a", 1)
	tc.Set("b", "b")
	tc.Set("c", 3.5)
	tc.SetWithTTL("d", 1, -1)
	tc.SetWithTTL("e", 1, time.Second)

	x, expiration, found := tc.GetWithTTL("a")
	if !found {
		t.Error("a was not found while getting a2")
	}
	if x == nil {
		t.Error("x for a is nil")
	} else if a2 := x.(int); a2+2 != 3 {
		t.Error("a2 (which should be 1) plus 2 does not equal 3; value:", a2)
	}
	if !expiration.IsZero() {
		t.Error("expiration for a is not a zeroed time")
	}

	x, expiration, found = tc.GetWithTTL("b")
	if !found {
		t.Error("b was not found while getting b2")
	}
	if x == nil {
		t.Error("x for b is nil")
	} else if b2 := x.(string); b2+"B" != "bB" {
		t.Error("b2 (which should be b) plus B does not equal bB; value:", b2)
	}
	if !expiration.IsZero() {
		t.Error("expiration for b is not a zeroed time")
	}

	x, expiration, found = tc.GetWithTTL("c")
	if !found {
		t.Error("c was not found while getting c2")
	}
	if x == nil {
		t.Error("x for c is nil")
	} else if c2 := x.(float64); c2+1.2 != 4.7 {
		t.Error("c2 (which should be 3.5) plus 1.2 does not equal 4.7; value:", c2)
	}
	if !expiration.IsZero() {
		t.Error("expiration for c is not a zeroed time")
	}

	x, expiration, found = tc.GetWithTTL("d")
	if !found {
		t.Error("d was not found while getting d2")
	}
	if x == nil {
		t.Error("x for d is nil")
	} else if d2 := x.(int); d2+2 != 3 {
		t.Error("d (which should be 1) plus 2 does not equal 3; value:", d2)
	}
	if !expiration.IsZero() {
		t.Error("expiration for d is not a zeroed time")
	}

	x, expiration, found = tc.GetWithTTL("e")
	if !found {
		t.Error("e was not found while getting e2")
	}
	if x == nil {
		t.Error("x for e is nil")
	} else if e2 := x.(int); e2+2 != 3 {
		t.Error("e (which should be 1) plus 2 does not equal 3; value:", e2)
	}
	if expiration.UnixMilli() != tc.items["e"].TTL {
		t.Error("expiration for e is not the correct time")
	}
	if expiration.UnixMilli() < time.Now().UnixMilli() {
		t.Error("expiration for e is in the past")
	}
}
