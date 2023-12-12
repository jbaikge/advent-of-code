package miragemaintenance

import (
	"bufio"
	"bytes"
	_ "embed"
	"strconv"
	"strings"

	"github.com/jbaikge/advent-of-code/solutions"
)

//go:embed test.txt
var testData []byte

//go:embed input.txt
var inputData []byte

func init() {
	solutions.Register(new(Solution))
}

type Solution struct {
	Sets [][]int
}

func (*Solution) Meta() solutions.Meta {
	return solutions.Meta{
		Name:    "mirage maintenance",
		Year:    2023,
		Problem: 9,
		Datas: []solutions.Data{
			{
				Name:    "Test 1",
				Input:   []byte(`0 3 6 9 12 15`),
				Expect1: 18,
				Expect2: -3,
			},
			{
				Name:    "Test 2",
				Input:   []byte(`1 3 6 10 15 21`),
				Expect1: 28,
				Expect2: 0,
			},
			{
				Name:    "Test 3",
				Input:   []byte(`10 13 16 21 30 45`),
				Expect1: 68,
				Expect2: 5,
			},
			{
				Name:    "Test 4",
				Input:   testData,
				Expect1: 114,
				Expect2: 2,
			},
			{
				Name:  "Input",
				Input: inputData,
			},
		},
	}
}

func (s *Solution) Parse(data []byte) (err error) {
	s.Sets = make([][]int, 0, 200)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		values := strings.Fields(scanner.Text())
		set := make([]int, len(values))
		for i, value := range values {
			if set[i], err = strconv.Atoi(value); err != nil {
				return
			}
		}
		s.Sets = append(s.Sets, set)
	}

	return
}

func (s *Solution) Part1() (answer int, err error) {
	for _, set := range s.Sets {
		answer += FindNext(set)
	}
	return
}

func (s *Solution) Part2() (answer int, err error) {
	for _, set := range s.Sets {
		answer += FindPrev(set)
	}
	return
}

func Tree(nums []int) [][]int {
	tree := make([][]int, 0, len(nums))
	tree = append(tree, nums)

	for i := 0; i < len(nums); i++ {
		allZeros := true

		values := tree[i]
		diffs := make([]int, len(nums)-len(tree))
		for i := range values[1:] {
			diffs[i] = values[i+1] - values[i]
			if diffs[i] != 0 {
				allZeros = false
			}
		}
		tree = append(tree, diffs)

		if allZeros {
			break
		}
	}

	return tree
}

func FindNext(nums []int) int {
	tree := Tree(nums)

	for i := len(tree) - 1; i >= 0; i-- {
		if i+1 == len(tree) {
			tree[i] = append(tree[i], 0)
			continue
		}
		belowTail := tree[i+1][len(tree[i+1])-1]
		tree[i] = append(tree[i], tree[i][len(tree[i])-1]+belowTail)
	}

	return tree[0][len(tree[0])-1]
}

func FindPrev(nums []int) int {
	tree := Tree(nums)

	for i := len(tree) - 1; i >= 0; i-- {
		if i+1 == len(tree) {
			tree[i] = append([]int{0}, tree[i]...)
			continue
		}
		belowHead := tree[i+1][0]
		tree[i] = append([]int{tree[i][0] - belowHead}, tree[i]...)
	}

	return tree[0][0]
}
