package library_puzzles

import (
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/library_puzzles/sudoku_classic"
)

func GetGenerator(t data.SudokuType) (data.SudokuGenerator, error) {
	switch t {
	case data.SudokuClassic:
		return sudoku_classic.Generator{}, nil
	default:
		return nil, data.ErrSudokuTypeUnknown
	}
}
