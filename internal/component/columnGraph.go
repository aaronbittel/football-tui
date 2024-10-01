package component

import (
	"fmt"
	"slices"
	"strings"
	"tui/internal/algorithm"
	term_utils "tui/internal/term-utils"
)

var debug = term_utils.GetDebugFunc()

type ColumnGraph struct {
	frames []algorithm.ColumnGraphData
	cursor int

	originalNums []int

	row    int
	col    int
	height int
	width  int

	algo *Algorithm

	legendBox *Box

	colParams columnParams
}

func NewColumnGraph(nums []int) *ColumnGraph {
	m := slices.Max(nums)
	spaces := len(fmt.Sprint(m))
	return &ColumnGraph{
		frames:       make([]algorithm.ColumnGraphData, 0, 100),
		originalNums: nums,
		height:       m + 1,
		width:        len(nums) + (len(nums)-1)*spaces,
		colParams: columnParams{
			maxVal: m,
			spaces: spaces,
		},
	}
}

func (cgf *ColumnGraph) Idle() string {
	var b strings.Builder

	var char string

	for row := cgf.colParams.maxVal; row >= 1; row-- {
		for _, n := range cgf.originalNums {
			if n >= row {
				char = term_utils.FullBlock
			} else {
				char = " "
			}
			b.WriteString(char + strings.Repeat(" ", cgf.colParams.spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range cgf.originalNums {
		s := cgf.colParams.spaces - len(fmt.Sprint(n)) + 1
		b.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}

	return b.String()
}

func (cgf ColumnGraph) Nums() []int {
	return slices.Clone(cgf.originalNums)
}

func (cgf *ColumnGraph) Init(algo *Algorithm) {
	columnCh := make(chan algorithm.ColumnGraphData)

	go algo.AlgorithmFn(columnCh, slices.Clone(cgf.originalNums))

	// Resetting from previous visualization
	cgf.frames = make([]algorithm.ColumnGraphData, 0, 100)
	cgf.cursor = 0
	for col := range columnCh {
		cgf.frames = append(cgf.frames, col)
	}

	cgf.algo = algo

	cgf.legendBox = cgf.createLegend(cgf.algo.legend)
	Print(cgf.legendBox)
}

func (cgf ColumnGraph) createLegend(legend []string) *Box {
	legendCol := cgf.col + cgf.width + 6

	return NewBox(legend...).
		WithRoundedCorners().
		WithTitle("Legend").
		At(cgf.row, legendCol)
}

func (cgf ColumnGraph) partialUpdate(prev, next algorithm.ColumnGraphData) string {
	var b strings.Builder
	for i := 0; i < len(prev.Nums()); i++ {
		// Only update the columns that have changed
		if prev.Nums()[i] == next.Nums()[i] && prev.Colors()[i] == next.Colors()[i] {
			continue
		}

		b.WriteString(cgf.removeColumn(prev.Nums()[i], i))
		b.WriteString(cgf.printNewCol(next.Nums()[i], i, next.Colors()[i]))
	}
	b.WriteString(cgf.updateDescription(prev.Desc(), next.Desc()))

	return b.String()
}

func (cgf *ColumnGraph) Next() string {
	if cgf.cursor+1 >= len(cgf.frames) {
		return ""
	}

	updateInstructions := cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor+1])
	cgf.cursor++
	return updateInstructions
}

func (cgf *ColumnGraph) Prev() {
	if cgf.cursor-1 < 0 {
		return
	}

	cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor-1])
	cgf.cursor--
}

