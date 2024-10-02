package component

import (
	"strings"
	term_utils "tui/internal/term-utils"
	"unicode/utf8"
)

type List struct {
	builder  strings.Builder
	items    []string
	Selected int
	padding  Padding
	row      int
	col      int
	maxLen   int
}

func NewList(items ...string) *List {
	return &List{
		builder: strings.Builder{},
		row:     1,
		col:     1,
		items:   items,
		padding: NewPadding(0, 1),
	}
}

func moveToNewLine(row *int, col int) string {
	(*row)++
	return term_utils.MoveCur(*row, col)
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
			b.WriteString(term_utils.ResetCode)
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

func (l *List) Next() string {
	if l.Selected+1 >= len(l.items) {
		return ""
	}

	removeInstructions := l.removeHighlight()
	l.Selected++
	addInstructions := l.addHighlight()
	return removeInstructions + addInstructions
}

func (l List) removeHighlight() string {
	defer l.builder.Reset()

	l.builder.WriteString(term_utils.MoveCur(l.row+l.Selected, l.col))
	l.builder.WriteString(strings.Repeat(" ", l.maxLen+l.padding.right+l.padding.left))

	l.builder.WriteString(term_utils.MoveCur(l.row+l.Selected, l.col))
	item := l.items[l.Selected]
	l.builder.WriteString(strings.Repeat(" ", l.padding.left))
	l.builder.WriteString(item)
	l.builder.WriteString(strings.Repeat(" ", l.maxLen-utf8.RuneCountInString(item)))
	l.builder.WriteString(strings.Repeat(" ", l.padding.right))

	return l.builder.String()
}

func (l List) addHighlight() string {
	defer l.builder.Reset()

	l.builder.WriteString(term_utils.MoveCur(l.row+l.Selected, l.col))
	item := l.items[l.Selected]
	l.builder.WriteString(term_utils.BgRedFgWhite)
	l.builder.WriteString(strings.Repeat(" ", l.padding.left))
	l.builder.WriteString(item)
	l.builder.WriteString(strings.Repeat(" ", l.maxLen-utf8.RuneCountInString(item)))
	l.builder.WriteString(strings.Repeat(" ", l.padding.right))
	l.builder.WriteString(term_utils.ResetCode)

	return l.builder.String()
}

func (l *List) Prev() string {
	if l.Selected-1 < 0 {
		return ""
	}

	removeInstructions := l.removeHighlight()
	l.Selected--
	addInstructions := l.addHighlight()

	return removeInstructions + addInstructions
}

func (l *List) Select(i int) *List {
	l.Selected = i
	return l
}

func (l List) SelectedValue() string {
	return l.items[l.Selected]
}
