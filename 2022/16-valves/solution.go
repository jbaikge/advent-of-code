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

type Search struct {
	Graph map[string]Valve
}

func (s Search) Path(start, goal string) (path []Valve) {
	stack := make([]Node, 0, len(s.Graph))
	stack = append(stack, Node{Name: start, Path: make([]string, 0, len(s.Graph))})

	visited := make(map[string]bool)
	for len(stack) > 0 {
		current := stack[0]
		stack = stack[1:]

		if current.Name == goal {
			path = make([]Valve, 0, len(current.Path))
			for _, name := range current.Path {
				path = append(path, s.Graph[name])
			}
			return
		}

		if _, found := visited[current.Name]; found {
			continue
		}
		visited[current.Name] = true

		current.Path = append(current.Path, current.Name)

		for _, to := range s.Graph[current.Name].TunnelsTo {
			node := Node{Name: to, Path: make([]string, len(current.Path))}
			copy(node.Path, current.Path)
			stack = append(stack, node)
		}
	}
	return
}

type Route struct {
	Open      []string
	Closed    []string
	TimeTaken int
	Flow      int
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
	const TimeLimit = 30

	graph := make(map[string]Valve)
	pressurized := make([]string, 0, len(s.Valves))
	for _, valve := range s.Valves {
		graph[valve.Name] = valve
		if valve.FlowRate > 0 {
			pressurized = append(pressurized, valve.Name)
		}
	}

	routes := make([]Route, 0, 1024)
	search := Search{
		Graph: graph,
	}

	// Initial routes from AA
	for _, name := range pressurized {
		valve := graph[name]
		path := search.Path("AA", name)
		timeTaken := len(path) + 1
		route := Route{
			Open:      make([]string, 0, len(pressurized)-1),
			Closed:    []string{name},
			TimeTaken: timeTaken,
			Flow:      (TimeLimit - timeTaken) * valve.FlowRate,
		}
		for _, open := range pressurized {
			if open == name {
				continue
			}
			route.Open = append(route.Open, open)
		}
		routes = append(routes, route)
	}

	var bestRoute Route

	stack := make([]Route, 0, 1024)
	stack = append(stack, routes...)
	for len(stack) > 0 {
		current := stack[0]
		stack = stack[1:]

		start := current.Closed[len(current.Closed)-1]
		for _, open := range current.Open {
			path := search.Path(start, open)
			timeTaken := current.TimeTaken + len(path) + 1
			route := Route{
				Open: make([]string, 0, len(current.Open)-1),
				// Closed:    append(current.Closed, open),
				Closed:    make([]string, 0, len(current.Closed)+1),
				TimeTaken: timeTaken,
				Flow:      current.Flow + (TimeLimit-timeTaken)*graph[open].FlowRate,
			}
			for _, name := range current.Open {
				if name == open {
					continue
				}
				route.Open = append(route.Open, name)
			}
			route.Closed = append(route.Closed, current.Closed...)
			route.Closed = append(route.Closed, open)

			if route.TimeTaken < TimeLimit {
				if route.Flow > bestRoute.Flow {
					bestRoute = route
				}
				if len(route.Open) > 0 {
					stack = append(stack, route)
				}
			}
		}
	}

	fmt.Fprintf(w, "Part 1: %d\n", bestRoute.Flow)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	fmt.Fprintf(w, "Part 2: %d\n", 0)
	return
}
