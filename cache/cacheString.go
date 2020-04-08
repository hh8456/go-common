package cache

import (
	"runtime"
	"sync"
	"time"
)

type itemString struct {
	obj        interface{}
	expiration int64
}

type CacheString struct {
	*cacheStr
}

type cacheStr struct {
	items map[string]itemString
	lock  sync.RWMutex
	stop  chan bool
}

func New() *CacheString {
	c := &cacheStr{
		items: make(map[string]itemString),
		stop:  make(chan bool),
	}

	C := &CacheString{c}
	runtime.SetFinalizer(C, stopCheckExpiredString)

	go checkExpired(c)

	return &CacheString{c}
}

func stopCheckExpiredString(c *CacheString) {
	c.stop <- true
}

func (c *cacheStr) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	c.lock.Lock()
	c.items[k] = itemString{
		obj:        x,
		expiration: e,
	}
	c.lock.Unlock()
}

func (c *cacheStr) Get(k string) (interface{}, bool) {
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

func (c *cacheStr) Delete(k string) {
	c.lock.Lock()
	delete(c.items, k)
	c.lock.Unlock()
}

func (c *cacheStr) Flush() {
	c.lock.Lock()
	c.items = map[string]itemString{}
	c.lock.Unlock()
}

func (c *cacheStr) deleteExpired() {
	now := time.Now().UnixNano()
	c.lock.Lock()
	for k, v := range c.items {
		if v.expiration > 0 && now > v.expiration {
			delete(c.items, k)
		}
	}
	c.lock.Unlock()
}

func checkExpired(c *cacheStr) {
	defer println("checkExpired over")
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
