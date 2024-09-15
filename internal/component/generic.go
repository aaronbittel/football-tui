package component

import "strings"

type Generic struct {
	content string
	row     int
	col     int
}

func NewGeneric(content string) *Generic {
	return &Generic{
		content: content,
	}
}

func (g *Generic) Update(newContent string) {
	g.content = newContent
	Clear(g)
	Print(g)
}

func (g Generic) String() string {
	return g.content
}

func (g Generic) Lines() []string {
	return strings.Split(g.content, "\n")
}

func (g Generic) Pos() (row, col int) {
	return g.row, g.col
}

func (g *Generic) At(row, col int) *Generic {
	g.row = row
	g.col = col
	return g
}
