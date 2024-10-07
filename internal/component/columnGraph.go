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
	instrCh chan<- string

	row    int
	col    int
	height int
	width  int

	builder strings.Builder

	frames []algorithm.ColumnGraphData
	cursor int

	originalNums []int

	algo *Algorithm

	colParams columnParams
}

func NewColumnGraph(instrCh chan<- string, nums []int) *ColumnGraph {
	m := slices.Max(nums)
	spaces := len(fmt.Sprint(m))
	return &ColumnGraph{
		instrCh:      instrCh,
		frames:       make([]algorithm.ColumnGraphData, 0, 100),
		originalNums: nums,
		builder:      strings.Builder{},
		height:       m + 1,
		width:        len(nums) + (len(nums)-1)*spaces,
		colParams: columnParams{
			maxVal: m,
			spaces: spaces,
		},
	}
}

type columnParams struct {
	maxVal int
	spaces int
}

func (cgf ColumnGraph) Chan() chan<- string {
	return cgf.instrCh
}

func (cgf *ColumnGraph) String() string {
	defer cgf.builder.Reset()

	var char string

	for row := cgf.colParams.maxVal; row >= 1; row-- {
		for _, n := range cgf.originalNums {
			if n >= row {
				char = term_utils.FullBlock
			} else {
				char = " "
			}
			cgf.builder.WriteString(char + strings.Repeat(" ", cgf.colParams.spaces))
		}
		cgf.builder.WriteString("\n")
	}

	for _, n := range cgf.originalNums {
		s := cgf.colParams.spaces - len(fmt.Sprint(n)) + 1
		cgf.builder.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}

	return cgf.builder.String()
}

func (cgf ColumnGraph) Pos() (row, col int) {
	return cgf.row, cgf.col
}

func (cgf ColumnGraph) Nums() []int {
	return slices.Clone(cgf.originalNums)
}

func (cgf *ColumnGraph) Init(algo Algorithm) {
	cgf.frames = make([]algorithm.ColumnGraphData, 0, 300)
	columnCh := make(chan algorithm.ColumnGraphData)

	// append "empty" frame because partialUpdate always compares prev with next
	//HACK: the first one gets skipped somehow, so add empty
	cgf.frames = append(
		cgf.frames,
		algorithm.NewColumnGraphData(cgf.originalNums, map[int]string{}, ""))

	// this is always the first frame of the visualization
	algoName := algo.String()
	cgf.frames = append(
		cgf.frames,
		algorithm.NewColumnGraphData(cgf.originalNums, map[int]string{},
			fmt.Sprintf("Starting %s visualization", algoName)))

	cgf.frames = append(cgf.frames, algo.GetFrames(columnCh, slices.Clone(cgf.originalNums))...)
	cgf.cursor = 0

	sortedNums := slices.Clone(cgf.originalNums)
	slices.Sort(sortedNums)

	// this is always the last frame of the visualization
	cgf.frames = append(
		cgf.frames,
		algorithm.NewColumnGraphData(sortedNums, map[int]string{},
			fmt.Sprintf(term_utils.Colorize(
				fmt.Sprintf("Finished %s visualization", algoName),
				term_utils.BoldGreen))))

	cgf.frames = slices.Clip(cgf.frames)
}

func (cgf ColumnGraph) createLegend(legend []string) *Box {
	legendCol := cgf.col + cgf.width + 6

	return NewBox(cgf.instrCh, legend...).
		WithRoundedCorners().
		WithTitle("Legend").
		At(cgf.row, legendCol)
}

func (cgf ColumnGraph) partialUpdate(prev, next algorithm.ColumnGraphData) string {
	defer cgf.builder.Reset()

	for i := 0; i < len(prev.Nums()); i++ {
		// Only update the columns that have changed
		if prev.Nums()[i] == next.Nums()[i] && prev.Colors()[i] == next.Colors()[i] {
			continue
		}

		cgf.builder.WriteString(cgf.removeColumn(prev.Nums()[i], i))
		cgf.builder.WriteString(cgf.printNewCol(next.Nums()[i], i, next.Colors()[i]))
	}

	cgf.builder.WriteString(cgf.updateDescription(prev.Desc(), next.Desc()))

	return cgf.builder.String()
}

