package component

import (
	"fmt"
	"slices"
	"strings"
	utils "tui/internal/term-utils"
)

type ColumnGraph struct {
	column Column
	row    int
	col    int
}

func NewColumnGraph(column Column) *ColumnGraph {
	return &ColumnGraph{
		column: column,
	}
}

func (c ColumnGraph) String() string {
	var b strings.Builder

	m := slices.Max(c.column.nums)
	spaces := len(fmt.Sprint(m))
	var char string

	for row := m; row >= 1; row-- {
		for i, n := range c.column.nums {
			if n >= row {
				if color, found := c.column.colors[i]; found {
					char = fmt.Sprint(color, utils.FullBlock, utils.Reset)
				} else {
					char = utils.FullBlock
				}
			} else {
				char = " "
			}
			b.WriteString(char + strings.Repeat(" ", spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range c.column.nums {
		s := spaces - len(fmt.Sprint(n)) + 1
		b.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}
	return b.String()
}

func (c *ColumnGraph) Update(column Column) {
	var (
		oldNums   = c.column.nums
		oldC      = c.column.colors
		newNums   = column.nums
		newC      = column.colors
		height, _ = c.Mask()
	)

	for i := 0; i < len(oldNums); i++ {
		if oldNums[i] == newNums[i] && oldC[i] == newC[i] {
			continue
		}

		//TODO: Only replace the minimum number of characters

		// if oldC[i] != newC[i] {
		// 	newColor := newC[i]
		// 	fmt.Print(newColor)
		// 	utils.MoveCursor(c.row+(height-oldNums[i]-1), c.col+i*3)
		// 	for range oldNums[i] {
		// 		fmt.Print(utils.FullBlock)
		// 		utils.MoveCursorDown()
		// 	}
		// 	fmt.Print(utils.Reset)
		// }
		//
		// if oldNums[i] != newNums[i] {
		// 	utils.MoveCursor(c.row+(height-oldNums[i]-1), c.col+i*3)
		// 	fmt.Print("X")
		// }

		utils.MoveCursor(c.row+(height-oldNums[i]-1), c.col+i*3)
		for range oldNums[i] {
			fmt.Print(" ")
			utils.MoveCursorDown()
		}
		fmt.Print("  ")

		utils.MoveCursorLeft()
		utils.MoveCursorLeft()
		fmt.Print(newNums[i])

		utils.MoveCursor(c.row+height-2, c.col+i*3)
		fmt.Print(newC[i])
		for range newNums[i] {
			fmt.Print(utils.FullBlock)
			utils.MoveCursorUp()
		}
		fmt.Print(utils.Reset)
		// for i := range height {
		// 	utils.MoveCursor(c.row+i, c.col)
		// 	fmt.Print("X")
		// }
		// fmt.Print("X")

	}
	c.column = column
}

func (c ColumnGraph) Lines() []string {
	return strings.Split(c.String(), "\n")
}

func (c ColumnGraph) Pos() (row, col int) {
	return c.row, c.col
}

func (c ColumnGraph) Mask() (height, width int) {
	lines := c.Lines()
	return len(lines), StringLen(lines[0])
}

func (c *ColumnGraph) At(row, col int) *ColumnGraph {
	c.row = row
	c.col = col
	return c
}

type Column struct {
	nums   []int
	colors map[int]string
}

func NewColumn(nums []int, colors map[int]string) Column {
	if colors == nil {
		colors = make(map[int]string)
	}
	return Column{
		nums:   nums,
		colors: colors,
	}
}

func (c ColumnGraph) Nums() []int {
	return c.column.nums
}
