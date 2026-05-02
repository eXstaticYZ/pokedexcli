package pokecache

import (
	"sync"
	"time"
)

// cacheEntry represents a cached entry with its creation time and value
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache provides thread-safe caching functionality
type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mutex    sync.Mutex
	ticker   *time.Ticker
	done     chan bool
}

// NewCache creates a new cache with the specified cleanup interval
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
		done:     make(chan bool, 1),
	}

	// Start the background reaping loop
	cache.ticker = time.NewTicker(interval)
	go cache.reapLoop()

	return cache
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       append([]byte{}, val...), // Create a copy to avoid external modifications
	}
}

// Get retrieves an entry from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, found := c.entries[key]
	if !found {
		return nil, false
	}

	return append([]byte{}, entry.val...), true // Return a copy to avoid external modifications
}

// reapLoop runs in the background and periodically removes expired entries
func (c *Cache) reapLoop() {
	for {
		select {
		case <-c.ticker.C:
			c.reap()
		case <-c.done:
			c.ticker.Stop()
			return
		}
	}
}

// reap removes all entries older than the cache interval
func (c *Cache) reap() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.entries, key)
		}
	}
}

// Close stops the background reaping loop
func (c *Cache) Close() error {
	c.done <- true
	return nil
}
