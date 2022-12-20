package tetris

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strings"

	"github.com/jbaikge/advent-of-code/util"
)

const (
	Air         = '.'
	RockSegment = '#'

	ArenaWidth = 7
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Point struct {
	X int
	Y int
}

type Rock struct {
	Name   string
	Fill   byte
	Width  int
	Height int
	Shape  []Point
}

type Position struct {
	Rock  Rock
	Point Point
}

type Arena struct {
	Rows [][ArenaWidth]byte
}

func NewArena() (a *Arena) {
	return &Arena{
		Rows: make([][ArenaWidth]byte, 0, 1024),
	}
}

func (a *Arena) Apply(p *Position) {
	// rowsNeeded := 0
	// fmt.Printf("%d - %d + %d = %d\n", p.Point.Y, a.Height(), p.Rock.Height, p.Point.Y-a.Height()+p.Rock.Height)
	for i := a.Height(); i <= p.Point.Y; i++ {
		row := [ArenaWidth]byte{}
		for r := range row {
			row[r] = Air
		}
		a.Rows = append(a.Rows, row)
	}
	for _, point := range p.Rock.Shape {
		x := p.Point.X + point.X
		y := p.Point.Y - point.Y
		// fmt.Printf("rows: %d; (%d, %d)\n", len(a.Rows), x, y)
		// a.Rows[y][x] = p.Rock.Fill
		a.Rows[y][x] = RockSegment
	}
}

func (a Arena) Drop(p *Position) (moved bool) {
	for _, point := range p.Rock.Shape {
		x := p.Point.X + point.X
		y := p.Point.Y - point.Y - 1
		if y < 0 {
			return false
		}
		if y >= a.Height() {
			continue
		}
		if a.Rows[y][x] != Air {
			return false
		}
	}

	p.Point.Y--
	return true
}

func (a Arena) Height() int {
	return len(a.Rows)
}

// Check the left or right side of the Rock to see if it will collide with the
// wall or another rock
func (a Arena) Push(p *Position, dir int) (moved bool) {
	for _, point := range p.Rock.Shape {
		x := p.Point.X + point.X + dir
		y := p.Point.Y - point.Y
		if x < 0 {
			return false
		}
		if x >= ArenaWidth {
			return false
		}
		if y >= a.Height() {
			continue
		}
		if a.Rows[y][x] != Air {
			return false
		}
	}

	p.Point.X += dir
	return true
}

func (a Arena) String() string {
	var sb strings.Builder
	for i := len(a.Rows) - 1; i >= 0; i-- {
		sb.Write(a.Rows[i][:])
		sb.WriteByte('\n')
	}
	return sb.String()
}

type Solution struct {
	Pattern []byte
	Rocks   []Rock
	Width   int
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Width = 7
	s.Rocks = []Rock{
		{
			Name:   "Bar",
			Fill:   '-',
			Width:  4,
			Height: 1,
			Shape: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 2, Y: 0},
				{X: 3, Y: 0},
			},
		},
		{
			Name:   "Plus",
			Fill:   '+',
			Width:  3,
			Height: 3,
			Shape: []Point{
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 1, Y: 2},
			},
		},
		{
			Name:   "Ell",
			Fill:   '%',
			Width:  3,
			Height: 3,
			Shape: []Point{
				{X: 2, Y: 0},
				{X: 2, Y: 1},
				{X: 0, Y: 2},
				{X: 1, Y: 2},
				{X: 2, Y: 2},
			},
		},
		{
			Name:   "Rod",
			Fill:   '|',
			Width:  1,
			Height: 4,
			Shape: []Point{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 0, Y: 2},
				{X: 0, Y: 3},
			},
		},
		{
			Name:   "Box",
			Fill:   '#',
			Width:  2,
			Height: 2,
			Shape: []Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
			},
		},
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s.Pattern = scanner.Bytes()
	}
	return
}

func (s Solution) DropRocks(num int) (arena *Arena) {
	var moveCounter int
	arena = NewArena()
	for n := 0; n < num; n++ {
		rock := s.Rocks[n%len(s.Rocks)]
		position := &Position{
			Rock:  rock,
			Point: Point{X: 2, Y: arena.Height() + rock.Height + 2},
		}
		// Push and drop until at rest
		for {
			direction := s.Pattern[moveCounter%len(s.Pattern)]
			push := 1
			if direction == '<' {
				push = -1
			}
			moveCounter++

			arena.Push(position, push)
			dropped := arena.Drop(position)

			if !dropped {
				break
			}
		}
		arena.Apply(position)
	}
	// return arena.Height()
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	arena := s.DropRocks(2022)
	fmt.Fprintf(w, "Part 1: %d\n", arena.Height())
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	drop := 1000000000000
	// Turns out len(pattern) and len(rocks) are prime so LCM is easy
	// lcm := len(s.Pattern) * len(s.Rocks)
	// stacks := drop / lcm
	// remain := drop % lcm
	// stacksHeight := s.DropRocks(lcm)
	// remainHeight := s.DropRocks(remain)

	// fmt.Printf("LCM: %d; Stacks: %d; Remain: %d; Stack Height: %d; Remain Height: %d\n",
	// 	lcm,
	// 	stacks,
	// 	remain,
	// 	stacksHeight,
	// 	remainHeight,
	// )

	// for n := 0; n < 5; n++ {
	// 	drop := remain + lcm*n
	// 	height := s.DropRocks(drop)
	// 	fmt.Printf("%d + %d * %d = %6d - %6d\n", remain, lcm, n, drop, height)
	// }

	// for i := 1; i < 100000; i++ {
	// 	single := s.DropRocks(i)
	// 	for _, n := range []int{2, 3, 4, 5} {
	// 		height := s.DropRocks(i * n)
	// 		if height%single > 0 {
	// 			break
	// 		}
	// 		fmt.Printf("%d * %d: %6d -> %6d %6d\n", i, n, i*n, height, height%single)
	// 	}
	// }

	// arena := s.DropRocks(10000)
	// rows := arena.Rows
	// maxNum := 0
	// for num := 100; num < 1000; num++ {
	// 	gotMatch := false
	// 	for skip := 0; skip < len(rows)-num; skip++ {
	// 		search := rows[skip : skip+num]
	// 		for i := skip + num; i < len(rows)-num; i++ {
	// 			match := 0
	// 			for s, row := range search {
	// 				if !bytes.Equal(row[:], rows[i+s][:]) {
	// 					break
	// 				}
	// 				match++
	// 			}
	// 			if match == num {
	// 				maxNum = num
	// 				gotMatch = true
	// 				break
	// 			}
	// 		}
	// 	}
	// 	if gotMatch && maxNum > 0 {
	// 		fmt.Println("Max Num:", maxNum)
	// 	}
	// 	if !gotMatch && maxNum > 0 {
	// 		break
	// 	}
	// }
	// fmt.Println("Max Num:", maxNum)

	fmt.Fprintf(w, "Part 2: %d\n", drop)
	return
}
