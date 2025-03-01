package puzzle

import "strings"

// Example:
//
// Chef:Freda,Karl,Sonia,Wade
// Dish:cashew tofu,lemon snapper,smoked pork,turkey soup
// Score:42,49,56,63
//
// Dish(Sonia) = cashew tofu
// Score(Karl) = Chef(smoked pork) + 7
// Score(Freda) = 49
// Score(Chef(turkey soup)) = Score(Sonia) - 7
//
// Answers:
// 42, Wade, lemon snapper
// 49, Freda, smoked pork
// 56, Karl, turkey soup
// 63, Sonia, cashew tofu

type category struct {
	name   string
	values []string
}

type LogicPuzzle struct {
	*Puzzle
	categories map[string]*category
}

func NewLogicPuzzle(s string) *LogicPuzzle {
	lines := strings.Split(s, "\n")
	result := &LogicPuzzle{}
	for _, line := range lines {
		if line == "" {
			break
		}
		parts := strings.Split(line, ":")
		result.categories[parts[0]] = &category{parts[0], strings.Split(parts[1], ",")}
	}
	return result
}
