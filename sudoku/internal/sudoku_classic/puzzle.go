package sudoku_classic

import (
	"github.com/cnblvr/sudoku/data"
)

type sudokuPuzzle [][]int8

func sudokuPuzzleFromString(str string) sudokuPuzzle {
	out := make(sudokuPuzzle, 9)
	for row := 0; row < 9; row++ {
		out[row] = make([]int8, 9)
	}
	for idx, ch := range str {
		if idx >= 81 {
			return out
		}
		switch {
		case '0' <= ch && ch <= '9':
			out[idx/9][idx%9] = int8(byte(ch) - '0')
		default:
			out[idx/9][idx%9] = 0
		}

	}
	return out
}

func (p sudokuPuzzle) solveBruteForce(breakOn int) []sudokuPuzzle {
	var solutions []sudokuPuzzle
	var recursion func(p sudokuPuzzle)
	recursion = func(p sudokuPuzzle) {
		if breakOn > 0 && breakOn <= len(solutions) {
			return
		}
		candidates := p.findCandidates()
		for row := 0; row < 9; row++ {
			for col := 0; col < 9; col++ {
				if p[row][col] > 0 {
					continue
				}
				point := data.Point{Row: row, Col: col}
				for _, c := range candidates.in(point) {
					p[row][col] = c
					recursion(p)
					p[row][col] = 0
				}
				return
			}
		}
		if !p.isCorrectSolve() {
			return
		}
		solutions = append(solutions, p.clone())
		return
	}

	puzzle := p.clone()
	recursion(puzzle)

	return solutions
}

func (p sudokuPuzzle) forEach(fn func(p data.Point, v int8, _break *bool), excludeCols ...int) {
	excludes := make(map[int]struct{})
	for _, e := range excludeCols {
		excludes[e] = struct{}{}
	}
	_break := false
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if _, isExcluded := excludes[col]; isExcluded {
				continue
			}
			fn(data.Point{Row: row, Col: col}, p[row][col], &_break)
		}
	}
}

func (p sudokuPuzzle) forEachInRow(row int, fn func(p data.Point, v int8, _break *bool), excludeCols ...int) {
	excludes := make(map[int]struct{})
	for _, e := range excludeCols {
		excludes[e] = struct{}{}
	}
	_break := false
	for col := 0; !_break && col < 9; col++ {
		if _, isExcluded := excludes[col]; isExcluded {
			continue
		}
		fn(data.Point{row, col}, p[row][col], &_break)
	}
}

func (p sudokuPuzzle) forEachInCol(col int, fn func(p data.Point, v int8, _break *bool), excludeRows ...int) {
	excludes := make(map[int]struct{})
	for _, e := range excludeRows {
		excludes[e] = struct{}{}
	}
	_break := false
	for row := 0; !_break && row < 9; row++ {
		if _, isExcluded := excludes[row]; isExcluded {
			continue
		}
		fn(data.Point{row, col}, p[row][col], &_break)
	}
}

func (p sudokuPuzzle) forEachInBox(point data.Point, fn func(p data.Point, v int8, _break *bool), excludePoints ...data.Point) {
	excludes := make(map[data.Point]struct{})
	for _, e := range excludePoints {
		excludes[e] = struct{}{}
	}
	_break := false
	pBox := data.Point{(point.Row / 3) * 3, (point.Col / 3) * 3}
	for row := 0; !_break && row < 3; row++ {
		for col := 0; !_break && col < 3; col++ {
			pCurrent := data.Point{pBox.Row + row, pBox.Col + col}
			if _, isExcluded := excludes[pCurrent]; isExcluded {
				continue
			}
			fn(pCurrent, p[pCurrent.Row][pCurrent.Col], &_break)
		}
	}
}

func (p sudokuPuzzle) findCandidates() sudokuCandidates {
	c := newSudokuCandidates()
	p.forEach(func(p1 data.Point, v1 int8, _ *bool) {
		if v1 > 0 {
			return
		}
		for i := int8(1); i <= 9; i++ {
			c[p1.Row][p1.Col][i] = struct{}{}
		}
	})
	// delete from vertical and horizontal lines and from boxes 3x3
	p.forEach(func(p1 data.Point, v1 int8, _ *bool) {
		if v1 == 0 {
			return
		}
		boxRow, boxCol := p1.Row/3*3, p1.Col/3*3
		for i := 0; i < 9; i++ {
			delete(c[p1.Row][i], v1)
			delete(c[i][p1.Col], v1)
			delete(c[boxRow+i%3][boxCol+i/3], v1)
		}
	})
	return c
}

func (p sudokuPuzzle) isCorrectSolve() bool {
	isCorrect := true
	p.forEach(func(p1 data.Point, v1 int8, _break1 *bool) {
		if p[p1.Row][p1.Col] == 0 {
			isCorrect = false
			*_break1 = true
			return
		}
		p.forEachInRow(p1.Row, func(p2 data.Point, v2 int8, _break2 *bool) {
			if v1 == v2 {
				isCorrect = false
				*_break1, *_break2 = true, true
				return
			}
		}, p1.Col)
		p.forEachInCol(p1.Col, func(p2 data.Point, v2 int8, _break2 *bool) {
			if v1 == v2 {
				isCorrect = false
				*_break1, *_break2 = true, true
				return
			}
		}, p1.Row)
		p.forEachInBox(p1, func(p2 data.Point, v2 int8, _break2 *bool) {
			if v1 == v2 {
				isCorrect = false
				*_break1, *_break2 = true, true
				return
			}
		}, p1)
	})
	return isCorrect
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

func (p sudokuPuzzle) clone() sudokuPuzzle {
	c := make([][]int8, len(p))
	for row := 0; row < len(p); row++ {
		c[row] = make([]int8, len(p[row]))
		copy(c[row], p[row])
	}
	return c
}

// Get the puzzle for user as a string of 81 characters, where the absence of a number is indicated by zero.
func (p sudokuPuzzle) String() string {
	return sudokuString(p)
}

func (p sudokuPuzzle) debug() string {
	return sudokuDebug(p)
}
