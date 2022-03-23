package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"strings"
)

func (r *RedisRepository) CreateUser(ctx context.Context, username string) (*data.User, error) {
	conn := r.pool.Get()
	defer conn.Close()

	if ok, err := redis.Bool(conn.Do("EXISTS", keyUsername(username))); err != nil {
		return nil, errors.Wrap(err, "failed to get owner of username")
	} else if ok {
		return nil, errors.WithStack(data.ErrUsernameIsBusy)
	}

	id, err := redis.Int64(conn.Do("INCR", keyLastUserID()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last user id")
	}

	user := &data.User{
		ID:       id,
		Username: username,
	}
	if err := r.putUser(ctx, conn, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *RedisRepository) GetUserByID(ctx context.Context, id int64) (*data.User, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return r.getUserByID(ctx, conn, id)
}

func (r *RedisRepository) GetUserByUsername(ctx context.Context, username string) (*data.User, error) {
	conn := r.pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("GET", keyUsername(username)))
	switch err {
	case redis.ErrNil:
		return nil, errors.WithStack(data.ErrUserNotFound)
	case nil:
	default:
		return nil, errors.Wrap(err, "failed to get user")
	}

	return r.getUserByID(ctx, conn, id)
}

func (r *RedisRepository) UpdateUser(ctx context.Context, user *data.User) error {
	conn := r.pool.Get()
	defer conn.Close()

	if err := r.checkExistenceUser(ctx, conn, user.ID); err != nil {
		return err
	}

	if err := r.putUser(ctx, conn, user); err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) DeleteUser(ctx context.Context, id int64) error {
	conn := r.pool.Get()
	defer conn.Close()

	if err := r.checkExistenceUser(ctx, conn, id); err != nil {
		return err
	}
	if _, err := conn.Do("DEL", keyUser(id)); err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}

func (r *RedisRepository) getUserByID(ctx context.Context, conn redis.Conn, id int64) (*data.User, error) {
	btsUser, err := redis.Bytes(conn.Do("GET", keyUser(id)))
	switch err {
	case redis.ErrNil:
		return nil, errors.WithStack(data.ErrUserNotFound)
	case nil:
	default:
		return nil, errors.Wrap(err, "failed to get user")
	}

	var user *data.User
	if err := json.Unmarshal(btsUser, user); err != nil {
		return nil, errors.Wrap(err, "failed to decode user")
	}

	return user, nil
}

func (r *RedisRepository) putUser(ctx context.Context, conn redis.Conn, user *data.User) error {
	btsUser, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "failed to encode user")
	}

	if _, err := conn.Do("SET", keyUser(user.ID), btsUser); err != nil {
		return errors.Wrap(err, "failed to set user")
	}

	return nil
}

func (r *RedisRepository) checkExistenceUser(ctx context.Context, conn redis.Conn, id int64) error {
	if ok, err := redis.Bool(conn.Do("EXISTS", keyUser(id))); err != nil {
		return errors.Wrap(err, "failed to check existence of user")
	} else if !ok {
		return errors.WithStack(data.ErrUserNotFound)
	}

	return nil
}

func keyLastUserID() string {
	return "last_user_id"
}

func keyUser(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func keyUsername(username string) string {
	return fmt.Sprintf("username:%s", strings.ToLower(username))
}