package component

import (
	"fmt"
	"strings"
	utils "tui/internal/term-utils"
	"unicode/utf8"
)

const ()

const sampleTable = `
┌───────┬─────┬─────┬────────┐
│Hallo  │Was  │ Geht│ Hiiiiii│
├───────┼─────┼─────┼────────┤
│  sdfa	│asdf │asdf │  asdf  │
└───────┴─────┴─────┴────────┘
`
const bubbleteaTable = `
┌────────────────────────────┐
│Rank    City      County    │
│────────────────────────────│
│                            │
└────────────────────────────┘
`

type Table struct {
	headers        []Header
	rows           [][]string
	row            int
	col            int
	minLengths     []int
	padding        int
	seperator      bool
	roundedCorners bool
}

func NewTable(headers ...Header) *Table {
	minLen := make([]int, len(headers))

	for i, h := range headers {
		minLen[i] = utf8.RuneCountInString(h.text)
	}

	return &Table{
		headers:    headers,
		minLengths: minLen,
		padding:    1,
	}
}

type Header struct {
	text     string
	centered bool
}

func NewHeader(text string, centered bool) Header {
	return Header{
		text:     text,
		centered: centered,
	}
}

func (t *Table) At(row, col int) *Table {
	t.row = row
	t.col = col
	return t
}

func (t Table) Pos() (height, width int) {
	return t.row, t.col
}

func (t *Table) AddRow(row []string) {
	lenH, lenR := len(t.headers), len(row)

	if lenR > lenH {
		t.rows = append(t.rows, row[:lenH])
		return
	}

	if lenR < lenH {
		for range lenH - lenR {
			row = append(row, "")
		}
	}

	for i, r := range row {
		length := utf8.RuneCountInString(r)
		if t.minLengths[i] < length {
			t.minLengths[i] = length
		}
	}

	t.rows = append(t.rows, row)
}

func (t Table) String() string {
	b := new(strings.Builder)

	b.WriteString(t.createTopLine())
	b.WriteString(t.createHeading())

	if len(t.rows) > 0 {
		b.WriteString(t.createSeperator())
	}

	for i := range t.rows {
		b.WriteString(t.createRow(i))
		if t.seperator {
			if i != len(t.rows)-1 {
				b.WriteString(t.createSeperator())
			}
		}
	}

	b.WriteString(t.createBottomLine())

	return b.String()
}

func (t Table) createRow(idx int) string {
	centerText := func(s string, length int) string {
		l := utf8.RuneCountInString(s)
		space := strings.Repeat(" ", (length-l)/2)

		return fmt.Sprintf("%s%s%s", space, s, strings.Repeat(" ", length-l-len(space)))
	}

	b := new(strings.Builder)

	for i, l := range t.minLengths {
		var (
			item   = t.rows[idx][i]
			length = utf8.RuneCountInString(item)
		)
		b.WriteString(verticalLine)
		b.WriteString(strings.Repeat(" ", t.padding))
		if t.headers[i].centered {
			b.WriteString(centerText(item, l))
		} else {
			b.WriteString(item)
			b.WriteString(strings.Repeat(" ", l-length))
		}
		b.WriteString(strings.Repeat(" ", t.padding))
	}

	b.WriteString(verticalLine)

	return b.String() + "\n"
}

func (t Table) createSeperator() string {
	b := new(strings.Builder)
	b.WriteString(utils.SquareRightVertial)

	for i, l := range t.minLengths {
		b.WriteString(strings.Repeat(horizontalLine, l+2*t.padding))
		if i != len(t.headers)-1 {
			b.WriteString(utils.SquareCross)
			continue
		}
		b.WriteString(utils.SquareLeftVertial)
	}

	return b.String() + "\n"
}

func (t Table) createBottomLine() string {
	b := new(strings.Builder)

	var (
		bottomLeft  = squareBottomLeft
		bottomRight = squareBottomRight
	)

	if t.roundedCorners {
		bottomLeft = roundedBottomLeft
		bottomRight = roundedBottomRight
	}

	b.WriteString(bottomLeft)
	for i, l := range t.minLengths {
		b.WriteString(strings.Repeat(horizontalLine, l+2*t.padding))
		if i != len(t.headers)-1 {
			b.WriteString(utils.SquareUpHorizontal)
			continue
		}
		b.WriteString(bottomRight)
	}

	return b.String()
}

func (t Table) createHeading() string {
	b := new(strings.Builder)

	for i, h := range t.headers {
		b.WriteString(verticalLine)
		b.WriteString(strings.Repeat(" ", t.padding))
		b.WriteString(h.text)
		b.WriteString(strings.Repeat(" ", t.minLengths[i]-utf8.RuneCountInString(h.text)))
		b.WriteString(strings.Repeat(" ", t.padding))
		if i == len(t.headers)-1 {
			b.WriteString(verticalLine)
		}
	}

	return b.String() + "\n"
}

func (t Table) createTopLine() string {
	b := new(strings.Builder)

	var (
		topLeft  = squareTopLeft
		topRight = squareTopRight
	)

	if t.roundedCorners {
		topLeft = roundedTopLeft
		topRight = roundedTopRight
	}

	b.WriteString(topLeft)
	for i := range t.headers {
		b.WriteString(strings.Repeat(horizontalLine, t.minLengths[i]+2*t.padding))
		if i != len(t.headers)-1 {
			b.WriteString(utils.SquareDownHorizontal)
			continue
		}
		b.WriteString(topRight)
	}

	return b.String() + "\n"
}

func (t *Table) WithRoundedCorners() *Table {
	t.roundedCorners = true
	return t
}

func (t *Table) WithSeperator() *Table {
	t.seperator = true
	return t
}
