package repository_test

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"testing"
	"time"
)

const (
	redisTestAddress  = "localhost:6379"
	redisTestPassword = ``
	redisTestDatabase = 9
)

func newRedisTestPool(t *testing.T) *redis.Pool {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp", redisTestAddress,
				redis.DialPassword(redisTestPassword),
				redis.DialDatabase(redisTestDatabase),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Second*10 {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: time.Minute,
	}
	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		t.Fatal(errors.Wrap(err, "failed to ping redis test"))
	}
	return pool
}

func flushDB(t *testing.T, pool *redis.Pool) *redis.Pool {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("FLUSHALL"); err != nil {
		t.Fatal(err)
	}
	return pool
}
