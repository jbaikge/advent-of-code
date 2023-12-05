package gearratios

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"strconv"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Solution struct {
	Chars [][]byte
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Chars = make([][]byte, 0, 140)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		data := scanner.Bytes()
		line := make([]byte, len(data))
		copy(line, data)
		s.Chars = append(s.Chars, line)
	}
	return
}

func (s Solution) Nums() (nums [][][2]int) {
	nums = make([][][2]int, len(s.Chars))
	var inNumber bool
	for i, line := range s.Chars {
		nums[i] = make([][2]int, 0, 16)
		inNumber = false
		for j, ch := range line {
			if ch < '0' || ch > '9' {
				inNumber = false
				continue
			}
			if inNumber {
				nums[i][len(nums[i])-1][1] = j + 1
			} else {
				nums[i] = append(nums[i], [2]int{j, j + 1})
			}
			inNumber = true
		}

	}
	return
}

// 540824 - too high
func (s Solution) Part1(w io.Writer) (err error) {
	var sum int
	for i, numSlices := range s.Nums() {
		line := s.Chars[i]
		for _, slice := range numSlices {
			start, end := slice[0], slice[1]
			chars := s.Chars[i][start:end]
			length := end - start
			var add bool
			// north
			if i > 0 && !bytes.Equal(s.Chars[i-1][start:end], bytes.Repeat([]byte{'.'}, length)) {
				add = true
			}
			// northeast
			if i > 0 && end < len(line) && s.Chars[i-1][end] != '.' {
				add = true
			}
			// east
			if end < len(line) && s.Chars[i][end] != '.' {
				add = true
			}
			// southeast
			if i+1 < len(s.Chars) && end < len(line) && s.Chars[i+1][end] != '.' {
				add = true
			}
			// south
			if i+1 < len(s.Chars) && !bytes.Equal(s.Chars[i+1][start:end], bytes.Repeat([]byte{'.'}, length)) {
				add = true
			}
			// southwest
			if i+1 < len(s.Chars) && start > 0 && s.Chars[i+1][start-1] != '.' {
				add = true
			}
			// west
			if start > 0 && s.Chars[i][start-1] != '.' {
				add = true
			}
			// northwest
			if i > 0 && start > 0 && s.Chars[i-1][start-1] != '.' {
				add = true
			}
			// Add to sum
			if add {
				num, _ := strconv.Atoi(string(chars))
				sum += num
			}
		}
	}

	fmt.Fprintf(w, "Part 1: %d\n", sum)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	var sum int

	gears := make(map[[2]int][]int)
	var key [2]int

	for i, numSlices := range s.Nums() {
		line := s.Chars[i]
		for _, slice := range numSlices {
			start, end := slice[0], slice[1]
			chars := s.Chars[i][start:end]
			key = [2]int{0, 0}

			// north
			if i > 0 && bytes.IndexByte(s.Chars[i-1][start:end], '*') > -1 {
				key[0] = i - 1
				key[1] = start + bytes.IndexByte(s.Chars[i-1][start:end], '*')
			}
			// northeast
			if i > 0 && end < len(line) && s.Chars[i-1][end] == '*' {
				key[0] = i - 1
				key[1] = end
			}
			// east
			if end < len(line) && s.Chars[i][end] == '*' {
				key[0] = i
				key[1] = end
			}
			// southeast
			if i+1 < len(s.Chars) && end < len(line) && s.Chars[i+1][end] == '*' {
				key[0] = i + 1
				key[1] = end
			}
			// south
			if i+1 < len(s.Chars) && bytes.IndexByte(s.Chars[i+1][start:end], '*') > -1 {
				key[0] = i + 1
				key[1] = start + bytes.IndexByte(s.Chars[i+1][start:end], '*')
			}
			// southwest
			if i+1 < len(s.Chars) && start > 0 && s.Chars[i+1][start-1] == '*' {
				key[0] = i + 1
				key[1] = start - 1
			}
			// west
			if start > 0 && s.Chars[i][start-1] == '*' {
				key[0] = i
				key[1] = start - 1
			}
			// northwest
			if i > 0 && start > 0 && s.Chars[i-1][start-1] == '*' {
				key[0] = i - 1
				key[1] = start - 1
			}
			// Add to gear map
			if key[0] > 0 && key[1] > 0 {
				num, _ := strconv.Atoi(string(chars))

				if _, found := gears[key]; !found {
					gears[key] = make([]int, 0, 4)
				}
				gears[key] = append(gears[key], num)
			}
		}
	}

	for _, nums := range gears {
		if len(nums) != 2 {
			continue
		}
		ratio := nums[0] * nums[1]
		sum += ratio
	}

	fmt.Fprintf(w, "Part 2: %d\n", sum)
	return
}
