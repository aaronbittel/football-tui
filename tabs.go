package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	lightgray = "\033[38;5;240m"
	white     = "\033[37m"
)

type Tabs struct {
	headers  []string
	selected int
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
			b.WriteString(squareDownHorizontal)
		}
	}
	b.WriteString(squareTopRight + "\n")

	b.WriteString(verticalLine)
	for i, h := range t.headers {
		if t.selected == i {
			b.WriteString(reset)
		}
		b.WriteString(fmt.Sprintf(" %s %s", h, reset))
		b.WriteString(lightgray)
		b.WriteString(verticalLine)
	}
	b.WriteString("\n")

	b.WriteString(squareBottomLeft)
	for i, h := range t.headers {
		b.WriteString(repeat(horizontalLine, utf8.RuneCountInString(h)+2))
		if i != len(t.headers)-1 {
			b.WriteString(squareUpHorizontal)
		}
	}
	b.WriteString(squareBottomRight)
	b.WriteString(reset)

	return b.String()
}

func (t *Tabs) SetTab(i int) *Tabs {
	if i < 0 || i >= len(t.headers) {
		return t
	}
	t.selected = i
	return t
}

func repeat(s string, i int) string {
	return strings.Repeat(s, i)
}
