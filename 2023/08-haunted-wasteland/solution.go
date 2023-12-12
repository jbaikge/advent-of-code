package hauntedwasteland

import (
	"bufio"
	"bytes"
	_ "embed"

	"github.com/jbaikge/advent-of-code/solutions"
)

//go:embed test1.txt
var test1Data []byte

//go:embed test2.txt
var test2Data []byte

//go:embed test3.txt
var test3Data []byte

//go:embed input.txt
var inputData []byte

func init() {
	solutions.Register(new(Solution))
}

type Node struct {
	Left  string
	Right string
}

type Solution struct {
	Instructions string
	Nodes        map[string]Node
}

func (*Solution) Meta() solutions.Meta {
	return solutions.Meta{
		Name:    "haunted wasteland",
		Year:    2023,
		Problem: 8,
		Datas: []solutions.Data{
			{
				Name:    "Test 1",
				Input:   test1Data,
				Expect1: 2,
				Expect2: 0,
			},
			{
				Name:    "Test 2",
				Input:   test2Data,
				Expect1: 6,
				Expect2: 0,
			},
			{
				Name:    "Test 3",
				Input:   test3Data,
				Expect1: 0,
				Expect2: 6,
			},
			{
				Name:  "Input",
				Input: inputData,
			},
		},
	}
}

func (s *Solution) Parse(data []byte) (err error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// First line contains instructions
	scanner.Scan()
	s.Instructions = scanner.Text()

	// Next line is blank
	scanner.Scan()

	// Remaining lines are nodes
	s.Nodes = make(map[string]Node)
	for scanner.Scan() {
		line := scanner.Text()
		s.Nodes[line[:3]] = Node{
			Left:  line[7:10],
			Right: line[12:15],
		}
	}

	return
}

func (s *Solution) Part1() (answer int, err error) {
	key := "AAA"

	// Skip test 3 data
	if _, found := s.Nodes[key]; !found {
		return
	}

	idx := 0
	for key != "ZZZ" {
		instruction := s.Instructions[idx]
		if instruction == 'L' {
			key = s.Nodes[key].Left
		} else {
			key = s.Nodes[key].Right
		}
		answer++
		idx++
		if idx >= len(s.Instructions) {
			idx = 0
		}
	}
	return
}

// Greatest Common Divisor
// @see https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/
func GCD(a int, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a int, b int, nums ...int) int {
	result := a * b / GCD(a, b)

	for _, num := range nums {
		result = LCM(result, num)
	}

	return result
}

func (s *Solution) Part2() (answer int, err error) {
	counts := make(map[string]int)
	keys := make([]string, 0, 6)
	for key := range s.Nodes {
		if key[2] == 'A' {
			keys = append(keys, key)
			counts[key] = 0
		}
	}

	for _, key := range keys {
		start := key
		idx := 0
		for key[2] != 'Z' {
			instruction := s.Instructions[idx]
			if instruction == 'L' {
				key = s.Nodes[key].Left
			} else {
				key = s.Nodes[key].Right
			}
			counts[start]++
			if idx++; idx >= len(s.Instructions) {
				idx = 0
			}
		}
	}

	nums := make([]int, 0, len(counts))
	for _, count := range counts {
		nums = append(nums, count)
	}

	if len(nums) > 2 {
		answer = LCM(nums[0], nums[1], nums[2:]...)
	} else if len(nums) == 2 {
		answer = LCM(nums[0], nums[1])
	}
	return
}
