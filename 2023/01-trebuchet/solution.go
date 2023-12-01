package trebuchet

import (
	"bytes"
	"embed"
	"fmt"
	"io"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Solution struct {
	Lines [][]byte
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}
	s.Lines = bytes.Split(bytes.TrimSpace(content), []byte{'\n'})
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	sum := 0
	for _, line := range s.Lines {
		var digits int
		// Find first digit by traversing forward
		for _, ch := range line {
			if ch >= '0' && ch <= '9' {
				digits = int(ch-'0') * 10
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			ch := line[i]
			if ch >= '0' && ch <= '9' {
				digits += int(ch - '0')
				break
			}
		}
		sum += digits
	}
	fmt.Fprintf(w, "Part 1: %d\n", sum)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	numbers := bytes.Fields([]byte(`0 1 2 3 4 5 6 7 8 9 zero one two three four five six seven eight nine`))
	sum := 0
	for _, line := range s.Lines {
		var digits int
		for i := range line {
			var found bool

			for n, num := range numbers {
				if bytes.HasPrefix(line[i:], num) {
					if n > 9 {
						n = n - 10
					}
					digits = n * 10
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		for i := range line {
			var found bool

			for n, num := range numbers {
				if bytes.HasSuffix(line[:len(line)-i], num) {
					if n > 9 {
						n -= 10
					}
					digits += n
					found = true
					break
				}
			}

			if found {
				break
			}
		}
		sum += digits
	}
	fmt.Fprintf(w, "Part 2: %d\n", sum)
	return
}
