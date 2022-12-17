package valves

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Valve struct {
	Name      string
	FlowRate  int
	TunnelsTo []string
}

func (v Valve) String() string {
	return fmt.Sprintf("%s.%d", v.Name, v.FlowRate)
}

type Node struct {
	Name string
	Path []string
}

type BFS struct {
	Graph     map[string]Valve
	TimeLimit int
}

func (b BFS) Path(start, end string) (path []Valve) {
	stack := make([]Node, 0, len(b.Graph))
	stack = append(stack, Node{Name: start, Path: make([]string, 0, len(b.Graph))})
	// neighbor -> current
	visited := make(map[string]bool)
	for len(stack) > 0 {
		// Pull from beginning of stack
		current := stack[0]
		stack = stack[1:]

		// Don't revisit valves
		if _, ok := visited[current.Name]; ok {
			continue
		}
		visited[current.Name] = true

		// Append to the current node's path
		current.Path = append(current.Path, current.Name)

		// Convert the path to Valves when we reach the target valve
		if current.Name == end {
			path = make([]Valve, 0, len(current.Path))
			for _, name := range current.Path {
				path = append(path, b.Graph[name])
			}
			return
		}

		// Generate new nodes for adjacent valves
		for _, to := range b.Graph[current.Name].TunnelsTo {
			node := Node{Name: to, Path: make([]string, len(current.Path))}
			copy(node.Path, current.Path)
			stack = append(stack, node)
		}
	}
	return
}

type Solution struct {
	Valves []Valve
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Valves = make([]Valve, 0, 64)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '=' || r == ',' || r == ';'
		})
		rate, err := strconv.Atoi(fields[5])
		if err != nil {
			return err
		}
		s.Valves = append(s.Valves, Valve{
			Name:      fields[1],
			FlowRate:  rate,
			TunnelsTo: fields[10:],
		})
	}
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	valveMap := make(map[string]Valve)
	pressurized := make([]string, 0, len(valveMap))
	for _, valve := range s.Valves {
		valveMap[valve.Name] = valve
		if valve.FlowRate > 0 {
			pressurized = append(pressurized, valve.Name)
		}
	}

	bfs := &BFS{
		Graph:     valveMap,
		TimeLimit: 30,
	}
	startNode := "AA"
	timeLimit := 30
	timeTaken := 0
	for len(pressurized) > 0 && timeTaken <= 30 {
		fmt.Printf("Pressurized: %d; Time Taken: %d\n", len(pressurized), timeTaken)
		var idx int
		minLength := len(s.Valves)
		paths := make(map[int][]Valve)
		for i, end := range pressurized {
			fmt.Printf("  %s -> %s\n", startNode, end)
			newPath := bfs.Path(startNode, end)

			score := (timeLimit - timeTaken - len(newPath) - 1) * newPath[len(newPath)-1].FlowRate
			fmt.Printf("    %3d %v\n", score, newPath)

			if length := len(newPath); length < minLength {
				minLength = length
			}

			path, found := paths[len(newPath)]
			if !found || path[len(path)-1].FlowRate < newPath[len(newPath)-1].FlowRate {
				paths[len(newPath)] = newPath
				if len(newPath) == minLength {
					idx = i
				}
			}
		}
		startNode = pressurized[idx]
		pressurized = append(pressurized[:idx], pressurized[idx+1:]...)
	}

	fmt.Fprintf(w, "Part 1: %d\n", 0)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	fmt.Fprintf(w, "Part 2: %d\n", 0)
	return
}
