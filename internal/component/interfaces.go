package component

import (
	"fmt"
	"strings"
	term_utils "tui/internal/term-utils"
)

// FIX: ? Feels weird
type Instructor interface {
	Chan() chan<- string
}

type Printer interface {
	Instructor
	fmt.Stringer
	Pos() (row, col int)
}

type Clearer interface {
	Instructor
	Pos() (row, col int)
	Size() (rows, cols int)
}

type Updater interface {
	Printer
	Clearer
}

// FIX: ? Feels weird
func Clear(c Clearer) {
	var b strings.Builder
	rows, cols := c.Size()
	row, col := c.Pos()
	for i := range rows {
		b.WriteString(term_utils.MoveCur(row+i, col))
		b.WriteString(strings.Repeat(" ", cols))
	}
	c.Chan() <- b.String()
}

// FIX: ? Feels weird
func Print(s Printer) {
	var b strings.Builder
	row, col := s.Pos()
	for i, line := range strings.Split(s.String(), "\n") {
		b.WriteString(term_utils.MoveCur(row+i, col))
		b.WriteString(line)
	}
	s.Chan() <- b.String()
}

func Update(u Updater) {
	Clear(u)
	Print(u)
}
