package puzzletypes

import (
	"strings"

	"github.com/offpath/puzzleutils/internal/constraints2"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/trie"
)

type CryptogramPuzzle struct {
	Words [][]*puzzle.Variable
	Codex map[string]*puzzle.Variable
}

func NewCryptogramPuzzle(p *puzzle.Puzzle2, input string, t *trie.Trie) *CryptogramPuzzle {
	result := &CryptogramPuzzle{
		Codex: make(map[string]*puzzle.Variable),
	}
	words := strings.Split(input, " ")

	alphabet := p.GetAlphabet()

	for _, w := range words {
		var word []*puzzle.Variable
		for _, c := range w {
			charStr := string(c)
			var v *puzzle.Variable
			var ok bool
			if v, ok = result.Codex[charStr]; !ok {
				v = p.NewVariable()
				v.SetValueRange(alphabet)
				result.Codex[charStr] = v
			}
			word = append(word, v)
		}
		p.AddConstraint(constraints2.NewValidWordConstraint(p, word, t, alphabet))
		result.Words = append(result.Words, word)
	}

	var allVars []*puzzle.Variable
	for _, v := range result.Codex {
		allVars = append(allVars, v)
	}
	p.AddConstraint(constraints2.NewUniqueConstraint(allVars, alphabet))

	return result
}
