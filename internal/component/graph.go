package component

import (
	"fmt"
	"slices"
	"strings"
)

type Graph struct {
	nums []int
	row  int
	col  int
}

func NewGraph(nums []int) *Graph {
	return &Graph{
		nums: nums,
	}
}

func (g Graph) String() string {
	var b strings.Builder

	m := slices.Max(g.nums)
	spaces := len(fmt.Sprint(m))
	var char string

	for i := m; i >= 1; i-- {
		for _, n := range g.nums {
			if n >= i {
				char = fullBlock
			} else {
				char = " "
			}
			b.WriteString(char + strings.Repeat(" ", spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range g.nums {
		s := spaces - len(fmt.Sprint(n)) + 1
		b.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}
	return b.String()
}

func (g *Graph) At(row, col int) *Graph {
	g.row = row
	g.col = col
	return g
}

func (g Graph) Pos() (row, col int) {
	return g.row, g.col
}

func (g Graph) Lines() []string {
	return strings.Split(g.String(), "\n")
}

func (g Graph) Mask() (height, width int) {
	lines := strings.Split(g.String(), "\n")
	return len(lines), StringLen(lines[0])
}
