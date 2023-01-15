package cache_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/Pauloo27/aryzona/internal/core/cache"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	resolveCalls := 0
	cache := cache.NewCache(func(key string) (*int, time.Duration) {
		resolveCalls++
		v, err := strconv.Atoi(key)
		// i should not be allowed to write tests that depends on time...
		ttl := time.Microsecond * 100
		if err == nil {
			return &v, ttl
		}
		return nil, ttl
	})

	t.Run("should get the right values", func(t *testing.T) {
		assert.Equal(t, 10, *cache.GetRawValue("10"))
		assert.Equal(t, 12, *cache.GetRawValue("12"))
		assert.Equal(t, -10, *cache.GetRawValue("-10"))
		assert.Equal(t, 0, *cache.GetRawValue("0"))
		assert.Equal(t, 0, *cache.GetRawValue("0000"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("-10.2"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("batata"))

		cache.Clear()
		resolveCalls = 0
	})

	t.Run("should clear the cache", func(t *testing.T) {
		assert.Equal(t, 10, *cache.GetRawValue("10"))
		assert.Equal(t, 12, *cache.GetRawValue("12"))
		assert.Equal(t, -10, *cache.GetRawValue("-10"))
		assert.Equal(t, 0, *cache.GetRawValue("0"))
		assert.Equal(t, 0, *cache.GetRawValue("0000"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("-10.2"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("batata"))

		assert.Equal(t, 7, cache.Size())
		cache.Clear()
		assert.Equal(t, 0, cache.Size())

		resolveCalls = 0
	})

	t.Run("should not call the resolve twice for the same key", func(t *testing.T) {
		assert.Equal(t, 10, *cache.GetRawValue("10"))
		assert.Equal(t, 10, *cache.GetRawValue("10"))

		assert.Equal(t, 12, *cache.GetRawValue("12"))
		assert.Equal(t, 12, *cache.GetRawValue("12"))

		assert.Equal(t, -10, *cache.GetRawValue("-10"))
		assert.Equal(t, -10, *cache.GetRawValue("-10"))

		assert.Equal(t, 0, *cache.GetRawValue("0"))
		assert.Equal(t, 0, *cache.GetRawValue("0"))

		assert.Equal(t, 0, *cache.GetRawValue("0000"))
		assert.Equal(t, 0, *cache.GetRawValue("0000"))

		assert.Equal(t, (*int)(nil), cache.GetRawValue("-10.2"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("-10.2"))

		assert.Equal(t, (*int)(nil), cache.GetRawValue("batata"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("batata"))

		assert.Equal(t, 7, resolveCalls)
		assert.Equal(t, 7, cache.Size())

		cache.Clear()
		resolveCalls = 0
	})

	t.Run("should call the resolve again after TTL is expired", func(t *testing.T) {
		assert.Equal(t, 10, *cache.GetRawValue("10"))
		assert.Equal(t, 12, *cache.GetRawValue("12"))
		assert.Equal(t, -10, *cache.GetRawValue("-10"))
		assert.Equal(t, 0, *cache.GetRawValue("0"))
		assert.Equal(t, 0, *cache.GetRawValue("0000"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("-10.2"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("batata"))

		assert.Equal(t, 7, cache.Size())
		time.Sleep(200 * time.Millisecond)

		assert.Equal(t, 10, *cache.GetRawValue("10"))
		assert.Equal(t, 12, *cache.GetRawValue("12"))
		assert.Equal(t, -10, *cache.GetRawValue("-10"))
		assert.Equal(t, 0, *cache.GetRawValue("0"))
		assert.Equal(t, 0, *cache.GetRawValue("0000"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("-10.2"))
		assert.Equal(t, (*int)(nil), cache.GetRawValue("batata"))

		assert.Equal(t, 2*7, resolveCalls)

		cache.Clear()
		resolveCalls = 0
	})
}
