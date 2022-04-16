package cache

import "time"

type CacheEntry[T any] struct {
	Value    T
	CachedAt time.Time
	TTL      time.Duration
}

type Cache[T any] struct {
	store   map[string]CacheEntry[T]
	resolve func(string) (T, time.Duration)
}

func NewCache[T any](resolve func(string) (T, time.Duration)) *Cache[T] {
	return &Cache[T]{
		store:   make(map[string]CacheEntry[T]),
		resolve: resolve,
	}
}

func (c *Cache[T]) Get(key string) (entry CacheEntry[T], cached bool) {
	entry, found := c.store[key]
	if found {
		if time.Since(entry.CachedAt) > entry.TTL {
			return entry, true
		}
	}
	value, ttl := c.resolve(key)
	entry = CacheEntry[T]{
		TTL:      ttl,
		Value:    value,
		CachedAt: time.Now(),
	}
	c.store[key] = entry
	return entry, false
}

func (c *Cache[T]) GetValue(key string) (value T, cached bool) {
	entry, cached := c.Get(key)
	return entry.Value, cached
}

func (c *Cache[T]) GetRawValue(key string) (value T) {
	entry, _ := c.Get(key)
	return entry.Value
}

func (c *Cache[T]) Remove(key string) {
	delete(c.store, key)
}

func (c *Cache[T]) Size() int {
	return len(c.store)
}

func (c *Cache[T]) Clear() {
	// how does the GC deal with this? i hope it does.
	// anyway, at one moment the old entries will be released from the memory,
	// either by the GB or by lack of power for the RAM.
	c.store = make(map[string]CacheEntry[T])
}
