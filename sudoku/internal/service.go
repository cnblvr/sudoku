package sudoku

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

// NewService initialize the service sudoku.
func NewService() (*Service, error) {
	srv := &Service{}

	// Connect to redis database
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}

	return srv, nil
}

// Service is a service struct.
type Service struct {
}
