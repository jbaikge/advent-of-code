package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
)

type State struct {
	Increasing int
	Values     []int
}

func NewState() *State {
	return &State{
		Values: make([]int, 0, 2000),
	}
}

func (s *State) AddValue(v int) {
	s.Values = append(s.Values, v)
}

func (s *State) SingleIncreasing() (increasing int) {
	lastVal := math.MaxInt
	for _, n := range s.Values {
		if n > lastVal {
			increasing++
		}
		lastVal = n
	}
	return
}

func (s *State) WindowedIncreasing() (increasing int) {
	lastVal := math.MaxInt
	for i := 0; i < len(s.Values)-2; i++ {
		sum := 0
		for _, n := range s.Values[i : i+3] {
			sum += n
		}
		if sum > lastVal {
			increasing++
		}
		lastVal = sum
	}
	return
}

func main() {
	s := NewState()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("Error converting to integer: %v", err)
		}
		s.AddValue(n)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Reading stdin: %v", err)
	}
	log.Printf("Single Increasing: %d", s.SingleIncreasing())
	log.Printf("Windowed Increasing: %d", s.WindowedIncreasing())
}
