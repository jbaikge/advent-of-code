package hillclimb

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"math"
	"sort"
	"sync"

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
	return fmt.Sprintf("(%2d, %2d)", p.X, p.Y)
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
	path := search.Path(search.Start)
	// for _, node := range path {
	// 	fmt.Fprintln(w, node)
	// }
	fmt.Fprintf(w, "Part 1: %d\n", len(path)-1)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	var end Point
	startNodes := make([]Node, 0, 1000)
	grid := make(map[Point]Node)
	for point, height := range s.HeightMap {
		if height == BestSignalPosition {
			end = point
		}
		if height == 'a' {
			startNodes = append(startNodes, Node{Point: point, Height: height})
		}
		grid[point] = Node{Point: point, Height: height}
	}

	search := &AStar{
		Grid: grid,
		Goal: Node{Point: end, Height: BestSignalPosition},
	}

	ch := make(chan int, len(startNodes))
	var wg sync.WaitGroup
	wg.Add(len(startNodes))
	for _, node := range startNodes {
		go func(node Node) {
			path := search.Path(node)
			wg.Done()
			length := len(path)
			if length == 0 {
				return
			}
			ch <- length - 1
		}(node)
	}

	wg.Wait()
	close(ch)

	lengths := make([]int, 0, len(startNodes))
	for length := range ch {
		lengths = append(lengths, length)
	}

	sort.Ints(lengths)
	fmt.Fprintf(w, "Part 2: %v\n", lengths[0])
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

func (a *AStar) Cost(n Node) (cost float64) {
	x := n.Point.X - a.Goal.Point.X
	y := n.Point.Y - a.Goal.Point.Y
	return math.Sqrt(float64(x*x + y*y))
}

func (a *AStar) Path(startNode Node) (path []Node) {
	open := []Node{startNode}
	from := make(map[Node]Node)
	gScore := map[Node]int{startNode: 0}
	fScore := map[Node]float64{startNode: a.Cost(startNode)}

	for len(open) > 0 {
		// Determine current node, or the node in open with the lowest fScore
		var current Node
		var currentIdx int
		minScore := math.MaxFloat64
		for i, node := range open {
			if score := fScore[node]; score < minScore {
				minScore = score
				current = node
				currentIdx = i
			}
		}

		if current == a.Goal {
			// fmt.Println("Found the goal!")
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
			if current.Height == CurrentPosition {
				valid = true
			}
			if current.Height > neighbor.Height {
				valid = true
			}
			if current.Height == neighbor.Height {
				valid = true
			}
			if current.Height+1 == neighbor.Height {
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
				fScore[neighbor] = float64(tentativeGScore) + a.Cost(neighbor)

				// "if neighbor not in open, open.add(neighbor)"
				openFound := false
				for _, node := range open {
					if node == neighbor {
						openFound = true
					}
				}
				if !openFound {
					open = append(open, neighbor)
				}
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
