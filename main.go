package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	chiton "github.com/jbaikge/advent-of-code/2021/15-chiton"
	hillclimb "github.com/jbaikge/advent-of-code/2022/12-hill-climb"
	distress "github.com/jbaikge/advent-of-code/2022/13-distress"
	reservoir "github.com/jbaikge/advent-of-code/2022/14-reservoir"
	sensors "github.com/jbaikge/advent-of-code/2022/15-sensors"
	valves "github.com/jbaikge/advent-of-code/2022/16-valves"
	tetris "github.com/jbaikge/advent-of-code/2022/17-tetris"
	trebuchet "github.com/jbaikge/advent-of-code/2023/01-trebuchet"
	cubeconundrum "github.com/jbaikge/advent-of-code/2023/02-cube-conundrum"
	gearratios "github.com/jbaikge/advent-of-code/2023/03-gear-ratios"
	scratchcards "github.com/jbaikge/advent-of-code/2023/04-scratchcards"
	"github.com/jbaikge/advent-of-code/util"
)

var solutions = map[int]util.Solution{
	202115: new(chiton.Solution),
	202212: new(hillclimb.Solution),
	202213: new(distress.Solution),
	202214: new(reservoir.Solution),
	202215: new(sensors.Solution),
	202216: new(valves.Solution),
	202217: new(tetris.Solution),
	202301: new(trebuchet.Solution),
	202302: new(cubeconundrum.Solution),
	202303: new(gearratios.Solution),
	202304: new(scratchcards.Solution),
}

func main() {
	flag.Parse()

	testKey := flag.Arg(0)
	if testKey == "" {
		fmt.Println("Valid tests are:")
		for k := range solutions {
			fmt.Printf(" - %d\n", k)
		}
		return
	}

	key, err := strconv.Atoi(testKey)
	if err != nil {
		log.Fatalf("Failed to parse first arg: %v", err)
	}

	solution, ok := solutions[key]
	if !ok {
		log.Fatalf("Solution not found for key: %d", key)
	}

	var r io.Reader = os.Stdin
	if basename := flag.Arg(1); basename != "" {
		filename := basename + ".txt"
		data, err := solution.Files().ReadFile(filename)
		if err != nil {
			log.Fatalf("Unable to get data for %s: %v", basename, err)
		}
		r = bytes.NewReader(data)
	}

	parseStart := time.Now()
	if err = solution.Parse(r); err != nil {
		log.Fatalf("Unable to parse input: %v", err)
	}
	parseTook := time.Since(parseStart)

	part1Start := time.Now()
	if err = solution.Part1(os.Stdout); err != nil {
		log.Fatalf("Unable to solve part 1: %v", err)
	}
	part1Took := time.Since(part1Start)

	part2Start := time.Now()
	if err = solution.Part2(os.Stdout); err != nil {
		log.Fatalf("Unable to solve part 2: %v", err)
	}
	part2Took := time.Since(part2Start)

	fmt.Printf("\nParse: %s Part 1: %s Part 2: %s\n", parseTook, part1Took, part2Took)
}
