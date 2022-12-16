package chiton

import (
	"embed"
	"fmt"
	"io"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Solution struct{}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	fmt.Fprintf(w, "Part 1: %d\n", 0)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	fmt.Fprintf(w, "Part 2: %d\n", 0)
	return
}
