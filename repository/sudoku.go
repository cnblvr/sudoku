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

func (r *RedisRepository) CreateSudoku(ctx context.Context, typ data.SudokuType, seed int64, puzzle, solution string) (*data.Sudoku, error) {
	conn := r.pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", keyLastSudokuID()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last sudoku id")
	}

	sudoku := &data.Sudoku{
		IDint64Getter: &idInt64{
			id: id,
		},
		Type:     typ,
		Seed:     seed,
		Puzzle:   puzzle,
		Solution: solution,
		CreatedAtGetter: &createdAt{
			createdAt: dateTime{time.Now().UTC()},
		},
	}
	if _, err := conn.Do("SET", keySudoku(id)); err != nil {
		return nil, errors.Wrap(err, "failed to set sudoku")
	}

	return sudoku, nil
}

func (r *RedisRepository) GetSudokuByID(ctx context.Context, id int64) (*data.Sudoku, error) {
	conn := r.pool.Get()
	defer conn.Close()

	btsSudoku, err := redis.Bytes(conn.Do("GET", keySudoku(id)))
	switch err {
	case redis.ErrNil:
		return nil, errors.WithStack(data.ErrSudokuNotFound)
	case nil:
	default:
		return nil, errors.Wrap(err, "failed to get sudoku")
	}

	sudoku := &data.Sudoku{
		IDint64Getter:   &idInt64{},
		CreatedAtGetter: &createdAt{},
	}
	if err := json.Unmarshal(btsSudoku, sudoku); err != nil {
		return nil, errors.Wrap(err, "failed to decode sudoku")
	}

	return sudoku, nil
}

func keyLastSudokuID() string {
	return "last_sudoku_id"
}

func keySudoku(id int64) string {
	return fmt.Sprintf("sudoku:%d", id)
}
