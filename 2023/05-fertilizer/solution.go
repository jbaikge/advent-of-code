package fertilizer

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"math"
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

type Range struct {
	Bounds []Bound
}

func (r *Range) Append(b Bound) {
	r.Bounds = append(r.Bounds, b)
}

func (r Range) Destination(source int) (destination int) {
	for _, b := range r.Bounds {
		if source >= b.Source && source <= b.Source+b.Length {
			return b.Destination + (source - b.Source)
		}
	}
	return source
}

type Solution struct {
	Seeds                 []int
	SeedToSoil            Range
	SoilToFertilizer      Range
	FertilizerToWater     Range
	WaterToLight          Range
	LightToTemperature    Range
	TemperatureToHumidity Range
	HumidityToLocation    Range
}

func (s Solution) Location(seed int) (loc int) {
	soil := s.SeedToSoil.Destination(seed)
	fert := s.SoilToFertilizer.Destination(soil)
	water := s.FertilizerToWater.Destination(fert)
	light := s.WaterToLight.Destination(water)
	temp := s.LightToTemperature.Destination(light)
	humid := s.TemperatureToHumidity.Destination(temp)
	loc = s.HumidityToLocation.Destination(humid)
	return
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
			target = nil
			continue
		}

		fields := strings.Fields(line)
		switch fields[0] {
		case "seeds:":
			s.Seeds = make([]int, len(fields[1:]))
			for i, v := range fields[1:] {
				if s.Seeds[i], err = strconv.Atoi(v); err != nil {
					return
				}
			}
		case "seed-to-soil":
			target = &s.SeedToSoil
		case "soil-to-fertilizer":
			target = &s.SoilToFertilizer
		case "fertilizer-to-water":
			target = &s.FertilizerToWater
		case "water-to-light":
			target = &s.WaterToLight
		case "light-to-temperature":
			target = &s.LightToTemperature
		case "temperature-to-humidity":
			target = &s.TemperatureToHumidity
		case "humidity-to-location":
			target = &s.HumidityToLocation
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
	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	min := math.MaxInt
	for _, seed := range s.Seeds {
		if loc := s.Location(seed); loc < min {
			min = loc
		}
	}
	fmt.Fprintf(w, "Part 1: %d\n", min)
	return
}

func (s Solution) Part2(w io.Writer) (err error) {
	min := math.MaxInt
	for i := 0; i < len(s.Seeds); i += 2 {
		fmt.Printf("%d: %d -> %d\n", i, s.Seeds[i], s.Seeds[i]+s.Seeds[i+1])
		for j := s.Seeds[i]; j < s.Seeds[i]+s.Seeds[i+1]; j++ {
			if loc := s.Location(j); loc < min {
				fmt.Printf("New Min: %d\n", loc)
				min = loc
			}
		}
	}
	fmt.Fprintf(w, "Part 2: %d\n", min)
	return
}