package sudoku_classic

import (
	"context"
	"github.com/cnblvr/sudoku/data"
	"github.com/rs/zerolog"
	"math/rand"
	"testing"
)

func TestTransformations(t *testing.T) {
	tests := []struct {
		name       string
		puzzle     string
		fn         func(p sudokuPuzzle)
		wantPuzzle string
	}{
		// SWAP LINES

		{
			name:   "swapLines a and b",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Horizontal, 0, 1)
			},
			wantPuzzle: "456789123123456789789123456891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "swapLines a and a",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Horizontal, 0, 0)
			},
			wantPuzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "swapLines b and a",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Horizontal, 1, 0)
			},
			wantPuzzle: "456789123123456789789123456891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "swapLines a and c",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Horizontal, 0, 2)
			},
			wantPuzzle: "789123456456789123123456789891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "swapLines h and i",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Horizontal, 7, 8)
			},
			wantPuzzle: "123456789456789123789123456891234567234567891567891234678912345345678912912345678",
		},
		{
			name:   "swapLines 1 and 2",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Vertical, 0, 1)
			},
			wantPuzzle: "213456789546789123879123456981234567324567891657891234768912345192345678435678912",
		},
		{
			name:   "swapLines 2 and 1",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Vertical, 1, 0)
			},
			wantPuzzle: "213456789546789123879123456981234567324567891657891234768912345192345678435678912",
		},
		{
			name:   "swapLines 8 and 9",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapLines(data.Vertical, 7, 8)
			},
			wantPuzzle: "123456798456789132789123465891234576234567819567891243678912354912345687345678921",
		},

		// SWAP "BIG" LINES

		{
			name:   "swapBigLines a-c and a-c",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapBigLines(data.Horizontal, 0, 0)
			},
			wantPuzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "swapBigLines a-c and d-f",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapBigLines(data.Horizontal, 0, 1)
			},
			wantPuzzle: "891234567234567891567891234123456789456789123789123456678912345912345678345678912",
		},
		{
			name:   "swapBigLines a-c and e-i",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapBigLines(data.Horizontal, 0, 2)
			},
			wantPuzzle: "678912345912345678345678912891234567234567891567891234123456789456789123789123456",
		},
		{
			name:   "swapBigLines e-i and a-c",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapBigLines(data.Horizontal, 2, 0)
			},
			wantPuzzle: "678912345912345678345678912891234567234567891567891234123456789456789123789123456",
		},
		{
			name:   "swapBigLines 1-3 and 4-6",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapBigLines(data.Vertical, 0, 1)
			},
			wantPuzzle: "456123789789456123123789456234891567567234891891567234912678345345912678678345912",
		},
		{
			name:   "swapBigLines 7-9 and 4-6",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapBigLines(data.Vertical, 2, 1)
			},
			wantPuzzle: "123789456456123789789456123891567234234891567567234891678345912912678345345912678",
		},

		// REFLECT

		{
			name:   "reflect horizontal",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.reflect(data.Horizontal)
			},
			wantPuzzle: "987654321321987654654321987765432198198765432432198765543219876876543219219876543",
		},
		{
			name:   "reflect horizontal double",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.reflect(data.Horizontal)
				p.reflect(data.Horizontal)
			},
			wantPuzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "reflect vertical",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.reflect(data.Vertical)
			},
			wantPuzzle: "345678912912345678678912345567891234234567891891234567789123456456789123123456789",
		},
		{
			name:   "reflect vertical double",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.reflect(data.Vertical)
				p.reflect(data.Vertical)
			},
			wantPuzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
		},

		// ROTATE

		{
			name:   "rotate 0 degrees",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.rotate(data.Rotate0)
			},
			wantPuzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "rotate 90 degrees",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.rotate(data.Rotate90)
			},
			wantPuzzle: "936714582825693471714582369693471258582369147471258936369147825258936714147825693",
		},
		{
			name:   "rotate 180 degrees",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.rotate(data.Rotate180)
			},
			wantPuzzle: "219876543876543219543219876432198765198765432765432198654321987321987654987654321",
		},
		{
			name:   "rotate 270 degrees",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.rotate(data.Rotate270)
			},
			wantPuzzle: "396528741417639852528741963639852174741963285852174396963285417174396528285417639",
		},
		{
			name:   "rotate 360 degrees",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.rotate(data.Rotate270 + 1)
			},
			wantPuzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
		},
		{
			name:   "rotate 450 degrees",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.rotate(data.Rotate270 + 2)
			},
			wantPuzzle: "936714582825693471714582369693471258582369147471258936369147825258936714147825693",
		},

		// SWAP DIGITS

		{
			name:   "swap 1 to 6",
			puzzle: "123456789456789123789123456891234567234567891567891234678912345912345678345678912",
			fn: func(p sudokuPuzzle) {
				p.swapDigits(1, 6)
			},
			wantPuzzle: "623451789451789623789623451896234517234517896517896234178962345962345178345178962",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := sudokuPuzzleFromString(test.puzzle)
			test.fn(p)
			got := p.String()
			if got != test.wantPuzzle {
				t.Errorf("not equal\nGot:  %s\nWant: %s", got, test.wantPuzzle)
			}
		})
	}
}

// Checking the puzzle generator for the uniqueness of the seed.
func TestSeed(t *testing.T) {
	seeds := []int64{
		0,
		1,
		238978,
		rand.Int63(),
	}
	for _, seed := range seeds {
		ctx := context.Background()
		puzzle, solution := Generator{}.Generate(ctx, seed)
		for i := 0; i < 5; i++ {
			ctx := context.Background()
			newPuzzle, newSolution := Generator{}.Generate(ctx, seed)
			if newPuzzle != puzzle {
				t.Errorf("seed generate various puzzles")
				continue
			}
			if newSolution != solution {
				t.Errorf("seed generate various solutions")
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
	for i := int64(0); i < 100; i++ {
		ctx := context.Background()
		puzzle, _ := Generator{}.Generate(ctx, i)
		if seed, exists := unique[puzzle]; exists {
			t.Errorf("seeds %d and %d generate same puzzles", seed, i)
			continue
		}
		unique[puzzle] = i
	}
}

//func TestSimpleGeneration(t *testing.T) {
//	const seed = 2
//
//	s := NewSudoku(seed)
//	t.Logf("base     %s\n%s", s.board.String(), s.board.debug())
//	t.Logf("%s\n%s", s.puzzle.String(), s.puzzle.debug())
//	t.Logf("count of hints = %d", s.puzzle.CountHints())
//
//}
//
//
