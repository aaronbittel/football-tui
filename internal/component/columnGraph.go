package component

import (
	"fmt"
	"slices"
	"strings"
	utils "tui/internal/term-utils"
)

type ColumnGraph struct {
	column    Column
	row       int
	col       int
	colParams columnParams
}

type columnParams struct {
	maxVal int
	spaces int
}

func NewColumnGraph(column Column) *ColumnGraph {
	m := slices.Max(column.nums)
	spaces := len(fmt.Sprint(m))

	return &ColumnGraph{
		column: column,
		colParams: columnParams{
			maxVal: m,
			spaces: spaces,
		},
	}
}

func (c ColumnGraph) String() string {
	var b strings.Builder

	var char string

	for row := c.colParams.maxVal; row >= 1; row-- {
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
			b.WriteString(char + strings.Repeat(" ", c.colParams.spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range c.column.nums {
		s := c.colParams.spaces - len(fmt.Sprint(n)) + 1
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

		// -1 because the number underneath the columns counts into the height
		utils.MoveCursor(c.row+(height-oldNums[i]-1), c.col+i*(c.colParams.spaces+1))
		for range oldNums[i] {
			fmt.Print(" ")
			utils.MoveCursorDown()
		}
		fmt.Print("  ")

		utils.MoveCursorLeft()
		utils.MoveCursorLeft()
		fmt.Print(newNums[i])

		// Move cursor to first column segment
		utils.MoveCursor(c.row+height-2, c.col+i*(c.colParams.spaces+1))
		fmt.Print(newC[i])
		for range newNums[i] {
			fmt.Print(utils.FullBlock)
			utils.MoveCursorUp()
		}
		fmt.Print(utils.Reset)

		//HACK: Remove from line 30, col 30 to end of line (remove description)
		//TODO: Make text center beneath graph, remove fixed numbers
		utils.ClearLine(30, 30)
		utils.MoveCursor(30, 30)
		fmt.Print(column.desc)

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

func (c ColumnGraph) Nums() []int {
	return c.column.nums
}

type Column struct {
	nums   []int
	colors map[int]string
	desc   string
}

func NewColumn(nums []int, colors map[int]string, desc string) Column {
	if colors == nil {
		colors = make(map[int]string)
	}
	return Column{
		nums:   nums,
		colors: colors,
		desc:   desc,
	}
}

func (c Column) Desc() string {
	return c.desc
}
