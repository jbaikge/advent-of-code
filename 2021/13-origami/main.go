package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	XAxis = 'x'
	YAxis = 'y'
)

type Fold struct {
	Axis  byte
	Value int
}

type Grid struct {
	Points map[[2]int]bool
	Folds  []Fold
}

func NewGrid() *Grid {
	return &Grid{
		Points: make(map[[2]int]bool),
		Folds:  make([]Fold, 0, 10),
	}
}

func (g *Grid) AddFold(axis byte, value int) {
	g.Folds = append(g.Folds, Fold{
		Axis:  axis,
		Value: value,
	})
}

func (src *Grid) Copy() (dst *Grid) {
	dst = &Grid{
		Points: make(map[[2]int]bool),
		Folds:  make([]Fold, len(src.Folds)),
	}
	copy(dst.Folds, src.Folds)
	for k, v := range src.Points {
		dst.Points[k] = v
	}
	return
}

func (g *Grid) SetPoint(x, y int) {
	g.Points[[2]int{x, y}] = true
}

func (g Grid) Size() (width, height int) {
	for coords := range g.Points {
		if coords[0] > width {
			width = coords[0]
		}
		if coords[1] > height {
			height = coords[1]
		}
	}
	width++
	height++
	return
}

func (g *Grid) Fold(times int) {
	if times == -1 {
		times = len(g.Folds)
	}
	for _, fold := range g.Folds[:times] {
		for point := range g.Points {
			if fold.Axis == XAxis && point[0] > fold.Value {
				delete(g.Points, point)
				g.Points[[2]int{fold.Value*2 - point[0], point[1]}] = true
			}
			if fold.Axis == YAxis && point[1] > fold.Value {
				delete(g.Points, point)
				g.Points[[2]int{point[0], fold.Value*2 - point[1]}] = true
			}
		}
	}
}

func (g Grid) String() (s string) {
	width, height := g.Size()
	grid := make([][]byte, height)
	for i := range grid {
		grid[i] = bytes.Repeat([]byte{'.'}, width)
	}
	for coords, active := range g.Points {
		// 	s += fmt.Sprintf("(%3d, %3d)\n", coords[0], coords[1])
		var ch byte = '#'
		if !active {
			ch = '%'
		}
		grid[coords[1]][coords[0]] = ch
	}
	for _, row := range grid {
		s += string(row) + "\n"
	}
	return
}

func part1(grid *Grid) (count int) {
	grid.Fold(1)
	return len(grid.Points)
}

func part2(grid *Grid) string {
	grid.Fold(-1)
	return grid.String()
}

func main() {
	grid := NewGrid()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ContainsRune(line, ',') {
			coords := strings.SplitN(line, ",", 2)
			x, _ := strconv.Atoi(coords[0])
			y, _ := strconv.Atoi(coords[1])
			grid.SetPoint(x, y)
			continue
		}
		if strings.HasPrefix(line, "fold") {
			fields := strings.Fields(line)
			parts := strings.SplitN(fields[2], "=", 2)
			value, _ := strconv.Atoi(parts[1])
			grid.AddFold(parts[0][0], value)
		}
	}

	fmt.Printf("Part 1: %d\n", part1(grid.Copy()))
	fmt.Printf("Part 2:\n%s\n", part2(grid.Copy()))
}
