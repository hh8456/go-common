package redisObj

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hh8456/redisSession"
)

var (
	pool *redis.Pool
)

func Init(redisAddr string, redisPwd string) {
	if pool != nil {
		return
	}

	var e error
	pool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		Wait:        true,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				e = err
				return nil, err
			}

			if redisPwd != "" {
				if _, err := c.Do("AUTH", redisPwd); err != nil {
					c.Close()
					e = err
					return nil, err
				}
			}
			return c, nil
		},
	}

	if e != nil {
		str := fmt.Sprintf("create redis pool has error: %v", e)
		panic(str)
	}
}

func NewSessionWithPrefix(prefix string) *redisSession.RedisSession {
	rdsSess := redisSession.NewRedisSessionWithPool(pool)
	rdsSess.SetPrefix(prefix)
	return rdsSess
}
