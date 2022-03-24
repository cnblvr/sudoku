package repository

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (r *RedisRepository) CreateSudokuGame(ctx context.Context, sudokuID, userID int64) (*data.SudokuGame, error) {
	conn := r.pool.Get()
	defer conn.Close()

	if ok, err := redis.Bool(conn.Do("EXISTS", keySudoku(sudokuID))); err != nil {
		return nil, errors.Wrap(err, "failed to check existence sudoku")
	} else if !ok {
		return nil, errors.WithStack(data.ErrSudokuNotFound)
	}
	if ok, err := redis.Bool(conn.Do("EXISTS", keyUser(userID))); err != nil {
		return nil, errors.Wrap(err, "failed to check existence user")
	} else if !ok {
		return nil, errors.WithStack(data.ErrUserNotFound)
	}

	seedID, err := redis.Int64(conn.Do("INCR", keyLastSudokuGameSeedID()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last sudoku game seed id")
	}
	id := generateSudokuGameID(seedID)

	game := &data.SudokuGame{
		ID:        id,
		SudokuID:  sudokuID,
		UserID:    userID,
		CreatedAt: data.DateTime{Time: time.Now().UTC()},
	}

	if err := r.putSudokuGame(ctx, conn, game); err != nil {
		return nil, errors.WithStack(err)
	}
	return game, nil
}

func (r *RedisRepository) GetSudokuGameByID(ctx context.Context, id uuid.UUID) (*data.SudokuGame, error) {
	conn := r.pool.Get()
	defer conn.Close()

	game, err := r.getSudokuGame(ctx, conn, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return game, nil
}

func (r *RedisRepository) GetSudokuGameState(ctx context.Context, id uuid.UUID) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	if err := r.isExistsSudokuGame(ctx, conn, id); err != nil {
		return "", errors.WithStack(err)
	}

	state, err := redis.String(conn.Do("GET", keySudokuGameState(id)))
	if err != nil {
		return "", errors.Wrap(err, "failed to get sudoku game state")
	}
	return state, nil
}

func (r *RedisRepository) AddSudokuStep(ctx context.Context, id uuid.UUID, step *data.SudokuStep) error {
	conn := r.pool.Get()
	defer conn.Close()

	if err := r.isExistsSudokuGame(ctx, conn, id); err != nil {
		return errors.WithStack(err)
	}

	bts, err := json.Marshal(step)
	if err != nil {
		return errors.Wrap(err, "failed to encode sudoku step")
	}

	if _, err := conn.Do("RPUSH", keySudokuGameSteps(id), bts); err != nil {
		return errors.Wrap(err, "failed to right push sudoku step")
	}
	return nil
}

func (r *RedisRepository) GetSudokuGameSteps(ctx context.Context, id uuid.UUID) ([]*data.SudokuStep, error) {
	conn := r.pool.Get()
	defer conn.Close()

	if err := r.isExistsSudokuGame(ctx, conn, id); err != nil {
		return nil, errors.WithStack(err)
	}

	btsSlice, err := redis.ByteSlices(conn.Do("LRANGE", keySudokuGameSteps(id), 0, -1))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of sudoku steps")
	}
	steps := make([]*data.SudokuStep, 0, len(btsSlice))
	for _, bts := range btsSlice {
		step := &data.SudokuStep{}
		if err := json.Unmarshal(bts, step); err != nil {
			return nil, errors.Wrap(err, "failed to decode sudoku step")
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func (r RedisRepository) isExistsSudokuGame(ctx context.Context, conn redis.Conn, id uuid.UUID) error {
	if ok, err := redis.Bool(conn.Do("EXISTS", keySudokuGame(id))); err != nil {
		return errors.Wrap(err, "failed to check existence sudoku game")
	} else if !ok {
		return errors.WithStack(data.ErrSudokuGameNotFound)
	}
	return nil
}

func (r *RedisRepository) getSudokuGame(ctx context.Context, conn redis.Conn, id uuid.UUID) (*data.SudokuGame, error) {
	bts, err := redis.Bytes(conn.Do("GET", keySudokuGame(id)))
	switch err {
	case redis.ErrNil:
		return nil, errors.WithStack(data.ErrSudokuGameNotFound)
	case nil:
	default:
		return nil, errors.Wrap(err, "failed to get sudoku game")
	}

	game := &data.SudokuGame{}
	if err := json.Unmarshal(bts, game); err != nil {
		return nil, errors.Wrap(err, "failed to decode sudoku game")
	}

	return game, nil
}

func (r *RedisRepository) putSudokuGame(ctx context.Context, conn redis.Conn, game *data.SudokuGame) error {
	bts, err := json.Marshal(game)
	if err != nil {
		return errors.Wrap(err, "failed to encode sudoku game")
	}

	if _, err := conn.Do("SET", keySudokuGame(game.ID), bts); err != nil {
		return errors.Wrap(err, "failed to set sudoku game")
	}

	return nil
}

func generateSudokuGameID(seedID int64) uuid.UUID {
	idHash := md5.Sum([]byte(fmt.Sprintf("sudoku_game_%d_emag_ukodus", seedID)))
	return uuid.FromBytesOrNil(idHash[:])
}

func keyLastSudokuGameSeedID() string {
	return "last_sudoku_game_seed_id"
}

func keySudokuGame(id uuid.UUID) string {
	return fmt.Sprintf("sudoku_game:%s", id.String())
}

func keySudokuGameState(id uuid.UUID) string {
	return fmt.Sprintf("%s:state", keySudokuGame(id))
}

func keySudokuGameSteps(id uuid.UUID) string {
	return fmt.Sprintf("%s:steps", keySudokuGame(id))
}
