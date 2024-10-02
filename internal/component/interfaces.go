package component

import (
	"fmt"
	"strings"
	term_utils "tui/internal/term-utils"
)

type Printer interface {
	fmt.Stringer
	Pos() (row, col int)
}

type Clearer interface {
	Pos() (row, col int)
	Size() (rows, cols int)
}

type Updater interface {
	Printer
	Clearer
}

func Clear(c Clearer) string {
	var b strings.Builder
	rows, cols := c.Size()
	row, col := c.Pos()
	for i := range rows {
		b.WriteString(term_utils.MoveCur(row+i, col))
		b.WriteString(strings.Repeat(" ", cols))
	}
	return b.String()
}

func Print(s Printer) string {
	var b strings.Builder
	row, col := s.Pos()
	for i, line := range strings.Split(s.String(), "\n") {
		b.WriteString(term_utils.MoveCur(row+i, col))
		b.WriteString(line)
	}
	return b.String()
}

func Update(u Updater) string {
	return Clear(u) + Print(u)
}
