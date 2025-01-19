package maze

import (
	"fmt"

	"github.com/offpath/puzzleutils/internal/pathplan"
)

type Player struct {
	sq *Square
}

func (p *Player) String() string {
	return p.sq.String()
}

func (p *Player) Clone() *Player {
	return &Player{p.sq}
}

type Square struct {
	isWall bool
	i, j   int
}

func (s *Square) String() string {
	return fmt.Sprintf("(%d, %d)", s.i, s.j)
}

type GridBoard struct {
	grid   [][]*Square
	player *Player
	exit   *Square
}

func (b *GridBoard) String() string {
	return b.player.String()
}

func (b *GridBoard) Clone() pathplan.Board {
	return &GridBoard{
		grid:   b.grid,
		player: b.player.Clone(),
		exit:   b.exit,
	}
}

func (b *GridBoard) IsSolved() bool {
	return b.player.sq == b.exit
}

func (b *GridBoard) IsFailed() bool {
	return false
}

func abs(i int) int {
	if i < 0 {
		return -1
	}
	return i
}

func distance(a, b *Square) int {
	return abs(a.i-b.i) + abs(a.j-b.j)
}

func (b *GridBoard) Heuristic() int {
	return distance(b.player.sq, b.exit)
}

var controls []string = []string{"N", "E", "S", "W"}

func (b *GridBoard) Controls() []string {
	return controls
}

func (b *GridBoard) DoControl(c string) {
	di, dj := 0, 0
	switch c {
	case "N":
		di = -1
	case "E":
		dj = 1
	case "S":
		di = 1
	case "W":
		dj = -1
	}
	i, j := b.player.sq.i+di, b.player.sq.j+dj
	if i < 0 || i >= len(b.grid) || j < 0 || j >= len(b.grid[i]) || b.grid[i][j].isWall {
		return
	}
	b.player.sq = b.grid[i][j]

}

func LoadBoard(boardText []string) *GridBoard {
	result := &GridBoard{nil, nil, nil}
	for i, s := range boardText {
		row := []*Square{}
		for j, c := range s {
			sq := &Square{c == 'X', i, j}
			row = append(row, sq)
			if c == 'E' {
				result.exit = sq
			} else if c == 'P' {
				result.player = &Player{sq}
			}
		}
		result.grid = append(result.grid, row)
	}
	return result
}
