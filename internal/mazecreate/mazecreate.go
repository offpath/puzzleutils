package mazecreate

import (
	"fmt"
	"math/rand"
	"time"
)

type square struct {
	open []bool
	visited bool
}

type coord struct {
	row, col int
}

type Maze struct {
	height, width int
	start, finish coord
	squares [][]*square
}

func (m *Maze) wall(i, j int, direction int) string {
	if m.squares[i][j].open[direction] {
		return " "
	}
	if direction % 2 == 0 {
		return "-"
	}
	return "|"
}

func (m *Maze) Print() {
	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j++ {
			fmt.Printf("+%s", m.wall(i, j, 0))
		}
		fmt.Printf("+\n")
		for j := 0; j < m.width; j++ {
			fmt.Printf("%s", m.wall(i, j, 3))
			if m.start.row == i && m.start.col == j {
				fmt.Printf("S")
			} else if m.finish.row == i && m.finish.col == j {
				fmt.Printf("F")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("|\n")
	}
	for j := 0; j < m.width; j++ {
		fmt.Printf("+-")
	}
	fmt.Printf("+\n")
}

func NewMaze(height, width int, branchChancePct int) {
	rand.Seed(time.Now().UTC().UnixNano())
	m := Maze{
		height: height,
		width: width,
		start: coord{0,rand.Intn(width)},
		finish: coord{height-1,rand.Intn(width)},
	}
	for i := 0; i < height; i++ {
		m.squares = append(m.squares, nil)
		for j := 0; j < width; j++ {
			m.squares[i] = append(m.squares[i],
				&square{
					open: []bool{false, false, false, false},
					visited: false,
				})
		}
	}

	m.squares[m.start.row][m.start.col].visited = true
	q := []coord{m.start}
	for len(q) > 0 {
		var pt coord
		pt, q = q[0], q[1:]
		if m.finish.row == pt.row && m.finish.col == pt.col {
			continue
		}
		var avail []int
		if pt.row > 0 && !m.squares[pt.row-1][pt.col].visited {
			avail = append(avail, 0)
		}
		if pt.row < m.height-1 && !m.squares[pt.row+1][pt.col].visited {
			avail = append(avail, 2)
		}
		if pt.col > 0 && !m.squares[pt.row][pt.col-1].visited {
			avail = append(avail, 3)
		}
		if pt.col < m.width-1 && !m.squares[pt.row][pt.col+1].visited {
			avail = append(avail, 1)
		}
		if len(avail) == 0 {
			continue
		}
		rand.Shuffle(len(avail), func(i, j int) {
			avail[i], avail[j] = avail[j], avail[i]
		})
		next := 1
		if len(avail) > 1 && rand.Intn(100) < branchChancePct {
			next = 2
		}
		for i := 0; i < next; i++ {
			pt2 := pt
			switch avail[i] {
			case 0:
				pt2.row -= 1
case 1:
				pt2.col += 1
			case 2:
				pt2.row += 1
			case 3:
				pt2.col -=1
			}
			m.squares[pt.row][pt.col].open[avail[i]] = true
			m.squares[pt2.row][pt2.col].open[(avail[i]+2)%4] = true
			m.squares[pt2.row][pt2.col].visited = true
			q = append(q, pt2)
		}
	}
	m.Print()
			
}
