package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Count struct {
	Zeros int
	Ones  int
}

func (c *Count) Increment(char byte) {
	if char == '0' {
		c.Zeros++
	} else {
		c.Ones++
	}
}

type Report struct {
	lines []string
}

func NewReport() *Report {
	return &Report{
		lines: make([]string, 0, 2000),
	}
}

func (r *Report) AddLine(line string) {
	r.lines = append(r.lines, line)
}

func (r *Report) Counts() (counts []Count) {
	counts = make([]Count, len(r.lines[0]))
	for i := range counts {
		c := &counts[i]
		for _, s := range r.lines {
			if s[i] == '0' {
				c.Zeros++
			} else {
				c.Ones++
			}
		}
	}
	return
}

func (r *Report) EpsilonRate() (int64, error) {
	counts := r.Counts()
	binary := make([]byte, len(counts))
	for i, c := range counts {
		binary[i] = '0'
		if c.Ones < c.Zeros {
			binary[i] = '1'
		}
	}
	return strconv.ParseInt(string(binary), 2, 64)
}

func (r *Report) GammaRate() (int64, error) {
	counts := r.Counts()
	binary := make([]byte, len(counts))
	for i, c := range counts {
		binary[i] = '0'
		if c.Ones > c.Zeros {
			binary[i] = '1'
		}
	}
	return strconv.ParseInt(string(binary), 2, 64)
}

func (r *Report) OxygenGeneratorRating() (int64, error) {
	prefix := make([]byte, 0, len(r.lines[0]))
	lines := r.lines
	for i := 0; i < len(r.lines[0]); i++ {
		// Count most frequent bit
		c := new(Count)
		for _, s := range lines {
			c.Increment(s[i])
		}

		// If ones >= zeros, choose 1
		ch := byte('1')
		if c.Zeros > c.Ones {
			ch = '0'
		}
		prefix = append(prefix, ch)

		keepLines := make([]string, 0, len(lines))
		for _, s := range lines {
			if strings.HasPrefix(s, string(prefix)) {
				keepLines = append(keepLines, s)
			}
		}

		if len(keepLines) == 0 {
			break
		}

		lines = keepLines
	}
	return strconv.ParseInt(lines[0], 2, 64)
}

func (r *Report) CO2ScrubberRating() (int64, error) {
	prefix := make([]byte, 0, len(r.lines[0]))
	lines := r.lines
	for i := 0; i < len(r.lines[0]); i++ {
		// Count most frequent bit
		c := new(Count)
		for _, s := range lines {
			c.Increment(s[i])
		}

		// If zeros <= ones, choose 0
		ch := byte('0')
		if c.Ones < c.Zeros {
			ch = '1'
		}
		prefix = append(prefix, ch)

		keepLines := make([]string, 0, len(lines))
		for _, s := range lines {
			if strings.HasPrefix(s, string(prefix)) {
				keepLines = append(keepLines, s)
			}
		}

		if len(keepLines) == 0 {
			break
		}

		lines = keepLines
	}
	return strconv.ParseInt(lines[0], 2, 64)
}

func main() {
	report := NewReport()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		report.AddLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	// Part 1
	gamma, err := report.GammaRate()
	if err != nil {
		log.Fatalf("Error during string conversion: %v", err)
	}
	log.Printf("Gamma: %d", gamma)

	epsilon, err := report.EpsilonRate()
	if err != nil {
		log.Fatalf("Error during string conversion: %v", err)
	}
	log.Printf("Epsilon: %d", epsilon)

	log.Printf("Gamma * Epsilon: %d", gamma*epsilon)

	// Part 2
	oxygen, err := report.OxygenGeneratorRating()
	if err != nil {
		log.Fatalf("Error during string conversion: %v", err)
	}
	log.Printf("Oxygen Generator Rating: %d", oxygen)

	co2, err := report.CO2ScrubberRating()
	if err != nil {
		log.Fatalf("Error during string conversion: %v", err)
	}
	log.Printf("CO2 Scrubber Rating: %d", co2)

	log.Printf("Oxygen Generator * CO2 Scrubber: %d", oxygen*co2)
}
