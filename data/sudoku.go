package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

var (
	ErrSudokuNotFound = fmt.Errorf("sudoku not found")
)

type SudokuRepository interface {
	CreateSudoku(ctx context.Context, typ SudokuType, seed int64, puzzle, solution string) (*Sudoku, error)
	GetSudokuByID(ctx context.Context, id int64) (*Sudoku, error)
}

type Sudoku struct {
	IDint64Getter   `json:"id"`
	Type            SudokuType `json:"type"`
	Seed            int64      `json:"seed"`
	Puzzle          string     `json:"puzzle"`
	Solution        string     `json:"solution"`
	CreatedAtGetter `json:"created_at"`
}

type SudokuGame struct {
	IDuuidGetter `json:"id"`
}

type SudokuType string

const (
	SudokuClassic SudokuType = "sudoku_classic"
)

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
