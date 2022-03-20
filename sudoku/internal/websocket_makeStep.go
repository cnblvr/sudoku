package sudoku

import (
	"context"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/model"
	"github.com/cnblvr/sudoku/sudoku/internal/sudoku_classic"
	uuid "github.com/satori/go.uuid"
	"sort"
)

func init() {
	websocketPool.Add((*websocketMakeStepRequest)(nil), (*websocketMakeStepResponse)(nil))
}

type websocketMakeStepRequest struct {
	SessionID string `json:"sessionID"`
	State     string `json:"state"`
}

func (websocketMakeStepRequest) Method() string {
	return "makeStep"
}

func (r websocketMakeStepRequest) Validate(ctx context.Context) error {
	if r.SessionID == "" {
		return fmt.Errorf("sessionID is empty")
	}
	if _, err := uuid.FromString(r.SessionID); err != nil {
		return fmt.Errorf("sessionID is not UUID")
	}
	if len(r.State) < 81 {
		return fmt.Errorf("state format invalid")
	}
	return nil
}

func (r websocketMakeStepRequest) Execute(ctx context.Context) (websocketResponse, error) {
	srv := ctx.Value("srv").(*Service)
	redis := srv.redis.Get()
	defer redis.Close()

	uniqueErrs := make(map[data.Point]struct{})

	session, err := model.SudokuSessionByIDString(redis, r.SessionID)
	if err != nil {
		return websocketMakeStepResponse{}, fmt.Errorf("internal server error")
	}

	boardStr, err := session.Sudoku().Board()
	if err != nil {
		return websocketMakeStepResponse{}, fmt.Errorf("internal server error")
	}
	if r.State == boardStr {
		// WIN
		return websocketMakeStepResponse{
			Win: true,
		}, nil
	}

	userState := sudoku_classic.PuzzleFromString(r.State)
	for _, p := range userState.FindUserErrors() {
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
	return resp, nil
}

// TODO handle and test
type websocketMakeStepResponse struct {
	Errors []data.Point `json:"errors,omitempty"`
	Win    bool         `json:"win,omitempty"`
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
