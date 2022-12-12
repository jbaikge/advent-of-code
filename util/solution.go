package util

import (
	"embed"
	"io"
)

type Solution interface {
	Files() embed.FS
	Parse(io.Reader) error
	Part1(io.Writer) error
	Part2(io.Writer) error
}
