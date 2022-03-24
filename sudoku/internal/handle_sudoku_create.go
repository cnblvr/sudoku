package sudoku

import (
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/library_puzzles"
	"net/http"
)

// HandleSudokuCreate is a puzzle generator handler/page(TODO).
func (srv *Service) HandleSudokuCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth, log := getAuth(ctx), getLogger(ctx)
	const seed = int64(3) // todo seed
	const typ = data.SudokuClassic

	var sudokuGame *data.SudokuGame
	status := func() int {
		generator, err := library_puzzles.GetGenerator(typ)
		if err != nil {
			log.Error().Err(err).Msg("failed to get generator")
			return http.StatusBadRequest
		}
		puzzle, solution := generator.Generate(ctx, seed)
		sudoku, err := srv.sudokuRepository.CreateSudoku(ctx, typ, seed, puzzle, solution)
		if err != nil {
			log.Error().Err(err).Msg("failed to create new sudoku")
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
