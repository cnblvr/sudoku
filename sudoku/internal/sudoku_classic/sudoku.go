package sudoku_classic

import (
	"fmt"
	"github.com/cnblvr/sudoku/data"
	stdlog "log"
	"math/rand"
	"strconv"
)

// Sudoku is the basic structure of a 9x9 Sudoku puzzle.
type Sudoku struct {
	// seed allows you to create a unique puzzle
	seed int64
	// board stores the solution to the puzzle
	board sudokuBoard
	// puzzle stores hints for the user
	puzzle sudokuPuzzle
}

func (s Sudoku) Board() data.SudokuBoard {
	return s.board
}

func (s Sudoku) Puzzle() data.SudokuPuzzle {
	return s.puzzle
}

// NewSudoku creates a new puzzle and removes some hints depending on the level.
// seed is used to create a unique puzzle.
func NewSudoku(seed int64) data.Sudoku {
	s := Sudoku{}
	s.seed = seed
	// randomizer for full puzzle generation
	rnd := rand.New(rand.NewSource(seed))

	// puzzle generation without shuffling
	s.board = generateSudokuBoard(rnd)

	// swap of horizontal or vertical lines within one "big" line
	// TODO: imperfect randomization
	for i := 0; i < (rnd.Int()%1024)+1024; i++ {
		typ := data.Horizontal
		if rnd.Int()%2 == 1 {
			typ = data.Vertical
		}
		line := rnd.Int() % 9
		s.board.swapLines(typ, line, neighborLine(line, rnd.Int()%2))
	}

	// TODO: swap "big" lines

	// horizontal reflection
	if rnd.Int()%2 == 1 {
		s.board.reflect(data.Horizontal)
	}
	// vertical reflection
	if rnd.Int()%2 == 1 {
		s.board.reflect(data.Vertical)
	}

	// rotate the puzzle by a random angle
	s.board.rotate(data.RotationType(rnd.Int() % 4))

	s.puzzle = make([][]int8, 9)
	for row := 0; row < 9; row++ {
		s.puzzle[row] = make([]int8, 9)
		copy(s.puzzle[row], s.board[row])
	}

	removes := 81
mainFor:
	for removes > 0 {
		removes = 0
		for _, p := range sudokuRandomPoints(rnd) {
			if removes >= 46 {
				break mainFor // todo level
			}
			digit := s.puzzle[p.Row][p.Col]
			if digit == 0 {
				continue
			}
			s.puzzle[p.Row][p.Col] = 0
			if len(s.puzzle.solveBruteForce(2)) != 1 {
				s.puzzle[p.Row][p.Col] = digit
				continue
			} else {
				removes++
			}
		}
		stdlog.Printf("iteration. removes: %d", removes)
	}

	return &s
}

// Puzzle generation without shuffling.
func generateSudokuBoard(rnd *rand.Rand) sudokuBoard {
	b := make([][]int8, 9)
	for i := 0; i < 9; i++ {
		b[i] = make([]int8, 9)
	}
	// Generate first line randomly
	digits := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9}
	line := make([]int8, 0, 9)
	for len(digits) > 0 {
		idx := rnd.Int() % len(digits)
		line = append(line, digits[idx])
		digits = append(digits[:idx], digits[idx+1:]...)
	}
	copy(b[0], line)

	// The second line is the offset of the first line to the left by 3
	line = append(line[3:9], line[:3]...)
	copy(b[1], line)

	// The third line is the offset of the second line to the left by 3
	line = append(line[3:9], line[:3]...)
	copy(b[2], line)

	// First "big" horizontally line is completed. Next lines generate by this algorithm:
	//  line n:   is offset of the previous line to the left by 1
	//  line n+1: is offset of the previous line to the left by 3
	//  line n+2: is offset of the previous line to the left by 3
	line = append(line[1:9], line[0]) // n
	copy(b[3], line)
	line = append(line[3:9], line[:3]...) // n+1
	copy(b[4], line)
	line = append(line[3:9], line[:3]...) // n+2
	copy(b[5], line)

	// Generation of third "big" horizontally line.
	line = append(line[1:9], line[0]) // n
	copy(b[6], line)
	line = append(line[3:9], line[:3]...) // n+1
	copy(b[7], line)
	line = append(line[3:9], line[:3]...) // n+2
	copy(b[8], line)

	return b
}

// Calculation of the neighboring line.
// lineIdx in the range [0,8].
// neighbor can take values {0,1}.
func neighborLine(lineIdx int, neighbor int) int {
	switch neighbor {
	case 0:
		switch lineIdx % 3 {
		case 0:
			return lineIdx + 1
		case 1:
			return lineIdx - 1
		case 2:
			return lineIdx - 2
		default:
			return lineIdx
		}
	case 1:
		switch lineIdx % 3 {
		case 0:
			return lineIdx + 2
		case 1:
			return lineIdx + 1
		case 2:
			return lineIdx - 1
		default:
			return lineIdx
		}
	default:
		return lineIdx
	}
}

func sudokuString(s [][]int8) (out string) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			val := strconv.Itoa(int(s[row][col]))
			if val == "0" {
				val = "."
			}
			out += val
		}
	}
	return
}

// ASCII representation of the puzzle when debugging.
func sudokuDebug(s [][]int8) (out string) {
	out += "╔═══════╤═══════╤═══════╗\n"
	for i := 0; i < 9; i++ {
		out += "║ "
		for j := 0; j < 9; j++ {
			space := " "
			if j%3 == 2 && j != 8 {
				space = " │ "
			}
			value := strconv.Itoa(int(s[i][j]))
			if value == "0" {
				value = " "
			}
			out += fmt.Sprintf("%s%s", value, space)
		}
		out += fmt.Sprintf("║ %s\n", string('a'+byte(i)))
		if i%3 == 2 && i != 8 {
			out += "╟───────┼───────┼───────╢\n"
		}
	}
	out += "╚═══════╧═══════╧═══════╝\n"
	out += "  1 2 3   4 5 6   7 8 9  "
	return out
}

// Get all puzzle points randomly.
func sudokuRandomPoints(rnd *rand.Rand) []data.Point {
	var points []data.Point
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			points = append(points, data.Point{row, col})
		}
	}
	rnd.Shuffle(len(points), func(i, j int) {
		points[i], points[j] = points[j], points[i]
	})
	return points
}
