package main

import (
	"fmt"

	"puzzleutils/internal/csp"
	"puzzleutils/internal/puzzle"
)

type Settings struct {
	count int
	multiple int
}

func (s *Settings) MakeDecision(p *csp.Problem) {
	s.count++
	if s.count % s.multiple == 0 {
		fmt.Printf("Count = %d\n", s.count)
	}
}

func (s *Settings) CaptureSolution(p *csp.Problem) {
	p.Print()
}

func (s *Settings) Decide(d []*csp.Decision) *csp.Decision {
	min := d[0]
	for i := 1; i < len(d); i++ {
		if d[i].Count() < min.Count() {
			min = d[i]
		}
	}
	return min
}

func main() {
	fmt.Println("Hello World")
	p := puzzle.NewGridPuzzle(3, 3, []string{"1", "2", "3"})
	for _, g := range p.ColumnGroups() {
		p.AddGroup(g, puzzle.NewUniqueConstraint{})
	}
	for _ , g := range p.RowGroups() {
		p.AddGroup(g, puzzle.NewUniqueConstraint{})
	}
	p.Init("2")
	s0 := &Settings{0, 1}
	p.Solve(s0)
	fmt.Printf("Decisions made: %d\n", s0.count)
	
	sudoku := puzzle.NewSudokuPuzzle()
	/*
	sudoku.Init(
		"53..7...." +
			"6..195..." +
			".98....6." +
			"8...6...3" +
			"4..8.3..1" +
			"7...2...6" +
			".6....28." +
			"...419..5" +
			"....8..79")
	*/
	sudoku.Init(
		"........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9")

	s1 := &Settings{0, 10000}
	sudoku.Solve(s1)
	fmt.Printf("Decisions made: %d\n", s1.count)
}
