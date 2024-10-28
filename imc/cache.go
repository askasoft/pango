package imc

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/askasoft/pango/cas"
)

type Item[T any] struct {
	Object  T     // Cache Object
	Expires int64 // UnixNano
}

// Returns true if the item has expired.
func (item Item[T]) Expired() bool {
	if item.Expires == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expires
}

type Cache[T any] struct {
	*cache[T]
}

// Return a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than 1,
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than 1, expired items are not
// deleted from the cache before calling c.DeleteExpired().
func New[T any](defaultExpires, cleanupInterval time.Duration) *Cache[T] {
	items := make(map[string]Item[T])
	return newCacheWithJanitor(defaultExpires, cleanupInterval, items)
}

// Return a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than 1,
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than 1, expired items are not
// deleted from the cache before calling c.DeleteExpired().
//
// NewFrom() also accepts an items map which will serve as the underlying map
// for the cache. This is useful for starting from a preallocated cache
// like make(map[string]Item, 500) to improve startup performance when the cache
// is expected to reach a certain minimum size.
//
// Only the cache's methods synchronize access to this map, so it is not
// recommended to keep any references to the map around after creating a cache.
// If need be, the map can be accessed at a later point using c.Items() (subject
// to the same caveat.)
func NewFrom[T any](defaultExpires, cleanupInterval time.Duration, items map[string]Item[T]) *Cache[T] {
	return newCacheWithJanitor(defaultExpires, cleanupInterval, items)
}

type cache[T any] struct {
	mu      sync.RWMutex
	de      time.Duration
	items   map[string]Item[T]
	janitor *janitor[T]
}

func (c *cache[T]) expires(ds ...time.Duration) int64 {
	var d time.Duration
	var e int64

	if len(ds) > 0 {
		d = ds[0]
	}

	if d == 0 {
		d = c.de
	}

	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	return e
}

// Add an item to the cache, replacing any existing item.
// If the duration is 0, the cache's default expiration time is used.
// If it < 0, the item never expires.
func (c *cache[T]) Set(k string, x T, d ...time.Duration) {
	e := c.expires(d...)

	c.mu.Lock()
	c.items[k] = Item[T]{
		Object:  x,
		Expires: e,
	}
	c.mu.Unlock()
}

func (c *cache[T]) set(k string, x T, e int64) {
	c.items[k] = Item[T]{
		Object:  x,
		Expires: e,
	}
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
func (c *cache[T]) Add(k string, x T, d ...time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("item '%s' already exists", k)
	}

	e := c.expires(d...)
	c.set(k, x, e)
	c.mu.Unlock()
	return nil
}

// Replace a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
func (c *cache[T]) Replace(k string, x T, d ...time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("item '%s' doesn't exist", k)
	}

	e := c.expires(d...)
	c.set(k, x, e)
	c.mu.Unlock()
	return nil
}

// Increment an item of type int, int8, int16, int32, int64, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found and default value 'x[0]' is not supplied,
// or if it is not possible to increment it by n.
func (c *cache[T]) Increment(k string, n T, x ...T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, found := c.items[k]
	if !found || v.Expired() {
		if len(x) > 0 {
			c.set(k, x[0], c.expires())
			return nil
		}
		return fmt.Errorf("item '%s' doesn't exist", k)
	}

	o, err := c.inc(v.Object, n)
	if err != nil {
		return fmt.Errorf("item '%s' cannot increment: %w", k, err)
	}

	v.Object = o.(T)
	c.items[k] = v
	return nil
}

// Decrement an item of type int, int8, int16, int32, int64, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found and default value 'x[0]' is not supplied,
// or if it is not possible to decrement it by n.
func (c *cache[T]) Decrement(k string, n T, x ...T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, found := c.items[k]
	if !found || v.Expired() {
		if len(x) > 0 {
			c.set(k, x[0], c.expires())
			return nil
		}
		return fmt.Errorf("item '%s' doesn't exist", k)
	}

	o, err := c.dec(v.Object, n)
	if err != nil {
		return fmt.Errorf("item '%s' cannot decrement: %w", k, err)
	}

	v.Object = o.(T)
	c.items[k] = v
	return nil
}

