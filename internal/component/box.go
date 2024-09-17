package component

import (
	"fmt"
	"strings"
	utils "tui/internal/term-utils"
	"unicode/utf8"
)

type Box struct {
	width          int
	height         int
	row            int
	col            int
	messages       []string
	title          string
	roundedCornors bool
	centeredText   bool
	padding        Padding
	borderColor    string
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

func NewBox(messages ...string) *Box {
	return &Box{
		messages: messages,
		padding:  NewPadding(0, 1),
	}
}

func (b *Box) Update(messages ...string) *Box {
	b.messages = messages
	return b
}

func (b *Box) String() string {
	colorize := func(s string) string {
		return fmt.Sprintf("%s%s%s", b.borderColor, s, utils.Reset)
	}

	multiColorize := func(strs ...string) string {
		out := new(strings.Builder)
		out.WriteString(b.borderColor)
		for _, s := range strs {
			out.WriteString(s)
		}
		out.WriteString(utils.Reset)
		return out.String()
	}

	var out strings.Builder
	var (
		topLeft     = utils.SquareTopLeft
		topRight    = utils.SquareTopRight
		bottomLeft  = utils.SquareBottomLeft
		bottomRight = utils.SquareBottomRight
	)

	if b.roundedCornors {
		topLeft = utils.RoundedTopLeft
		topRight = utils.RoundedTopRight
		bottomLeft = utils.RoundedBottomLeft
		bottomRight = utils.RoundedBottomRight
	}

	maxLength := MaxLength(b.messages)
	totalLength := maxLength + b.padding.left + b.padding.right

	titleLen := len(b.title)
	spaceTitle := (totalLength - titleLen) / 2

	out.WriteString(
		multiColorize(
			topLeft,
			strings.Repeat(utils.HorizontalLine, spaceTitle),
			b.title,
			strings.Repeat(utils.HorizontalLine, totalLength-titleLen-spaceTitle),
			topRight+"\n",
		),
	)

	for range b.padding.top {
		out.WriteString(
			multiColorize(
				utils.VerticalLine,
				strings.Repeat(" ", maxLength+b.padding.left+b.padding.right),
				utils.VerticalLine+"\n",
			),
		)
	}

	for _, message := range b.messages {
		mLen := StringLen(message)
		out.WriteString(colorize(utils.VerticalLine))
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
		out.WriteString(colorize(utils.VerticalLine + "\n"))
	}

	for range b.padding.bottom {
		out.WriteString(
			multiColorize(
				utils.VerticalLine,
				strings.Repeat(" ", maxLength+b.padding.left+b.padding.right),
				utils.VerticalLine+"\n",
			),
		)
	}

	out.WriteString(
		multiColorize(
			bottomLeft,
			strings.Repeat(utils.HorizontalLine, maxLength+b.padding.left+b.padding.right),
			bottomRight,
		),
	)

	return out.String()
}

func (b *Box) Lines() []string {
	return strings.Split(b.String(), "\n")
}

func (b *Box) At(row, col int) *Box {
	b.row = row
	b.col = col
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

func (b *Box) WithSize(width, height int) *Box {
	if width <= 5 || height <= 5 {
		return b
	}
	b.width = width
	b.height = height
	return b
}

func MaxLength(messages []string) int {
	maxLength := 0
	for _, m := range messages {
		length := utf8.RuneCountInString(StripAnsi(m))
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

func (b Box) Pos() (row, col int) {
	return b.row, b.col
}

func Clear(c Clearer) {
	height, width := Mask(c)
	row, col := c.Pos()
	for i := range height {
		utils.MoveCursor(row+i, col)
		fmt.Print(strings.Repeat(" ", width))
	}
}
