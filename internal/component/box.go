package component

import (
	"fmt"
	"strings"
	term_utils "tui/internal/term-utils"
	"unicode/utf8"
)

type Box struct {
	builder        strings.Builder
	row            int
	col            int
	messages       []string
	title          string
	roundedCornors bool
	centeredText   bool
	padding        Padding
	borderColor    string
}

func NewBox(messages ...string) *Box {
	return &Box{
		builder:  strings.Builder{},
		row:      1,
		col:      1,
		messages: messages,
		padding:  NewPadding(0, 1),
	}
}

func (b *Box) String() string {
	colorize := func(s string) string {
		return fmt.Sprintf("%s%s%s", b.borderColor, s, term_utils.ResetCode)
	}

	multiColorize := func(strs ...string) string {
		out := new(strings.Builder)
		out.WriteString(b.borderColor)
		for _, s := range strs {
			out.WriteString(s)
		}
		out.WriteString(term_utils.ResetCode)
		return out.String()
	}

	var out strings.Builder
	var (
		topLeft     = term_utils.SquareTopLeft
		topRight    = term_utils.SquareTopRight
		bottomLeft  = term_utils.SquareBottomLeft
		bottomRight = term_utils.SquareBottomRight
	)

	if b.roundedCornors {
		topLeft = term_utils.RoundedTopLeft
		topRight = term_utils.RoundedTopRight
		bottomLeft = term_utils.RoundedBottomLeft
		bottomRight = term_utils.RoundedBottomRight
	}

	maxLength := MaxLength(b.messages)
	totalLength := maxLength + b.padding.left + b.padding.right

	titleLen := len(b.title)
	spaceTitle := (totalLength - titleLen) / 2

	out.WriteString(
		multiColorize(
			topLeft,
			strings.Repeat(term_utils.HorizontalLine, spaceTitle),
			b.title,
			strings.Repeat(term_utils.HorizontalLine, totalLength-titleLen-spaceTitle),
			topRight+"\n",
		),
	)

	for range b.padding.top {
		out.WriteString(
			multiColorize(
				term_utils.VerticalLine,
				strings.Repeat(" ", maxLength+b.padding.left+b.padding.right),
				term_utils.VerticalLine+"\n",
			),
		)
	}

	for _, message := range b.messages {
		mLen := term_utils.StringLen(message)
		out.WriteString(colorize(term_utils.VerticalLine))
		out.WriteString(strings.Repeat(" ", b.padding.left))
		if !b.centeredText {
			rightSpace := maxLength - mLen
			out.WriteString(message)
			out.WriteString(strings.Repeat(" ", rightSpace))
		} else {
			leftSpaceLen := (maxLength - mLen) / 2
			out.WriteString(strings.Repeat(" ", leftSpaceLen))
			out.WriteString(message)
			out.WriteString(strings.Repeat(" ", leftSpaceLen))
		}
		out.WriteString(strings.Repeat(" ", b.padding.right))
		out.WriteString(colorize(term_utils.VerticalLine + "\n"))
	}

	for range b.padding.bottom {
		out.WriteString(
			multiColorize(
				term_utils.VerticalLine,
				strings.Repeat(" ", maxLength+b.padding.left+b.padding.right),
				term_utils.VerticalLine+"\n",
			),
		)
	}

	out.WriteString(
		multiColorize(
			bottomLeft,
			strings.Repeat(term_utils.HorizontalLine, maxLength+b.padding.left+b.padding.right),
			bottomRight,
		),
	)

	return out.String()
}

func (b *Box) Lines() []string {
	return strings.Split(b.String(), "\n")
}

func (b Box) Size() (rows, cols int) {
	lines := b.Lines()
	return len(lines), term_utils.StringLen(lines[0])
}

func (b *Box) Update(updated int) string {
	//HACK: improve this
	// Right now this only works for speed box
	return fmt.Sprintf("%s%d", term_utils.MoveCur(b.row+1, b.col+2), updated)
}

func (b *Box) At(row int, optCol ...int) *Box {
	var col int = 1
	if len(optCol) >= 1 {
		col = max(optCol[0], 1)
	}

	b.row = row
	b.col = col
	return b
}

func (b *Box) Centered(width int) *Box {
	maxLength := MaxLength(b.messages)
	totalLength := maxLength + b.padding.left + b.padding.right
	b.col = (width - totalLength) / 2
	return b
}

func (b *Box) WithTitle(title string) *Box {
	b.title = fmt.Sprintf(" %s ", title)
	return b
}

func (b *Box) WithPadding(p ...int) *Box {
	if len(p) == 0 || len(p) > 4 {
		return b
	}

	b.padding = NewPadding(p...)

	return b
}

func (b *Box) WithRoundedCorners() *Box {
	b.roundedCornors = true
	return b
}

func (b *Box) WithCenteredText() *Box {
	b.centeredText = true
	return b
}

func (b *Box) WithColoredBorder(color string) *Box {
	b.borderColor = color
	return b
}

func MaxLength(messages []string) int {
	maxLength := 0
	for _, m := range messages {
		length := utf8.RuneCountInString(term_utils.StripAnsi(m))
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

func (b Box) Pos() (row, col int) {
	return b.row, b.col
}

type Padding struct {
	top    int
	right  int
	bottom int
	left   int
}

func NewPadding(ps ...int) Padding {
	var p Padding

	switch len(ps) {
	case 4:
		p.top = ps[0]
		p.right = ps[1]
		p.bottom = ps[2]
		p.left = ps[3]
	case 3:
		p.top = ps[0]
		p.right = ps[1]
		p.bottom = ps[1]
		p.left = ps[2]
	case 2:
		p.top = ps[0]
		p.right = ps[1]
		p.bottom = ps[0]
		p.left = ps[1]
	case 1:
		p.top = ps[0]
		p.right = ps[0]
		p.bottom = ps[0]
		p.left = ps[0]
	}
	return p
}
