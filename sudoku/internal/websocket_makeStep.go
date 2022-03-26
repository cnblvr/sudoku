package sudoku

import (
	"context"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/library_puzzles"
	uuid "github.com/satori/go.uuid"
	"sort"
)

func init() {
	websocketPool.Add((*websocketMakeStepRequest)(nil), (*websocketMakeStepResponse)(nil))
}

type websocketMakeStepRequest struct {
	GameID         string `json:"game_id"`
	State          string `json:"state"`
	NeedCandidates bool   `json:"need_candidates,omitempty"`
}

func (websocketMakeStepRequest) Method() string {
	return "makeStep"
}

func (r websocketMakeStepRequest) Validate(ctx context.Context) error {
	if r.GameID == "" {
		return fmt.Errorf("game_id is empty")
	}
	if _, err := uuid.FromString(r.GameID); err != nil {
		return fmt.Errorf("game_id is not UUID")
	}
	if len(r.State) < 81 {
		return fmt.Errorf("state format invalid")
	}
	return nil
}

func (r websocketMakeStepRequest) Execute(ctx context.Context) (websocketResponse, error) {
	srv := ctx.Value("srv").(*Service)

	uniqueErrs := make(map[data.Point]struct{})

	game, err := srv.sudokuRepository.GetSudokuGameByID(ctx, uuid.FromStringOrNil(r.GameID))
	if err != nil {
		return websocketMakeStepResponse{}, fmt.Errorf("internal server error")
	}
	sudoku, err := srv.sudokuRepository.GetSudokuByID(ctx, game.SudokuID)
	if err != nil {
		return websocketMakeStepResponse{}, fmt.Errorf("internal server error")
	}

	if r.State == sudoku.Solution {
		// WIN
		return websocketMakeStepResponse{
			Win: true,
		}, nil
	}

	generator, err := library_puzzles.GetGenerator(sudoku.Type)
	if err != nil {
		return websocketMakeStepResponse{}, fmt.Errorf("internal server error")
	}
	for _, p := range generator.FindUserErrors(ctx, r.State) {
		uniqueErrs[p] = struct{}{}
	}

	// TODO new method "compare with answer" and use this function
	//board := sudoku_classic.PuzzleFromString(boardStr)
	//for _, p := range board.FindErrors(userState) {
	//	uniqueErrs[p] = struct{}{}
	//}

	resp := websocketMakeStepResponse{}
	for p := range uniqueErrs {
		resp.Errors = append(resp.Errors, p)
	}
	sort.Slice(resp.Errors, func(i, j int) bool {
		if resp.Errors[i].Row != resp.Errors[j].Row {
			return resp.Errors[i].Row < resp.Errors[j].Row
		}
		return resp.Errors[i].Col < resp.Errors[j].Col
	})
	if r.NeedCandidates {
		resp.Candidates = generator.GetCandidates(ctx, r.State)
	}
	return resp, nil
}

// TODO handle and test
type websocketMakeStepResponse struct {
	Errors     []data.Point          `json:"errors,omitempty"`
	Win        bool                  `json:"win,omitempty"`
	Candidates data.SudokuCandidates `json:"candidates,omitempty"`
}

func (websocketMakeStepResponse) Method() string {
	return "makeStep"
}

func (r websocketMakeStepResponse) Validate(ctx context.Context) error {
	return nil
}

func (r websocketMakeStepResponse) Execute(ctx context.Context) error {
	return nil
}