func (c *cache[T]) inc(o any, n any) (any, error) {
	switch s := o.(type) {
	case int:
		d, err := cas.ToInt(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case int8:
		d, err := cas.ToInt8(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case int16:
		d, err := cas.ToInt16(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case int32:
		d, err := cas.ToInt32(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case int64:
		d, err := cas.ToInt64(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case uint:
		d, err := cas.ToUint(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case uint8:
		d, err := cas.ToUint8(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case uint16:
		d, err := cas.ToUint16(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case uint32:
		d, err := cas.ToUint32(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case uint64:
		d, err := cas.ToUint64(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case float32:
		d, err := cas.ToFloat32(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	case float64:
		d, err := cas.ToFloat64(n)
		if err != nil {
			return o, err
		}
		return s + d, nil
	default:
		return o, fmt.Errorf("'%T' is not number", o)
	}
}

func (c *cache[T]) dec(o any, n any) (any, error) {
	switch s := o.(type) {
	case int:
		d, err := cas.ToInt(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case int8:
		d, err := cas.ToInt8(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case int16:
		d, err := cas.ToInt16(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case int32:
		d, err := cas.ToInt32(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case int64:
		d, err := cas.ToInt64(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case uint:
		d, err := cas.ToUint(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case uint8:
		d, err := cas.ToUint8(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case uint16:
		d, err := cas.ToUint16(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case uint32:
		d, err := cas.ToUint32(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case uint64:
		d, err := cas.ToUint64(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case float32:
		d, err := cas.ToFloat32(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	case float64:
		d, err := cas.ToFloat64(n)
		if err != nil {
			return o, err
		}
		return s - d, nil
	default:
		return o, fmt.Errorf("'%T' is not number", s)
	}
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *cache[T]) Get(k string) (T, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		var d T
		return d, false
	}

	if item.Expires > 0 {
		if time.Now().UnixNano() > item.Expires {
			c.mu.RUnlock()
			var d T
			return d, false
		}
	}

	c.mu.RUnlock()
	return item.Object, true
}

// GetWithExpires returns an item and its expiration time from the cache.
// It returns the item or nil, the expiration time if one is set (if the item
// never expires a zero value for time.Time is returned), and a bool indicating
// whether the key was found.
func (c *cache[T]) GetWithExpires(k string) (T, time.Time, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		var d T
		return d, time.Time{}, false
	}

	if item.Expires > 0 {
		if time.Now().UnixNano() > item.Expires {
			c.mu.RUnlock()
			var d T
			return d, time.Time{}, false
		}

		// Return the item and the expiration time
		c.mu.RUnlock()
		return item.Object, time.Unix(0, item.Expires), true
	}

	// If expiration <= 0 (i.e. no expiration time set) then return the item
	// and a zeroed time.Time
	c.mu.RUnlock()
	return item.Object, time.Time{}, true
}

func (c *cache[T]) get(k string) (T, bool) {
	item, found := c.items[k]
	if !found {
		var d T
		return d, false
	}

	if item.Expires > 0 {
		if time.Now().UnixNano() > item.Expires {
			var d T
			return d, false
		}
	}
	return item.Object, true
}

// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *cache[T]) Delete(k string) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}

// Delete all expired items from the cache.
func (c *cache[T]) DeleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.items {
		if v.Expires > 0 && now > v.Expires {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Copies all unexpired items in the cache into a new map and returns it.
func (c *cache[T]) Items() map[string]Item[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m := make(map[string]Item[T], len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expires > 0 {
			if now > v.Expires {
				continue
			}
		}
		m[k] = v
	}
	return m
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func (c *cache[T]) Count() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

// Delete all items from the cache.
func (c *cache[T]) Clear() {
	c.mu.Lock()
	c.items = map[string]Item[T]{}
	c.mu.Unlock()
}

type janitor[T any] struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor[T]) Run(c *cache[T]) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor[T any](c *Cache[T]) {
	c.janitor.stop <- true
}

func runJanitor[T any](c *cache[T], ci time.Duration) {
	j := &janitor[T]{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func newCache[T any](de time.Duration, m map[string]Item[T]) *cache[T] {
	if de == 0 {
		de = -1
	}
	c := &cache[T]{
		de:    de,
		items: m,
	}
	return c
}

func newCacheWithJanitor[T any](de time.Duration, ci time.Duration, m map[string]Item[T]) *Cache[T] {
	c := newCache(de, m)

	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running DeleteExpired on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &Cache[T]{c}
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(C, stopJanitor[T])
	}
	return C
}
