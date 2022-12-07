package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type File struct {
	Name string
	Size int
}

type Dir struct {
	Files []File
}

func (d Dir) Size() (size int) {
	for _, f := range d.Files {
		size += f.Size
	}
	return
}

type Filesystem map[string]*Dir

func NewFilesystem() Filesystem {
	f := make(Filesystem)
	f["/"] = &Dir{
		Files: make([]File, 0, 16),
	}
	return f
}

func (f Filesystem) DirSizes() (sizes map[string]int) {
	sizes = make(map[string]int)
	for parent := range f {
		for path := range f {
			if strings.HasPrefix(path, parent) {
				sizes[parent] += f[path].Size()
			}
		}
	}
	return
}

// Find all of the directories with a total size of at most 100000, then
// calculate the sum of their total sizes
func part1(f Filesystem) (total int) {
	const Max = 100000

	for _, size := range f.DirSizes() {
		if size <= Max {
			total += size
		}
	}
	return
}

func part2(f Filesystem) (dirSize int) {
	const DiskSize = 70000000
	const DiskRequired = 30000000

	// Move sizes into a slice for sorting purposes
	dirSizes := f.DirSizes()
	sizes := make([]int, 0, len(dirSizes))
	for _, size := range dirSizes {
		sizes = append(sizes, size)
	}
	sort.Ints(sizes)

	spaceLeft := DiskSize - dirSizes["/"]
	needed := DiskRequired - spaceLeft
	for _, size := range sizes {
		if size >= needed {
			return size
		}
	}

	return
}

func main() {
	filesystem := NewFilesystem()
	currentPath := "/"

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch fields[0] {
		case "$":
			switch fields[1] {
			case "cd":
				switch fields[2] {
				case "/":
					// Absolute Path
					currentPath = fields[2]
				case "..":
					// Up a directory
					currentPath = filepath.Dir(currentPath)
				default:
					// Into a directory
					currentPath = filepath.Join(currentPath, fields[2])
				}
			case "ls":
				// NOOP
			}
		case "dir":
			filesystem[filepath.Join(currentPath, fields[1])] = &Dir{
				Files: make([]File, 0, 16),
			}
		default:
			size, _ := strconv.Atoi(fields[0])
			name := fields[1]
			dir := filesystem[currentPath]
			dir.Files = append(dir.Files, File{
				Name: name,
				Size: size,
			})
		}
	}

	fmt.Printf("Part 1: %d\n", part1(filesystem))
	fmt.Printf("Part 2: %d\n", part2(filesystem))
}
