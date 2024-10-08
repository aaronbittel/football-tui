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

func (cg ColumnGraph) Chan() chan<- string {
	return cg.instrCh
}

func (cg *ColumnGraph) String() string {
	defer cg.builder.Reset()

	var char string

	for row := cg.colParams.maxVal; row >= 1; row-- {
		for _, n := range cg.originalNums {
			if n >= row {
				char = term_utils.FullBlock
			} else {
				char = " "
			}
			cg.builder.WriteString(char + strings.Repeat(" ", cg.colParams.spaces))
		}
		cg.builder.WriteString("\n")
	}

	for _, n := range cg.originalNums {
		s := cg.colParams.spaces - len(fmt.Sprint(n)) + 1
		cg.builder.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}

	return cg.builder.String()
}

func (cg ColumnGraph) Pos() (row, col int) {
	return cg.row, cg.col
}

func (cg ColumnGraph) Nums() []int {
	return slices.Clone(cg.originalNums)
}

func (cg *ColumnGraph) Init(algo Algorithm) {
	cg.frames = make([]algorithm.ColumnGraphData, 0, 300)
	columnCh := make(chan algorithm.ColumnGraphData)

	// append "empty" frame because partialUpdate always compares prev with next
	//HACK: the first one gets skipped somehow, so add empty
	cg.frames = append(
		cg.frames,
		algorithm.NewColumnGraphData(cg.originalNums, map[int]string{}, ""))

	// this is always the first frame of the visualization
	algoName := algo.String()
	cg.frames = append(
		cg.frames,
		algorithm.NewColumnGraphData(cg.originalNums, map[int]string{},
			fmt.Sprintf("Starting %s visualization", algoName)))

	cg.frames = append(cg.frames, algo.GetFrames(columnCh, slices.Clone(cg.originalNums))...)
	cg.cursor = 0

	sortedNums := slices.Clone(cg.originalNums)
	slices.Sort(sortedNums)

	// this is always the last frame of the visualization
	cg.frames = append(
		cg.frames,
		algorithm.NewColumnGraphData(sortedNums, map[int]string{},
			fmt.Sprintf(term_utils.Colorize(
				fmt.Sprintf("Finished %s visualization", algoName),
				term_utils.BoldGreen))))

	cg.frames = slices.Clip(cg.frames)
}

func (cg ColumnGraph) createLegend(legend []string) *Box {
	legendCol := cg.col + cg.width + 6

	return NewBox(cg.instrCh, legend...).
		WithRoundedCorners().
		WithTitle("Legend").
		At(cg.row, legendCol)
}

func (cg ColumnGraph) partialUpdate(prev, next algorithm.ColumnGraphData) string {
	defer cg.builder.Reset()

	for i := 0; i < len(prev.Nums()); i++ {
		// Only update the columns that have changed
		if prev.Nums()[i] == next.Nums()[i] && prev.Colors()[i] == next.Colors()[i] {
			continue
		}

		cg.builder.WriteString(cg.removeColumn(prev.Nums()[i], i))
		cg.builder.WriteString(cg.printNewCol(next.Nums()[i], i, next.Colors()[i]))
	}

	cg.builder.WriteString(cg.updateDescription(prev.Desc(), next.Desc()))

	return cg.builder.String()
}

func (cg *ColumnGraph) Next() {
	if cg.cursor+1 >= len(cg.frames) {
		return
	}

	updateInst := cg.partialUpdate(cg.frames[cg.cursor], cg.frames[cg.cursor+1])
	cg.cursor++
	cg.instrCh <- updateInst
}

func (cg ColumnGraph) Size() (rows, cols int) {
	lines := strings.Split(cg.String(), "\n")
	// +2 for description
	return len(lines) + 2, len(lines[0])
}

func (cg *ColumnGraph) Prev() {
	if cg.cursor-1 < 0 {
		return
	}

	cg.instrCh <- cg.partialUpdate(cg.frames[cg.cursor], cg.frames[cg.cursor-1])
	cg.cursor--
}

func (cg *ColumnGraph) At(row, col int) {
	cg.row = row
	cg.col = col
}

func (cg ColumnGraph) moveCursorTopColumn(colHeight, colIdx int) string {
	return term_utils.MoveCur(cg.row+(cg.height-colHeight-1), cg.col+colIdx*(cg.colParams.spaces+1))
}

func (cg *ColumnGraph) removeColumn(colHeight, colIdx int) string {
	defer cg.builder.Reset()

	// -1 because the number underneath the columns counts into the cg.height

	cg.builder.WriteString(cg.moveCursorTopColumn(colHeight, colIdx))
	for range colHeight {
		cg.builder.WriteString(" ")
		cg.builder.WriteString(term_utils.MoveCurDown())
	}
	cg.builder.WriteString("  ")

	cg.builder.WriteString(term_utils.MoveCurLeft())
	cg.builder.WriteString(term_utils.MoveCurLeft())

	return cg.builder.String()
}

func (cg *ColumnGraph) printNewCol(colHeight, colIdx int, color string) string {
	defer cg.builder.Reset()

	cg.builder.WriteString(fmt.Sprintf("%d", colHeight))

	// Move cursor to first column segment
	cg.builder.WriteString(term_utils.MoveCur(cg.row+cg.height-2, cg.col+colIdx*(cg.colParams.spaces+1)))

	cg.builder.WriteString(color)
	for range colHeight {
		cg.builder.WriteString(term_utils.FullBlock)
		cg.builder.WriteString(term_utils.MoveCurUp())
	}
	cg.builder.WriteString(term_utils.ResetCode)

	return cg.builder.String()
}

func (cg *ColumnGraph) updateDescription(prev, next string) string {
	defer cg.builder.Reset()

	leftPad := max((cg.width-term_utils.StringLen(next))/2, 0)
	cg.builder.WriteString(term_utils.MoveCur(cg.DescCol(), cg.col))

	cg.builder.WriteString(strings.Repeat(" ", leftPad))
	cg.builder.WriteString(next)

	diff := term_utils.StringLen(prev) - term_utils.StringLen(next)
	if diff > 0 {
		cg.builder.WriteString(strings.Repeat(" ", diff))
	}

	return cg.builder.String()
}

func (cg *ColumnGraph) clearGraph() string {
	defer cg.builder.Reset()

	for i := range cg.height {
		cg.builder.WriteString(term_utils.MoveCur(cg.row+i, cg.col))
		//HACK: +1 because if last number has more than 1 digit the width is
		// greater than cg.width
		cg.builder.WriteString(fmt.Sprint(strings.Repeat(" ", cg.width+1)))
	}
	return cg.builder.String()
}

func (cg *ColumnGraph) clearDescription() string {
	defer cg.builder.Reset()

	cg.builder.WriteString(term_utils.MoveCur(cg.DescCol(), cg.col))
	cg.builder.WriteString(fmt.Sprint(strings.Repeat(" ", cg.width)))

	return cg.builder.String()
}

func (cg *ColumnGraph) Reset() {
	clearGraphInst := cg.clearGraph()
	clearDescInst := cg.clearDescription()
	printGraphIdleInst := cg.String()
	cg.instrCh <- clearGraphInst + clearDescInst + printGraphIdleInst
}

func (cg ColumnGraph) DescCol() int {
	return cg.row + cg.height + 1
}
