package cache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

//func TestCacheParallel(t *testing.T) {
//t.Parallel()
//tc := NewInt64()
//tc.Set(1, "bar", time.Minute)
//tc.Delete(1)
//}

func TestMultiCache(t *testing.T) {
	t.Parallel()
	mm := map[int64]*CacheInt64{}
	for i := 0; i < 10; i++ {
		mm[int64(i)] = NewInt64()
	}

	var wg sync.WaitGroup
	wg.Add(100)
	cnt := int64(0)

	for i := 0; i < 100; i++ {
		go func(i int) {
			for j := 0; j < 100000; j++ {
				k := atomic.AddInt64(&cnt, 1)
				mm[int64(i%10)].Set(k, "bar", time.Minute)
				mm[int64(i%10)].Delete(k)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkMultiCacheInt64WithLock(b *testing.B) {
	b.StopTimer()
	mm := map[int64]*CacheInt64{}
	for i := 0; i < 10; i++ {
		mm[int64(i)] = NewInt64()
	}

	var locks [10]sync.RWMutex
	cnt := int64(0)
	n := 10
	var wg sync.WaitGroup
	wg.Add(n)

	b.SetParallelism(2)

	b.StartTimer()
	for i := 0; i < n; i++ {
		go func(i int) {
			for j := 0; j < b.N; j++ {
				k := atomic.AddInt64(&cnt, 1)
				lock := &locks[i%10]
				lock.Lock()
				mm[int64(i%10)].Set(k, "bar", time.Minute)
				mm[int64(i%10)].Delete(k)
				lock.Unlock()
			}

			wg.Done()
		}(i)

	}

	wg.Wait()
}

func BenchmarkMultiCacheInt64(b *testing.B) {
	b.StopTimer()
	mm := map[int64]*CacheInt64{}
	for i := 0; i < 10; i++ {
		mm[int64(i)] = NewInt64()
	}

	n := 10
	var wg sync.WaitGroup
	wg.Add(n)

	cnt := int64(0)

	b.StartTimer()
	for i := 0; i < n; i++ {
		go func(i int) {
			for j := 0; j < b.N; j++ {
				k := atomic.AddInt64(&cnt, 1)
				mm[int64(i%10)].Set(k, "bar", time.Minute)
				mm[int64(i%10)].Delete(k)
			}

			wg.Done()
		}(i)

	}

	wg.Wait()
}

func BenchmarkCacheInt64SetDelete(b *testing.B) {
	b.StopTimer()
	tc := NewInt64()

	n := 10
	var wg sync.WaitGroup
	wg.Add(n)
	cnt := int64(0)
	b.SetParallelism(2)

	b.StartTimer()
	for i := 0; i < n; i++ {
		go func(i int) {
			for j := 0; j < b.N; j++ {
				k := atomic.AddInt64(&cnt, 1)
				tc.Set(k, "bar", time.Minute)
				tc.Delete(int64(j))
			}

			wg.Done()
		}(i)

	}

	wg.Wait()

}
