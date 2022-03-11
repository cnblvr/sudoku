package sudoku

import (
	"fmt"
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

type (
	sudokuBoard  [][]uint8
	sudokuPuzzle [][]uint8
)

// NewSudoku creates a new puzzle and removes some hints depending on the level.
// seed is used to create a unique puzzle.
func NewSudoku(seed int64, level Level) *Sudoku {
	s := Sudoku{}
	s.seed = seed
	// randomizer for full puzzle generation
	rnd := rand.New(rand.NewSource(seed))

	// puzzle generation without shuffling
	s.board = generateSudokuBoard(rnd)

	// swap of horizontal or vertical lines within one "big" line
	// TODO: imperfect randomization
	for i := 0; i < (rnd.Int()%1024)+1024; i++ {
		typ := horizontal
		if rnd.Int()%2 == 1 {
			typ = vertical
		}
		line := rnd.Int() % 9
		s.board.swapLines(typ, line, neighborLine(line, rnd.Int()%2))
	}

	// TODO: swap "big" lines

	// horizontal reflection
	if rnd.Int()%2 == 1 {
		s.board.reflect(horizontal)
	}
	// vertical reflection
	if rnd.Int()%2 == 1 {
		s.board.reflect(vertical)
	}

	// rotate the puzzle by a random angle
	s.board.rotate(rotationType(rnd.Int() % 4))

	s.puzzle = make([][]uint8, 9)
	for row := 0; row < 9; row++ {
		s.puzzle[row] = make([]uint8, 9)
		copy(s.puzzle[row], s.board[row])
	}
	// remove hints based on given level
	if level.hasStrategy(strategyNakedSingle) {
		s.puzzle.nakedSingle(rnd, level)
	}

	return &s
}

// Level of the game is based on the strategies used in solving.
// Strategies for solving are given by a bitmask.
// TODO: add other strategies https://www.sudokuwiki.org/Strategy_Families
type Level uint64

const (
	Beginner = Level(strategyNakedSingle)
	Easy     = Level(strategyNakedSingle | strategyNakedPairs)
)

// Checks if the level contains the s strategy.
func (l Level) hasStrategy(s sudokuStrategy) bool {
	return sudokuStrategy(l)&s != 0
}

type sudokuStrategy uint64

const (
	// https://www.sudokuwiki.org/Getting_Started
	strategyNakedSingle sudokuStrategy = 1 << iota
	// https://www.sudokuwiki.org/Naked_Candidates#NP
	strategyNakedPairs
)

// Puzzle generation without shuffling.
func generateSudokuBoard(rnd *rand.Rand) sudokuBoard {
	b := make([][]uint8, 9)
	for i := 0; i < 9; i++ {
		b[i] = make([]uint8, 9)
	}
	// Generate first line randomly
	digits := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}
	line := make([]uint8, 0, 9)
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

// Direction of line/"big" line/some kind of field change.
type lineType uint8

const (
	horizontal lineType = iota
	vertical
)

// Swap of lines within one "big" line.
func (b sudokuBoard) swapLines(typ lineType, lineA, lineB int) {
	switch typ {

	case horizontal:
		// Swap two horizontal lines
		temp := make([]uint8, 9)
		copy(temp, b[lineA])
		copy(b[lineA], b[lineB])
		copy(b[lineB], temp)

	case vertical:
		// Swap two vertical lines
		for row := 0; row < 9; row++ {
			b[row][lineA], b[row][lineB] = b[row][lineB], b[row][lineA]
		}
	}
}

// Reflect the entire puzzle in the typ direction.
func (b sudokuBoard) reflect(typ lineType) {
	switch typ {

	case horizontal:
		for col := 0; col < 4; col++ {
			for row := 0; row < 9; row++ {
				b[row][col], b[row][8-col] = b[row][8-col], b[row][col]
			}
		}

	case vertical:
		for row := 0; row < 4; row++ {
			b[row], b[8-row] = b[8-row], b[row]
		}
	}
}

// Angle of rotation.
type rotationType uint8

const (
	rotate0 rotationType = iota
	rotate90
	rotate180
	rotate270
)

// Rotate the entire puzzle in the rotation r.
func (b sudokuBoard) rotate(r rotationType) {
	switch r {

	case rotate90:
		temp := make([][]uint8, 9)
		for col := 8; col >= 0; col-- {
			for row := 0; row < 9; row++ {
				temp[8-col] = append(temp[8-col], b[row][col])
			}
		}
		for row := 0; row < 9; row++ {
			b[row] = temp[row]
		}

	case rotate180:
		for row := 0; row < 4; row++ {
			for col := 0; col < 9; col++ {
				b[row][col], b[8-row][8-col] = b[8-row][8-col], b[row][col]
			}
		}
		for col := 0; col < 4; col++ {
			b[4][col], b[4][8-col] = b[4][8-col], b[4][col]
		}

	case rotate270:
		temp := make([][]uint8, 9)
		for col := 0; col < 9; col++ {
			for row := 8; row >= 0; row-- {
				temp[col] = append(temp[col], b[row][col])
			}
		}
		for row := 0; row < 9; row++ {
			b[row] = temp[row]
		}
	}
}

// Swap the digits digitA and digitB in the entire puzzle.
func (b sudokuBoard) swapDigits(digitA, digitB int) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			switch int(b[row][col]) {
			case digitA:
				b[row][col] = uint8(digitB)
			case digitB:
				b[row][col] = uint8(digitA)
			}
		}
	}
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

