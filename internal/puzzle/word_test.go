package puzzle

import (
	"path/filepath"
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/tracker"
	"github.com/offpath/puzzleutils/internal/trie"
)

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
	{
		name: "extra letters",
		input: `...X.....X..X.......
..X..X...X....X..X..
....X...X......X...X
.....XXXXXXXXXXXXXXX

OEETA
DXHPQ
EILZ
IES
WTNN
OW
TRHP
OHL
DE
SR
IMH
SAO
IUL
NPL
DO
PB
UYNQ
OL
TPA
TE
`,
		want: "THEWORLDISPOPULATEDINTHEMAINBYPEOPLEWHOSHOULDNOTEXIST",
	},
}

func TestDropquote(t *testing.T) {
	tr := trie.New()
	tr.AddFile(filepath.Join("testdata", "ospd2.txt"))
	for _, tt := range dropquoteTests {
		dropquote := NewDropquotePuzzle(tt.input, tr)
		if !dropquote.Solve(csp.Settings{Decider: &decide.First{}, DecisionTracker: tracker.PrintEveryN(1)}) {
			t.Errorf("test: %s, failed to solve!\n", tt.name)
		}
		got := dropquote.String()
		if got != tt.want {
			t.Errorf("test: %s, got: %s, want: %s\n", tt.name, got, tt.want)
		}
	}
}

var cryptegramTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name:  "trivial",
		input: "ABBCCEEFEG",
		want:  "BOKEPR",
	},
}

func TestCryptogram(t *testing.T) {
	tr := trie.New()
	tr.AddFile(filepath.Join("testdata", "ospd2.txt"))
	for _, tt := range cryptegramTests {
		cryptogram := NewCryptogramPuzzle(tt.input, tr)
		if !cryptogram.Solve(csp.Settings{Decider: &decide.First{}, DecisionTracker: tracker.PrintEveryN(1)}) {
			t.Errorf("test: %s, failed to solve!\n", tt.name)
		}
		got := cryptogram.String()
		if got != tt.want {
			t.Errorf("test: %s, got: %s, want: %s\n", tt.name, got, tt.want)
		}
	}
}
