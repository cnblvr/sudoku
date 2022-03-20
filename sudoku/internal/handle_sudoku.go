package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/model"
	"github.com/cnblvr/sudoku/sudoku/static"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	ErrorSudokuNotFound = "Sudoku not found."
)

// HandleSudoku renders page with puzzle.
func (srv *Service) HandleSudoku(w http.ResponseWriter, r *http.Request) {
	redis := srv.redis.Get()
	defer redis.Close()

	var d struct {
		Session      string
		ErrorMessage string
	}

	d.ErrorMessage = func() string {
		var sudokuSession model.SudokuSession
		var err error
		sessionID, ok := mux.Vars(r)["session_id"]
		if !ok {
			log.Error().Msg("'session_id' not found in mux.Vars")
			return ErrorBadRequest
		}
		if sudokuSession, err = model.SudokuSessionByIDString(redis, sessionID); err != nil {
			log.Warn().Err(err).Msgf("sudoku session '%s' not found", sessionID)
			return ErrorSudokuNotFound
		}
		d.Session = sudokuSession.ID().String()

		// TODO sudokuSession.Sudoku().AddUserID()
		_ = sudokuSession

		return ""
	}()

	args := templates.Args{
		Header: templates.Header{
			Title: fmt.Sprintf("sudoku"),
			Css:   []string{static.CssSudoku},
		},
		Data: d,
		Footer: templates.Footer{
			Js: []string{static.JsSudoku},
		},
	}
	srv.executeTemplate(w, "page_sudoku", args)
}
