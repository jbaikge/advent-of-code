package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Very helpful writeup, even if I still don't understand DFS
// https://skarlso.github.io/2021/12/17/aoc-day12-updated/

const Start = "start"
const End = "end"

type Cave struct {
	Value       string
	PathSoFar   []string
	VistedTwice bool
}

func (c Cave) Seen(cave string) bool {
	for _, v := range c.PathSoFar {
		if cave == v {
			return true
		}
	}
	return false
}

func part1(caves map[string][]string) (count int) {
	start := Cave{
		Value:     Start,
		PathSoFar: []string{Start},
	}
	queue := []Cave{start}
	var current Cave
	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]
		if current.Value == End {
			count++
			continue
		}
		for _, next := range caves[current.Value] {
			if current.Seen(next) {
				continue
			}
			path := make([]string, 0, len(current.PathSoFar))
			path = append(path, current.PathSoFar...)
			// Can only visit small (lowercase) caves once
			if strings.ToLower(next) == next {
				path = append(path, next)
			}
			queue = append(queue, Cave{
				Value:     next,
				PathSoFar: path,
			})
		}
	}

	return
}

func part2(caves map[string][]string) (count int) {
	startEnd := Cave{
		PathSoFar: []string{Start, End},
	}
	start := Cave{
		Value:     Start,
		PathSoFar: []string{Start},
	}
	queue := []Cave{start}
	var current Cave
	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]
		if current.Value == End {
			count++
			continue
		}
		for _, next := range caves[current.Value] {
			seen := current.Seen(next)
			if !seen {
				path := make([]string, 0, len(current.PathSoFar))
				path = append(path, current.PathSoFar...)
				if strings.ToLower(next) == next {
					path = append(path, next)
				}
				queue = append(queue, Cave{
					Value:       next,
					PathSoFar:   path,
					VistedTwice: current.VistedTwice,
				})
				continue
			}
			if seen && !current.VistedTwice && !startEnd.Seen(next) {
				queue = append(queue, Cave{
					Value:       next,
					PathSoFar:   current.PathSoFar,
					VistedTwice: true,
				})
			}
		}
	}
	return
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	paths := make([][2]string, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "-", 2)
		paths = append(paths, [2]string{parts[0], parts[1]})
	}

	caves := make(map[string][]string)
	for _, path := range paths {
		caves[path[0]] = append(caves[path[0]], path[1])
		caves[path[1]] = append(caves[path[1]], path[0])
	}

	fmt.Printf("Part 1: %d\n", part1(caves))
	fmt.Printf("Part 2: %d\n", part2(caves))
}
