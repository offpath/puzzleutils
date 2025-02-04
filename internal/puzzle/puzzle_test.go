package puzzle

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/tracker"
	"github.com/offpath/puzzleutils/internal/trie"
)

func TestNonogram(t *testing.T) {
	rows := [][]int{
		{8, 7, 5, 7},
		{5, 4, 3, 3},
		{3, 3, 2, 3},
		{4, 3, 2, 2},
		{3, 3, 2, 2},
		{3, 4, 2, 2},
		{4, 5, 2},
		{3, 5, 1},
		{4, 3, 2},
		{3, 4, 2},
		{4, 4, 2},
		{3, 6, 2},
		{3, 2, 3, 1},
		{4, 3, 4, 2},
		{3, 2, 3, 2},
		{6, 5},
		{4, 5},
		{3, 3},
		{3, 3},
		{1, 1},
	}
	cols := [][]int{
		{1},
		{1},
		{2},
		{4},
		{7},
		{9},
		{2, 8},
		{1, 8},
		{8},
		{1, 9},
		{2, 7},
		{3, 4},
		{6, 4},
		{8, 5},
		{1, 11},
		{1, 7},
		{8},
		{1, 4, 8},
		{6, 8},
		{4, 7},
		{2, 4},
		{1, 4},
		{5},
		{1, 4},
		{1, 5},
		{7},
		{5},
		{3},
		{1},
		{1},
	}
	nonogram := NewNonogramPuzzle(rows, cols)
	nonogram.problem.Solve(csp.Settings{Decider: &decide.First{}})
	got := nonogram.ToString()
	want := `XXXXXXXX.XXXXXXX.XXXXX.XXXXXXX
..XXXXX...XXXX....XXX....XXX..
...XXX.....XXX....XX.....XXX..
...XXXX.....XXX...XX.....XX...
....XXX.....XXX..XX......XX...
....XXX.....XXXX.XX.....XX....
....XXXX.....XXXXX......XX....
.....XXX.....XXXXX......X.....
.....XXXX.....XXX......XX.....
......XXX.....XXXX.....XX.....
......XXXX....XXXX....XX......
.......XXX...XXXXXX...XX......
.......XXX...XX.XXX...X.......
.......XXXX.XXX.XXXX.XX.......
........XXX.XX...XXX.XX.......
........XXXXXX...XXXXX........
.........XXXX....XXXXX........
.........XXX......XXX.........
.........XXX......XXX.........
..........X........X..........
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s\n", got, want)
	}
}

var dropquoteTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name: "trivial",
		input: `.....

H
E
L
L
O`,
		want: "HELLO",
	},
	{
		name: "mostly trivial",
		input: `.....X.....

H
E
L
L
O

W
O
R
L
D`,
		want: "HELLOWORLD",
	},
	{
		name: "short",
		input: `...x
....

TB
HE
ES
T`,
		want: "THEBEST",
	},
	{
		name: "long",
		input: `...X.....X..X.......
..X..X...X....X..X..
....X...X......X...X
.....XXXXXXXXXXXXXXX

OEET
DXHP
EIL
IES
WTN
OW
TRH
OHL
DE
S
IMH
SAO
IU
NPL
DO
PB
UYN
OL
TPA
TE
`,
		want: "THEWORLDISPOPULATEDINTHEMAINBYPEOPLEWHOSHOULDNOTEXIST",
	},
}

func TestDropquote(t *testing.T) {
	tr := trie.New()
	tr.Add("hello")
	tr.Add("world")
	tr.Add("the")
	tr.Add("best")
	tr.Add("is")
	tr.Add("populated")
	tr.Add("in")
	tr.Add("main")
	tr.Add("by")
	tr.Add("people")
	tr.Add("who")
	tr.Add("should")
	tr.Add("not")
	tr.Add("exist")
	// TODO(dneal): tr.AddFile("")
	for _, tt := range dropquoteTests {
		dropquote := NewDropquotePuzzle(tt.input, tr)
		if !dropquote.Solve(csp.Settings{Decider: &decide.First{}, DecisionTracker: tracker.PrintEveryN(1)}) {
			t.Errorf("test: %s, failed to solve!\n", tt.name)
		}
		got := dropquote.ToString()
		if got != tt.want {
			t.Errorf("test: %s, got: %s, want: %s\n", tt.name, got, tt.want)
		}
	}
}
