package cache

import (
	"runtime"
	"sync"
	"time"
)

type itemInt64 struct {
	obj        interface{}
	expiration int64
}

type CacheInt64 struct {
	*cacheInt64
}

type cacheInt64 struct {
	items map[int64]itemInt64
	lock  sync.RWMutex
	stop  chan bool
}

func NewInt64() *CacheInt64 {
	c := &cacheInt64{
		items: make(map[int64]itemInt64),
		stop:  make(chan bool),
	}

	C := &CacheInt64{c}
	runtime.SetFinalizer(C, stopCheckExpiredInt64)

	go checkExpiredInt64(c)

	return &CacheInt64{c}
}

func stopCheckExpiredInt64(c *CacheInt64) {
	c.stop <- true
}

func (c *cacheInt64) Set(k int64, x interface{}, d time.Duration) {
	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	c.lock.Lock()
	c.items[k] = itemInt64{
		obj:        x,
		expiration: e,
	}
	c.lock.Unlock()
}

func (c *cacheInt64) Get(k int64) (interface{}, bool) {
	c.lock.RLock()
	item, found := c.items[k]
	if !found {
		c.lock.RUnlock()
		return nil, false
	}

	if item.expiration > 0 {
		if time.Now().UnixNano() > item.expiration {
			c.lock.RUnlock()
			return nil, false
		}
	}

	c.lock.RUnlock()
	return item.obj, true
}

func (c *cacheInt64) Delete(k int64) {
	c.lock.Lock()
	delete(c.items, k)
	c.lock.Unlock()
}

func (c *cacheInt64) Flush() {
	c.lock.Lock()
	c.items = map[int64]itemInt64{}
	c.lock.Unlock()
}

func (c *cacheInt64) deleteExpired() {
	now := time.Now().UnixNano()
	c.lock.Lock()
	for k, v := range c.items {
		if v.expiration > 0 && now > v.expiration {
			delete(c.items, k)
		}
	}
	c.lock.Unlock()
}

func checkExpiredInt64(c *cacheInt64) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()

		case <-c.stop:
			ticker.Stop()
			return
		}
	}

}
