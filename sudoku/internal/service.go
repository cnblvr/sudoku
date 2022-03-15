package sudoku

import (
	"context"
	"encoding/base64"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/securecookie"
	"github.com/rs/zerolog/log"
	"html/template"
	"os"
)

// Service is a service structure.
type Service struct {
	// Connection of database
	redis *redis.Client
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
	srv.redis = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // todo port from env vars
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // todo index of database from env vars
	})
	if _, err := srv.redis.Ping(context.TODO()).Result(); err != nil {
		log.Error().Err(err).Msg("failed to ping to redis database")
		return nil, err
	}

	return srv, nil
}
