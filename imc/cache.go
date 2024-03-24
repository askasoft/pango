package imc

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/askasoft/pango/cas"
)

type Item struct {
	Object any
	Expiry int64
}

// Returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiry == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiry
}

type Cache struct {
	*cache
}

// Return a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than 1,
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than 1, expired items are not
// deleted from the cache before calling c.DeleteExpired().
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
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
func NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache {
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}

type cache struct {
	mu      sync.RWMutex
	de      time.Duration
	items   map[string]Item
	janitor *janitor
}

func (c *cache) expire(ds ...time.Duration) int64 {
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
func (c *cache) Set(k string, x any, d ...time.Duration) {
	e := c.expire(d...)

	c.mu.Lock()
	c.items[k] = Item{
		Object: x,
		Expiry: e,
	}
	c.mu.Unlock()
}

func (c *cache) set(k string, x any, e int64) {
	c.items[k] = Item{
		Object: x,
		Expiry: e,
	}
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
func (c *cache) Add(k string, x any, d ...time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("item '%s' already exists", k)
	}

	e := c.expire(d...)
	c.set(k, x, e)
	c.mu.Unlock()
	return nil
}

// Replace a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
func (c *cache) Replace(k string, x any, d ...time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("item %s doesn't exist", k)
	}

	e := c.expire(d...)
	c.set(k, x, e)
	c.mu.Unlock()
	return nil
}

// Increment an item of type int, int8, int16, int32, int64, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found and default value 'x[0]' is not supplied,
// or if it is not possible to increment it by n.
func (c *cache) Increment(k string, n any, x ...any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, found := c.items[k]
	if !found || v.Expired() {
		if len(x) > 0 {
			c.set(k, x[0], c.expire())
			return nil
		}
		return fmt.Errorf("item %s doesn't exist", k)
	}

	switch s := v.Object.(type) {
	case int:
		d, err := cas.ToInt(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case int8:
		d, err := cas.ToInt8(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case int16:
		d, err := cas.ToInt16(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case int32:
		d, err := cas.ToInt32(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case int64:
		d, err := cas.ToInt64(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case uint:
		d, err := cas.ToUint(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case uint8:
		d, err := cas.ToUint8(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case uint16:
		d, err := cas.ToUint16(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case uint32:
		d, err := cas.ToUint32(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case uint64:
		d, err := cas.ToUint64(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case float32:
		d, err := cas.ToFloat32(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	case float64:
		d, err := cas.ToFloat64(n)
		if err != nil {
			return err
		}
		v.Object = s + d
	default:
		return fmt.Errorf("item '%s' is not number", k)
	}

	c.items[k] = v
	return nil
}

// Decrement an item of type int, int8, int16, int32, int64, uint,
// uint8, uint32, or uint64, float32 or float64 by n. Returns an error if the
// item's value is not an integer, if it was not found and default value 'x[0]' is not supplied,
// or if it is not possible to decrement it by n.
func (c *cache) Decrement(k string, n any, x ...any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, found := c.items[k]
	if !found || v.Expired() {
		if len(x) > 0 {
			c.set(k, x[0], c.expire())
			return nil
		}
		return fmt.Errorf("item %s doesn't exist", k)
	}

	switch s := v.Object.(type) {
	case int:
		d, err := cas.ToInt(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case int8:
		d, err := cas.ToInt8(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case int16:
		d, err := cas.ToInt16(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case int32:
		d, err := cas.ToInt32(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case int64:
		d, err := cas.ToInt64(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case uint:
		d, err := cas.ToUint(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case uint8:
		d, err := cas.ToUint8(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case uint16:
		d, err := cas.ToUint16(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case uint32:
		d, err := cas.ToUint32(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case uint64:
		d, err := cas.ToUint64(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case float32:
		d, err := cas.ToFloat32(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	case float64:
		d, err := cas.ToFloat64(n)
		if err != nil {
			return err
		}
		v.Object = s - d
	default:
		return fmt.Errorf("item '%s' is not number", k)
	}

	c.items[k] = v
	return nil
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *cache) Get(k string) (any, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	if item.Expiry > 0 {
		if time.Now().UnixNano() > item.Expiry {
			c.mu.RUnlock()
			return nil, false
		}
	}

	c.mu.RUnlock()
	return item.Object, true
}

// GetWithExpiry returns an item and its expiration time from the cache.
// It returns the item or nil, the expiration time if one is set (if the item
// never expires a zero value for time.Time is returned), and a bool indicating
// whether the key was found.
func (c *cache) GetWithExpiry(k string) (any, time.Time, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, time.Time{}, false
	}

	if item.Expiry > 0 {
		if time.Now().UnixNano() > item.Expiry {
			c.mu.RUnlock()
			return nil, time.Time{}, false
		}

		// Return the item and the expiration time
		c.mu.RUnlock()
		return item.Object, time.Unix(0, item.Expiry), true
	}

	// If expiration <= 0 (i.e. no expiration time set) then return the item
	// and a zeroed time.Time
	c.mu.RUnlock()
	return item.Object, time.Time{}, true
}

func (c *cache) get(k string) (any, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}

	if item.Expiry > 0 {
		if time.Now().UnixNano() > item.Expiry {
			return nil, false
		}
	}
	return item.Object, true
}

// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *cache) Delete(k string) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}

// Delete all expired items from the cache.
func (c *cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiry > 0 && now > v.Expiry {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Copies all unexpired items in the cache into a new map and returns it.
func (c *cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiry > 0 {
			if now > v.Expiry {
				continue
			}
		}
		m[k] = v
	}
	return m
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func (c *cache) Count() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

// Delete all items from the cache.
func (c *cache) Clear() {
	c.mu.Lock()
	c.items = map[string]Item{}
	c.mu.Unlock()
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache) {
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

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		de:    de,
		items: m,
	}
	return c
}

func newCacheWithJanitor(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(de, m)

	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running DeleteExpired on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}
