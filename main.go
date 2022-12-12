package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	hillclimb "github.com/jbaikge/advent-of-code/2022/12-hill-climb"
	"github.com/jbaikge/advent-of-code/util"
)

var solutions = map[int]util.Solution{
	202212: new(hillclimb.Solution),
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

	if err = solution.Parse(r); err != nil {
		log.Fatalf("Unable to parse input: %v", err)
	}

	if err = solution.Part1(os.Stdout); err != nil {
		log.Fatalf("Unable to solve part 1: %v", err)
	}

	if err = solution.Part2(os.Stdout); err != nil {
		log.Fatalf("Unable to solve part 2: %v", err)
	}
}
