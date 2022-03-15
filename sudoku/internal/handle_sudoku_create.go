package sudoku

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

// HandleSudokuCreate is a puzzle generator page/handler.
func (srv *Service) HandleSudokuCreate(w http.ResponseWriter, r *http.Request) {
	a := getAuth(r)

	// TODO get sudoku puzzle
	_ = a

	redirectPath := "/sudoku/qwe"
	log.Debug().Str("redirect", redirectPath).Msg("success HandleSudokuCreate")
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}
