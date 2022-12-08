package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Grid struct {
	Rows  int
	Cols  int
	Trees map[[2]int]int
}

func NewGrid() Grid {
	return Grid{
		Trees: make(map[[2]int]int),
	}
}

func (g Grid) Tree(x, y int) int {
	return g.Trees[[2]int{x, y}]
}

func part1(g Grid) (total int) {
	total = (g.Rows - 1 + g.Cols - 1) * 2
	var north, east, south, west bool
	for pos, height := range g.Trees {
		x, y := pos[0], pos[1]
		// Skip edges
		if x == 0 || x == g.Cols-1 || y == 0 || y == g.Rows-1 {
			continue
		}

		north, east, south, west = true, true, true, true

		// Check North
		for i := 0; i < y; i++ {
			if g.Tree(x, i) >= height {
				north = false
				break
			}
		}

		// Check East
		for i := x + 1; i < g.Cols; i++ {
			if g.Tree(i, y) >= height {
				east = false
				break
			}
		}

		// Check South
		for i := y + 1; i < g.Rows; i++ {
			if g.Tree(x, i) >= height {
				south = false
				break
			}
		}

		// Check West
		for i := 0; i < x; i++ {
			if g.Tree(i, y) >= height {
				west = false
				break
			}
		}

		if north || east || south || west {
			total++
		}
	}
	return
}

func part2(g Grid) (score int) {
	var north, east, south, west int
	for pos, height := range g.Trees {
		x, y := pos[0], pos[1]

		// Skip edges - they will always have a zero value and
		// cause a zero score
		if x == 0 || x == g.Cols-1 || y == 0 || y == g.Rows-1 {
			continue
		}

		north, east, south, west = 0, 0, 0, 0

		// Look North
		for i := y - 1; i >= 0; i-- {
			north++
			if g.Tree(x, i) >= height {
				break
			}
		}

		// Look East
		for i := x + 1; i < g.Cols; i++ {
			east++
			if g.Tree(i, y) >= height {
				break
			}
		}

		// Look South
		for i := y + 1; i < g.Rows; i++ {
			south++
			if g.Tree(x, i) >= height {
				break
			}
		}

		// Look West
		for i := x - 1; i >= 0; i-- {
			west++
			if g.Tree(i, y) >= height {
				break
			}
		}

		if treeScore := north * east * south * west; treeScore > score {
			score = treeScore
		}
	}
	return
}

func main() {
	grid := NewGrid()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		grid.Cols = len(line)
		for col, ch := range strings.Split(line, "") {
			n, _ := strconv.Atoi(ch)
			grid.Trees[[2]int{col, grid.Rows}] = n
		}
		grid.Rows++
	}

	fmt.Printf("Part 1: %d\n", part1(grid))
	fmt.Printf("Part 2: %d\n", part2(grid))
}
