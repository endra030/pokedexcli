package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu   sync.Mutex
	e    map[string]cacheEntry
	d    time.Duration
	quit chan struct{}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		d: interval,
		e: map[string]cacheEntry{},
	}
	go c.reapLoop()
	return &c

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.e[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, b bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var retVal []byte
	ce, ok := c.e[key]
	if ok {
		b = true
		retVal = ce.val
	}
	return retVal, b

}

func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.d)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			expiredKeys := []string{}
			for k, v := range c.e {
				if time.Since(v.createdAt) > c.d {
					expiredKeys = append(expiredKeys, k)
				}
			}
			for _, ek := range expiredKeys {
				delete(c.e, ek)
			}
			c.mu.Unlock()
		case <-c.quit:
			return

		}

	}

}

func (c *Cache) Stop() {
	close(c.quit)
}
