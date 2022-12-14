package distress

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strings"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)
var _ sort.Interface = make(Parsed, 0)

type Packet struct {
	Raw string
}

func (p Packet) Parse() (out []interface{}) {
	out = make([]interface{}, 0)
	if err := json.NewDecoder(strings.NewReader(p.Raw)).Decode(&out); err != nil {
		log.Fatalf("Error parsing %s: %v", p.Raw, err)
	}
	return
}

type Parsed [][]interface{}

func (p Parsed) Len() int {
	return len(p)
}

func (p Parsed) Less(i, j int) bool {
	left := p[i]
	right := p[j]
	return p.compare(left, right, 0) == 1
}

func (p Parsed) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// 1  = in order
// 0  = not in order
// -1 = undetermined
func (p Parsed) compare(left, right []interface{}, depth int) (ordered int) {
	minLength := len(left)
	if rLen := len(right); rLen < minLength {
		minLength = rLen
	}

	for i := 0; i < minLength; i++ {
		var leftIsNum, rightIsNum bool
		var leftNum, rightNum int
		var leftSlice, rightSlice []interface{}

		switch v := left[i].(type) {
		case float64:
			leftIsNum = true
			leftNum = int(math.Round(v))
		case int:
			leftIsNum = true
			leftNum = v
		case []interface{}:
			leftSlice = v
		default:
			log.Fatalf("Not sure how to parse left value: left[%d]: %T %+v\n", i, v, v)
		}

		switch v := right[i].(type) {
		case float64:
			rightIsNum = true
			rightNum = int(math.Round(v))
		case int:
			rightIsNum = true
			rightNum = v
		case []interface{}:
			rightSlice = v
		default:
			log.Fatalf("Not sure how to parse right value: right[%d]: %T %+v\n", i, v, v)
		}

		// Actually compare numerical values
		if leftIsNum && rightIsNum {
			// Same, no need to do anything
			if leftNum == rightNum {
				continue
			}
			// Left is smaller, right order
			if leftNum < rightNum {
				return 1
			} else {
				return 0
			}
		}

		// Wrap scalar values in array and recurse
		if leftIsNum {
			leftSlice = []interface{}{leftNum}
		}
		if rightIsNum {
			rightSlice = []interface{}{rightNum}
		}
		switch p.compare(leftSlice, rightSlice, depth+1) {
		case 1:
			return 1
		case 0:
			return 0
		case -1:
			continue
		}
	}

	// If the left side runs out of items before the right side,
	// they are in order
	if len(right) > minLength {
		return 1
	}

	// If the right side runs out of items before the left side,
	// they are in the wrong order
	if len(left) > minLength {
		return 0
	}

	// Everything above is equal
	return -1
}

type Solution struct {
	Packets []Packet
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Packets = make([]Packet, 0, 300)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		s.Packets = append(s.Packets, Packet{Raw: line})
	}
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	parsed := make(Parsed, 0, len(s.Packets))
	for _, packet := range s.Packets {
		parsed = append(parsed, packet.Parse())
	}

	var correct int
	for i := 0; i < len(parsed); i += 2 {
		inOrder := parsed.Less(i, i+1)
		if inOrder {
			correct += i/2 + 1
		}
	}

	fmt.Fprintf(w, "Part 1: %d\n", correct)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	parsed := make(Parsed, 0, len(s.Packets))
	for _, packet := range s.Packets {
		parsed = append(parsed, packet.Parse())
	}

	// Add dividers
	dividers := [][]interface{}{
		{[]interface{}{2}},
		{[]interface{}{6}},
	}
	parsed = append(parsed, dividers...)

	sort.Sort(parsed)

	// Find indexes of dividers by comparing the string representation. Lots of
	// code required to break down and cast the interfaces for comparison
	decoderKey := 1
	for i, p := range parsed {
		str := fmt.Sprintf("%v", p)
		if str == "[[2]]" || str == "[[6]]" {
			decoderKey *= i + 1
		}
	}

	fmt.Fprintf(w, "Part 2: %d\n", decoderKey)
	return
}
