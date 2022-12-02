package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	ScoreLose = 0
	ScoreDraw = 3
	ScoreWin  = 6
)

type Shape struct {
	Name  string
	Value int
}

var (
	Rock     = Shape{Name: "Rock", Value: 1}
	Paper    = Shape{Name: "Paper", Value: 2}
	Scissors = Shape{Name: "Scissors", Value: 3}
)

type Strategy struct {
	Map map[byte]Shape
}

func (s Strategy) Outcome(opponent byte, self byte) (score int) {
	opponentShape := s.Map[opponent]
	selfShape := s.Map[self]

	switch {
	case opponentShape == selfShape:
		return selfShape.Value + ScoreDraw
	// Rock beats Scissors
	case opponentShape == Scissors && selfShape == Rock:
		return selfShape.Value + ScoreWin
	// Scissors beats Paper
	case opponentShape == Paper && selfShape == Scissors:
		return selfShape.Value + ScoreWin
	// Paper beats Rock
	case opponentShape == Rock && selfShape == Paper:
		return selfShape.Value + ScoreWin
	default:
		return selfShape.Value + ScoreLose
	}
}

func main() {
	rounds := make([][2]byte, 0, 2500)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		rounds = append(rounds, [2]byte{line[0], line[2]})
	}

	part1 := Strategy{
		Map: map[byte]Shape{
			'A': Rock,
			'B': Paper,
			'C': Scissors,
			'X': Rock,
			'Y': Paper,
			'Z': Scissors,
		},
	}
	part1Total := 0
	for _, round := range rounds {
		outcome := part1.Outcome(round[0], round[1])
		fmt.Printf("%s %s %d\n", string(round[0]), string(round[1]), outcome)
		part1Total += outcome
	}
	fmt.Printf("Part1: %d\n", part1Total)
}