func (cgf *ColumnGraph) Next() {
	if cgf.cursor+1 >= len(cgf.frames) {
		return
	}

	updateInst := cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor+1])
	cgf.cursor++
	cgf.instrCh <- updateInst
}

func (cgf ColumnGraph) Size() (rows, cols int) {
	lines := strings.Split(cgf.String(), "\n")
	// +2 for description
	return len(lines) + 2, len(lines[0])
}

func (cgf *ColumnGraph) Prev() {
	if cgf.cursor-1 < 0 {
		return
	}

	cgf.instrCh <- cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor-1])
	cgf.cursor--
}

func (cgf *ColumnGraph) At(row, col int) {
	cgf.row = row
	cgf.col = col
}

func (cgf ColumnGraph) moveCursorTopColumn(colHeight, colIdx int) string {
	return term_utils.MoveCur(cgf.row+(cgf.height-colHeight-1), cgf.col+colIdx*(cgf.colParams.spaces+1))
}

func (cgf *ColumnGraph) removeColumn(colHeight, colIdx int) string {
	defer cgf.builder.Reset()

	// -1 because the number underneath the columns counts into the cgf.height

	cgf.builder.WriteString(cgf.moveCursorTopColumn(colHeight, colIdx))
	for range colHeight {
		cgf.builder.WriteString(" ")
		cgf.builder.WriteString(term_utils.MoveCurDown())
	}
	cgf.builder.WriteString("  ")

	cgf.builder.WriteString(term_utils.MoveCurLeft())
	cgf.builder.WriteString(term_utils.MoveCurLeft())

	return cgf.builder.String()
}

func (cgf *ColumnGraph) printNewCol(colHeight, colIdx int, color string) string {
	defer cgf.builder.Reset()

	cgf.builder.WriteString(fmt.Sprintf("%d", colHeight))

	// Move cursor to first column segment
	cgf.builder.WriteString(term_utils.MoveCur(cgf.row+cgf.height-2, cgf.col+colIdx*(cgf.colParams.spaces+1)))

	cgf.builder.WriteString(color)
	for range colHeight {
		cgf.builder.WriteString(term_utils.FullBlock)
		cgf.builder.WriteString(term_utils.MoveCurUp())
	}
	cgf.builder.WriteString(term_utils.ResetCode)

	return cgf.builder.String()
}

func (cgf *ColumnGraph) updateDescription(prev, next string) string {
	defer cgf.builder.Reset()

	leftPad := max((cgf.width-term_utils.StringLen(next))/2, 0)
	cgf.builder.WriteString(term_utils.MoveCur(cgf.DescCol(), cgf.col))

	cgf.builder.WriteString(strings.Repeat(" ", leftPad))
	cgf.builder.WriteString(next)

	diff := term_utils.StringLen(prev) - term_utils.StringLen(next)
	if diff > 0 {
		cgf.builder.WriteString(strings.Repeat(" ", diff))
	}

	return cgf.builder.String()
}

func (cgf *ColumnGraph) clearGraph() string {
	defer cgf.builder.Reset()

	for i := range cgf.height {
		cgf.builder.WriteString(term_utils.MoveCur(cgf.row+i, cgf.col))
		//HACK: +1 because if last number has more than 1 digit the width is
		// greater than cgf.width
		cgf.builder.WriteString(fmt.Sprint(strings.Repeat(" ", cgf.width+1)))
	}
	return cgf.builder.String()
}

func (cgf *ColumnGraph) clearDescription() string {
	defer cgf.builder.Reset()

	cgf.builder.WriteString(term_utils.MoveCur(cgf.DescCol(), cgf.col))
	cgf.builder.WriteString(fmt.Sprint(strings.Repeat(" ", cgf.width)))

	return cgf.builder.String()
}

func (cgf *ColumnGraph) Reset() {
	clearGraphInst := cgf.clearGraph()
	clearDescInst := cgf.clearDescription()
	printGraphIdleInst := cgf.String()
	cgf.instrCh <- clearGraphInst + clearDescInst + printGraphIdleInst
}

func (cgf ColumnGraph) DescCol() int {
	return cgf.row + cgf.height + 1
}
