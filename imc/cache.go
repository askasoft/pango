package imc

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/askasoft/pango/num/mathx"
)

// Worker interface for cached item to control it's expiration
type Worker interface {
	// Working returns true to prevent expired
	Working() bool
}

type Item[T any] struct {
	Val T     // Cache Value
	TTL int64 // Time-To-Live (time.Unix)
}

// Working returns true to prevent expired
func (item Item[T]) Working() bool {
	var a any = item.Val
	if w, ok := a.(Worker); ok {
		return w.Working()
	}
	return false
}

// Expired Returns true if the item has expired.
func (item Item[T]) Expired() bool {
	if item.TTL <= 0 || item.Working() {
		return false
	}

	return time.Now().Unix() > item.TTL
}

// ExpiredAt Returns true if the item has expired at time `t`.
func (item Item[T]) ExpiredAt(t int64) bool {
	if item.TTL <= 0 || item.Working() {
		return false
	}

	return t > item.TTL
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
	return newCache(defaultExpires, cleanupInterval, items)
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
	return newCache(defaultExpires, cleanupInterval, items)
}

func newCache[T any](de time.Duration, ci time.Duration, m map[string]Item[T]) *Cache[T] {
	c := &cache[T]{
		de:    de,
		items: m,
	}

	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running DeleteExpired on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &Cache[T]{c}
	if ci > 0 {
		startJanitor(c, ci)
		runtime.SetFinalizer(C, stopJanitor[T])
	}
	return C
}

type janitor[T any] struct {
	interval time.Duration
	stopChan chan bool
}

func (j *janitor[T]) Run(c *cache[T]) {
	timer := time.NewTimer(j.interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.DeleteExpired()
			timer.Reset(j.interval)
		case <-j.stopChan:
			return
		}
	}
}

func stopJanitor[T any](c *Cache[T]) {
	c.janitor.stopChan <- true
}

func startJanitor[T any](c *cache[T], interval time.Duration) {
	j := &janitor[T]{
		interval: interval,
		stopChan: make(chan bool),
	}
	c.janitor = j

	go j.Run(c)
}

type cache[T any] struct {
	mu      sync.RWMutex
	de      time.Duration
	items   map[string]Item[T]
	janitor *janitor[T]
}

func (c *cache[T]) timeToLive(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	return time.Now().Add(d).Unix()
}

// Add an item to the cache, replacing any existing item.
// The cache's default expiration time is used.
func (c *cache[T]) Set(k string, x T) {
	c.SetWithExpires(k, x, c.de)
}

// Add an item to the cache, replacing any existing item.
// If the duration is 0, the cache's default expiration time is used.
// If it < 0, the item never expires.
func (c *cache[T]) SetWithExpires(k string, x T, d time.Duration) {
	t := c.timeToLive(d)

	c.mu.Lock()
	c.set(k, x, t)
	c.mu.Unlock()
}

func (c *cache[T]) set(k string, x T, t int64) {
	c.items[k] = Item[T]{Val: x, TTL: t}
}

func (c *cache[T]) get(k string) (T, bool) {
	item, found := c.items[k]
	if !found {
		var d T
		return d, false
	}

	if item.Expired() {
		var d T
		return d, false
	}
	return item.Val, true
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
func (c *cache[T]) Add(k string, x T) error {
	return c.AddWithExpires(k, x, c.de)
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
func (c *cache[T]) AddWithExpires(k string, x T, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("item '%s' already exists", k)
	}

	t := c.timeToLive(d)
	c.set(k, x, t)
	c.mu.Unlock()
	return nil
}

// Replace a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
func (c *cache[T]) Replace(k string, x T) error {
	return c.ReplaceWithExpires(k, x, c.de)
}

// Replace a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
func (c *cache[T]) ReplaceWithExpires(k string, x T, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("item '%s' doesn't exist", k)
	}

	t := c.timeToLive(d)
	c.set(k, x, t)
	c.mu.Unlock()
	return nil
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

	if item.Expired() {
		c.mu.RUnlock()
		var d T
		return d, false
	}

	c.mu.RUnlock()
	return item.Val, true
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

	if item.TTL <= 0 {
		// If expiration <= 0 (i.e. no expiration time set) then return the item
		// and a zeroed time.Time
		c.mu.RUnlock()
		return item.Val, time.Time{}, true
	}

	if item.Expired() {
		c.mu.RUnlock()
		var d T
		return d, time.Time{}, false
	}

	// Return the item and the expiration time
	c.mu.RUnlock()
	return item.Val, time.Unix(item.TTL, 0), true
}

// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *cache[T]) Delete(k string) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}

// Delete all expired items from the cache.
func (c *cache[T]) DeleteExpired() {
	c.mu.Lock()
	now := time.Now().Unix()
	for k, v := range c.items {
		if v.ExpiredAt(now) {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Items Copies all unexpired items in the cache into a new map and returns it.
func (c *cache[T]) Items() map[string]Item[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now().Unix()

	m := make(map[string]Item[T], len(c.items))
	for k, v := range c.items {
		if v.ExpiredAt(now) {
			continue
		}
		m[k] = v
	}
	return m
}

// Each Iterate all unexpired items in the cache to call function `f`.
// Function `f` returns false can break th iteration.
func (c *cache[T]) Each(f func(string, Item[T]) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now().Unix()

	for k, v := range c.items {
		if v.ExpiredAt(now) {
			continue
		}
		if !f(k, v) {
			break
		}
	}
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

// Increment an item of type int, int8, int16, int32, int64, uint,
// uint8, uint32, or uint64, float32 or float64 by n.
// Returns the incremented value.
func (c *cache[T]) Increment(k string, n T) T {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, c.timeToLive(c.de))
		return n
	}

	o, err := mathx.Add(v.Val, n)
	if err != nil {
		panic(fmt.Errorf("item '%s' cannot increment: %w", k, err))
	}

	v.Val = o.(T)
	c.items[k] = v
	return v.Val
}
