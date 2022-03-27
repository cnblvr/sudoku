package generator

import (
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/library_puzzles"
	"github.com/cnblvr/sudoku/repository"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"time"
)

// Service is a service structure.
type Service struct {
	sudokuRepository   data.SudokuRepository
	generateRepository data.SudokuGenerateRepository
	generator          data.SudokuGenerator
}

func NewService() (*Service, error) {
	srv := &Service{}

	rand.Seed(time.Now().UnixNano())

	// Connect to redis database
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp", "redis:6379", // todo port from env vars
				redis.DialPassword(os.Getenv("REDIS_PASSWORD")),
				redis.DialDatabase(0), // todo index of database from env vars
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
	conn := redisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		log.Error().Err(err).Msg("failed to ping to redis database")
		return nil, err
	}
	repo := repository.New(redisPool)
	srv.sudokuRepository = repo
	srv.generateRepository = repo

	var err error
	srv.generator, err = library_puzzles.GetGenerator(data.SudokuClassic)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
