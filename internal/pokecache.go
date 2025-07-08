package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		entries: make(map[string]cacheEntry),
		mu:      &sync.Mutex{},
	}

	cache.reapLoop(interval)

	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	newEntry := cacheEntry{val: val, createdAt: time.Now()}
	c.entries[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	value, exists := c.entries[key]

	if !exists {
		return nil, false
	}

	return value.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	timer := time.NewTicker(interval)

	go func() {
		for range timer.C {
			c.mu.Lock()
			for key, value := range c.entries {
				if time.Since(value.createdAt) > interval {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
