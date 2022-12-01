package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Cavern struct {
	octopuses [][]Octopus
}

type Octopus struct {
	Level   int
	Flashed bool
}

func (o *Octopus) Reset() {
	o.Level = 0
	o.Flashed = false
}

func NewCavern() *Cavern {
	return &Cavern{
		octopuses: make([][]Octopus, 0, 10),
	}
}

func (c *Cavern) AddLine(line string) {
	row := make([]Octopus, len(line))
	for i, ch := range strings.Split(line, "") {
		n, err := strconv.Atoi(ch)
		if err != nil {
			log.Fatalf("c.AddLine: %v", err)
		}
		row[i].Level = n
	}
	c.octopuses = append(c.octopuses, row)
}

func (c *Cavern) AllFlashing() bool {
	for i := 0; i < len(c.octopuses); i++ {
		for j := 0; j < len(c.octopuses[i]); j++ {
			if !c.octopuses[i][j].Flashed {
				return false
			}
		}
	}
	return true
}

func (c *Cavern) Flashes(steps int) (flashes int) {
	for i := 0; i < steps; i++ {
		flashes += c.Step()
	}
	return
}

func (c *Cavern) SimultaneousFlash() (step int) {
	for i := 0; i < 10000; i++ {
		c.Step()
		if c.AllFlashing() {
			step = i + 101
			return
		}
	}
	return
}

func (c *Cavern) Step() (flashes int) {
	// 0. Unset flash flag
	for i := range c.octopuses {
		for j := range c.octopuses[i] {
			c.octopuses[i][j].Flashed = false
		}
	}

	// 1. Increment the energy level by one
	for i := range c.octopuses {
		for j := range c.octopuses[i] {
			c.octopuses[i][j].Level++
		}
	}

	// 2. Any octopus with an energy level greater than 9 flashes
	for i := range c.octopuses {
		for j := range c.octopuses[i] {
			c.flash(i, j)
		}
	}

	// 3. Count flashes and reset to zero
	for i := range c.octopuses {
		for j := range c.octopuses[i] {
			if c.octopuses[i][j].Flashed {
				flashes++
				c.octopuses[i][j].Level = 0
			}
		}
	}

	return
}

func (c *Cavern) String() string {
	var builder strings.Builder
	for _, row := range c.octopuses {
		for _, octopus := range row {
			format := "%d"
			if octopus.Flashed {
				format = "\033[1;31m%d\033[0m"
			}
			fmt.Fprintf(&builder, format, octopus.Level)
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

func (c *Cavern) flash(i, j int) {
	if c.octopuses[i][j].Level < 10 {
		return
	}
	if c.octopuses[i][j].Flashed {
		return
	}
	c.octopuses[i][j].Flashed = true
	for _, y := range []int{i - 1, i, i + 1} {
		for _, x := range []int{j - 1, j, j + 1} {
			// Skip "self"
			if y == i && x == j {
				continue
			}
			// Skip out-of-bounds
			if y < 0 || y >= len(c.octopuses) || x < 0 || x >= len(c.octopuses[y]) {
				continue
			}
			c.octopuses[y][x].Level++
			c.flash(y, x)
		}
	}
}

func main() {
	cavern := NewCavern()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cavern.AddLine(scanner.Text())
	}

	// flashes := 0
	// fmt.Printf("Before any steps: (%d)\n%s\n", flashes, cavern)
	// flashes = cavern.Step()
	// fmt.Printf("After step 1: (%d)\n%s\n", flashes, cavern)
	// flashes = cavern.Step()
	// fmt.Printf("After step 2: (%d)\n%s", flashes, cavern)

	// Part 1
	log.Printf("Flashes: %d", cavern.Flashes(100))

	// Part 2
	log.Printf("Simultaneous Flash Step: %d", cavern.SimultaneousFlash())
}
