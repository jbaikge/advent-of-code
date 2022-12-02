package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	ScoreLose     = 0
	ScoreDraw     = 3
	ScoreWin      = 6
	ShapeRock     = 1
	ShapePaper    = 2
	ShapeScissors = 3
)

// A, B, C = Rock, Paper, Scissors
// X, Y, Z assumed to by Rock, Paper, Scissors
// Score each round accordingly
func part1(rounds [][2]byte) (total int) {
	roundMap := map[byte]map[byte]int{
		// Opponent plays Rock
		'A': {
			'X': ShapeRock + ScoreDraw,
			'Y': ShapePaper + ScoreWin,
			'Z': ShapeScissors + ScoreLose,
		},
		// Opponent plays Paper
		'B': {
			'X': ShapeRock + ScoreLose,
			'Y': ShapePaper + ScoreDraw,
			'Z': ShapeScissors + ScoreWin,
		},
		// Opponent plays Scissors
		'C': {
			'X': ShapeRock + ScoreWin,
			'Y': ShapePaper + ScoreLose,
			'Z': ShapeScissors + ScoreDraw,
		},
	}

	for _, round := range rounds {
		score := roundMap[round[0]][round[1]]
		total += score
	}

	return
}

// A, B, C = Rock, Paper, Scissors
// X, Y, Z = Lose, Draw, Win
// Score each round accordingly
func part2(rounds [][2]byte) (total int) {
	roundMap := map[byte]map[byte]int{
		// Opponent plays Rock
		'A': {
			'X': ShapeScissors + ScoreLose,
			'Y': ShapeRock + ScoreDraw,
			'Z': ShapePaper + ScoreWin,
		},
		// Opponent plays Paper
		'B': {
			'X': ShapeRock + ScoreLose,
			'Y': ShapePaper + ScoreDraw,
			'Z': ShapeScissors + ScoreWin,
		},
		// Opponent plays Scissors
		'C': {
			'X': ShapePaper + ScoreLose,
			'Y': ShapeScissors + ScoreDraw,
			'Z': ShapeRock + ScoreWin,
		},
	}

	for _, round := range rounds {
		score := roundMap[round[0]][round[1]]
		total += score
	}

	return
}

func main() {
	rounds := make([][2]byte, 0, 2500)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		rounds = append(rounds, [2]byte{line[0], line[2]})
	}

	fmt.Printf("Part 1: %d\n", part1(rounds))
	fmt.Printf("Part 2: %d\n", part2(rounds))
}
