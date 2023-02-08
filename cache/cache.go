package cache

import (
	"fmt"
	"log"
	"sync"

	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache{
	return &Cache {
		data: make(map[string][]byte),
	}
}

func (c *Cache) Delete(key []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	k := string(key)

	delete(c.data, k)
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(key)]

	return ok
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	k := string(key)

	val, ok := c.data[k]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", k)
	}

	log.Printf("GET %s = %s\n", k, string(val))

	return val, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	go func(){
		<-time.After(ttl)
		c.Delete(key)
	}()

	k := string(key)
	_, ok := c.data[k]
	if ok {
		return fmt.Errorf("key (%s) already exists", k)
	}

	c.data[k] = value
	log.Printf("SET %s to %s\n", k, string(value))

	return nil
}
