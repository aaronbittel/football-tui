package component

import (
	"fmt"
	"strings"
	term_utils "tui/internal/term-utils"
	"unicode/utf8"
)

type Tabs struct {
	builder  strings.Builder
	headers  []string
	Selected int
	row      int
	col      int
}

func NewTabs(headers ...string) *Tabs {
	return &Tabs{
		builder: strings.Builder{},
		row:     1,
		col:     1,
		headers: headers,
	}
}

func (t Tabs) String() string {
	defer t.builder.Reset()

	t.builder.WriteString(term_utils.Lightgray)
	t.builder.WriteString(term_utils.SquareTopLeft)
	for i, h := range t.headers {
		t.builder.WriteString(repeat(term_utils.HorizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			t.builder.WriteString(term_utils.SquareDownHorizontal)
		}
	}
	t.builder.WriteString(term_utils.SquareTopRight + "\n")

	t.builder.WriteString(term_utils.VerticalLine)
	for i, h := range t.headers {
		if t.Selected == i {
			t.builder.WriteString(term_utils.ResetCode)
		}
		t.builder.WriteString(fmt.Sprintf(" %s %s", h, term_utils.ResetCode))
		t.builder.WriteString(term_utils.Lightgray)
		t.builder.WriteString(term_utils.VerticalLine)
	}
	t.builder.WriteString("\n")

	t.builder.WriteString(term_utils.SquareBottomLeft)
	for i, h := range t.headers {
		t.builder.WriteString(repeat(term_utils.HorizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			t.builder.WriteString(term_utils.SquareUpHorizontal)
		}
	}
	t.builder.WriteString(term_utils.SquareBottomRight)
	t.builder.WriteString(term_utils.ResetCode)

	return t.builder.String()
}

func (t *Tabs) Select(i int) *Tabs {
	if i < 0 || i >= len(t.headers) {
		return t
	}
	t.Selected = i
	return t
}

func (t *Tabs) Next() string {
	defer t.builder.Reset()

	t.builder.WriteString(t.moveToStartofWord())
	t.builder.WriteString(fmt.Sprintf("%s%s%s", term_utils.Lightgray, t.headers[t.Selected], term_utils.ResetCode))

	t.Selected++
	if t.Selected >= len(t.headers) {
		t.Selected = 0
	}

	t.builder.WriteString(t.moveToStartofWord())
	t.builder.WriteString(fmt.Sprintf("%s%s", term_utils.ResetCode, t.headers[t.Selected]))

	return t.builder.String()
}

func (t *Tabs) moveToStartofWord() string {
	defer t.builder.Reset()

	col := t.col + 2
	for i := range t.Selected {
		col += term_utils.StringLen(t.headers[i]) + 3
	}

	t.builder.WriteString(term_utils.MoveCur(t.row+1, col))

	return t.builder.String()
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
