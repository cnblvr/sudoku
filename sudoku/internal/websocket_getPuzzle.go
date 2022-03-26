package sudoku

import (
	"context"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/library_puzzles"
	uuid "github.com/satori/go.uuid"
)

func init() {
	websocketPool.Add((*websocketGetPuzzleRequest)(nil), (*websocketGetPuzzleResponse)(nil))
}

type websocketGetPuzzleRequest struct {
	GameID         string `json:"game_id"`
	NeedCandidates bool   `json:"need_candidates,omitempty"`
}

func (websocketGetPuzzleRequest) Method() string {
	return "getPuzzle"
}

func (r websocketGetPuzzleRequest) Validate(ctx context.Context) error {
	if r.GameID == "" {
		return fmt.Errorf("game_id is empty")
	}
	_, err := uuid.FromString(r.GameID)
	if err != nil {
		return fmt.Errorf("game_id is not UUID")
	}
	return nil
}

func (r websocketGetPuzzleRequest) Execute(ctx context.Context) (websocketResponse, error) {
	srv := ctx.Value("srv").(*Service)

	game, err := srv.sudokuRepository.GetSudokuGameByID(ctx, uuid.FromStringOrNil(r.GameID))
	if err != nil {
		return websocketGetPuzzleResponse{}, fmt.Errorf("internal server error")
	}
	sudoku, err := srv.sudokuRepository.GetSudokuByID(ctx, game.SudokuID)
	if err != nil {
		return websocketGetPuzzleResponse{}, fmt.Errorf("internal server error")
	}

	resp := websocketGetPuzzleResponse{
		Puzzle: sudoku.Puzzle,
	}

	generator, err := library_puzzles.GetGenerator(sudoku.Type)
	if err != nil {
		return websocketGetPuzzleResponse{}, fmt.Errorf("internal server error")
	}
	if r.NeedCandidates {
		resp.Candidates = generator.GetCandidates(ctx, sudoku.Puzzle)
	}

	return resp, nil
}

// TODO handle and test
type websocketGetPuzzleResponse struct {
	Puzzle     string                `json:"puzzle"`
	Candidates data.SudokuCandidates `json:"candidates,omitempty"`
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
