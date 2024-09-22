package component

import (
	"fmt"
	"strings"
	term_utils "tui/internal/term-utils"
	"unicode/utf8"
)

type Tabs struct {
	headers  []string
	Selected int
	row      int
	col      int
}

func NewTabs(headers ...string) *Tabs {
	return &Tabs{
		row:     1,
		col:     1,
		headers: headers,
	}
}

func (t Tabs) String() string {
	b := new(strings.Builder)

	b.WriteString(term_utils.Lightgray)
	b.WriteString(term_utils.SquareTopLeft)
	for i, h := range t.headers {
		b.WriteString(repeat(term_utils.HorizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			b.WriteString(term_utils.SquareDownHorizontal)
		}
	}
	b.WriteString(term_utils.SquareTopRight + "\n")

	b.WriteString(term_utils.VerticalLine)
	for i, h := range t.headers {
		if t.Selected == i {
			b.WriteString(term_utils.Reset)
		}
		b.WriteString(fmt.Sprintf(" %s %s", h, term_utils.Reset))
		b.WriteString(term_utils.Lightgray)
		b.WriteString(term_utils.VerticalLine)
	}
	b.WriteString("\n")

	b.WriteString(term_utils.SquareBottomLeft)
	for i, h := range t.headers {
		b.WriteString(repeat(term_utils.HorizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			b.WriteString(term_utils.SquareUpHorizontal)
		}
	}
	b.WriteString(term_utils.SquareBottomRight)
	b.WriteString(term_utils.Reset)

	return b.String()
}

func (t *Tabs) Select(i int) *Tabs {
	if i < 0 || i >= len(t.headers) {
		return t
	}
	t.Selected = i
	return t
}

func (t *Tabs) Next() {

	t.moveToStartofWord()
	fmt.Print(fmt.Sprintf("%s%s%s", term_utils.Lightgray, t.headers[t.Selected], term_utils.Reset))

	t.Selected++
	if t.Selected >= len(t.headers) {
		t.Selected = 0
	}

	t.moveToStartofWord()
	fmt.Print(fmt.Sprintf("%s%s", term_utils.Reset, t.headers[t.Selected]))

}

func (t Tabs) moveToStartofWord() {
	col := t.col + 2
	for i := range t.Selected {
		col += StringLen(t.headers[i]) + 3
	}

	term_utils.MoveCursor(t.row+1, col)
}

func (b *Tabs) Lines() []string {
	return strings.Split(b.String(), "\n")
}

func (t *Tabs) GetSelected() string {
	return t.headers[t.Selected]
}

func (t *Tabs) At(row, col int) *Tabs {
	t.row = row
	t.col = col
	return t
}

func (t *Tabs) Pos() (row, col int) {
	return t.row, t.col
}

func repeat(s string, i int) string {
	return strings.Repeat(s, i)
}
