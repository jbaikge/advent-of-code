package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type Polymer struct {
	Template string
	Pairs    map[string]byte
}

/*
 * Old naive solution:
func (p Polymer) Apply(steps int) (result string) {
	str := []byte(p.Template)
	chars := make([]byte, 0)
	combined := make([]byte, 0)
	for step := 0; step < steps; step++ {
		chars = chars[0:0]
		combined = combined[0:0]

		fmt.Printf("[%2d] Str len: %d\n", step+1, len(str))
		for i := 0; i < len(str)-1; i++ {
			key := str[i : i+2]
			chars = append(chars, p.Pairs[string(key)])
		}

		for i, ch := range chars {
			combined = append(combined, str[i])
			combined = append(combined, ch)
		}
		combined = append(combined, str[len(str)-1])
		str = make([]byte, len(combined))
		copy(str, combined)
	}
	return string(combined)
}

func (p Polymer) CharDiff(steps int) int {
	str := p.Apply(steps)
	counts := make(map[rune]int)
	for _, ch := range str {
		counts[ch]++
	}
	min := len(str)
	max := 0
	for _, count := range counts {
		if count > max {
			max = count
		}
		if count < min {
			min = count
		}
	}
	return max - min
}
*/

type Polymerization struct {
	Polymer      Polymer
	Steps        int
	PairCounts   map[string]int
	LetterCounts map[byte]int
}

func NewPolymerization(steps int, polymer Polymer) (p Polymerization) {
	p = Polymerization{
		Polymer:      polymer,
		Steps:        steps,
		LetterCounts: make(map[byte]int),
	}
	// Initialize pair counts
	p.PairCounts = p.NewPairCounts()
	for i := 0; i < len(p.Polymer.Template)-1; i++ {
		p.PairCounts[p.Polymer.Template[i:i+2]]++
	}
	// Initialize letter counts
	for _, ch := range p.Polymer.Template {
		p.LetterCounts[byte(ch)]++
	}
	return
}

func (p Polymerization) NewPairCounts() (counts map[string]int) {
	counts = make(map[string]int)
	for pair := range p.Polymer.Pairs {
		counts[pair] = 0
	}
	return
}

func (p Polymerization) Apply() {
	for i := 0; i < p.Steps; i++ {
		p.Step()
	}
}

func (p Polymerization) MinMax() (min, max int) {
	min = math.MaxInt
	for _, count := range p.LetterCounts {
		if count > max {
			max = count
		}
		if count < min {
			min = count
		}
	}
	return
}

func (p Polymerization) Step() {
	newCounts := p.NewPairCounts()
	for pair, count := range p.PairCounts {
		ch := p.Polymer.Pairs[pair]
		left := string([]byte{pair[0], ch})
		right := string([]byte{ch, pair[1]})
		newCounts[left] += count
		newCounts[right] += count
		p.LetterCounts[ch] += count
	}
	// Copy values
	for pair, count := range newCounts {
		p.PairCounts[pair] = count
	}
}

func part1(polymer Polymer) (result int) {
	p := NewPolymerization(10, polymer)
	p.Apply()
	min, max := p.MinMax()
	return max - min
}

func part2(polymer Polymer) (result int) {
	p := NewPolymerization(40, polymer)
	p.Apply()
	min, max := p.MinMax()
	return max - min
}

func main() {
	polymer := Polymer{
		Pairs: make(map[string]byte),
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if polymer.Template == "" {
			polymer.Template = line
		}
		if strings.Contains(line, "->") {
			fields := strings.Fields(line)
			polymer.Pairs[fields[0]] = fields[2][0]
		}
	}

	fmt.Printf("Part 1: %d\n", part1(polymer))
	fmt.Printf("Part 2: %d\n", part2(polymer))
}
