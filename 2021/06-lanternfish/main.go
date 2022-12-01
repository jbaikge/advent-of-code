package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Lanternfish struct {
	Timer int
}

func NewLanternfish() *Lanternfish {
	return &Lanternfish{8}
}

func (l *Lanternfish) GiveBirth() bool {
	l.Timer--
	if l.Timer == -1 {
		l.Timer = 6
		return true
	}
	return false
}

func (l *Lanternfish) String() string {
	return fmt.Sprint(l.Timer)
}

func population(init []int, days int) int {
	fish := make([]*Lanternfish, 0, 1000000)
	for _, age := range init {
		fish = append(fish, &Lanternfish{age})
	}
	for d := 0; d < days; d++ {
		fishLen := len(fish)
		for i := 0; i < fishLen; i++ {
			if fish[i].GiveBirth() {
				fish = append(fish, NewLanternfish())
			}
		}

		// ages := make([]string, len(fish))
		// for i := range ages {
		// 	ages[i] = fish[i].String()
		// }
		// log.Printf("Day %d: %s", d, strings.Join(ages, ", "))
	}
	return len(fish)
}

func populationV2(init []int, days int) (total uint64) {
	gestation := make([]uint64, 9)
	for _, v := range init {
		gestation[v]++
	}
	for day := 1; day <= days; day++ {
		pregnant := gestation[0]
		for i := 1; i < len(gestation); i++ {
			gestation[i-1] = gestation[i]
		}
		gestation[6] += pregnant
		gestation[8] = pregnant
	}
	for _, v := range gestation {
		total += v
	}
	return
}

func main() {
	var ages []int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), ",")
		ages = make([]int, 0, len(values))
		for _, value := range values {
			i, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Atoi failure with %s: %v", value, err)
			}
			ages = append(ages, i)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	log.Printf("%+v", ages)

	// Part 1
	log.Printf("Fish after 18 days: %d", population(ages, 18))
	log.Printf("Fish after 80 days: %d", population(ages, 80))
	log.Printf("Fish after 18 days v2: %d", populationV2(ages, 18))
	log.Printf("Fish after 80 days v2: %d", populationV2(ages, 80))
	// Part 2
	start := time.Now()
	log.Printf("Fish after 256 days: %d", populationV2(ages, 256))
	log.Printf("Took: %s", time.Since(start))
}
