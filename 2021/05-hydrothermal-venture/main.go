package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Line struct {
	A Point
	B Point
}

func (l Line) IsVertical() bool {
	return l.A.Y == l.B.Y
}

func (l Line) IsHorizontal() bool {
	return l.A.X == l.B.X
}

func (l Line) IsStraight() bool {
	return l.IsHorizontal() || l.IsVertical()
}

func (l Line) Points() (points []Point) {
	deltaX := l.A.X - l.B.X
	if deltaX < 0 {
		deltaX *= -1
	}

	deltaY := l.A.Y - l.B.Y
	if deltaY < 0 {
		deltaY *= -1
	}

	delta := deltaX
	if deltaY > delta {
		delta = deltaY
	}

	// +1 to include the last point
	points = make([]Point, 0, delta+1)
	for i := 0; i < cap(points); i++ {
		xMul, yMul := 0, 0

		if l.A.X > l.B.X {
			xMul = -1
		} else if l.A.X < l.B.X {
			xMul = 1
		}

		if l.A.Y > l.B.Y {
			yMul = -1
		} else if l.A.Y < l.B.Y {
			yMul = 1
		}

		points = append(points, Point{
			X: l.A.X + i*xMul,
			Y: l.A.Y + i*yMul,
		})
	}

	return points
}

func overlappingPoints(lines []Line) int {
	type metric struct {
		Point Point
		Count int
	}

	metrics := make([]*metric, 0, 1000)
	for _, line := range lines {
		for _, point := range line.Points() {
			add := true
			for _, m := range metrics {
				if m.Point.X == point.X && m.Point.Y == point.Y {
					m.Count++
					add = false
					break
				}
			}
			if add {
				metrics = append(metrics, &metric{
					Point: point,
					Count: 1,
				})
			}
		}
	}

	overlapping := 0
	for _, m := range metrics {
		if m.Count > 1 {
			overlapping++
		}
	}
	return overlapping
}

func main() {
	lines := make([]Line, 0, 500)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Fields(text)
		numbers := make([]int, 0, 4)
		for _, field := range []string{fields[0], fields[2]} {
			for _, value := range strings.Split(field, ",") {
				i, err := strconv.Atoi(value)
				if err != nil {
					log.Fatalf("Atoi failure with %s: %v", value, err)
				}
				numbers = append(numbers, i)
			}
		}
		lines = append(lines, Line{
			A: Point{
				X: numbers[0],
				Y: numbers[1],
			},
			B: Point{
				X: numbers[2],
				Y: numbers[3],
			},
		})
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	// Part 1
	straightLines := make([]Line, 0, len(lines))
	for _, line := range lines {
		if line.IsStraight() {
			straightLines = append(straightLines, line)
		}
	}
	log.Printf("Straight Overlapping: %d", overlappingPoints(straightLines))

	// Part 2
	log.Printf("All Overlapping: %d", overlappingPoints(lines))
}
