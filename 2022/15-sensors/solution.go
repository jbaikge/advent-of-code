package sensors

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jbaikge/advent-of-code/util"
)

const (
	MarkerSensor = 'S'
	MarkerBeacon = 'B'
	MarkerSignal = '#'
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

type Sensor struct {
	Position Point
	Beacon   Point
}

func (s Sensor) Distance() int {
	x := s.Position.X - s.Beacon.X
	if x < 0 {
		x *= -1
	}

	y := s.Position.Y - s.Beacon.Y
	if y < 0 {
		y *= -1
	}

	return x + y
}

type Solution struct {
	Sensors []Sensor
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Sensors = make([]Sensor, 0, 33)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// Sensor at x=2, y=18: closest beacon is at x=-2, y=15
		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '=' || r == ',' || r == ':'
		})
		var sensor Sensor
		sensor.Position.X, _ = strconv.Atoi(fields[3])
		sensor.Position.Y, _ = strconv.Atoi(fields[5])
		sensor.Beacon.X, _ = strconv.Atoi(fields[11])
		sensor.Beacon.Y, _ = strconv.Atoi(fields[13])
		s.Sensors = append(s.Sensors, sensor)
	}
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	targetY := 10
	if s.Sensors[0].Position.X == 2302110 {
		targetY = 2000000
	}

	beacons := make(map[int]bool)
	for _, sensor := range s.Sensors {
		if sensor.Beacon.Y == targetY {
			beacons[sensor.Beacon.X] = true
		}
	}

	points := make(map[int]bool)
	for _, sensor := range s.Sensors {
		dist := sensor.Distance()
		y := sensor.Position.Y
		fromAbove := y <= targetY && y+dist >= targetY
		fromBelow := y >= targetY && y-dist <= targetY
		crossesTarget := fromAbove || fromBelow
		if !crossesTarget {
			continue
		}

		targetDistance := targetY - y
		if targetDistance < 0 {
			targetDistance *= -1
		}
		xSpread := dist - targetDistance
		for x := sensor.Position.X - xSpread; x <= sensor.Position.X+xSpread; x++ {
			if _, ok := beacons[x]; ok {
				continue
			}
			points[x] = true
		}
	}

	fmt.Fprintf(w, "Part 1: %d\n", len(points))
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	return
}
