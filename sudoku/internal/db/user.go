package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

// Key to user data in the database.
func userKey(username string) string {
	return fmt.Sprintf("user\u0000%s", strings.ToLower(username))
}

// CreateUser creates an entity in the database with user data.
func CreateUser(ctx context.Context, r redis.Cmdable, user data.User) error {
	user.SignupTimestamp = time.Now().UTC().Unix()
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.Set(ctx, userKey(user.Username), body, 0).Err()
}

// GetUser returns user data.
func GetUser(ctx context.Context, r redis.Cmdable, username string) (data.User, bool, error) {
	body, err := r.Get(ctx, userKey(username)).Result()
	if err != nil {
		if err == redis.Nil {
			return data.User{}, false, nil
		}
		return data.User{}, false, err
	}
	var u data.User
	if err := json.Unmarshal([]byte(body), &u); err != nil {
		return data.User{}, false, err
	}
	return u, true, nil
}

// IsExistsUser checks if the user exists.
func IsExistsUser(ctx context.Context, r redis.Cmdable, username string) (bool, error) {
	n, err := r.Exists(ctx, userKey(username)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
