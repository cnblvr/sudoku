package sudoku

import (
	"encoding/base64"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/repository"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"os"
	"time"
)

// Service is a service structure.
type Service struct {
	userRepository   data.UserRepository
	sudokuRepository data.SudokuRepository
	// Storage for templates
	templates *template.Template
	// Object to generate hash from password.
	// Object settings are stored in environment variables: SECURECOOKIE_HASH_KEY and SECURECOOKIE_BLOCK_KEY
	securecookie *securecookie.SecureCookie
	// Pepper for password hashing
	passwordPepper []byte
	// gorilla/websocket object
	upgrader websocket.Upgrader
}

// NewService initialize the service sudoku.
func NewService() (*Service, error) {
	srv := &Service{}

	// Load securecookie environment variables
	hashKey, err := base64.StdEncoding.DecodeString(os.Getenv("SECURECOOKIE_HASH_KEY"))
	if err != nil {
		log.Error().Err(err).Msg("failed to decode 'SECURECOOKIE_HASH_KEY' env variable")
		return nil, err
	}
	blockKey, err := base64.StdEncoding.DecodeString(os.Getenv("SECURECOOKIE_BLOCK_KEY"))
	if err != nil {
		log.Error().Err(err).Msg("failed to decode 'SECURECOOKIE_BLOCK_KEY' env variable")
		return nil, err
	}
	srv.securecookie = securecookie.New(hashKey, blockKey)
	if _, err := srv.securecookie.Encode("test", "test"); err != nil {
		log.Error().Err(err).Msg("failed to test encode securecookie")
		return nil, err
	}

	// Load pepper for hash passwords from environment variable
	srv.passwordPepper, err = base64.StdEncoding.DecodeString(os.Getenv("PASSWORD_PEPPER"))
	if err != nil {
		log.Error().Err(err).Msg("failed to decode 'PASSWORD_PEPPER' env variable")
		return nil, err
	}

	// Load templates
	srv.templates, err = template.ParseFS(templates.Templates, append(templates.Common(), "*.gohtml")...)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse FS of templates")
		return nil, err
	}

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
	srv.userRepository = repo
	srv.sudokuRepository = repo

	// init upgrader
	srv.upgrader = websocket.Upgrader{}

	return srv, nil
}

func (srv *Service) executeTemplate(w http.ResponseWriter, name string, args templates.Args) {
	err := srv.templates.ExecuteTemplate(w, name, args)
	if err != nil {
		log.Error().Err(err).Str("template", name).Msg("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
