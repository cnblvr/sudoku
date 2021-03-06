package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/gomodule/redigo/redis"
	"time"
)

type User struct {
	conn redis.Conn
	id   int64
}

const dateTimeFormat = "2006-01-02T15:04:05.999999Z07:00"

func UserByID(conn redis.Conn, id int64) (User, bool, error) {
	_, err := redis.String(conn.Do("GET", keyUser(id)))
	switch err {
	case nil:
	case redis.ErrNil:
		return User{}, false, nil
	default:
		return User{}, false, err
	}
	return User{
		conn: conn,
		id:   id,
	}, true, nil
}

func UserByUsername(conn redis.Conn, username string) (User, bool, error) {
	id, err := redis.Int64(conn.Do("GET", keyUsername(username)))
	switch err {
	case nil:
	case redis.ErrNil:
		return User{}, false, nil
	default:
		return User{}, false, err
	}
	return User{
		conn: conn,
		id:   id,
	}, true, nil
}

func NewUser(conn redis.Conn, username string) (User, error) {
	// TODO transactions
	isVacant, err := IsUsernameVacant(conn, username)
	if err != nil {
		return User{}, err
	}
	if !isVacant {
		return User{}, fmt.Errorf("username is busy")
	}

	id, err := redis.Int64(conn.Do("INCR", keyLastUserID()))
	if err != nil {
		return User{}, err
	}

	if err := occupyUsername(conn, username, id); err != nil {
		return User{}, err
	}

	if _, err := conn.Do("SET", keyUser(id), username); err != nil {
		return User{}, err
	}

	if _, err = conn.Do("SET", keyUserCreatedAt(id), time.Now().UTC().Format(dateTimeFormat)); err != nil {
		return User{}, err
	}

	infoBts, err := json.Marshal(data.UserInfo{})
	if err != nil {
		return User{}, err
	}
	if _, err = conn.Do("SET", keyUserInfo(id), infoBts); err != nil {
		return User{}, err
	}

	return User{
		conn: conn,
		id:   id,
	}, nil
}

func (u User) UserInfo() (data.UserInfo, error) {
	username, err := u.Username()
	if err != nil {
		return data.UserInfo{}, err
	}
	infoBts, err := redis.Bytes(u.conn.Do("GET", keyUserInfo(u.id)))
	if err != nil {
		if err == redis.ErrNil {
			return data.UserInfo{Username: username}, nil
		}
		return data.UserInfo{}, err
	}
	var info data.UserInfo
	if err := json.Unmarshal(infoBts, &info); err != nil {
		return data.UserInfo{}, err
	}
	info.Username = username
	return info, nil
}

func (u User) SetUserInfo(info data.UserInfo) error {
	infoBts, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = u.conn.Do("SET", keyUserInfo(u.id), infoBts)
	return err
}

// ID returns user's id.
func (u User) ID() int64 {
	return u.id
}

func (u User) IsNull() bool {
	return u.id <= 0
}

// Username returns the username.
func (u User) Username() (string, error) {
	return redis.String(u.conn.Do("GET", keyUser(u.id)))
}

func IsUsernameVacant(conn redis.Conn, username string) (bool, error) {
	isExists, err := redis.Bool(conn.Do("EXISTS", keyUsername(username)))
	if err != nil {
		return false, err
	}
	if isExists {
		return false, nil
	}
	return true, nil
}

func occupyUsername(conn redis.Conn, username string, id int64) error {
	_, err := conn.Do("SET", keyUsername(username), id)
	return err
}

// SetUsername changes the username, leaving the old one vacant.
func (u User) SetUsername(username string) error {
	isVacant, err := IsUsernameVacant(u.conn, username)
	if err != nil {
		return err
	}
	if !isVacant {
		return fmt.Errorf("username is busy")
	}

	oldUsername, err := u.Username()
	if err != nil {
		return err
	}

	if err := occupyUsername(u.conn, username, u.id); err != nil {
		return err
	}

	if _, err = u.conn.Do("SET", keyUser(u.id), username); err != nil {
		return err
	}

	if _, err = u.conn.Do("DEL", keyUsername(oldUsername)); err != nil {
		return err
	}

	return nil
}

// PasswordHash returns a hash generated by the bcrypt algorithm.
func (u User) PasswordHash() ([]byte, error) {
	return redis.Bytes(u.conn.Do("GET", keyUserPasswordHash(u.id)))
}

// SetPasswordHash changes a hash.
func (u User) SetPasswordHash(hash []byte) error {
	_, err := u.conn.Do("SET", keyUserPasswordHash(u.id), hash)
	return err
}

// PasswordSalt returns the salt for hashing the user's password.
func (u User) PasswordSalt() ([]byte, error) {
	return redis.Bytes(u.conn.Do("GET", keyUserPasswordSalt(u.id)))
}

// SetPasswordSalt changes the salt.
func (u User) SetPasswordSalt(salt []byte) error {
	_, err := u.conn.Do("SET", keyUserPasswordSalt(u.id), salt)
	return err
}

// CreatedAt returns the date and time the user was registered.
func (u User) CreatedAt() (time.Time, error) {
	str, err := redis.String(u.conn.Do("GET", keyUserCreatedAt(u.id)))
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(dateTimeFormat, str)
}

func keyUser(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func keyUsername(username string) string {
	return fmt.Sprintf("username:%s", base64.StdEncoding.EncodeToString([]byte(username)))
}

func keyLastUserID() string {
	return "last_user_id"
}

func keyUserPasswordHash(id int64) string {
	return fmt.Sprintf("%s:password_hash", keyUser(id))
}

func keyUserPasswordSalt(id int64) string {
	return fmt.Sprintf("%s:password_salt", keyUser(id))
}

func keyUserCreatedAt(id int64) string {
	return fmt.Sprintf("%s:created_at", keyUser(id))
}

func keyUserInfo(id int64) string {
	return fmt.Sprintf("%s:info", keyUser(id))
}
