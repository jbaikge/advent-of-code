package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	From int
	To   int
}

func (r Range) Len() int {
	return r.To - r.From + 1
}

type Pair struct {
	A Range
	B Range
}

func NewPair(a, b, c, d int) (p Pair) {
	p.A.From = a
	p.A.To = b
	p.B.From = c
	p.B.To = d
	return
}

func (p Pair) FullOverlap() bool {
	a, b := p.A, p.B
	if a.Len() < b.Len() {
		a, b = b, a
	}
	return a.From <= b.From && a.To >= b.To
}

func (p Pair) PartialOverlap() bool {
	a, b := p.A, p.B
	if a.From > b.From {
		a, b = b, a
	}
	return a.To >= b.From
}

func part1(pairs []Pair) (count int) {
	for _, pair := range pairs {
		if pair.FullOverlap() {
			count++
		}
	}
	return
}

func part2(pairs []Pair) (count int) {
	for _, pair := range pairs {
		if pair.PartialOverlap() {
			count++
		}
	}
	return
}

func main() {
	pairs := make([]Pair, 0, 1000)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// Format: a-b,c-d
		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == '-' || r == ','
		})
		n := make([]int, len(fields))
		var err error
		for i, f := range fields {
			if n[i], err = strconv.Atoi(f); err != nil {
				log.Fatalf("could not convert %s to int: %v", f, err)
			}
		}
		pairs = append(pairs, NewPair(n[0], n[1], n[2], n[3]))
	}

	fmt.Printf("Part 1: %d\n", part1(pairs))
	fmt.Printf("Part 2: %d\n", part2(pairs))
}
