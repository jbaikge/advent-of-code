package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

type Subsystem struct {
	lines    []string
	replacer *strings.Replacer
	opening  string
	closing  string
	cScores  []int // corrupt scores
	iScores  []int // incomplete scores
}

func NewSubsystem() *Subsystem {
	s := &Subsystem{
		lines:   make([]string, 0, 100),
		opening: "([{<",
		closing: ")]}>",
		cScores: []int{3, 57, 1197, 25137},
		iScores: []int{1, 2, 3, 4},
	}
	oldnew := make([]string, 0, len(s.opening)*2)
	for i := range s.opening {
		old := s.opening[i:i+1] + s.closing[i:i+1]
		oldnew = append(oldnew, old, "")
	}
	s.replacer = strings.NewReplacer(oldnew...)
	return s
}

func (s *Subsystem) AddLine(line string) {
	s.lines = append(s.lines, line)
}

func (s *Subsystem) CorruptedScore() (score int) {
	for _, reduced := range s.reduced() {
		if corrupt, lineScore := s.isCorrupted(reduced); corrupt {
			score += lineScore
		}
	}
	return
}

func (s *Subsystem) IncompleteScore() (score int) {
	scores := make([]int, 0, len(s.lines))
	for _, reduced := range s.reduced() {
		if incomplete, lineScore := s.isIncomplete(reduced); incomplete {
			scores = append(scores, lineScore)
		}
	}
	sort.Ints(scores)
	middle := len(scores) / 2
	return scores[middle]
}

// A corrupted line is one where a chunk closes with the wrong character
func (s *Subsystem) isCorrupted(line string) (corrupted bool, score int) {
	for _, ch := range line {
		for i, closing := range s.closing {
			if ch != closing {
				continue
			}
			return true, s.cScores[i]
		}
	}
	return
}

// An incomplete (reduced) line is one with only opening characters
func (s *Subsystem) isIncomplete(line string) (incomplete bool, score int) {
	if corrupt, _ := s.isCorrupted(line); corrupt {
		return
	}
	incomplete = true
	for i := len(line) - 1; i >= 0; i-- {
		ch := rune(line[i])
		for j, open := range s.opening {
			if ch == open {
				score *= 5
				score += s.iScores[j]
			}
		}
	}
	return
}

func (s *Subsystem) reduced() (reduced []string) {
	reduced = make([]string, len(s.lines))
	for i, line := range s.lines {
		reduced[i] = s.reduceLine(line)
	}
	return
}

func (s *Subsystem) reduceLine(in string) (out string) {
	out = s.replacer.Replace(in)
	if in == out {
		return
	}
	return s.reduceLine(out)
}

func main() {
	subsystem := NewSubsystem()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		subsystem.AddLine(scanner.Text())
	}

	// Part 1
	log.Printf("Corrupted lines score: %d", subsystem.CorruptedScore())
	// Part 2
	log.Printf("Incomplete line score: %d", subsystem.IncompleteScore())
}
