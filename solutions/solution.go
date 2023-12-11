package solutions

import "fmt"

var registered []Solution

type Data struct {
	Name    string
	Input   []byte
	Expect1 int
	Expect2 int
}

type Meta struct {
	Name    string
	Year    int
	Problem int
	Datas   []Data
}

type Solution interface {
	Meta() Meta
	Parse([]byte) error
	Part1() (int, error)
	Part2() (int, error)
}

func Get(year int, problem int) (Solution, error) {
	for _, s := range registered {
		meta := s.Meta()
		if meta.Year == year && meta.Problem == problem {
			return s, nil
		}
	}
	return nil, fmt.Errorf("unable to find solution for year:%d problem: %d", year, problem)
}

func Register(s Solution) {
	registered = append(registered, s)
}