// Get the puzzle solution as a string of 81 characters.
func (b sudokuBoard) String() string {
	return sudokuString(b)
}

func (b sudokuBoard) debug() string {
	return sudokuDebug(b)
}

// Delete easy hints.
// https://www.sudokuwiki.org/Getting_Started
func (p sudokuPuzzle) nakedSingle(rnd *rand.Rand, level Level) {
	points := sudokuPoints(rnd)
	for _, point := range points {
		digit := p[point.row][point.col]
		if digit == 0 {
			continue
		}
		p[point.row][point.col] = 0
		guesses := p.searchGuesses(level)
		if _, isGuessed := guesses[point.row][point.col][digit]; len(guesses[point.row][point.col]) != 1 || !isGuessed {
			p[point.row][point.col] = digit
			continue
		}
	}
}

type sudokuGuesses [][]sudokuCellGuesses

type sudokuCellGuesses map[uint8]struct{}

// Generation of puzzle guesses in the current state.
func (p sudokuPuzzle) searchGuesses(level Level) sudokuGuesses {
	guesses := make(sudokuGuesses, 9)
	for row := 0; row < 9; row++ {
		guesses[row] = make([]sudokuCellGuesses, 9)
		for col := 0; col < 9; col++ {
			guesses[row][col] = make(sudokuCellGuesses)

			if p[row][col] != 0 {
				continue
			}
		guessFor:
			for guess := uint8(1); guess <= 9; guess++ {
				// check horizontal line
				for i := 0; i < 9; i++ {
					if p[row][i] == guess {
						continue guessFor
					}
				}
				// check vertical line
				for i := 0; i < 9; i++ {
					if p[i][col] == guess {
						continue guessFor
					}
				}
				// check 3x3 box
				boxRow, boxCol := row/3, col/3
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if p[boxRow*3+i][boxCol*3+j] == guess {
							continue guessFor
						}
					}
				}
				guesses[row][col][guess] = struct{}{}
			}
		}
	}

	// Filtering guesses based on the level of the game.
	switch {
	case level.hasStrategy(strategyNakedPairs):
		guesses.nakedPairs()
	}

	return guesses
}

// Filtering guesses by the strategy Naked Pairs.
func (g sudokuGuesses) nakedPairs() {
	equalPairs := func(p1, p2 sudokuCellGuesses) bool {
		for g := range p1 {
			if _, ok := p2[g]; !ok {
				return false
			}
		}
		return true
	}

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if len(g[row][col]) != 2 {
				continue
			}
			// horizontally search a pair
			for i := 0; i < 9; i++ {
				if i == col || len(g[row][i]) != 2 {
					continue
				}
				if equalPairs(g[row][col], g[row][i]) {
					for j := 0; j < 9; j++ {
						for del := range g[row][col] {
							if _, ok := g[row][j][del]; ok {
								delete(g[row][j], del)
							}
						}
					}
				}
			}
			// vertically search a pair
			for i := 0; i < 9; i++ {
				if i == row || len(g[i][col]) != 2 {
					continue
				}
				if equalPairs(g[row][col], g[i][col]) {
					for j := 0; j < 9; j++ {
						for del := range g[row][col] {
							if _, ok := g[j][col][del]; ok {
								delete(g[j][col], del)
							}
						}
					}
				}
			}
			// search a pair in a box
			boxRow, boxCol := row/3, col/3
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if (boxRow*3+i == row && boxCol*3+j == col) || len(g[boxRow*3+i][boxCol*3+j]) != 2 {
						continue
					}
					if equalPairs(g[row][col], g[boxRow*3+i][boxCol*3+j]) {
						for ii := 0; ii < 3; ii++ {
							for jj := 0; jj < 3; jj++ {
								for del := range g[row][col] {
									if _, ok := g[boxRow*3+ii][boxCol*3+jj][del]; ok {
										delete(g[boxRow*3+ii][boxCol*3+jj], del)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

// CountHints returns the number of hints in the current state.
func (p sudokuPuzzle) CountHints() (c int) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if p[row][col] != 0 {
				c++
			}
		}
	}
	return
}

// Get the puzzle for user as a string of 81 characters, where the absence of a number is indicated by zero.
func (p sudokuPuzzle) String() string {
	return sudokuString(p)
}

func (p sudokuPuzzle) debug() string {
	return sudokuDebug(p)
}

func sudokuString(s [][]uint8) (out string) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			out += strconv.Itoa(int(s[row][col]))
		}
	}
	return
}

// ASCII representation of the puzzle when debugging.
func sudokuDebug(s [][]uint8) (out string) {
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
		out += "║\n"
		if i%3 == 2 && i != 8 {
			out += "║───────┼───────┼───────║\n"
		}
	}
	out += "╚═══════╧═══════╧═══════╝"
	return out
}

type sudokuPoint struct {
	row, col int
}

// Get all puzzle points randomly.
func sudokuPoints(rnd *rand.Rand) []sudokuPoint {
	var points []sudokuPoint
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			points = append(points, sudokuPoint{row, col})
		}
	}
	rnd.Shuffle(len(points), func(i, j int) {
		points[i], points[j] = points[j], points[i]
	})
	return points
}
