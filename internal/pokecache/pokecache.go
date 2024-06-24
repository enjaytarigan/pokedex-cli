package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	expireAt time.Time
	val      []byte
}

func (c cacheEntry) isExpired() bool {
	return time.Now().After(c.expireAt)
}

type cache struct {
	items      map[string]cacheEntry
	mu         sync.Mutex
	expiryTime time.Duration
}

func NewCache(expiryTime time.Duration) *cache {
	c := &cache{
		items: make(map[string]cacheEntry),
	}

	go func() {
		for range time.Tick(expiryTime) {
			c.mu.Lock()

			for key, item := range c.items {
				if item.isExpired() {
					delete(c.items, key)
				}
			}

			c.mu.Unlock()
		}
	}()

	return c
}

func (c *cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheEntry{
		val:      value,
		expireAt: time.Now().Add(c.expiryTime),
	}
}

func (c *cache) Get(key string) (value []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.items[key]

	if !ok {
		// Default value
		return v.val, false
	}

	if v.isExpired() {
		return v.val, false
	}

	value = v.val
	found = true

	return value, found
}
