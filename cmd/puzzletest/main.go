package main

import (
	"fmt"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/tracker"
)

type Printer struct {
	p *puzzle.GridPuzzle
}

func (s *Printer) CaptureSolution(p *csp.Problem) {
	s.p.Print()
}

func main() {
	fmt.Println("Hello World!!")
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

	p := puzzle.NewPuzzle2()
	s := puzzle.NewSudoku(p)
	s.Get(0, 0).Value(p.GetIntValue(6))

	sudoku0 := puzzle.NewSudokuPuzzle()
	sudoku0.Init(
		"........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9")

	p0 := Printer{sudoku0}
	s0 := csp.Settings{
		DecisionTracker: tracker.PrintEveryLogN(10),
		SolutionTracker: &p0,
		Decider:         &decide.First{},
	}
	sudoku0.Solve(s0)
	//fmt.Printf("Decisions made: %d\n", p0.count)

	sudoku1 := puzzle.NewSudokuPuzzle()
	sudoku1.Init(
		"........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9")

	p1 := Printer{sudoku1}
	s1 := csp.Settings{
		DecisionTracker: tracker.PrintEveryLogN(10),
		SolutionTracker: &p1,
		Decider:         &decide.Min{},
	}
	sudoku1.Solve(s1)
	//fmt.Printf("Decisions made: %d\n", p1.count)

	sudoku2 := puzzle.NewSudokuPuzzle()
	sudoku2.Init(
		"........." +
			".....3.85" +
			"..1.2...." +
			"...5.7..." +
			"..4...1.." +
			".9......." +
			"5......73" +
			"..2.1...." +
			"....4...9")

	p2 := Printer{sudoku2}
	s2 := csp.Settings{
		DecisionTracker: tracker.PrintEveryLogN(10),
		SolutionTracker: &p2,
		Decider:         &decide.MinMin{},
	}
	sudoku2.Solve(s2)
	//fmt.Printf("Decisions made: %d\n", p2.count)
}
