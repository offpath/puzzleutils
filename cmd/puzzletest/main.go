package main

import (
	"fmt"

	"puzzleutils/internal/csp"
	"puzzleutils/internal/decide"
	"puzzleutils/internal/puzzle"
)

type Printer struct {
	count int
	multiple int
	p *puzzle.GridPuzzle
}

func (s *Printer) MakeDecision(p *csp.Problem) {
	s.count++
	if s.count % s.multiple == 0 {
		fmt.Printf("Count = %d\n", s.count)
	}
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

	p0 := Printer{0, 100000, sudoku0}
	s0 := csp.Settings{&p0, &decide.First{}}
	sudoku0.Solve(s0)
	fmt.Printf("Decisions made: %d\n", p0.count)

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

	p1 := Printer{0, 100000, sudoku1}
	s1 := csp.Settings{&p1, &decide.Min{}}
	sudoku1.Solve(s1)
	fmt.Printf("Decisions made: %d\n", p1.count)	

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

	p2 := Printer{0, 100000, sudoku2}
	s2 := csp.Settings{&p2, &decide.MinMin{}}
	sudoku2.Solve(s2)
	fmt.Printf("Decisions made: %d\n", p2.count)
}
