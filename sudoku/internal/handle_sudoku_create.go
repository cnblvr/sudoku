package sudoku

import (
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/model"
	"github.com/cnblvr/sudoku/sudoku/internal/sudoku_classic"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HandleSudokuCreate is a puzzle generator handler/page(TODO).
func (srv *Service) HandleSudokuCreate(w http.ResponseWriter, r *http.Request) {
	//a := getAuth(r)
	redis := srv.redis.Get()
	defer redis.Close()
	const seed = int64(2) // todo seed
	log := log.With().Int64("seed", seed).Logger()

	var sudokuSession model.SudokuSession
	status := func() int {
		sudoku := sudoku_classic.NewSudoku(seed)
		mSudoku, err := model.NewSudoku(redis,
			sudoku.Board().String(),
			sudoku.Puzzle().String(),
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to create new sudoku")
			return http.StatusInternalServerError
		}
		sudokuSession, err = model.NewSudokuSession(redis, mSudoku, model.User{})
		if err != nil {
			log.Error().Err(err).Msg("failed to create new sudoku session")
			return http.StatusInternalServerError
		}
		return http.StatusOK
	}()
	if status == http.StatusOK {
		redirectPath := data.EndpointSudoku(sudokuSession.ID().String())
		log.Debug().Str("redirect", redirectPath).Msg("success HandleSudokuCreate")
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
		return
	}

	// render error
	http.Error(w, http.StatusText(status), status)
}
