package logRedisWriter

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hh8456/redisSession"
)

var (
	pid  = os.Getpid()
	host = "unknowhost"
)

func init() {
	h, e := os.Hostname()
	if e == nil {
		host = shortHostname(h)
	}
}

func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}

	return hostname
}

type logRedisWriter struct {
	pid      int
	host     string
	rdsSess  *redisSession.RedisSession
	chanData chan []byte
	done     chan interface{}
}

func NewLogFileWriter(rdsSess *redisSession.RedisSession) *logRedisWriter {
	lfw := &logRedisWriter{rdsSess: rdsSess, host: host}
	lfw.chanData = make(chan []byte, 10000)
	lfw.done = make(chan interface{}, 1)

	go lfw.saveToRedis()

	return lfw
}

func (w *logRedisWriter) Write(p []byte) (int, error) {
	select {
	case w.chanData <- p:

	default:
		return 0, errors.New("logRedisWriter.chanData is full")
	}

	return len(p), nil
}

func (w *logRedisWriter) Close() {
	close(w.done)
}

func (w *logRedisWriter) saveToRedis() {

	redisKey := "central-log"

	if w.rdsSess.GetPrefix() != "" {
		redisKey = w.rdsSess.GetPrefix()
	}

	localData := make([][]byte, 10)
	for {
		cnt := 0
		select {
		case <-w.done:
			return

		case data1, ok1 := <-w.chanData:
			if !ok1 {
				return
			}
			localData[cnt] = data1
			cnt++

			exitLoop := false
			for {
				select {
				case data2, ok2 := <-w.chanData:
					if !ok2 {
						return
					}

					localData[cnt] = data2
					cnt++

				default:
					exitLoop = true
				}

				if exitLoop || cnt == 10 {
					break
				}
			}
		}

		var err error
		// 写 redis
		switch cnt {
		case 1:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0])

		case 2:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1])

		case 3:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2])

		case 4:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2], localData[3])

		case 5:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2], localData[3], localData[4])

		case 6:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2], localData[3], localData[4], localData[5])

		case 7:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2], localData[3], localData[4], localData[5],
				localData[6])

		case 8:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2], localData[3], localData[4], localData[5],
				localData[6], localData[7])

		case 9:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2], localData[3], localData[4], localData[5],
				localData[6], localData[7], localData[8])

		case 10:
			_, err = w.rdsSess.Do("lpush", redisKey, localData[0], localData[1],
				localData[2], localData[3], localData[4], localData[5],
				localData[6], localData[7], localData[8], localData[9])

		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to access redis , %v\n", err)
		}
	}
}
