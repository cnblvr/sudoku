package sudoku_classic

import (
	"github.com/cnblvr/sudoku/data"
	"github.com/rs/zerolog"
	"math/rand"
	"testing"
)

// Step by step generation and shuffling of the puzzle with output to the console.
func TestManualGenerate(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	b := generateSudokuBoard(rnd)
	t.Logf("base\n%s", b.debug())

	b.swapLines(data.Horizontal, 0, 1)
	t.Logf("swap horizontal 0 and 1\n%s", b.debug())

	b.swapLines(data.Vertical, 0, 1)
	t.Logf("swap vertical 0 and 1\n%s", b.debug())

	b.reflect(data.Horizontal)
	t.Logf("reflect horizontal\n%s", b.debug())

	b.reflect(data.Vertical)
	t.Logf("reflect vertical\n%s", b.debug())

	b.rotate(data.Rotate0)
	t.Logf("rotate to 0\n%s", b.debug())

	b.rotate(data.Rotate90)
	t.Logf("rotate to 90\n%s", b.debug())

	b.rotate(data.Rotate180)
	t.Logf("rotate to 180\n%s", b.debug())

	b.rotate(data.Rotate270)
	t.Logf("rotate to 270\n%s", b.debug())
}

func TestSimpleGeneration(t *testing.T) {
	const seed = 2

	s := NewSudoku(seed)
	t.Logf("base     %s\n%s", s.board.String(), s.board.debug())
	t.Logf("%s\n%s", s.puzzle.String(), s.puzzle.debug())
	t.Logf("count of hints = %d", s.puzzle.CountHints())

}

// Checking the puzzle generator for the uniqueness of the seed.
func TestSeed(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	seeds := []int64{
		0,
		1,
		238978,
		rand.Int63(),
	}
	for _, seed := range seeds {
		board := NewSudoku(seed).board.String()
		for i := 0; i < 10000; i++ {
			if NewSudoku(seed).board.String() != board {
				t.Errorf("seed generate various puzzles")
				continue
			}
		}
	}
}

// Checking the uniqueness of puzzles for many seeds.
// TODO: increase the uniqueness of puzzles
func TestUnique(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	unique := make(map[string]int64)
	for i := int64(0); i < 1000000; i++ {
		s := NewSudoku(i)
		if seed, exists := unique[s.board.String()]; exists {
			t.Errorf("seeds %d and %d generate same boards", seed, i)
			continue
		}
		unique[s.board.String()] = i
	}
}