func (cgf ColumnGraph) Print() {
	var b strings.Builder

	var char string
	frame := cgf.frames[cgf.cursor]

	for row := cgf.colParams.maxVal; row >= 1; row-- {
		for i, n := range frame.Nums() {
			if n >= row {
				if color, found := frame.Colors()[i]; found {
					char = fmt.Sprint(color, term_utils.FullBlock, term_utils.ResetCode)
				} else {
					char = term_utils.FullBlock
				}
			} else {
				char = " "
			}
			b.WriteString(char + strings.Repeat(" ", cgf.colParams.spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range frame.Nums() {
		s := cgf.colParams.spaces - len(fmt.Sprint(n)) + 1
		b.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}

	width := len(cgf.originalNums) + (len(cgf.originalNums)-1)*cgf.colParams.spaces

	b.WriteString("\n")
	leftPad := (width - term_utils.StringLen(frame.Desc())) / 2
	if leftPad > 0 {
		b.WriteString(strings.Repeat(" ", leftPad))
	}
	b.WriteString(frame.Desc())

	for i, line := range strings.Split(b.String(), "\n") {
		term_utils.MoveCursor(cgf.row+i, cgf.col)
		fmt.Print(line)
	}
}

func (cgf *ColumnGraph) At(row, col int) {
	cgf.row = row
	cgf.col = col
}

func (cgf ColumnGraph) PrintIdle() {
	for i, line := range strings.Split(cgf.Idle(), "\n") {
		term_utils.MoveCursor(cgf.row+i, cgf.col)
		fmt.Print(line)
	}
}

func (cgf ColumnGraph) Clear() {
	width := len(cgf.originalNums) + (len(cgf.originalNums)-1)*cgf.colParams.spaces
	for i := range strings.Split(cgf.Idle(), "\n") {
		term_utils.MoveCursor(cgf.row+i, cgf.col)
		fmt.Print(strings.Repeat(" ", width))
	}
}

func (cgf ColumnGraph) moveCursorTopColumn(colHeight, colIdx int) string {
	return term_utils.MoveCur(cgf.row+(cgf.height-colHeight-1), cgf.col+colIdx*(cgf.colParams.spaces+1))
}

func (cgf ColumnGraph) removeColumn(colHeight, colIdx int) string {
	var b strings.Builder

	// -1 because the number underneath the columns counts into the cgf.height

	b.WriteString(cgf.moveCursorTopColumn(colHeight, colIdx))
	for range colHeight {
		b.WriteString(" ")
		b.WriteString(term_utils.MoveCurDown())
	}
	b.WriteString("  ")

	b.WriteString(term_utils.MoveCurLeft())
	b.WriteString(term_utils.MoveCurLeft())

	return b.String()
}

func (cgf ColumnGraph) printNewCol(colHeight, colIdx int, color string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%d", colHeight))

	// Move cursor to first column segment
	b.WriteString(term_utils.MoveCur(cgf.row+cgf.height-2, cgf.col+colIdx*(cgf.colParams.spaces+1)))
	b.WriteString(color)
	for range colHeight {
		b.WriteString(term_utils.FullBlock)
		b.WriteString(term_utils.MoveCurUp())
	}
	b.WriteString(term_utils.ResetCode)

	return b.String()
}

func (cgf ColumnGraph) updateDescription(prev, next string) string {
	var b strings.Builder

	leftPad := max((cgf.width-term_utils.StringLen(next))/2, 0)
	b.WriteString(term_utils.MoveCur(cgf.DescCol(), cgf.col))

	b.WriteString(strings.Repeat(" ", leftPad))
	b.WriteString(next)

	diff := term_utils.StringLen(prev) - term_utils.StringLen(next)
	if diff > 0 {
		b.WriteString(strings.Repeat(" ", diff))
	}

	return b.String()
}

func (cgf ColumnGraph) clearGraph() {
	for i := range cgf.height {
		term_utils.MoveCursor(cgf.row+i, cgf.col)
		//FIX: +1 because if last number has more than 1 digit the width is
		// greater than cgf.width
		fmt.Print(strings.Repeat(" ", cgf.width+1))
	}
}

func (cgf ColumnGraph) clearDescription() {
	term_utils.MoveCursor(cgf.DescCol(), cgf.col)
	fmt.Print(strings.Repeat(" ", cgf.width))
}

func (cgf ColumnGraph) Reset() {
	cgf.clearGraph()
	cgf.clearDescription()
	Print(cgf.legendBox)
	cgf.PrintIdle()
}

func (cgf ColumnGraph) DescCol() int {
	return cgf.row + cgf.height + 1
}
