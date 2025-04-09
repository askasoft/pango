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

type Item[V any] struct {
	Val V     // Cache Value
	TTL int64 // Time-To-Live (time.Unix)
}

// Working returns true to prevent expired
func (item Item[V]) Working() bool {
	var a any = item.Val
	if w, ok := a.(Worker); ok {
		return w.Working()
	}
	return false
}

// Expired Returns true if the item has expired.
func (item Item[V]) Expired() bool {
	if item.TTL <= 0 || item.Working() {
		return false
	}

	return time.Now().Unix() > item.TTL
}

// ExpiredAt Returns true if the item has expired at time `t`.
func (item Item[V]) ExpiredAt(t int64) bool {
	if item.TTL <= 0 || item.Working() {
		return false
	}

	return t > item.TTL
}

type Cache[K comparable, V any] struct {
	*cache[K, V]
}

// Return a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than 1,
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than 1, expired items are not
// deleted from the cache before calling c.Clean().
func New[K comparable, V any](defaultExpires, cleanupInterval time.Duration) *Cache[K, V] {
	return newCache(defaultExpires, cleanupInterval, make(map[K]Item[V]))
}

// Return a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than 1,
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than 1, expired items are not
// deleted from the cache before calling c.Clean().
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
func NewFrom[K comparable, V any](defaultExpires, cleanupInterval time.Duration, items map[K]Item[V]) *Cache[K, V] {
	return newCache(defaultExpires, cleanupInterval, items)
}

func newCache[K comparable, V any](de time.Duration, ci time.Duration, m map[K]Item[V]) *Cache[K, V] {
	c := &cache[K, V]{
		ttl:   de,
		items: m,
	}

	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running Clean on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &Cache[K, V]{c}
	if ci > 0 {
		startJanitor(C, ci)
		runtime.SetFinalizer(C, stopJanitor[K, V])
	}
	return C
}

func startJanitor[K comparable, V any](c *Cache[K, V], ci time.Duration) {
	j := &janitor{
		interval: ci,
		stopChan: make(chan bool),
	}
	c.janitor = j

	go j.Run(c.Clean)
}

func stopJanitor[K comparable, V any](c *Cache[K, V]) {
	c.janitor.stopChan <- true
}

type janitor struct {
	interval time.Duration
	stopChan chan bool
}

func (j *janitor) Run(f func()) {
	timer := time.NewTimer(j.interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			f()
			timer.Reset(j.interval)
		case <-j.stopChan:
			return
		}
	}
}

type cache[K comparable, V any] struct {
	mu      sync.RWMutex
	ttl     time.Duration
	items   map[K]Item[V]
	janitor *janitor
}

func (c *cache[K, V]) timeToLive(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	return time.Now().Add(d).Unix()
}

func (c *cache[K, V]) set(k K, v V, t int64) {
	c.items[k] = Item[V]{Val: v, TTL: t}
}

func (c *cache[K, V]) get(k K) (V, bool) {
	item, found := c.items[k]
	if !found {
		var d V
		return d, false
	}

	if item.Expired() {
		var d V
		return d, false
	}
	return item.Val, true
}

// Len Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func (c *cache[K, V]) Len() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

// Set an item to the cache, replacing any existing item.
// The cache's default expiration time is used.
func (c *cache[K, V]) Set(k K, v V) {
	c.SetWithTTL(k, v, c.ttl)
}

// Set an item to the cache, replacing any existing item.
// If the duration is 0, the cache's default expiration time is used.
// If it < 0, the item never expires.
func (c *cache[K, V]) SetWithTTL(k K, v V, d time.Duration) {
	t := c.timeToLive(d)

	c.mu.Lock()
	c.set(k, v, t)
	c.mu.Unlock()
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns false otherwise.
func (c *cache[K, V]) Add(k K, v V) bool {
	return c.AddWithTTL(k, v, c.ttl)
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns false otherwise.
func (c *cache[K, V]) AddWithTTL(k K, v V, d time.Duration) bool {
	c.mu.Lock()
	_, ok := c.get(k)
	if !ok {
		t := c.timeToLive(d)
		c.set(k, v, t)
	}

	c.mu.Unlock()
	return !ok
}

// Replace a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns false otherwise.
func (c *cache[K, V]) Replace(k K, v V) bool {
	return c.ReplaceWithTTL(k, v, c.ttl)
}

// Replace a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns false otherwise.
func (c *cache[K, V]) ReplaceWithTTL(k K, v V, d time.Duration) bool {
	c.mu.Lock()
	_, ok := c.get(k)
	if ok {
		t := c.timeToLive(d)
		c.set(k, v, t)
	}

	c.mu.Unlock()
	return ok
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *cache[K, V]) Get(k K) (V, bool) {
	c.mu.RLock()

	item, found := c.items[k]
	if found && !item.Expired() {
		c.mu.RUnlock()
		return item.Val, true
	}

	c.mu.RUnlock()

	var v V
	return v, false
}

// GetWithTTL returns an item and its expiration time from the cache.
// It returns the item or nil, the expiration time if one is set (if the item
// never expires a zero value for time.Time is returned), and a bool indicating
// whether the key was found.
func (c *cache[K, V]) GetWithTTL(k K) (V, time.Time, bool) {
	c.mu.RLock()

	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		var v V
		return v, time.Time{}, false
	}

	if item.TTL <= 0 {
		// If expiration <= 0 (i.e. no expiration time set) then return the item
		// and a zeroed time.Time
		c.mu.RUnlock()
		return item.Val, time.Time{}, true
	}

	if item.Expired() {
		c.mu.RUnlock()
		var v V
		return v, time.Time{}, false
	}

	// Return the item and the expiration time
	c.mu.RUnlock()
	return item.Val, time.Unix(item.TTL, 0), true
}

// Remove an item from the cache. Does nothing if the key is not in the cache.
func (c *cache[K, V]) Remove(k K) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}

// Clean Remove all expired items from the cache.
func (c *cache[K, V]) Clean() {
	c.mu.Lock()
	now := time.Now().Unix()
	for k, v := range c.items {
		if v.ExpiredAt(now) {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Clear Remove all items from the cache.
func (c *cache[K, V]) Clear() {
	c.mu.Lock()
	c.items = map[K]Item[V]{}
	c.mu.Unlock()
}

// Items Copies all unexpired items in the cache into a new map and returns it.
func (c *cache[K, V]) Items() map[K]Item[V] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now().Unix()

	m := make(map[K]Item[V], len(c.items))
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
func (c *cache[K, V]) Each(f func(K, Item[V]) bool) {
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

// Increment an item of type int, int8, int16, int32, int64, uint,
// uint8, uint32, or uint64, float32 or float64 by n.
// Returns the incremented value.
func (c *cache[K, V]) Increment(k K, n V) V {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[k]
	if !found || item.Expired() {
		c.set(k, n, c.timeToLive(c.ttl))
		return n
	}

	v, err := mathx.Add(item.Val, n)
	if err != nil {
		panic(fmt.Errorf("item '%v' cannot increment: %w", k, err))
	}

	item.Val = v.(V)
	c.items[k] = item
	return item.Val
}
