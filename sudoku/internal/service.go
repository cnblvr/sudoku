package sudoku

import (
	"encoding/base64"
	"fmt"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"os"
)

// Service is a service structure.
type Service struct {
	// Connection of database
	redis redis.Conn // TODO: stop redis and test
	// Storage for templates
	templates *template.Template
	// Object to generate hash from password.
	// Object settings are stored in environment variables: SECURECOOKIE_HASH_KEY and SECURECOOKIE_BLOCK_KEY
	securecookie *securecookie.SecureCookie
	// Pepper for password hashing
	passwordPepper []byte
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
	conn, err := redis.Dial(
		"tcp", "redis:6379", // todo port from env vars
		redis.DialPassword(os.Getenv("REDIS_PASSWORD")),
		redis.DialDatabase(0), // todo index of database from env vars
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to redis")
	}
	srv.redis = conn
	if pong, err := srv.redis.Do("PING"); err != nil {
		log.Error().Err(err).Msg("failed to ping to redis database")
		return nil, err
	} else if pong != "PONG" {
		log.Error().Msg("ping, not pong")
		return nil, fmt.Errorf("redis ping error")
	}

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
