package scratchcards

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Ticket struct {
	Revealed []int
	Winning  []int
}

type Solution struct {
	Tickets []Ticket
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Tickets = make([]Ticket, 0, 250)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		colon := strings.IndexByte(line, ':')
		parts := strings.Split(line[colon+1:], "|")

		revealed := strings.Fields(strings.TrimSpace(parts[0]))
		winning := strings.Fields(strings.TrimSpace(parts[1]))

		ticket := Ticket{
			Revealed: make([]int, len(revealed)),
			Winning:  make([]int, len(winning)),
		}
		for i, r := range revealed {
			if ticket.Revealed[i], err = strconv.Atoi(r); err != nil {
				return
			}
		}
		for i, w := range winning {
			if ticket.Winning[i], err = strconv.Atoi(w); err != nil {
				return
			}
		}

		slices.Sort[[]int](ticket.Revealed)
		slices.Sort[[]int](ticket.Winning)
		s.Tickets = append(s.Tickets, ticket)
	}
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	var sum int
	for _, ticket := range s.Tickets {
		var matches int
		// just aliases to cut down on typing
		rev := ticket.Revealed
		win := ticket.Winning
		for r, w := 0, 0; r < len(rev) && w < len(win); {
			if rev[r] == win[w] {
				matches++
				r++
				w++
				continue
			}
			if rev[r] > win[w] {
				w++
				continue
			}
			if win[w] > rev[r] {
				r++
				continue
			}
		}
		if matches == 0 {
			continue
		}
		sum += int(math.Pow(2, float64(matches-1)))
	}
	fmt.Fprintf(w, "Part 1: %d\n", sum)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	var sum int
	matches := make([]int, len(s.Tickets))
	for i, ticket := range s.Tickets {
		// just aliases to cut down on typing
		rev := ticket.Revealed
		win := ticket.Winning
		for r, w := 0, 0; r < len(rev) && w < len(win); {
			if rev[r] == win[w] {
				matches[i]++
				r++
				w++
				continue
			}
			if rev[r] > win[w] {
				w++
				continue
			}
			if win[w] > rev[r] {
				r++
				continue
			}
		}
	}

	copies := make([]int, len(matches))
	for i := range copies {
		copies[i] = 1
	}

	for i, count := range matches {
		for j := i + 1; j <= i+count; j++ {
			copies[j] += copies[i]
		}
	}

	for _, count := range copies {
		sum += count
	}

	fmt.Fprintf(w, "Part 2: %d\n", sum)
	return
}
