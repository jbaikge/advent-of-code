package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Stack struct {
	Crates []string
}

func (s *Stack) Pop() string {
	length := len(s.Crates)
	if length == 0 {
		log.Fatalf("empty stack of crates!")
	}
	last := s.Crates[length-1]
	s.Crates = s.Crates[:length-1]
	return last
}

func (s *Stack) Push(crate string) {
	s.Crates = append(s.Crates, crate)
}

func (s *Stack) Top() string {
	return s.Crates[len(s.Crates)-1]
}

type Move struct {
	Quantity int
	From     int
	To       int
}

// line expects the following format:
// move 1 from 2 to 1
func NewMove(line string) (m Move) {
	fields := strings.Fields(line)
	if len(fields) != 6 {
		log.Fatalf("invalid move line: %s", line)
	}
	m.Quantity, _ = strconv.Atoi(fields[1])
	m.From, _ = strconv.Atoi(fields[3])
	m.To, _ = strconv.Atoi(fields[5])
	return
}

type Ship struct {
	Stacks      []Stack
	Moves       []Move
	CrateBuffer []string
}

func NewShip() *Ship {
	return &Ship{
		Moves:       make([]Move, 0, 500),
		CrateBuffer: make([]string, 0, 8),
	}
}

func (s *Ship) AddMove(line string) {
	s.Moves = append(s.Moves, NewMove(line))
}

// Adds a line from the input into a buffer, trimming the space on the right
// for consistency since my editor is also trimming the space on the right.
func (s *Ship) BufferCrates(line string) {
	s.CrateBuffer = append(s.CrateBuffer, strings.TrimRight(line, " "))
}

// Make copies of the ship in order to operate the crane differently between
// parts 1 & 2
func (src Ship) Copy() (dst Ship) {
	dst.Moves = make([]Move, len(src.Moves))
	copy(dst.Moves, src.Moves)
	dst.Stacks = make([]Stack, len(src.Stacks))
	copy(dst.Stacks, src.Stacks)
	for i, stack := range src.Stacks {
		dst.Stacks[i].Crates = make([]string, len(stack.Crates))
		copy(dst.Stacks[i].Crates, stack.Crates)
	}
	return
}

// Parses the line with 1 2 3 ...
func (s *Ship) InitStacks(line string) {
	fields := strings.Fields(strings.TrimSpace(line))
	stacks, _ := strconv.Atoi(fields[len(fields)-1])
	s.Stacks = make([]Stack, stacks)
	s.parseCrateBuffer()
}

// CrateMover 9000 only moves one crate at a time
func (s *Ship) MoveCrates9000() {
	for _, move := range s.Moves {
		from := move.From - 1
		to := move.To - 1
		for i := 0; i < move.Quantity; i++ {
			s.Stacks[to].Push(s.Stacks[from].Pop())
		}
	}
}

// CrateMover 9001 can move entire stacks of crates at a time
func (s *Ship) MoveCrates9001() {
	for _, move := range s.Moves {
		from := move.From - 1
		to := move.To - 1
		buffer := make([]string, move.Quantity)
		for i := move.Quantity - 1; i >= 0; i-- {
			buffer[i] = s.Stacks[from].Pop()
		}
		for _, crate := range buffer {
			s.Stacks[to].Push(crate)
		}
	}
}

func (s Ship) String() (str string) {
	maxLen := 0
	for _, stack := range s.Stacks {
		if length := len(stack.Crates); length > maxLen {
			maxLen = length
		}
	}

	for i := 0; i < maxLen; i++ {
		crates := make([]string, 0, len(s.Stacks))
		for _, stack := range s.Stacks {
			if i+1 > len(stack.Crates) {
				crates = append(crates, "   ")
				continue
			}
			crates = append(crates, stack.Crates[i])
		}
		str = strings.Join(crates, " ") + "\n" + str
	}
	return
}

func (s Ship) Top() (top string) {
	for _, stack := range s.Stacks {
		top += stack.Top()[1:2]
	}
	return
}

func (s *Ship) parseCrateBuffer() {
	for i := len(s.CrateBuffer) - 1; i >= 0; i-- {
		line := s.CrateBuffer[i]
		// Put placeholders in between stacks
		line = strings.ReplaceAll(line, "]     ", "] --- ")
		// Put placeholders along the leading edge
		line = strings.ReplaceAll(line, "    ", "--- ")
		crates := strings.Fields(line)
		for i, crate := range crates {
			if crate == "---" {
				continue
			}
			s.Stacks[i].Push(crate)
		}
	}
}

func part1(ship Ship) (top string) {
	fmt.Println(ship)
	ship.MoveCrates9000()
	fmt.Println(ship)
	return ship.Top()
}

func part2(ship Ship) (top string) {
	fmt.Println(ship)
	ship.MoveCrates9001()
	fmt.Println(ship)
	return ship.Top()
}

func main() {
	ship := NewShip()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ContainsRune(line, '[') {
			ship.BufferCrates(line)
			continue
		}
		if strings.HasPrefix(line, " 1") {
			ship.InitStacks(line)
			continue
		}
		if strings.HasPrefix(line, "move") {
			ship.AddMove(line)
			continue
		}
	}

	fmt.Printf("Part 1: %s\n", part1(ship.Copy()))
	fmt.Printf("Part 2: %s\n", part2(ship.Copy()))
}
