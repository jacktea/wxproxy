package cache

import (
	//"github.com/jacktea/golang-lru/simplelru"
	"github.com/bluele/gcache"
	"time"
)

type Cache interface {
	Get(interface{}) (interface{}, bool)
	Set(interface{}, interface{}) bool
	Add(interface{}, interface{}) bool
}

type zCache struct {
	cache gcache.Cache
}

func NewCache(size int, expiration time.Duration) *zCache {
	return &zCache{
		cache: gcache.New(size).Expiration(expiration).LRU().Build(),
	}
}

func (c *zCache) Get(key interface{}) (interface{}, bool) {
	v, err := c.cache.Get(key)
	if err != nil {
		return nil, false
	} else {
		return v, true
	}
}

func (c *zCache) Set(key, value interface{}) bool {
	return c.cache.Set(key, value) == nil
}

func (c *zCache) Add(key, value interface{}) bool {
	return c.cache.Set(key, value) == nil
}

/*
func Purge() {
	cache.Purge()
}

func Add(key, value interface{}) bool {
	return cache.Add(key,value)
}

func Get(key interface{}) (value interface{}, ok bool) {
	return cache.Get(key)
}

func Contains(key interface{}) (ok bool) {
	return cache.Contains(key)
}

func Peek(key interface{}) (value interface{}, ok bool) {
	return cache.Peek(key)
}

func Remove(key interface{}) bool {
	return cache.Remove(key)
}

func RemoveOldest() (interface{}, interface{}, bool) {
	return cache.RemoveOldest()
}

func GetOldest() (interface{}, interface{}, bool) {
	return cache.GetOldest()
}

func Keys() []interface{} {
	return cache.Keys()
}

func Len() int {
	return cache.Len()
}
*/
