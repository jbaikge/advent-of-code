package reservoir

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jbaikge/advent-of-code/util"
)

const (
	Rock = '#'
	Air  = '.'
	Sand = '+'
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

type Path []Point

func ParsePath(line string) (p Path) {
	points := strings.Split(line, " -> ")
	p = make(Path, len(points))
	for i, point := range points {
		values := strings.SplitN(point, ",", 2)
		p[i].X, _ = strconv.Atoi(values[0])
		p[i].Y, _ = strconv.Atoi(values[1])
	}
	return
}

type Cave struct {
	Tiles map[Point]byte
	Floor int // Y-height of floor (maxY + 2)
}

func NewCave() *Cave {
	return &Cave{
		Tiles: make(map[Point]byte),
	}
}

func (c *Cave) BuildRocks(paths []Path) {
	for _, path := range paths {
		for i := 0; i < len(path)-1; i++ {
			start, end := path[i], path[i+1]
			minX, maxX := start.X, end.X
			if minX > maxX {
				minX, maxX = maxX, minX
			}

			minY, maxY := start.Y, end.Y
			if minY > maxY {
				minY, maxY = maxY, minY
			}

			for x := minX; x <= maxX; x++ {
				for y := minY; y <= maxY; y++ {
					c.Tiles[Point{X: x, Y: y}] = Rock
				}
			}

			if floor := maxY + 2; floor > c.Floor {
				c.Floor = floor
			}
		}
	}
}

func (c *Cave) Count() (total int) {
	for _, tile := range c.Tiles {
		if tile == Sand {
			total++
		}
	}
	return
}

// Ideally a new grain of sand will start at 500,0, then fall to 500,1
// then 500,1 will get fed back into the method to drop to 500, 2 until
// currentPos == newPos to represent "At Rest"
func (c *Cave) DropSand(currentPos Point) (newPos Point) {
	possibleMoves := []Point{
		{X: currentPos.X, Y: currentPos.Y + 1},
		{X: currentPos.X - 1, Y: currentPos.Y + 1},
		{X: currentPos.X + 1, Y: currentPos.Y + 1},
	}
	for _, move := range possibleMoves {
		if tile, ok := c.Tiles[move]; ok && tile != Air {
			continue
		}
		c.Tiles[currentPos] = Air
		c.Tiles[move] = Sand
		return move
	}
	return currentPos
}

func (c Cave) String() (str string) {
	var min, max Point
	// min.X, min.Y = math.MaxInt, math.MaxInt
	min.X = math.MaxInt

	for p := range c.Tiles {
		if p.X < min.X {
			min.X = p.X
		}
		if p.X > max.X {
			max.X = p.X
		}
		// if p.Y < min.Y {
		// 	min.Y = p.Y
		// }
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}

	grid := make([][]byte, max.Y-min.Y+1)
	for i := range grid {
		grid[i] = bytes.Repeat([]byte{Air}, max.X-min.X+1)
	}

	for p, tile := range c.Tiles {
		grid[p.Y-min.Y][p.X-min.X] = tile
	}

	return string(bytes.Join(grid, []byte{'\n'}))
}

type Solution struct {
	Pour  Point
	Paths []Path
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Pour.X = 500

	s.Paths = make([]Path, 0, 141)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s.Paths = append(s.Paths, ParsePath(scanner.Text()))
	}

	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	cave := NewCave()
	cave.BuildRocks(s.Paths)

	var grain, lastGrain Point
	for steps := 0; steps < 200000; steps++ {
		// Hit the floor
		if grain.Y == cave.Floor-1 {
			delete(cave.Tiles, grain)
			break
		}

		// At rest, start a new grain
		if lastGrain == grain {
			grain.X, grain.Y = 500, 0
			cave.Tiles[grain] = Sand
		}

		lastGrain = grain
		grain = cave.DropSand(grain)
	}

	fmt.Fprintf(w, "Part 1: %d\n", cave.Count())
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	cave := NewCave()
	cave.BuildRocks(s.Paths)

	var grain, lastGrain Point
	startNew := true
	for steps := 0; steps < 20000000; steps++ {
		// Hit the floor, start a new grain
		if grain.Y == cave.Floor-1 {
			startNew = true
		}

		// At rest, start a new grain
		if !startNew && lastGrain == grain {
			if grain.Y == 0 {
				break
			}
			startNew = true
		}

		if startNew {
			grain.X, grain.Y = s.Pour.X, s.Pour.Y
			cave.Tiles[grain] = Sand
			startNew = false
		}

		lastGrain = grain
		grain = cave.DropSand(grain)
	}

	fmt.Fprintf(w, "Part 2: %d\n", cave.Count())
	return
}
