package fertilizer

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

type Bound struct {
	Destination int
	Source      int
	Length      int
}

func (b Bound) Lower() int {
	return b.Source
}

func (b Bound) Upper() int {
	return b.Source + b.Length - 1
}

func (b Bound) Delta() int {
	return b.Destination - b.Source
}

func (b Bound) Overlaps(a Bound) bool {
	// b fully contains a
	// b ... a ... a ... b
	if b.Lower() <= a.Lower() && b.Upper() >= a.Upper() {
		return true
	}
	// b starts before a and there is overlap
	// b ... a ... b ... a
	if b.Lower() <= a.Lower() && b.Upper() >= a.Lower() && b.Upper() <= a.Upper() {
		return true
	}
	// b starts after a and there is overlap
	// a ... b ... a ... b
	if b.Lower() >= a.Lower() && b.Lower() <= a.Upper() && b.Upper() >= a.Upper() {
		return true
	}
	return false
}

func (b Bound) Split(a Bound) (chunks []Bound) {
	chunks = make([]Bound, 1, 2)
	chunks[0] = a
	if b.Upper() < a.Upper() {
		chunks[0].Length = b.Upper() - chunks[0].Lower() + 1
		chunks = append(chunks, Bound{
			Source: b.Upper() + 1,
			Length: a.Length - chunks[0].Length,
		})
	}
	return
}

func (b Bound) String() string {
	return fmt.Sprintf("<%d, %d>", b.Lower(), b.Upper())
}

type Range struct {
	Name   string
	Bounds []Bound
}

func (r *Range) AddCaps() {
	r.Sort()

	if r.Bounds[0].Source > 0 {
		r.Bounds = append([]Bound{
			{
				Destination: 0,
				Source:      0,
				Length:      r.Bounds[0].Source,
			},
		}, r.Bounds...)
	}

	lastBound := r.Bounds[len(r.Bounds)-1]
	if lastBound.Upper() < math.MaxUint32 {
		r.Bounds = append(r.Bounds, Bound{
			Destination: lastBound.Upper(),
			Source:      lastBound.Upper(),
			Length:      math.MaxUint32 - lastBound.Upper(),
		})
	}

}

func (r *Range) Append(b Bound) {
	r.Bounds = append(r.Bounds, b)
}

func (r *Range) FillGaps() {
	for i, current := range r.Bounds[1:] {
		prev := r.Bounds[i]
		if prev.Upper()+1 == current.Lower() {
			continue
		}
		filler := Bound{
			Source:      prev.Upper() + 1,
			Length:      current.Lower() - prev.Upper() - 1,
			Destination: prev.Upper() + 1,
		}
		// Expand the slice by 1
		r.Bounds = append(r.Bounds[:i+1], r.Bounds[i:]...)
		// i+1 because i is ahead by 1 with the slice operation
		r.Bounds[i+1] = filler
	}
}

func (r *Range) Sort() {
	slices.SortFunc[[]Bound, Bound](r.Bounds, func(a, b Bound) int {
		return a.Source - b.Source
	})
}

func (r Range) Destinations(b Bound) (destinations []Bound) {
	destinations = make([]Bound, 0, 16)
	for _, bound := range r.Bounds {
		if !bound.Overlaps(b) {
			continue
		}
		chunks := bound.Split(b)
		destinations = append(destinations, Bound{
			Source: chunks[0].Source + bound.Delta(),
			Length: chunks[0].Length,
		})
		if len(chunks) == 2 {
			b = chunks[1]
		}
	}
	return
}

type Solution struct {
	Seeds  []int
	Ranges []Range
}

func (s Solution) Location(seed Bound) (locs []Bound) {
	destinations := []Bound{seed}

	for _, r := range s.Ranges {
		sources := make([]Bound, len(destinations))
		copy(sources, destinations)
		destinations = destinations[:0]
		for _, source := range sources {
			destinations = append(destinations, r.Destinations(source)...)
		}
	}
	return destinations
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	scanner := bufio.NewScanner(r)
	var target *Range
	for scanner.Scan() {
		line := scanner.Text()

		// Skip blank lines
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		switch {
		case fields[0] == "seeds:":
			s.Seeds = make([]int, len(fields[1:]))
			for i, v := range fields[1:] {
				if s.Seeds[i], err = strconv.Atoi(v); err != nil {
					return
				}
			}
		case fields[1] == "map:":
			s.Ranges = append(s.Ranges, Range{
				Name:   fields[0],
				Bounds: make([]Bound, 0, 16),
			})
			target = &s.Ranges[len(s.Ranges)-1]
		default:
			if target == nil {
				return fmt.Errorf("target is nil")
			}
			var b Bound
			if b.Destination, err = strconv.Atoi(fields[0]); err != nil {
				return
			}
			if b.Source, err = strconv.Atoi(fields[1]); err != nil {
				return
			}
			if b.Length, err = strconv.Atoi(fields[2]); err != nil {
				return
			}
			target.Append(b)
		}
	}

	// Add caps
	// Also fill in gaps!
	for i := range s.Ranges {
		s.Ranges[i].AddCaps()
		s.Ranges[i].FillGaps()
	}

	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	min := math.MaxUint32
	for _, source := range s.Seeds {
		seed := Bound{Source: source, Length: 1}
		for _, loc := range s.Location(seed) {
			if loc.Lower() < min {
				min = loc.Lower()
			}
		}
	}
	fmt.Fprintf(w, "Part 1: %d\n", min)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	min := math.MaxUint32
	for i := 0; i < len(s.Seeds); i += 2 {
		seed := Bound{Source: s.Seeds[i], Length: s.Seeds[i+1]}
		for _, loc := range s.Location(seed) {
			if loc.Lower() < min {
				min = loc.Lower()
			}
		}
	}
	fmt.Fprintf(w, "Part 2: %d\n", min)
	return
}
