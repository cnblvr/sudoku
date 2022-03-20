package sudoku

import (
	"context"
	"fmt"
	"github.com/cnblvr/sudoku/model"
	uuid "github.com/satori/go.uuid"
)

func init() {
	websocketPool.Add((*websocketGetPuzzleRequest)(nil), (*websocketGetPuzzleResponse)(nil))
}

type websocketGetPuzzleRequest struct {
	SessionID string `json:"sessionID"`
}

func (websocketGetPuzzleRequest) Method() string {
	return "getPuzzle"
}

func (r websocketGetPuzzleRequest) Validate(ctx context.Context) error {
	if r.SessionID == "" {
		return fmt.Errorf("sessionID is empty")
	}
	_, err := uuid.FromString(r.SessionID)
	if err != nil {
		return fmt.Errorf("sessionID is not UUID")
	}
	return nil
}

func (r websocketGetPuzzleRequest) Execute(ctx context.Context) (websocketResponse, error) {
	srv := ctx.Value("srv").(*Service)
	redis := srv.redis.Get()
	defer redis.Close()

	session, err := model.SudokuSessionByIDString(redis, r.SessionID)
	if err != nil {
		return websocketGetPuzzleResponse{}, fmt.Errorf("internal server error")
	}
	puzzle, err := session.Sudoku().Puzzle()
	if err != nil {
		return websocketGetPuzzleResponse{}, fmt.Errorf("internal server error")
	}

	return websocketGetPuzzleResponse{
		Puzzle: puzzle,
	}, nil
}

// TODO handle and test
type websocketGetPuzzleResponse struct {
	Puzzle string `json:"puzzle"`
}

func (websocketGetPuzzleResponse) Method() string {
	return "getPuzzle"
}

func (r websocketGetPuzzleResponse) Validate(ctx context.Context) error {
	return nil
}

func (r websocketGetPuzzleResponse) Execute(ctx context.Context) error {
	return nil
}
