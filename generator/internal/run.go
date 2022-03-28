package generator

import (
	"context"
	"github.com/cnblvr/sudoku/data"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
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

	// Generate sudoku
	ctxGen, cancelGen := context.WithTimeout(ctx, time.Hour)
	defer cancelGen()
	generatedChan := make(chan data.GeneratedSudoku, 3)
	go srv.generator.Generate(ctxGen, seed, generatedChan)
	for generated := range generatedChan {
		// Save sudoku
		sudoku, err := srv.sudokuRepository.CreateSudoku(ctx, srv.generator.Type(), seed, generated.Level, generated.Puzzle, generated.Solution)
		if err != nil {
			log.Error().Err(err).Msg("failed to create new sudoku in db")
			return err
		}
		log.Info().Int64("id", sudoku.ID).Str("sudoku_level", generated.Level.String()).Msg("new sudoku created and saved")
	}

	return nil
}
