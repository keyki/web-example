package cache

import (
	"sync"
	"time"
	"web-example/log"
)

type Item struct {
	Value      any
	Expiration time.Time
}

type Cache struct {
	data            map[string]Item
	mutex           sync.RWMutex
	cleanupInterval time.Duration
}

func NewCache(cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		data:            make(map[string]Item),
		cleanupInterval: cleanupInterval,
	}

	go cache.cleanupExpiredItems()

	return cache
}

func (c *Cache) Set(key string, item any, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = Item{
		Value:      item,
		Expiration: time.Now().Add(duration),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, exists := c.data[key]
	if !exists || time.Now().After(item.Expiration) {
		return nil, false
	}
	log.BaseLogger().Infof("Cache hit for key: %v", key)
	return item.Value, true
}

func (c *Cache) cleanupExpiredItems() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			now := time.Now()
			for key, item := range c.data {
				if now.After(item.Expiration) {
					log.BaseLogger().Infof("Deleting expired user from cache: %s", key)
					delete(c.data, key)
				}
			}
			c.mutex.Unlock()
		}
	}
}
