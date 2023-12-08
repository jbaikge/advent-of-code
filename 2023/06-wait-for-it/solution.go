package waitforit

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

type Race struct {
	Time     int
	Distance int
}

type Solution struct {
	Races []Race
}

func (s Solution) WaysToWin(r Race) int {
	// y = ax^2 + bx + c
	// <traveled_past> = -x^2 + <time>x - <winning_distance>
	a := float64(-1)
	b := float64(r.Time)
	c := float64(-r.Distance - 1)

	x1 := (-b + math.Sqrt(b*b-4*a*c)) / (2 * a)
	x2 := (-b - math.Sqrt(b*b-4*a*c)) / (2 * a)
	x1 = math.Ceil(x1)
	x2 = math.Floor(x2)
	return int(x2 - x1 + 1)
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	scanner := bufio.NewScanner(r)

	// Parse Times
	scanner.Scan()
	timeNums := strings.Fields(scanner.Text())[1:]
	s.Races = make([]Race, len(timeNums))
	for i, num := range timeNums {
		if s.Races[i].Time, err = strconv.Atoi(num); err != nil {
			return
		}
	}

	// Parse distances
	scanner.Scan()
	for i, num := range strings.Fields(scanner.Text())[1:] {
		if s.Races[i].Distance, err = strconv.Atoi(num); err != nil {
			return
		}
	}

	return
}

// Observation: The numbers make a curve up and then down
// This means you can make a quadratic equation to position
// the curve where anything above the X-axis wins the race.
// Casting the numbers to integers will give the number of
// ways to win
func (s Solution) Part1(w io.Writer) (err error) {
	ways := 1

	// for _, race := range s.Races {
	// 	wins := 0
	// 	for hold := 1; hold < race.Time; hold++ {
	// 		remaining := race.Time - hold
	// 		distance := hold * remaining
	// 		if distance > race.Distance {
	// 			wins++
	// 		}
	// 	}
	// 	ways *= wins
	// }

	for _, race := range s.Races {
		ways *= s.WaysToWin(race)
	}
	fmt.Fprintf(w, "Part 1: %d\n", ways)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	// Combine the numbers into a single number
	race := s.Races[0]
	for _, r := range s.Races[1:] {
		// Determine the multiplier to make room
		// for the next number
		var m int
		for m = 1; m < r.Time; m *= 10 {
			// NOOP
		}
		race.Time = race.Time*m + r.Time

		// Do the same for distance
		for m = 1; m < r.Distance; m *= 10 {
			// NOOP
		}
		race.Distance = race.Distance*m + r.Distance
	}

	fmt.Fprintf(w, "Part 2: %d\n", s.WaysToWin(race))
	return
}
