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
	Length   int
}

func NewTabs(headers ...string) *Tabs {
	tabs := &Tabs{
		builder: strings.Builder{},
		row:     1,
		col:     1,
		headers: headers,
	}

	tabs.Length = term_utils.StringLen(strings.Split(tabs.String(), "\n")[0])
	return tabs
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
			b.WriteString(term_utils.ResetCode)
		}
		b.WriteString(fmt.Sprintf(" %s %s", h, term_utils.ResetCode))
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
	b.WriteString(term_utils.ResetCode)

	return b.String()
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

func (t *Tabs) Centered(width int) *Tabs {
	t.col = (width - t.Length) / 2
	return t
}

func (t *Tabs) At(row int, optCol ...int) *Tabs {
	var col int = 1
	if len(optCol) > 1 {
		col = max(optCol[0], 1)
	}

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
