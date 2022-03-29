package frontend

import (
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/frontend/static"
	"github.com/cnblvr/sudoku/frontend/templates"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

const (
	ErrorSudokuNotFound = "Sudoku not found."
	ErrorAccessDenied   = "Access denied."
)

// HandleSudoku renders page with puzzle.
func (srv *Service) HandleSudoku(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth, log := data.GetCtxAuth(ctx), data.GetCtxLog(ctx)
	var d struct {
		GameID       string
		ErrorMessage string
	}

	d.ErrorMessage = func() string {
		var sudokuGame *data.SudokuGame
		var err error
		var gameID uuid.UUID
		if gameIDStr, ok := mux.Vars(r)["game_id"]; !ok {
			log.Error().Msg("'game_id' not found in mux.Vars")
			return ErrorBadRequest
		} else if gameID, err = uuid.FromString(gameIDStr); err != nil {
			log.Error().Err(err).Msg("failed to parse sudoku game id as uuid")
			return ErrorBadRequest
		}
		if sudokuGame, err = srv.sudokuRepository.GetSudokuGameByID(ctx, gameID); err != nil {
			log.Warn().Err(err).Msgf("sudoku session '%s' not found", gameID.String())
			return ErrorSudokuNotFound
		}
		if !sudokuGame.ValidateByUser(auth) {
			return ErrorAccessDenied
		}
		d.GameID = sudokuGame.ID.String()
		return ""
	}()

	// TODO sudokuSession.Sudoku().AddUserID()

	args := templates.Args{
		Header: templates.NewHeader(ctx, templates.Header{
			Title: fmt.Sprintf("sudoku"),
			Css:   []string{static.CssSudoku},
			Js:    []string{static.JsSudoku, static.JsWs},
		}),
		Data:   d,
		Footer: templates.Footer{},
	}
	srv.executeTemplate(w, "page_sudoku", args)
}
