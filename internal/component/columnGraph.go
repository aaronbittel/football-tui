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

func (cgf *ColumnGraph) String() string {
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

func (cgf ColumnGraph) Pos() (row, col int) {
	return cgf.row, cgf.col
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

	// append "empty" frame because partialUpdate always compares prev with next
	cgf.frames = append(
		cgf.frames,
		algorithm.NewColumnGraphData(cgf.originalNums, map[int]string{}, ""),
	)

	// this is always the first frame of the visualization
	algoName := algo.Name.String()
	cgf.frames = append(
		cgf.frames,
		algorithm.NewColumnGraphData(cgf.originalNums, map[int]string{},
			fmt.Sprintf("Starting %s visualization", algoName)),
	)

	for col := range columnCh {
		cgf.frames = append(cgf.frames, col)
	}

	sortedNums := slices.Clone(cgf.originalNums)
	slices.Sort(sortedNums)

	// this is always the last frame of the visualization
	cgf.frames = append(
		cgf.frames,
		algorithm.NewColumnGraphData(sortedNums, map[int]string{},
			fmt.Sprintf(term_utils.Colorize(
				fmt.Sprintf("Finished %s visualization", algoName),
				term_utils.BoldGreen,
			))),
	)

	cgf.algo = algo
	cgf.legendBox = cgf.createLegend(cgf.algo.Legend)
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

	updateInst := cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor+1])
	cgf.cursor++
	return updateInst
}

func (cgf ColumnGraph) Size() (rows, cols int) {
	lines := strings.Split(cgf.String(), "\n")
	// +2 for description
	return len(lines) + 2, len(lines[0])
}

func (cgf *ColumnGraph) Prev() string {
	if cgf.cursor-1 < 0 {
		return ""
	}

	updateInst := cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor-1])
	cgf.cursor--
	return updateInst
}

func (cgf *ColumnGraph) At(row, col int) {
	cgf.row = row
	cgf.col = col
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

func (cgf ColumnGraph) ClearGraph() string {
	var b strings.Builder
	for i := range cgf.height {
		b.WriteString(term_utils.MoveCur(cgf.row+i, cgf.col))
		//HACK: +1 because if last number has more than 1 digit the width is
		// greater than cgf.width
		b.WriteString(fmt.Sprint(strings.Repeat(" ", cgf.width+1)))
	}
	return b.String()
}

func (cgf ColumnGraph) ClearDescription() string {
	var b strings.Builder
	b.WriteString(term_utils.MoveCur(cgf.DescCol(), cgf.col))
	b.WriteString(fmt.Sprint(strings.Repeat(" ", cgf.width)))
	return b.String()
}

func (cgf ColumnGraph) Reset() string {
	clearGraphInst := cgf.ClearGraph()
	clearDescInst := cgf.ClearDescription()
	// printGraphInst := cgf.String()
	return clearGraphInst + clearDescInst
}

func (cgf ColumnGraph) DescCol() int {
	return cgf.row + cgf.height + 1
}
