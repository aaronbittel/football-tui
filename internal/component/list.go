package component

import (
	"strings"
	term_utils "tui/internal/term-utils"
	"unicode/utf8"
)

type List struct {
	instrCh chan<- string

	row int
	col int

	builder  strings.Builder
	items    []string
	Selected int
	padding  Padding
	maxLen   int
}

func NewList(instrCh chan<- string, items ...string) *List {
	return &List{
		instrCh: instrCh,
		builder: strings.Builder{},

		row:     1,
		col:     1,
		items:   items,
		padding: NewPadding(0, 1),
	}
}

func (l *List) PrintIdle() {
	l.instrCh <- Print(l)
}

func (l *List) String() string {
	defer l.builder.Reset()

	for _, item := range l.items {
		length := utf8.RuneCountInString(item)
		if length > l.maxLen {
			l.maxLen = length
		}
	}

	for i, item := range l.items {
		if i == l.Selected {
			l.builder.WriteString(term_utils.BgRedFgWhite)
		}
		l.builder.WriteString(strings.Repeat(" ", l.padding.left))
		l.builder.WriteString(item)
		l.builder.WriteString(strings.Repeat(" ", l.maxLen-utf8.RuneCountInString(item)))
		l.builder.WriteString(strings.Repeat(" ", l.padding.right))
		if i == l.Selected {
			l.builder.WriteString(term_utils.ResetCode)
		}
		l.builder.WriteString("\n")
	}

	return l.builder.String()
}

func (l *List) At(row, col int) *List {
	l.row = row
	l.col = col
	return l
}

func (l *List) Pos() (row, col int) {
	return l.row, l.col
}

func (l *List) Next() {
	if l.Selected+1 >= len(l.items) {
		return
	}

	removeInstructions := l.removeHighlight()
	l.Selected++
	addInstructions := l.addHighlight()
	l.instrCh <- removeInstructions + addInstructions
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

func (l *List) Prev() {
	if l.Selected-1 < 0 {
		return
	}

	removeInstructions := l.removeHighlight()
	l.Selected--
	addInstructions := l.addHighlight()

	l.instrCh <- removeInstructions + addInstructions
}

func (l *List) Select(i int) *List {
	l.Selected = i
	return l
}

func (l List) SelectedValue() string {
	return l.items[l.Selected]
}
