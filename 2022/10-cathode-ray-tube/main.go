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
	OpAddX = "addx"
	OpNoop = "noop"
)

type Instruction struct {
	Op    string
	Value int
}

func NewInstruction(line string) (inst Instruction) {
	fields := strings.Fields(line)
	inst.Op = fields[0]
	if inst.Op == OpAddX && len(fields) == 2 {
		inst.Value, _ = strconv.Atoi(fields[1])
	}
	return
}

func (i Instruction) Cycles() int {
	switch i.Op {
	case OpAddX:
		return 2
	case OpNoop:
		return 1
	default:
		return 0
	}
}

type Machine struct {
	PC int // Program Counter
	X  int // Register X
}

func NewMachine() *Machine {
	return &Machine{
		X: 1,
	}
}

func (m *Machine) Apply(inst Instruction) {
	if inst.Op == OpAddX {
		m.X += inst.Value
	}
}

func part1(instructions []Instruction) (total int) {
	breakpoints := make([]int, 0, 10)
	for i := 20; i < 1000; i += 40 {
		breakpoints = append(breakpoints, i)
	}
	bp := 0
	m := NewMachine()
	for _, inst := range instructions {
		cycles := inst.Cycles()
		for cycle := 0; cycle < cycles; cycle++ {
			// Bump program counter
			m.PC++
			// Read register (possibly) between instructions
			if breakpoint := breakpoints[bp]; m.PC == breakpoint {
				total += m.X * breakpoint
				bp++
			}
			// Apply instruction on last cycle
			if cycle == cycles-1 {
				m.Apply(inst)
			}
		}
	}
	return
}

func part2(instructions []Instruction) (screen string) {
	// Set up the screen pixels
	rows, cols := 6, 40
	pixels := make([][]byte, rows)
	for i := range pixels {
		pixels[i] = bytes.Repeat([]byte{'.'}, cols)
	}

	sprite := bytes.Repeat([]byte{'.'}, cols)
	illuminated := make([]int, 3)
	m := NewMachine()
	for _, inst := range instructions {
		cycles := inst.Cycles()
		for cycle := 0; cycle < cycles; cycle++ {
			row := m.PC / cols
			col := m.PC % cols
			// Illuminate sprite
			for i := -1; i < 2; i++ {
				if pos := m.X + i; pos >= 0 && pos < len(sprite) {
					sprite[pos] = '#'
					illuminated = append(illuminated, pos)
				}
			}
			// Transfer sprite state
			pixels[row][col] = sprite[col]
			// Debug
			// fmt.Printf("PC: %3d Row: %d Pos: %3d Sprite: %s\n", m.PC, row, m.X, string(sprite))
			// Bump program counter
			m.PC++
			// Apply instruction on last cycle
			if cycle == cycles-1 {
				m.Apply(inst)
			}
			// Darken lit sprite pixels
			for _, pos := range illuminated {
				sprite[pos] = '.'
			}
			illuminated = illuminated[0:0]
		}
	}

	for _, row := range pixels {
		screen += string(row) + "\n"
	}
	return
}

func main() {
	instructions := make([]Instruction, 0, 200)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		instructions = append(instructions, NewInstruction(scanner.Text()))
	}

	fmt.Printf("Part 1: %d\n", part1(instructions))
	fmt.Printf("Part 2:\n%s\n", part2(instructions))
}
