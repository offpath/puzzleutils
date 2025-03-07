package puzzle

import "testing"

var logicTests = []struct {
	name  string
	input string
	want  string
}{
	{
		name: "basic",
		input: `Chef:Alice,Bob
Dish:Apple,Blueberry

eq(Chef(Apple), Bob)`,
		want: "",
	},
}

func TestLogic(t *testing.T) {
	for _, tt := range logicTests {
		NewLogicPuzzle(tt.input)
	}
}
