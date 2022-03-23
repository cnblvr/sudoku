package repository

import (
	"github.com/gomodule/redigo/redis"
)

type RedisRepository struct {
	pool *redis.Pool
}

func New(pool *redis.Pool) *RedisRepository {
	return &RedisRepository{
		pool: pool,
	}
}
