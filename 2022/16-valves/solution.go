package valves

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"math"
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

type AStar struct {
	Map       map[string]Valve
	Start     Valve
	TimeLimit int
}

func (a *AStar) Cost(v Valve) int {
	return v.FlowRate
}

func (a *AStar) Path() (path []Valve) {
	startValve := a.Start
	open := []Valve{startValve}
	// Map of valve name to valve name
	from := make(map[string]string)
	gScore := map[string]int{startValve.Name: 0}
	fScore := map[string]int{startValve.Name: a.Cost(startValve)}
	// Only have a certain amount of time to complete the traversal
	timeTaken := 0

	// Zero-Flow-Rate valves can be considered "already open" since no action
	// will happen.
	openValves := make(map[string]bool)
	for _, valve := range a.Map {
		openValves[valve.Name] = valve.FlowRate > 0
	}

	var current Valve
	for timeTaken < a.TimeLimit && len(open) > 0 {
		var currentIdx int
		minScore := math.MaxInt
		for i, valve := range open {
			if score := fScore[valve.Name]; score < minScore {
				minScore = score
				current = valve
				currentIdx = i
			}
		}

		open = append(open[:currentIdx], open[currentIdx+1:]...)

		timeTaken++
		if openValves[current.Name] {
			timeTaken++
			openValves[current.Name] = false
		}

		for _, tunnel := range current.TunnelsTo {
			tentativeGScore := gScore[current.Name] + 1
			tunnelGScore, found := gScore[tunnel]
			if !found {
				tunnelGScore = math.MaxInt
			}
			fmt.Printf("%d < %d\n", tentativeGScore, tunnelGScore)
			if tentativeGScore < tunnelGScore {
				from[tunnel] = current.Name
				gScore[tunnel] = tentativeGScore
				fScore[tunnel] = tentativeGScore + a.Cost(a.Map[tunnel])

				// Explore all the sub-tunnels
				for _, subTunnel := range a.Map[tunnel].TunnelsTo {
					open = append(open, a.Map[subTunnel])
				}
			}
		}
	}

	return a.Reconstruct(from, current)
}

func (a *AStar) Reconstruct(from map[string]string, current Valve) (path []Valve) {
	path = make([]Valve, 0, len(from)+1)
	path = append(path, current)
	for {
		newName, found := from[current.Name]
		if !found {
			break
		}
		current = a.Map[newName]
		path = append([]Valve{current}, path...)
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
	for _, valve := range s.Valves {
		valveMap[valve.Name] = valve
	}

	aStar := &AStar{
		Map:       valveMap,
		Start:     s.Valves[0],
		TimeLimit: 30,
	}
	path := aStar.Path()

	fmt.Printf("%+v\n", path)

	fmt.Fprintf(w, "Part 1: %d\n", 0)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	fmt.Fprintf(w, "Part 2: %d\n", 0)
	return
}
