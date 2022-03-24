package repository

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type RedisRepository struct {
	pool *redis.Pool
}

func New(pool *redis.Pool) *RedisRepository {
	return &RedisRepository{
		pool: pool,
	}
}

type idInt64 struct {
	id int64
}

func (i idInt64) ID() int64 {
	return i.id
}

func (i idInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.id)
}

func (i *idInt64) UnmarshalJSON(data []byte) error {

	id, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	i.id = id
	return nil
}

const dateTimeFormat = "2006-01-02T15:04:05.000000Z07:00"

type dateTime struct {
	time.Time
}

type createdAt struct {
	createdAt dateTime
}

func (dt createdAt) CreatedAt() time.Time {
	return dt.createdAt.Time
}

func (dt createdAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.createdAt.Format(dateTimeFormat))
}

func (dt *createdAt) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	t, err := time.Parse(dateTimeFormat, str)
	if err != nil {
		return err
	}
	(*dt).createdAt = dateTime{t}
	return nil
}

type updatedAt struct {
	updatedAt dateTime
}

func (dt updatedAt) UpdatedAt() time.Time {
	return dt.updatedAt.Time
}

func (dt updatedAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.updatedAt.Format(dateTimeFormat))
}

func (dt *updatedAt) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	t, err := time.Parse(dateTimeFormat, str)
	if err != nil {
		return err
	}
	(*dt).updatedAt = dateTime{t}
	return nil
}
