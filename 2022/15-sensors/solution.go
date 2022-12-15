package sensors

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

func (s Sensor) BeaconDistance() int {
	return s.Distance(s.Beacon)
}

func (s Sensor) Distance(to Point) int {
	x := s.Position.X - to.X
	if x < 0 {
		x *= -1
	}

	y := s.Position.Y - to.Y
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
		dist := sensor.BeaconDistance()
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
	// Calculate a box around all the sensor locations
	var min, max Point
	min.X, min.Y = math.MaxInt, math.MaxInt
	// Cache beacon distances for the next part
	beaconDistances := make([]int, len(s.Sensors))
	for i, sensor := range s.Sensors {
		x, y := sensor.Position.X, sensor.Position.Y
		if x < min.X {
			min.X = x
		}
		if x > max.X {
			max.X = x
		}
		if y < min.Y {
			min.Y = y
		}
		if y > max.Y {
			max.Y = y
		}

		beaconDistances[i] = sensor.BeaconDistance()
	}

	// Walk 1 tile outside the perimeter of each sensor's scanning radius
	// Then compare each point on the perimeter to the other sensors' scanning radius
	// If a point does not fall in any radius, it is the location of the missing beacon
	var found Point
	for _, sensor := range s.Sensors {
		beaconDistance := sensor.BeaconDistance()
		perimeter := make([]Point, 0, beaconDistance*2)
		// NE quadrant
		for x, y := sensor.Position.X, sensor.Position.Y-beaconDistance-1; x <= sensor.Position.X+beaconDistance+1; x, y = x+1, y+1 {
			perimeter = append(perimeter, Point{X: x, Y: y})
		}
		// SE quadrant
		for x, y := sensor.Position.X, sensor.Position.Y+beaconDistance+1; x <= sensor.Position.X+beaconDistance+1; x, y = x+1, y-1 {
			perimeter = append(perimeter, Point{X: x, Y: y})
		}
		// NW quadrant
		for x, y := sensor.Position.X-beaconDistance-1, sensor.Position.Y; x <= sensor.Position.X; x, y = x+1, y-1 {
			perimeter = append(perimeter, Point{X: x, Y: y})
		}
		// SW quadrant
		for x, y := sensor.Position.X-beaconDistance-1, sensor.Position.Y; x <= sensor.Position.X; x, y = x+1, y+1 {
			perimeter = append(perimeter, Point{X: x, Y: y})
		}
		// Walk around the outside and find out if any sensor
		// radii overlap
		for _, point := range perimeter {
			if point.X < min.X || point.X > max.X || point.Y < min.Y || point.Y > max.Y {
				continue
			}
			overlaps := false
			for i, dist := range beaconDistances {
				if s.Sensors[i].Distance(point) <= dist {
					overlaps = true
					break
				}
			}
			if !overlaps {
				found = point
				break
			}
		}
		if found.X != 0 && found.Y != 0 {
			break
		}
	}

	fmt.Fprintf(w, "Part 2: %d\n", found.X*4000000+found.Y)
	return
}
