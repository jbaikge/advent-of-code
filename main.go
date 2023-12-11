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
	fertilizer "github.com/jbaikge/advent-of-code/2023/05-fertilizer"
	waitforit "github.com/jbaikge/advent-of-code/2023/06-wait-for-it"
	"github.com/jbaikge/advent-of-code/solutions"
	"github.com/jbaikge/advent-of-code/util"

	_ "github.com/jbaikge/advent-of-code/2023/07-camel-cards"
)

var utilSolutions = map[int]util.Solution{
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
	202305: new(fertilizer.Solution),
	202306: new(waitforit.Solution),
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("Must provide a year and a problem number")
		os.Exit(1)
	}

	year, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Printf("Invalid year: %s\n", flag.Arg(0))
		os.Exit(1)
	}

	problem, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		fmt.Printf("Invalid problem: %s\n", flag.Arg(1))
		os.Exit(1)
	}

	solution, err := solutions.Get(year, problem)
	if err != nil {
		fmt.Printf("Could not find problem!\n")
		os.Exit(1)
	}

	meta := solution.Meta()
	fmt.Printf("%s\n", meta.Name)

	for _, data := range meta.Datas {
		fmt.Printf("\nProcessing data: %s\n", data.Name)

		parseStart := time.Now()
		if err := solution.Parse(data.Input); err != nil {
			log.Fatalf("Unable to parse data: %v", err)
		}
		parseTook := time.Since(parseStart)

		part1Start := time.Now()
		if answer, err := solution.Part1(); err != nil {
			log.Fatalf("Part 1 failed: %v", err)
		} else {
			fmt.Printf("  Part 1: %d", answer)
			if data.Expect1 != 0 {
				if answer == data.Expect1 {
					fmt.Print(" (Pass)")
				} else {
					fmt.Print(" (Fail)")
				}
			}
			fmt.Println()
		}
		part1Took := time.Since(part1Start)

		part2Start := time.Now()
		if answer, err := solution.Part2(); err != nil {
			log.Fatalf("Part 2 failed: %v", err)
		} else {
			fmt.Printf("  Part 2: %d", answer)
			if data.Expect2 != 0 {
				if answer == data.Expect2 {
					fmt.Println(" (Pass)")
				} else {
					fmt.Println(" (Fail)")
				}
			}
			fmt.Println()
		}
		part2Took := time.Since(part2Start)

		fmt.Printf("Parse: %s Part 1: %s Part 2: %s\n", parseTook, part1Took, part2Took)
	}
}

func utilMain() {
	flag.Parse()

	testKey := flag.Arg(0)
	if testKey == "" {
		fmt.Println("Valid tests are:")
		for k := range utilSolutions {
			fmt.Printf(" - %d\n", k)
		}
		return
	}

	key, err := strconv.Atoi(testKey)
	if err != nil {
		log.Fatalf("Failed to parse first arg: %v", err)
	}

	solution, ok := utilSolutions[key]
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
