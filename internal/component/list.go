package component

import (
	"strings"
	utils "tui/internal/term-utils"
	"unicode/utf8"
)

const (
	bgRed        = "\033[48;2;255;95;95m"
	bgRedFgWhite = "\033[38;2;255;255;255;48;2;255;95;95m"
)

type List struct {
	items    []string
	Selected int
	padding  Padding
	row      int
	col      int
}

func NewList(items ...string) *List {
	return &List{
		items:   items,
		padding: NewPadding(0, 1),
	}
}

func (l *List) String() string {
	b := new(strings.Builder)

	maxLen := 0
	for _, item := range l.items {
		length := utf8.RuneCountInString(item)
		if length > maxLen {
			maxLen = length
		}
	}

	for i, item := range l.items {
		if i == l.Selected {
			b.WriteString(bgRedFgWhite)
		}
		b.WriteString(strings.Repeat(" ", l.padding.left))
		b.WriteString(item)
		b.WriteString(strings.Repeat(" ", maxLen-utf8.RuneCountInString(item)))
		b.WriteString(strings.Repeat(" ", l.padding.right))
		if i == l.Selected {
			b.WriteString(utils.Reset)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (l *List) At(row, col int) *List {
	l.row = row
	l.col = col
	return l
}

func (l *List) Pos() (row, col int) {
	return l.row, l.col
}

func (b *List) Lines() []string {
	return strings.Split(b.String(), "\n")
}

func (l *List) Next() {
	if l.Selected+1 < len(l.items) {
		l.Selected++
	}
}

func (l *List) Prev() {
	if l.Selected-1 < 0 {
		return
	}
	l.Selected--
	Clear(l)
	Print(l)
}

func (l *List) Select(i int) *List {
	l.Selected = i
	return l
}
