package cache

import (
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewInt64()

	cnt := int64(0)
	go func() {
		for m := 0; m < 10; m++ {
			go func() {
				for i := 0; i < 10000; i++ {
					str := strconv.Itoa(i)
					k := atomic.AddInt64(&cnt, 1)
					c.Set(k, str, time.Minute)
				}
			}()
		}
	}()

}

func BenchmarkMapInt64SetDelete(b *testing.B) {
	b.StopTimer()
	m := map[int64]string{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[int64(i)] = "bar"
		delete(m, int64(i))
	}
}

func BenchmarkCacheInt64SetDelete(b *testing.B) {
	b.StopTimer()
	tc := NewInt64()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Set(int64(i), "bar", time.Minute)
		tc.Delete(int64(i))
	}
}
