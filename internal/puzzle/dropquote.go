package puzzle

import (
	"strings"

	"github.com/offpath/puzzleutils/internal/constraints"
	"github.com/offpath/puzzleutils/internal/trie"
)

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func NewDropquotePuzzle(input string, t *trie.Trie) *Puzzle {
	lines := strings.Split(input, "\n")
	numCols := len(lines[0])
	numLetters := 0
	var cols [][]int
	var colSets []map[int]int
	var isCovering []bool
	for i := 0; i < numCols; i++ {
		cols = append(cols, nil)
	}
	var wordLengths []int
	currentLength := 0
	letterBanks := false
	for _, line := range lines {
		if line == "" && !letterBanks {
			letterBanks = true
			continue
		}
		if !letterBanks {
			for i, c := range line {
				if c == '.' {
					currentLength++
					cols[i] = append(cols[i], numLetters)
					numLetters++
				} else {
					if currentLength > 0 {
						wordLengths = append(wordLengths, currentLength)
						currentLength = 0
					}
				}
			}
		} else {
			s := map[int]int{}
			for _, c := range line {
				s[int(c)-int('A')]++
			}
			colSets = append(colSets, s)
			isCovering = append(isCovering, len(line) == len(cols[len(colSets)-1]))
			if len(colSets) == len(cols) {
				break
			}
		}
	}
	if currentLength > 0 {
		wordLengths = append(wordLengths, currentLength)
	}
	result := NewPuzzle(numLetters, alphabet)
	offset := 0
	for _, length := range wordLengths {
		var group []int
		for j := 0; j < length; j++ {
			group = append(group, offset+j)
		}
		result.problem.AddGroup(group, constraints.ValidWord(t, alphabet))
		offset += length
	}
	for i := 0; i < numCols; i++ {
		result.problem.AddGroup(cols[i], constraints.SetCount(colSets[i], isCovering[i]))
	}
	return result
}
