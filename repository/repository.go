package repository

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisRepository struct {
	pool *redis.Pool
}

func New(pool *redis.Pool) *RedisRepository {
	return &RedisRepository{
		pool: pool,
	}
}

func redisExpiration(dur time.Duration) int {
	return int(dur.Truncate(time.Second).Seconds())
}
