package main

import (
	"strings"
	"unicode/utf8"
)

const (
	bgRed        = "\033[48;2;255;95;95m"
	bgRedFgWhite = "\033[38;2;255;255;255;48;2;255;95;95m"
)

type List struct {
	items    []string
	selected int
	padding  Padding
}

func NewList(items ...string) *List {
	return &List{
		items:   items,
		padding: NewPadding(0, 1),
	}
}

func (l *List) String() string {
	b := new(strings.Builder)

	maxLen := 0
	for _, item := range l.items {
		length := utf8.RuneCountInString(item)
		if length > maxLen {
			maxLen = length
		}
	}

	for i, item := range l.items {
		if i == l.selected {
			b.WriteString(bgRedFgWhite)
		}
		b.WriteString(strings.Repeat(" ", l.padding.left))
		b.WriteString(item)
		b.WriteString(strings.Repeat(" ", maxLen-utf8.RuneCountInString(item)))
		b.WriteString(strings.Repeat(" ", l.padding.right))
		if i == l.selected {
			b.WriteString(reset)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (l *List) Select(i int) *List {
	l.selected = i
	return l
}
