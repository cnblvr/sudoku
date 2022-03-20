package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Sudoku struct {
	conn redis.Conn
	id   int64
}

func SudokuByID(conn redis.Conn, id int64) (Sudoku, bool, error) {
	_, err := redis.String(conn.Do("GET", keySudoku(id)))
	switch err {
	case nil:
	case redis.ErrNil:
		return Sudoku{}, false, nil
	default:
		return Sudoku{}, false, err
	}
	return Sudoku{
		conn: conn,
		id:   id,
	}, true, nil
}

func (s Sudoku) ID() int64 {
	return s.id
}

func (s Sudoku) IsNull() bool {
	return s.id <= 0
}

func (s Sudoku) Puzzle() (string, error) {
	return redis.String(s.conn.Do("GET", keySudoku(s.id)))
}

func (s Sudoku) Board() (string, error) {
	return redis.String(s.conn.Do("GET", keySudokuBoard(s.id)))
}

func NewSudoku(conn redis.Conn, board string, puzzle string) (Sudoku, error) {
	id, err := redis.Int64(conn.Do("INCR", keyLastSudokuID()))
	if err != nil {
		return Sudoku{}, err
	}

	if _, err := conn.Do("SET", keySudoku(id), puzzle); err != nil {
		return Sudoku{}, err
	}
	if _, err := conn.Do("SET", keySudokuBoard(id), board); err != nil {
		return Sudoku{}, err
	}

	return Sudoku{
		conn: conn,
		id:   id,
	}, nil
}

func keySudoku(id int64) string {
	return fmt.Sprintf("sudoku:%d", id)
}

func keyLastSudokuID() string {
	return "last_sudoku_id"
}

func keySudokuBoard(id int64) string {
	return fmt.Sprintf("%s:board", keySudoku(id))
}
