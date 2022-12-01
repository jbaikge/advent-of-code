package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Map struct {
	points  [][]int
	visited [][]bool
}

func NewMap() *Map {
	return &Map{
		points:  make([][]int, 0, 100),
		visited: make([][]bool, 0, 100),
	}
}

func (m *Map) AddRow(line string) (err error) {
	row := make([]int, len(line))
	values := strings.Split(line, "")
	for i, v := range values {
		if row[i], err = strconv.Atoi(v); err != nil {
			return
		}
	}
	m.points = append(m.points, row)
	m.visited = append(m.visited, make([]bool, len(row)))
	return
}

func (m *Map) IsLowest(x, y int) bool {
	v := m.points[y][x]
	isLowest := true

	// Left
	if X := x - 1; X >= 0 && m.points[y][X] <= v {
		isLowest = isLowest && m.points[y][X] > v
	}
	// Right
	if X := x + 1; X < len(m.points[y]) {
		isLowest = isLowest && m.points[y][X] > v
	}
	// Above
	if Y := y - 1; Y >= 0 {
		isLowest = isLowest && m.points[Y][x] > v
	}
	// Below
	if Y := y + 1; Y < len(m.points) {
		isLowest = isLowest && m.points[Y][x] > v
	}

	return isLowest
}

func (m *Map) FindLowPoints() (low []int) {
	low = make([]int, 0, 1024)
	for y, row := range m.points {
		for x, v := range row {
			if m.IsLowest(x, y) {
				low = append(low, v)
			}
		}
	}
	return
}

func (m *Map) Crawl(i, j int) (size int) {
	if m.points[i][j] == 9 {
		return
	}
	if m.visited[i][j] {
		return
	}

	// "this" point counts as 1
	m.visited[i][j] = true
	size++

	// Scan left
	if j-1 >= 0 {
		size += m.Crawl(i, j-1)
	}
	// Scan right
	if j+1 < len(m.points[i]) {
		size += m.Crawl(i, j+1)
	}
	// Scan up
	if i-1 >= 0 {
		size += m.Crawl(i-1, j)
	}
	// Scan down
	if i+1 < len(m.points) {
		size += m.Crawl(i+1, j)
	}
	return
}

func (m *Map) FindBasins() (basins []int) {
	for i := range m.points {
		for j := range m.points[i] {
			size := m.Crawl(i, j)
			if size == 0 {
				continue
			}
			basins = append(basins, size)
		}
	}
	return
}

func (m *Map) LowPointRisk() (risk int) {
	for _, p := range m.FindLowPoints() {
		risk += p + 1
	}
	return
}

func (m *Map) PrintPoints() {
	for _, row := range m.points {
		log.Printf("%v", row)
	}
}

func (m *Map) PrintVisited() {
	for _, row := range m.visited {
		log.Printf("%5v", row)
	}
}

func largestBasinProduct(basins []int, top int) (product int) {
	product = 1
	sort.Ints(basins)
	log.Printf("%v", basins[len(basins)-top:])
	for _, v := range basins[len(basins)-top:] {
		product *= v
	}
	return
}

func main() {
	m := NewMap()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := m.AddRow(scanner.Text()); err != nil {
			log.Fatalf("AddRow failure: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	// Part 1
	start := time.Now()
	log.Printf("%d %s", m.LowPointRisk(), time.Since(start))

	// Part 2
	start = time.Now()
	log.Printf("%d %s", largestBasinProduct(m.FindBasins(), 3), time.Since(start))
}
