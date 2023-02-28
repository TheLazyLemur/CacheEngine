package cache

import (
	"fmt"
	"sync"

	"time"
)

// TODO: Add database to back up cache
// TODO: Restore cache from database on startup
// TODO: Add a cache warmup
// TODO: Sync cache with database when shutdown

type CacheImpl struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *CacheImpl {
	return &CacheImpl{
		data: make(map[string][]byte),
	}
}

func (c *CacheImpl) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	k := string(key)

	delete(c.data, k)
	// TODO: Remove from database
	return nil
}

func (c *CacheImpl) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(key)]

	return ok
}

func (c *CacheImpl) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	k := string(key)

	val, ok := c.data[k]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", k)
	}

	return val, nil
}

func (c *CacheImpl) Set(key, value []byte, ttl int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ttl > 0 {
		go func() {
			//TODO: Convert this to a go routine that runs separately and will clean up the cache after a specified time
			<-time.After(time.Duration(ttl))
			_ = c.Delete(key)
		}()
	}

	k := string(key)
	_, ok := c.data[k]
	if ok {
		return fmt.Errorf("key (%s) already exists", k)
	}

	c.data[k] = value
	//TODO: Store in database

	return nil
}

func (c *CacheImpl) All() ([][]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	ks := make([][]byte, 0)

	for k := range c.data {
		ks = append(ks, []byte(k))
	}

	return ks, nil
}
