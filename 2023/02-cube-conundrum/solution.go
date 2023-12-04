package cubeconundrum

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jbaikge/advent-of-code/util"
)

//go:embed *.txt
var Files embed.FS

var _ util.Solution = new(Solution)

type Set struct {
	Red   int
	Green int
	Blue  int
}

type Game struct {
	Id   int
	Sets []Set
}

func NewGame(line string) (g Game) {
	var err error
	colonIdx := strings.IndexByte(line, ':')
	if colonIdx < 0 {
		panic("could not find colon! " + line)
	}

	g.Id, err = strconv.Atoi(line[5:colonIdx])
	if err != nil {
		panic("could not atoi: " + line[5:colonIdx])
	}

	sets := strings.Split(line[colonIdx+2:], "; ")
	g.Sets = make([]Set, len(sets))
	for i, set := range sets {
		cubes := strings.Split(set, ", ")
		for _, cube := range cubes {
			spaceIdx := strings.IndexByte(cube, ' ')
			num, err := strconv.Atoi(cube[:spaceIdx])
			if err != nil {
				panic("could not parse cube count: " + err.Error())
			}
			switch cube[spaceIdx+1:] {
			case "red":
				g.Sets[i].Red = num
			case "green":
				g.Sets[i].Green = num
			case "blue":
				g.Sets[i].Blue = num
			}
		}
	}
	return
}

type Solution struct {
	Games []Game
}

func (s Solution) Files() embed.FS {
	return Files
}

func (s *Solution) Parse(r io.Reader) (err error) {
	s.Games = make([]Game, 0, 100)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s.Games = append(s.Games, NewGame(scanner.Text()))
	}
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	bag := Set{
		Red:   12,
		Green: 13,
		Blue:  14,
	}
	var sum int
	for _, game := range s.Games {
		var impossible bool
		for _, set := range game.Sets {
			if set.Red > bag.Red || set.Green > bag.Green || set.Blue > bag.Blue {
				impossible = true
			}
		}
		if !impossible {
			sum += game.Id
		}
	}
	fmt.Fprintf(w, "Part 1: %d\n", sum)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	var sum int
	for _, game := range s.Games {
		min := Set{}
		for _, set := range game.Sets {
			if set.Red > min.Red {
				min.Red = set.Red
			}
			if set.Green > min.Green {
				min.Green = set.Green
			}
			if set.Blue > min.Blue {
				min.Blue = set.Blue
			}
		}
		power := min.Red * min.Green * min.Blue
		sum += power
	}
	fmt.Fprintf(w, "Part 2: %d\n", sum)
	return
}
