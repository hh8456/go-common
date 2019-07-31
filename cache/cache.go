package cache

import (
	"runtime"
	"sync"
	"time"
)

type item struct {
	obj        interface{}
	expiration int64
}

type Cache struct {
	*cache
}

type cache struct {
	items map[int64]item
	lock  sync.RWMutex
	stop  chan bool
}

func New() *Cache {
	c := &cache{
		items: make(map[int64]item),
		stop:  make(chan bool),
	}

	C := &Cache{c}
	runtime.SetFinalizer(C, stopCheckExpired)

	go checkExpired(c)

	return &Cache{c}
}

func stopCheckExpired(c *Cache) {
	c.stop <- true
}

func (c *cache) Set(k int64, x interface{}, d time.Duration) {
	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	c.lock.Lock()
	c.items[k] = item{
		obj:        x,
		expiration: e,
	}
	c.lock.Unlock()
}

func (c *cache) Get(k int64) (interface{}, bool) {
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

func (c *cache) Delete(k int64) {
	c.lock.Lock()
	delete(c.items, k)
	c.lock.Unlock()
}

func (c *cache) Flush() {
	c.lock.Lock()
	c.items = map[int64]item{}
	c.lock.Unlock()
}

func (c *cache) deleteExpired() {
	now := time.Now().UnixNano()
	c.lock.Lock()
	for k, v := range c.items {
		if v.expiration > 0 && now > v.expiration {
			delete(c.items, k)
		}
	}
	c.lock.Unlock()
}

func checkExpired(c *cache) {
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
