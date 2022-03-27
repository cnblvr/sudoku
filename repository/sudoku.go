package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"time"
)

func (r *RedisRepository) CreateSudoku(ctx context.Context, typ data.SudokuType, seed int64, level data.SudokuLevel, puzzle, solution string) (*data.Sudoku, error) {
	conn := r.pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", keyLastSudokuID()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last sudoku id")
	}

	sudoku := &data.Sudoku{
		ID:        id,
		Type:      typ,
		Seed:      seed,
		Level:     level,
		Puzzle:    puzzle,
		Solution:  solution,
		CreatedAt: data.DateTime{Time: time.Now().UTC()},
	}

	if err := r.putSudoku(ctx, conn, sudoku); err != nil {
		return nil, errors.WithStack(err)
	}

	if _, err := conn.Do("SADD", keySudokuLevel(sudoku.Level), sudoku.ID); err != nil {
		return nil, errors.Wrap(err, "failed to save sudoku to level store")
	}
	if _, err := conn.Do("HSET", keySudokuSeeds(), seed, sudoku.ID); err != nil {
		return nil, errors.Wrap(err, "failed to save sudoku to seed store")
	}

	return sudoku, nil
}

func (r *RedisRepository) GetSudokuByID(ctx context.Context, id int64) (*data.Sudoku, error) {
	conn := r.pool.Get()
	defer conn.Close()

	sudoku, err := r.getSudoku(ctx, conn, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sudoku, nil
}

func (r *RedisRepository) IsExistsSeed(ctx context.Context, seed int64) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	isExists, err := redis.Bool(conn.Do("HEXISTS", keySudokuSeeds(), seed))
	if err != nil {
		return false, errors.Wrap(err, "failed to check existence of seed")
	}

	return isExists, nil
}

func (r *RedisRepository) GetRandomSudokuByLevel(ctx context.Context, level data.SudokuLevel) (*data.Sudoku, error) {
	conn := r.pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("SRANDMEMBER", keySudokuLevel(level)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get random sudoku by level")
	}
	sudoku, err := r.getSudoku(ctx, conn, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sudoku, nil
}

func (r *RedisRepository) getSudoku(ctx context.Context, conn redis.Conn, id int64) (*data.Sudoku, error) {
	bts, err := redis.Bytes(conn.Do("GET", keySudoku(id)))
	switch err {
	case redis.ErrNil:
		return nil, errors.WithStack(data.ErrSudokuNotFound)
	case nil:
	default:
		return nil, errors.Wrap(err, "failed to get sudoku")
	}

	sudoku := &data.Sudoku{}
	if err := json.Unmarshal(bts, sudoku); err != nil {
		return nil, errors.Wrap(err, "failed to decode sudoku")
	}
	return sudoku, nil
}

func (r *RedisRepository) putSudoku(ctx context.Context, conn redis.Conn, sudoku *data.Sudoku) error {
	bts, err := json.Marshal(sudoku)
	if err != nil {
		return errors.Wrap(err, "failed to encode sudoku")
	}

	if _, err := conn.Do("SET", keySudoku(sudoku.ID), bts); err != nil {
		return errors.Wrap(err, "failed to set sudoku")
	}

	return nil
}

func keyLastSudokuID() string {
	return "last_sudoku_id"
}

func keySudoku(id int64) string {
	return fmt.Sprintf("sudoku:%d", id)
}

func keySudokuSeeds() string {
	return "sudoku_seeds"
}

func keySudokuLevel(level data.SudokuLevel) string {
	return fmt.Sprintf("sudoku_level:%s", level.String())
}
