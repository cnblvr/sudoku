package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

func (r *RedisRepository) CreateToken(ctx context.Context, userID int64, expiration time.Duration) (*data.TokenInfo, error) {
	conn := r.pool.Get()
	defer conn.Close()

	if expiration == 0 {
		expiration = data.DefaultTokenExpiration
	}

	token := &data.TokenInfo{
		ID:         uuid.NewV4(),
		Expiration: expiration,
		UserID:     userID,
	}
	btsToken, err := json.Marshal(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode token info")
	}
	if _, err := conn.Do("SETEX", keyToken(token.ID), redisExpiration(token.Expiration), btsToken); err != nil {
		return nil, errors.Wrap(err, "failed to set token info")
	}
	return token, nil
}

func (r *RedisRepository) GetToken(ctx context.Context, id uuid.UUID) (*data.TokenInfo, error) {
	conn := r.pool.Get()
	defer conn.Close()

	btsToken, err := redis.Bytes(conn.Do("GET", keyToken(id)))
	switch err {
	case nil:
	case redis.ErrNil:
		return nil, errors.WithStack(data.ErrTokenNotFound)
	default:
		return nil, errors.Wrap(err, "failed to get token info")
	}
	var token data.TokenInfo
	if err := json.Unmarshal(btsToken, &token); err != nil {
		return nil, errors.Wrap(err, "failed to decode token info")
	}
	if _, err := conn.Do("EXPIRE", keyToken(id), redisExpiration(token.Expiration)); err != nil {
		return nil, errors.Wrap(err, "failed to set expiration on token info")
	}
	return &token, nil
}

func (r *RedisRepository) DeleteToken(ctx context.Context, id uuid.UUID) error {
	conn := r.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", keyToken(id)); err != nil {
		return errors.Wrap(err, "failed to delete token info")
	}
	return nil
}

func (r *RedisRepository) CreateUser(ctx context.Context, username string) (*data.User, error) {
	conn := r.pool.Get()
	defer conn.Close()

	if ok, err := redis.Bool(conn.Do("EXISTS", keyUsername(username))); err != nil {
		return nil, errors.Wrap(err, "failed to check existence owner of username")
	} else if ok {
		return nil, errors.WithStack(data.ErrUsernameIsBusy)
	}

	id, err := redis.Int64(conn.Do("INCR", keyLastUserID()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last user id")
	}

	if _, err := conn.Do("SET", keyUsername(username), id); err != nil {
		return nil, errors.Wrap(err, "failed to register username")
	}

	user := &data.User{
		ID:        id,
		Username:  username,
		CreatedAt: data.DateTime{Time: time.Now().UTC()},
	}

	if err := r.putUser(ctx, conn, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *RedisRepository) GetUserByID(ctx context.Context, id int64) (*data.User, error) {
	conn := r.pool.Get()
	defer conn.Close()

	user, err := r.getUser(ctx, conn, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
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

	user, err := r.getUser(ctx, conn, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

func (r *RedisRepository) UpdateUser(ctx context.Context, user *data.User) error {
	conn := r.pool.Get()
	defer conn.Close()

	oldUser, err := r.getUser(ctx, conn, user.ID)
	if err != nil {
		return err
	}
	if oldUser.Username != user.Username {
		if err := r.occupyUsername(ctx, conn, oldUser.Username, user.Username, user.ID); err != nil {
			return err
		}
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

func (r *RedisRepository) getUser(ctx context.Context, conn redis.Conn, id int64) (*data.User, error) {
	btsUser, err := redis.Bytes(conn.Do("GET", keyUser(id)))
	switch err {
	case redis.ErrNil:
		return nil, errors.WithStack(data.ErrUserNotFound)
	case nil:
	default:
		return nil, errors.Wrap(err, "failed to get user")
	}

	user := &data.User{}
	if err := json.Unmarshal(btsUser, user); err != nil {
		return nil, errors.Wrap(err, "failed to decode user")
	}
	return user, nil
}

func (r *RedisRepository) putUser(ctx context.Context, conn redis.Conn, user *data.User) error {
	user.UpdatedAt = data.DateTime{Time: time.Now().UTC()}
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

func (r *RedisRepository) occupyUsername(ctx context.Context, conn redis.Conn, oldUsername, newUsername string, id int64) error {
	owner, err := redis.Int64(conn.Do("GET", keyUsername(newUsername)))
	switch err {
	case redis.ErrNil:
	case nil:
		if owner != id {
			return errors.WithStack(data.ErrUsernameIsBusy)
		}
	default:
		return errors.Wrap(err, "failed to get owner of username")
	}

	if _, err := conn.Do("SET", keyUsername(newUsername), id); err != nil {
		return errors.Wrap(err, "failed to register new username")
	}
	if _, err := conn.Do("DEL", keyUsername(oldUsername)); err != nil {
		return errors.Wrap(err, "failed to unregister old username")
	}
	return nil
}

func keyToken(id uuid.UUID) string {
	return fmt.Sprintf("token:%s", id.String())
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
