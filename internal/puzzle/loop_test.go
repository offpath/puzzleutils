package puzzle

import (
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/tracker"
)

var slitherlinkTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name:  "Trivial 1x1",
		input: "4",
		want: `.-.
| |
.-.`,
	},
	{
		name:  "Empty 1x1",
		input: "0",
		want: `.X.
X X
.X.`,
	},
	{
		name: "Basic 5x5 (loop constraint required)",
		input: `...1.
32.2.
.22..
.223.
.22.3`,
		want: `.-.-.-.-.-.
| X X X X |
.X.-.-.X.-.
| | X | | X
.-.X.-.X.-.
X X | X X |
.X.-.X.-.X.
X | X | | |
.-.X.-.X.X.
| X | X | |
.-.-.X.X.-.`,
	},
}

func TestSlitherlink(t *testing.T) {
	for _, tt := range slitherlinkTests {
		slitherlink := NewSlitherlinkPuzzle(tt.input)
		if !slitherlink.Solve(csp.Settings{Decider: &decide.First{}, DecisionTracker: tracker.PrintEveryN(1)}) {
			t.Errorf("test: %s, failed to solve!\n", tt.name)
		}
		got := slitherlink.String()
		if got != tt.want {
			t.Errorf("test: %s, got: \n%s\n want: \n%s\n", tt.name, got, tt.want)
		}
	}
}
