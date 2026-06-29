package puzzletypes

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/tracker"
	"github.com/offpath/puzzleutils/internal/trie"
)

var cryptogramTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name:  "trivial",
		input: "ABBCCEEFEG",
		want:  "BOOKKEEPER",
	},
}

func TestCryptogram(t *testing.T) {
	tr := trie.New()
	tr.AddFile(filepath.Join("testdata", "ospd2.txt"))
	for _, tt := range cryptogramTests {
		p := puzzle.NewPuzzle()
		result := NewCryptogramPuzzle(p, tt.input, tr)
		if !p.Solve(csp.Settings{Decider: &decide.First{}, DecisionTracker: tracker.PrintEveryN(1)}) {
			t.Errorf("test: %s, failed to solve!\n", tt.name)
			continue
		}

		var gotWords []string
		for _, group := range result.Words {
			var wordGot string
			for _, v := range group {
				if val := v.Values().Value(); val != nil {
					wordGot += val.Str()
				} else {
					wordGot += "?"
				}
			}
			gotWords = append(gotWords, wordGot)
		}
		got := strings.Join(gotWords, " ")

		if got != tt.want {
			t.Errorf("test: %s, got: %s, want: %s\n", tt.name, got, tt.want)
		}
	}
}
