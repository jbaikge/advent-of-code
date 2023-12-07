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

type Range struct {
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
	lastSource := lastBound.Source + lastBound.Length
	if lastSource < math.MaxUint32 {
		r.Bounds = append(r.Bounds, Bound{
			Destination: lastSource,
			Source:      lastSource,
			Length:      math.MaxUint32 - lastSource,
		})
	}

}

func (r *Range) Append(b Bound) {
	r.Bounds = append(r.Bounds, b)
}

func (r *Range) Sort() {
	slices.SortFunc[[]Bound, Bound](r.Bounds, func(a, b Bound) int {
		return a.Source - b.Source
	})
}

func (r Range) Destination(source int, length int) (destinations [][2]int) {
	destinations = make([][2]int, 0, 16)
	for _, b := range r.Bounds {
		if source < b.Source || source >= b.Source+b.Length {
			continue
		}

		dest := [2]int{b.Destination + (source - b.Source), length}

		if source+length-1 > b.Source+b.Length {
			dest[1] = b.Length - (source - b.Source)
			source = b.Source + b.Length
			length = length - dest[1]
		}

		destinations = append(destinations, dest)
	}
	return
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

func (s Solution) Location(seed int, length int) (locs [][2]int) {
	soils := s.SeedToSoil.Destination(seed, length)
	// fmt.Printf("Soils: %v\n", soils)

	ferts := make([][2]int, 0, 32)
	for _, soil := range soils {
		ferts = append(ferts, s.SoilToFertilizer.Destination(soil[0], soil[1])...)
	}
	// fmt.Printf("Fers: %v\n", ferts)

	waters := make([][2]int, 0, 32)
	for _, fert := range ferts {
		waters = append(waters, s.FertilizerToWater.Destination(fert[0], fert[1])...)
	}
	// fmt.Printf("Waters: %v\n", waters)

	lights := make([][2]int, 0, 32)
	for _, water := range waters {
		lights = append(lights, s.WaterToLight.Destination(water[0], water[1])...)
	}
	// fmt.Printf("Lights: %v\n", lights)

	temps := make([][2]int, 0, 32)
	for _, light := range lights {
		temps = append(temps, s.LightToTemperature.Destination(light[0], light[1])...)
	}
	// fmt.Printf("Temps: %v\n", temps)

	humids := make([][2]int, 0, 32)
	for _, temp := range temps {
		humids = append(humids, s.TemperatureToHumidity.Destination(temp[0], temp[1])...)
	}
	// fmt.Printf("Humids: %v\n", humids)

	locs = make([][2]int, 0, 32)
	for _, humid := range humids {
		locs = append(locs, s.HumidityToLocation.Destination(humid[0], humid[1])...)
	}

	fmt.Println(locs)
	fmt.Println()
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

		// Add bounds to beginning and end of range
		if line == "" && target != nil {
			target.AddCaps()
		}

		// Skip blank lines
		if line == "" {
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

	// Need to add caps to final range
	target.AddCaps()

	for _, b := range s.SeedToSoil.Bounds {
		fmt.Printf("%d -> %d\n", b.Source, b.Source+b.Length-1)
	}

	return
}

func (s Solution) Part1(w io.Writer) (err error) {
	min := math.MaxInt
	for _, seed := range s.Seeds {
		for _, loc := range s.Location(seed, 1) {
			if loc[0] < min {
				min = loc[0]
			}
		}
	}
	fmt.Fprintf(w, "Part 1: %d\n", min)
	return
}

// Too high: 225749547
// Goal:     17729182
func (s Solution) Part2(w io.Writer) (err error) {
	min := math.MaxUint32
	for i := 0; i < len(s.Seeds); i += 2 {
		for _, loc := range s.Location(s.Seeds[i], s.Seeds[i+1]) {
			if loc[0] < min {
				min = loc[0]
			}
		}
	}
	fmt.Fprintf(w, "Part 2: %d\n", min)
	return
}
