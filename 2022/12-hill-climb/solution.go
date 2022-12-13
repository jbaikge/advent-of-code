package hillclimb

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"math"

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
	var start, end Point
	grid := make(map[Point]Node)
	for point, height := range s.HeightMap {
		if height == CurrentPosition {
			start = point
		}
		if height == BestSignalPosition {
			end = point
		}
		grid[point] = Node{Point: point, Height: height}
	}

	search := &AStar{
		Grid:  grid,
		Start: Node{Point: start, Height: CurrentPosition},
		Goal:  Node{Point: end, Height: BestSignalPosition},
	}
	path := search.Path()
	for _, node := range path {
		fmt.Fprintln(w, node)
	}
	fmt.Fprintf(w, "Part 1: %d\n", len(path)-1)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	return
}

type Node struct {
	Point  Point
	Height byte
}

func (n Node) String() string {
	return fmt.Sprintf("%s %s", n.Point.String(), string(n.Height))
}

// Adapted from pseudocode found on Wikipedia
// Ref: https://en.wikipedia.org/wiki/A*_search_algorithm
type AStar struct {
	Grid  map[Point]Node
	Start Node
	Goal  Node
}

func (a *AStar) Cost(n Node) (cost int) {
	x := a.Start.Point.X - a.Goal.Point.X
	if x < 0 {
		x *= -1
	}

	y := a.Start.Point.Y - a.Goal.Point.Y
	if y < 0 {
		y *= -1
	}

	return x + y
}

func (a *AStar) Path() (path []Node) {
	open := []Node{a.Start}
	from := make(map[Node]Node)
	gScore := map[Node]int{a.Start: 0}
	fScore := map[Node]int{a.Start: a.Cost(a.Start)}

	for len(open) > 0 {
		// Determine current node, or the node in open with the lowest fScore
		var current Node
		var currentIdx int
		minScore := math.MaxInt
		for i, node := range open {
			if score := fScore[node]; score < minScore {
				minScore = score
				current = node
				currentIdx = i
			}
		}

		if current == a.Goal {
			fmt.Println("Found the goal!")
			return a.Reconstruct(from, current)
		}

		// Remove node from open
		open = append(open[:currentIdx], open[currentIdx+1:]...)

		points := []Point{
			{current.Point.X, current.Point.Y + 1},
			{current.Point.X + 1, current.Point.Y},
			{current.Point.X, current.Point.Y - 1},
			{current.Point.X - 1, current.Point.Y},
		}
		for _, point := range points {
			neighbor, found := a.Grid[point]
			if !found {
				continue
			}

			valid := false
			if current == a.Start {
				valid = true
			}
			if current.Height == neighbor.Height {
				valid = true
			}
			if current.Height+1 == neighbor.Height {
				fmt.Println(current, "->", neighbor)
				valid = true
			}
			if current.Height == 'z' && neighbor.Height == a.Goal.Height {
				valid = true
			}
			if !valid {
				continue
			}

			tentativeGScore := gScore[current] + 1
			neighborGScore, found := gScore[neighbor]
			if !found {
				neighborGScore = math.MaxInt
			}
			if tentativeGScore < neighborGScore {
				from[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + a.Cost(neighbor)

				// "if neighbor not in open, open.add(neighbor)"
				openFound := false
				for _, node := range open {
					if node == neighbor {
						openFound = true
					}
				}
				if !openFound {
					// fmt.Println("Appending", neighbor)
					open = append(open, neighbor)
				}
			} else {
				fmt.Println(neighbor, tentativeGScore, "<=", neighborGScore)
			}
		}
	}

	return
}

func (a *AStar) Reconstruct(from map[Node]Node, current Node) (path []Node) {
	path = make([]Node, 0, len(from)+1)
	path = append(path, current)
	for {
		newNode, found := from[current]
		if !found {
			break
		}
		current = newNode
		path = append([]Node{current}, path...)
	}
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
