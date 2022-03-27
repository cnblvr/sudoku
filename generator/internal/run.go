package generator

import (
	"context"
	"github.com/cnblvr/sudoku/data"
	"github.com/rs/zerolog/log"
	"math/rand"
)

func (srv *Service) Run() error {
	for {
		if err := srv.GenerateSudokuByRandomSeed(); err != nil {
			log.Error().Err(err).Msg("GenerateRandomSeed failed")
		}
	}
}

func (srv *Service) GenerateSudokuByRandomSeed() error {
	ctx := context.Background()

	// Search random seed for generating sudoku
	var seed int64
	for {
		seed = rand.Int63()
		ok, err := srv.generateRepository.IsExistsSeed(ctx, seed)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}

	// Generate level
	levels := []data.SudokuLevel{data.SudokuRandomEasy, data.SudokuRandomMedium, data.SudokuRandomHard}
	level := levels[rand.Int()%len(levels)]

	// Generate sudoku
	puzzle, solution, err := srv.generator.Generate(ctx, seed, level)

	// Save sudoku
	sudoku, err := srv.sudokuRepository.CreateSudoku(ctx, srv.generator.Type(), seed, level, puzzle, solution)
	if err != nil {
		log.Error().Err(err).Msg("failed to create new sudoku in db")
		return err
	}
	log.Info().Int64("id", sudoku.ID).Str("sudoku_level", level.String()).Msg("new sudoku created and saved")

	return nil
}
