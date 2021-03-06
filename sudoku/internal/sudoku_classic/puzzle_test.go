package sudoku_classic

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_sudokuPuzzle_findCandidates(t *testing.T) {
	tests := []struct {
		name string
		p    sudokuPuzzle
		want string
	}{
		{
			name: "#1",
			p:    sudokuPuzzleFromString("400000938032094100095300240370609004529001673604703090957008300003900400240030709"),
			want: `{"a2":[1,6],"a3":[1,6],"a4":[1,2,5],"a5":[1,2,5,6,7],"a6":[2,5,6,7],"b1":[7,8],"b4":[5,8],"b8":[5,6],"b9":[5,6,7],"c1":[1,7,8],"c5":[1,6,7,8],"c6":[6,7],"c9":[6,7],"d3":[1,8],"d5":[2,5,8],"d7":[5,8],"d8":[1,2,5,8],"e4":[4,8],"e5":[4,8],"f2":[1,8],"f5":[2,5,8],"f7":[5,8],"f9":[1,2,5],"g4":[1,2,4],"g5":[1,2,4,6],"g8":[1,2,6],"g9":[1,2,6],"h1":[1,8],"h2":[1,6,8],"h5":[1,2,5,6,7],"h6":[2,5,6,7],"h8":[1,2,5,6,8],"h9":[1,2,5,6],"i3":[1,6,8],"i4":[1,5],"i6":[5,6],"i8":[1,5,6,8]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.findCandidates()
			gotBts, err := json.Marshal(got)
			if err != nil {
				t.Errorf("failed to marshal candidates: %v", err)
			}
			if string(gotBts) != tt.want {
				t.Errorf("findCandidates() =\n%s,\nwant\n%s", gotBts, tt.want)
			}
		})
	}
}

func Test_sudokuPuzzle_solveBruteForce(t *testing.T) {
	tests := []struct {
		name     string
		p        string
		breakOn  int
		want     []string
		duration bool
	}{
		{
			name: "#1",
			p:    "...1.5...14....67..8...24...63.7..1.9.......3.1..9.52...72...8..26....35...4.9...",
			want: []string{
				"672145398145983672389762451263574819958621743714398526597236184426817935831459267",
			},
		},
		{
			name: "#2",
			p:    ".....5...14....67..8...24...63.7..1.9.......3.1..9.52...72...8..26....35...4.9...",
			want: []string{
				"672145398145983672389762451263574819958621743714398526597236184426817935831459267",
			},
		},
		{
			name: "#3",
			p:    ".........14....67..8...24...63.7..1.9.......3.1..9.52...72...8..26....35...4.9...",
			want: []string{
				"672145398145983672389762451263574819958621743714398526597236184426817935831459267",
			},
		},
		{
			name: "#4",
			p:    "..........4....67..8...24...63.7..1.9.......3.1..9.52...72...8..26....35...4.9...",
			want: []string{
				"172645398345981672689732451263574819958126743714398526597263184426817935831459267",
				"172946358349185672685732491263574819958621743714398526597263184426817935831459267",
				"172946358349815672685732491263574819958621743714398526597263184426187935831459267",
				"175846392342915678689732451263574819958621743714398526597263184426187935831459267",
				"179846352342915678685732491263574819958621743714398526597263184426187935831459267",
				"672145398145983672389762451263574819958621743714398526597236184426817935831459267",
			},
			duration: true,
		},
		{
			name:    "#4 break",
			p:       "..........4....67..8...24...63.7..1.9.......3.1..9.52...72...8..26....35...4.9...",
			breakOn: 3,
			want: []string{
				"172645398345981672689732451263574819958126743714398526597263184426817935831459267",
				"172946358349185672685732491263574819958621743714398526597263184426817935831459267",
				"172946358349815672685732491263574819958621743714398526597263184426187935831459267",
			},
			duration: true,
		},
		{
			name:    "#5",
			p:       "000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			breakOn: 2,
			want: []string{
				"123456789456789123789123456214365897365897214897214365531642978642978531978531642",
				"123456789456789123789123456214365897365897214897214365531642978648971532972538641",
			},
			duration: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			solutions := sudokuPuzzleFromString(tt.p).solveBruteForce(tt.breakOn)
			if tt.duration {
				t.Logf("solve time: %s", time.Since(start).Truncate(time.Microsecond).String())
			}
			solutionsMap := make(map[string]struct{})
			for _, s := range solutions {
				solutionsMap[s.String()] = struct{}{}
			}
			if len(solutionsMap) != len(solutions) {
				t.Errorf("got(len=%d) contains same solutions", len(solutions))
			}
			if len(solutionsMap) != len(tt.want) {
				t.Errorf("want num solutions = %d, got = %d", len(tt.want), len(solutionsMap))
			}
			for _, w := range tt.want {
				if _, ok := solutionsMap[w]; !ok {
					t.Errorf("want solution\n%s\nnot found in got solutions", w)
				}
			}
			if t.Failed() {
				t.Logf("got solutions:\n%v", solutions)
			}
		})
	}
}

func Test_sudokuPuzzle_isSolve(t *testing.T) {
	tests := []struct {
		name string
		p    string
		want bool
	}{
		{
			name: "#1",
			p:    "400000938032094100095300240370609004529001673604703090957008300003900400240030709",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sudokuPuzzleFromString(tt.p).isCorrectSolve(); got != tt.want {
				t.Errorf("isSolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
