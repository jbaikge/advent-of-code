package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Command struct {
	Direction string
	Units     int
}

func NewCommand(line string) (c Command, err error) {
	fields := strings.Fields(line)
	if l := len(fields); l != 2 {
		err = fmt.Errorf("invalid field count: %d for line: `%s'", l, line)
		return
	}

	c = Command{
		Direction: fields[0],
	}
	c.Units, err = strconv.Atoi(fields[1])
	return
}

func (c Command) X() int {
	if c.Direction == "forward" {
		return c.Units
	}
	return 0
}

func (c Command) Depth() int {
	if c.Direction == "up" {
		return -c.Units
	}
	if c.Direction == "down" {
		return c.Units
	}
	return 0
}

type Commands struct {
	commands []Command
}

func NewCommands() *Commands {
	return &Commands{
		commands: make([]Command, 0, 2000),
	}
}

func (cs *Commands) AddCommand(c Command) {
	cs.commands = append(cs.commands, c)
}

func (cs Commands) Position() (x int, depth int) {
	for _, c := range cs.commands {
		x += c.X()
		depth += c.Depth()
	}
	return
}

func (cs Commands) PositionWithAim() (x int, depth int) {
	aim := 0
	for _, c := range cs.commands {
		aim += c.Depth()
		if cx := c.X(); cx > 0 {
			x += cx
			depth += aim * cx
		}
	}
	return
}

func main() {
	commands := NewCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		c, err := NewCommand(line)
		if err != nil {
			log.Fatalf("Error creating command: %v", err)
		}
		commands.AddCommand(c)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during scanning: %v", err)
	}

	x, depth := commands.Position()
	log.Printf("Final Position: %d, %d (%d)", x, depth, x*depth)

	x, depth = commands.PositionWithAim()
	log.Printf("Position With Aim: %d, %d (%d)", x, depth, x*depth)
}
