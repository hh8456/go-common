package cache

import (
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New()

	go func() {
		for m := 0; m < 10; m++ {
			go func() {
				for i := 0; i < 10000; i++ {
					str := strconv.Itoa(i)

					c.Set(str, str, time.Minute)
				}
			}()
		}
	}()

	go func() {
		for m := 0; m < 10; m++ {
			go func() {
				for i := 0; i < 10000; i++ {
					str := strconv.Itoa(i)
					c.Get(str)
				}
			}()
		}

	}()
}
