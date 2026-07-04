package main

import (
	"fmt"

	"github.com/offpath/puzzleutils/internal/csp"
	"github.com/offpath/puzzleutils/internal/decide"
	"github.com/offpath/puzzleutils/internal/layouts"
	"github.com/offpath/puzzleutils/internal/puzzle"
	"github.com/offpath/puzzleutils/internal/puzzletypes"
	"github.com/offpath/puzzleutils/internal/tracker"
)

type Printer struct {
	p *layouts.Grid
}

func (s *Printer) CaptureSolution(p *csp.Problem) bool {
	fmt.Print(s.p)
	return true
}

func main() {
	fmt.Println("Hello World!!")

	p0 := puzzle.NewPuzzle()
	sudoku0 := puzzletypes.NewSudoku(p0)
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

	printer0 := Printer{sudoku0}
	s0 := csp.Settings{
		DecisionTracker: tracker.PrintEveryLogN(10),
		SolutionTracker: &printer0,
		Decider:         &decide.First{},
	}
	p0.Solve(s0)
	//fmt.Printf("Decisions made: %d\n", p0.count)

	p1 := puzzle.NewPuzzle()
	sudoku1 := puzzletypes.NewSudoku(p1)
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

	printer1 := Printer{sudoku1}
	s1 := csp.Settings{
		DecisionTracker: tracker.PrintEveryLogN(10),
		SolutionTracker: &printer1,
		Decider:         &decide.Min{},
	}
	p1.Solve(s1)
	//fmt.Printf("Decisions made: %d\n", p1.count)

	p2 := puzzle.NewPuzzle()
	sudoku2 := puzzletypes.NewSudoku(p2)
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

	printer2 := Printer{sudoku2}
	s2 := csp.Settings{
		DecisionTracker: tracker.PrintEveryLogN(10),
		SolutionTracker: &printer2,
		Decider:         &decide.MinMin{},
	}
	p2.Solve(s2)
	//fmt.Printf("Decisions made: %d\n", p2.count)
}
