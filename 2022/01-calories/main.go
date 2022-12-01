package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const Part2TopElves = 3

func main() {
	// Part 1 - Total calories held by top elf
	elfCalories := make([]int, 1, 100)
	elf := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			elf++
			elfCalories = append(elfCalories, 0)
			continue
		}
		calories, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Error converting %s to int: %v\n", line, err)
		}
		elfCalories[elf] += calories
	}

	sort.Sort(sort.Reverse(sort.IntSlice(elfCalories)))

	fmt.Printf("Elf with the most calories: %d\n", elfCalories[0])

	// Part 2 - Total calories held by top 3 elves
	topElvesCalories := 0
	for i := 0; i < Part2TopElves; i++ {
		topElvesCalories += elfCalories[i]
	}
	fmt.Printf("Top %d elves are holding %d calories\n", Part2TopElves, topElvesCalories)
}
