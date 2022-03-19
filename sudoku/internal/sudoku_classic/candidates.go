package sudoku_classic

import (
	"encoding/json"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"sort"
)

type sudokuCandidates [9][9]map[int8]struct{}

func newSudokuCandidates() (c sudokuCandidates) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			c[row][col] = make(map[int8]struct{})
		}
	}
	return c
}

func (c *sudokuCandidates) UnmarshalJSON(bts []byte) error {
	var in map[string][]int8
	if err := json.Unmarshal(bts, &in); err != nil {
		return err
	}
	*c = newSudokuCandidates()
	for pointStr, candidates := range in {
		p, err := data.PointFromString(pointStr)
		if err != nil {
			return fmt.Errorf("failed to parse point '%s': %v", pointStr, err)
		}
		for _, candidate := range candidates {
			(*c)[p.Row][p.Col][candidate] = struct{}{}
		}
	}
	return nil
}

func (c sudokuCandidates) MarshalJSON() ([]byte, error) {
	out := make(map[string][]int8)
	c.forEach(func(p data.Point, candidates []int8) {
		if len(candidates) > 0 {
			out[p.String()] = candidates
		}
	})
	return json.Marshal(out)
}

// todo break and excludes
func (c sudokuCandidates) forEach(fn func(p data.Point, candidates []int8)) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			fn(data.Point{row, col}, c.in(data.Point{row, col}))
		}
	}
}

// todo break and excludes
func (c sudokuCandidates) forEachInRow(row int, fn func(p data.Point, candidates []int8)) {
	for col := 0; col < 9; col++ {
		fn(data.Point{row, col}, c.in(data.Point{row, col}))
	}
}

// todo break and excludes
func (c sudokuCandidates) forEachInCol(col int, fn func(p data.Point, candidates []int8)) {
	for row := 0; row < 9; row++ {
		fn(data.Point{row, col}, c.in(data.Point{row, col}))
	}
}

// todo break and excludes
func (c sudokuCandidates) forEachInBox(p data.Point, fn func(p data.Point, candidates []int8)) {
	pBox := data.Point{(p.Row / 3) * 3, (p.Col / 3) * 3}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			pCurrent := data.Point{pBox.Row + row, pBox.Col + col}
			fn(pCurrent, c.in(pCurrent))
		}
	}
}

func (c sudokuCandidates) in(p data.Point) []int8 {
	out := make([]int8, 0, len(c[p.Row][p.Col]))
	for candidate, _ := range c[p.Row][p.Col] {
		out = append(out, candidate)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}
