package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const OperationOld = -1

type Operation struct {
	A  int
	B  int
	Op byte
}

func NewOperation(raw string) (o Operation) {
	fields := strings.Fields(raw)
	if len(fields) != 5 {
		log.Fatalf("Invalid number of fields in %s", raw)
	}
	if v := fields[2]; v == "old" {
		o.A = OperationOld
	} else {
		o.A, _ = strconv.Atoi(v)
	}
	if v := fields[4]; v == "old" {
		o.B = OperationOld
	} else {
		o.B, _ = strconv.Atoi(v)
	}
	o.Op = fields[3][0]
	return
}

func (o Operation) Calculate(old int) int {
	a, b := o.A, o.B
	if a == OperationOld {
		a = old
	}
	if b == OperationOld {
		b = old
	}
	switch o.Op {
	case '+':
		return a + b
	case '*':
		return a * b
	default:
		log.Fatalf("No idea what to do with operation: %s", string(o.Op))
	}
	return 0
}

type Test struct {
	DivisibleBy int
	IfTrue      int // Throw to Monkey (num)
	IfFalse     int // Throw to Monkey (num)
}

func (t Test) ThrowTo(v int) int {
	if v%t.DivisibleBy == 0 {
		return t.IfTrue
	}
	return t.IfFalse
}

type Monkey struct {
	Num           int
	StartingItems []int
	Operation     Operation
	Test          Test
}

func NewMonkey(num int) *Monkey {
	return &Monkey{
		Num:           num,
		StartingItems: make([]int, 0, 8),
	}
}

func (m *Monkey) CalculateWorry(old int) int {
	return m.Operation.Calculate(old)
}

type MonkeyStats struct {
	Items     []int
	Inspected int
}

func part1(monkeys []*Monkey) (total int) {
	const Rounds = 20

	stats := make([]MonkeyStats, len(monkeys))
	for i, monkey := range monkeys {
		stats[i].Items = append(stats[i].Items, monkey.StartingItems...)
	}

	for round := 1; round <= Rounds; round++ {
		// fmt.Printf("Round %d:\n", round)
		for m, monkey := range monkeys {
			// fmt.Printf("  Monkey %d:\n", monkey.Num)
			for _, item := range stats[m].Items {
				stats[m].Inspected++
				// fmt.Printf("    Inspecting item: %d\n", item)
				newItem := monkey.CalculateWorry(item)
				// fmt.Printf("      Worry level increased to %d\n", newItem)
				newItem /= 3
				// fmt.Printf("      Bored, worry level decreased to %d\n", newItem)
				idx := monkey.Test.ThrowTo(newItem)
				// fmt.Printf("      Throwing to %d\n", idx)
				stats[idx].Items = append(stats[idx].Items, newItem)
			}
			// Empty current monkey's hands
			stats[m].Items = stats[m].Items[:0]
		}
	}

	inspected := make([]int, len(monkeys))
	for m, stat := range stats {
		inspected[m] = stat.Inspected
		fmt.Printf("Monkey %d: %+v\n", m, stat)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(inspected)))

	return inspected[0] * inspected[1]
}

func part2(monkeys []*Monkey) (total int) {
	const Rounds = 10000

	// LCM is easy to calculate since all the monkeys try to divide
	// by a number which is always prime.
	lcm := 1
	stats := make([]MonkeyStats, len(monkeys))
	for i, monkey := range monkeys {
		stats[i].Items = append(stats[i].Items, monkey.StartingItems...)
		lcm *= monkey.Test.DivisibleBy
	}

	for round := 1; round <= Rounds; round++ {
		// fmt.Printf("Round %d:\n", round)
		for m, monkey := range monkeys {
			// fmt.Printf("  Monkey %d:\n", monkey.Num)
			for _, item := range stats[m].Items {
				stats[m].Inspected++
				// fmt.Printf("    Inspecting item: %d\n", item)
				newItem := monkey.CalculateWorry(item)
				// fmt.Printf("      Worry level increased to %d\n", newItem)
				newItem %= lcm
				// fmt.Printf("      Decrease by modulo with %d: %d\n", lcm, newItem)
				idx := monkey.Test.ThrowTo(newItem)
				// fmt.Printf("      Throwing to %d\n", idx)
				stats[idx].Items = append(stats[idx].Items, newItem)
			}
			// Empty current monkey's hands
			stats[m].Items = stats[m].Items[:0]
		}
	}

	inspected := make([]int, len(monkeys))
	for m, stat := range stats {
		inspected[m] = stat.Inspected
		fmt.Printf("Monkey %d: %+v\n", m, stat)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(inspected)))

	return inspected[0] * inspected[1]
}

func main() {
	monkeys := make([]*Monkey, 0, 8)
	var monkey *Monkey
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Monkey"):
			fields := strings.Fields(line)
			num, _ := strconv.Atoi(fields[1][:len(fields[1])-1])
			monkey = NewMonkey(num)
			monkeys = append(monkeys, monkey)
		case strings.HasPrefix(line, "  Starting items"):
			split := strings.SplitN(line, ": ", 2)
			nums := strings.Split(split[1], ", ")
			for _, num := range nums {
				n, _ := strconv.Atoi(num)
				monkey.StartingItems = append(monkey.StartingItems, n)
			}
		case strings.HasPrefix(line, "  Operation"):
			split := strings.SplitN(line, ": ", 2)
			monkey.Operation = NewOperation(split[1])
		case strings.HasPrefix(line, "  Test"):
			split := strings.SplitN(line, ": ", 2)
			fields := strings.Fields(split[1])
			if fields[0] != "divisible" {
				log.Fatalf("Uh oh, unexpected test: %s", fields[0])
			}
			monkey.Test.DivisibleBy, _ = strconv.Atoi(fields[2])
		case strings.HasPrefix(line, "    If"):
			split := strings.SplitN(line, ": ", 2)
			fields := strings.Fields(split[1])
			num, _ := strconv.Atoi(fields[3])
			if strings.Contains(split[0], "true") {
				monkey.Test.IfTrue = num
			} else {
				monkey.Test.IfFalse = num
			}
		}
	}

	for _, monkey := range monkeys {
		fmt.Printf("%+v\n", monkey)
	}

	fmt.Printf("Part 1: %d\n", part1(monkeys))
	fmt.Printf("Part 2: %d\n", part2(monkeys))
}
