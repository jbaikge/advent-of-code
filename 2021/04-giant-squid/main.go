package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Board struct {
	numbers         []int
	hits            []bool
	lastCallApplied int
}

func NewBoard() *Board {
	return &Board{
		numbers: make([]int, 0, 25),
		hits:    make([]bool, 25),
	}
}

func (b *Board) AddNumber(n int) {
	b.numbers = append(b.numbers, n)
}

func (b *Board) Call(number int) {
	for i, n := range b.numbers {
		if n == number {
			b.hits[i] = true
			b.lastCallApplied = number
		}
	}
}

func (b *Board) LastCall() int {
	return b.lastCallApplied
}

func (b *Board) HasVerticalWin() bool {
	// Columns
	for i := 0; i < 5; i++ {
		hits := 0
		// Rows
		for j := i; j < 25; j += 5 {
			if b.hits[j] {
				hits++
			}
		}
		if hits == 5 {
			return true
		}
	}
	return false
}

func (b *Board) HasHorizontalWin() bool {
	// Rows
	for i := 0; i < 25; i += 5 {
		hits := 0
		for _, hit := range b.hits[i : i+5] {
			if hit {
				hits++
			}
		}
		if hits == 5 {
			return true
		}
	}
	return false
}

func (b *Board) HasWin() bool {
	return b.HasHorizontalWin() || b.HasVerticalWin()
}

func (b *Board) UnmarkedSum() (sum int) {
	for i, n := range b.numbers {
		if !b.hits[i] {
			sum += n
		}
	}
	return
}

func (b *Board) ToString() string {
	str := ""
	for i := 0; i < 25; i += 5 {
		for _, n := range b.numbers[i : i+5] {
			str += fmt.Sprintf("%3d", n)
		}
		str += "\n"
	}
	return str
}

func firstToWin(boards []*Board, calls []int) *Board {
	for _, call := range calls {
		for _, board := range boards {
			board.Call(call)
			if board.HasWin() {
				return board
			}
		}
	}
	return nil
}

func lastToWin(boards []*Board, calls []int) *Board {
	var lastWinner *Board
	for _, call := range calls {
		for _, board := range boards {
			// Skip won boards
			if board.HasWin() {
				continue
			}
			board.Call(call)
			if board.HasWin() {
				lastWinner = board
			}
		}
	}
	return lastWinner
}

func main() {
	calls := make([]int, 0, 100)
	boards := make([]*Board, 0, 100)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(calls) == 0 {
			for _, call := range strings.Split(line, ",") {
				i, err := strconv.Atoi(call)
				if err != nil {
					log.Fatalf("Atoi failure: %v", err)
				}
				calls = append(calls, i)
			}
			continue
		}

		if line == "" {
			boards = append(boards, NewBoard())
			continue
		}

		for _, n := range strings.Fields(line) {
			i, err := strconv.Atoi(n)
			if err != nil {
				log.Fatalf("Atoi failure: %v", err)
			}
			boards[len(boards)-1].AddNumber(i)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	// Part 1
	first := firstToWin(boards, calls)
	if first == nil {
		log.Fatalf("No winning board found!")
	}
	log.Printf("*** First to win ***")
	log.Printf("Winning Call: %d", first.LastCall())
	log.Printf("Board Sum:    %d", first.UnmarkedSum())
	log.Printf("Final Score:  %d", first.LastCall()*first.UnmarkedSum())
	log.Printf("Board:\n%s", first.ToString())

	// Part 2
	last := lastToWin(boards, calls)
	if last == nil {
		log.Fatalf("No winning board found!")
	}
	log.Printf("*** Last to win ***")
	log.Printf("Winning Call: %d", last.LastCall())
	log.Printf("Board Sum:    %d", last.UnmarkedSum())
	log.Printf("Final Score:  %d", last.LastCall()*last.UnmarkedSum())
	log.Printf("Board:\n%s", last.ToString())

}
