package component

import (
	"fmt"
	"strings"
	term_utils "tui/internal/term-utils"
	"unicode/utf8"
)

type List struct {
	items    []string
	Selected int
	padding  Padding
	row      int
	col      int
	maxLen   int
}

func NewList(items ...string) *List {
	return &List{
		row:     1,
		col:     1,
		items:   items,
		padding: NewPadding(0, 1),
	}
}

func (l *List) String() string {
	b := new(strings.Builder)

	for _, item := range l.items {
		length := utf8.RuneCountInString(item)
		if length > l.maxLen {
			l.maxLen = length
		}
	}

	for i, item := range l.items {
		if i == l.Selected {
			b.WriteString(term_utils.BgRedFgWhite)
		}
		b.WriteString(strings.Repeat(" ", l.padding.left))
		b.WriteString(item)
		b.WriteString(strings.Repeat(" ", l.maxLen-utf8.RuneCountInString(item)))
		b.WriteString(strings.Repeat(" ", l.padding.right))
		if i == l.Selected {
			b.WriteString(term_utils.Reset)
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
	if l.Selected+1 >= len(l.items) {
		return
	}

	l.removeHighlight()
	l.Selected++
	l.addHighlight()

}

func (l List) removeHighlight() {
	term_utils.MoveCursor(l.row+l.Selected, l.col)
	fmt.Print(strings.Repeat(" ", l.maxLen+l.padding.right+l.padding.left))

	term_utils.MoveCursor(l.row+l.Selected, l.col)
	item := l.items[l.Selected]
	fmt.Print(strings.Repeat(" ", l.padding.left))
	fmt.Print(item)
	fmt.Print(strings.Repeat(" ", l.maxLen-utf8.RuneCountInString(item)))
	fmt.Print(strings.Repeat(" ", l.padding.right))
}

func (l List) addHighlight() {
	term_utils.MoveCursor(l.row+l.Selected, l.col)
	item := l.items[l.Selected]
	fmt.Print(term_utils.BgRedFgWhite)
	fmt.Print(strings.Repeat(" ", l.padding.left))
	fmt.Print(item)
	fmt.Print(strings.Repeat(" ", l.maxLen-utf8.RuneCountInString(item)))
	fmt.Print(strings.Repeat(" ", l.padding.right))
	fmt.Print(term_utils.Reset)
}

func (l *List) Prev() {
	if l.Selected-1 < 0 {
		return
	}

	l.removeHighlight()
	l.Selected--
	l.addHighlight()
}

func (l *List) Select(i int) *List {
	l.Selected = i
	return l
}

func (l List) SelectedValue() string {
	return l.items[l.Selected]
}
