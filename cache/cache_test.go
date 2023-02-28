package cache

import (
	"testing"
	"time"
)

func getNewCache(t *testing.T) *CacheImpl {
	c := New()
	if c == nil {
		t.Errorf("new cache returned nil")
	}

	return c
}

func getTestValues() ([]byte, []byte, int64) {
	key := []byte("test_key")
	value := []byte("test_val")
	ttl := 1000

	return key, value, int64(ttl)
}

func TestNewCache(t *testing.T) {
	getNewCache(t)
}

func TestSetValueInCache(t *testing.T) {
	c := getNewCache(t)

	key, value, ttl := getTestValues()

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Errorf("error setting cache")
	}
}

func TestFailSetDuplicateValueInCache(t *testing.T) {
	c := getNewCache(t)

	key, value, ttl := getTestValues()

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Errorf("error setting cache")
	}

	err = c.Set(key, value, ttl)
	if err == nil {
		t.Errorf("cache didnt fail on duplicate value")
	}
}

func TestGetValueFromCache(t *testing.T) {
	c := getNewCache(t)

	key, value, ttl := getTestValues()

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Errorf("error setting cache")
	}

	returnedVal, err := c.Get(key)
	if err != nil {
		t.Errorf("could not retrieve key")
	}

	if string(returnedVal) != string(value) {
		t.Errorf("retrieved key was incorrect")
	}
}

func TestHasKey(t *testing.T) {
	c := getNewCache(t)

	key, value, ttl := getTestValues()

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Errorf("error setting cache")
	}

	has := c.Has(key)
	if !has {
		t.Errorf("key was not in cache as expected")
	}
}

func TestDeleteFromCache(t *testing.T) {
	c := getNewCache(t)

	key, value, ttl := getTestValues()

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Errorf("error setting cache")
	}

	has := c.Has(key)
	if !has {
		t.Errorf("key was not in cache as expected")
	}

	_ = c.Delete(key)
	hasAfterDelete := c.Has(key)
	if hasAfterDelete {
		t.Errorf("key was not deleted from cache")
	}
}

func TestValueExpiresAfterTime(t *testing.T) {
	c := getNewCache(t)

	key, value, ttl := getTestValues()

	err := c.Set(key, value, ttl)
	if err != nil {
		t.Errorf("error setting cache")
	}

	has := c.Has(key)
	if !has {
		t.Errorf("key was not in cache as expected")
	}

	time.Sleep(2 * time.Second)

	hasAfterExpire := c.Has(key)
	if hasAfterExpire {
		t.Errorf("key did not expire from cache")
	}
}
