package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	squareTopLeft     = "┌"
	squareTopRight    = "┐"
	squareBottomLeft  = "└"
	squareBottomRight = "┘"

	roundedTopLeft     = "╭"
	roundedTopRight    = "╮"
	roundedBottomLeft  = "╰"
	roundedBottomRight = "╯"

	horizontalLine = "─"
	verticalLine   = "│"
)

type Box struct {
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

func (b *Box) String() string {
	colorize := func(s string) string {
		return fmt.Sprintf("%s%s%s", b.borderColor, s, reset)
	}

	multiColorize := func(strs ...string) string {
		out := new(strings.Builder)
		out.WriteString(b.borderColor)
		for _, s := range strs {
			out.WriteString(s)
		}
		out.WriteString(reset)
		return out.String()
	}

	var out strings.Builder
	var (
		topLeft     = squareTopLeft
		topRight    = squareTopRight
		bottomLeft  = squareBottomLeft
		bottomRight = squareBottomRight
	)

	if b.roundedCornors {
		topLeft = roundedTopLeft
		topRight = roundedTopRight
		bottomLeft = roundedBottomLeft
		bottomRight = roundedBottomRight
	}

	maxLength := MaxLength(b.messages)
	totalLength := maxLength + b.padding.left + b.padding.right

	titleLen := len(b.title)
	spaceTitle := (totalLength - titleLen) / 2

	out.WriteString(
		multiColorize(
			topLeft,
			strings.Repeat(horizontalLine, spaceTitle),
			b.title,
			strings.Repeat(horizontalLine, totalLength-titleLen-spaceTitle),
			topRight+"\n",
		),
	)

	for range b.padding.top {
		out.WriteString(
			multiColorize(
				verticalLine,
				strings.Repeat(" ", maxLength+b.padding.left+b.padding.right),
				verticalLine+"\n",
			),
		)
	}

	for _, message := range b.messages {
		mLen := utf8.RuneCountInString(StripAnsi(message))
		out.WriteString(colorize(verticalLine))
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
		out.WriteString(colorize(verticalLine + "\n"))
	}

	for range b.padding.bottom {
		out.WriteString(
			multiColorize(
				verticalLine,
				strings.Repeat(" ", maxLength+b.padding.left+b.padding.right),
				verticalLine+"\n",
			),
		)
	}

	out.WriteString(
		multiColorize(
			bottomLeft,
			strings.Repeat(horizontalLine, maxLength+b.padding.left+b.padding.right),
			bottomRight,
		),
	)

	return out.String()
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
		length := utf8.RuneCountInString(StripAnsi(m))
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

// REFERENCE: https://github.com/acarl005/stripansi/blob/master/stripansi.go
func StripAnsi(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	re := regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}
