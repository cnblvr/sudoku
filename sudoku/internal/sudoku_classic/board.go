package sudoku_classic

import (
	"github.com/cnblvr/sudoku/data"
)

// sudokuBoard is 9 lines with 9 digits on each line. A correct puzzle satisfies the condition that each column, each
// row, and each of the nine box 3x3 contain all of the digits from 1 to 9.
//  ╔═══════╤═══════╤═══════╗
//  ║ 3 1 6 │ 5 7 9 │ 2 4 8 ║ a
//  ║ 5 7 9 │ 2 4 8 │ 3 1 6 ║ b
//  ║ 2 4 8 │ 3 1 6 │ 5 7 9 ║ c
//  ╟───────┼───────┼───────╢
//  ║ 4 8 3 │ 1 6 5 │ 7 9 2 ║ d
//  ║ 1 6 5 │ 7 9 2 │ 4 8 3 ║ e
//  ║ 7 9 2 │ 4 8 3 │ 1 6 5 ║ f
//  ╟───────┼───────┼───────╢
//  ║ 9 2 4 │ 8 3 1 │ 6 5 7 ║ g
//  ║ 8 3 1 │ 6 5 7 │ 9 2 4 ║ h
//  ║ 6 5 7 │ 9 2 4 │ 8 3 1 ║ i
//  ╚═══════╧═══════╧═══════╝
//    1 2 3   4 5 6   7 8 9
// Context of methods:
//  - the lines a-i are data.Horizontal lines;
//  - the lines 1-9 are data.Vertical lines;
//  - the lines [1-3], [4-6], [7-9], [a-c], [d-f], [g-i] are "big" lines;
//  - box 3x3 is a matrix with 3 rows and 3 columns, for example:
//    [[a1,a2,a3],[b1,b2,b3],[c1,c2,c3]].
type sudokuBoard [][]int8

// Swap of lines within one "big" line.
func (b sudokuBoard) swapLines(dir data.DirectionType, lineA, lineB int) {
	switch dir {

	case data.Horizontal:
		// Swap two horizontal lines
		temp := make([]int8, 9)
		copy(temp, b[lineA])
		copy(b[lineA], b[lineB])
		copy(b[lineB], temp)

	case data.Vertical:
		// Swap two vertical lines
		for row := 0; row < 9; row++ {
			b[row][lineA], b[row][lineB] = b[row][lineB], b[row][lineA]
		}
	}
}

// Reflect the entire puzzle in the typ direction.
func (b sudokuBoard) reflect(typ data.DirectionType) {
	switch typ {

	case data.Horizontal:
		for col := 0; col < 4; col++ {
			for row := 0; row < 9; row++ {
				b[row][col], b[row][8-col] = b[row][8-col], b[row][col]
			}
		}

	case data.Vertical:
		for row := 0; row < 4; row++ {
			b[row], b[8-row] = b[8-row], b[row]
		}
	}
}

// Rotate the entire puzzle in the rotation r.
func (b sudokuBoard) rotate(r data.RotationType) {
	switch r {

	case data.Rotate90:
		temp := make([][]int8, 9)
		for col := 8; col >= 0; col-- {
			for row := 0; row < 9; row++ {
				temp[8-col] = append(temp[8-col], b[row][col])
			}
		}
		for row := 0; row < 9; row++ {
			b[row] = temp[row]
		}

	case data.Rotate180:
		for row := 0; row < 4; row++ {
			for col := 0; col < 9; col++ {
				b[row][col], b[8-row][8-col] = b[8-row][8-col], b[row][col]
			}
		}
		for col := 0; col < 4; col++ {
			b[4][col], b[4][8-col] = b[4][8-col], b[4][col]
		}

	case data.Rotate270:
		temp := make([][]int8, 9)
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
				b[row][col] = int8(digitB)
			case digitB:
				b[row][col] = int8(digitA)
			}
		}
	}
}

// Get the puzzle solution as a string of 81 characters.
func (b sudokuBoard) String() string {
	return sudokuString(b)
}

func (b sudokuBoard) debug() string {
	return sudokuDebug(b)
}
