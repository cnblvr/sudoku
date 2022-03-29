package frontend

import (
	"github.com/cnblvr/sudoku/data"
	"net/http"
)

// HandleSudokuCreate is a puzzle generator handler/page(TODO).
func (srv *Service) HandleSudokuCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth, log := data.GetAuth(ctx), getLogger(ctx)

	var sudokuGame *data.SudokuGame
	status := func() int {
		sudoku, err := srv.sudokuRepository.GetRandomSudokuByLevel(ctx, data.SudokuLevel(r.URL.Query().Get("level")))
		if err != nil {
			log.Error().Err(err).Msg("failed to get sudoku")
			return http.StatusInternalServerError
		}
		sudokuGame, err = srv.sudokuRepository.CreateSudokuGame(ctx, sudoku.ID, auth.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to create new sudoku game")
			return http.StatusInternalServerError
		}
		return http.StatusOK
	}()
	if status == http.StatusOK {
		redirectPath := data.EndpointSudoku(sudokuGame.ID.String())
		log.Debug().Str("redirect", redirectPath).Msg("success HandleSudokuCreate")
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
		return
	}

	// render error
	http.Error(w, http.StatusText(status), status)
}
