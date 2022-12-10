package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	Up    = 'U'
	Right = 'R'
	Down  = 'D'
	Left  = 'L'
)

type Motion struct {
	Direction byte
	Distance  int
}

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func move(knot *Point, m Motion) {
	switch m.Direction {
	case Up:
		knot.Y++
	case Down:
		knot.Y--
	case Right:
		knot.X++
	case Left:
		knot.X--
	}
}

func tug(head, tail *Point) {
	switch {
	// H/T in the same row
	case head.X == tail.X:
		switch head.Y - tail.Y {
		case 2:
			tail.Y++
		case -2:
			tail.Y--
		}
	// H/T in the same column
	case head.Y == tail.Y:
		switch head.X - tail.X {
		case 2:
			tail.X++
		case -2:
			tail.X--
		}
	// Tail needs to move diagonally to catch up to Head
	case math.Pow(float64(head.X-tail.X), 2)+math.Pow(float64(head.Y-tail.Y), 2) > 2:
		switch head.X - tail.X {
		case 2:
			tail.X++
			if head.Y > tail.Y {
				tail.Y++
			} else {
				tail.Y--
			}
		case -2:
			tail.X--
			if head.Y > tail.Y {
				tail.Y++
			} else {
				tail.Y--
			}
		}
		switch head.Y - tail.Y {
		case 2:
			tail.Y++
			if head.X > tail.X {
				tail.X++
			} else {
				tail.X--
			}
		case -2:
			tail.Y--
			if head.X > tail.X {
				tail.X++
			} else {
				tail.X--
			}
		}
	default:
		// NOOP - don't know why this is necessary
	}
}

func walk(knotCount int, motions []Motion) int {
	if knotCount < 1 {
		log.Fatalf("knot count must be at least 1; %d given", knotCount)
	}
	knots := make([]Point, knotCount)
	visted := make(map[Point]bool)
	for _, m := range motions {
		for step := 0; step < m.Distance; step++ {
			move(&knots[0], m)
			for k := 1; k < len(knots); k++ {
				tug(&knots[k-1], &knots[k])
			}
			visted[knots[len(knots)-1]] = true
		}
	}
	return len(visted)
}

func part1(motions []Motion) (visited int) {
	return walk(2, motions)
}

func part2(motions []Motion) (visited int) {
	return walk(10, motions)
}

func main() {
	motions := make([]Motion, 0, 2000)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		m := Motion{
			Direction: fields[0][0],
		}
		m.Distance, _ = strconv.Atoi(fields[1])
		motions = append(motions, m)
	}

	fmt.Printf("Part 1: %d\n", part1(motions))
	fmt.Printf("Part 2: %d\n", part2(motions))
}
