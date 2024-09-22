package component

import (
	"fmt"
	"strings"
	utils "tui/internal/term-utils"
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

	b.WriteString(utils.Lightgray)
	b.WriteString(utils.SquareTopLeft)
	for i, h := range t.headers {
		b.WriteString(repeat(utils.HorizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			b.WriteString(utils.SquareDownHorizontal)
		}
	}
	b.WriteString(utils.SquareTopRight + "\n")

	b.WriteString(utils.VerticalLine)
	for i, h := range t.headers {
		if t.Selected == i {
			b.WriteString(utils.Reset)
		}
		b.WriteString(fmt.Sprintf(" %s %s", h, utils.Reset))
		b.WriteString(utils.Lightgray)
		b.WriteString(utils.VerticalLine)
	}
	b.WriteString("\n")

	b.WriteString(utils.SquareBottomLeft)
	for i, h := range t.headers {
		b.WriteString(repeat(utils.HorizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			b.WriteString(utils.SquareUpHorizontal)
		}
	}
	b.WriteString(utils.SquareBottomRight)
	b.WriteString(utils.Reset)

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
	fmt.Print(fmt.Sprintf("%s%s%s", utils.Lightgray, t.headers[t.Selected], utils.Reset))

	t.Selected++
	if t.Selected >= len(t.headers) {
		t.Selected = 0
	}

	t.moveToStartofWord()
	fmt.Print(fmt.Sprintf("%s%s", utils.Reset, t.headers[t.Selected]))

}

func (t Tabs) moveToStartofWord() {
	col := t.col + 2
	for i := range t.Selected {
		col += StringLen(t.headers[i]) + 3
	}

	utils.MoveCursor(t.row+1, col)
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
