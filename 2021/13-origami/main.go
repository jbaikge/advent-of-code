package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Axis string

const (
	X Axis = "x"
	Y Axis = "y"
)

type Fold struct {
	Axis  Axis
	Value int
}

func dimensions(coords [][2]int) (width, height int) {
	for _, coord := range coords {
		if coord[0] > width {
			width = coord[0]
		}
		if coord[1] > height {
			height = coord[1]
		}
	}
	// zero indexed coords means max x and y need one more
	// to prevent index out of range errors
	width++
	height++
	return
}

func printGrid(grid [][]bool) {
	lines := make([]string, len(grid[0]))
	for x := range grid {
		for y := range grid[x] {
			ch := "."
			if grid[x][y] {
				ch = "#"
			}
			lines[y] += ch
		}
	}
	fmt.Println(strings.Join(lines, "\n"))
}

func count(grid [][]bool) (total int) {
	for _, col := range grid {
		for _, dot := range col {
			if dot {
				total++
			}
		}
	}
	return
}

func foldLeft(grid [][]bool, along int) (folded [][]bool) {
	width, height := along, len(grid[0])
	gridWidth := len(grid)
	folded = make([][]bool, width)
	for i := range folded {
		folded[i] = make([]bool, height)
	}

	for y := 0; y < height; y++ {
		// Copy existing dots
		for x := 0; x < width; x++ {
			folded[x][y] = grid[x][y]
		}
		// Transfer dots from fold
		for x := gridWidth - 1; x > width; x-- {
			newX := gridWidth - x - 1
			folded[newX][y] = folded[newX][y] || grid[x][y]
		}
	}
	return
}

func foldUp(grid [][]bool, along int) (folded [][]bool) {
	width, height := len(grid), along
	gridHeight := len(grid[0])
	folded = make([][]bool, width)
	for i := range folded {
		folded[i] = make([]bool, height)
	}
	for x := 0; x < width; x++ {
		// Copy existing dots
		for y := 0; y < height; y++ {
			folded[x][y] = grid[x][y]
		}
		// Transfer dots from fold
		for y := gridHeight - 1; y > height; y-- {
			newY := gridHeight - y - 1
			folded[x][newY] = folded[x][newY] || grid[x][y]
		}
	}
	return
}

func part1(coords [][2]int, folds []Fold) (n int) {
	width, height := dimensions(coords)
	grid := make([][]bool, width)
	for i := range grid {
		grid[i] = make([]bool, height)
	}
	for _, coord := range coords {
		grid[coord[0]][coord[1]] = true
	}
	fold := folds[0]
	if fold.Axis == Y {
		grid = foldUp(grid, fold.Value)
	} else {
		grid = foldLeft(grid, fold.Value)
	}
	return count(grid)
}

func part2(coords [][2]int, folds []Fold) (n int) {
	width, height := dimensions(coords)
	grid := make([][]bool, width)
	for i := range grid {
		grid[i] = make([]bool, height)
	}
	for _, coord := range coords {
		grid[coord[0]][coord[1]] = true
	}
	for _, fold := range folds {
		if fold.Axis == Y {
			grid = foldUp(grid, fold.Value)
		} else {
			grid = foldLeft(grid, fold.Value)
		}
	}
	printGrid(grid)
	return count(grid)
}

func main() {
	coords := make([][2]int, 0, 1000)
	folds := make([]Fold, 0, 10)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// Blank line separates coordinates from folds
		if line == "" {
			continue
		}
		// Parse fold directives
		if strings.HasPrefix(line, "fold") {
			fields := strings.Fields(line)
			def := strings.SplitN(fields[2], "=", 2)
			v, err := strconv.Atoi(def[1])
			if err != nil {
				log.Fatalf("could not convert %s to int: %v", def[1], err)
			}
			folds = append(folds, Fold{
				Axis:  Axis(def[0]),
				Value: v,
			})
			continue
		}
		// Parse coordinates
		coord := strings.SplitN(line, ",", 2)
		x, err := strconv.Atoi(coord[0])
		if err != nil {
			log.Fatalf("could not convert %s to int: %v", coord[0], err)
		}
		y, err := strconv.Atoi(coord[1])
		if err != nil {
			log.Fatalf("could not convert %s to int: %v", coord[1], err)
		}
		coords = append(coords, [2]int{x, y})
	}

	fmt.Printf("Part 1: %d\n", part1(coords, folds))
	fmt.Printf("Part 2: %d\n", part2(coords, folds))
}
