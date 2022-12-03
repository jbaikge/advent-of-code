package main

import (
	"bufio"
	"fmt"
	"os"
)

var priorities map[rune]int

func init() {
	priorities = make(map[rune]int)
	for i, j := 'a', 'A'; i <= 'z'; i, j = i+1, j+1 {
		priorities[i] = int(i-'a') + 1
		priorities[j] = int(j-'A') + 27
	}
}

func sackCommon(a, b string) (common rune) {
	for _, aChar := range a {
		for _, bChar := range b {
			if aChar == bChar {
				return aChar
			}
		}
	}
	panic("no common rune found")
}

// Find item common to both rucksacks and total the item's priority
func part1(lines []string) (total int) {
	for _, line := range lines {
		midpoint := len(line) / 2
		sack1, sack2 := line[:midpoint], line[midpoint:]
		common := sackCommon(sack1, sack2)
		priority := priorities[common]
		total += priority
	}
	return
}

func groupCommmon(a, b, c string) (common rune) {
	for _, aChar := range a {
		for _, bChar := range b {
			if aChar != bChar {
				continue
			}
			for _, cChar := range c {
				if aChar == cChar {
					return aChar
				}
			}
		}
	}
	panic("no common rune found")
}

func part2(lines []string) (total int) {
	for i := 0; i < len(lines); i += 3 {
		group := lines[i : i+3]
		common := groupCommmon(group[0], group[1], group[2])
		priority := priorities[common]
		total += priority
	}
	return
}

func main() {
	lines := make([]string, 0, 300)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	fmt.Printf("Part 1: %d\n", part1(lines))
	fmt.Printf("Part 2: %d\n", part2(lines))
}
