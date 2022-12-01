package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

var segments = []string{
	"abcefg",  // 0
	"cf",      // 1*
	"acdeg",   // 2
	"acdfg",   // 3
	"bcdf",    // 4*
	"abdfg",   // 5
	"abdefg",  // 6
	"acf",     // 7*
	"abcdefg", // 8*
	"abcdfg",  // 9
}

func instances1478(output []string) (count int) {
	checks := []int{
		len(segments[1]),
		len(segments[4]),
		len(segments[7]),
		len(segments[8]),
	}
	for _, o := range output {
		for _, check := range checks {
			if len(o) == check {
				count++
				break
			}
		}
	}
	return
}

func sortChars(in []string) (out []string) {
	out = make([]string, len(in))
	for i := range in {
		chars := strings.Split(in[i], "")
		sort.Strings(chars)
		out[i] = strings.Join(chars, "")
	}
	return
}

func hasOverlay(s, overlay string) bool {
	sChars := strings.Split(s, "")
	overlayChars := strings.Split(overlay, "")
	for _, c := range overlayChars {
		if x := sort.SearchStrings(sChars, c); x >= len(sChars) || sChars[x] != c {
			return false
		}
	}
	return true
}

func decode(patterns []string, output []string) (num int) {
	sortedPatterns := sortChars(patterns)
	sortedOutput := sortChars(output)
	mapping := make([]string, 10)

	// Find the unique-length instances
	for _, i := range []int{1, 4, 7, 8} {
		for _, p := range sortedPatterns {
			if len(p) == len(segments[i]) {
				mapping[i] = p
				break
			}
		}
	}

	for _, p := range sortedPatterns {
		// 3 is the only one of len 5 with 1 or 7 overlaid
		if len(p) == 5 && hasOverlay(p, mapping[7]) {
			mapping[3] = p
			continue
		}
		// 6 is the only one of len 6 with 1 or 7 NOT overlaid
		if len(p) == 6 && !hasOverlay(p, mapping[7]) {
			mapping[6] = p
			continue
		}
		// 9 is the only one of len 6 with 4 overlaid
		if len(p) == 6 && hasOverlay(p, mapping[4]) {
			mapping[9] = p
			continue
		}
		// Now check for zero - the only len 6 with 1 or 7 overlaid
		if len(p) == 6 && hasOverlay(p, mapping[7]) {
			mapping[0] = p
			continue
		}
	}

	// Have to start another loop to ensure 6 is filled in
	for _, p := range sortedPatterns {
		// 5 overlays 6
		if len(p) == 5 && hasOverlay(mapping[6], p) {
			mapping[5] = p
			continue
		}
		// 2 does NOT overlay 6
		if len(p) == 5 && p != mapping[3] && !hasOverlay(mapping[6], p) {
			mapping[2] = p
			continue
		}
	}

	for i, o := range sortedOutput {
		for n, m := range mapping {
			if m == o {
				num += n * int(math.Pow10(3-i))
				continue
			}
		}
	}

	return
}

func main() {
	start := time.Now()
	totalInstances := 0
	totalSum := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		blocks := strings.Split(scanner.Text(), " | ")
		patterns := strings.Fields(blocks[0])
		output := strings.Fields(blocks[1])
		totalInstances += instances1478(output)
		totalSum += decode(patterns, output)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	// Part 1
	log.Printf("Instances of 1, 4, 7, and 8: %d", totalInstances)
	log.Printf("Total Sum: %d", totalSum)
	log.Printf("Took %s", time.Since(start))
}
