package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Expects positions to be presorted
// part 1 uses simple subtraction
func minFuel(positions []int) (minFuel int) {
	minFuel = math.MaxInt
	max := positions[len(positions)-1]
	for p := positions[0]; p < max; p++ {
		fuel := 0
		for _, pos := range positions {
			cost := pos - p
			if cost < 0 {
				cost *= -1
			}
			fuel += cost
		}
		if fuel < minFuel {
			minFuel = fuel
		}
	}
	return
}

// Expects positions to be presorted
// Part 2 uses a summation formula which pretty cleanly translates to
// (n^2 + n) / 2
func minFuelSummation(positions []int) (minFuel int) {
	minFuel = math.MaxInt
	max := positions[len(positions)-1]
	for p := positions[0]; p < max; p++ {
		fuel := 0
		for _, pos := range positions {
			n := pos - p
			if n < 0 {
				n *= -1
			}
			fuel += (n*n + n) / 2
		}
		if fuel < minFuel {
			minFuel = fuel
		}
	}
	return
}

func main() {
	var positions []int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), ",")
		positions = make([]int, 0, len(values))
		for _, value := range values {
			i, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Atoi failure with %s: %v", value, err)
			}
			positions = append(positions, i)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	log.Printf("Raw: %+v", positions)

	sort.Ints(positions)
	log.Printf("Sorted: %+v", positions)

	// Part 1
	start := time.Now()
	log.Printf("Min Fuel: %d (%s)", minFuel(positions), time.Since(start))

	// Part 2
	start = time.Now()
	log.Printf("Min Fuel (Summation): %d (%s)", minFuelSummation(positions), time.Since(start))
}
