package data

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Sudoku interface {
	Board() SudokuBoard
	Puzzle() SudokuPuzzle
}

type SudokuBoard interface {
	String() string
}

type SudokuPuzzle interface {
	String() string
	FindUserErrors() []Point
	FindErrors(target SudokuPuzzle) []Point
	In(point Point) int8
}

// DirectionType is a direction of line/"big" line/some kind of field change.
type DirectionType uint8

const (
	Horizontal DirectionType = iota
	Vertical
)

// RotationType is an angle of rotation.
type RotationType uint8

const (
	Rotate0 RotationType = iota
	Rotate90
	Rotate180
	Rotate270
)

type Point struct {
	Row, Col int
}

func (p Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p Point) InSameBox(points ...Point) bool {
	boxRow, boxCol := p.Row/3, p.Col/3
	for _, ip := range points {
		if boxRow != ip.Row {
			return false
		}
		if boxCol != ip.Col {
			return false
		}
	}
	return true
}

func (p Point) InSameRow(points ...Point) bool {
	row := p.Row
	for _, ip := range points {
		if row != ip.Row {
			return false
		}
	}
	return true
}

func (p Point) InSameCol(points ...Point) bool {
	col := p.Col
	for _, ip := range points {
		if col != ip.Col {
			return false
		}
	}
	return true
}

func (p Point) String() string {
	return fmt.Sprintf("%s%d", string('a'+byte(p.Row)), p.Col+1)
}

func PointFromString(s string) (Point, error) {
	if len(s) < 2 {
		return Point{}, fmt.Errorf("unknown format Point")
	}
	p := Point{}

	switch ch := s[0]; {
	case 'a' <= ch && ch <= 'z':
		p.Row = int(ch) - 'a'
	case 'A' <= ch && ch <= 'Z':
		p.Row = int(ch) - 'A'
	default:
		return Point{}, fmt.Errorf("unknown format Point")
	}

	var err error
	p.Col, err = strconv.Atoi(s[1:])
	if err != nil {
		return Point{}, err
	}
	p.Col--

	return p, nil
}
