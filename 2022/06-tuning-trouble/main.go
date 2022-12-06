package main

import (
	"bufio"
	"fmt"
	"os"
)

func part1(line string) (pos int) {
	return uniqueWindow(line, 4)
}

func part2(line string) (pos int) {
	return uniqueWindow(line, 14)
}

func uniqueWindow(line string, width int) (pos int) {
	var window string
	var ch rune
	set := make(map[rune]struct{})
	for i := 0; i <= len(line)-width; i++ {
		window = line[i : width+i]
		for _, ch = range window {
			set[ch] = struct{}{}
		}
		if len(set) == width {
			return i + width
		}
		for ch = range set {
			delete(set, ch)
		}
	}
	return
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Part 1: %d\n", part1(line))
		fmt.Printf("Part 2: %d\n", part2(line))
	}
}
