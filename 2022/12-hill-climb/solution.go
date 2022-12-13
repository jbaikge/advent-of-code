package hillclimb

import (
	"bufio"
	"embed"
	"fmt"
	"io"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var files embed.FS

const (
	CurrentPosition    = 'S'
	BestSignalPosition = 'E'
)

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

var _ util.Solution = new(Solution)

type Solution struct {
	HeightMap map[Point]byte
}

func (s *Solution) Files() embed.FS {
	return files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	y := 0
	s.HeightMap = make(map[Point]byte)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		for x, height := range scanner.Bytes() {
			s.HeightMap[Point{x, y}] = height
		}
		y++
	}
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	var start Point
	for point, height := range s.HeightMap {
		if height == CurrentPosition {
			start = point
		}
	}
	visited := make(map[Point]bool)
	steps := dfs(s.HeightMap, visited, start, start)
	fmt.Fprintf(w, "Part 1: %d\n", steps)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	return
}

func dfs(grid map[Point]byte, visited map[Point]bool, p Point, last Point) (steps int) {
	// Does this point exist in the grid?
	height, ok := grid[p]
	if !ok {
		return
	}

	// Have we visited here before?
	if visited[p] {
		// fmt.Printf("%s Already visited\n", p)
		return
	}

	if height == BestSignalPosition {
		fmt.Printf("%s %s\n", string(grid[last]), string(BestSignalPosition))
	}

	// New position is lower than current position
	if height < grid[last] {
		// fmt.Printf("%s %s < %s\n", p, string(height), string(grid[last]))
		return
	}

	visited[p] = true

	points := []Point{
		{p.X, p.Y + 1},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X - 1, p.Y},
	}
	for _, point := range points {
		if height == CurrentPosition {
			visited = map[Point]bool{}
		}
		steps = dfs(grid, visited, point, p)
		if height == CurrentPosition {
			fmt.Println(point, len(visited))
		}
	}

	return
}
