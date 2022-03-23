package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

func (r *RedisRepository) CreateSudoku(ctx context.Context, typ data.SudokuType, seed int64, puzzle, solution string) (*data.Sudoku, error) {
	conn := r.pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", keyLastSudokuID()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last sudoku id")
	}

	sudoku := &data.Sudoku{
		ID:       id,
		Type:     typ,
		Seed:     seed,
		Puzzle:   puzzle,
		Solution: solution,
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
		return nil, errors.Wrap(err, data.SudokuNotFound)
	case nil:
	default:
		return nil, errors.Wrap(err, "failed to get sudoku")
	}

	var sudoku *data.Sudoku
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

//type Sudoku struct {
//	conn redis.Conn
//	id   int64
//}
//
//func SudokuByID(conn redis.Conn, id int64) (Sudoku, bool, error) {
//	_, err := redis.String(conn.Do("GET", keySudoku(id)))
//	switch err {
//	case nil:
//	case redis.ErrNil:
//		return Sudoku{}, false, nil
//	default:
//		return Sudoku{}, false, err
//	}
//	return Sudoku{
//		conn: conn,
//		id:   id,
//	}, true, nil
//}
//
//func (s Sudoku) ID() int64 {
//	return s.id
//}
//
//func (s Sudoku) IsNull() bool {
//	return s.id <= 0
//}
//
//func (s Sudoku) Puzzle() (string, error) {
//	return redis.String(s.conn.Do("GET", keySudoku(s.id)))
//}
//
//func (s Sudoku) Board() (string, error) {
//	return redis.String(s.conn.Do("GET", keySudokuBoard(s.id)))
//}
//
//func NewSudoku(conn redis.Conn, board string, puzzle string) (Sudoku, error) {
//	id, err := redis.Int64(conn.Do("INCR", keyLastSudokuID()))
//	if err != nil {
//		return Sudoku{}, err
//	}
//
//	if _, err := conn.Do("SET", keySudoku(id), puzzle); err != nil {
//		return Sudoku{}, err
//	}
//	if _, err := conn.Do("SET", keySudokuBoard(id), board); err != nil {
//		return Sudoku{}, err
//	}
//
//	return Sudoku{
//		conn: conn,
//		id:   id,
//	}, nil
//}
//

//
//func keySudokuBoard(id int64) string {
//	return fmt.Sprintf("%s:board", keySudoku(id))
//}
