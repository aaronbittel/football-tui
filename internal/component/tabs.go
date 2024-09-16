package component

import (
	"fmt"
	"strings"
	utils "tui/internal/term-utils"
	"unicode/utf8"
)

const (
	lightgray = "\033[38;5;240m"
	white     = "\033[37m"
)

type Tabs struct {
	headers  []string
	Selected int
	row      int
	col      int
}

func NewTabs(headers ...string) *Tabs {
	return &Tabs{
		headers: headers,
	}
}

func (t Tabs) String() string {
	b := new(strings.Builder)

	b.WriteString(lightgray)
	b.WriteString(squareTopLeft)
	for i, h := range t.headers {
		b.WriteString(repeat(horizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			b.WriteString(utils.SquareDownHorizontal)
		}
	}
	b.WriteString(squareTopRight + "\n")

	b.WriteString(verticalLine)
	for i, h := range t.headers {
		if t.Selected == i {
			b.WriteString(utils.Reset)
		}
		b.WriteString(fmt.Sprintf(" %s %s", h, utils.Reset))
		b.WriteString(lightgray)
		b.WriteString(verticalLine)
	}
	b.WriteString("\n")

	b.WriteString(squareBottomLeft)
	for i, h := range t.headers {
		b.WriteString(repeat(horizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			b.WriteString(utils.SquareUpHorizontal)
		}
	}
	b.WriteString(squareBottomRight)
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
	t.Selected++
	if t.Selected >= len(t.headers) {
		t.Selected = 0
	}
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